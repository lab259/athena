package athena

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lab259/athena/config"

	cli "github.com/jawher/mow.cli"
	rscsrv "github.com/lab259/go-rscsrv"
	"github.com/lab259/rlog/v2"
)

type CLIOptions struct {
	Version     string
	Build       string
	Environment string

	BindAddress *string
	Wait        *int
	Hostname    *string
	IsDryRun    *bool
}

type cliBuilder struct {
	cli            *cli.Cli
	env            string
	version        string
	build          string
	wait           int
	bindAddress    string
	hostname       string
	serviceStarter rscsrv.ServiceStarter
	simple         bool
}

func NewCLI(name, description string) *cliBuilder {
	hostname, _ := os.Hostname()
	return &cliBuilder{
		cli:         cli.App(name, description),
		env:         config.Environment(),
		bindAddress: "127.0.0.1:3000",
		wait:        0,
		hostname:    hostname,
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

func (b *cliBuilder) ServiceStarter(serviceStarter rscsrv.ServiceStarter) *cliBuilder {
	b.serviceStarter = serviceStarter
	return b
}

func (b *cliBuilder) BindAddress(bindAddress string) *cliBuilder {
	b.bindAddress = bindAddress
	return b
}

func (b *cliBuilder) Wait(wait int) *cliBuilder {
	b.wait = wait
	return b
}

func (b *cliBuilder) Hostname(hostname string) *cliBuilder {
	b.hostname = hostname
	return b
}

// Simple will disable all flags and the bootstrap Before.
func (b *cliBuilder) Simple() *cliBuilder {
	b.simple = true
	return b
}

func (b *cliBuilder) Build() (*cli.Cli, *CLIOptions) {
	var options CLIOptions

	options.Version = b.version
	options.Build = b.build
	options.Environment = b.env

	var version strings.Builder
	version.WriteString(fmt.Sprintf("Version: %s", b.version))
	if b.build != "" {
		version.WriteString(fmt.Sprintf("\n  Build: %s", b.build))
	}
	b.cli.Version("v version", version.String())

	if !b.simple {
		options.BindAddress = b.cli.String(cli.StringOpt{
			Name:   "B bind-address",
			Value:  b.bindAddress,
			Desc:   "The bind address will be used on the HTTP server",
			EnvVar: "BIND_ADDR",
		})

		options.Wait = b.cli.Int(cli.IntOpt{
			Name:   "w wait",
			Value:  b.wait,
			Desc:   "Delay in seconds before the initialization",
			EnvVar: "WAIT",
		})

		options.Hostname = b.cli.String(cli.StringOpt{
			Name:   "H hostname",
			Value:  b.hostname,
			Desc:   "The name of the station running the app instance",
			EnvVar: "HOSTNAME",
		})

		options.IsDryRun = b.cli.BoolOpt("d dry-run", false, "Loads the configuration and check if the dependencies are working (such as database connections)")

		b.cli.Before = func() {
			os.Setenv("HOSTNAME", *options.Hostname)

			if *options.IsDryRun {
				rlog.Trace(1, "This is a dry run!")
			}

			rlog.Infof("Version: %s (%s)", b.version, b.build)
			rlog.Infof("Environment: %s", b.env)

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

			if *options.IsDryRun {
				if b.serviceStarter != nil {
					err := b.serviceStarter.Stop(false)
					if err != nil {
						rlog.Critical(err)
						os.Exit(2)
					}
				}

				if *options.IsDryRun {
					rlog.Trace(1, "Everything looks fine!")
				}
				os.Exit(0)
			}
		}
	}

	return b.cli, &options
}
