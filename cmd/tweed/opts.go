package main

var opts struct {
	Help    bool `cli:"-h, --help"`
	Version bool `cli:"-v, --version"`
	Quiet   bool `cli:"-q, --quiet"`
	JSON    bool `cli:"--json"`
	Debug   bool `cli:"-D, --debug"`

	LogLevel string `cli:"--log-level", env:"TWEED_LOG_LEVEL"`

	Oops struct{} `cli:"oops"`

	Broker struct {
		Config     string `cli:"-c, --config"   env:"TWEED_CONFIG_FILE"`
		ConfigJSON string `                     env:"TWEED_CONFIG"`

		Listen           string `cli:"-L, --listen"        env:"TWEED_LISTEN"`
		Root             string `cli:"-r, --root"          env:"TWEED_ROOT"`
		HTTPAuthUsername string `cli:"-U, --http-username" env:"TWEED_HTTP_USERNAME"`
		HTTPAuthPassword string `cli:"-P, --http-password" env:"TWEED_HTTP_PASSWORD"`
		HTTPAuthRealm    string `cli:"--http-realm"        env:"TWEED_HTTP_REALM"`
		KeepErrors       int    `cli:"--keep-errors"       env:"TWEED_ERRORS"`
	} `cli:"broker"`

	Username string `cli:"-u, --username"   env:"TWEED_USERNAME"`
	Password string `cli:"-p, --password"   env:"TWEED_PASSWORD"`
	Tweed    string `cli:"-T, --tweed"      env:"TWEED_URL"`

	Catalog   struct{} `cli:"catalog,cat"`
	Instances struct{} `cli:"instances,ls"`
	Bindings  struct{} `cli:"bindings"`
	Binding   struct{} `cli:"binding"`

	Wait struct {
		Instance bool   `cli:"-i, --instance"`
		Task     bool   `cli:"-t, --task"`
		State    string `cli:"-s, --state"`
		Not      bool   `cli:"--not"`
		Sleep    int    `cli:"-S, --sleep"`
		Max      int    `cli:"-m, --max"`
	} `cli:"wait"`

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

	Bind struct {
		ID   string `cli:"--as"`
		Wait bool   `cli:"--wait, --no-wait"`
	} `cli:"bind"`

	Unbind struct {
		Wait bool `cli:"--wait, --no-wait"`
	} `cli:"unbind"`
}

func init() {
	opts.LogLevel = "info"

	opts.Broker.Listen = ":5000"
	opts.Broker.HTTPAuthUsername = "tweed"
	opts.Broker.HTTPAuthPassword = "tweed"
	opts.Broker.HTTPAuthRealm = "Tweed"
	opts.Broker.KeepErrors = 1000

	opts.Wait.Sleep = 2

	opts.Reset.State = "quiet"

	opts.Provision.Wait = true
	opts.Bind.Wait = true
	opts.Unbind.Wait = true
	opts.Deprovision.Wait = true
}
