package main

var opts struct {
	Help    bool `cli:"-h, --help"`
	Version bool `cli:"-v, --version"`
	Quiet   bool `cli:"-q, --quiet"`
	JSON    bool `cli:"--json"`
	Debug   bool `cli:"-D, --debug"`

	LogLevel string `cli:"--log-level", env:"TWEED_LOG_LEVEL"`

	Username string `cli:"-u, --username"   env:"TWEED_USERNAME"`
	Password string `cli:"-p, --password"   env:"TWEED_PASSWORD"`
	Tweed    string `cli:"-T, --tweed"      env:"TWEED_URL"`

	Reset struct {
		State string `cli:"-s, --state"`
	} `cli:"reset"`

	Check struct{} `cli:"check"`

	Provision struct {
		ID     string   `cli:"--as"`
		Params []string `cli:"-P, --param"`
		Wait   bool     `cli:"--wait, --no-wait"`
	} `cli:"provision,prov,create-service,new-service"`

	Task struct {
		Wait bool `cli:"--wait, --no-wait"`
	} `cli:"task"`

	Deprovision struct {
		Wait bool `cli:"--wait, --no-wait"`
	} `cli:"deprovision,deprov"`

	Purge struct {
		All bool `cli:"-a, --all"`
	} `cli:"purge"`

	Instance struct{} `cli:"instance,inst"`

	Log struct {
		Tail   int  `cli:"-n, --tail"`
		Follow bool `cli:"-f, --follow"`
	} `cli:"log,logs"`

	Files struct{} `cli:"files"`
	File  struct{} `cli:"file"`

	Unbind struct {
		Wait bool `cli:"--wait, --no-wait"`
	} `cli:"unbind"`
}

func init() {
	opts.LogLevel = "info"

	opts.Provision.Wait = true
	opts.Unbind.Wait = true
	opts.Deprovision.Wait = true
}
