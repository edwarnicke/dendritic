// Copyright 2018 Ed Warnicke

//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at

//      http://www.apache.org/licenses/LICENSE-2.0

//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/edwarnicke/dendritic/cmd/dendritic/libs/ads1299"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dendritic",
	Short: "Dendritic an EEG control CLI",
	Long:  `An EEG control CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		ads := ads1299.New()
		if err := ads.Init(); err != nil {
			fmt.Printf("Initialization error: %s\n ", err)
			os.Exit(1)
		}
		defer ads.Close()
		id, _ := ads.ReadReg(ads1299.ID)
		fmt.Printf("ID: 0x%x\n", id)
		regs, _ := ads.DumpRegs()
		fmt.Printf("Regdump: 0x%x\n", regs)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
