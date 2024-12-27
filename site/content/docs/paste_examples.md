---
linkTitle: "Paste Command Examples"
title: "Paste Command Examples"
weight: 5
---

Here are practical examples demonstrating how to use **tabtomd** with the
`paste` command to convert tab-delimited data from the clipboard into
Markdown tables.

### Convert tab-delimited text from the clipboard to a Markdown table and print it to the screen

Suppose you have copied the following tab-delimited data to the clipboard:

```text
Name	Species	Gender
Ariel	Mermaid	Female
Sebastian	Crab	Male
Chef Louis	Human	Male
Prince Eric	Human	Male
Ursula	Octopus	Female
```

To convert this clipboard text into a Markdown table and print it to the
screen, use the `paste` [command](/docs/commands#command-list) with the
`--print` [flag](/docs/commands#flag-list):

```bash
tabtomd paste --print
```

**tabtomd** reads the tab-delimited text from the clipboard, treats the
first row as a header, and outputs the following Markdown table:

```text
| Name | Species | Gender |
| ----- | ----- | ----- |
| Ariel | Mermaid | Female |
| Sebastian | Crab | Male |
| Chef Louis | Human | Male |
| Prince Eric | Human | Male |
| Ursula | Octopus | Female |
```

{{< callout type="info" >}}
Without the `--pretty` flag, the Markdown table is functional but
not aligned. Columns are spaced according to the length of the text, which
may affect readability.
{{< /callout >}}

### Convert tab-delimited text from the clipboard to an aligned Markdown table and print it to the screen

To create a neatly aligned Markdown table with evenly spaced columns, use
the `--pretty` flag along with the `--print` flag:

```bash
tabtomd paste --pretty --print
```

This produces an aligned Markdown table as shown below:

```text
| Name        | Species | Gender |
| ----------- | ------- | ------ |
| Ariel       | Mermaid | Female |
| Sebastian   | Crab    | Male   |
| Chef Louis  | Human   | Male   |
| Prince Eric | Human   | Male   |
| Ursula      | Octopus | Female |
```

{{< callout type="info" >}}
With the `--pretty` flag, the Markdown table is aligned.
{{< /callout >}}

### Convert tab-delimited text from the clipboard to an aligned Markdown table and save it to a file

To save the converted table to a file instead of printing it to the screen,
use the `--output` [flag](/docs/commands#flag-list). When this flag is
used, the output is redirected to the specified file.

For example:

```bash
tabtomd paste --output characters.md --pretty
```

This creates a new file, `characters.md`, containing the aligned Markdown
table. To verify, you can view the file's content:

```bash
cat characters.md
```

Output:


{{< tabs items="characters.md" defaultIndex="0" >}}
{{< tab >}}
```text
| Name        | Species | Gender |
| ----------- | ------- | ------ |
| Ariel       | Mermaid | Female |
| Sebastian   | Crab    | Male   |
| Chef Louis  | Human   | Male   |
| Prince Eric | Human   | Male   |
| Ursula      | Octopus | Female |
```
{{< /tab >}}
{{< /tabs >}}

{{< callout type="info" >}}
Use the `--output` flag to save the Markdown table to a file.
{{< /callout >}}

These examples demonstrate how **tabtomd** can easily convert tab-delimited
data from the clipboard into Markdown format, offering both quick display
options and the ability to save the output for later use. Using the
`--pretty` flag ensures the results are well-aligned and visually appealing
for improved readability.
