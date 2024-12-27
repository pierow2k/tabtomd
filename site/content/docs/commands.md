---
title: "Commands and Flags"
weight: 3
---

**`tabtomd`** follows a standard command-line application structure:

```bash
tabtomd [COMMAND] [OPTIONS]
```

- **`[COMMAND]`**: Specifies the action to perform (see [Command List](#command-list)).
- **`[OPTIONS]`**: Optional flags that modify the behavior of the command (see [Flag List](#flag-list)).

### Command List

The most commonly used commands are **`file`** and **`paste`**, which are
often combined with [flags](#flag-list) to control the output. The table
below lists all available commands and their descriptions:

| Command          | Description                                                         |
| ---------------- | ------------------------------------------------------------------- |
| **`file`**       | Convert a tab-delimited _file_ to a Markdown table.                 |
| **`paste`**      | Convert tab-delimited data _from the system clipboard_ to Markdown. |
|                  |                                                                     |
| **`completion`** | Generate an autocompletion script for a specific shell.             |
| **`help`**       | Display detailed help information for any command.                  |
| **`version`**    | Display the **tabtomd** version and build details.                  |

### Flag List

Flags are optional parameters that modify how commands behave. The table
below describes the flags available for use with **tabtomd**:

| Flag                       | Description                                                        |
| -------------------------- | ------------------------------------------------------------------ |
| **`--output`** *`string`*  | Specify a file to save the generated Markdown table.               |
| **`--pretty`**             | Align vertical bars in the Markdown table for improved readability. |
| **`--print`**              | Print the Markdown table directly to standard output.       |

### Examples

For practical demonstrations of how to use commands and flags, visit the
examples pages:

- [Examples for the **`file`** command](/docs/file_examples)
- [Examples for the **`paste`** command](/docs/paste_examples)
- [**`help`** command examples](/docs/help/#the-help-flag)
