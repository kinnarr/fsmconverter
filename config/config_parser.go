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

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/felixangell/toml"
	"github.com/imdario/mergo"
	"go.uber.org/zap"
)

var FsmConfigRootDir string
var mainInputConfig fsmConverterInputConfig

func ParseConfig() bool {
	returnValue := true

	_ = filepath.Walk(FsmConfigRootDir, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".toml") {
			var currentFileConfig fsmConverterInputConfig
			if _, err := toml.DecodeFile(path, &currentFileConfig); err != nil {
				zap.S().Errorf("Error parsing toml file %s: %s", path, err)
				returnValue = false
			}
			mergo.Merge(&mainInputConfig, currentFileConfig)
		}
		return nil
	})

	MainConfig.Defaults = mainInputConfig.Defaults
	MainConfig.Inputs = mainInputConfig.Inputs
	MainConfig.Outputs = mainInputConfig.Outputs
	MainConfig.States = make(map[string]State)

	for stateName, inputState := range mainInputConfig.States {
		var state State
		state.DefaultSuccessor = inputState.DefaultSuccessor
		state.Outputs = inputState.Outputs
		state.Preserve = inputState.Preserve
		state.Successors = make(map[string]Condition)
		for nextName, inputCondition := range inputState.Successors {
			var condition Condition
			if inputCondition.And != nil && inputCondition.Or != nil {
				zap.S().Errorf("Root condition can't contain 'and' and 'or' part: %s", nextName)
				returnValue = false
			}
			if inputCondition.And != nil {
				condition.Type = ConditionType_And
				condition.Subconditions = make([]Condition, 1)
				condition.Subconditions[0] = inputConditionToCondition(*inputCondition.And)
			}
			if inputCondition.Or != nil {
				condition.Type = ConditionType_Or
				condition.Subconditions = make([]Condition, 1)
				condition.Subconditions[0] = inputConditionToCondition(*inputCondition.Or)
			}
			state.Successors[nextName] = condition
		}
		MainConfig.States[stateName] = state
	}

	zap.S().Infof("Length of internal config states %d", len(MainConfig.States))

	return returnValue
}

func inputConditionToCondition(inputCondition inputCondition) Condition {
	var condition Condition
	if inputCondition.And != nil && inputCondition.Or != nil {
		zap.S().Errorf("Condition can't contain 'and' and 'or' part")
	}
	if inputCondition.And != nil {
		condition.Type = ConditionType_And
		condition.Conditions = inputCondition.Conditions
		condition.Subconditions = make([]Condition, 1)
		condition.Subconditions[0] = inputConditionToCondition(*inputCondition.And)
	}
	if inputCondition.Or != nil {
		condition.Type = ConditionType_Or
		condition.Conditions = inputCondition.Conditions
		condition.Subconditions = make([]Condition, 1)
		condition.Subconditions[0] = inputConditionToCondition(*inputCondition.Or)
	}
	condition.Conditions = inputCondition.Conditions
	return condition
}
