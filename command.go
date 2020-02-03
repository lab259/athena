package athena

import (
	"context"
	"fmt"
	"os"
	"os/signal"
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
	IsDryRun    bool
}

type commandBuilder struct {
	wait     int
	hostname string
	services []rscsrv.Service
}

func NewCommand(services ...rscsrv.Service) *commandBuilder {
	hostname, _ := os.Hostname()
	return &commandBuilder{
		services: services,
		wait:     0,
		hostname: hostname,
	}
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
	serviceStarter := rscsrv.DefaultServiceStarter(b.services...)

	return cli.CmdInitializer(func(cmd *cli.Cmd) {
		var options CommandOptions

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

		cmd.BoolOptPtr(&options.IsDryRun, "d dry-run", false, "Loads the configuration and check if the dependencies are working (such as database connections)")

		cmd.Before = func() {
			os.Setenv("HOSTNAME", options.Hostname)

			if options.IsDryRun {
				rlog.Trace(1, "This is a dry run!")
			}

			if options.Wait > 0 {
				rlog.Infof("  Waiting %d seconds before continue ...", options.Wait)
				time.Sleep(time.Duration(options.Wait) * time.Second)
				rlog.Info(fmt.Sprintf("    > Waiting %s", "DONE"))
			}

			if options.IsDryRun {
				// TODO(felipemfp): figure out how to use dry run (if it is useful? maybe go thru services and try to .Load and .ApplyConfiguration?)
				if options.IsDryRun {
					rlog.Trace(1, "Everything looks fine!")
				}
				os.Exit(0)
			}
		}

		cmd.Action = func() {
			var exitCode int

			signals := make(chan os.Signal, 1)
			signal.Notify(signals, os.Interrupt)

			go func() {
				<-signals
				serviceStarter.Stop(true)
			}()

			if err := serviceStarter.Start(); err != context.Canceled {
				exitCode = 2
				serviceStarter.Stop(true)
			}

			serviceStarter.Wait()
			os.Exit(exitCode)
		}
	})

}
