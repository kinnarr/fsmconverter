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
	"strings"

	"github.com/kinnarr/fsmconverter/config"
)

func RootConditionToString(rc config.Condition, nextName string) string {
	var conditionBuffer bytes.Buffer

	if len(rc.Subconditions) == 0 {
		conditionBuffer.WriteString(fmt.Sprintf("\talways: next <= %s\n", nextName))
	} else {
		conditionBuffer.WriteString("\tif ")
		if rc.Type == config.ConditionType_And {
			conditionBuffer.WriteString(conditionAndToString(rc))
		}
		if rc.Type == config.ConditionType_Or {
			conditionBuffer.WriteString(conditionOrToString(rc))
		}
		conditionBuffer.WriteString(fmt.Sprintf(": next <= %s\n", nextName))
	}
	return conditionBuffer.String()
}

func ConditionAutoToString(c config.Condition) string {
	if len(c.Subconditions) == 0 {
		return "always"
	} else {
		if c.Type == config.ConditionType_And {
			return conditionAndToString(c)
		}
		if c.Type == config.ConditionType_Or {
			return conditionOrToString(c)
		}
	}
	return "never"
}

func conditionAndToString(c config.Condition) string {
	return conditionToString(c, config.AndString)
}

func conditionOrToString(c config.Condition) string {
	return conditionToString(c, config.OrString)
}

func conditionToString(c config.Condition, logicalOp string) string {
	return conditionToStringOptBinary(c, logicalOp, true)
}

func conditionToStringOptBinary(c config.Condition, logicalOp string, printBinary bool) string {
	conditionStrings := make([]string, 0)
	for _, condition := range c.Conditions {
		for conditionName, conditionValue := range condition {
			if _, ok := config.MainConfig.Inputs[conditionName]; ok {
				if printBinary {
					conditionStrings = append(conditionStrings, fmt.Sprintf("%s == %d'b%b", conditionName, config.MainConfig.Inputs[conditionName], conditionValue))
				} else {
					conditionStrings = append(conditionStrings, fmt.Sprintf("%s == %d", conditionName, conditionValue))
				}
			}
		}
	}
	for _, subcondition := range c.Subconditions {
		if subcondition.Type == config.ConditionType_And {
			conditionStrings = append(conditionStrings, fmt.Sprintf("(%s)", conditionAndToString(subcondition)))
		}
		if subcondition.Type == config.ConditionType_Or {
			conditionStrings = append(conditionStrings, fmt.Sprintf("(%s)", conditionOrToString(subcondition)))
		}
	}
	returnString := strings.Join(conditionStrings, fmt.Sprintf(" %s ", logicalOp))
	return returnString
}
