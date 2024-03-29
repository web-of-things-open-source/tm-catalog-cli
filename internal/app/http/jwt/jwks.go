package jwt

import (
	"fmt"
	"net/url"

	"github.com/MicahParks/keyfunc/v3"
)

type JWTValidationOpts struct {
	JWTServiceID  string
	JWKSURLString string
}

func validateOptions(opts JWTValidationOpts) {
	// JWKS URL must be initialized
	var err error
	_, err = url.ParseRequestURI(opts.JWKSURLString)
	if err != nil {
		msg := fmt.Sprintf(`jwt validation activated, but jwks URL could not be initialized:
%s
Check documentation to "serve" to configure correctly
`, err)
		panic(msg)
	}
}

func startJWKSFetch(opts JWTValidationOpts) keyfunc.Keyfunc {
	validateOptions(opts)
	// start a new go routine fetching the jwks periodically
	k, err := keyfunc.NewDefault([]string{opts.JWKSURLString})
	if err != nil {
		msg := fmt.Sprintf("Failed to create a keyfunc.Keyfunc from the server's URL.\nError: %s", err)
		panic(msg)
	}
	return k
}
