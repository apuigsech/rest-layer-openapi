package openapi

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jinzhu/inflection"
	"github.com/rs/rest-layer/resource"
	"strings"
)

func addResource(doc *openapi3.Swagger, prevRscList []*resource.Resource, rsc *resource.Resource) {
	schemaNamePlural := rsc.Name()
	schemaNameSingular := inflection.Singular(rsc.Name())
	schemaIdParameter := schemaNameSingular + "Id"

	doc.Components.Schemas[schemaNameSingular] = &openapi3.SchemaRef{
		Value: generateSchema(rsc.Schema()),
	}

	doc.Components.Parameters[schemaIdParameter] = &openapi3.ParameterRef{
		Value: &openapi3.Parameter{
			Name:        schemaIdParameter,
			Description: fmt.Sprintf("The %s's ID", schemaNameSingular),
			In:          "path",
			Required:    true,
			Schema: &openapi3.SchemaRef{
				Ref: fmt.Sprintf("#/components/schemas/%s/properties/id", schemaNameSingular),
			},
		},
	}

	var path string
	var operationSufix string
	var params []*openapi3.ParameterRef
	for _, prevRsc := range prevRscList {
		prevSchemaNameSingular := inflection.Singular(prevRsc.Name())
		prevSchemaIdParameter := prevSchemaNameSingular + "Id"

		path = path + fmt.Sprintf("/%s/{%s}", prevSchemaNameSingular, prevSchemaIdParameter)
		operationSufix = operationSufix + fmt.Sprintf("On%s", strings.Title(prevSchemaNameSingular))

		param := &openapi3.ParameterRef{
			Ref: fmt.Sprintf("#/components/parameters/%s", prevSchemaIdParameter),
		}

		params = append(params, param)
	}

	path = path + fmt.Sprintf("/%s", schemaNamePlural)

	if rsc.Conf().IsModeAllowed(resource.List) {
		op := &openapi3.Operation{
			OperationID: "List" + strings.Title(schemaNamePlural) + operationSufix,
			Parameters: append(
				[]*openapi3.ParameterRef{
					{Ref: "#/components/parameters/filter"},
					{Ref: "#/components/parameters/fields"},
					{Ref: "#/components/parameters/limit"},
					{Ref: "#/components/parameters/page"},
					{Ref: "#/components/parameters/skip"},
					{Ref: "#/components/parameters/total"},
				},
				params...
			),
			Responses: map[string]*openapi3.ResponseRef{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: fmt.Sprintf("List of %s", schemaNamePlural),
						Headers: map[string]*openapi3.HeaderRef{
							"Date":    {Ref: "#/components/headers/Date"}, // TODO: Verify
							"X-Total": {Ref: "#/components/headers/X-Total"},
						},
						Content: map[string]*openapi3.MediaType{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: "array",
										Items: &openapi3.SchemaRef{
											Ref: fmt.Sprintf("#/components/schemas/%s", schemaNameSingular),
										},
									},
								},
							},
						},
					},
				},
				"default": &openapi3.ResponseRef{
					Ref: "#/components/responses/Error",
				},
			},
		}
		doc.AddOperation(path, "GET", op)
	}

	if rsc.Conf().IsModeAllowed(resource.Create) {
		op := &openapi3.Operation{
			OperationID: "Create" + strings.Title(schemaNameSingular) + operationSufix,
			Parameters: append(
				[]*openapi3.ParameterRef{},
				params...
			),
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Content: map[string]*openapi3.MediaType{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{
								Ref: fmt.Sprintf("#/components/schemas/%s", schemaNameSingular),
							},
						},
					},
				},
			},
			Responses: map[string]*openapi3.ResponseRef{
				"201": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: fmt.Sprintf("Create %s", schemaNameSingular),
						Headers: map[string]*openapi3.HeaderRef{
							"Etag":          {Ref: "#/components/headers/Etag"},
							"Last-Modified": {Ref: "#/components/headers/Last-Modified"},
						},
						Content: map[string]*openapi3.MediaType{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: "array",
										Items: &openapi3.SchemaRef{
											Ref: fmt.Sprintf("#/components/schemas/%s", schemaNameSingular),
										},
									},
								},
							},
						},
					},
				},
				"422": &openapi3.ResponseRef{
					Ref: "#/components/responses/ValidationError",
				},
				"default": &openapi3.ResponseRef{
					Ref: "#/components/responses/Error",
				},
			},
		}
		doc.AddOperation(path, "POST", op)
	}

	if rsc.Conf().IsModeAllowed(resource.Clear) {
		op := &openapi3.Operation{
			OperationID: "Clear" + strings.Title(rsc.Name()) + operationSufix,
			Parameters: append(
				[]*openapi3.ParameterRef{
					{Ref: "#/components/parameters/filter"},
				},
				params...
			),
			Responses: map[string]*openapi3.ResponseRef{
				"204": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: fmt.Sprintf("Clear %s", rsc.Name()),
						Headers: map[string]*openapi3.HeaderRef{
							"Date":    {Ref: "#/components/headers/Date"}, // TODO: Verify
							"X-Total": {Ref: "#/components/headers/X-Total"},
						},
					},
				},
				"default": &openapi3.ResponseRef{
					Ref: "#/components/responses/Error",
				},
			},
		}
		doc.AddOperation(path, "DELETE", op)
	}

	path = path + fmt.Sprintf("/{%s}", schemaIdParameter)

	if rsc.Conf().IsModeAllowed(resource.Read) {
		op := &openapi3.Operation{
			OperationID: "Read" + strings.Title(schemaNameSingular) + operationSufix,
			Parameters: append(
				[]*openapi3.ParameterRef{
					{Ref: "#/components/parameters/fields"},
					{Ref: fmt.Sprintf("#/components/parameters/%s", schemaIdParameter)},
				},
				params...
			),
			Responses: map[string]*openapi3.ResponseRef{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: fmt.Sprintf("Get %s", rsc.Name()),
						Headers: map[string]*openapi3.HeaderRef{
							"Date":    {Ref: "#/components/headers/Date"}, // TODO: Verify
							"X-Total": {Ref: "#/components/headers/X-Total"},
						},
						Content: map[string]*openapi3.MediaType{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: "array",
										Items: &openapi3.SchemaRef{
											Ref: fmt.Sprintf("#/components/schemas/%s", schemaNameSingular),
										},
									},
								},
							},
						},
					},
				},
				"default": &openapi3.ResponseRef{
					Ref: "#/components/responses/Error",
				},
			},
		}
		doc.AddOperation(path, "GET", op)
	}

	if rsc.Conf().IsModeAllowed(resource.Replace) {
		op := &openapi3.Operation{
			OperationID: "Replace" + strings.Title(schemaNameSingular) + operationSufix,
			Parameters: append(
				[]*openapi3.ParameterRef{
					{Ref: fmt.Sprintf("#/components/parameters/%s", schemaIdParameter)},
				},
				params...
			),
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Content: map[string]*openapi3.MediaType{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{
								Ref: fmt.Sprintf("#/components/schemas/%s", schemaNameSingular),
							},
						},
					},
				},
			},
			Responses: map[string]*openapi3.ResponseRef{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: fmt.Sprintf("Replace %s", rsc.Name()),
						Headers: map[string]*openapi3.HeaderRef{
							"Etag":          {Ref: "#/components/headers/Etag"},
							"Last-Modified": {Ref: "#/components/headers/Last-Modified"},
						},
						Content: map[string]*openapi3.MediaType{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: "array",
										Items: &openapi3.SchemaRef{
											Ref: fmt.Sprintf("#/components/schemas/%s", schemaNameSingular),
										},
									},
								},
							},
						},
					},
				},
				"422": &openapi3.ResponseRef{
					Ref: "#/components/responses/ValidationError",
				},
				"default": &openapi3.ResponseRef{
					Ref: "#/components/responses/Error",
				},
			},
		}
		doc.AddOperation(path, "PUT", op)
	}

	if rsc.Conf().IsModeAllowed(resource.Update) {
		op := &openapi3.Operation{
			OperationID: "Update" + strings.Title(schemaNameSingular) + operationSufix,
			Parameters: append(
				[]*openapi3.ParameterRef{
					{Ref: fmt.Sprintf("#/components/parameters/%s", schemaIdParameter)},
				},
				params...
			),
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Content: map[string]*openapi3.MediaType{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{
								Ref: fmt.Sprintf("#/components/schemas/%s", schemaNameSingular),
							},
						},
					},
				},
			},
			Responses: map[string]*openapi3.ResponseRef{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: fmt.Sprintf("Update %s", rsc.Name()),
						Headers: map[string]*openapi3.HeaderRef{
							"Etag":          {Ref: "#/components/headers/Etag"},
							"Last-Modified": {Ref: "#/components/headers/Last-Modified"},
						},
						Content: map[string]*openapi3.MediaType{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: "array",
										Items: &openapi3.SchemaRef{
											Ref: fmt.Sprintf("#/components/schemas/%s", schemaNameSingular),
										},
									},
								},
							},
						},
					},
				},
				"422": &openapi3.ResponseRef{
					Ref: "#/components/responses/ValidationError",
				},
				"default": &openapi3.ResponseRef{
					Ref: "#/components/responses/Error",
				},
			},
		}
		doc.AddOperation(path, "PATCH", op)
	}

	if rsc.Conf().IsModeAllowed(resource.Delete) {
		op := &openapi3.Operation{
			OperationID: "Delete" + strings.Title(schemaNameSingular) + operationSufix,
			Parameters: append(
				[]*openapi3.ParameterRef{
					{Ref: fmt.Sprintf("#/components/parameters/%s", schemaIdParameter)},
				},
				params...
			),
			Responses: map[string]*openapi3.ResponseRef{
				"204": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: fmt.Sprintf("Delete %s", rsc.Name()),
					},
				},
				"422": &openapi3.ResponseRef{
					Ref: "#/components/responses/ValidationError",
				},
				"default": &openapi3.ResponseRef{
					Ref: "#/components/responses/Error",
				},
			},
		}
		doc.AddOperation(path, "DELETE", op)
	}

	for _, subRsc := range rsc.GetResources() {
		addResource(doc, append(prevRscList, rsc), subRsc)
	}
}
