package main

import (
	"fmt"
	"os"

	"github.com/SvenDowideit/gendoc/allprojects"
	"github.com/SvenDowideit/gendoc/commands"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

// Set from the go build commandline
var Version, CommitHash string

type Exit struct {
	Code int
}

func main() {
	// We want our defer functions to be run when calling fatal()
	defer func() {
		if e := recover(); e != nil {
			if ex, ok := e.(Exit); ok == true {
				os.Exit(ex.Code)
			}
			panic(e)
		}
	}()
	app := cli.NewApp()
	app.Name = "gendoc"
	app.Version = Version
	app.Usage = "Generate documentation from multiple GitHub repositories"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		// TODO: add a debug / info log file
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output in the logs",
		},
		cli.StringFlag{
			Name:        "ghtoken",
			Usage:       "GITHUB_TOKEN for git and GitHub API",
			EnvVar:      "GITHUB_TOKEN",
			Destination: &allprojects.GithubToken,
		},
	}
	app.Commands = []cli.Command{
		versionCommand,
		commands.Clone,
		commands.Checkout,
		commands.Install,
		commands.Release,
		commands.Remote,
		commands.Render,
		commands.Status,
	}
	app.Before = func(context *cli.Context) error {
		if context.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

var versionCommand = cli.Command{
	Name:  "version",
	Usage: "return the version",
	Action: func(context *cli.Context) error {
		fmt.Println(context.App.Version)
		fmt.Println(CommitHash)
		return nil
	},
}

func fatal(err string, code int) {
	fmt.Fprintf(os.Stderr, "[ctr] %s\n", err)
	panic(Exit{code})
}
