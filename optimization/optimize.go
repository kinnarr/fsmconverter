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

package optimization

import (
	"github.com/kinnarr/fsmconverter/config"
	"github.com/kinnarr/fsmconverter/generation"
	"go.uber.org/zap"
)

var deletedStates map[string]interface{}
var notToDeleteStates map[string]interface{}

func OptimizeConfig() {
	OptimizedConfig.NoModification = false
	OptimizedConfig.States = make(map[string]config.State)
	for stateName, state := range config.MainConfig.States {
		OptimizedConfig.States[stateName] = state
	}

	deletedStates = make(map[string]interface{})
	notToDeleteStates = make(map[string]interface{})

	for !OptimizedConfig.NoModification {
		zap.S().Info("Start optimize cycle")
		OptimizedConfig.NoModification = true
		zap.S().Infof("Number of states: %d", len(OptimizedConfig.States))
		var count = 0
		for stateName := range OptimizedConfig.States {
			OptimizedConfig.NoModification = !optimizeState(stateName) && OptimizedConfig.NoModification
			count++
			zap.S().Debugf("Number of states done: %d, currently: %s", count, stateName)
		}
	}

	config.MainConfig.States = make(map[string]config.State)

	for state := range deletedStates {
		if _, ok := notToDeleteStates[state]; !ok {
			zap.S().Debugf("Delete state: %s", state)
			delete(OptimizedConfig.States, state)
		}
	}

	zap.S().Infof("Number of states: %d", len(OptimizedConfig.States))

	for oStateName, oState := range OptimizedConfig.States {
		config.MainConfig.States[oStateName] = oState
	}

}

func optimizeState(stateName string) bool {
	stateModified := false
	state := OptimizedConfig.States[stateName]
	zap.S().Debugf("Optimizing state: %s, successors: %d", stateName, len(state.Successors))
	for nextName, nextRC := range state.Successors {
		if nextName == stateName {
			zap.S().Debugf("Don't loop")
			continue
		}
		nextState := OptimizedConfig.States[nextName]
		if len(nextState.Outputs) == 0 && !nextState.Preserve {
			found := false
			for nextNextName, nextNextRC := range nextState.Successors {
				found = true
				if stateRc, ok := state.Successors[nextNextName]; ok {
					zap.S().Infof("%s already in %s, or-ing", nextNextName, stateName)
					delete(state.Successors, nextNextName)
					mergedNextRC := mergeRootConditionsAnd(nextRC, nextNextRC)
					newRoot := mergeRootConditionsOr(stateRc, mergedNextRC)
					state.Successors[nextNextName] = newRoot
					zap.S().Debugf("State: %s, cond: %s", stateName, generation.ConditionAutoToString(stateRc))
					zap.S().Debugf("NextState %s, cond: %s", nextName, generation.ConditionAutoToString(nextRC))
					zap.S().Debugf("NextNextState: %s, cond: %s", nextNextName, generation.ConditionAutoToString(nextNextRC))
					zap.S().Debugf("Merged next cond: %s, New cond: %s", generation.ConditionAutoToString(mergedNextRC), generation.ConditionAutoToString(newRoot))
					zap.S().Debugf("Merged next cond: %s, New cond: %v", generation.ConditionAutoToString(mergedNextRC), newRoot)
				} else {
					state.Successors[nextNextName] = mergeRootConditionsAnd(nextRC, nextNextRC)
				}
			}
			for nextElseName := range nextState.DefaultSuccessor {
				found = true
				if stateRc, ok := state.Successors[nextElseName]; ok {
					zap.S().Debugf("%s already in %s, or-ing", nextElseName, stateName)
					delete(state.Successors, nextElseName)
					state.Successors[nextElseName] = mergeRootConditionsOr(stateRc, nextRC)
				} else {
					if len(nextRC.Conditions) == 0 && len(nextRC.Subconditions) == 0 {
						if state.DefaultSuccessor == nil {
							state.DefaultSuccessor = make(map[string]interface{})
						}
						state.DefaultSuccessor[nextElseName] = struct{}{}
					} else {
						state.Successors[nextElseName] = optimizeCondition(nextRC)
					}
				}
			}
			if found {
				delete(state.Successors, nextName)
				deletedStates[nextName] = struct{}{}
				stateModified = true
			} else {
				zap.S().Infof("No follow ups found for %s at %s", nextName, stateName)
				notToDeleteStates[nextName] = struct{}{}
			}
		} else {
			zap.S().Debugf("Next state %s has %d outputs and preserve is %v", nextName, len(nextState.Outputs), nextState.Preserve)
		}
	}
	for elseName := range state.DefaultSuccessor {
		if elseName == stateName {
			zap.S().Debugf("Don't loop")
			continue
		}
		elseState := OptimizedConfig.States[elseName]
		if len(elseState.Outputs) == 0 && !elseState.Preserve {
			found := false
			for nextNextName, nextNextRC := range elseState.Successors {
				found = true
				if stateRc, ok := state.Successors[nextNextName]; ok {
					zap.S().Debugf("%s already in %s, or-ing", nextNextName, stateName)
					delete(state.Successors, nextNextName)
					state.Successors[nextNextName] = mergeRootConditionsOr(stateRc, nextNextRC)
				} else {
					state.Successors[nextNextName] = optimizeCondition(nextNextRC)
				}
			}
			for nextElseName := range elseState.DefaultSuccessor {
				found = true
				if stateRc, ok := state.Successors[nextElseName]; ok {
					zap.S().Debugf("%s already in %s, or-ing", nextElseName, stateName)
					delete(state.Successors, nextElseName)
					state.Successors[nextElseName] = optimizeCondition(stateRc)
				} else {
					state.DefaultSuccessor[nextElseName] = struct{}{}
				}
			}
			if found {
				delete(state.DefaultSuccessor, elseName)
				deletedStates[elseName] = struct{}{}
				stateModified = true
			} else {
				zap.S().Infof("No follow ups found for %s at %s", elseName, stateName)
				notToDeleteStates[elseName] = struct{}{}
			}
		} else {
			zap.S().Debugf("Next state %s has %d outputs and preserve is %v", elseName, len(elseState.Outputs), elseState.Preserve)
		}
	}

	OptimizedConfig.States[stateName] = state

	return stateModified
}

