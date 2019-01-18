package generator

import (
	"bytes"
	"errors"
	"sort"

	gen "github.com/nasa9084/go-genutils"
	openapi "github.com/nasa9084/go-openapi"
	"golang.org/x/tools/imports"
)

var httpMethods = []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE"}

func GenerateHandlers(spec *openapi.Document) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("package main")
	// TODO: import by tools/imports
	buf.WriteString(gen.Imports{gen.Import{ImportPath: "net/http"}, gen.Import{ImportPath: "github.com/gorilla/mux"}}.String())
	buf.WriteString("\n// code generated by genserver. DO NOT EDIT.")
	// Generate NewRouter()
	buf.WriteString("\n\nfunc NewRouter() http.Handler {")
	buf.WriteString("\nr := mux.NewRouter()")
	var paths []string
	for path := range spec.Paths {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		pathItem := spec.Paths[path]
		for _, method := range httpMethods {
			op := pathItem.GetOperationByMethod(method)
			if op == nil {
				continue
			}
			if op.OperationID == "" {
				return nil, errors.New("operationId is required: " + path)
			}
			buf.WriteString("\nr.HandleFunc(\"")
			buf.WriteString(path)
			buf.WriteString(`", `)
			buf.WriteString(op.OperationID)
			buf.WriteString(`Handler).Methods("`)
			buf.WriteString(method)
			buf.WriteString(`")`)
		}
	}
	buf.WriteString("\nreturn r")
	buf.WriteString("\n}")

	//Generate Handlers
	for _, path := range paths {
		pathItem := spec.Paths[path]
		for _, method := range httpMethods {
			op := pathItem.GetOperationByMethod(method)
			if op == nil {
				continue
			}
			if op.OperationID == "" {
				return nil, errors.New("operationId is required: " + path)
			}
			buf.WriteString("\n\nfunc ")
			buf.WriteString(op.OperationID)
			buf.WriteString("Handler(w http.ResponseWriter, r *http.Request) {")
			buf.WriteString("\nst, hdr, res, err := ")
			buf.WriteString(op.OperationID)
			buf.WriteString("(r)")
			buf.WriteString("\nif err != nil {")
			buf.WriteString("\nreturn")
			buf.WriteString("\n}")
			buf.WriteString("\nvar buf bytes.Buffer")
			buf.WriteString("\nif err := json.NewDecoder(&buf).Decode(res); err != nil {")
			buf.WriteString("\nw.WriteHeader(http.StatusInternalServerError)")
			buf.WriteString("\nreturn")
			buf.WriteString("\n}")
			buf.WriteString("\nif _, err := buf.WriteTo(w); err != nil {")
			buf.WriteString("\nw.WriteHeader(http.StatusInternalServerError)")
			buf.WriteString("\nreturn")
			buf.WriteString("\n}")
			buf.WriteString("\nfor k, v := range hdr {")
			buf.WriteString("\nw.Header().Add(k, v)")
			buf.WriteString("\n}")
			buf.WriteString("\nw.WriteHeader(st)")
			buf.WriteString("\n}")
		}
	}

	src, err := imports.Process("", buf.Bytes(), &imports.Options{Comments: true, Fragment: true, FormatOnly: true})
	if err != nil {
		return nil, err
	}
	return src, nil
}
