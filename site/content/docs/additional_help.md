---
linkTitle: 'Additional Help'
title: 'Additional Help'
weight: 7
---

### Manual

In addition to the built-in [help flags and commands](/docs/help_examples),
**tabtomd** comes with a complete man page in both troff (man page) and in
PDF formats.

- **Man Page in troff Format** [troff format](/files/tabtomd.1)
- **man Page in PDF Format** [PDF](/files/tabtomd.1.pdf)

```text
tabtomd(1)                   General Commands Manual                  tabtomd(1)

NAME
       tabtomd - Convert tab delimited data to Markdown table

SYNOPSIS
       tabtomd [COMMAND] [OPTIONS]

DESCRIPTION
       tabtomd is a versatile command-line utility for converting tab-delimited
       text into Markdown table format.  It supports input from either the
       system clipboard or a specified file and allows the output to be written
       to standard output or saved to a file.

       The tool includes an optional --pretty flag to align vertical bars in the
       Markdown table, improving readability and presentation.  It is especially
       useful for generating Markdown tables for documentation, reports, or
       other structured data presentations.

COMMANDS
       completion
              Generate an autocompletion script for the specified shell.

       file   Convert a tab-delimited file to a Markdown table.

       help   Display help information for any command.

       paste  Convert tab-delimited data from the system clipboard.

       version
              Display the tabtomd version and build information.

FLAGS
       --output string
              Specify a file to save the generated Markdown table.

       --pretty
              Align vertical bars in the Markdown table for better readability.

       --print
              Print the Markdown table to standard output.
…
```

### Open an Issue

Have an idea for a new feature or noticed something that isn’t working
quite right? [Open an issue](ttps://github.com/pierow2k/tabtomd/issues)
to let us know. Your feedback helps us keep tabtomd reliable and
feature-rich.
