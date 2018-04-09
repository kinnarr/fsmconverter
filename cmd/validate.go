// Copyright 2018 Franz Schmidt
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/felixangell/toml"

	"github.com/kinnarr/fsmconverter/config"
	"github.com/kinnarr/fsmconverter/validation"
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validates the fsm config",
	Run: func(cmd *cobra.Command, args []string) {
		searchDir := "fsm"

		fileList := []string{}
		_ = filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".toml") {
				fileList = append(fileList, path)
			}
			return nil
		})

		var tomlBuffer bytes.Buffer
		for _, file := range fileList {
			if readedBytes, err := ioutil.ReadFile(file); err == nil {
				tomlBuffer.Write(readedBytes)
			}
		}

		if _, err := toml.Decode(tomlBuffer.String(), &config.MainConfig); err != nil {
			fmt.Println(err)
			return
		}

		validation.ValidateStates()
		validation.ValidateOutputs()
	},
}
