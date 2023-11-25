// Package http provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get the contained authors of the catalog
	// (GET /authors)
	GetAuthors(w http.ResponseWriter, r *http.Request, params GetAuthorsParams)
	// Get the inventory of the catalog
	// (GET /inventory)
	GetInventory(w http.ResponseWriter, r *http.Request, params GetInventoryParams)
	// Get an inventory entry by inventory ID
	// (GET /inventory/{inventoryId})
	GetInventoryById(w http.ResponseWriter, r *http.Request, inventoryId string)
	// Get the versions of an inventory entry
	// (GET /inventory/{inventoryId}/versions)
	GetInventoryVersionsById(w http.ResponseWriter, r *http.Request, inventoryId string)
	// Get the contained manufacturers of the catalog
	// (GET /manufacturers)
	GetManufacturers(w http.ResponseWriter, r *http.Request, params GetManufacturersParams)
	// Get the contained mpns (manufacturer part numbers) of the catalog
	// (GET /mpns)
	GetMpns(w http.ResponseWriter, r *http.Request, params GetMpnsParams)
	// Get the content of a Thing Model by it's ID
	// (GET /thing-models/{tmId})
	GetThingModelById(w http.ResponseWriter, r *http.Request, tmId string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetAuthors operation middleware
func (siw *ServerInterfaceWrapper) GetAuthors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, Api_keyScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAuthorsParams

	// ------------- Optional query parameter "filter[manufacturer]" -------------

	err = runtime.BindQueryParameter("form", true, false, "filter[manufacturer]", r.URL.Query(), &params.FilterManufacturer)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filter[manufacturer]", Err: err})
		return
	}

	// ------------- Optional query parameter "filter[mpn]" -------------

	err = runtime.BindQueryParameter("form", true, false, "filter[mpn]", r.URL.Query(), &params.FilterMpn)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filter[mpn]", Err: err})
		return
	}

	// ------------- Optional query parameter "filter[original]" -------------

	err = runtime.BindQueryParameter("form", true, false, "filter[original]", r.URL.Query(), &params.FilterOriginal)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filter[original]", Err: err})
		return
	}

	// ------------- Optional query parameter "search[content]" -------------

	err = runtime.BindQueryParameter("form", true, false, "search[content]", r.URL.Query(), &params.SearchContent)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "search[content]", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetAuthors(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetInventory operation middleware
func (siw *ServerInterfaceWrapper) GetInventory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, Api_keyScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetInventoryParams

	// ------------- Optional query parameter "filter[author]" -------------

	err = runtime.BindQueryParameter("form", true, false, "filter[author]", r.URL.Query(), &params.FilterAuthor)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filter[author]", Err: err})
		return
	}

	// ------------- Optional query parameter "filter[manufacturer]" -------------

	err = runtime.BindQueryParameter("form", true, false, "filter[manufacturer]", r.URL.Query(), &params.FilterManufacturer)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filter[manufacturer]", Err: err})
		return
	}

	// ------------- Optional query parameter "filter[mpn]" -------------

	err = runtime.BindQueryParameter("form", true, false, "filter[mpn]", r.URL.Query(), &params.FilterMpn)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filter[mpn]", Err: err})
		return
	}

	// ------------- Optional query parameter "filter[original]" -------------

	err = runtime.BindQueryParameter("form", true, false, "filter[original]", r.URL.Query(), &params.FilterOriginal)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filter[original]", Err: err})
		return
	}

	// ------------- Optional query parameter "search[content]" -------------

	err = runtime.BindQueryParameter("form", true, false, "search[content]", r.URL.Query(), &params.SearchContent)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "search[content]", Err: err})
		return
	}

	// ------------- Optional query parameter "sort" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort", r.URL.Query(), &params.Sort)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "sort", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetInventory(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetInventoryById operation middleware
