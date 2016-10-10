# linker
A golang based cli application that takes a SOURCE path and sets up symlinks to its whole hierarchy from a TARGET path.

## Background
I needed a simple tool to quickly setup a linux user configuration by symlinking to the actual configuration files kept in a git repository.

## TODOs
### Features
- [x] CLI handling
- [x] Read in source files
- [x] Switch on existingFile, existingSymlink, newFile
- [x] Diff files on content
- [x] Create directory structure if missing
- [x] Create symlink
- [ ] Backup existing files
- [ ] Add dry-run mode (no changes performed to FS)
- [ ] Add force mode (no question)

### Verification
- [ ] ENV & ARGS (a whole bunch of variations)
- [ ] new file
- [ ] new file without directories
- [ ] existing symlink
- [ ] existing file
- [ ] dry-run
- [ ] Benchmark: Mimic each "handle" test case with a benchmark

### Improvements
- [ ] Double printouts of short/long version arguments in helper (as well as double handling in the code)
- [ ] Integrate debug control into the logger instead of having to have ugly if statements directly in the code
- [x] Wrap up or simplify response handling (2 cases atm)
- [ ] What is the correct FileMode to use when making directories instead of 0755?
- [x] Fix so that a single file can be given as source instead of a directory
- [ ] Fix so that the prompter accepts an empty string (just newline)
