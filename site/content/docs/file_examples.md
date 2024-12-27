---
linkTitle: "File Command Examples"
title: "File Command Examples"
weight: 4
---

Here are practical examples demonstrating how to use **tabtomd** with the
`file` command to convert tab-delimited files into Markdown tables.

### Convert a tab-delimited file to a Markdown table and print it to the screen

Suppose you have a tab-delimited file, `characters.tsv`, containing the
following text:

{{< tabs items="characters.tsv" defaultIndex="0" >}}
{{< tab >}}
```text
Name	Species	Gender
Ariel	Mermaid	Female
Sebastian	Crab	Male
Chef Louis	Human	Male
Prince Eric	Human	Male
Ursula	Octopus	Female
```
{{< /tab >}}
{{< /tabs >}}


To convert this file to a Markdown table and print the result to the
screen, use the `file` [command](/docs/commands#command-list) and the
`--print` [flag](/docs/commands#flag-list):

```bash
tabtomd file characters.tsv --print
```

**tabtomd** reads the tab-delimited content from the file, treats the first
row as a header, and outputs the following Markdown table:

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

### Convert a tab-delimited file to an aligned Markdown table and print it to the screen

To produce a neatly aligned Markdown table with evenly spaced columns, use
the `--pretty` flag along with the `--print` flag:

```bash
tabtomd file characters.tsv --pretty --print
```

This generates a table with aligned columns, as shown below:

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

### Convert a tab-delimited file to an aligned Markdown table and save it to a file

To save the converted table to a file instead of printing it to the screen,
use the `--output` [flag](/docs/commands#flag-list). When this flag is
used, the output is redirected to the specified file.

For example:

```bash
tabtomd file characters.tsv --output characters.md --pretty
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

These examples showcase how **tabtomd** simplifies the conversion of
tab-delimited data into Markdown format, whether you need to display the
results immediately or save them for later use. Using the `--pretty` flag
ensures the output is well-aligned for improved readability.
