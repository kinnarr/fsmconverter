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
	"os"
	"text/template"

	"github.com/kinnarr/fsmconverter/config"
)

func GenerateDot() {
	var conditionBuffer bytes.Buffer

	funcMap := template.FuncMap{
		"conditionToString": conditionToString,
	}

	templ, err := template.New("dot.tpl").Funcs(funcMap).ParseFiles("tmpl/dot.tpl")
	if err != nil {
		panic(err)
	}
	err = templ.Execute(os.Stdout, config.MainConfig)
	if err != nil {
		panic(err)
	}

	fmt.Print(conditionBuffer.String())
}
