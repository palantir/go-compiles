// Copyright 2016 Palantir Technologies, Inc.
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

package compiles_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/gofiles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/go-compiles/compiles"
)

func TestCompilesPassCases(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir(".", "")
	require.NoError(t, err)
	defer cleanup()

	for i, tc := range []struct {
		files []gofiles.GoFileSpec
	}{
		{
			files: []gofiles.GoFileSpec{
				{
					RelPath: "foo/foo.go",
					Src: `package foo
import "github.com/inner"
func Foo() {
	inner.Inner()
}`,
				},
				{
					RelPath: "foo/foo_test.go",
					Src: `package foo_test
import "testing"
import "github.com/inner"
func TestFoo(t *testing.T) {
	inner.Inner()
}`,
				},
				{
					RelPath: "foo/vendor/github.com/inner/inner.go",
					Src: `package inner
func Inner() {}`,
				},
			},
		},
	} {
		projectDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		buf := bytes.Buffer{}
		_, err = gofiles.Write(projectDir, tc.files)
		require.NoError(t, err)

		pkgPath := "./" + path.Join(projectDir, "foo")
		err = compiles.Run([]string{pkgPath}, &buf)
		require.NoError(t, err, "Case %d: %v", i, buf.String())
	}
}

func TestCompilesErrorCases(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir(".", "")
	require.NoError(t, err)
	defer cleanup()

	wd, err := os.Getwd()
	require.NoError(t, err)

	for i, tc := range []struct {
		files     []gofiles.GoFileSpec
		inputPkgs []string
		want      func(baseDir string) string
	}{
		{
			[]gofiles.GoFileSpec{
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
					baseDir + "/foo/foo.go" + `:3:9: no result values expected`,
					baseDir + "/bar/bar.go" + `:2:8: "fmt" imported but not used`,
					baseDir + "/bar/bar.go" + `:2:8: "fmt" imported but not used`,
					"",
				}
				return strings.Join(lines, "\n")
			},
		},
		{
			[]gofiles.GoFileSpec{
				{
					RelPath: "foo/foo.go",
					Src: `package foo
func Foo() string {
	return "Foo"
}`,
				},
				{
					RelPath: "foo/foo_test.go",
					Src: `package foo
import (
	"testing"
	"{{index . "foo/foo.go"}}"
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
	"{{index . "foo/foo.go"}}"
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
					"found packages bar (bar.go) and foo (foo.go) in " + baseDir + "/foo",
					"",
				}
				return strings.Join(lines, "\n")
			},
		},
	} {
		projectDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		buf := bytes.Buffer{}
		_, err = gofiles.Write(projectDir, tc.files)
		require.NoError(t, err)

		var pkgs []string
		for _, inputPkg := range tc.inputPkgs {
			pkgs = append(pkgs, "./"+path.Join(projectDir, inputPkg))
		}
		err = compiles.Run(pkgs, &buf)
		require.Error(t, err, fmt.Sprintf("Case %d", i))
		assert.Equal(t, tc.want(path.Join(wd, projectDir)), buf.String(), "Case %d", i)
	}
}
