package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/d7561985/redshift-test/cmd/redshift"
	"github.com/d7561985/redshift-test/cmd/redshift/migrate"
	"github.com/d7561985/tel/v2"
	"github.com/urfave/cli/v2" // imports as package "cli"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	//// Some text we want to compress.
	//original := "bird and frog"
	//
	//// Open a file for writing.
	//f, _ := os.Create("file.gz")
	//
	//// Create gzip writer.
	//w := gzip.NewWriter(f)
	//
	//// Write bytes in compressed form to the file.
	//w.Write([]byte(original))
	//
	//// Close the file.
	//w.Close()
	//
	//fmt.Println("DONE")

	l := tel.NewSimple(tel.DefaultDebugConfig())
	ctx, cancel := context.WithCancel(l.Ctx())
	defer cancel()

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, os.Interrupt)

		<-ch
		l.Info("stop application")
		cancel()
	}()

	pgCli := redshift.New()

	app := &cli.App{
		Name:     pgCli.Name,
		Usage:    pgCli.Usage,
		Action:   pgCli.Action,
		Flags:    pgCli.Flags,
		Commands: []*cli.Command{migrate.New().Cli()},
	}

	err := app.RunContext(ctx, os.Args)
	if err != nil {
		l.Error("execution", tel.Error(err))
	}
}
