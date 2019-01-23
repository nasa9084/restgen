package generator

import (
	"bytes"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"

	gen "github.com/nasa9084/go-genutils"
	openapi "github.com/nasa9084/go-openapi"
	"golang.org/x/tools/imports"
)

var httpMethods = []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE"}

func GenerateHandlers(spec *openapi.Document) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("package main")
	i := gen.Imports{
		gen.Import{ImportPath: "encoding/json"},
		gen.Import{ImportPath: "net/http"},
	}
	buf.WriteString(i.String())
	buf.WriteString("\n// code generated by genserver. DO NOT EDIT.")
	// Generate NewRouter()
	buf.WriteString("\n\nfunc NewRouter() http.Handler {")
	buf.WriteString("\nr := http.NewServeMux()")
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
			buf.WriteString(`Handler)`)
		}
	}
	buf.WriteString("\nreturn r")
	buf.WriteString("\n}")

	// Generate errorStatus()
	buf.WriteString("\n\nfunc errorStatus(err error) int {")
	buf.WriteString("return http.StatusInternalServerError")
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
			resp, st, ok := op.SuccessResponse()
			if !ok {
				return nil, errors.New("there is no success response: " + path)
			}
			if resp.Ref == "" {
				return nil, errors.New("response have to be defined as reference object: " + path)
			}
			if method != http.MethodGet {
				buf.WriteString("\nvar payload model.")
				ref := strings.Split(resp.Ref, "/")
				buf.WriteString(ref[len(ref)-1])
				buf.WriteString("\nif err := json.NewDecoder(r.Body).Decode(&payload); err != nil {")
				buf.WriteString("\nw.WriteHeader(http.StatusBadRequest)")
				buf.WriteString("\n}")
			}
			buf.WriteString("\nhdr, res, err := ")
			buf.WriteString(op.OperationID)
			buf.WriteString("(r")
			if method != http.MethodGet {
				buf.WriteString(", payload")
			}
			buf.WriteString(")")
			buf.WriteString("\nif err != nil {")
			buf.WriteString("\nw.WriteHeader(errorStatus(err))")
			buf.WriteString("\nreturn")
			buf.WriteString("\n}")
			buf.WriteString("\nbuf := getBuffer()")
			buf.WriteString("\ndefer releaseBuffer(buf)")
			buf.WriteString("\nif err := json.NewEncoder(buf).Encode(res); err != nil {")
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
			buf.WriteString("\nw.WriteHeader(")
			buf.WriteString(strconv.Itoa(st))
			buf.WriteString(")")
			buf.WriteString("\n}")
		}
	}
	src, err := imports.Process("", buf.Bytes(), &imports.Options{Fragment: true, Comments: true})
	if err != nil {
		return nil, err
	}
	return src, nil
}
