package exchangeWrappers

//ExchangeConfig Represents a configuration for an API Connection to an exchange.
//
//Can be used to generate an ExchangeWrapper
type ExchangeConfig struct {
	ExchangeName string `json:"exchange"`   //Represents the exchange name.
	PublicKey    string `json:"public-key"` //Represents the public key used to connect to Exchange API.
	SecretKey    string `json:"secret-key"` //Represents the secret key used to connect to Exchange API.
}

type Configs []ExchangeConfig
