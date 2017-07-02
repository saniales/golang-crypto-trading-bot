// Copyright © 2017 Alessandro Sanino <saninoale@gmail.com>
// {{.copyright}}
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
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"

	"github.com/AlessandroSanino1994/gobot/environment"
	"github.com/spf13/cobra"
)

var initFlags struct {
	ConfigFile string
	Exchange   string
	Strategies []struct {
		Market   string
		Strategy string
	}
	BTCAddress string
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the bot to trade",
	Long: `Initializes the trading bot: it will ask several questions to properly create a conf file.
	It must be run prior any other command if config file is not present.`,
	Run: executeInitCommand,
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&initFlags.ConfigFile, "import", "", "imports configuration from a file.")
}

func executeInitCommand(cmd *cobra.Command, args []string) {
	initConfig()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if initFlags.ConfigFile != "" {
		//try first to unmarshal the file to check if it is correct format.
		content, err := ioutil.ReadFile(initFlags.ConfigFile)
		if err != nil {
			fmt.Print("Error while opening the config file provided")
			if GlobalFlags.Verbose > 0 {
				fmt.Printf(": %s", err.Error())
			}
			fmt.Println()
			return
		}
		var checker environment.BotConfig
		err = yaml.Unmarshal(content, &checker)
		if err != nil {
			fmt.Print("Cannot load provided configuration file")
			if GlobalFlags.Verbose > 0 {
				fmt.Printf(": %s", err.Error())
			}
			fmt.Println()
			return
		}
		err = ioutil.WriteFile("./.gobot", content, 0666)
		if err != nil {
			fmt.Print("Cannot write new configuration file")
			if GlobalFlags.Verbose > 0 {
				fmt.Printf(": %s", err.Error())
			}
			fmt.Println()
			return
		}
	} else {
		generateInitFile()
	}
}

func generateInitFile() {
	configs := environment.BotConfig{}
	fmt.Println("Which exchange are you going to use?")
	fmt.Scanln(&configs.Exchange.ExchangeName)
	fmt.Println("Please provide Public Key for that exchange.")
	fmt.Scanln(&configs.Exchange.PublicKey)
	fmt.Println("Please provide Secret Key for that exchange.")
	fmt.Scanln(&configs.Exchange.SecretKey)

	for {
		var YesNo string
		for YesNo != "Y" && YesNo != "n" {
			fmt.Println("Do you want to add a strategy binded to a market? (Y/n)")
			fmt.Scanln(&YesNo)
		}
		if YesNo == "n" {
			break
		}
		var tempStrategyAppliance environment.StrategyConfig
		fmt.Println("Please Enter Market Name using short notation " +
			"(e.g. BTC-ETH for a Bitcoin-Ethereum market).")
		fmt.Scanln(&tempStrategyAppliance.Market)
		fmt.Println("Please Enter The Name of the strategy you want to use\n" +
			"in this market (must be in the system)")
		fmt.Scanln(&tempStrategyAppliance.Strategy)
		configs.Strategies = append(configs.Strategies, tempStrategyAppliance)
	}

	//preview the contents of the file to be written, then creates a new file in .
	contentToBeWritten, err := yaml.Marshal(configs)
	if err != nil {
		fmt.Print("Error while creating the content for the new config file")
		if GlobalFlags.Verbose > 0 {
			fmt.Printf(": %s", err.Error())
		}
		fmt.Println()
		return
	}
	fmt.Println("The following content:")
	fmt.Println(string(contentToBeWritten))
	fmt.Println("is going to be written on ./.gobot, is it ok? (Y/n)")

	var YesNo string
	for YesNo != "Y" && YesNo != "n" {
		fmt.Scanln(&YesNo)
	}
	if YesNo == "Y" {
		err := ioutil.WriteFile("./.gobot", contentToBeWritten, 0666)
		if err != nil {
			fmt.Print("Error while writing content to new config file")
			if GlobalFlags.Verbose > 0 {
				fmt.Printf(": %s", err.Error())
			}
			fmt.Println()
		} else {
			fmt.Println("Config file created on ./.gobot\nNow you can use gobot with this configuration.\nHappy Trading, folk :)")
		}
		return
	}
	fmt.Println("You chose not to write the content to configuration file.\n" +
		"You can relaunch this command again to create another configuration.\n" +
		"This bot won't work until it has a valid configuration file.")
}
