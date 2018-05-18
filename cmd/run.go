// Copyright © 2018 Charles Corebtt <nafredy@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/asciifaceman/p2p2p/service"
	"github.com/spf13/cobra"
)

// Initialization variables and constants
var (
	nodeName string
	nodePort int
	nodePool string
)

const (
	defaultName string = "node"
	defaultPort int    = 5000
	defaultHost string = "0.0.0.0"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// instantiate & run server
		server := service.Server{Host: defaultHost, Name: nodeName, Port: nodePort}
		err := server.Start(nodePool)
		if err != nil {
			fmt.Println("Error starting node: ", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Initialize global random source
	rand.Seed(time.Now().UnixNano())

	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&nodeName, "name", "n", fmt.Sprintf("%s-%d", defaultName, rand.Intn(100)), "The name of the node being initialized.")
	runCmd.Flags().StringVarP(&nodePool, "bootnodes", "b", "", "The bootnode pool (if applicable). Format: -b localhost:3031,localhost3032,...")
	runCmd.Flags().IntVarP(&nodePort, "port", "p", defaultPort, "The port of the node being initialized.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
