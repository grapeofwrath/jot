package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"gopkg.in/yaml.v3"
)

// TODO
// change Template to read from a template file
type Config struct {
	NotesDir string `yaml:"notesdir"`
	Editor   string `yaml:"editor"`
	Template string `yaml:"template"`
}

func main() {
	var err error
	var homeDir string
	var cfg Config
	var allNoteNames []string

	homeDir, err = os.UserHomeDir()
	if err != nil {
		log.Printf("ERROR: %v", err)
	}
	// change to config
	notesDir := filepath.Join(homeDir, "notes")
	editor := "nvim"
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Printf("ERROR: %v", err)
	}
	configCheck, err := filepath.Glob(filepath.Join(configDir, "jot", "config.yaml"))
	if err != nil {
		log.Printf("ERROR: %v", err)
	}

	if len(configCheck) > 0 {
		data, err := os.ReadFile(configCheck[0])
		if err != nil {
			log.Printf("ERROR: %v", err)
		}
		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			log.Printf("ERROR: %v", err)
		}
		if cfg.NotesDir != "" {
			notesDir = cfg.NotesDir
		}
		if cfg.Editor != "" {
			editor = cfg.Editor
		}
	}

	if len(os.Args) < 2 {
		openEditor := exec.Command(editor)
		openEditor.Dir = notesDir
		openEditor.Stdin = os.Stdin
		openEditor.Stdout = os.Stdout
		openEditor.Stderr = os.Stderr

		err = openEditor.Run()
		if err != nil {
			log.Printf("ERROR: %v", err)
		}
		os.Exit(0)
	}

	allNotes, err := os.ReadDir(notesDir)
	if err != nil {
		log.Printf("ERROR: %v", err)
	}

	for _, file := range allNotes {
		allNoteNames = append(allNoteNames, file.Name())
	}

	filename := strings.Join(os.Args[1:], "-") + ".md"
	for slices.Contains(allNoteNames, filename) {
		matches := fuzzy.RankFind(filename, allNoteNames)
		sort.Sort(matches)
		log.Println("ERROR: A note with this filename already exists")
		fmt.Println("\033[1mMatching notes:\033[0m")
		for _, match := range matches {
			fmt.Println(" - \033[32m" + match.Target + "\033[0m")
		}
		fmt.Println("\033[1mPlease input a new filename:\033[0m")
		newFilename, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Printf("ERROR: %v", err)
		}
		filename = strings.TrimSpace(strings.ReplaceAll(newFilename, " ", "-")) + ".md"
		for strings.HasPrefix(filename, "-") || strings.HasPrefix(filename, ".") {
			log.Println("ERROR: Filename may not be null")
			log.Println("\033[1mPlease input a new filename:\033[0m")
			newFilename, err = bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				log.Printf("ERROR: %v", err)
			}
			filename = strings.TrimSpace(strings.ReplaceAll(newFilename, " ", "-")) + ".md"
		}
	}

	// title := strings.ReplaceAll(strings.TrimSuffix(filename, ".md"), "-", " ")
	title := filename
	date := time.Now().Format("2006-01-02")
	template := fmt.Sprintf(`---
title: %s
date: %s
source:
related:
tags:
  - marcus
---

`, title, date)
	notePath := filepath.Join(notesDir, filename)
	if cfg.Template != "" {
		template = strings.ReplaceAll(strings.ReplaceAll(cfg.Template, "$title", title), "$date", date)
	}

	openNote := exec.Command(editor, notePath)
	openNote.Stdin = os.Stdin
	openNote.Stdout = os.Stdout
	openNote.Stderr = os.Stderr

	err = os.WriteFile(notePath, []byte(template), 0666)
	if err != nil {
		log.Printf("ERROR: %v", err)
	}
	log.Println("Creating note...", notePath)
	err = openNote.Run()
	if err != nil {
		log.Printf("ERROR: %v", err)
	}
}
