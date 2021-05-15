package bouncer

import (
	"errors"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

var (
	// the version of this tool
	version = "unknown"
	// the arguments handed over to the cli
	args    = os.Args
	// function the get the users home directory
	getUserHomeDir        = os.UserHomeDir
	// function to enforce the chia executable
	enforceChiaExecutable = enforceExists
	// function to provide file info
	getFileInfo           = os.Stat
)

// Context describes the environment of the tool execution
type Context struct {
	chiaExecutable string
	location       string
}

const DefaultChiaExecutableSuffix = "chia-blockchain/venv/bin/chia"

// defaultChiaExecutable to get the default chia executable from the home directory of the current user
func defaultChiaExecutable() (string, error) {
	dirname, err := getUserHomeDir()
	if err != nil {
		return "", err
	}
	return dirname + "/" + DefaultChiaExecutableSuffix, nil
}

// enforceExists enforces that the chia executable can be used
func enforceExists(chiaExecutable string) error {
	info, err := getFileInfo(chiaExecutable)
	if os.IsNotExist(err) {
		return errors.New("chia executable does not exist")
	}
	if info.IsDir() {
		return errors.New("chia executable can not be a directory")
	}

	// TODO could add check if file is executable for current user

	return nil
}

// Run starts the cli which includes validation of parameters.
// the returned context consists of chia executable and location to filter for
func Run() (*Context, error) {
	var chiaExecutable string
	var location string

	app := &cli.App{
		Name:      "chia-bouncer",
		Usage:     "remove nodes by given location from your connections",
		UsageText: "chia-bouncer -ce /home/steffen/chia-blockchain/venv/bin/chia mars",
		ArgsUsage: "LOCATION",
		Description: "Tool will lookup connections via 'chia show -c', get ip locations via geoiplookup and " +
			"remove nodes from specified location via 'chia show -r' ",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "chiaexec",
				Aliases:     []string{"e"},
				Required:    false,
				DefaultText: "$HOME/chia-blockchain/venv/bin/chia",
				Usage:       "`CHIA-EXECUTABLE`. normally located inside the bin folder of your venv directory",
				Destination: &chiaExecutable,
			},
			cli.VersionFlag,
		},
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return errors.New("LOCATION is missing")
			}
			if chiaExecutable == "" {
				defaultExecutable, err := defaultChiaExecutable()
				if err != nil {
					return err
				}
				chiaExecutable = defaultExecutable
			}
			if err := enforceChiaExecutable(chiaExecutable); err != nil {
				return err
			}

			location = strings.Join(c.Args().Slice(), " ")
			return nil
		},
		Authors: []*cli.Author{
			{
				Name:  "st3ffn",
				Email: "funk.up.up@gmail.com",
			},
		},
		Copyright: "GNU GPLv3",
		Version:   version,
	}

	err := app.Run(args)
	if err != nil {
		return nil, err
	}

	return &Context{
		chiaExecutable: chiaExecutable,
		location:       location,
	}, nil
}
