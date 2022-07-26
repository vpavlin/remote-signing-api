// Package signer provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package signer

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// Address defines model for Address.
type Address = string

// NewSigner200 defines model for NewSigner200.
type NewSigner200 struct {
	PublicKey *Address `json:"publicKey,omitempty"`
}

// SignBytes defines model for SignBytes.
type SignBytes struct {
	Bytes *[]byte `json:"bytes,omitempty"`
}

// SignBytes200 defines model for SignBytes200.
type SignBytes200 struct {
	SignedData *[]byte `json:"signedData,omitempty"`
}

// SignerKey defines model for SignerKey.
type SignerKey struct {
	Key *string `json:"key,omitempty"`
}

// NewSignerJSONBody defines parameters for NewSigner.
type NewSignerJSONBody = SignerKey

// SignBytesJSONBody defines parameters for SignBytes.
type SignBytesJSONBody = SignBytes

// SignBytesParams defines parameters for SignBytes.
type SignBytesParams struct {
	Authorization string `json:"Authorization"`
}

// ReplaceKeyJSONBody defines parameters for ReplaceKey.
type ReplaceKeyJSONBody = SignerKey

// ReplaceKeyParams defines parameters for ReplaceKey.
type ReplaceKeyParams struct {
	Authorization string `json:"Authorization"`
}

// NewSignerJSONRequestBody defines body for NewSigner for application/json ContentType.
type NewSignerJSONRequestBody = NewSignerJSONBody

// SignBytesJSONRequestBody defines body for SignBytes for application/json ContentType.
type SignBytesJSONRequestBody = SignBytesJSONBody

