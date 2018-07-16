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

package validation

import (
	"fmt"

	"github.com/kinnarr/fsmconverter/config"
	"go.uber.org/zap"
)

func ValidateDefaults() bool {
	returnValue := true
	if config.MainConfig.Defaults.State == "" {
		zap.S().Errorf("Default state not defined!")
		returnValue = false
	} else {
		if _, ok := config.MainConfig.States[config.MainConfig.Defaults.State]; !ok && !config.IgnoreUnknownStates {
			zap.S().Errorf("Unknown default state %s", config.MainConfig.Defaults.State)
			returnValue = false
		}
	}
	for defaultName := range config.MainConfig.Defaults.Outputs {
		if _, ok := config.MainConfig.Outputs[defaultName]; !ok {
			zap.S().Errorf("Default for unknown output %s defined!", defaultName)
			returnValue = false
		}
	}
	return returnValue
}

func ValidateStates() bool {
	returnValue := true
	for stateName, state := range config.MainConfig.States {
		for nextName, next := range state.Successors {
			if _, ok := config.MainConfig.States[nextName]; !ok && !config.IgnoreUnknownStates {
				zap.S().Errorf("Unknown state %s for state %s", nextName, stateName)
				returnValue = false
			} else {
				returnValue = validateRootCondition(next, nextName) && returnValue
			}
		}
		if len(state.DefaultSuccessor) > 1 {
			zap.S().Errorf("Only one else state allowed! State %s", stateName)
			returnValue = false
		} else {
			for elseName := range state.DefaultSuccessor {
				if _, ok := config.MainConfig.States[elseName]; !ok && !config.IgnoreUnknownStates {
					zap.S().Errorf("Unknown else state %s for state %s", elseName, stateName)
					returnValue = false
				}
			}
		}
		for outputName, outputValue := range state.Outputs {
			if outputSize, ok := config.MainConfig.Outputs[outputName]; !ok {
				zap.S().Errorf("Could not find output %s from state '%s'", outputName, stateName)
				returnValue = false
			} else {
				if outputSize < len(fmt.Sprintf("%b", outputValue)) {
					zap.S().Errorf("Value for output %s from state '%s' is too large", outputName, stateName)
					returnValue = false
				}
			}
		}
	}
	return returnValue
}

func validateRootCondition(rc config.RootCondition, nextName string) bool {
	if rc.And != nil && rc.Or != nil {
		zap.S().Errorf("Root condition can't contain 'and' and 'or' part: %s", nextName)
		return false
	}
	if rc.And != nil {
		return validateCondition(*rc.And, nextName)
	}
	if rc.Or != nil {
		return validateCondition(*rc.Or, nextName)
	}
	return true
}

func validateCondition(c config.Condition, nextName string) bool {
	returnValue := true
	for _, condition := range c.Conditions {
		for conditionName, conditionValue := range condition {
			if inputSize, ok := config.MainConfig.Inputs[conditionName]; !ok {
				zap.S().Errorf("Could not find input %s from condition for next state '%s'", conditionName, nextName)
				returnValue = false
			} else {
				if inputSize < len(fmt.Sprintf("%b", conditionValue)) {
					zap.S().Errorf("Value for input %s from condition for next state '%s' is too large", conditionName, nextName)
					returnValue = false
				}
			}
		}
		if c.And != nil {
			returnValue = validateCondition(*c.And, nextName) && returnValue
		}
		if c.Or != nil {
			returnValue = validateCondition(*c.Or, nextName) && returnValue
		}
		if len(c.Conditions) == 0 && c.And == nil && c.Or == nil {
			zap.S().Errorf("No conditions found for state! Maybe you forgot an .condition")
			returnValue = false
		}
	}
	return returnValue
}
