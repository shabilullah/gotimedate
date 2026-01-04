package docs

import (
	"testing"
)

func TestSwaggerInfo(t *testing.T) {
	if SwaggerInfo.Title == "" {
		t.Error("Swagger title should not be empty")
	}
	if SwaggerInfo.Version == "" {
		t.Error("Swagger version should not be empty")
	}
	if SwaggerInfo.BasePath == "" {
		t.Error("Swagger base path should not be empty")
	}
}
