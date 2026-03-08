package schema

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// ValidateArtifact validates a DIP artifact against a JSON schema.
//
// This ensures the artifact structure conforms to the DIP protocol
// before signing or verification.
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
		return nil
	}

	var errors []string

	for _, desc := range result.Errors() {
		errors = append(errors, desc.String())
	}

	return fmt.Errorf(
		"artifact schema validation failed: %s",
		strings.Join(errors, "; "),
	)
}