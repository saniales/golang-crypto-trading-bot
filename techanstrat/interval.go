package techanstrat

import (
	"fmt"
	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/saniales/golang-crypto-trading-bot/exchanges"
	"github.com/saniales/golang-crypto-trading-bot/strategies"
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	// SMA Trends
	SMA_long  = 1000
	SMA_short = 50

	// BULL
	BULL_RSI      = 10
	BULL_RSI_high = 80
	BULL_RSI_low  = 60

	// BEAR
	BEAR_RSI      = 15
	BEAR_RSI_high = 50
	BEAR_RSI_low  = 20

	// BULL/BEAR is defined by the longer SMA trends
	// if SHORT over LONG = BULL
	// if SHORT under LONG = BEAR
)

var indicator *LastNIndicator = &LastNIndicator{N: SMA_long}
var record *techan.TradingRecord = techan.NewTradingRecord()
var currentIndex = 0

type LastNIndicator struct {
	N         int
	datum     []float64
	baseIndex int
}

func (i *LastNIndicator) Calculate(index int) big.Decimal {
	return big.NewDecimal(i.datum[index-i.baseIndex])
}

func (i *LastNIndicator) addNewestCandle(candle environment.CandleStick) {}

func (i *LastNIndicator) addNewestBBO(bbo *environment.MarketSummary) {
	price, _ := bbo.Last.Float64()
	i.datum = append(i.datum, price)
	if len(i.datum) > i.N {
		i.datum = i.datum[1:]
		i.baseIndex++
	}
}

type relativeStrengthIndexIndicator struct {
	indicator techan.Indicator
	timeframe int
}

func NewRelativeStrengthIndexIndicator(indicator techan.Indicator, timeframe int) techan.Indicator {
	return relativeStrengthIndexIndicator{
		indicator: indicator,
		timeframe: timeframe,
	}
}

func (rsi relativeStrengthIndexIndicator) Calculate(index int) big.Decimal {
	if index == 0 {
		return big.ZERO
	}
	if index < rsi.timeframe {
		return big.ZERO
	}

	gain, loss := big.ZERO, big.ZERO
	for i := 0; i < rsi.timeframe; i++ {
		today := rsi.indicator.Calculate(index - i)
		yesterday := rsi.indicator.Calculate(index - i - 1)
		if today.GT(yesterday) {
			gain = gain.Add(today.Sub(yesterday))
		} else {
			loss = loss.Add(yesterday.Sub(today))
		}
		//log.Println(today, yesterday, gain, loss)
	}
	avgGain := gain.Div(big.NewDecimal(float64(rsi.timeframe)))
	avgLoss := loss.Div(big.NewDecimal(float64(rsi.timeframe)))

	oneHundred := big.NewDecimal(100)
	if loss.EQ(big.ZERO) {
		return oneHundred
	}
	return avgGain.Div(avgGain.Add(avgLoss)).Mul(oneHundred)
}

var Interval = strategies.IntervalStrategy{
	Model: strategies.StrategyModel{
		Name: "RSIBullBear",
		Setup: func(wrappers []exchanges.ExchangeWrapper, markets []*environment.Market) error {
			fmt.Println("RSIBullBear starting")
			return nil
		},
		OnUpdate: func(wrappers []exchanges.ExchangeWrapper, markets []*environment.Market) (err error) {
			//candles, err := wrappers[0].GetCandles(markets[0])
			//if err != nil {
			//	return err
			//}
			defer func() { if err == nil {currentIndex++} else {err = nil} }()

			m, err := wrappers[0].GetMarketSummary(markets[0])
			if err != nil {
				logrus.Error("OnUpdate", err)
				return
			}

			indicator.addNewestBBO(m)
			maFast := techan.NewSimpleMovingAverage(indicator, SMA_short).Calculate(currentIndex)
			maSlow := techan.NewSimpleMovingAverage(indicator, SMA_long).Calculate(currentIndex)

			var strategy = techan.RuleStrategy{
				UnstablePeriod: SMA_long,
			}
			var situation string
			var rsi big.Decimal
			if maFast.LT(maSlow) {
				// bear
				situation = "bear"

				rsiIndicator := NewRelativeStrengthIndexIndicator(indicator, BEAR_RSI)
				strategy.EntryRule = techan.And(
					techan.UnderIndicatorRule{First: rsiIndicator, Second: techan.NewConstantIndicator(BEAR_RSI_low)},
					techan.PositionNewRule{})

				strategy.ExitRule = techan.And(
					techan.OverIndicatorRule{First: rsiIndicator, Second: techan.NewConstantIndicator(BEAR_RSI_high)},
					techan.PositionOpenRule{})

				rsi = rsiIndicator.Calculate(currentIndex)
			} else {
				// bull
				situation = "bull"

				rsiIndicator := NewRelativeStrengthIndexIndicator(indicator, BULL_RSI)
				strategy.EntryRule = techan.And(
					techan.UnderIndicatorRule{First: rsiIndicator, Second: techan.NewConstantIndicator(BULL_RSI_low)},
					techan.PositionNewRule{})

				strategy.ExitRule = techan.And(
					techan.OverIndicatorRule{First: rsiIndicator, Second: techan.NewConstantIndicator(BULL_RSI_high)},
					techan.PositionOpenRule{})

				rsi = rsiIndicator.Calculate(currentIndex)
			}
			fmt.Fprintln(os.Stdout, fmt.Sprintf("%v, %s, %s, %s, %s, %s, %s",
				currentIndex, time.Now().Format(time.RFC3339), m.Last, maFast, maSlow, rsi, situation))

			if strategy.ShouldEnter(currentIndex, record) {
				_, err = wrappers[0].BuyMarket(markets[0], 1)
				if err != nil {
					logrus.Error("OnUpdate: BuyMarket", err)
					return
				}

				order := techan.Order{
					Side: techan.BUY,
				}
				record.Operate(order)
			}
			if strategy.ShouldExit(currentIndex, record) {
				_, err = wrappers[0].SellMarket(markets[0], 1)
				if err != nil {
					logrus.Error("OnUpdate: SellMarket", err)
					return
				}
				record.Operate(techan.Order{
					Side: techan.SELL,
				})
			}
			return nil
		},
		OnError: func(err error) {
			fmt.Println("OnError", err)
		},
		TearDown: func(wrappers []exchanges.ExchangeWrapper, markets []*environment.Market) error {
			fmt.Println("Watch5Sec exited")
			return nil
		},
	},
	Interval: time.Second * 1,
}
