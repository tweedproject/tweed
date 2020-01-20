package main

type TweedCommand struct {
	Quiet bool `short:"q" long:"quiet"`
	JSON  bool `long:"json"`
	Debug bool `short:"D" long:"debug"`

	LogLevel string `long:"log-level"  env:"TWEED_LOG_LEVEL" default:"info"`

	Version func()        `short:"v" long:"version" description:"Print the version of Tweed and exit"`
	Broker  BrokerCommand `command:"broker" description:"Run the tweed broker"`

	Username string `short:"u" long:"username"  env:"TWEED_USERNAME"`
	Password string `short:"p" long:"password"  env:"TWEED_PASSWORD"`
	Tweed    string `short:"T" long:"tweed"     env:"TWEED_URL"`

	Catalog   CatalogCommand   `command:"catalog" alias:"cat"  description:"Show service catalog"`
	Instances InstancesCommand `command:"instances" alias:"ls" description:"List created service instances"`
	Bind      BindCommand      `command:"bind"                 description:"Create a binding to a given instance"`
	Bindings  BindingsCommand  `command:"bindings"             description:"List created bindings for a given instance"`
	Binding   BindingCommand   `command:"binding"              description:"Show binding for a given instance by binding id"`
	Wait      WaitCommand      `command:"wait"                 description:"Wait for a task or instance to be in a given state"`
	Oops      OopsCommand      `command:"oops"`
}
