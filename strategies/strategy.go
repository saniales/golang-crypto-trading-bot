package strategies

//Action provides which action should the bot take with the current configuration.
type Action int16

const (
	Buy       Action = iota //Represents a BUY action.
	Sell      Action = iota //Represents a SELL action.
	DoNothing Action = iota //Represents a DO-NOTHING action.
	Invalid   Action = iota //Represents an invalid action.
)

//Strategy represents a strategy to attach a bot on a market.
type Strategy interface {
	OnCandleUpdate() Action //Represents what to do when new data has been synced.
	SetupStrategy()         //Represents what to do when strategy is attached.
	TearDownStrategy()      //Represents what to do when strategy is detached.
	//OnNewCandle() Action    //Represents what to do when new candle is added to the graph.
}
