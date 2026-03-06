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

	schemaLoader := gojsonschema.NewReferenceLoader("file:///" + filepath.ToSlash(absSchema))
	documentLoader := gojsonschema.NewReferenceLoader("file:///" + filepath.ToSlash(absRecord))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		for _, desc := range result.Errors() {
			return fmt.Errorf(desc.String())
		}
	}

	return nil
}