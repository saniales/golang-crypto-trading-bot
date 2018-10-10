// Copyright © 2017 Alessandro Sanino <saninoale@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package bot

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

const (
	versionNumber = "0.0.1-pre-alpha"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gobot",
	Short: fmt.Sprintf("USAGE %s [OPTIONS]", os.Args[0]),
	Long:  fmt.Sprintf(`USAGE %s [OPTIONS] : see --help for details`, os.Args[0]),
	Run:   executeRootCommand,
}

var signals chan os.Signal

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		<-signals
		signal.Stop(signals)
		fmt.Println()
		fmt.Println("CTRL-C command received. Exiting...")
		os.Exit(0)
	}()

	RootCmd.Flags().BoolVarP(&rootFlags.Version, "version", "V", false, "show version information.")

	RootCmd.PersistentFlags().CountVarP(&GlobalFlags.Verbose, "verbose", "v", "show verbose information when trading : use multiple times to increase verbosity level.")
	RootCmd.PersistentFlags().StringVar(&GlobalFlags.ConfigFile, "config-file", "./.bot_config.yaml", "Config file path (default : ./.bot_config.yaml)")
}

func executeRootCommand(cmd *cobra.Command, args []string) {
	if rootFlags.Version {
		fmt.Printf("Golang Crypto Trading Bot v. %s\n", versionNumber)
	} else {
		cmd.Help()
	}
}
