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
	"os"

	"github.com/kinnarr/fsmconverter/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var debugLoggerConfig bool
var LoggerConfig zap.Config

var rootCmd = &cobra.Command{
	Use:   "fsmconverter",
	Short: "FSMconverter a toml to verilog compiler",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debugLoggerConfig {
			LoggerConfig.Level.SetLevel(zap.DebugLevel)
		} else {
			LoggerConfig.Level.SetLevel(zap.InfoLevel)
		}
		zap.S().Infof("Current config dir is %s", config.FsmConfigRootDir)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&config.FsmConfigRootDir, "fsm-config-dir", "fsm", "search here for fsm config")
	rootCmd.PersistentFlags().BoolVar(&config.IgnoreUnknownStates, "ignore-unknown-states", false, "ignores unknown states in validation")
	rootCmd.PersistentFlags().BoolVar(&config.Optimize, "optimize", false, "optimize state of fsm")
	rootCmd.PersistentFlags().BoolVar(&debugLoggerConfig, "debug", false, "debugging output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
