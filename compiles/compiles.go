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

package compiles

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/loader"
)

// Run runs the "compiles" checks on the provided packages. The provided packages should be specified as relative paths
// from the current working directory. Dot expansion ("...") is not supported.
func Run(pkgs []string, w io.Writer) error {
	cfg := loader.Config{}

	rest, err := cfg.FromArgs(pkgs, true)
	if err != nil {
		return errors.Wrapf(err, "failed to parse packages")
	}
	if len(rest) > 0 {
		return errors.Errorf("failed to parse arguments as packages: %v", rest)
	}

	cfg.TypeChecker.Error = func(e error) {
		fmt.Fprintln(w, e.Error())
	}

	if _, err := cfg.Load(); err != nil {
		// return blank error if any errors were encountered during load. Load function prints errors to writer
		// in proper format as they are encountered so no need to create any other output.
		return fmt.Errorf("")
	}
	return nil
}
