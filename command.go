package athena

import (
	"fmt"
	"os"
	"time"

	cli "github.com/jawher/mow.cli"
	rscsrv "github.com/lab259/go-rscsrv"
	"github.com/lab259/rlog/v2"
)

type CommandAction func(opt *CommandOptions)

type CommandOptions struct {
	BindAddress string
	Wait        int
	Hostname    string
	isDryRun    bool
}

type commandBuilder struct {
	wait           int
	bindAddress    string
	hostname       string
	serviceStarter rscsrv.ServiceStarter
	action         CommandAction
}

func NewCommand(action CommandAction) *commandBuilder {
	hostname, _ := os.Hostname()
	return &commandBuilder{
		action:      action,
		bindAddress: "127.0.0.1:3000",
		wait:        0,
		hostname:    hostname,
	}
}

func (b *commandBuilder) ServiceStarter(serviceStarter rscsrv.ServiceStarter) *commandBuilder {
	b.serviceStarter = serviceStarter
	return b
}

func (b *commandBuilder) BindAddress(bindAddress string) *commandBuilder {
	b.bindAddress = bindAddress
	return b
}

func (b *commandBuilder) Wait(wait int) *commandBuilder {
	b.wait = wait
	return b
}

func (b *commandBuilder) Hostname(hostname string) *commandBuilder {
	b.hostname = hostname
	return b
}

func (b *commandBuilder) Build() cli.CmdInitializer {
	return cli.CmdInitializer(func(cmd *cli.Cmd) {
		var options CommandOptions

		cmd.StringPtr(&options.BindAddress, cli.StringOpt{
			Name:   "B bind-address",
			Value:  b.bindAddress,
			Desc:   "The bind address will be used on the HTTP server",
			EnvVar: "BIND_ADDR",
		})

		cmd.IntPtr(&options.Wait, cli.IntOpt{
			Name:   "w wait",
			Value:  b.wait,
			Desc:   "Delay in seconds before the initialization",
			EnvVar: "WAIT",
		})

		cmd.StringPtr(&options.Hostname, cli.StringOpt{
			Name:   "H hostname",
			Value:  b.hostname,
			Desc:   "The name of the station running the app instance",
			EnvVar: "HOSTNAME",
		})

		cmd.BoolOptPtr(&options.isDryRun, "d dry-run", false, "Loads the configuration and check if the dependencies are working (such as database connections)")

		cmd.Before = func() {
			os.Setenv("HOSTNAME", options.Hostname)

			if options.isDryRun {
				rlog.Trace(1, "This is a dry run!")
			}

			if options.Wait > 0 {
				rlog.Infof("  Waiting %d seconds before continue ...", options.Wait)
				time.Sleep(time.Duration(options.Wait) * time.Second)
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

			if options.isDryRun {
				if b.serviceStarter != nil {
					err := b.serviceStarter.Stop(false)
					if err != nil {
						rlog.Critical(err)
						os.Exit(2)
					}
				}

				if options.isDryRun {
					rlog.Trace(1, "Everything looks fine!")
				}
				os.Exit(0)
			}
		}

		cmd.Action = func() {
			b.action(&options)
		}
	})

}
