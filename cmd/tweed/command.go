package main

type TweedCommand struct {
	Quiet bool `short:"q" long:"quiet" description:"Silence human readable output"`
	JSON  bool `long:"json"            description:"Output json"`
	Debug bool `short:"D" long:"debug" description:"Enable debug logging"`

	LogLevel string `long:"log-level"  env:"TWEED_LOG_LEVEL" default:"info" description:"Log verbosity level"`

	Version func()        `short:"v" long:"version" description:"Print the version of Tweed and exit"`
	Broker  BrokerCommand `command:"broker" description:"Run the tweed broker"`

	Username string `short:"u" long:"username"  env:"TWEED_USERNAME" description:"Username used to authenticate with Tweed broker"`
	Password string `short:"p" long:"password"  env:"TWEED_PASSWORD" description:"Password used to authenticate with Tweed broker"`
	Tweed    string `short:"T" long:"tweed"     env:"TWEED_URL"      description:"Address of the Tweed broker"`

	Catalog     CatalogCommand     `command:"catalog" alias:"cat"        description:"Show service catalog"`
	Instances   InstancesCommand   `command:"instances" alias:"ls"       description:"List created service instances"`
	Instance    InstanceCommand    `command:"instance" alias:"inst"      description:"Show a given service instance"`
	Provision   ProvisionCommand   `command:"provision" alias:"prov"     description:"Provision a new service instance"`
	Deprovision DeprovisionCommand `command:"deprovision" alias:"deprov" description:"Deprovision a given service instance"`
	Purge       PurgeCommand       `command:"purge"                      description:"Purge references of given instance from Tweed broker"`
	Bind        BindCommand        `command:"bind"                       description:"Create a binding to a given instance"`
	Unbind      UnbindCommand      `command:"unbind"                     description:"Unbind a given service instance binding"`
	Bindings    BindingsCommand    `command:"bindings"                   description:"List created bindings for a given instance"`
	Binding     BindingCommand     `command:"binding"                    description:"Show binding for a given instance by binding id"`
	Wait        WaitCommand        `command:"wait"                       description:"Wait for a task or instance to be in a given state"`
	Task        TaskCommand        `command:"task"                       description:"Show the output of a given task"`
	Reset       ResetCommand       `command:"reset"                      description:"Reset instance into given state"`
	Check       CheckCommand       `command:"check"                      description:"Check if the given service plan can be provisioned"`
	Log         LogCommand         `command:"log" alias:"logs"           description:"Show logs for a given service instance"`
	Files       FilesCommand       `command:"files"                      description:"List available files for given service instance"`
	File        FileCommand        `command:"file"                       description:"Retrieve file for given service instance"`
	Oops        OopsCommand        `command:"oops"                       hidden:"yes"`
}
