package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
)

var staticComponents = openapi3.Components{
	Parameters: map[string]*openapi3.ParameterRef{
		"filter": {
			Value: &openapi3.Parameter{
				Description: "[Filter](http://rest-layer.io/#filtering) which entries to show. Allows a MongoDB-like query syntax.",
				Name:        "filter",
				In:          "query",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "string",
					},
				},
			},
		},
		"fields": {
			Value: &openapi3.Parameter{
				Description: "[Select](http://rest-layer.io/#field-selection) which fields to show, including [embedding](http://rest-layer.io/#embedding) of related resources.",
				Name:        "fields",
				In:          "query",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "string",
					},
				},
			},
		},
		"limit": {
			Value: &openapi3.Parameter{
				Description: "Limit maximum entries per [page](http://rest-layer.io/#paginatio).",
				Name:        "limit",
				In:          "query",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "integer",
						Min:  openapi3.Float64Ptr(0),
					},
				},
			},
		},
		"skip": {
			Value: &openapi3.Parameter{
				Description: "[Skip](http://rest-layer.io/#skipping) the first N entries.",
				Name:        "skip",
				In:          "query",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "integer",
						Min:  openapi3.Float64Ptr(0),
					},
				},
			},
		},
		"page": {
			Value: &openapi3.Parameter{
				Description: "The [page](http://rest-layer.io/#pagination) number to display, starting at 1.",
				Name:        "page",
				In:          "query",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:    "integer",
						Default: openapi3.Float64Ptr(1),
						Min:     openapi3.Float64Ptr(1),
					},
				},
			},
		},
		"total": {
			Value: &openapi3.Parameter{
				Description: "Force total number of entries to be included in the response header. This could have performance implications.",
				Name:        "total",
				In:          "query",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						/*
							Type: "integer",
							Default:  openapi3.Float64Ptr(0),
						*/
						Type:    "boolean",
						Default: false,
					},
				},
			},
		},
	},
	Headers: map[string]*openapi3.HeaderRef{
		"Date": {
			Value: &openapi3.Header{
				Description: "The time this request was served.",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:   "string",
						Format: "date-time",
					},
				},
			},
		},
		"Etag": {
			Value: &openapi3.Header{
				Description: "Provides [concurrency-control](http://rest-layer.io/#data-integrity-and-concurrency-control) down to the storage layer.",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "string",
					},
				},
			},
		},
		"Last-Modified": {
			Value: &openapi3.Header{
				Description: "When this resource was last modified.",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:   "string",
						//Format: "date-time",
					},
				},
			},
		},
		"X-Total": {
			Value: &openapi3.Header{
				Description: "Total number of entries matching the supplied filter.",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "integer",
					},
				},
			},
		},
	},
	Schemas: map[string]*openapi3.SchemaRef{
		"Error": {
			Value: &openapi3.Schema{
				Type:     "object",
				Required: []string{"code", "message"},
				Properties: map[string]*openapi3.SchemaRef{
					"code": &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Description: "HTTP Status code",
							Type:        "integer",
						},
					},
					"message": &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Description: "Error message",
							Type:        "string",
						},
					},
				},
			},
		},
		"ValidationError": {
			Value: &openapi3.Schema{
				Type:     "object",
				Required: []string{"code", "message"},
				Properties: map[string]*openapi3.SchemaRef{
					"code": &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Description: "HTTP Status code",
							Type:        "integer",
						},
					},
					"message": &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Description: "Error message",
							Type:        "string",
						},
					},
					/* TODO: complete
					"issues": &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Description: "Error message",
							Type:        "array",
							Items: &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: "string",
								},
							},
						},
					},
					*/
				},
			},
		},
	},
	Responses: map[string]*openapi3.ResponseRef{
		"Error": &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: "Error",
				Content: map[string]*openapi3.MediaType{
					"application/json": &openapi3.MediaType{
						Schema: &openapi3.SchemaRef{
							Ref: "#/components/schemas/Error",
						},
					},
				},
			},
		},
		"ValidationError": &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: "Validation Error",
				Content: map[string]*openapi3.MediaType{
					"application/json": &openapi3.MediaType{
						Schema: &openapi3.SchemaRef{
							Ref: "#/components/schemas/ValidationError",
						},
					},
				},
			},
		},
	},
}