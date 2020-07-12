package views

import (
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/markbates/pkger"

	"github.com/sirupsen/logrus"
)

const templateExtension = ".html.tmpl"

type IndexViewModel struct {
	Table string
}

type TableViewModel struct {
	Table string
}

func Parse() (*template.Template, error) {
	t := template.New("")
	return t, pkger.Walk("/hosting/views", func(p string, info os.FileInfo, _ error) error {
		if !strings.HasSuffix(p, templateExtension) || info == nil || info.IsDir() {
			return nil
		}

		_, name := path.Split(p)
		templateName := strings.TrimSuffix(name, templateExtension)

		f, err := pkger.Open(p)
		if err != nil {
			return err
		}

		src, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		logrus.WithFields(logrus.Fields{"prefix": "templates", "template": templateName}).Debug("Loading Template")
		_, err = t.New(templateName).Parse(string(src))
		return err
	})
}
