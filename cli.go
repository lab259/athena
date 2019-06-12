package athena

import (
	"fmt"
	"os"
	"time"

	cli "github.com/jawher/mow.cli"
	rscsrv "github.com/lab259/go-rscsrv"
	"github.com/lab259/rlog"
)

type CLIOptions struct {
	BindAddress *string
	Wait        *int
	HostName    *string
	isDryRun    *bool
	isHelp      *bool
	isVersion   *bool
}

type cliBuilder struct {
	cli            *cli.Cli
	env            string
	version        string
	build          string
	serviceStarter *rscsrv.ServiceStarter
}

func NewCLI(name, description string) *cliBuilder {
	return &cliBuilder{
		cli: cli.App(name, description),
	}
}

func (b *cliBuilder) Version(version, build string) *cliBuilder {
	b.version = version
	b.build = build
	return b
}

func (b *cliBuilder) Environment(env string) *cliBuilder {
	b.env = env
	return b
}

func (b *cliBuilder) ServiceStarter(serviceStarter *rscsrv.ServiceStarter) *cliBuilder {
	b.serviceStarter = serviceStarter
	return b

}

func (b *cliBuilder) Build() (*cli.Cli, *CLIOptions) {
	var options CLIOptions

	options.BindAddress = b.cli.String(cli.StringOpt{
		Name:   "B bind-address",
		Value:  "127.0.0.1:3000",
		Desc:   "The bind address will be used on the HTTP server",
		EnvVar: "BIND_ADDR",
	})

	options.Wait = b.cli.Int(cli.IntOpt{
		Name:   "w wait",
		Value:  0,
		Desc:   "Delay in seconds before the initialization",
		EnvVar: "WAIT",
	})

	hostName, _ := os.Hostname()
	options.HostName = b.cli.String(cli.StringOpt{
		Name:   "H hostname",
		Value:  hostName,
		Desc:   "The name of the station running the app instance",
		EnvVar: "HOSTNAME",
	})
	os.Setenv("HOSTNAME", *options.HostName)

	options.isDryRun = b.cli.BoolOpt("d dry-run", false, "Loads the configuration and check if the dependencies are working (such as database connections)")
	options.isHelp = b.cli.BoolOpt("h help", false, "Displays this help message")
	options.isVersion = b.cli.BoolOpt("v version", false, "Displays the version")

	b.cli.Before = func() {
		if *options.isHelp {
			b.cli.PrintLongHelp()
			os.Exit(0)
		}

		if *options.isVersion {
			fmt.Printf("Version: %s\n", b.version)
			if b.build != "" {
				fmt.Printf("  Build: %s\n", b.build)
			}
			os.Exit(0)
		}

		if *options.isDryRun {
			rlog.Info("Dry run!")
		}

		rlog.Infof("Version: %s (%s)", b.version, b.build)
		rlog.Infof("Loading configuration from '%s' ...", b.env)

		if *options.Wait > 0 {
			rlog.Infof("  Waiting %d seconds before continue ...", *options.Wait)
			time.Sleep(time.Duration(*options.Wait) * time.Second)
			rlog.Info(fmt.Sprintf("    > Waiting %s", "DONE"))
		}

		if b.serviceStarter != nil {
			err := b.serviceStarter.Start()
			if err != nil {
				b.serviceStarter.Stop(true)
				rlog.Critical(err)
				os.Exit(2)
			}
		}

		if *options.isDryRun {
			if b.serviceStarter != nil {
				err := b.serviceStarter.Stop(false)
				if err != nil {
					rlog.Critical(err)
					os.Exit(2)
				}
			}
			os.Exit(0)
		}
	}

	return b.cli, &options
}
