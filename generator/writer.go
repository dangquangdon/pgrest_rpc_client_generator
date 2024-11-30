package generator

import (
	"bytes"
	"html/template"
	"os"
	"path"
)

type TypeProperties struct {
	Name     string
	JsonName string
	Type     string
}

type TypeData struct {
	TypeName   string
	Properties []TypeProperties
}

type PostRPC struct {
	Path            string
	RequestType     TypeData
	RequestTypeName string
}

func GetDataToWrite(data any, tmpl string) (string, error) {
	var buf bytes.Buffer
	t := template.Must(template.New("tmpl").Parse(tmpl))

	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func WriteToPackage(data string, rootPath string, fileName string, fileHeader string) error {
	dir := path.Join(rootPath, "client")
	filePath := path.Join(dir, fileName)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(fileHeader)
	file.WriteString(data)
	return nil
}
