package main

type TweedCommand struct {
	Quiet bool `short:"q" long:"quiet"`
	JSON  bool `long:"json"`
	Debug bool `short:"D" long:"debug"`

	LogLevel string `long:"log-level"  env:"TWEED_LOG_LEVEL" default:"info"`

	Version func()        `short:"v" long:"version" description:"Print the version of Tweed and exit"`
	Broker  BrokerCommand `command:"broker" description:"Run the tweed broker."`

	Username string `short:"u" long:"username"  env:"TWEED_USERNAME"`
	Password string `short:"p" long:"password"  env:"TWEED_PASSWORD"`
	Tweed    string `short:"T" long:"tweed"     env:"TWEED_URL"`

	Oops OopsCommand `command:"oops"`
}
