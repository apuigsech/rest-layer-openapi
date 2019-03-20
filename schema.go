package openapi

import (
	"fmt"
	"reflect"
	"github.com/rs/rest-layer/schema"
	"github.com/getkin/kin-openapi/openapi3"
)

func generateSchema(s schema.Schema) *openapi3.Schema {
	ret := &openapi3.Schema{
		Type:        "object",
		Description: s.Description,
		Properties:  map[string]*openapi3.SchemaRef{},
	}

	for fieldName, field := range s.Fields {
		ret.Properties[fieldName] = &openapi3.SchemaRef{}
		ret.Properties[fieldName].Value = generateSchemaFromField(field)
	}

	return ret
}

func generateSchemaFromField(field schema.Field) *openapi3.Schema {
	switch t := field.Validator.(type) {
	case *schema.String:
		return generateSchemaFromFieldString(field)
	case *schema.Array:
		return generateSchemaFromFieldArray(field)
	case *schema.Reference:
		return generateSchemaFromFieldReference(field)
	default:
		fmt.Println("TYPE > ", reflect.TypeOf(t))
		return nil
	}
	return nil
}

func generateSchemaFromFieldString(f schema.Field) *openapi3.Schema {
	v := f.Validator.(*schema.String)
	ret := &openapi3.Schema{
		Type:      "string",
		MinLength: uint64(v.MinLen),
		Pattern:   v.Regexp,
		//Enum: []interface{}(v.Allowed),
	}
	if v.MaxLen > 0 {
		ret.MaxLength = openapi3.Uint64Ptr(uint64(v.MaxLen))
	}

	return ret
}

func generateSchemaFromFieldArray(f schema.Field) *openapi3.Schema {
	v := f.Validator.(*schema.Array)
	ret := &openapi3.Schema{
		Type:     "array",
		MinItems: uint64(v.MinLen),
		Items: &openapi3.SchemaRef{
			Value: generateSchemaFromField(v.Values),
		},
	}
	if v.MaxLen > 0 {
		ret.MaxItems = openapi3.Uint64Ptr(uint64(v.MaxLen))
	}

	return ret
}

func generateSchemaFromFieldReference(f schema.Field) *openapi3.Schema {
	ret := &openapi3.Schema{
		Type: "string",
	}

	return ret
}