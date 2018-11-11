package strategies

import (
	"errors"

	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/saniales/golang-crypto-trading-bot/exchanges"
)

// WebsocketStrategy polls data from a websocket in real-time.
//
//     NOTE: The update function must be handled in the websocket strategy-model.
type WebsocketStrategy struct {
	Model StrategyModel
}

// Name returns the name of the strategy.
func (wss WebsocketStrategy) Name() string {
	return wss.Model.Name
}

// String returns a string representation of the object.
func (wss WebsocketStrategy) String() string {
	return wss.Name()
}

// Apply executes Cyclically the On Update, basing on provided interval.
func (wss WebsocketStrategy) Apply(wrappers []exchanges.ExchangeWrapper, markets []*environment.Market) {
	var err error

	hasSetupFunc := wss.Model.Setup != nil
	hasTearDownFunc := wss.Model.TearDown != nil
	hasUpdateFunc := wss.Model.OnUpdate != nil
	hasErrorFunc := wss.Model.OnError != nil

	if hasSetupFunc {
		err = wss.Model.Setup(wrappers, markets)
		if err != nil && hasErrorFunc {
			wss.Model.OnError(err)
		}
	}

	// update is handled by the developer externally, here we just checked for existence.
	if !hasUpdateFunc {
		_err := errors.New("OnUpdate func cannot be empty")
		if hasErrorFunc {
			wss.Model.OnError(_err)
		} else {
			panic(_err)
		}
	}

	if hasTearDownFunc {
		err = wss.Model.TearDown(wrappers, markets)
		if err != nil && hasErrorFunc {
			wss.Model.OnError(err)
		}
	}
}
