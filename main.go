package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "os/exec"
    "time"
)

func main() {
    homeDir, err := os.UserHomeDir()
    if err != nil {log.Fatal(err)}

    newCmd := flag.NewFlagSet("new", flag.ExitOnError)
    notesDir := newCmd.String("dir", homeDir+"/notes", "Notes directory")
    editor := newCmd.String("editor", "nvim", "Editor application")
    switch os.Args[1] {
    case "new":
        newCmd.Parse(os.Args[2:])

        notePath := *notesDir+"/"+os.Args[2]+".md"
        date := time.Now()
        template := "---\n"+date.String()+"\n---\n\n# "+os.Args[2]
        edit := exec.Command(*editor, notePath)
        edit.Stdin = os.Stdin
        edit.Stdout = os.Stdout
        edit.Stderr = os.Stderr

        err := os.WriteFile(notePath, []byte(template), 0666)
        if err != nil {log.Fatal(err)}
        fmt.Println("Created note: "+notePath)
        err = edit.Run()
        if err != nil {log.Fatal(err)}
    default:
        fmt.Println("Expected 'new' command")
    }
}
