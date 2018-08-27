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

package generation

import (
	"fmt"
	"math/bits"
	"os"
	"path/filepath"
	"text/template"

	"github.com/kinnarr/fsmconverter/config"
	"go.uber.org/zap"
)

func GenerateVerilog(outputDir string) {
	config.AndString = "&&"
	config.OrString = "||"

	funcMap := template.FuncMap{
		"conditionToString": conditionToString,
		"enumerateKeys": func(mymap map[string]config.State) []string {
			keys := make([]string, 0, len(mymap))
			for k := range mymap {
				keys = append(keys, k)
			}
			return keys
		},
		"convertBinary": func(v int, len int) string {
			return fmt.Sprintf("%d'b%b", len, v)
		},
		"getBinarySize": func(v int) int {
			return bits.Len(uint(v))
		},
		"minus": func(a, b int) int {
			return a - b
		},
	}

	absPath, err := filepath.Abs(outputDir)
	if err != nil {
		zap.S().Fatal("Create fsm path for output", zap.Error(err))
	}
	err = os.MkdirAll(absPath, os.ModePerm)
	if err != nil {
		zap.S().Fatal("Create fsm output directory", zap.Error(err))
	}

	cuPath := filepath.Join(absPath, "cu.v")
	fsmPath := filepath.Join(absPath, "fsm.v")

	cuTempl, err := template.New("cu.tpl").Funcs(funcMap).ParseFiles("tmpl/cu.tpl")
	if err != nil {
		zap.S().Fatal("Create control unit template", zap.Error(err))
	}
	cuFile, err := os.Create(cuPath)
	if err != nil {
		zap.S().Fatal("Create control unit file", zap.Error(err))
	}
	err = cuTempl.Execute(cuFile, config.MainConfig)
	if err != nil {
		zap.S().Fatal("Render control unit template to file", zap.Error(err))
	}

	fsmTempl, err := template.New("fsm.tpl").Funcs(funcMap).ParseFiles("tmpl/fsm.tpl")
	if err != nil {
		zap.S().Fatal("Create fsm template", zap.Error(err))
	}
	fsmFile, err := os.Create(fsmPath)
	if err != nil {
		zap.S().Fatal("Create fsm file", zap.Error(err))
	}
	err = fsmTempl.Execute(fsmFile, config.MainConfig)
	if err != nil {
		zap.S().Fatal("Render fsm template to file", zap.Error(err))
	}
}
