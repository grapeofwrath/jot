# Jot

**A Zettelkasten CLI helper written in Go**

Currently, this assumes a few defaults:

 - `$HOME/notes` is the directory used to store the Zettelkasten
 - Neovim is the editor to open

If the executable is ran without any additional arguments, it will open Neovim inside the notes directory.
To quickly create a new note, the executable can be run using all additional arguments as the title.

```shell
jot this is a zettel title

Output:

2022/06/27 02:30:30 INFO Creating note... /home/user/notes/this_is_a_zettel_title.md
```

A note is then created using the following template, and then is the opened in the editor:

```markdown
---
title:
date: 06/27/2024 02:30 AM
tags:
---

#
```

This project is created to help with the Zettelkasten method, but it does not have to be used solely for that purpose.
However, I don't plan on adding anything more than needed to help with that method.
This will hopefully keep the project fairly simple and focused.

If you're wondering *what the heck is a Zettelkasten*, you can learn more about it here:

[Zettelkasten Introduction](https://zettelkasten.de/introduction/)
