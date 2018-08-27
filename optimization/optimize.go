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
					zap.S().Debugf("%s already in %s, or-ing", nextNextName, stateName)
					delete(state.Successors, nextNextName)
					state.Successors[nextNextName] = mergeRootConditionsOr(stateRc, mergeRootConditionsAnd(nextRC, nextNextRC))
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
					state.Successors[nextElseName] = nextRC
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
					state.Successors[nextNextName] = nextNextRC
				}
			}
			for nextElseName := range elseState.DefaultSuccessor {
				found = true
				if stateRc, ok := state.Successors[nextElseName]; ok {
					zap.S().Debugf("%s already in %s, or-ing", nextElseName, stateName)
					delete(state.Successors, nextElseName)
					state.Successors[nextElseName] = stateRc
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

	return stateModified
}

func mergeRootConditionsAnd(rc ...config.RootCondition) config.RootCondition {
	var newRootCondition config.RootCondition
	newCondition := new(config.Condition)
	newCondition.And = mergeConditions(rc[0].And, rc[1].And)
	newCondition.Or = mergeConditions(rc[0].Or, rc[1].Or)
	newRootCondition.And = newCondition
	return newRootCondition
}

func mergeRootConditionsOr(rc ...config.RootCondition) config.RootCondition {
	var newRootCondition config.RootCondition
	newCondition := new(config.Condition)
	newCondition.And = mergeConditions(rc[0].And, rc[1].And)
	newCondition.Or = mergeConditions(rc[0].Or, rc[1].Or)
	newRootCondition.Or = newCondition
	return newRootCondition
}

func mergeConditions(c1, c2 *config.Condition) *config.Condition {
	if c1 == nil {
		return c2
	}
	if c2 == nil {
		return c1
	}
	newCondition := new(config.Condition)
	newCondition.And = mergeConditions(c1.And, c2.And)
	newCondition.Or = mergeConditions(c1.Or, c2.Or)
	newCondition.Conditions = append(c1.Conditions, c2.Conditions...)
	return newCondition
}
