package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/rest-layer/resource"
)

func NewOpenapiFromIndex(index resource.Index, info *openapi3.Info) *openapi3.Swagger {
	doc := &openapi3.Swagger{
		OpenAPI:    "3.0.0",
		Info:       info,
		Components: staticComponents,
	}

	for _, rsc := range index.GetResources() {
		addResource(doc, []*resource.Resource{}, rsc)
	}

	return doc
}