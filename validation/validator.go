package validation

import (
	"github.com/kinnarr/fsmconverter/config"
	"go.uber.org/zap"
)

func ValidateStates() bool {
	returnValue := true
	for stateName, state := range config.MainConfig.States {
		for nextName, next := range state.Successors {
			if _, ok := config.MainConfig.States[nextName]; !ok {
				zap.S().Errorf("Unknown state %s for state %s\n", nextName, stateName)
				returnValue = false
			} else {
				returnValue = validateRootCondition(next, nextName) && returnValue
			}
		}
		if len(state.DefaultSuccessor) > 1 {
			zap.S().Errorf("Only one else state allowed! State %s\n", stateName)
			returnValue = false
		} else {
			for elseName := range state.DefaultSuccessor {
				if _, ok := config.MainConfig.States[elseName]; !ok {
					zap.S().Errorf("Unknown else state %s for state %s\n", elseName, stateName)
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
	for conditionName := range c.Conditions {
		if _, ok := config.MainConfig.Inputs[conditionName]; !ok {
			zap.S().Errorf("Could not find input from condition for next state '%s': %s", nextName, conditionName)
			returnValue = false
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
	return returnValue
}
