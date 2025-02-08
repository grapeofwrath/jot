package cmd

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/grapeofwrath/jot/helper"
)

func Down(args []string, c helper.Config, template string) error {
	var templateName string

	flagSet := flag.NewFlagSet("new", flag.ExitOnError)
	flagSet.StringVar(
		&templateName,
		"t",
		"",
		"Specify a template to use from template.directory")
	flagSet.Usage = func() {
		helper.Logger.Print(`Create a new note in your vault.

Usage:

    jot down [flags] [args]

Flags:

    -t  STRING
        Specify the template file name to use from the template directory.

Example:

    jot down -t zettel guten tag
        The template 'zettel.md' will be used to create 'guten-tag.md'.`)
	}
	flagSet.Parse(args)

	var title = strings.Join(flagSet.Args()[0:], " ")

	title, err := helper.FormatTitle(c.Template.Format.Title, title)
	helper.Check(err)

	fileName, title, err := helper.DuplicateCheck(c.Vault.Directory, c.Template.Format.Title, title)
	helper.Check(err)

	if templateName != "" {
		tPath := filepath.Join(c.Template.Directory, templateName+".md")
		tByte, err := os.ReadFile(tPath)
		helper.Check(err)
		template = string(tByte)
	}

	template, err = helper.MkTemplate(c, template, title)
	helper.Check(err)

	filePath := filepath.Join(c.Vault.Directory, fileName)

	file, err := os.Create(filePath)
	helper.Check(err)
	defer file.Close()

	_, err = file.Write([]byte(template))
	helper.Check(err)

	editCmd := exec.Command(c.Editor, filePath)
	editCmd.Stdin = os.Stdin
	editCmd.Stdout = os.Stdout
	editCmd.Stderr = os.Stderr

	err = editCmd.Run()
	helper.Check(err)

	return nil
}
