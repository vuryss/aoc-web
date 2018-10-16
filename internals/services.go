package internals

import "./service"

var Services = map[string]interface{} {
	"Index": (*service.IndexService)(nil),
	"Solver": (*service.SolverService)(nil),
}