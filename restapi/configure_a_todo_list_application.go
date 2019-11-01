// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"go-swagger-example-in-gitpod/models"
	"crypto/tls"
	"net/http"
	"sync"
	"sync/atomic"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"go-swagger-example-in-gitpod/restapi/operations"
	"go-swagger-example-in-gitpod/restapi/operations/todos"
)

//go:generate swagger generate server --target ../../go-swagger-example-in-gitpod --name ATodoListApplication --spec ../swagger.yml

func configureFlags(api *operations.ATodoListApplicationAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

// store items in a map
var items = make(map[int64]*models.Item)
var lastID int64

// Lock for manipulate items
var itemsLock = &sync.Mutex{}

func newItemID() int64 {
	return atomic.AddInt64(&lastID, 1)
}

func allItems(since int64, limit int32)( []*models.Item, errors.Error){
	if limit < 0 {
		return nil, errors.New(400, "Limit must be greater than 0")
	}
	result := make([]*models.Item, 0)
	for id, item := range items {
		if len(result) >= int(limit){
			return result, nil
		}
		if id >= since{
			result = append(result, item)
		}
	}
	return result, nil
}

func addItem(item *models.Item) errors.Error{
	if item == nil{
		return errors.New(400, "Item is not present")
	}
	itemsLock.Lock();
	defer itemsLock.Unlock()

	newID := newItemID()
	item.ID = newID
	items[newID] = item

	return nil
}

func deleteItem(id int64) errors.Error{
	itemsLock.Lock()
	defer itemsLock.Unlock()

	_, exists := items[id]
	if !exists {
		return errors.New(404, "Item %d not found", id)
	}

	delete(items, id)
	return nil
}

func updateItem(id int64, item *models.Item) errors.Error {
	if item == nil {
	 	return errors.New(400, "Item is not present")
	}

	itemsLock.Lock()
	defer itemsLock.Unlock()

	_, exists := items[id]
	if !exists {
		return errors.New(404, "Item %d not found", id)
	}

	item.ID = id
	items[id] = item
	return nil
}

func configureAPI(api *operations.ATodoListApplicationAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.TodosAddOneHandler = todos.AddOneHandlerFunc(func(params todos.AddOneParams) middleware.Responder {
		if err := addItem(params.Body); err != nil {
			return todos.NewAddOneDefault(int(err.Code())).WithPayload(&models.Error{Code: int64(err.Code()), Message: swag.String(err.Error())})
		}
		return todos.NewAddOneCreated().WithPayload(params.Body)
	})

	api.TodosDestroyOneHandler = todos.DestroyOneHandlerFunc(func(params todos.DestroyOneParams) middleware.Responder {
		if err := deleteItem(params.ID); err != nil {
			return todos.NewDestroyOneDefault(int(err.Code())).WithPayload(&models.Error{Code: int64(err.Code()), Message: swag.String(err.Error())})
		}
		return todos.NewDestroyOneNoContent()
	})

	api.TodosFindTodosHandler = todos.FindTodosHandlerFunc(func(params todos.FindTodosParams) middleware.Responder {
		mergedParams := todos.NewFindTodosParams()
		mergedParams.Since = swag.Int64(0)
		if params.Since != nil {
			mergedParams.Since = params.Since
		}
		if params.Limit != nil {
			mergedParams.Limit = params.Limit
		}
		result, err := allItems(*mergedParams.Since, *mergedParams.Limit)
		if err != nil {
			return todos.NewFindTodosDefault(int(err.Code())).WithPayload(&models.Error{Code: int64(err.Code()), Message: swag.String(err.Error())})
		}
		return todos.NewFindTodosOK().WithPayload(result)
	})

	api.TodosUpdateOneHandler = todos.UpdateOneHandlerFunc(func(params todos.UpdateOneParams) middleware.Responder {
		if err := updateItem(params.ID, params.Body); err != nil {
			return todos.NewUpdateOneDefault(int(err.Code())).WithPayload(&models.Error{Code: int64(err.Code()), Message: swag.String(err.Error())})
		}
		return todos.NewUpdateOneOK().WithPayload(params.Body)
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
