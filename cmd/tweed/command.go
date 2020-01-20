package main

type TweedCommand struct {
	Broker BrokerCommand `command:"broker" description:"Run the tweed broker."`
}
