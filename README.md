textlint - a collection of text validation tools
================================================

Textlint is a collection of text validation tools. It tries to do things in parallel when it makes sense.

Usage
-----

### Detect the presence of null characters

```
Check for the presence of null bytes in a file.
It uses multiple threads to process the file in chunks and validate the content.

Usage:
  textlint null [flags]

Flags:
  -b, --batches int   Number of batches to use (default 64)
  -h, --help          help for null
  -j, --threads int   Number of cores to use (default 4)
```

### Count the newline, tab and null characters

```
Count the number of newline, tab and null characters in a file.

Usage:
  textlint lnc [flags]

Flags:
  -b, --batches int   Number of batches to use (default 64)
  -h, --help          help for lnc
  -j, --threads int   Number of cores to use (default 4)
```