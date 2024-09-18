// Copyright 2023 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jug

import "net/http"

func Default() Engine {
	return defaultGinEngine()
}

func New() Engine {
	return newGinEngine()
}

type Validatable interface {
	Validate() error
}

type Engine interface {
	RouterGroup

	NoMethod(handlers ...HandlerFunc)
	NoRoute(handlers ...HandlerFunc)

	// ExpandMethods expands each non-configured method for each path to return 405 Method not allowed
	ExpandMethods()

	Run(addr ...string) error

	EnableDebugMode()

	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type RouterGroup interface {
	Router
	Group(relativePath string, handlers ...HandlerFunc) RouterGroup
}

type Router interface {
	Use(middleware ...HandlerFunc) Router
	Any(relativePath string, handlers ...HandlerFunc) Router
	GET(relativePath string, handlers ...HandlerFunc) Router
	POST(relativePath string, handlers ...HandlerFunc) Router
	PUT(relativePath string, handlers ...HandlerFunc) Router
	DELETE(relativePath string, handlers ...HandlerFunc) Router
	PATCH(relativePath string, handlers ...HandlerFunc) Router
	OPTIONS(relativePath string, handlers ...HandlerFunc) Router
	HEAD(relativePath string, handlers ...HandlerFunc) Router
}

func MethodNotAllowed(c Context) {
	c.Status(http.StatusMethodNotAllowed)
}
