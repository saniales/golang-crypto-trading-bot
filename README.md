# golang-crypto-trading-bot

[![Release](https://img.shields.io/github/downloads/saniales/golang-crypto-trading-bot/total.svg)](https://github.com/yangwenmai/how-to-add-badge-in-github-readme/releases)
[![Documentation](https://godoc.org/github.com/saniales/golang-crypto-trading-bot?status.svg)](https://github.com/saniales/golang-crypto-trading-bot)
[![Travis CI](https://img.shields.io/travis/saniales/golang-crypto-trading-bot.svg)]((https://github.com/saniales/golang-crypto-trading-bot))
[![Go Report Card](https://goreportcard.com/badge/github.com/saniales/golang-crypto-trading-bot.svg?branch=master)](https://github.com/saniales/golang-crypto-trading-bot)
[![GitHub release](https://img.shields.io/github/release/saniales/golang-crypto-trading-bot.svg)](https://github.com/saniales/golang-crypto-trading-bot/releases)
[![license](https://img.shields.io/github/license/saniales/golang-crypto-trading-bot.svg?maxAge=2592000)](https://github.com/saniales/golang-crypto-trading-bot/LICENSE)


A golang implementation of a console-based trading bot for cryptocurrency exchanges, can be deployed to heroku too. 

# Supported Exchanges
Bittrex, Poloniex, Binance, Bitfinex

# Usage

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

For strategy reference see the [godocs](https://godoc.org/github.com/saniales/golang-crypto-trading-bot).

# Current version
[pre-alpha 0.0.0.1]

# Donate
Feel free to donate:

| METHOD 	| ADDRESS                                   	|
|--------	|--------------------------------------------	|
| Paypal 	| https://paypal.me/AlessandroSanino         	|
| BTC    	| 1DVgmv6jkUiGrnuEv1swdGRyhQsZjX9MT3         	|
| XVG    	| DFstPiWFXjX8UCyUCxfeVpk6JkgaLBSNvS         	|
| ETH    	| 0x2fe7bd8a41e91e9284aada0055dbb15ecececf02 	|
| USDT   	| 18obCEVmbT6MHXDcPoFwnUuCmkttLbK5Xo         	|
