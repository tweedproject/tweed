package main

import (
	"log"
	"net/http"
	"os"

	fmt "github.com/jhunt/go-ansi"

	"github.com/tweedproject/tweed"
	"github.com/tweedproject/tweed/stencil"

	"code.cloudfoundry.org/lager"
	"github.com/jessevdk/go-flags"

	"github.com/tweedproject/tweed/creds"
	_ "github.com/tweedproject/tweed/creds/kubernetes"
	_ "github.com/tweedproject/tweed/creds/vault"
)

type BrokerCommand struct {
	Config     string `short:"c" long:"config"  env:"TWEED_CONFIG_FILE" description:"Location of service config file"`
	ConfigJSON string `long:"inline-config"     env:"TWEED_CONFIG" description:"Inline service config in JSON format"`

	Listen           string `short:"L" long:"listen"        env:"TWEED_LISTEN" description:"TCP port to listen on" default:":5000"`
	Root             string `short:"r" long:"root"          env:"TWEED_ROOT"   description:"Location of root infrastructure config files" required:"true"`
	HTTPAuthUsername string `short:"U" long:"http-username" env:"TWEED_HTTP_USERNAME" description:"Username to use for HTTP Basic Auth of the API" default:"tweed"`
	HTTPAuthPassword string `short:"P" long:"http-password" env:"TWEED_HTTP_PASSWORD" description:"Password to use for HTTP Basic Auth of the API" default:"tweed"`
	HTTPAuthRealm    string `long:"http-realm"              env:"TWEED_HTTP_REALM"    description:"Realm name to use for HTTP Basic Auth of the API" default:"Tweed"`
	KeepErrors       int    `long:"keep-errors"             env:"TWEED_ERRORS"        default:"1000"`

	CredentialManagement creds.CredentialManagementConfig `group:"Credential Management"`
	CredentialManagers   creds.Managers
}

