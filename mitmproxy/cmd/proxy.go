package main

import (
	"context"
	"fmt"
	"github.com/kataras/golog"
	"github.com/urfave/cli"
	"online/common/log"
	"online/mitmproxy"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

var (
	sigExitOnce = new(sync.Once)
)

func init() {
	go sigExitOnce.Do(func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		defer signal.Stop(c)

		for {
			select {
			case <-c:
				fmt.Printf("exit by signal [SIGTERM/SIGINT/SIGKILL]")
				os.Exit(1)
				return
			}
		}
	})
}

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{}

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port",
			Value: 8088,
		},
	}

	app.Before = func(context *cli.Context) error {
		log.SetLevel(golog.InfoLevel)
		return nil
	}

	app.Action = func(c *cli.Context) error {
		log.Info("start to startup mitmproxy...")
		proxy, err := mitmproxy.NewMITMProxy(mitmproxy.WithAutoCa())
		if err != nil {
			return err
		}

		return proxy.Run(context.Background())
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("command: [%v] failed: %v\n", strings.Join(os.Args, " "), err)
		return
	}
}
