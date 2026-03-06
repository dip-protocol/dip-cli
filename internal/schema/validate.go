package schema

import (
	"fmt"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
)

func ValidateArtifact(artifactPath string, schemaPath string) error {

	// Convert Windows paths to forward slash format
	artifactPath = filepath.ToSlash(artifactPath)
	schemaPath = filepath.ToSlash(schemaPath)

	artifactURL := "file:///" + artifactPath
	schemaURL := "file:///" + schemaPath

	artifactLoader := gojsonschema.NewReferenceLoader(artifactURL)
	schemaLoader := gojsonschema.NewReferenceLoader(schemaURL)

	result, err := gojsonschema.Validate(schemaLoader, artifactLoader)

	if err != nil {
		return err
	}

	if result.Valid() {

		fmt.Println("Schema validation successful")
		return nil

	} else {

		fmt.Println("Schema validation failed:")

		for _, desc := range result.Errors() {
			fmt.Println("-", desc)
		}

		return fmt.Errorf("artifact schema invalid")
	}
}