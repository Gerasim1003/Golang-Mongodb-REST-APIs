package store

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

//Route ...
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes ...
type Routes []Route

var controller *Controller = &Controller{Repository: NewRepository()}

var routes = Routes{
	Route{
		"Authentication",
		"POST",
		"/get-token",
		controller.GetToken,
	},
	Route{
		"Index",
		"GET",
		"/getHeroes",
		AuthenticationMiddleware(controller.Index),
	},
	Route{
		"AddHero",
		"POST",
		"/addHero",
		AuthenticationMiddleware(controller.AddHero),
	},
}

//NewRouter ...
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.HandlerFunc
		log.Println(route.Name)
		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
