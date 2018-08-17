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
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"text/template"

	"github.com/kinnarr/fsmconverter/config"
)

func GenerateVerilog() {
	var conditionBuffer bytes.Buffer
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

	templ, err := template.New("fsm.tpl").Funcs(funcMap).ParseFiles("tmpl/fsm.tpl")
	if err != nil {
		panic(err)
	}
	err = templ.Execute(os.Stdout, config.MainConfig)
	if err != nil {
		panic(err)
	}

	fmt.Print(conditionBuffer.String())
}
