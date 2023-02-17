package html

import (
	"bytes"
	"fmt"
)

func (t *HTML) ParseTemplate(filename string, data map[string]interface{}) (HTMLData, error) {

	template := t.templates[filename]
	var body string
	buf := new(bytes.Buffer)
	if err := template.Execute(buf, data); err != nil {
		err = fmt.Errorf("failed to parse template with data : %s", err.Error())
		return HTMLData{}, err
	}

	body = buf.String()
	return HTMLData{Data: body}, nil
}
