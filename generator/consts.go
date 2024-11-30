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

const TmplForRequest string = `{{- range . }}

func Do{{ .RequestTypeName }}Request(data {{ .RequestTypeName }}) (*http.Request, error) {
	rData, err := data.ToReader()
	if err != nil {
		return nil, err
	}
	return http.NewRequest("POST", "{{ .Path }}", rData)
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

type BaseType struct {}

func (base *BaseType) ToReader() (io.Reader, error) {
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
