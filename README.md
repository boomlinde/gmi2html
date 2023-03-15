# gmi2html

Generate HTML files from gmi (text/gemini) documents.

Makes an opinionated, non-standard interpretation of text/gemini: consecutive
non-empty text lines are considered paragraphs. Line breaks are still inserted
between non-empty text lines.

## Usage

    gmi2html [-template TEMPLATE] [FILE]...

* If more than one file is given, the files are concatenated before converted.
* If no files are given, `gmi2html` will instead read the document on stdin
* `gmi2html` will output the resulting document to stdout
* The template is a Go `html/template` where `.Content` and `.Title` may be
   used
* The `.Title` parameter is the first line of the input if that is a `#`
  heading, "Untitled document" otherwise
