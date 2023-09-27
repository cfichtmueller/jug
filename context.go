// Copyright 2023 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jug

import (
	"io"
	"time"
)

type Context interface {
	// Get gets a value from the context.
	Get(key string) (any, bool)
	// MustGet tries to get a value from the context. If the value cannot be found the request is aborted.
	MustGet(key string) any
	// Set sets a context value.
	Set(key string, value any)

	// Query gets a raw query value
	Query(key string) string
	// QueryArray gets an array query value
	QueryArray(key string) []string
	// IntQuery gets a query value as int
	IntQuery(key string) (int, error)
	// BoolQuery gets a query value as bool
	BoolQuery(key string) (bool, error)
	// Iso8601DateQuery gets a query value as ISO 8601 Date
	Iso8601DateQuery(key string) (*time.Time, error)
	// Iso8601DateTimeQuery gets a query values as ISO 8601 DateTime
	Iso8601DateTimeQuery(key string) (*time.Time, error)
	// StringQuery gets a query value as string. This method performs unescaping.
	StringQuery(key string) (string, error)
	// DefaultQuery gets a query value. If the value cannot be found a default value is returned.
	DefaultQuery(key string, defaultValue string) string
	// DefaultIntQuery gets a query value as int. If the value cannot be found a default value is returned.
	DefaultIntQuery(key string, defaultValue int) (int, error)
	// DefaultBoolQuery gets a query value as bool. If the value cannot be found a default value is returned.
	DefaultBoolQuery(key string, defaultValue bool) (bool, error)
	// DefaultStringQuery gets a query value as string. If the value cannot be found a default value is returned.
	DefaultStringQuery(key string, defaultValue string) (string, error)
	// GetHeader gets a request header
	GetHeader(key string) string

	// Param gets a request param (aka path parameter)
	Param(key string) string

	//TODO: ParamAsInt

	// GetRawData gets the raw request body
	GetRawData() ([]byte, error)
	// MayBindJSON tries to bind the request body from JSON to the given object.
	MayBindJSON(obj any) bool
	// MayBindJSONV tries to bind the request body from JSON to the given object. It will then invoke the provided validator function.
	MayBindJSONV(obj any, validator func() error) bool
	// MustBindJSON tries to bind the request body from JSON to the given object. If that fails the request is aborted with 400.
	MustBindJSON(obj any) bool
	// MustBindJSONV tries to bind the request body from JSON to the given object. If that fails the request is aborted with 400.
	// If it succeeds the provided validator function is invoked.
	MustBindJSONV(obj any, validator func() error) bool

	// Status sets the response status code.
	Status(code int) Context
	// String sets the response status code and writes a string response.
	String(code int, format string, values ...any) Context

	// SetHeader sets a response header.
	SetHeader(key string, value string)
	// SetContentType sets the response content type.
	SetContentType(value string)

	// Cookie returns the named cookie provided in the request or false if not found.
	// If multiple cookies match the given name, only one cookie will be returned.
	Cookie(name string) (string, bool)
	// SetCookie sets a cookie.
	SetCookie(name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool)

	// Stream writes a stream response.
	Stream(step func(w io.Writer) bool) bool
	// SSEvent writes a server sent event.
	SSEvent(name string, message any)

	// Data sets the response status code and writes the given data as is.
	Data(code int, contentType string, data []byte)

	// RespondOk sets status 200, marshals obj to JSON
	RespondOk(obj any)
	// RespondNoContent sets status 204, no response body
	RespondNoContent()
	// RespondCreated sets status 201, marshals obj to JSON
	RespondCreated(obj any)
	// RespondForbidden sets status 403, marshals obj to JSON
	RespondForbidden(obj any)
	// RespondForbiddenE sets status 403, writes error as error response (JSON)
	RespondForbiddenE(err error)
	// RespondUnauthorized sets status 401, marshals obj to JSON
	RespondUnauthorized(obj any)
	// RespondUnauthorizedE sets status 401, writes error as error response (JSON)
	RespondUnauthorizedE(err error)
	// RespondBadRequest sets status 400, marshals obj to JSON
	RespondBadRequest(obj any)
	// RespondBadRequestE sets status 400, writes error as error response (JSON)
	RespondBadRequestE(err error)
	// RespondNotFound sets status 404, marshals obj to JSON
	RespondNotFound(obj any)
	// RespondNotFoundE sets status 404, writes error as error response (JSON)
	RespondNotFoundE(err error)
	// RespondConflict sets status 409, marshals obj to JSON
	RespondConflict(obj any)
	// RespondConflictE sets status 409, writes error as error response (JSON)
	RespondConflictE(err error)
	// RespondInternalServerError sets status 500, marshals obj to JSON
	RespondInternalServerError(obj any)
	// RespondInternalServerErrorE sets status 500, writes error as error response (JSON)
	RespondInternalServerErrorE(err error)

	// RespondMissingRequestBody sets status 400, writes error response (JSON)
	RespondMissingRequestBody()

	// Abort prevents pending handlers being called. This will not stop the current handler.
	Abort()
	// AbortWithError prevents pending handlers being called. This will not stop the current handler.
	// The response status code is set to the given value.
	// Writes error as error response (JSON).
	AbortWithError(code int, err error)

	// Next should only be used in middlewares. It executes pending handlers in the chain inside the current handler.
	Next()

	// HandleError inspects the given error and writes an appropriate response.
	HandleError(err error)

	// Deadline returns that there is no deadline (ok==false) when c.Request has no Context.
	Deadline() (deadline time.Time, ok bool)
	// Done returns nil (chan which will wait forever) when c.Request has no Context.
	Done() <-chan struct{}
	// Err returns nil when c.Request has no Context.
	Err() error
	// Value returns the value associated with this context for key, or nil if no value is associated with key.
	// Successive calls to Value with the same key returns the same result.
	Value(key any) any
}
