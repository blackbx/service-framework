package middleware

import (
	"net/http"

	"github.com/BlackBX/service-framework/dependency"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

// CORSService allows the the cors middleware to be registered
// with an application
var CORSService = dependency.Service{
	Name: "cors",
	ConfigFunc: func(set dependency.FlagSet) {
		set.StringSlice(
			"cors-allowed-headers",
			[]string{
				":authority",
				":method",
				":path",
				":scheme",
				"Accept",
				"Accept-Encoding",
				"Accept-Language",
				"Authorization",
				"Origin",
				"Referer",
				"Sec-Fetch-Mode",
				"Sec-Fetch-Site",
				"User-Agent",
				"X-Forwarded-For",
				"X-Real-IP",
				"X-Forwarded-Proto",
				"X-Requested-With",
			},
			"The headers allowed to be passed from a CORS request",
		)
		set.StringSlice(
			"cors-allowed-methods",
			[]string{
				http.MethodGet,
				http.MethodHead,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
				http.MethodConnect,
				http.MethodOptions,
				http.MethodTrace,
			},
			"The headers to allow cross origin requests with",
		)
		set.StringSlice(
			"cors-allowed-origins",
			[]string{
				"localhost",
			},
			"The origins to allow requests from",
		)
		set.Bool(
			"cors-allow-credentials",
			true,
			"Whether to allow credentials over cross origin",
		)
	},
	Constructor: fx.Annotated{
		Group:  "middleware",
		Target: NewCORS,
	},
}

// NewCORS creates a new cors middleware configured from the app
func NewCORS(config dependency.ConfigGetter) mux.MiddlewareFunc {
	baseOptions := []handlers.CORSOption{
		handlers.AllowedHeaders(config.GetStringSlice("cors-allowed-headers")),
		handlers.AllowedMethods(config.GetStringSlice("cors-allowed-methods")),
		handlers.AllowedOrigins(config.GetStringSlice("cors-allowed-origins")),
	}
	options := make([]handlers.CORSOption, 0, 4)
	if config.GetBool("cors-allow-credentials") {
		options = append(options, handlers.AllowCredentials())
	}
	options = append(options, baseOptions...)
	middleware := handlers.CORS(options...)
	return middleware
}
