package cmd

import (
	"flag"
	"os"

	"github.com/grapeofwrath/jot/helper"
	"github.com/spf13/viper"
)

func Help(_ []string, _ helper.Config, _ string) error {
	flag.Usage()
	return nil
}

func Usage() {
	helper.Logger.Print(`jot is a CLI note helper.

Although this is primarily designed for a zettelkasten with Obsidian, Jot can also
be used with your own templates.

Run 'jot' to open your $EDITOR in the vault directory defined in its config.

Usage:

    jot [command] [flags] [args...]`)
	helper.Logger.Print("\nCommands:\n")
	for _, cmd := range commands {
		helper.Logger.Printf("    %-8s %s", cmd.Name, cmd.Help)
	}
	helper.Logger.Print(`
Examples:

    jot down wo ist mein kaffee
        'wo-ist-mein-kaffee.md' will be created in your vault and opened with $EDITOR.

    jot down -t zettel mein kaffee ist weg
        The template 'zettel.md' will be used to create 'mein-kaffee-ist-weg.md'
        in your vault. It will then be opened with $EDITOR.

Run 'jot <command> -h' to get help for a specific command.`)
}

func Config(args []string, c helper.Config, template string) error {
	var t bool

	flagSet := flag.NewFlagSet("config", flag.ExitOnError)
	flagSet.BoolVar(&t, "t", false, "Print the default template.")
	flagSet.Parse(args)

	switch t {
	case false:
		cPath := viper.ConfigFileUsed()
		cContent, err := os.ReadFile(cPath)
		helper.Check(err)

		helper.Logger.Printf("%s:\n\n%s\n", cPath, cContent)
		helper.Logger.Info("Run 'jot config -t to print the set template.")

	case true:
		helper.Logger.Printf("%s\n", template)
		helper.Logger.Info("Run 'jot config' to print the config.")
	}

	return nil
}
