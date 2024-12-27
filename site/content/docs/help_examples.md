---
linkTitle: 'Getting Help'
title: 'Getting Help'
weight: 6
---

### The Help Flag

tabtomd's `--help` flag provides an overview of the application's
usage, the application [commands](/docs/commands), and the various
[flags](/docs/flags) that can be used.

```bash
tabtomd --help
```

```text
Convert tab delimited data to a Markdown table

Usage:
  tabtomd [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  file        convert a tab-delimited file to a Markdown table
  help        Help about any command
  paste       paste tab delimited data from the clipboard
  version     Print the tabtomd version

Flags:
  -h, --help   help for tabtomd

Use "tabtomd [command] --help" for more information about a command.
```

### Command Help

Each command has additional help information available. Help for a command
can be viewed either by using the special `help` command or using the
`--help` flag after the command.

#### Using the `help` command

The example below uses the `help` command to view help for the `file` command.

```bash
$ tabtomd help file
```

```text
The file command reads a tab-delimited file and converts its content into a
Markdown formatted table.

Usage:
  tabtomd file [filename] [flags]

Flags:
  -h, --help            help for file
      --output string   Specify the output file to save the Markdown table
      --pretty          Use pretty Markdown table formatting
      --print           Print to screen
```

#### Using the `--help` flag

The example below uses the `--help` flag to view help for the `paste` command.

```bash
tabtomd paste --help
```

```text
The paste command converts tab delimited data from the clipboard into a
Markdown formatted table.

Usage:
  tabtomd paste [flags]

Flags:
  -h, --help            help for paste
      --output string   Specify the output file to save the Markdown table
      --pretty          Use pretty Markdown table formatting
      --print           Print to screen
```
