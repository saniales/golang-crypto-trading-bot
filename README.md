# golang-crypto-trading-bot

[![GoDoc](https://godoc.org/github.com/saniales/golang-crypto-trading-bot?status.svg)](https://godoc.org/github.com/saniales/golang-crypto-trading-bot)
[![Travis CI](https://img.shields.io/travis/saniales/golang-crypto-trading-bot.svg)]((https://travis-ci.org/saniales/golang-crypto-trading-bot))
[![GitHub release](https://img.shields.io/github/release/saniales/golang-crypto-trading-bot.svg)](https://github.com/saniales/golang-crypto-trading-bot/releases)
[![license](https://img.shields.io/github/license/saniales/golang-crypto-trading-bot.svg?maxAge=2592000)](https://github.com/saniales/golang-crypto-trading-bot/LICENSE)


A golang implementation of a console-based trading bot for cryptocurrency exchanges. 

## Supported Exchanges
Bittrex, Poloniex, Binance, Bitfinex and Kraken, other in progress.

## Usage

Download a release or directly build the code from this repository.
``` bash
$ go get github.com/saniales/golang-crypto-trading-bot
```

If you need to, you can create a strategy and bind it to the bot:
``` go
import bot "github.com/saniales/golang-crypto-trading-bot/cmd"

bot.AddCustomStrategy(myStrategy)
bot.Execute()
```

For strategy reference see the [Godoc documentation](https://godoc.org/github.com/saniales/golang-crypto-trading-bot).

# Configuration file template
Create a configuration file from this example or run the `init` command of the compiled executable.
``` yaml
exchange_configs: 
  - exchange: bittrex
    public_key: your_bittrex_public_key
    secret_key: your_bittrex_secret_key
    websocket_enabled: true
  - exchange: binance
    public_key: your_binance_public_key
    secret_key: your_binance_secret_key
    websocket_enabled: true
  - exchange: bitfinex
    public_key: your_bitfinex_public_key
    secret_key: your_bitfinex_secret_key
    websocket_enabled: true
strategies:
  - strategy: your_strategy_name
    markets:
      - market: market_logical_name
        bindings:
        - exchange: bittrex
          market_name: market_name_on_bittrex
        - exchange: binance
          market_name: market_name_on_binance
        - exchange: bitfinex
          market_name: market_name_on_bitfinex
      - market: another_market_logical_name
        bindings:
        - exchange: bittrex
          market_name: market_name_on_bittrex
        - exchange: binance
          market_name: market_name_on_binance
        - exchange: bitfinex
          market_name: market_name_on_bitfinex
```

# Donate
Feel free to donate:

| METHOD 	| ADDRESS                                   	|
|--------	|--------------------------------------------	|
| Paypal 	| https://paypal.me/AlessandroSanino         	|
| BTC    	| 1DVgmv6jkUiGrnuEv1swdGRyhQsZjX9MT3         	|
| XVG    	| DFstPiWFXjX8UCyUCxfeVpk6JkgaLBSNvS         	|
| ETH    	| 0x2fe7bd8a41e91e9284aada0055dbb15ecececf02 	|
| USDT   	| 18obCEVmbT6MHXDcPoFwnUuCmkttLbK5Xo         	|
