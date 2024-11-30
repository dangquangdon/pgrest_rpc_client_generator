package generator

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func SnakeToCamel(snake string) string {
	// Create a Title casing transformer for proper Unicode handling
	caser := cases.Title(language.Und)
	// Split the string by underscores
	parts := strings.Split(snake, "_")
	// Capitalize each part
	for i := range parts {
		parts[i] = caser.String(parts[i])
	}
	// Join the parts to form CamelCase
	return strings.Join(parts, "")
}

func GetDataFromMapByKey(data map[string]interface{}, key string) (map[string]interface{}, error) {
	value, ok := data[key]
	if !ok {
		errMessage := fmt.Sprintf("%s does not exist in %s", key, data)
		return nil, errors.New(errMessage)
	}

	return value.(map[string]interface{}), nil
}