// ReplaceKeyJSONRequestBody defines body for ReplaceKey for application/json ContentType.
type ReplaceKeyJSONRequestBody = ReplaceKeyJSONBody

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// Health request
	Health(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// NewSigner request with any body
	NewSignerWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	NewSigner(ctx context.Context, body NewSignerJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// SignBytes request with any body
	SignBytesWithBody(ctx context.Context, address Address, params *SignBytesParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	SignBytes(ctx context.Context, address Address, params *SignBytesParams, body SignBytesJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ReplaceKey request with any body
	ReplaceKeyWithBody(ctx context.Context, address Address, params *ReplaceKeyParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	ReplaceKey(ctx context.Context, address Address, params *ReplaceKeyParams, body ReplaceKeyJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) Health(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewHealthRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) NewSignerWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewNewSignerRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) NewSigner(ctx context.Context, body NewSignerJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewNewSignerRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SignBytesWithBody(ctx context.Context, address Address, params *SignBytesParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSignBytesRequestWithBody(c.Server, address, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SignBytes(ctx context.Context, address Address, params *SignBytesParams, body SignBytesJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSignBytesRequest(c.Server, address, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ReplaceKeyWithBody(ctx context.Context, address Address, params *ReplaceKeyParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewReplaceKeyRequestWithBody(c.Server, address, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ReplaceKey(ctx context.Context, address Address, params *ReplaceKeyParams, body ReplaceKeyJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewReplaceKeyRequest(c.Server, address, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewHealthRequest generates requests for Health
func NewHealthRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/signer/health")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewNewSignerRequest calls the generic NewSigner builder with application/json body
func NewNewSignerRequest(server string, body NewSignerJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewNewSignerRequestWithBody(server, "application/json", bodyReader)
}

// NewNewSignerRequestWithBody generates requests for NewSigner with any type of body
func NewNewSignerRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/signer/new")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewSignBytesRequest calls the generic SignBytes builder with application/json body
func NewSignBytesRequest(server string, address Address, params *SignBytesParams, body SignBytesJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewSignBytesRequestWithBody(server, address, params, "application/json", bodyReader)
}

// NewSignBytesRequestWithBody generates requests for SignBytes with any type of body
func NewSignBytesRequestWithBody(server string, address Address, params *SignBytesParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "address", runtime.ParamLocationPath, address)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/signer/%s/bytes", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	var headerParam0 string

	headerParam0, err = runtime.StyleParamWithLocation("simple", false, "Authorization", runtime.ParamLocationHeader, params.Authorization)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", headerParam0)

	return req, nil
}

// NewReplaceKeyRequest calls the generic ReplaceKey builder with application/json body
func NewReplaceKeyRequest(server string, address Address, params *ReplaceKeyParams, body ReplaceKeyJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewReplaceKeyRequestWithBody(server, address, params, "application/json", bodyReader)
}

// NewReplaceKeyRequestWithBody generates requests for ReplaceKey with any type of body
func NewReplaceKeyRequestWithBody(server string, address Address, params *ReplaceKeyParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "address", runtime.ParamLocationPath, address)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/signer/%s/key", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	var headerParam0 string

	headerParam0, err = runtime.StyleParamWithLocation("simple", false, "Authorization", runtime.ParamLocationHeader, params.Authorization)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", headerParam0)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// Health request
	HealthWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*HealthResponse, error)

	// NewSigner request with any body
	NewSignerWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*NewSignerResponse, error)

	NewSignerWithResponse(ctx context.Context, body NewSignerJSONRequestBody, reqEditors ...RequestEditorFn) (*NewSignerResponse, error)

	// SignBytes request with any body
	SignBytesWithBodyWithResponse(ctx context.Context, address Address, params *SignBytesParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SignBytesResponse, error)

	SignBytesWithResponse(ctx context.Context, address Address, params *SignBytesParams, body SignBytesJSONRequestBody, reqEditors ...RequestEditorFn) (*SignBytesResponse, error)

	// ReplaceKey request with any body
	ReplaceKeyWithBodyWithResponse(ctx context.Context, address Address, params *ReplaceKeyParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ReplaceKeyResponse, error)

	ReplaceKeyWithResponse(ctx context.Context, address Address, params *ReplaceKeyParams, body ReplaceKeyJSONRequestBody, reqEditors ...RequestEditorFn) (*ReplaceKeyResponse, error)
}

type HealthResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r HealthResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r HealthResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type NewSignerResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *NewSigner200
}

// Status returns HTTPResponse.Status
func (r NewSignerResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r NewSignerResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type SignBytesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *SignBytes200
}

// Status returns HTTPResponse.Status
func (r SignBytesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SignBytesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ReplaceKeyResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r ReplaceKeyResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ReplaceKeyResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// HealthWithResponse request returning *HealthResponse
func (c *ClientWithResponses) HealthWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*HealthResponse, error) {
	rsp, err := c.Health(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseHealthResponse(rsp)
}

// NewSignerWithBodyWithResponse request with arbitrary body returning *NewSignerResponse
func (c *ClientWithResponses) NewSignerWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*NewSignerResponse, error) {
	rsp, err := c.NewSignerWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseNewSignerResponse(rsp)
}

func (c *ClientWithResponses) NewSignerWithResponse(ctx context.Context, body NewSignerJSONRequestBody, reqEditors ...RequestEditorFn) (*NewSignerResponse, error) {
	rsp, err := c.NewSigner(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseNewSignerResponse(rsp)
}

// SignBytesWithBodyWithResponse request with arbitrary body returning *SignBytesResponse
func (c *ClientWithResponses) SignBytesWithBodyWithResponse(ctx context.Context, address Address, params *SignBytesParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SignBytesResponse, error) {
	rsp, err := c.SignBytesWithBody(ctx, address, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSignBytesResponse(rsp)
}

func (c *ClientWithResponses) SignBytesWithResponse(ctx context.Context, address Address, params *SignBytesParams, body SignBytesJSONRequestBody, reqEditors ...RequestEditorFn) (*SignBytesResponse, error) {
	rsp, err := c.SignBytes(ctx, address, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSignBytesResponse(rsp)
}

// ReplaceKeyWithBodyWithResponse request with arbitrary body returning *ReplaceKeyResponse
func (c *ClientWithResponses) ReplaceKeyWithBodyWithResponse(ctx context.Context, address Address, params *ReplaceKeyParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ReplaceKeyResponse, error) {
	rsp, err := c.ReplaceKeyWithBody(ctx, address, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseReplaceKeyResponse(rsp)
}

func (c *ClientWithResponses) ReplaceKeyWithResponse(ctx context.Context, address Address, params *ReplaceKeyParams, body ReplaceKeyJSONRequestBody, reqEditors ...RequestEditorFn) (*ReplaceKeyResponse, error) {
	rsp, err := c.ReplaceKey(ctx, address, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseReplaceKeyResponse(rsp)
}

// ParseHealthResponse parses an HTTP response from a HealthWithResponse call
func ParseHealthResponse(rsp *http.Response) (*HealthResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &HealthResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseNewSignerResponse parses an HTTP response from a NewSignerWithResponse call
func ParseNewSignerResponse(rsp *http.Response) (*NewSignerResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &NewSignerResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest NewSigner200
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseSignBytesResponse parses an HTTP response from a SignBytesWithResponse call
func ParseSignBytesResponse(rsp *http.Response) (*SignBytesResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &SignBytesResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest SignBytes200
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseReplaceKeyResponse parses an HTTP response from a ReplaceKeyWithResponse call
func ParseReplaceKeyResponse(rsp *http.Response) (*ReplaceKeyResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ReplaceKeyResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Health check endpoint
	// (GET /signer/health)
	Health(ctx echo.Context) error
	// Creates new signer
	// (POST /signer/new)
	NewSigner(ctx echo.Context) error
	// Signes bytes
	// (POST /signer/{address}/bytes)
	SignBytes(ctx echo.Context, address Address, params SignBytesParams) error
	// Replace the API key
	// (PUT /signer/{address}/key)
	ReplaceKey(ctx echo.Context, address Address, params ReplaceKeyParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// Health converts echo context to params.
func (w *ServerInterfaceWrapper) Health(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.Health(ctx)
	return err
}

// NewSigner converts echo context to params.
func (w *ServerInterfaceWrapper) NewSigner(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.NewSigner(ctx)
	return err
}

// SignBytes converts echo context to params.
func (w *ServerInterfaceWrapper) SignBytes(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "address" -------------
	var address Address

	err = runtime.BindStyledParameterWithLocation("simple", false, "address", runtime.ParamLocationPath, ctx.Param("address"), &address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params SignBytesParams

	headers := ctx.Request().Header
	// ------------- Required header parameter "Authorization" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("Authorization")]; found {
		var Authorization string
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for Authorization, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "Authorization", runtime.ParamLocationHeader, valueList[0], &Authorization)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter Authorization: %s", err))
		}

		params.Authorization = Authorization
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter Authorization is required, but not found"))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.SignBytes(ctx, address, params)
	return err
}

// ReplaceKey converts echo context to params.
func (w *ServerInterfaceWrapper) ReplaceKey(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "address" -------------
	var address Address

	err = runtime.BindStyledParameterWithLocation("simple", false, "address", runtime.ParamLocationPath, ctx.Param("address"), &address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params ReplaceKeyParams

	headers := ctx.Request().Header
	// ------------- Required header parameter "Authorization" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("Authorization")]; found {
		var Authorization string
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for Authorization, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "Authorization", runtime.ParamLocationHeader, valueList[0], &Authorization)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter Authorization: %s", err))
		}

		params.Authorization = Authorization
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter Authorization is required, but not found"))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ReplaceKey(ctx, address, params)
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

	router.GET(baseURL+"/signer/health", wrapper.Health)
	router.POST(baseURL+"/signer/new", wrapper.NewSigner)
	router.POST(baseURL+"/signer/:address/bytes", wrapper.SignBytes)
	router.PUT(baseURL+"/signer/:address/key", wrapper.ReplaceKey)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xV0W7jNhD8FWJ7j4qtyL4Ep6faF197SNEe2qJ9OKTAWlpLvEgkS9J2XEMf0w/oV9yP",
	"FUtJ9tkxEgNN+9Q3hzFnZmeH4y1kujZakfIO0i24rKQaw8dJnlty4SM9YG0qghTih2kyHr1OXo/eXl++",
	"yybJ9Xw2vkziq9k0uX4zmo3HdHmzuJpf3UAEBr0nqyCF3+KHj3ixmFy8iy/e3G3HcfMKIvAbw5jOW6kK",
	"aCL4ntY/yUKRTeKYeY3VhqyXFFSY5byS2S1t+I9XlhaQwlfDvf5hJ37YK2+aHYeef6LMMwcTTDe+gzwg",
	"mPfHC21r9JCGk8dCn4Q9Kd3xVPkNevwH8GS72Q+x79vDWqrvSBW+hHSUnIHJR1ItNN/NtPKY+bDqGmUF",
	"Kawwq3A1MLiqpPq64NNBpmuIQGHNOL98/pO/IT7gqvr8l2KRObnMSuOl5p3PfEmWlrX4FauKvGhHEDUq",
	"LMiJdTh1AlUufEnSCmPlCj2Je9q0x+yaE46UFzlbF4GXPsSwxYIIVmRdSxcP4sElq9CGFBoJKYwG8SBp",
	"c1gGq4ZhD3ZYElZs1BYKClMfKp+Imnypc+G1yErK7oVcsEYx+fBeSCcsYVYSzisKMlu0DQRqiwzyPocU",
	"vm1ZIrDkjFauXVeXj0PGH27DityyrtFudnc7dlK50VJ5NgALB+nHNlAW7vhWP5WidUiHdidm+oYUayMn",
	"FK0784P6Go0T0vOsKAq5IhXG5FQdD7R7nmGm35fk/FTnmz5BpAIvGlPJLNwafnJM3tfKc+92H/KQTqaQ",
	"lnJIvV1Sc9rHFyE+6J3A/cx63lraeel6S57azRbbTmqGu5I5vaegw4nwrTb6Uok5u3y8jX2RccIt1uTJ",
	"soAtSEYqCfOgq3uwk6UvtZV/BAQ4tjf6wqpH3dEh8jva43UTPYl0Xkvf/Xtxag36j+N08FtwTpy+XPrZ",
	"Qep63yxPpOjnUrq+xLCq9NqF/tJrbmCvhdWem7bvtFOP/UcyFWZ0G/71f75eqK6eCUJn+tFiHuehaf4O",
	"AAD//3kkHZa5CQAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
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
	var res = make(map[string]func() ([]byte, error))
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
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
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
