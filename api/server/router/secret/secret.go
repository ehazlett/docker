package secret

import "github.com/docker/docker/api/server/router"

// secretRouter is a router to talk with the secret controller
type secretRouter struct {
	backend Backend
	routes  []router.Route
}

// NewRouter initializes a new secret router
func NewRouter(b Backend) router.Router {
	r := &secretRouter{
		backend: b,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routers to the secret controller
func (sr *secretRouter) Routes() []router.Route {
	return sr.routes
}

func (sr *secretRouter) initRoutes() {
	sr.routes = []router.Route{
		// GET
		router.NewGetRoute("/secrets/{name:.*}", sr.inspectSecret),
		router.NewGetRoute("/secrets", sr.listSecrets),
		// PUT
		router.NewPutRoute("/secrets/{name:.*}", sr.updateSecret),
		// POST
		router.NewPostRoute("/secrets", sr.createSecret),
		// DELETE
		router.NewDeleteRoute("/secrets/{name:.*}", sr.removeSecret),
	}
}
