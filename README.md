# linker
A golang based cli application that takes a SOURCE path and sets up symlinks to its whole hierarchy from a TARGET path.

## Background
I needed a simple tool to quickly setup a linux user configuration by symlinking to the actual configuration files kept in a git repository.

## TODOs
* CLI handling
* Read in source files
* Switch on existingFile, existingSymlink, newFile
* Diff files on content
* Create directory structure if missing
* Create symlink
* Backup existing files
* Add dry-run mode (no changes performed to FS)
* Add force mode (no question)
