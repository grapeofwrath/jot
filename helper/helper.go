package helper

import (
	"bufio"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Config struct {
	Editor   string
	Template struct {
		Directory string
		Name      string
		Format    struct {
			Date  string
			Time  int
			Title string
		}
	}
	Vault struct {
		Directory string
	}
}

var Logger = log.New(os.Stderr)

func Check(e error) {
	if e != nil {
		Logger.Fatal(e)
	}
}

func MkConfig() (Config, string, error) {
	configDir, err := os.UserConfigDir()
	Check(err)

	jConfigDir := filepath.Join(configDir, "jot")
	jConfigPath := filepath.Join(jConfigDir, "config.json")

	viper.SetConfigName("config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(jConfigDir)
	viper.AddConfigPath(".")

	viper.SetDefault("editor", os.Getenv("EDITOR"))
	viper.SetDefault("template.directory", "")
	viper.SetDefault("template.name", "")
	viper.SetDefault("template.format.date", "YYYY-MM-DD")
	viper.SetDefault("template.format.time", 24)
	viper.SetDefault("template.format.title", "kebab-case")
	viper.SetDefault("vault.directory", "")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.WriteConfigAs(jConfigPath)
		} else {
			Check(err)
		}
	}

	if viper.GetString("vault.directory") == "" {
		Logger.Error("vault missing")
		Logger.Print("\033[1mPlease provide the path to your vault:\033[0m")

		vaultDirectory, err := bufio.NewReader(os.Stdin).ReadString('\n')
		Check(err)

		viper.Set("vault.directory", strings.TrimSpace(vaultDirectory))
		templateDir := filepath.Join(viper.GetString("vault.directory"), "Templates")
		viper.Set("template.directory", templateDir)

		Logger.Infof(
			"Vault Directory set as %s\nTemplate Directory set as %s",
			viper.GetString("vault.directory"),
			viper.GetString("template.directory"))

		viper.WriteConfigAs(jConfigPath)
	}

	var config Config

	err = viper.Unmarshal(&config)
	Check(err)

	var template = `---
title: {{title}}
date: {{date}}
source:
related:
tags:
---`

	if config.Template.Name != "" {
		tPath := filepath.Join(
			config.Template.Directory,
			config.Template.Name+".md")
		tByte, err := os.ReadFile(tPath)
		Check(err)

		template = string(tByte)
	}

	return config, template, nil
}

func MkTemplate(c Config, template string, title string) (string, error) {
	template = strings.ReplaceAll(template, "{{title}}", title)
	template = strings.ReplaceAll(template, "{{date}}", formatDate(c.Template.Format.Date))
	template = strings.ReplaceAll(template, "{{time}}", formatTime(c.Template.Format.Time))

	return template, nil
}

func formatDate(format string) string {
	var date string

	switch format {
	case "YYYY-MM-DD":
		date = time.Now().Format("2006-01-02")
	case "YYYY/MM/DD":
		date = time.Now().Format("2006-01-02")
	case "DD-MM-YYYY":
		date = time.Now().Format("02-01-2006")
	case "DD/MM/YYYY":
		date = time.Now().Format("02-01-2006")
	case "MM-DD-YYYY":
		date = time.Now().Format("01-02-2006")
	case "MM/DD/YYYY":
		date = time.Now().Format("01-02-2006")
	}

	return date
}

func formatTime(format int) string {
	var t string

	switch format {
	case 12:
		t = time.Now().Format("3:04 PM")
	case 24:
		t = time.Now().Format("15:04")
	}

	return t
}

func FormatTitle(format string, title string) (string, error) {
	switch format {
	case "camelCase":
		title = cases.Title(language.AmericanEnglish, cases.NoLower).String(title)
		title = strings.ReplaceAll(title, " ", "")
		if len(title) > 0 {
			title = strings.ToLower(title[:1]) + title[1:]
		}
	case "PascalCase":
		title = cases.Title(language.AmericanEnglish, cases.NoLower).String(title)
		title = strings.ReplaceAll(title, " ", "")
	case "kebab-case":
		title = strings.ReplaceAll(title, " ", "-")
	case "snake_case":
		title = strings.ReplaceAll(title, " ", "_")
	}

	return title, nil
}

func DuplicateCheck(vaultDirectory string, format string, title string) (string, string, error) {
	var existingNames []string

	vaultNotes, err := os.ReadDir(vaultDirectory)
	Check(err)

	for _, n := range vaultNotes {
		existingNames = append(existingNames, n.Name())
	}

	fileName := title + ".md"

	for slices.Contains(existingNames, fileName) {
		titleMatches := fuzzy.RankFind(fileName, existingNames)
		sort.Sort(titleMatches)

		Logger.Errorf("%s already exists", fileName)
		Logger.Info("Matching notes:")
		for _, m := range titleMatches {
			Logger.Print(" - " + m.Target)
		}

		Logger.Print("Please input a new title:")
		nTitle, err := bufio.NewReader(os.Stdin).ReadString('\n')
		Check(err)

		nTitle = strings.TrimSpace(nTitle)
		title, err = FormatTitle(format, nTitle)
		Check(err)
		fileName = title + ".md"

		for strings.HasPrefix(fileName, "-") || strings.HasPrefix(fileName, "_") || strings.HasPrefix(fileName, ".") {
			Logger.Error("Title may not be null")
			Logger.Print("Please input a new title:")
			nTitle, err := bufio.NewReader(os.Stdin).ReadString('\n')
			Check(err)

			nTitle = strings.TrimSpace(nTitle)
			title, err = FormatTitle(format, nTitle)
			Check(err)
			fileName = title + ".md"
		}
	}

	title = strings.TrimSuffix(fileName, ".md")

	return fileName, title, nil
}
