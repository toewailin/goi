package templates

// route_template.go - Template for generating routes
const RouteTemplate = `package routes

import "github.com/gorilla/mux"

// SetupRoutes sets up the routes for the application
func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/{{.RouteName}}", func(w http.ResponseWriter, r *http.Request) {
		// Handler logic here
	}).Methods("GET")

	return r
}
`
