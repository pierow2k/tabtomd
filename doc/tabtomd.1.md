% tabtomd(1) Version {{version}}| General Commands Manual
%
% {{version}}


NAME
====
**tabtomd** - Convert tab delimited data to Markdown table

SYNOPSIS
========

| **tabtomd** \[**COMMAND**\] \[**OPTIONS**\]

DESCRIPTION
===========

`tabtomd` is a versatile command-line utility for converting tab-delimited
text into Markdown table format. It supports input from either the system
clipboard or a specified file and allows the output to be written to
standard output or saved to a file.

The tool includes an optional `--pretty` flag to align vertical bars in the
Markdown table, improving readability and presentation. It is especially
useful for generating Markdown tables for documentation, reports, or other
structured data presentations.

COMMANDS
========

**completion**
:   Generate an autocompletion script for the specified shell.

**file**
:   Convert a tab-delimited file to a Markdown table.

**help**
:   Display help information for any command.

**paste**
:   Convert tab-delimited data from the system clipboard.

**version**
:   Display the `tabtomd` version and build information.

FLAGS
=====

**\-\-output *string***
:   Specify a file to save the generated Markdown table.

**\-\-pretty**
:   Align vertical bars in the Markdown table for better readability.

**\-\-print**
:   Print the Markdown table to standard output.

EXAMPLES
========

### Convert a tab-delimited file to a Markdown table and print it to the screen

`tabtomd file tsv_file.txt --print`

### Convert a tab-delimited file to an aligned Markdown table and print it to the screen

`tabtomd file tsv_file.txt --pretty --print`

### Convert a tab-delimited file to an aligned Markdown table and save it to a file

`tabtomd file tsv_file.txt --output markdown_table.md --pretty`

### Convert tab-delimited text from the clipboard to a Markdown table and print it to the screen

`tabtomd paste --print`

### Convert tab-delimited text from the clipboard to an aligned Markdown table and print it to the screen

`tabtomd paste --pretty --print`

### Convert tab-delimited text from the clipboard to an aligned Markdown table and save it to a file

`tabtomd paste --output markdown_table.md --pretty`

### Display help for the `file` command

`tabtomd file --help` **OR** `tabtomd help file`

COPYRIGHT
=========

Copyright (C) {{copyright_date}} Pierow2K. Released under the MIT License.
