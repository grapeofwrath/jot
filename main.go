package main

import (
    "fmt"
    "os"
    "os/exec"
    //"sort"
    "strings"
    "time"

    "github.com/charmbracelet/log"
    "github.com/lithammer/fuzzysearch/fuzzy"
)

func main() {
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
        var fileName = strings.Join(os.Args[1:], "_")
        var notePath = notesDir+"/"+fileName+".md"

        if _, err := os.Stat(notePath); err == nil {
            var allNoteNames []string
            var allNotes, err = os.ReadDir(notesDir)
                if err != nil {log.Fatal(err)}
            for _, file := range allNotes {
                allNoteNames = append(allNoteNames, file.Name())
            }
            var matches = fuzzy.Find(fileName, allNoteNames)
            // TODO
            // sort.Sort(matches)

            log.Error("Note exists: "+notePath)
            fmt.Println("Matching files:")
            for _, match := range matches {
                fmt.Println(" - "+match)
            }
            //fmt.Println("Please input a new file name:")
            os.Exit(1)
            // TODO
            // prompt to rename note
            // while loop until new name is input
        }

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
