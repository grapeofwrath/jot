package cmd

import (
	"flag"
	"os"
	"os/exec"
	"slices"

	"github.com/grapeofwrath/jot/helper"
)

type Command struct {
	Name string
	Help string
	Run  func(args []string, c helper.Config, template string) error
}

var commands = []Command{
	{
		Name: "config",
		Help: "Print the current configuration",
		Run:  Config,
	},
	{
		Name: "down",
		Help: "Create a new note in your vault.",
		Run:  Down,
	},
	{
		Name: "help",
		Help: "Print this help message.",
		Run:  Help,
	},
}

func Execute() error {
	config, template, err := helper.MkConfig()
	helper.Check(err)

	flag.Usage = Usage
	flag.Parse()

	if len(flag.Args()) < 1 {
		root(config)
		os.Exit(0)
	}

	subCmd := flag.Arg(0)
	subCmdArgs := flag.Args()[1:]

	runCommand(subCmd, subCmdArgs, config, template)

	return nil
}

func runCommand(name string, args []string, c helper.Config, template string) {
	commandIDX := slices.IndexFunc(commands, func(cmd Command) bool {
		return cmd.Name == name
	})

	if commandIDX < 0 {
		helper.Logger.Errorf("Command \"%s\" not found\n\n", name)
		flag.Usage()
		os.Exit(1)
	}

	if err := commands[commandIDX].Run(args, c, template); err != nil {
		helper.Logger.Error(err)
		os.Exit(1)
	}
}

func root(c helper.Config) error {
	// TODO
	// open fzf to search vault and pick file to open
	err := os.Chdir(c.Vault.Directory)
	helper.Check(err)

	cmd := exec.Command(c.Editor)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	helper.Check(err)

	return nil
}
