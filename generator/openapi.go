package generator

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func GetOpenApiSpecFromUrl(url string, clientId string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", clientId)

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf("error: status code %d", response.StatusCode)
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func ReadOpenApiResponse(url string, clientId string) (raw map[string]interface{}, err error) {
	data, err := GetOpenApiSpecFromUrl(url, clientId)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	return
}

func GetRPCRequestDataTypes(paths map[string]interface{}, baseUrl string) ([]TypeData, []PostRPC) {
	types := []TypeData{}
	posts := []PostRPC{}

	for key, val := range paths {
		if strings.HasPrefix(key, "/rpc/") {
			post, err := GetDataFromMapByKey(val.(map[string]interface{}), "post")
			if err != nil {
				log.Printf("%s error: %s. Skip it!", key, err.Error())
				continue
			}

			parameters, ok := post["parameters"]
			if !ok {
				log.Printf("cannot find `parameters` in post of %s\n", key)
				continue
			}

			params := parameters.([]interface{})
			if len(params) < 1 {
				log.Printf("parameters in post of %s is empty\n", key)
				continue
			}

			schema, err := GetDataFromMapByKey(params[0].(map[string]interface{}), "schema")
			if err != nil {
				log.Printf("%s error: %s. Skip it!", key, err.Error())
				continue
			}

			properties, err := GetDataFromMapByKey(schema, "properties")
			if err != nil {
				log.Printf("%s error: %s. Skip it!", key, err.Error())
				continue
			}

			name := SnakeToCamel(strings.ReplaceAll(key, "/rpc/", ""))
			td := TypeData{
				TypeName:   name,
				Properties: []TypeProperties{},
			}
			for pKey, pVal := range properties {
				if pKey == "" {
					continue
				}

				stringMapVal := make(map[string]string)
				for vK, vV := range pVal.(map[string]interface{}) {
					sVal, ok := vV.(string)
					if !ok {
						panic(fmt.Sprintf("%s error: cannot convert %s to string", key, vV))
					}

					stringMapVal[vK] = sVal
				}
				pType, err := GetPropertyType(stringMapVal)
				if err != nil {
					panic(err)
				}

				tp := TypeProperties{
					Name:     SnakeToCamel(pKey),
					JsonName: pKey,
					Type:     pType,
				}
				td.Properties = append(td.Properties, tp)
			}

			types = append(types, td)
			posts = append(posts, PostRPC{
				Path:            baseUrl + key,
				RequestType:     td,
				RequestTypeName: td.TypeName,
			})

		}
	}

	return types, posts
}

func GetPropertyType(prop map[string]string) (string, error) {
	t, ok := prop["type"]
	if !ok {
		// maybe check "format"
		format, ok := prop["format"]
		if !ok {
			return "", fmt.Errorf("no format or type specified")
		}

		return JsonFormatTypeToGoMap[format], nil
	}

	return JsonTypesToGoMap[t], nil
}