func mergeRootConditionsAnd(rc ...config.Condition) config.Condition {
	var newRootCondition config.Condition
	newRootCondition.Subconditions = make([]config.Condition, len(rc))
	for index, condition := range rc {
		newRootCondition.Subconditions[index] = condition
	}
	newRootCondition.Type = config.ConditionType_And
	newRootCondition = optimizeCondition(newRootCondition)
	return newRootCondition
}

func mergeRootConditionsOr(rc ...config.Condition) config.Condition {
	var newRootCondition config.Condition
	newRootCondition.Subconditions = make([]config.Condition, len(rc))
	for index, condition := range rc {
		newRootCondition.Subconditions[index] = condition
	}
	newRootCondition.Type = config.ConditionType_Or
	newRootCondition = optimizeCondition(newRootCondition)
	return newRootCondition
}

func optimizeCondition(cond config.Condition) config.Condition {
	var optimizedSubconditions []config.Condition
	optimizedSubconditions = make([]config.Condition, 0)
	for _, subcond := range cond.Subconditions {
		if subcond.Type == cond.Type {
			cond.Conditions = mergeConditionInputs(cond.Conditions, subcond.Conditions)
			for _, subsubcond := range subcond.Subconditions {
				optimizedSubsubcond := optimizeCondition(subsubcond)
				if len(optimizedSubsubcond.Conditions) != 0 || len(optimizedSubsubcond.Subconditions) != 0 {
					optimizedSubconditions = append(optimizedSubconditions, optimizedSubsubcond)
				}
			}
		} else {
			optimizedSubcond := optimizeCondition(subcond)
			if len(optimizedSubcond.Conditions) != 0 || len(optimizedSubcond.Subconditions) != 0 {
				if len(optimizedSubcond.Subconditions) == 1 && len(optimizedSubcond.Conditions) == 0 {
					optimizedSubconditions = append(optimizedSubconditions, optimizedSubcond.Subconditions[0])
				} else if len(optimizedSubcond.Subconditions) == 0 && len(optimizedSubcond.Conditions) == 1 && len(optimizedSubcond.Conditions[0]) == 1 {
					cond.Conditions = mergeConditionInputs(cond.Conditions, optimizedSubcond.Conditions)
				} else {
					optimizedSubconditions = append(optimizedSubconditions, optimizedSubcond)
				}
			}
		}
	}
	cond.Subconditions = optimizedSubconditions
	if len(cond.Subconditions) == 1 {
		emptyConditions := true
		for _, condList := range cond.Conditions {
			emptyConditions = emptyConditions && (len(condList) == 0)
		}
		if emptyConditions {
			return optimizeCondition(cond.Subconditions[0])
		}
	}
	return cond
}

func mergeConditionInputs(condInputParameters ...[]map[string]int) []map[string]int {
	returnValue := make([]map[string]int, 1)
	returnValue[0] = make(map[string]int)
	for _, condInputParameter := range condInputParameters {
		for _, condInputs := range condInputParameter {
			for condInput, condValue := range condInputs {
				for retIndex, retList := range returnValue {
					if condExistingValue, ok := retList[condInput]; ok {
						if condExistingValue == condValue {
							break
						} else if len(returnValue) == (retIndex + 1) {
							newMap := make(map[string]int)
							newMap[condInput] = condValue
							returnValue = append(returnValue, newMap)
						}
					} else {
						returnValue[retIndex][condInput] = condValue
						break
					}
				}
			}
		}
	}
	if len(returnValue) == 1 && len(returnValue[0]) == 0 {
		return []map[string]int{}
	}
	return returnValue
}
