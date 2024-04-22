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
    var homeDir, err = os.UserHomeDir()
        if err != nil {log.Fatal(err)}
    // TODO
    // add flags for notesDir and editor
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

            log.Error("A note with this name already exists")
            fmt.Println("Matching notes:")
            for _, match := range matches {
                fmt.Println(" - "+match.Target)
            }
            fmt.Println("Please input a new file name:")
            var newFileName, err = bufio.NewReader(os.Stdin).ReadString('\n')
                if err != nil {log.Fatal(err)}
            fileName = strings.TrimSpace(strings.ReplaceAll(newFileName, " ", "_"))+".md"
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
