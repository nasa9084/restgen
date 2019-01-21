package generator

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	openapi "github.com/nasa9084/go-openapi"
	"golang.org/x/tools/imports"
)

// GenerateSchemaTypes generates Go's type definitions from .components.schemas.
// The generated source is formatted by goimports.
func GenerateSchemaTypes(spec *openapi.Document) ([]byte, error) {
	components := spec.Components
	if components == nil {
		return nil, nil
	}
	schemas := components.Schemas
	if schemas == nil {
		return nil, nil
	}
	var schemaNames []string
	for name := range schemas {
		schemaNames = append(schemaNames, name)
	}
	sort.Strings(schemaNames)

	var buf bytes.Buffer
	buf.WriteString("package models")
	buf.WriteString("\n// code generated by genserver. DO NOT EDIT.")
	for _, name := range schemaNames {
		schema := schemas[name]
		t, err := generateSchemaType(name, schema)
		if err != nil {
			return nil, err
		}
		buf.Write(t)
	}
	src, err := imports.Process("", buf.Bytes(), &imports.Options{Fragment: true, Comments: true})
	if err != nil {
		return nil, err
	}
	return src, nil
}

func generateSchemaType(name string, schema *openapi.Schema) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("\n")
	if schema.Description != "" {
		for _, line := range strings.Split(schema.Description, "\n") {
			fmt.Fprintf(&buf, "\n// %s", line)
		}
	}
	fmt.Fprintf(&buf, "\ntype %s %s", name, (*OpenAPISchema)(schema).GoType())
	return buf.Bytes(), nil
}

// GenerateResponseTypes generates Go's type definitions from .components.responses.
// The generated source is formatted by goimports.
func GenerateResponseTypes(spec *openapi.Document) ([]byte, error) {
	components := spec.Components
	if components == nil {
		return nil, nil
	}
	responses := components.Responses
	if responses == nil {
		return nil, nil
	}
	var responseNames []string
	for name := range responses {
		responseNames = append(responseNames, name)
	}
	sort.Strings(responseNames)

	var buf bytes.Buffer
	buf.WriteString("package models")
	buf.WriteString("\n// code generated by genserver. DO NOT EDIT.")
	for _, name := range responseNames {
		response := responses[name]
		t, err := generateResponseType(name, response)
		if err != nil {
			return nil, err
		}
		buf.Write(t)
	}
	src, err := imports.Process("", buf.Bytes(), &imports.Options{Fragment: true, Comments: true})
	if err != nil {
		return nil, err
	}
	return src, nil
}

func generateResponseType(name string, response *openapi.Response) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("\n")
	for _, line := range strings.Split(response.Description, "\n") {
		fmt.Fprintf(&buf, "\n// %s", line)
	}
	fmt.Fprintf(&buf, "\ntype %s", name)
	content := response.Content
	if content == nil {
		return nil, fmt.Errorf(".components.responses[%s].content is nil", name)
	}
	mediaType, ok := content["application/json"]
	if !ok {
		return nil, fmt.Errorf(".components.responses[%s].content[application/json] is nil", name)
	}
	buf.WriteString(" ")
	buf.WriteString((*OpenAPISchema)(mediaType.Schema).GoType())
	return buf.Bytes(), nil
}

// GenerateRequestTypes generates Go's type definitions from .components.requestBodies.
// The generated source is formatted by goimports.
func GenerateRequestTypes(spec *openapi.Document) ([]byte, error) {
	components := spec.Components
	if components == nil {
		return nil, nil
	}
	requestBodies := components.RequestBodies
	if requestBodies == nil {
		return nil, nil
	}
	var requestNames []string
	for name := range requestBodies {
		requestNames = append(requestNames, name)
	}
	sort.Strings(requestNames)
	var buf bytes.Buffer
	buf.WriteString("package models")
	buf.WriteString("\n// code generated by genserver. DO NOT EDIT.")
	for _, name := range requestNames {
		requestBody := requestBodies[name]
		t, err := generateRequestType(name, requestBody)
		if err != nil {
			return nil, err
		}
		buf.Write(t)
	}
	src, err := imports.Process("", buf.Bytes(), &imports.Options{Fragment: true, Comments: true})
	if err != nil {
		return nil, err
	}
	return src, nil
}

func generateRequestType(name string, requestBody *openapi.RequestBody) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("\n")
	for _, line := range strings.Split(requestBody.Description, "\n") {
		fmt.Fprintf(&buf, "\n// %s", line)
	}
	fmt.Fprintf(&buf, "\ntype %s", name)
	// requestBody.Content is required parameter and checked by Validate().
	mediaType, ok := requestBody.Content["application/json"]
	if !ok {
		return nil, fmt.Errorf(".components.requestBodies[%s].content[application/json] is nil", name)
	}
	buf.WriteString(" ")
	buf.WriteString((*OpenAPISchema)(mediaType.Schema).GoType())
	return buf.Bytes(), nil
}
