// Package pow provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package pow

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// VerifySolutionJSONBody defines parameters for VerifySolution.
type VerifySolutionJSONBody struct {
	// Challenge The challenge string provided by the api
	Challenge *string `json:"challenge,omitempty"`

	// Nonce The solution (nonce) provided by the client
	Nonce *string `json:"nonce,omitempty"`
}

// VerifySolutionParams defines parameters for VerifySolution.
type VerifySolutionParams struct {
	// XRemark A unique identifier for the challenge to verify
	XRemark openapi_types.UUID `json:"X-Remark"`
}

// VerifySolutionJSONRequestBody defines body for VerifySolution for application/json ContentType.
type VerifySolutionJSONRequestBody VerifySolutionJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get PoW Hashcash challenge
	// (GET /challenge)
	GetChallenge(ctx echo.Context) error
	// Verify Hashcash PoW solution
	// (POST /verify)
	VerifySolution(ctx echo.Context, params VerifySolutionParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetChallenge converts echo context to params.
func (w *ServerInterfaceWrapper) GetChallenge(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetChallenge(ctx)
	return err
}

// VerifySolution converts echo context to params.
func (w *ServerInterfaceWrapper) VerifySolution(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params VerifySolutionParams

	headers := ctx.Request().Header
	// ------------- Required header parameter "X-Remark" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Remark")]; found {
		var XRemark openapi_types.UUID
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Remark, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Remark", runtime.ParamLocationHeader, valueList[0], &XRemark)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Remark: %s", err))
		}

		params.XRemark = XRemark
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter X-Remark is required, but not found"))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.VerifySolution(ctx, params)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/challenge", wrapper.GetChallenge)
	router.POST(baseURL+"/verify", wrapper.VerifySolution)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RVTW/zRBD+K6s5vZVMEz4OlW9QIajEIaIVQUI9bL3jeNv1jjs7GxQq/3c06yiJqZGg",
	"8J6cjOfreeaZ8Rv42BLUb+AwNewH8RShhm/Nw+3G2MGbgUmwEXTm6WB+tKlrbOrMholaQ63ZEr+YTxva",
	"XhnprBhGyRyT+d0nR715zSSYjG0F2TSdDQHjDs0e2be+sVrtGioQLwGhhi2xK1mnaO3hHnmPDBXskdPU",
	"3JfX6+s1jBXQgNEOHmr4upgqGKx0SeGsTsX03w5FHzQgl5p3Dmr4AeX25FQBYxooJizhX63X+mgoCsYS",
	"a4chHDtePSft4w1S02Fv9dfAmlv8FD2rPSf2oUOzoe0FFy2xkQ5NEzxGMUImUdhrR3IYlJQk7ONO8Trf",
	"tr7JQQ4LE+spR1HyAlrn4878gUyYjI+GMeUgauts6s6ZfRTcIcM4nkz09IyNwKimeYETV2aHUXlEZ1Ju",
	"GkypzSEcoIIOrUMuFPz6xc/YW35ZklaO/jWj8Q6j+NYjnzm4mMeZ3Ja4twI15Ozde17G0mzKfW/5MM21",
	"UHwS6zmrOq6K+AqDA6UFXfxS3t9TyKVjFRXbHqUg++1DcHSqx7IVeA2bqIIKou0VzIkvFeJr9owOauGM",
	"/4aIxykYk3xH7vD5BHzGNZXWI7H3broRClyXckG/kWLzNynTkW7zqThdvUs5rccC6EXpzjkc/9ftLjdt",
	"GcXl1TueQnTH63feldn5+6eI5sVU3yfKpnR/Wcdq8RTrFL75T+iRmfg9+u/VbHpMye7wY5ju4t4G78wl",
	"Npiv9rSaF5+huav6ls/F0qL+RI0NxuEeAw29ntpJpZmD7qPIUK9WQZ06SlLfrG/WMD6OfwYAAP//VxKW",
	"TSIHAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
