package webserver

import (
	"encoding/json"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/j4rv/hsr-tct/pkg/hsrtct"
)

var (
	templates   *template.Template
	templateDir = "web/template"
	mu          sync.RWMutex
)

type TemplateData struct {
	ContextJSON string
	Stats       []hsrtct.Stat
	AttackTags  []hsrtct.DamageTag
	Elements    []hsrtct.Element
}

func newTemplateData(contextObj interface{}) TemplateData {
	jsonData, _ := json.Marshal(contextObj)
	return TemplateData{
		ContextJSON: string(jsonData),
		Stats:       hsrtct.StatKeys(),
		Elements:    hsrtct.ElementKeys(),
		AttackTags:  hsrtct.AttackTagKeys(),
	}
}

func loadTemplates() {
	mu.Lock()
	defer mu.Unlock()

	templates = template.Must(template.New("").Funcs(template.FuncMap{
		"safeJS": func(s string) template.HTML {
			return template.HTML(s)
		},
	}).ParseFiles(getTemplateFiles()...))

	for _, t := range templates.Templates() {
		log.Printf("Loaded template: %s", t.Name())
	}
}

func getTemplateFiles() []string {
	var files []string

	// Load templates from the root directory
	rootFiles, err := filepath.Glob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		log.Fatal(err)
	}
	files = append(files, rootFiles...)

	// Load templates from subdirectories
	subDirFiles, err := filepath.Glob(filepath.Join(templateDir, "**/*.html"))
	if err != nil {
		log.Fatal(err)
	}
	files = append(files, subDirFiles...)

	return files
}

func watchTemplateChanges() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			err = watcher.Add(path)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("Reloading templates due to change:", event.Name)
				time.Sleep(100 * time.Millisecond) // Give some time for the file system to settle
				loadTemplates()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error watching templates:", err)
		}
	}
}
