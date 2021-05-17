package bouncer

import (
	"errors"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

var (
	// the arguments handed over to the cli
	args = os.Args
	// function the get the users home directory
	getUserHomeDir = os.UserHomeDir
	// function to enforce the chia executable
	enforceChiaExecutable = enforceExists
	// function to provide file info
	getFileInfo = os.Stat
)

// Context describes the environment of the tool execution
type Context struct {
	// ChiaExecutable the chia executable e.g. /home/steffen/chia-blockchain/venv/bin/chia
	ChiaExecutable string
	// Location is the location to filter for
	Location string
	// IsDownThreshold is true in case it was set by the user
	IsDownThreshold bool
	// DownThreshold is the down speed threshold to filter for
	DownThreshold float64
	// Done indicates that we are done (--help, --version...)
	Done bool
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

// RunCli starts the cli which includes validation of parameters.
// The returned context consists of a chia executable and location to filter for
func RunCli() (*Context, error) {
	var chiaExecutable string
	var location string
	var done bool
	var isDownThreshold bool
	var downThreshold float64

	cli.HelpFlag = &cli.BoolFlag{
		Name:        "help",
		Aliases:     []string{"h"},
		Usage:       "show help",
		Destination: &done,
	}

	app := &cli.App{
		Name:      "chia-bouncer",
		Usage:     "remove unwanted connections from your Chia Node based on Geo IP Location.",
		UsageText: "chia-bouncer [-e CHIA-EXECUTABLE] [-d DOWN-THRESHOLD] LOCATION\n\t chia-bouncer -e /chia-blockchain/venv/bin/chia -d 0.2 mars",
		ArgsUsage: "LOCATION",
		Description: "Tool will lookup connections via 'chia show -c', get ip locations via geoiplookup and " +
			"remove nodes from specified LOCATION via 'chia show -r' ",
		EnableBashCompletion: true,
		HideHelpCommand:      true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "chia-exec",
				Aliases:     []string{"e"},
				Required:    false,
				DefaultText: "$HOME/chia-blockchain/venv/bin/chia",
				Usage:       "`CHIA-EXECUTABLE`. normally located inside the bin folder of your venv directory",
				Destination: &chiaExecutable,
			},
			&cli.Float64Flag{
				Name:        "down-threshold",
				Aliases:     []string{"d"},
				Required:    false,
				DefaultText: "not active",
				Usage:       "`DOWN-THRESHOLD` defines the additional filter for minimal down speed in MiB for filtering.",
				Destination: &downThreshold,
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return errors.New("LOCATION is missing")
			}
			if c.IsSet("down-threshold") {
				isDownThreshold = true
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

			location = strings.TrimSpace(strings.Join(c.Args().Slice(), " "))
			return nil
		},
		Copyright: "GNU GPLv3",
	}

	err := app.Run(args)
	if err != nil {
		return nil, err
	}

	return &Context{
		ChiaExecutable:  chiaExecutable,
		Location:        location,
		IsDownThreshold: isDownThreshold,
		DownThreshold:   downThreshold,
		Done:            done,
	}, nil
}
