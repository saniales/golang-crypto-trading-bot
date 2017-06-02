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

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

//GlobalFlags provides flag definitions valid for the whole system.
var GlobalFlags struct {
	Verbose bool //Tells the program to print everything to screen.
}

//rootFlags provides flag definitions valid for root command.
var rootFlags struct {
	Version bool
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gobot",
	Short: "USAGE gobot [OPTIONS].",
	Long:  `USAGE gobot [OPTIONS] : see --help for details`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVar(&GlobalFlags.Verbose, "verbose", false, "show verbose information (default = false) when trading.")

	RootCmd.Flags().BoolVar(&rootFlags.Version, "version", false, "show version information (default = false).")
}
