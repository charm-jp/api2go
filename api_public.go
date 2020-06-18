package api2go

import (
	"context"
	"gitlab.dev.charm.internal/charm/api2go/jsonapi"
	"gitlab.dev.charm.internal/charm/api2go/routing"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

// HandlerFunc for api2go beforeMiddlewares
type BeforeHandlerFunc func(context.Context, http.ResponseWriter, *http.Request) (context.Context, error)
type AfterHandlerFunc func(context.Context, http.ResponseWriter, *http.Request, error) context.Context

// API is a REST JSONAPI.
type API struct {
	ContentType       string
	router            routing.Routeable
	info              information
	resources         []resource
	beforeMiddlewares []BeforeHandlerFunc
	afterMiddlewares  []AfterHandlerFunc
	contextPool       sync.Pool
	URLResolver       URLResolver
}

// Handler returns the http.Handler instance for the API.
func (api API) Handler() http.Handler {
	return api.router.Handler()
}

// Lookup the model associated with a named resource
func (api API) LookupResourceType(name string) reflect.Type {
	for _, resourceRecord := range api.resources {
		if resourceRecord.name == name {
			return resourceRecord.resourceType
		}
	}

	return nil
}

//Router returns the specified router on an api instance
func (api API) Router() routing.Routeable {
	return api.router
}

// AddResource registers a data source for the given resource
// At least the CRUD interface must be implemented, all the other interfaces are optional.
// `resource` should be either an empty struct instance such as `Post{}` or a pointer to
// a struct such as `&Post{}`. The same type will be used for constructing new elements.
func (api *API) AddResource(prototype jsonapi.MarshalIdentifier, source interface{}) {
	api.addResource(prototype, source)
}

// UseBeforeMiddleware registers beforeMiddlewares that implement the api2go.HandlerFunc
// Middleware is run before any generated routes.
func (api *API) UseBeforeMiddleware(middleware ...BeforeHandlerFunc) {
	api.beforeMiddlewares = append(api.beforeMiddlewares, middleware...)
}

// UseAfterMiddleware registers AfterMiddlewares that implement the api2go.HandlerFunc
// Middleware is run After any generated routes.
func (api *API) UseAfterMiddleware(middleware ...AfterHandlerFunc) {
	api.afterMiddlewares = append(api.afterMiddlewares, middleware...)
}

// NewAPIVersion can be used to chain an additional API version to the routing of a previous
// one. Use this if you have multiple version prefixes and want to combine all
// your different API versions. This reuses the baseURL or URLResolver
func (api *API) NewAPIVersion(prefix string) *API {
	return newAPI(prefix, api.info.resolver, api.router)
}

// NewAPIWithResolver can be used to create an API with a custom URL resolver.
func NewAPIWithResolver(prefix string, resolver URLResolver) *API {
	handler := notAllowedHandler{}
	r := routing.NewHTTPRouter(prefix, &handler)
	api := newAPI(prefix, resolver, r)
	handler.API = api
	return api
}

// NewAPIWithBaseURL does the same as NewAPI with the addition of
// a baseURL which get's added in front of all generated URLs.
// For example http://localhost/v1/myResource/abc instead of /v1/myResource/abc
func NewAPIWithBaseURL(prefix string, baseURL string) *API {
	handler := notAllowedHandler{}
	staticResolver := NewStaticResolver(baseURL)
	r := routing.NewHTTPRouter(prefix, &handler)
	api := newAPI(prefix, staticResolver, r)
	handler.API = api
	return api
}

// NewAPI returns an initialized API instance
// `prefix` is added in front of all endpoints.
func NewAPI(prefix string) *API {
	handler := notAllowedHandler{}
	staticResolver := NewStaticResolver("")
	r := routing.NewHTTPRouter(prefix, &handler)
	api := newAPI(prefix, staticResolver, r)
	handler.API = api
	return api
}

// NewAPIWithRouting allows you to use a custom URLResolver, marshalers and custom routing
// if you want to use the default routing, you should use another constructor.
//
// If you don't need any of the parameters you can skip them with the defaults:
// the default for `prefix` would be `""`, which means there is no namespace for your api.
// although we suggest using one.
//
// if your api only answers to one url you can use a NewStaticResolver() as  `resolver`
func NewAPIWithRouting(prefix string, resolver URLResolver, router routing.Routeable) *API {
	return newAPI(prefix, resolver, router)
}

// newAPI is now an internal method that can be changed if params are changing
func newAPI(prefix string, resolver URLResolver, router routing.Routeable) *API {
	// Add initial and trailing slash to prefix
	prefixSlashes := strings.Trim(prefix, "/")
	if len(prefixSlashes) > 0 {
		prefixSlashes = "/" + prefixSlashes + "/"
	} else {
		prefixSlashes = "/"
	}

	info := information{prefix: prefix, resolver: resolver}

	api := &API{
		ContentType:       defaultContentTypHeader,
		router:            router,
		info:              info,
		beforeMiddlewares: make([]BeforeHandlerFunc, 0),
		afterMiddlewares:  make([]AfterHandlerFunc, 0),
		URLResolver:       info.resolver,
	}

	return api
}
