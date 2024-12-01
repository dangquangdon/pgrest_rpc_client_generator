package generator

import (
	"log"
	"strings"
)

type Generator struct {
	RawData      map[string]interface{}
	RequestTypes []TypeData
	Posts        []PostRPC
	RootPath     string
	BaseUrl      string
	ClientID     string
}

func NewGeneartor(url string, rootPath string, clientId string) *Generator {
	raw, err := ReadOpenApiResponse(url, clientId)
	if err != nil {
		log.Fatalf("failed to process raw data: %s", err)
	}

	paths, ok := raw["paths"]
	if !ok {
		log.Fatal("no `paths` key found.")
	}

	var baseUrl string
	if strings.HasSuffix(url, "/") {
		baseUrl = strings.ReplaceAll(url, "/", "")
	} else {
		baseUrl = url
	}

	types, posts := GetRPCRequestDataTypes(paths.(map[string]interface{}), baseUrl)

	return &Generator{
		RequestTypes: types,
		RawData:      raw,
		Posts:        posts,
		RootPath:     rootPath,
		BaseUrl:      baseUrl,
		ClientID:     clientId,
	}
}

func (gen Generator) GenerateTypes() error {
	types, err := GetDataToWrite(gen.RequestTypes, TmplForType)
	if err != nil {
		return err
	}
	return WriteToPackage(types, gen.RootPath, "types.go", TypesFileHeaderToWrite)
}

func (gen Generator) GenerateRequests() error {
	type WriteContext struct {
		UserAgent string
		Posts     []PostRPC
	}

	posts, err := GetDataToWrite(WriteContext{
		UserAgent: gen.ClientID,
		Posts:     gen.Posts,
	}, TmplForRequest)
	if err != nil {
		return err
	}

	return WriteToPackage(posts, gen.RootPath, "requests.go", RequestsFileHeaderToWrite)
}
