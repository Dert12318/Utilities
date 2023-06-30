package html

import (
	"fmt"
	"html/template"
	"os"
)

type HTMLData struct {
	Data string
}

type HTML struct {
	templates map[string]*template.Template
}

func NewHTMLGenerator(filepath map[string]string) (*HTML, error) {

	templates, err := loadTemplates(filepath)
	if err != nil {
		err = fmt.Errorf("FAILED TO LOAD TEMPLATES %s", err.Error())
		return &HTML{}, err
	}

	return &HTML{
		templates: templates,
	}, nil
}

func loadTemplates(filepath map[string]string) (map[string]*template.Template, error) {

	htmlTemplate := make(map[string]*template.Template)

	for name, path := range filepath {
		pwd, _ := os.Getwd()
		pathTemplate := pwd + path
		emailTemplate, err := template.ParseFiles(pathTemplate)

		if err != nil {
			err = fmt.Errorf("failed to load template %s : %v", name, err.Error())
			return htmlTemplate, err
		}
		htmlTemplate[name] = emailTemplate
	}

	return htmlTemplate, nil
}