func (siw *ServerInterfaceWrapper) GetInventoryById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "inventoryId" -------------
	var inventoryId string

	err = runtime.BindStyledParameter("simple", false, "inventoryId", mux.Vars(r)["inventoryId"], &inventoryId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "inventoryId", Err: err})
		return
	}

	ctx = context.WithValue(ctx, Api_keyScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetInventoryById(w, r, inventoryId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetInventoryVersionsById operation middleware
func (siw *ServerInterfaceWrapper) GetInventoryVersionsById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "inventoryId" -------------
	var inventoryId string

	err = runtime.BindStyledParameter("simple", false, "inventoryId", mux.Vars(r)["inventoryId"], &inventoryId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "inventoryId", Err: err})
		return
	}

	ctx = context.WithValue(ctx, Api_keyScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetInventoryVersionsById(w, r, inventoryId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetManufacturers operation middleware
func (siw *ServerInterfaceWrapper) GetManufacturers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, Api_keyScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetManufacturersParams

	// ------------- Optional query parameter "filter[author]" -------------

	err = runtime.BindQueryParameter("form", true, false, "filter[author]", r.URL.Query(), &params.FilterAuthor)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filter[author]", Err: err})
		return
	}

	// ------------- Optional query parameter "filter[mpn]" -------------

	err = runtime.BindQueryParameter("form", true, false, "filter[mpn]", r.URL.Query(), &params.FilterMpn)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filter[mpn]", Err: err})
		return
	}

	// ------------- Optional query parameter "filter[original]" -------------

	err = runtime.BindQueryParameter("form", true, false, "filter[original]", r.URL.Query(), &params.FilterOriginal)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filter[original]", Err: err})
		return
	}

	// ------------- Optional query parameter "search[content]" -------------

	err = runtime.BindQueryParameter("form", true, false, "search[content]", r.URL.Query(), &params.SearchContent)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "search[content]", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetManufacturers(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetMpns operation middleware
func (siw *ServerInterfaceWrapper) GetMpns(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, Api_keyScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetMpnsParams

	// ------------- Optional query parameter "filter[author]" -------------

	err = runtime.BindQueryParameter("form", true, false, "filter[author]", r.URL.Query(), &params.FilterAuthor)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filter[author]", Err: err})
		return
	}

	// ------------- Optional query parameter "filter[manufacturer]" -------------

	err = runtime.BindQueryParameter("form", true, false, "filter[manufacturer]", r.URL.Query(), &params.FilterManufacturer)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filter[manufacturer]", Err: err})
		return
	}

	// ------------- Optional query parameter "filter[original]" -------------

	err = runtime.BindQueryParameter("form", true, false, "filter[original]", r.URL.Query(), &params.FilterOriginal)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filter[original]", Err: err})
		return
	}

	// ------------- Optional query parameter "search[content]" -------------

	err = runtime.BindQueryParameter("form", true, false, "search[content]", r.URL.Query(), &params.SearchContent)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "search[content]", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetMpns(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetThingModelById operation middleware
func (siw *ServerInterfaceWrapper) GetThingModelById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "tmId" -------------
	var tmId string

	err = runtime.BindStyledParameter("simple", false, "tmId", mux.Vars(r)["tmId"], &tmId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "tmId", Err: err})
		return
	}

	ctx = context.WithValue(ctx, Api_keyScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetThingModelById(w, r, tmId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{})
}

type GorillaServerOptions struct {
	BaseURL          string
	BaseRouter       *mux.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r *mux.Router) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r *mux.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options GorillaServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = mux.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.HandleFunc(options.BaseURL+"/authors", wrapper.GetAuthors).Methods("GET")

	r.HandleFunc(options.BaseURL+"/manufacturers", wrapper.GetManufacturers).Methods("GET")

	r.HandleFunc(options.BaseURL+"/mpns", wrapper.GetMpns).Methods("GET")

	r.HandleFunc(options.BaseURL+"/inventory", wrapper.GetInventory).Methods("GET")

	r.HandleFunc(options.BaseURL+"/inventory/{inventoryId:.+}/versions", wrapper.GetInventoryVersionsById).Methods("GET")

	r.HandleFunc(options.BaseURL+"/inventory/{inventoryId:.+}", wrapper.GetInventoryById).Methods("GET")

	r.HandleFunc(options.BaseURL+"/thing-models/{tmId:.+}", wrapper.GetThingModelById).Methods("GET")

	return r
}
