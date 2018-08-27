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

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"bytes"

	"go.uber.org/zap"

	"github.com/kinnarr/fsmconverter/config"
	"github.com/kinnarr/fsmconverter/generation"
	"github.com/kinnarr/fsmconverter/optimization"
	"github.com/kinnarr/fsmconverter/validation"
)

func init() {
	rootCmd.AddCommand(printStateCmd)
}

var printStateCmd = &cobra.Command{
	Use:   "printstate [state]",
	Short: "Prints the given fsm state",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !config.ParseConfig() {
			return
		}

		if !validation.ValidateStates() || !validation.ValidateDefaults() {
			zap.S().Errorf("Validation failed! See errors above!\n")
			return
		}

		if config.Optimize {
			optimization.OptimizeConfig()

			if !validation.ValidateStates() || !validation.ValidateDefaults() {
				zap.S().Errorf("Validation failed! See errors above!\n")
				return
			}
		}

		var fsmOuputBuffer bytes.Buffer

		if state, ok := config.MainConfig.States[args[0]]; ok {
			fsmOuputBuffer.WriteString(fmt.Sprintf("State: %s (#Outputs: %d)\n", args[0], len(state.Outputs)))
			for nextName, next := range state.Successors {
				fsmOuputBuffer.WriteString(generation.RootConditionToString(next, nextName))
			}
			for elseName := range state.DefaultSuccessor {
				fsmOuputBuffer.WriteString(fmt.Sprintf("\telse: next <= %s\n", elseName))
			}
		}

		fmt.Print(fsmOuputBuffer.String())

	},
}
