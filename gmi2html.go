package main

import (
	"flag"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/boomlinde/gemini/gemini"
)

func main() {
	templateFlag := flag.String("template", "", "template to use")
	flag.Parse()

	var tplText = defaultTemplate
	if *templateFlag != "" {
		content, err := ioutil.ReadFile(*templateFlag)
		fatal("reading template", err)
		tplText = string(content)
	}
	tpl, err := template.New("gmi").Parse(tplText)
	fatal("parsing template", err)

	lines := []gemini.Line{}
	args := flag.Args()
	if len(args) > 0 {
		for _, arg := range args {
			f, err := os.Open(arg)
			fatal("opening gmi file", err)
			defer f.Close()
			l, err := gemini.Itemize(f)
			fatal("parsing gmi file", err)
			f.Close()
			lines = append(lines, l...)
		}
	} else {
		l, err := gemini.Itemize(os.Stdin)
		fatal("parsing gmi on stdin", err)
		lines = append(lines, l...)
	}

	buf := make([]byte, 0, 100000)
	w := func(format string, a ...interface{}) {
		buf = append(buf, fmt.Sprintf(format, a...)...)
	}
	e := html.EscapeString

	paragraph := false
	lastLine := gemini.Line{Type: -1}
	for i, l := range lines {
		nextLine := gemini.Line{Type: -1}
		if i < len(lines)-1 {
			nextLine = lines[i+1]
		}

		if (l.Type != gemini.TextLine || l.Display == "") && paragraph {
			paragraph = false
			w("</p>\n")
		}

		switch l.Type {
		case gemini.LinkLine:
			if lastLine.Type != gemini.LinkLine {
				w("<ul class=\"linklist\">\n")
			}
			w("<li><a href=\"%s\">%s</a></li>\n", e(l.Link), e(l.Display))
			if nextLine.Type != gemini.LinkLine {
				w("</ul>\n")
			}
		case gemini.PreLine:
			if lastLine.Type != gemini.PreLine {
				w("<pre>")
			}
			w(e(l.Raw))
			if nextLine.Type != gemini.PreLine {
				w("</pre>")
			}
			w("\n")
		case gemini.TextLine:
			if l.Display == "" {
				break
			}

			if !paragraph {
				w("<p>\n")
				paragraph = true
			}
			w(e(l.Display))
			if nextLine.Type == gemini.TextLine && nextLine.Display != "" {
				w("<br>")
			}
			w("\n")
		case gemini.ListLine:
			if lastLine.Type != gemini.ListLine {
				w("<ul>\n")
			}
			w("<li>%s</li>\n", e(l.Display))
			if nextLine.Type != gemini.ListLine {
				w("</ul>\n")
			}
		case gemini.H1Line:
			w("<h1>%s</h1>\n", e(l.Display))
		case gemini.H2Line:
			w("<h2>%s</h2>\n", e(l.Display))
		case gemini.H3Line:
			w("<h3>%s</h3>\n", e(l.Display))
		case gemini.QuoteLine:
			if lastLine.Type != gemini.QuoteLine {
				w("<blockquote>\n")
			}
			w("%s<br>\n", e(l.Display))
			if nextLine.Type != gemini.QuoteLine {
				w("</blockquote>\n")
			}
		default:
			fatal("emitting HTML", fmt.Errorf("invalid line type %d", l.Type))
		}

		lastLine = l
	}
	if paragraph {
		paragraph = false
		w("</p>\n")
	}

	var title = "Untitled document"
	if len(lines) > 0 && lines[0].Type == gemini.H1Line {
		title = lines[0].Display
	}

	tplData := struct {
		Title   string
		Content template.HTML
	}{
		Title:   title,
		Content: template.HTML(buf),
	}

	fatal("executing template", tpl.Execute(os.Stdout, tplData))
}

func fatal(prefix string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s: %v", prefix, err))
		os.Exit(1)
	}
}
