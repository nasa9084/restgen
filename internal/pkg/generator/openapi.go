package generator

import (
	"sort"
	"strings"
	"sync"

	openapi "github.com/nasa9084/go-openapi"
)

var typeMap sync.Map

// RegisterType registers type and format to go-type mapping for schema.GoType function.
func RegisterType(typ, format, gotype string) {
	formatMap, ok := typeMap.Load(typ)
	if !ok {
		formatMap = &sync.Map{}
		typeMap.Store(typ, formatMap)
	}
	formatMap.(*sync.Map).Store(format, gotype)
}

// DeregisterType removes type and format to go-type mapping.
func DeregisterType(typ, format string) {
	formatMap, ok := typeMap.Load(typ)
	if !ok {
		return
	}
	formatMap.(*sync.Map).Delete(format)
}

// LoadType loads registered go-type from type and format.
func LoadType(typ, format string) string {
	formatMap, ok := typeMap.Load(typ)
	if !ok {
		return ""
	}
	gotype, ok := formatMap.(*sync.Map).Load(format)
	if !ok {
		gotype, ok = formatMap.(*sync.Map).Load("")
		if !ok {
			return ""
		}
	}
	return gotype.(string)
}

func init() {
	RegisterType("integer", "", "int")
	RegisterType("integer", "int32", "int32")
	RegisterType("integer", "int64", "int64")
	RegisterType("number", "", "float64")
	RegisterType("number", "float", "float32")
	RegisterType("double", "double", "float64")
	RegisterType("string", "", "string")
	RegisterType("string", "byte", "[]byte")
	RegisterType("string", "binary", "[]byte")
	RegisterType("boolean", "", "bool")
	RegisterType("string", "date", "time.Time")
	RegisterType("string", "date-time", "time.Time")
	RegisterType("string", "password", "string")
}

type OpenAPISchema openapi.Schema

func (schema OpenAPISchema) GoType() string {
	if schema.Ref != "" {
		ref := strings.Split(schema.Ref, "/")
		return ref[len(ref)-1]
	}
	switch schema.Type {
	case "":
		return ""
	case "object":
		var buf strings.Builder
		buf.WriteString("struct {")
		var propnames []string
		for name := range schema.Properties {
			propnames = append(propnames, name)
		}
		sort.Strings(propnames)
		for _, name := range propnames {
			prop := schema.Properties[name]
			// call 4 times WriteString is faster than fmt.Fprintf
			buf.WriteString("\n")
			buf.WriteString(name)
			buf.WriteString(" ")
			buf.WriteString(prop.GoType())
		}
		buf.WriteString("\n}")
		return buf.String()
	case "array":
		return "[]" + schema.Items.GoType()
	default:
		t := LoadType(schema.Type, schema.Format)
		if t != "" {
			return t
		}
		return "interface{}"
	}
}
