package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"online/common/log"
	"online/server/api"
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
			Value: 80,
		},
		cli.StringFlag{
			Name: "fe",
		},
	}

	app.Before = func(context *cli.Context) error {
		level, err := log.ParseLevel("info")
		if err != nil {
			return err
		}
		log.SetLevel(level)
		log.SetOutput(os.Stdout)
		return nil
	}

	app.Action = func(c *cli.Context) error {
		println("------------Yakit Online Banner-----------")
		pg := api.GeneratePostgresParams("127.0.0.1", 5432, "root", "password")
		println("------------Yakit Online Start------------")
		log.Info("start to run server")
		err := api.StartServer(
			pg,
			c.Int("port"),  // web port
			c.String("fe"), // frontend dir
		)
		if err != nil {
			log.Errorf("start yaklang.online service failed: %v", err)
			return err
		}

		ctx := context.Background()
		select {
		case <-ctx.Done():
		}

		//grpcTrans := grpc.NewServer(
		//	grpc.MaxRecvMsgSize(100*1024*1024),
		//	grpc.MaxSendMsgSize(100*1024*1024),
		//)
		//pb.RegisterOnlineServer(grpcTrans, nil)

		return errors.New("server finished")
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("command: [%v] failed: %v\n", strings.Join(os.Args, " "), err)
		return
	}
}
