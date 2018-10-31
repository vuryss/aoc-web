package internals

import (
	"./core"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func NewRouter(config *core.Config) http.Handler {
	routersFile, found := config.GetString("routes_file")

	if !found {
		panic("ROUTER: Cannot find routes file in config")
	}

	router := &router{config: config, configFile: routersFile}

	// Add allowed methods
	router.supportedMethods = make(map[string]bool)
	router.supportedMethods[http.MethodGet] = true
	router.supportedMethods[http.MethodPost] = true
	router.supportedMethods[http.MethodPut] = true
	router.supportedMethods[http.MethodDelete] = true
	router.supportedMethods[http.MethodOptions] = true

	// Regex for path validation (very restrictive, I like it that way - deal with it)
	router.pathRegex = regexp.MustCompile(`^[{}0-9a-z\-/]+$`)

	// Init routes
	router.routes = make(map[string][]route)

	// Parse config file
	router.parseConfiguration()

	log.Print("Literal routes:")
	log.Print(router.literalRoutes)
	log.Print("Parameter routes:")
	log.Print(router.parameterRoutes)

	return router
}

/**
 * Route definition
 */
type route struct {
	method 		string
	path   		string
	parts  		[]routePart
	parameters 	map[string]string
	controller	string
	action 		string
}

type routePart struct {
	partType 	int
	value 		string
}

/**
 * Destination Parameters
 */
type RouteParameters struct {
	response 	http.ResponseWriter
	request 	*http.Request
	parameters  map[string]string
}

/**
 * Router definition
 */
type router struct {
	config           *core.Config
	configFile 		 string
	supportedMethods map[string]bool
	routes			 map[string][]route
	pathRegex		 *regexp.Regexp
	parses 			 int

	/**
	 * Literal paths match the input path fully, meaning no user-input parameters are parsed.
	 * They are in format of [GET][Path] -> Route
	 */
	literalRoutes		map[string]map[string]route

	/**
	 * Parameter paths match urls with user-input parameters, which should go into given parameter placeholder or name.
	 */
	parameterRoutes		map[string][]route
}

func (router *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	route, isMatched := router.match(r.Method, r.URL.Path)
	log.Printf("Route matched in %v", time.Since(start))

	if !isMatched {
		w.WriteHeader(404)
		return
	}

	reflectInstance := reflect.New(reflect.TypeOf(Services[route.controller]).Elem())

	abstractService := reflect.Indirect(reflectInstance).FieldByName("Service")
	abstractService.Set(reflect.ValueOf(&core.Service{
		Request		: r,
		Response	: w,
		Parameters	: route.parameters,
		Config		: router.config,
		View        : core.NewView(router.config),
	}))

	methodRef := reflectInstance.MethodByName(route.action)
	responseValue := methodRef.Call(nil)[0]
	methodRef = responseValue.MethodByName("GetBody")
	responseBodyValue := methodRef.Call(nil)[0]

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseBodyValue.String()))
}

func (router *router) match(method string, path string) (route, bool) {
	log.Printf("Method: %v | Path: %v", method, path)

	if _, exists := router.literalRoutes[method]; exists {
		if _, exists := router.literalRoutes[method][path]; exists {
			return router.literalRoutes[method][path], true
		}
	}

	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	pathLength := len(pathParts)

	if _, exists := router.parameterRoutes[method]; exists {
		RouteMatch:
		for i := range router.parameterRoutes[method] {
			route := router.parameterRoutes[method][i]
			route.parameters = make(map[string]string)

			if len(route.parts) != pathLength {
				continue
			}

			for j := range route.parts {
				if route.parts[j].partType == 0 && pathParts[j] != route.parts[j].value {
					continue RouteMatch
				}

				if route.parts[j].partType == 1 {
					route.parameters[route.parts[j].value] = pathParts[j]
				}
			}

			return route, true
		}
	}

	return route{}, false
}

/**
 * Parses router configuration from the routes config file and generates ready to use routes in categories.
 */
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

/**
 * Generates a single route based on definition in configuration file
 */
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

	if !router.pathRegex.MatchString(routeParts[1]) {
		log.Printf("ROUTER: Unsupported path in route definition: %v", routeParts[1])
		return
	}

	trimmedPath := "/" + strings.Trim(routeParts[1], "/")

	// Parameter route
	if strings.ContainsRune(routeParts[1], '{') {
		router.createParameterPath(routeParts[0], trimmedPath, destination, parameters...)
		return
	}

	// Literal route
	router.createLiteralRoute(routeParts[0], trimmedPath, destination, parameters...)
}

/**
 * Creates literal route, which must fully match given URL to be activated
 */
func (router *router) createLiteralRoute(method string, path string, destination string, parameters ...string) {
	if router.literalRoutes == nil {
		router.literalRoutes = make(map[string]map[string]route)
	}

	if _, exists := router.literalRoutes[method]; !exists {
		router.literalRoutes[method] = make(map[string]route)
	}

	if _, exists := router.literalRoutes[method][path]; exists {
		log.Printf("ROUTER: Duplicated literal route: " + path)
	}

	destinationParts := strings.Split(destination, "::")

	if len(destinationParts) != 2 {
		log.Printf("ROUTER: Invalid route target: %v", destination)
	}

	router.literalRoutes[method][path] = route{
		method: method,
		path: path,
		controller: strings.TrimSpace(destinationParts[0]),
		action: strings.TrimSpace(destinationParts[1]),
	}
}

/**
 * Creates a parametrized route, which matches a dynamic URL with a parameter in it
 */
func (router *router) createParameterPath(method string, path string, destination string, parameters ...string) {
	pathParts := strings.Split(path, "/")
	routeParts := make([]routePart, 0)

	for i := range pathParts {
		if i == 0 {
			continue
		}

		// Is this a parameter?
		if strings.HasPrefix(pathParts[i], "{") && strings.HasSuffix(pathParts[i], "}") {
			if len(pathParts[i]) < 3 {
				log.Printf("ROUTER: Empty parameter name, skipping route: %v:%v", method, path)
			}

			paramName := pathParts[i][1:len(pathParts[i]) - 1]

			routeParts = append(routeParts, routePart{
				partType: 1,
				value: paramName,
			})

			continue
		}

		// Normal literal parameter
		routeParts = append(routeParts, routePart{
			partType: 0,
			value: pathParts[i],
		})
	}

	if router.parameterRoutes == nil {
		router.parameterRoutes = make(map[string][]route)
	}

	if _, exists := router.parameterRoutes[method]; !exists {
		router.parameterRoutes[method] = make([]route, 0)
	}

	destinationParts := strings.Split(destination, "::")

	if len(destinationParts) != 2 {
		log.Printf("ROUTER: Invalid route target: %v", destination)
	}

	router.parameterRoutes[method] = append(router.parameterRoutes[method], route{
		method: method,
		path: path,
		parts: routeParts,
		controller: strings.TrimSpace(destinationParts[0]),
		action: strings.TrimSpace(destinationParts[1]),
	})
}