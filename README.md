<p align="center"><img src="https://res.cloudinary.com/saniales-github/image/upload/v1541856660/saniales/golang-crypto-trading-bot/logo.png" width="360"></p>
<p align="center">
  <a href="https://github.com/avelino/awesome-go#other-software"><img src="https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg" alt="Listed on Awesome Go"></img></a>
  <a href="https://github.com/saniales/golang-crypto-trading-bot/actions"><img src="https://github.com/saniales/golang-crypto-trading-bot/workflows/Go%20Build/badge.svg?branch=develop" alt="Develop Branch Build Status"></img></a>
  <a href="https://godoc.org/github.com/saniales/golang-crypto-trading-bot"><img src="https://godoc.org/github.com/saniales/golang-crypto-trading-bot?status.svg" alt="Godoc reference"></a>
  <a href="https://github.com/saniales/golang-crypto-trading-bot/releases"><img src="https://img.shields.io/github/release/saniales/golang-crypto-trading-bot.svg" alt="Last Release"></a>
  <a href="https://github.com/saniales/golang-crypto-trading-bot/LICENSE"><img src="https://img.shields.io/github/license/saniales/golang-crypto-trading-bot.svg?maxAge=2592000" alt="License"></a>
  <a href="https://goreportcard.com/report/github.com/saniales/golang-crypto-trading-bot"><img src="https://goreportcard.com/badge/github.com/saniales/golang-crypto-trading-bot" alt="Goreportcard" /></a>
</p>

# Golang Crypto Trading Bot

A golang implementation of a console-based trading bot for cryptocurrency exchanges.

## Usage

Download a release or directly build the code from this repository.

``` bash
go get github.com/saniales/golang-crypto-trading-bot
```

If you need to, you can create a strategy and bind it to the bot:

``` go
import bot "github.com/saniales/golang-crypto-trading-bot/cmd"

func main() {
    bot.AddCustomStrategy(examples.MyStrategy)
    bot.Execute()
}
```

For strategy reference see the [Godoc documentation](https://godoc.org/github.com/saniales/golang-crypto-trading-bot).

## Simulation Mode

If enabled, the bot will do paper trading, as it will execute fake orders in a sandbox environment.

A Fake balance for each coin must be specified for each exchange if simulation mode is enabled.

## Supported Exchanges

| Exchange Name | REST Supported    | Websocket Support |
| ------------- |------------------ | ----------------- |
| Bittrex       | Yes               | No                |
| Poloniex      | Yes               | Yes               |
| Kraken        | Yes (no withdraw) | No                |
| Bitfinex      | Yes               | Yes               |
| Binance       | Yes               | Yes               |
| Kucoin        | Yes               | No                |
| HitBtc        | Yes               | Yes               |

## Configuration file template

Create a configuration file from this example or run the `init` command of the compiled executable.

``` yaml
simulation_mode: true # if you want to enable simulation mode.
exchange_configs:
  - exchange: bitfinex
    public_key: bitfinex_public_key
    secret_key: bitfinex_secret_key
    deposit_addresses:
      BTC: bitfinex_deposit_address_btc
      ETH: bitfinex_deposit_address_eth
      ZEC: bitfinex_deposit_address_zec
    fake_balances: # used only if simulation mode is enabled, can be omitted if not enabled.
      BTC: 100
      ETH: 100
      ZEC: 100
      ETC: 100
  - exchange: hitbtc
    public_key: hitbtc_public_key
    secret_key: hitbtc_secret_key
    deposit_addresses:
      BTC : hitbtc_deposit_address_btc
      ETH: hitbtc_deposit_address_eth
      ZEC: hitbtc_deposit_address_zec
    fake_balances:
      BTC: 100
      ETH: 100
      ZEC: 100
      ETC: 100
strategies:
  - strategy: strategy_name
    markets:
      - market: ETH-BTC
        bindings:
        - exchange: bitfinex
          market_name: ETHBTC
        - exchange: hitbtc
          market_name: ETHBTC
      - market: ZEC-BTC
        bindings:
        - exchange: bitfinex
          market_name: ZECBTC
        - exchange: hitbtc
          market_name: ZECBTC
      - market: ETC-BTC
        bindings:
        - exchange: bitfinex
          market_name: ETCBTC
        - exchange: hitbtc
          market_name: ETCBTC
```

## Donate

Feel free to donate:

| METHOD  | ADDRESS                                     |
|-------- |-------------------------------------------- |
| Paypal  | https://paypal.me/AlessandroSanino          |
| BTC     | 1DVgmv6jkUiGrnuEv1swdGRyhQsZjX9MT3          |
| XVG     | DFstPiWFXjX8UCyUCxfeVpk6JkgaLBSNvS          |
| ETH     | 0x2fe7bd8a41e91e9284aada0055dbb15ecececf02  |
| USDT    | 18obCEVmbT6MHXDcPoFwnUuCmkttLbK5Xo          |
