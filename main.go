package main

import (
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type Variables struct {
	TemplateFile string
}

func getWriter() io.Writer {
	if out == "STDOUT" {
		return os.Stdout
	}

	w, err := os.Create(out)
	if err != nil {
		panic(err)
	}

	return w
}

func cat(f string) string {
	d, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}
	return string(d)
}

func diff(f1, f2 string) string {
	out, _ := exec.Command("diff", f1, f2).CombinedOutput()
	return string(out)
}

func safeHTML(val string) template.HTML {
	return template.HTML(val)
}

func makeReadme(tmpl string, w io.Writer, v *Variables) {
	funcMap := template.FuncMap{
		"cat":      cat,
		"diff":     diff,
		"safeHTML": safeHTML,
	}

	t := template.Must(template.New(tmpl).Funcs(funcMap).ParseFiles(tmpl))
	err := t.Execute(w, v)
	if err != nil {
		panic(err)
	}
}

var (
	tmpl string
	out  string
)

func init() {
	flag.StringVar(&tmpl, "t", "README.md.tmpl", "Use <file> as the tmplate instead of README.md.tmpl")
	flag.StringVar(&out, "o", "README.md", "Write output to <file> instead of README.md")
	flag.Parse()
}

func main() {
	var v = &Variables{
		TemplateFile: tmpl,
	}

	w := getWriter()
	makeReadme(tmpl, w, v)
}
