package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "slices"
    "sort"
    "strings"
    "time"

    "github.com/charmbracelet/log"
    "github.com/lithammer/fuzzysearch/fuzzy"
)

func main() {
    // TODO
    // add config file for primary variables
    var homeDir, err = os.UserHomeDir()
        if err != nil {log.Fatal(err)}
    var notesDir = homeDir+"/notes"
    var editor = "nvim"
    var date = time.Now().Format("01/02/2006 3:04 PM")
    var template = "---\ntitle:\ndate: "+date+"\ntags:\n---\n\n#"

    if len(os.Args) < 2 {
        var openEditor = exec.Command(editor)
        openEditor.Dir = notesDir
        openEditor.Stdin = os.Stdin
        openEditor.Stdout = os.Stdout
        openEditor.Stderr = os.Stderr

        err = openEditor.Run()
            if err != nil {log.Fatal(err)}
    } else {
        var allNoteNames []string
        var allNotes, err = os.ReadDir(notesDir)
            if err != nil {log.Fatal(err)}
        for _, file := range allNotes {
            allNoteNames = append(allNoteNames, file.Name())
        }
        var fileName = strings.Join(os.Args[1:], "_")+".md"
        for slices.Contains(allNoteNames, fileName) {
            var matches = fuzzy.RankFind(fileName, allNoteNames)
            sort.Sort(matches)

            log.Error("A note with this filename already exists")
            fmt.Println("\033[1mMatching notes:\033[0m")
            for _, match := range matches {
                fmt.Println(" - \033[32m"+match.Target+"\033[0m")
            }
            fmt.Println("\033[1mPlease input a new filename:\033[0m")
            var newFileName, err = bufio.NewReader(os.Stdin).ReadString('\n')
                if err != nil {log.Fatal(err)}
            fileName = strings.TrimSpace(strings.ReplaceAll(newFileName, " ", "_"))+".md"
            for strings.HasPrefix(fileName, "_") || strings.HasPrefix(fileName, ".") {
                log.Error("Filename may not be null")
                fmt.Println("\033[1mPlease input a new filename:\033[0m")
                newFileName, err = bufio.NewReader(os.Stdin).ReadString('\n')
                    if err != nil {log.Fatal(err)}
                fileName = strings.TrimSpace(strings.ReplaceAll(newFileName, " ", "_"))+".md"
            }
        }
        var notePath = notesDir+"/"+fileName

        var openNote = exec.Command(editor, notePath)
        openNote.Stdin = os.Stdin
        openNote.Stdout = os.Stdout
        openNote.Stderr = os.Stderr

        err = os.WriteFile(notePath, []byte(template), 0666)
            if err != nil {log.Fatal(err)}
        log.Info("Creating note... "+notePath)
        err = openNote.Run()
            if err != nil {log.Fatal(err)}
    }
}
