package main

import (
	fmt "github.com/jhunt/go-ansi"
	"net/http"
	"os"

	"github.com/tweedproject/tweed"
)

func Broker(args []string) {
	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "ERROR: extra arguments found in invocation.\n")
		fmt.Fprintf(os.Stderr, "tweed service broker SHUTTING DOWN.\n")
		os.Exit(1)
	}

	ok := true
	if opts.Broker.Listen == "" {
		fmt.Fprintf(os.Stderr, "@R{(error)} No @R{--listen} flag given, and @W{$TWEED_LISTEN} not set.\n")
		ok = false
	}
	if opts.Broker.Root == "" {
		fmt.Fprintf(os.Stderr, "@R{(error)} No @R{--root} flag given, and @W{$TWEED_ROOT} not set.\n")
		ok = false
	}
	if opts.Broker.Config == "" && opts.Broker.ConfigJSON == "" {
		fmt.Fprintf(os.Stderr, "@R{(error)} No @R{--config} flag given, and neither @W{$TWEED_CONFIG_FILE} nor @W{$TWEED_CONFIG} were set.\n")
		ok = false
	}
	if !ok {
		os.Exit(1)
	}

	core := tweed.Core{
		Root:             opts.Broker.Root,
		HTTPAuthUsername: opts.Broker.HTTPAuthUsername,
		HTTPAuthPassword: opts.Broker.HTTPAuthPassword,
		HTTPAuthRealm:    opts.Broker.HTTPAuthRealm,
	}
	if opts.Broker.ConfigJSON != "" {
		opts.Broker.Config = "{{json literal from environment}}"
		c, err := tweed.ParseConfig([]byte(opts.Broker.ConfigJSON))
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{(error)} config JSON (from @W{$TWEED_CONFIG}) was invalid:\n")
			fmt.Fprintf(os.Stderr, "        @R{%s}\n", err)
			os.Exit(1)
		}
		core.Config = c

	} else {
		c, err := tweed.ReadConfig(opts.Broker.Config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{(error)} failed to read config from @Y{%s}:\n", opts.Broker.Config)
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

`, opts.Broker.Config, opts.Broker.Root, opts.Broker.Listen, core.Config.Prefix, core.Config.Vault.Prefix)

	fmt.Fprintf(os.Stderr, "waiting for vault to open up for business...\n")
	core.WaitForVault()

	core.KeepErrors(opts.Broker.KeepErrors)

	fmt.Fprintf(os.Stderr, "tweed broker API spinning up...\n")
	http.Handle("/b/", core.API())
	http.ListenAndServe(opts.Broker.Listen, nil)
}
