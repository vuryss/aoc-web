package internals

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func NewRouter(configFile string) http.Handler {
	router := &router{configFile: configFile}

	// Add allowed methods
	router.supportedMethods = make(map[string]bool)
	router.supportedMethods[http.MethodGet] = true
	router.supportedMethods[http.MethodPost] = true
	router.supportedMethods[http.MethodPut] = true
	router.supportedMethods[http.MethodDelete] = true
	router.supportedMethods[http.MethodOptions] = true

	// Regex for path validation (very restrictive, I like it that way - deal with it)
	router.pathRegex = regexp.MustCompile(`^[0-9a-z\-\\]$`)

	// Init routes
	router.routes = make(map[string][]route)

	// Parse config file
	router.parseConfiguration()

	log.Printf("Router parses: %v", router.parses)
	log.Print(router.routes)

	return router
}

/**
 * Route definition
 */
type route struct {
	method string
	path   string
}

/**
 * Router definition
 */
type router struct {
	configFile 		 string
	supportedMethods map[string]bool
	routes			 map[string][]route
	pathRegex		 *regexp.Regexp
	parses 			 int
}

func (router *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done"))
}

func (router *router) parseConfiguration() {
	contents, err := ioutil.ReadFile(router.configFile)

	if err != nil {
		panic("ROUTER: Could not read router configuration file: " + err.Error())
	}

	lines := strings.Split(string(contents), "\n")

	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])

		if len(lines[i]) == 0 {
			continue
		}

		routeParts := strings.Split(lines[i], "|")

		if len(routeParts) < 2 {
			log.Printf("ROUTER: Invalid route definition found: %v", lines[i])
			continue
		}

		router.generateRoute(routeParts[0], routeParts[1], routeParts[2:]...)
	}

	router.parses++
}

func (router *router) generateRoute(routeDef string, destination string, parameters ...string) {
	routeParts := strings.Split(strings.TrimSpace(routeDef), ":")

	if len(routeParts) != 2 {
		log.Printf("ROUTER: Unable to parse route definition: %v", routeDef)
		return
	}

	if _, exists := router.supportedMethods[routeParts[0]]; !exists {
		log.Printf("ROUTER: Unsupported HTTP Method: %v", routeParts[0])
		return
	}

	if router.pathRegex.MatchString(routeParts[1]) {
		log.Printf("ROUTER: Unsupported path in route definition: %v", routeParts[1])
		return
	}

	router.routes[routeParts[0]] = append(
		router.routes[routeParts[0]],
		route{
			method: routeParts[0],
			path: routeParts[1],
		},
	)
}