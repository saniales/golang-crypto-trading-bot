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

package environment

//ExchangeConfig Represents a configuration for an API Connection to an exchange.
//
//Can be used to generate an ExchangeWrapper.
type ExchangeConfig struct {
	ExchangeName string `yaml:"exchange"`   //Represents the exchange name.
	PublicKey    string `yaml:"public_key"` //Represents the public key used to connect to Exchange API.
	SecretKey    string `yaml:"secret_key"` //Represents the secret key used to connect to Exchange API.
}

// StrategyConfig contains where a strategy will be applied in the specified exchange.
type StrategyConfig struct {
	Market   string `yaml:"market"`   //Represents the market where the strategy is applied.
	Strategy string `yaml:"strategy"` //Represents the applied strategy name: must be one in the system.
}

// BotConfig contains all config data of the bot, which can be also loaded from config file.
type BotConfig struct {
	Exchange   ExchangeConfig   `yaml:"exchange_config"` //Represents the current exchange configuration.
	Strategies []StrategyConfig `yaml:"strategies"`      //Represents the current strategies adopted by the bot.
}

//type Configs []ExchangeConfig
