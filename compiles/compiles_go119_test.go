// Copyright 2023 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !go1.20
// +build !go1.20

package compiles_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/go-compiles/compiles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompilesPassCases(t *testing.T) {
	// set GOFLAGS to blank for test so that it does not use vendor mode
	origValue := os.Getenv("GOFLAGS")
	_ = os.Setenv("GOFLAGS", "")
	defer func() {
		_ = os.Setenv("GOFLAGS", origValue)
	}()

	for i, tc := range []struct {
		files []gofiles.GoFileSpec
	}{
		{
			files: []gofiles.GoFileSpec{
				{
					RelPath: "go.mod",
					Src: `module github.com/foo
go 1.16

require github.com/inner v1.0.0

replace github.com/inner => ./inner
`,
				},
				{
					RelPath: "foo.go",
					Src: `package foo
import "github.com/inner"
func Foo() {
	inner.Inner()
}`,
				},
				{
					RelPath: "foo_test.go",
					Src: `package foo_test
import "testing"
import "github.com/inner"
func TestFoo(t *testing.T) {
	inner.Inner()
}`,
				},
				{
					RelPath: "inner/go.mod",
					Src:     `module github.com/inner`,
				},
				{
					RelPath: "inner/inner.go",
					Src: `package inner
func Inner() {}`,
				},
				{
					RelPath: "vendor/modules.txt",
					Src: `# github.com/inner v1.0.0 => ./inner
## explicit
github.com/inner
# github.com/inner => ./inner
`,
				},
				{
					RelPath: "vendor/github.com/inner/go.mod",
					Src:     `module github.com/inner`,
				},
				{
					RelPath: "vendor/github.com/inner/inner.go",
					Src: `package inner
func Inner() {}`,
				},
			},
		},
	} {
		projectDir := t.TempDir()

		buf := bytes.Buffer{}
		out, err := gofiles.Write(projectDir, tc.files)
		_ = out
		require.NoError(t, err)

		err = compiles.RunInDir([]string{"."}, projectDir, &buf)
		require.NoError(t, err, "Case %d: %v", i, buf.String())
	}
}

func TestCompilesErrorCases(t *testing.T) {
	for i, tc := range []struct {
		files     []gofiles.GoFileSpec
		inputPkgs []string
		want      func(baseDir string) string
	}{
		{
			[]gofiles.GoFileSpec{
				{
					RelPath: "go.mod",
					Src:     "module github.com/go-compiles-tester",
				},
				{
					RelPath: "foo/foo.go",
					Src: `package foo
func Foo() {
	return "Foo"
}`,
				},
				{
					RelPath: "bar/bar.go",
					Src: `package bar
import "fmt"`,
				},
			},
			[]string{
				"foo",
				"bar",
			},
			func(baseDir string) string {
				lines := []string{
					baseDir + "/foo/foo.go" + ":3:9: too many return values\n\thave (string)\n\twant ()",
					baseDir + "/bar/bar.go" + `:2:8: "fmt" imported but not used`,
					"",
				}
				return strings.Join(lines, "\n")
			},
		},
		{
			[]gofiles.GoFileSpec{
				{
					RelPath: "go.mod",
					Src:     "module github.com/go-compiles-tester",
				},
				{
					RelPath: "foo/foo.go",
					Src: `package foo
func Foo() string {
	return "Foo"
}`,
				},
				{
					RelPath: "foo/foo_test.go",
					Src: `package foo_test
import (
	"testing"
	"github.com/go-compiles-tester/foo"
)
func TestFoo(t *testing.T) {
	bar := foo.Foo()
}`,
				},
			},
			[]string{
				"foo",
			},
			func(baseDir string) string {
				lines := []string{
					baseDir + "/foo/foo_test.go" + `:7:2: bar declared but not used`,
					"",
				}
				return strings.Join(lines, "\n")
			},
		},
		{
			[]gofiles.GoFileSpec{
				{
					RelPath: "go.mod",
					Src:     "module github.com/go-compiles-tester",
				},
				{
					RelPath: "foo/foo.go",
					Src: `package foo
func Foo() string {
	return "Foo"
}`,
				},
				{
					RelPath: "foo/foo_test.go",
					Src: `package foo_test
import (
	"testing"
	"github.com/go-compiles-tester/foo"
)
func TestFoo(t *testing.T) {
	bar := foo.Foo()
}`,
				},
			},
			[]string{
				"foo",
			},
			func(baseDir string) string {
				lines := []string{
					baseDir + "/foo/foo_test.go" + `:7:2: bar declared but not used`,
					"",
				}
				return strings.Join(lines, "\n")
			},
		},
		{
			[]gofiles.GoFileSpec{
				{
					RelPath: "go.mod",
					Src:     "module github.com/go-compiles-tester",
				},
				{
					RelPath: "foo/foo.go",
					Src:     `package foo`,
				},
				{
					RelPath: "foo/bar.go",
					Src:     `package bar`,
				},
			},
			[]string{
				"foo",
			},
			func(baseDir string) string {
				lines := []string{
					"-: found packages bar (bar.go) and foo (foo.go) in " + baseDir + "/foo",
					baseDir + "/foo/foo.go" + `:1:1: package foo; expected bar`,
					"",
				}
				return strings.Join(lines, "\n")
			},
		},
	} {
		projectDir := t.TempDir()

		buf := bytes.Buffer{}
		_, err := gofiles.Write(projectDir, tc.files)
		require.NoError(t, err)

		var pkgs []string
		for _, inputPkg := range tc.inputPkgs {
			pkgs = append(pkgs, "./"+inputPkg)
		}
		err = compiles.RunInDir(pkgs, projectDir, &buf)
		require.Error(t, err, fmt.Sprintf("Case %d", i))
		assert.Equal(t, tc.want(projectDir), buf.String(), "Case %d", i)
	}
}
