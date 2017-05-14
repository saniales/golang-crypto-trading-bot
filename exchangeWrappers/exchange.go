package exchangeWrappers

//Represents an exchange.
type Exchange interface {
	Name() string      //Represents the name of the exchange.
	PublicKey() string //Represents the public key to connect to exchange.
	SecretKey() string //Represents the secret key to connect to exchange.

	//TODO: wrap BUY, SELL, GETTRADES, GETMARKETS, GETORDERBOOK
}
