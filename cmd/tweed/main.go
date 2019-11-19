package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/go-cli"
	env "github.com/jhunt/go-envirotron"
	"github.com/jhunt/go-log"
)

var (
	Version     string
	BuildNumber string
)

func version(prefix string) string {
	if Version == "" {
		return "(development)"
	} else {
		return prefix + Version
	}
}

func build() string {
	return BuildNumber
}

func main() {
	env.Override(&opts)
	command, args, err := cli.Parse(&opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "!! @R{%s}\n", err)
		os.Exit(1)
	}

	log.SetupLogging(log.LogConfig{
		Type:  "console",
		Level: opts.LogLevel,
	})

	if len(args) == 1 && args[0] == "help" {
		args = args[1:]
		opts.Help = true
	}
	if opts.Help {
		fmt.Printf("tweed %s %s\n", version("v"), build())
		if command == "broker" {
			fmt.Printf("USAGE: tweed broker --root /path/to/root --listen :5000 [options]\n\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf(" -h, --help                 Show this help screen.\n")
			fmt.Printf(" -v, --version              Print Tweed version and exit.\n")
			fmt.Printf("\n")
			fmt.Printf(" -L, --listen @Y{:PORT}         What TCP port to listen on, and broker API on.\n")
			fmt.Printf("                            May be set via @W{$TWEED_LISTEN}\n")
			fmt.Printf("\n")
			fmt.Printf(" -r, --root @Y{/PATH}           Where can the root infrastructure / stencil\n")
			fmt.Printf("                            configuration files be found?\n")
			fmt.Printf("                            May be set via @W{TWEED_ROOT}\n")
			fmt.Printf("\n")
			fmt.Printf(" -c, --config @Y{/PATH}         Where can the service config configuration files\n")
			fmt.Printf("                            be found?  May be set via @W{TWEED_CATALOG_FILE}\n")
			fmt.Printf("\n")
			fmt.Printf(" -U @K{USERNAME}                Username to use for HTTP Basic Authentication.\n")
			fmt.Printf(" --http-username @K{USERNAME}   Defaults to @C{tweed}.\n")
			fmt.Printf("                            May be set via @W{$TWEED_HTTP_USERNAME}\n")
			fmt.Printf("\n")
			fmt.Printf(" -P @K{PASSWORD}                Username to use for HTTP Basic Authentication.\n")
			fmt.Printf(" --http-password @K{PASSSWORD}  Defaults to @C{tweed}.\n")
			fmt.Printf("                            May be set via @W{$TWEED_HTTP_PASSWORD}\n")
			fmt.Printf("\n")
			fmt.Printf(" --http-realm @K{REALM}         Realm name to use for HTTP Basic Authentication.\n")
			fmt.Printf("                            May be set via @W{$TWEED_HTTP_REALM}\n")
			fmt.Printf("\n")
			os.Exit(0)

		} else {
			fmt.Printf("USAGE: tweed COMMAND [options]\n\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf(" -h, --help                 Show this help screen.\n")
			fmt.Printf(" -v, --version              Print Tweed version and exit.\n")
			fmt.Printf("\n")
			fmt.Printf("COMMANDS\n")
			fmt.Printf(" broker                     Run the Tweed Service Broker.\n")
			fmt.Printf("                            (see @W{tweed broker -h} for details)\n")
			fmt.Printf("\n")
			fmt.Printf(" catalog                    Print the Tweed Catalog\n")
			fmt.Printf("                            (see @W{tweed catalog -h} for details)\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
	}
	if opts.Version {
		fmt.Printf("tweed %s %s\n", version("v"), build())
		os.Exit(0)
	}

	if command == "broker" {
		Broker(args)
		os.Exit(0)
	}

	if command == "catalog" {
		Catalog(args)
		os.Exit(0)
	}

	if command == "wait" {
		Wait(args)
		os.Exit(0)
	}

	if command == "reset" {
		Reset(args)
		os.Exit(0)
	}

	if command == "check" {
		Check(args)
		os.Exit(0)
	}

	if command == "provision" {
		Provision(args)
		os.Exit(0)
	}

	if command == "task" {
		Task(args)
		os.Exit(0)
	}

	if command == "instances" {
		Instances(args)
		os.Exit(0)
	}

	if command == "instance" {
		Instance(args)
		os.Exit(0)
	}

	if command == "log" {
		Log(args)
		os.Exit(0)
	}

	if command == "files" {
		Files(args)
		os.Exit(0)
	}

	if command == "file" {
		File(args)
		os.Exit(0)
	}

	if command == "bind" {
		Bind(args)
		os.Exit(0)
	}

	if command == "bindings" {
		Bindings(args)
		os.Exit(0)
	}

	if command == "binding" {
		Binding(args)
		os.Exit(0)
	}

	if command == "unbind" {
		Unbind(args)
		os.Exit(0)
	}

	if command == "deprovision" {
		Deprovision(args)
		os.Exit(0)
	}

	if command == "purge" {
		Purge(args)
		os.Exit(0)
	}

	if command == "oops" {
		Oops(args)
		os.Exit(0)
	}

	if command != "" {
		fmt.Fprintf(os.Stderr, "@R{(error)} unrecognized command: '%s'\n", command)
	} else if len(args) > 0 {
		fmt.Fprintf(os.Stderr, "@R{(error)} unrecognized command: '%s'\n", args[0])
	} else {
		fmt.Fprintf(os.Stderr, "@R{(error)} no command specified.\n")
	}
	os.Exit(3)

}
