package botHelpers

import (
	"github.com/AlessandroSanino1994/gobot/environment"
	"github.com/AlessandroSanino1994/gobot/exchangeWrappers"
	flags "github.com/jessevdk/go-flags"
)

//CommandLineOptions represents the command line args for the bot.
type CommandLineOptions struct {
	Verbose    bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
	Simulate   bool   `short:"s" long:"simulate" description:"Run the bot in simulation mode (only simulate trades)."`
	Exchange   string `short:"e" long:"exchange" description:"Exchange to connect to trade." required:"true"`
	ConfigFile string `short:"c" long:"config-file" description:"Path to the bot configuration file" required:"true"`
}

//InitArgs initializes the arguments from command line, or returns error.
func InitArgs() (CommandLineOptions, error) {
	var commandOptions CommandLineOptions
	_, err := flags.Parse(&commandOptions)
	return commandOptions, err
}

//InitExchange initialize a new ExchangeWrapper binded to the specified exchange provided.
func InitExchange(exchangeConfig exchangeWrappers.ExchangeConfig) exchangeWrappers.ExchangeWrapper {
	switch exchangeConfig.ExchangeName {
	case "bittrex":
		return exchangeWrappers.NewBittrexWrapper(exchangeConfig.PublicKey, exchangeConfig.SecretKey)
	case "poloniex":
		return nil
	default:
		return nil
	}
}

//InitMarkets uses ExchangeWrapper to find info about markets and initialize them.
func InitMarkets(exchange exchangeWrappers.ExchangeWrapper) ([]environment.Market, error) {
	return exchange.GetMarkets()
}
