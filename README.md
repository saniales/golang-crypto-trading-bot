# golang-crypto-trading-bot
A golang implementation of a console-based trading bot for cryptocurrency exchanges, can be deployed to heroku too. 

# Supported Exchanges
Bittrex, Poloniex (both works in progress)

# Usage
Create a strategy by implementing Strategy interface and add it to the bot, then compile and execute the bot specifying the --strategy flag and --exchange flag

`go run gobot.go --exchange bittrex --strategy myCustomStrategyName`

You can use the --simulate flag to avoid trade but just simulate them.
There is also a --watch flag to print info about markets to the console continuously.

Kill the bot by using CTRL+C  or SIGINT signal (on linux, Command+C on MacOsX).

# Current version
[pre-alpha 0.0.0.1]
