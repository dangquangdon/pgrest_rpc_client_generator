package generator

var JsonTypesToGoMap = map[string]string{
	"string":  "string",
	"number":  "float64",
	"integer": "int",
	"boolean": "bool",
	"null":    "nil",
}

var JsonFormatTypeToGoMap = map[string]string{
	"text": "string",
	"json": "interface{}",
}

const TmplForType string = `{{- range . }}

type {{ .TypeName }} struct {
	BaseType
	{{- range .Properties }}
	{{ .Name }} {{ .Type }} ` + "`json:\"{{ .JsonName }}\"`" + `
	{{- end }}
}
	{{- end }}`

const TmplForRequest string = `

func doRequest[T BaseTypeInterface](data T, path string) (*http.Request, error) {
	rData, err := data.ToReader()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", path, rData)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("User-Agent", "{{ .UserAgent }}")
	
	return req, nil
}

{{- range .Posts }}

func Do{{ .RequestTypeName }}Request(data {{ .RequestTypeName }}) (*http.Request, error) {
	return doRequest(data, "{{ .Path }}")
}
	{{- end }}
	`

const TypesFileHeaderToWrite string = `
package client

import (
	"bytes"
	"encoding/json"
	"io"
)

type BaseTypeInterface interface {
	ToReader() (io.Reader, error)
}

type BaseType struct {}

func (base BaseType) ToReader() (io.Reader, error) {
	data, err := json.Marshal(base)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}
`

const RequestsFileHeaderToWrite string = `
package client

import "net/http"

`
