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
				returnValue = validateCondition(next, nextName) && returnValue
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
					zap.S().Errorf("Value %b for output %s from state '%s' is too large (size: %d)", outputValue, outputName, stateName, outputSize)
					returnValue = false
				}
			}
		}
	}
	return returnValue
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
	}
	for _, subcondition := range c.Subconditions {
		returnValue = validateCondition(subcondition, nextName) && returnValue
	}
	if len(c.Conditions) == 0 && len(c.Subconditions) == 0 {
		zap.S().Debugf("No condition found! Assume always: %s", nextName)
	}
	return returnValue
}
