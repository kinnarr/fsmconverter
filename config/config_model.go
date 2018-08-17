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

package config

var MainConfig FsmCreatorConfig
var IgnoreUnknownStates bool
var Optimize bool
var AndString string = "and"
var OrString string = "or"

type FsmCreatorConfig struct {
	States   map[string]State `toml:"state"`
	Inputs   map[string]int
	Outputs  map[string]int
	Defaults defaults
}

type defaults struct {
	Outputs map[string]int
	State   string
}

type State struct {
	Successors       map[string]RootCondition `toml:"next"`
	DefaultSuccessor map[string]interface{}   `toml:"else"`
	Outputs          map[string]int           `toml:"output"`
	Preserve         bool
}

type RootCondition struct {
	And *Condition
	Or  *Condition
}

type Condition struct {
	And        *Condition
	Or         *Condition
	Conditions []map[string]int `toml:"condition"`
}
