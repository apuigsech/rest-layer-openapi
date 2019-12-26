// +build go1.7

package main

import (
	"fmt"
	"net/url"

	"github.com/rs/rest-layer/resource/testing/mem"
	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/schema"
	"github.com/rs/rest-layer/schema/query"

	"github.com/apuigsech/rest-layer-openapi"
	"github.com/getkin/kin-openapi/openapi3"
)

var (
	// Define a user resource schema
	user = schema.Schema{
		Fields: schema.Fields{
			"id": {
				Required: true,
				// The Filterable and Sortable allows usage of filter and sort
				// on this field in requests.
				Filterable: true,
				Sortable:   true,
				Validator: &schema.String{
					Regexp: "^[0-9a-z]{2,20}$",
				},
			},
			"created": {
				Required:   true,
				ReadOnly:   true,
				Filterable: true,
				Sortable:   true,
				OnInit:     schema.Now,
				Validator:  &schema.Time{},
			},
			"updated": {
				Required:   true,
				ReadOnly:   true,
				Filterable: true,
				Sortable:   true,
				OnInit:     schema.Now,
				// The OnUpdate hook is called when the item is edited. Here we use
				// provided Now hook which just return the current time.
				OnUpdate:  schema.Now,
				Validator: &schema.Time{},
			},
			// Define a name field as required with a string validator
			"name": {
				Required:   true,
				Filterable: true,
				Validator: &schema.String{
					MaxLen: 150,
				},
			},
		},
	}

	// Define a post resource schema
	post = schema.Schema{
		Fields: schema.Fields{
			// schema.*Field are shortcuts for common fields (identical to users' same fields)
			"id":      schema.IDField,
			"created": schema.CreatedField,
			"updated": schema.UpdatedField,
			// Define a user field which references the user owning the post.
			// See bellow, the content of this field is enforced by the fact
			// that posts is a sub-resource of users.
			"user": {
				Required:   true,
				Filterable: true,
				ReadOnly:   true,
				Validator: &schema.Reference{
					Path: "users",
				},
			},
			"published": {
				Filterable: true,
				Default:    false,
				Validator:  &schema.Bool{},
			},
			"title": {
				Required: true,
				Validator: &schema.String{
					MaxLen: 150,
				},
				// Dependency defines that body field can't be changed if
				// the published field is not "false".
				Dependency: query.MustParsePredicate(`{published: false}`),
			},
			"body": {
				Validator: &schema.String{
					MaxLen: 100000,
				},
				Dependency: query.MustParsePredicate(`{published: false}`),
			},
		},
	}
)

func main() {
	// Create a REST API resource index
	index := resource.NewIndex()

	// Add a resource on /users[/:user_id]
	users := index.Bind("users", user, mem.NewHandler(), resource.Conf{
		// We allow all REST methods
		// (rest.ReadWrite is a shortcut for []resource.Mode{resource.Create, resource.Read, resource.Update, resource.Delete, resource,List})
		AllowedModes: resource.ReadWrite,
	})

	// Bind a sub resource on /users/:user_id/posts[/:post_id]
	// and reference the user on each post using the "user" field of the posts resource.
	posts := users.Bind("posts", "user", post, mem.NewHandler(), resource.Conf{
		AllowedModes: resource.ReadWrite,
	})

	// Add a friendly alias to public posts
	// (equivalent to /users/:user_id/posts?filter={"published":true})
	posts.Alias("public", url.Values{"filter": []string{"{\"published\":true}"}})


	doc := openapi.NewOpenapiFromIndex(index, &openapi3.Info{
		Title:   "ApiName",
		Version: "ApiVersion",
	})
	b, _ := doc.MarshalJSON()
	fmt.Println(string(b))
}
