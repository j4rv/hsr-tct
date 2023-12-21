package webserver

import (
	"html/template"
	"os"
	"path/filepath"
)

var templates = template.Must(template.New("").Parse(""))

func loadTemplates() {
	templateDir := "web/template"
	filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			templates = template.Must(templates.ParseFiles(path))
		}
		return nil
	})
}
