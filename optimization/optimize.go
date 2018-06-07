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

func OptimizeConfig() {
	OptimizedConfig.NoModification = false
	OptimizedConfig.States = make(map[string]config.State)
	for stateName, state := range config.MainConfig.States {
		OptimizedConfig.States[stateName] = state
	}

	for !OptimizedConfig.NoModification {
		zap.S().Info("Start optimize cycle")
		OptimizedConfig.NoModification = true
		zap.S().Infof("Number of states: %d", len(OptimizedConfig.States))
		var count = 0
		for stateName := range OptimizedConfig.States {
			OptimizedConfig.NoModification = !optimizeState(stateName) && OptimizedConfig.NoModification
			count++
			zap.S().Infof("Number of states done: %d, currently: %s", count, stateName)
		}

		for oStateName, oState := range OptimizedConfig.States {
			if len(oState.Successors) == 0 {
				delete(OptimizedConfig.States, oStateName)
			}
		}
	}

	config.MainConfig.States = make(map[string]config.State)

	for oStateName, oState := range OptimizedConfig.States {
		if len(oState.Successors) > 0 {
			config.MainConfig.States[oStateName] = oState
		}
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
			// TODO: or combination if nextNextState already in state.Successors
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
			if found {
				delete(state.Successors, nextName)
				stateModified = true
			}
		} else {
			zap.S().Debugf("Next state %s has %d outputs and preserve is %v", nextName, len(nextState.Outputs), nextState.Preserve)
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
