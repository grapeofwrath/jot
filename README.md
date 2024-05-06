# Jot

**A Zettelkasten CLI helper written in Go**

<img alt="Welcome to Jot" src="./assets/jot_demo.gif" width="600" />

When Jot is ran without any additional arguments, it will open Neovim inside the ~/notes directory.
To quickly create a new note, the executable can be run using all additional arguments as the title.
A note is then created using the following template and opened in Neovim:

```markdown
---
title: This Is A Zettel Title
date: 06/27/2024 02:30 AM
tags:
---

# This Is A Zettel Title
```

## Configuration

You can create `~/.config/jot/config.yaml` in order to adjust Jot how you like.

```yaml
NotesDir: "/home/user/notes"
Editor: "nvim"
Template: "---\ntitle: $title\ndate: $date\ntags:\n---\n\n# $title"
```

Configuration Options:

`NotesDir`: The directory for your notes
`Editor`: The command to launch your editor (ie nvim not neovim)
`Template`: The template to use when creating a new note (use `$date` and `$title` if you wish to include those)

This project is created to help with the Zettelkasten method, but it does not have to be used solely for that purpose.
However, I don't plan on adding anything more than needed to help with that method.
This will hopefully keep the project fairly simple and focused.

If you're wondering *what the heck is a Zettelkasten*, you can learn more about it here:

[Zettelkasten Introduction](https://zettelkasten.de/introduction/)
