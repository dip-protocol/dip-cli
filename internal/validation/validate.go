package validation

import (
	"fmt"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
)

func Validate(recordPath string, schemaPath string) error {

	absRecord, err := filepath.Abs(recordPath)
	if err != nil {
		return err
	}

	absSchema, err := filepath.Abs(schemaPath)
	if err != nil {
		return err
	}

	// Convert Windows paths to URI-compatible format
	absRecord = filepath.ToSlash(absRecord)
	absSchema = filepath.ToSlash(absSchema)

	schemaLoader := gojsonschema.NewReferenceLoader("file:///" + absSchema)
	documentLoader := gojsonschema.NewReferenceLoader("file:///" + absRecord)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	if result.Valid() {
		fmt.Println("Record is valid according to schema")
	} else {
		fmt.Println("Record is invalid:")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}

	return nil
}