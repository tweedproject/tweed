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

	if opts.Version {
		fmt.Printf("tweed %s %s\n", version("v"), build())
		os.Exit(0)
	}

	if opts.Help {
		fmt.Printf("tweed %s %s\n", version("v"), build())
		if command == "" {
			fmt.Printf("USAGE: tweed COMMAND [options]\n\n")
			fmt.Printf("GLOBAL OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" -h, --help                 Show the help screen.\n")
			fmt.Printf(" -v, --version              Print Tweed version and exit.\n")
			fmt.Printf("\n")
			fmt.Printf(" -q, --quiet                Try not to print anything.\n")
			fmt.Printf(" -D, --debug                Enable extra debugging output.\n")
			fmt.Printf("                            (Not too compatible with `--quiet')\n")
			fmt.Printf("\n")
			fmt.Printf("     --json                 If possible, format output in JSON.\n")
			fmt.Printf("\n")
			fmt.Printf(" -T, --tweed URL            The URL of the Tweed Data Services Broker to use.\n")
			fmt.Printf("                            May be set via @W{TWEED_URL}\n")
			fmt.Printf(" -u, --username USERNAME    Username for authenticating to the broker.\n")
			fmt.Printf("                            May be set via @W{TWEED_USERNAME}\n")
			fmt.Printf(" -p, --password PASSWORD    Password for authenticating to the broker.\n")
			fmt.Printf("                            May be set via @W{TWEED_PASSWORD}\n")
			fmt.Printf("\n")
			fmt.Printf("COMMANDS\n")
			fmt.Printf("\n")
			fmt.Printf(" catalog        Print the Tweed Catalog.             (see @W{tweed catalog -h})\n")
			fmt.Printf("\n")
			fmt.Printf(" instances      List known Data Service Instances.   (see @W{tweed instances -h})\n")
			fmt.Printf(" instance       View Data Service Instance details.  (see @W{tweed instance -h})\n")
			fmt.Printf(" bindings       List service instance bindings.      (see @W{tweed bindings -h})\n")
			fmt.Printf(" binding        View service a binding's details.    (see @W{tweed binding -h})\n")
			fmt.Printf("\n")
			fmt.Printf(" wait           Wait for state transitions.          (see @W{tweed wait -h})\n")
			fmt.Printf(" reset          Set the state of an instance.        (see @W{tweed reset -h})\n")
			fmt.Printf(" task           View task details.                   (see @W{tweed task -h})\n")
			fmt.Printf("\n")
			fmt.Printf(" log            View service instance logs.          (see @W{tweed log -h})\n")
			fmt.Printf(" files          List files for a service instance.   (see @W{tweed files -h})\n")
			fmt.Printf(" file           Retrieve a file from an instance.    (see @W{tweed file -h})\n")
			fmt.Printf("\n")
			fmt.Printf(" provision      Deploy a new service instance.       (see @W{tweed provision -h})\n")
			fmt.Printf(" check          Check an instance's viability.       (see @W{tweed check -h})\n")
			fmt.Printf(" bind           Bind credentials for an instance.    (see @W{tweed bind -h})\n")
			fmt.Printf(" unbind         Revoke a service instance binding.   (see @W{tweed unbind -h})\n")
			fmt.Printf(" deprovision    Tear down a service instance.        (see @W{tweed deprovision -h})\n")
			fmt.Printf(" purge          Drop instance logs, tasks, etc.      (see @W{tweed purge -h})\n")
			fmt.Printf("\n")
			fmt.Printf(" broker         Run the Tweed Service Broker.        (see @W{tweed broker -h})\n")
			fmt.Printf(" oops           Dump the broker's error buffer.      (see @W{tweed oops -h})\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
	}

	if command == "broker" {
		if opts.Help {
			fmt.Printf("tweed broker - The Tweed Data Services Broker\n\n")
			fmt.Printf("USAGE: tweed broker --root /path/to/root --listen :5000 [options]\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			//	fmt.Printf(" -h, --help                 Show this help screen.\n")
			//	fmt.Printf(" -v, --version              Print Tweed version and exit.\n")
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
			fmt.Printf(" --keep-errors N            The size of the Tweed broker's error buffer,\n")
			fmt.Printf("                            as a number of errors.  Defaults to 1000.\n")
			fmt.Printf("                            May be set via @W{$TWEED_ERRORS}\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
		Broker(args)
		os.Exit(0)
	}

	if command == "catalog" {
		if opts.Help {
			fmt.Printf("tweed catalog - Display the full Tweed Data Service Offering Catalog.\n\n")
			fmt.Printf("USAGE: tweed catalog [--json]\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" --json       Print the catalog in JSON format, instead of a more\n")
			fmt.Printf("              human-readable format (which is the default)\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
		Catalog(args)
		os.Exit(0)
	}

	if command == "wait" {
		if opts.Help {
			fmt.Printf("tweed wait - Pause until something interesting happens.\n\n")
			fmt.Printf("USAGE: tweed wait -i INSTANCE-ID [-s STATE] [--not] [--sleep N] [--max M]\n")
			fmt.Printf("       tweed wait -t INSTANCE-ID TASK-ID [-s STATE] [--not] [--sleep N] [--max M]\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("SELECTOR\n")
			fmt.Printf(" (you must choose exactly one of these.)\n")
			fmt.Printf("\n")
			fmt.Printf(" -i, --instance   Wait for a state transition on a service instance.\n")
			fmt.Printf("\n")
			fmt.Printf(" -t, --task       Wait for a task to run to completion.\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" -s, --state STATE   What instance state to wait for.\n")
			fmt.Printf("                     Only makes sense in `--instance' mode.\n")
			fmt.Printf("\n")
			fmt.Printf("          STATE must be one of:\n")
			fmt.Printf("\n")
			fmt.Printf("            @W{quiet}           nothing is happening\n")
			fmt.Printf("            @W{provisioning}    new instance is being set up\n")
			fmt.Printf("            @W{binding}         new credentials are being bound\n")
			fmt.Printf("            @W{unbinding}       credentials are being unbound\n")
			fmt.Printf("            @W{deprovisioning}  instance is being torn down\n")
			fmt.Printf("            @W{gone}            instance deprovisioned\n")
			fmt.Printf("\n")
			fmt.Printf(" --not             Negate the --state match predicate.\n")
			fmt.Printf("\n")
			fmt.Printf(" -S, --sleep N     How long to sleep in between checks against the\n")
			fmt.Printf("                   Tweed API.  Defaults to 2 seconds, like `watch'.\n")
			fmt.Printf("\n")
			fmt.Printf(" -m, --max M       Maximum number of seconds to wait for the state transition.\n")
			fmt.Printf("                   By default, there is no upper limit.\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
		Wait(args)
		os.Exit(0)
	}

	if command == "reset" {
		if opts.Help {
			fmt.Printf("tweed reset - Set the state of a Tweed Data Service Instance.\n\n")
			fmt.Printf("USAGE: tweed reset [--state STATE] INSTANCE-ID\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" -s, --state STATE   What state to set for the instance:\n")
			fmt.Printf("\n")
			fmt.Printf("            @W{quiet}           nothing is happening\n")
			fmt.Printf("            @W{provisioning}    new instance is being set up\n")
			fmt.Printf("            @W{binding}         new credentials are being bound\n")
			fmt.Printf("            @W{unbinding}       credentials are being unbound\n")
			fmt.Printf("            @W{deprovisioning}  instance is being torn down\n")
			fmt.Printf("            @W{gone}            instance deprovisioned\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
		Reset(args)
		os.Exit(0)
	}

	if command == "check" {
		if opts.Help {
			fmt.Printf("tweed check - Check the viability of a Tweed Data Service Instance.\n\n")
			fmt.Printf("USAGE: tweed check INSTANCE-ID\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" There aren't really any options.  Perhaps check the global help?\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
		Check(args)
		os.Exit(0)
	}

	if command == "provision" {
		if opts.Help {
			fmt.Printf("tweed provision - Deploy a new Tweed Data Service Instance\n\n")
			fmt.Printf("USAGE: tweed provision SERVICE PLAN [options]\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" --as INSTANCE-ID         Assign an explicit instance ID.\n")
			fmt.Printf("\n")
			fmt.Printf(" -P, --param KEY=VALUE    Specify user-configurable service instance\n")
			fmt.Printf("                          parameters.  These vary depending on catalog.\n")
			fmt.Printf("                          This option can be used more than once.\n")
			fmt.Printf("\n")
			fmt.Printf(" --wait                   Whether or not to wait for the instance to\n")
			fmt.Printf(" --no-wait                finish provisioning and enter the quiet state.\n")
			fmt.Printf("                          (defaults to `--no-wait')\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
		Provision(args)
		os.Exit(0)
	}

	if command == "task" {
		if opts.Help {
			fmt.Printf("tweed task - View Task details.\n\n")
			fmt.Printf("USAGE: tweed task INSTANCE-ID TASK-ID [--wait | --no-wait] [--json]\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" --wait       Whether or not to wait for the task to complete and\n")
			fmt.Printf(" --no-wait    enter the quiet state. (defaults to `--no-wait')\n")
			fmt.Printf("\n")
			fmt.Printf(" --json       Print task details in JSON format, instead of a more\n")
			fmt.Printf("              human-readable format (which is the default)\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
		Task(args)
		os.Exit(0)
	}

	if command == "instances" {
		if opts.Help {
			fmt.Printf("tweed instances - List Tweed Data Service Instances.\n\n")
			fmt.Printf("USAGE: tweed instances [--json]\n")
			fmt.Printf("       tweed ls [--json]\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" --json       Print the instances in JSON format, instead of a more\n")
			fmt.Printf("              human-readable format (which is the default)\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
		Instances(args)
		os.Exit(0)
	}

	if command == "instance" {
		if opts.Help {
			fmt.Printf("tweed instance - View Tweed Data Service Instance details.\n\n")
			fmt.Printf("USAGE: tweed instance INSTANCE-ID [--json]\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" --json       Print instance details in JSON format, instead of a more\n")
			fmt.Printf("              human-readable format (which is the default)\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
		Instance(args)
		os.Exit(0)
	}

	if command == "log" {
		if opts.Help {
			fmt.Printf("tweed log - Review logs for a Tweed Data Service Instance.\n\n")
			fmt.Printf("USAGE: tweed log INSTANCE-ID [--tail LINES] [--follow]\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" -n, --tail LINES    How many lines of 'recent' log messages to display.\n")
			fmt.Printf("                     Defaults to zero (0).\n")
			fmt.Printf("\n")
			fmt.Printf(" -f, --follow        Follow the log, displaying new log messages as they\n")
			fmt.Printf("                     are logged.  Similar to `tail -f'\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
		Log(args)
		os.Exit(0)
	}

	if command == "files" {
		if opts.Help {
			fmt.Printf("tweed files - List files pertinent to a Tweed Data Service Instance.\n\n")
			fmt.Printf("USAGE: tweed files INSTANCE-ID [--json]\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" --json       Print the file list in JSON format, instead of a more\n")
			fmt.Printf("              human-readable format (which is the default)\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
		Files(args)
		os.Exit(0)
	}

	if command == "file" {
		if opts.Help {
			fmt.Printf("tweed file - Retrieve a file pertinent to a Tweed Data Service Instance.\n\n")
			fmt.Printf("USAGE: tweed file INSTANCE-ID FILENAME [--json]\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" --json       Print the file details in JSON format, instead of a more\n")
			fmt.Printf("              human-readable format (which is the default)\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
		File(args)
		os.Exit(0)
	}

	if command == "bind" {
		if opts.Help {
			fmt.Printf("tweed bind - Bind a Tweed Data Service Instance, to get credentials.\n\n")
			fmt.Printf("USAGE: tweed bind INSTANCE-ID [--as BINDING-ID] [--wait | --no-wait]\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" --as BINDING-ID   Assign an explicit ID to the new binding.\n")
			fmt.Printf("\n")
			fmt.Printf(" --wait            Whether or not to wait for the bind to complete and\n")
			fmt.Printf(" --no-wait         enter the quiet state. (defaults to `--no-wait')\n")
			fmt.Printf("\n")
		}
		Bind(args)
		os.Exit(0)
	}

	if command == "bindings" {
		if opts.Help {
			fmt.Printf("tweed bindings - List bindings for a Tweed Data Service Instance.\n\n")
			fmt.Printf("USAGE: tweed bindings INSTANCE-ID [--json]\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" --json       Print the bindings in JSON format, instead of a more\n")
			fmt.Printf("              human-readable format (which is the default)\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
		Bindings(args)
		os.Exit(0)
	}

	if command == "binding" {
		if opts.Help {
			fmt.Printf("tweed binding - View Tweed Data Service Instance Binding details.\n\n")
			fmt.Printf("USAGE: tweed binding INSTANCE-ID BINDING-ID [--json]\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" --json       Print the catalog in JSON format, instead of a more\n")
			fmt.Printf("              human-readable format (which is the default)\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
		Binding(args)
		os.Exit(0)
	}

	if command == "unbind" {
		if opts.Help {
			fmt.Printf("tweed unbind - Remove a Tweed Data Service Instance Binding.\n\n")
			fmt.Printf("USAGE: tweed unbind INSTANCE-ID BINDING-ID [--wait | --no-wait]\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" --wait       Whether or not to wait for the unbind to complete and\n")
			fmt.Printf(" --no-wait    enter the gone state. (defaults to `--no-wait')\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
		Unbind(args)
		os.Exit(0)
	}

	if command == "deprovision" {
		if opts.Help {
			fmt.Printf("tweed deprovision - Tear down a Tweed Data Service Instance.\n\n")
			fmt.Printf("USAGE: tweed deprovision INSTANCE-ID [--wait | --no-wait]\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" --wait       Whether or not to wait for the instance to be fully torn down,\n")
			fmt.Printf(" --no-wait    and enter the gone state. (defaults to `--no-wait')\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
		Deprovision(args)
		os.Exit(0)
	}

	if command == "purge" {
		if opts.Help {
			fmt.Printf("tweed purge - Remove traces of deprovisioned Tweed Data Service Instance.\n\n")
			fmt.Printf("USAGE: tweed purge [--all | INSTANCE-ID]\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" -a, --all      Find all service instances in the 'gone' state, and\n")
			fmt.Printf("                remove their files, logs, etc.\n")
			fmt.Printf("\n")
			fmt.Printf(" INSTANCE-ID    Remove all traces of the identified service instance.\n")
			fmt.Printf("                That service must be in the @W{gone} state.\n")
			fmt.Printf("                Incompatible with the `--all' option.\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
		Purge(args)
		os.Exit(0)
	}

	if command == "oops" {
		if opts.Help {
			fmt.Printf("tweed oops - Dump the Tweed Broker error buffer.\n\n")
			fmt.Printf("USAGE: tweed oops [--json]\n")
			fmt.Printf("       (see `tweed --help' for global options / other commands)\n")
			fmt.Printf("\n")
			fmt.Printf("OPTIONS\n")
			fmt.Printf("\n")
			fmt.Printf(" --json       Print the errors in JSON format, instead of a more\n")
			fmt.Printf("              human-readable format (which is the default)\n")
			fmt.Printf("\n")
			os.Exit(0)
		}
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