func (cmd *BrokerCommand) Execute(args []string) error {
	SetupLogging()

	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "ERROR: extra arguments found in invocation.\n")
		fmt.Fprintf(os.Stderr, "tweed service broker SHUTTING DOWN.\n")
		os.Exit(1)
	}

	ok := true
	if cmd.Listen == "" {
		fmt.Fprintf(os.Stderr, "@R{(error)} No @R{--listen} flag given, and @W{$TWEED_LISTEN} not set.\n")
		ok = false
	}
	if cmd.Config == "" && cmd.ConfigJSON == "" {
		fmt.Fprintf(os.Stderr, "@R{(error)} No @R{--config} flag given, and neither @W{$TWEED_CONFIG_FILE} nor @W{$TWEED_CONFIG} were set.\n")
		ok = false
	}
	if !ok {
		os.Exit(1)
	}

	logger := log.New(log.Writer(), "", log.LstdFlags)

	credsLogger := lager.NewLogger("broker")
	credsLogger.RegisterSink(lager.NewWriterSink(log.Writer(), lager.INFO))

	secretManager, err := cmd.secretManager(credsLogger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{(error)} failed to configure the Credential Manager:\n")
		fmt.Fprintf(os.Stderr, "        @R{%s}\n", err)
		os.Exit(1)
	}

	stencilFactory := stencil.NewFactory(cmd.Root, logger)
	core := tweed.Core{
		Root:             cmd.Root,
		HTTPAuthUsername: cmd.HTTPAuthUsername,
		HTTPAuthPassword: cmd.HTTPAuthPassword,
		HTTPAuthRealm:    cmd.HTTPAuthRealm,
		StencilFactory:   stencilFactory,
		SecretManager:    secretManager,
	}
	if cmd.ConfigJSON != "" {
		cmd.Config = "{{json literal from environment}}"
		c, err := tweed.ParseConfig([]byte(cmd.ConfigJSON))
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{(error)} config JSON (from @W{$TWEED_CONFIG}) was invalid:\n")
			fmt.Fprintf(os.Stderr, "        @R{%s}\n", err)
			os.Exit(1)
		}
		core.Config = c

	} else {
		c, err := tweed.ReadConfig(cmd.Config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{(error)} failed to read config from @Y{%s}:\n", cmd.Config)
			fmt.Fprintf(os.Stderr, "        @R{%s}\n", err)
			os.Exit(1)
		}
		core.Config = c
	}

	if err := tweed.ValidInstancePrefix(core.Config.Prefix); err != nil {
		fmt.Fprintf(os.Stderr, "@R{(error)} failed to configure the broker:\n")
		fmt.Fprintf(os.Stderr, "        @R{%s}\n", err)
		os.Exit(1)
	}

	if err := core.SetupVault(); err != nil {
		fmt.Fprintf(os.Stderr, "@R{(error)} failed to configure safe / vault:\n")
		fmt.Fprintf(os.Stderr, "        @R{%s}\n", err)
		os.Exit(1)
	}

	if err := core.SetupInfrastructures(); err != nil {
		fmt.Fprintf(os.Stderr, "@R{(error)} failed to configure infrastructure(s):\n")
		fmt.Fprintf(os.Stderr, "        @R{%s}\n", err)
		os.Exit(1)
	}

	if err := core.ValidateCatalog(); err != nil {
		fmt.Fprintf(os.Stderr, "@R{(error)} failed to validate catalog:\n")
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "loading catalog stencil images ...\n")
	core.LoadCatalogStencils()

	if err := core.Scan(); err != nil {
		fmt.Fprintf(os.Stderr, "@R{(error)} failed to detect pre-existing service instances / bindings:\n")
		fmt.Fprintf(os.Stderr, "        @R{%s}\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, `

     |  ###  |   |  ###  |     ######## ##      ## ######## ######## ########
    -----------------------       ##    ##  ##  ## ##       ##       ##     ##
     |  ###  |   |  ###  |        ##    ##  ##  ## ##       ##       ##     ##
    --- ### ------- ### ---       ##    ##  ##  ## ######   ######   ##     ##
    #######################       ##    ##  ##  ## ##       ##       ##     ##
    #######################       ##    ##  ##  ## ##       ##       ##     ##
    --- ### ------- ### ---       ##     ###  ###  ######## ######## ########
     |  ###  |   |  ###  |
    -----------------------    @M{Tweed Service Broker}
     |  ###  |   |  ###  |     @K{v}@G{`+version("")+`} @K{`+build()+`}

    config  :: @W{%s}
    root    :: @W{%s}
    binding :: @G{%s}
    prefix  :: @C{%s}
    vault   :: @C{%s}

`, cmd.Config, cmd.Root, cmd.Listen, core.Config.Prefix, core.Config.Vault.Prefix)

	fmt.Fprintf(os.Stderr, "waiting for vault to open up for business...\n")
	core.WaitForVault()

	core.KeepErrors(cmd.KeepErrors)

	fmt.Fprintf(os.Stderr, "tweed broker API spinning up...\n")
	http.Handle("/b/", core.API())
	http.ListenAndServe(cmd.Listen, nil)

	return nil
}

func (cmd *BrokerCommand) WireDynamicFlags(commandFlags *flags.Command) {
	var credsGroup *flags.Group
	groups := commandFlags.Groups()

	for i := 0; i < len(groups); i++ {
		group := groups[i]

		if credsGroup == nil && group.ShortDescription == "Credential Management" {
			credsGroup = group
		}

		groups = append(groups, group.Groups()...)
	}

	if credsGroup == nil {
		panic("could not find Credential Management group for registering managers")
	}

	managerConfigs := make(creds.Managers)
	for name, p := range creds.ManagerFactories() {
		managerConfigs[name] = p.AddConfig(credsGroup)
	}
	cmd.CredentialManagers = managerConfigs

}

func (cmd *BrokerCommand) secretManager(logger lager.Logger) (creds.Secrets, error) {
	var secretsFactory creds.SecretsFactory
	for name, manager := range cmd.CredentialManagers {
		if !manager.IsConfigured() {
			continue
		}

		credsLogger := logger.Session("credential-manager", lager.Data{
			"name": name,
		})

		credsLogger.Info("configured credentials manager")

		err := manager.Init(credsLogger)
		if err != nil {
			return nil, err
		}

		err = manager.Validate()
		if err != nil {
			return nil, fmt.Errorf("credential manager '%s' misconfigured: %s", name, err)
		}

		secretsFactory, err = manager.NewSecretsFactory(credsLogger)
		if err != nil {
			return nil, err
		}

		break
	}

	if secretsFactory == nil {
		return nil, fmt.Errorf("Missing configured Credential Manager")
	}

	result := secretsFactory.NewSecrets()
	result = creds.NewRetryableSecrets(result, cmd.CredentialManagement.RetryConfig)
	return result, nil
}
