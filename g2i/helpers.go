package g2i

import (
	"io"
	"text/template"
)

func (c *Config) generateTemplate(name, path string, data interface{}, wr io.Writer) error {
	t, err := template.New(name).ParseFiles(c.IDEC.HelloMessageTemplatePath)
	if err != nil {
		return err
	}
	return t.Execute(wr, data)
}
