package main

import (
	"os"

	"github.com/codegangsta/cli"
)

const version = "0.0.1"

var (
	hostname string
)

func main() {

	var err error
	hostname, err = os.Hostname()
	assertExit("Could not get the hostname", err, 3)

	app := cli.NewApp()
	app.Author = "sx team @ bq"
	app.Email = "sx@bq.com"
	app.Name = "fwrules"
	app.Usage = "fwrules maintains firewall configuration for a group of machines"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Value:  "config.toml",
			Usage:  "configuration file path",
			EnvVar: "CONFIG_FILE",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "service",
			Usage:  "runs the service that configures firewall on demand",
			Action: fwrulesService,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "template",
					Usage:  "firewall template for the machine",
					EnvVar: "TEMPLATE",
				},
			},
		},
		{
			Name:   "test",
			Usage:  "show the network interfaces of the current machine",
			Action: fwrulesTest,
		},
		{
			Name:   "interfaces",
			Usage:  "show the network interfaces of the current machine",
			Action: fwrulesInterfaces,
		},
		{
			Name:  "kv",
			Usage: "key/value",
			Subcommands: []cli.Command{
				{
					Name:   "set",
					Usage:  "sets a value",
					Action: kvSet,
				},
				{
					Name:   "get",
					Usage:  "gets a value",
					Action: kvGet,
				},
				{
					Name:   "delete",
					Usage:  "deletes a value",
					Action: kvDelete,
				},
			},
		},
	}
	app.Run(os.Args)
}
