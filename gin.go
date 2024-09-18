// Copyright 2023 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jug

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ginEngine struct {
	engine       *gin.Engine
	pathRegistry *PathRegistry
	groups       []*ginRouterGroup
}

func defaultGinEngine() Engine {
	return &ginEngine{
		engine:       gin.Default(),
		pathRegistry: NewPathRegistry(),
		groups:       make([]*ginRouterGroup, 0),
	}
}

func newGinEngine() Engine {
	gin.SetMode(gin.ReleaseMode)
	return &ginEngine{
		engine:       gin.New(),
		pathRegistry: NewPathRegistry(),
		groups:       make([]*ginRouterGroup, 0),
	}
}

func (r *ginEngine) EnableDebugMode() {
	gin.SetMode(gin.DebugMode)
}

func (r *ginEngine) Use(middleware ...HandlerFunc) Router {
	return &ginRoutesRouter{routes: r.engine.Use(MapMany(middleware, wrapHandler)...)}
}

func (r *ginEngine) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	g := newGinRouterGroup(r.engine.Group(relativePath, MapMany(handlers, wrapHandler)...))
	r.groups = append(r.groups, g)
	return g
}

func (r *ginEngine) Any(relativePath string, handlers ...HandlerFunc) Router {
	r.pathRegistry.Add(relativePath, "GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD")
	return &ginRoutesRouter{routes: r.engine.Any(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginEngine) GET(relativePath string, handlers ...HandlerFunc) Router {
	r.pathRegistry.Add(relativePath, "GET")
	return &ginRoutesRouter{routes: r.engine.GET(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginEngine) POST(relativePath string, handlers ...HandlerFunc) Router {
	r.pathRegistry.Add(relativePath, "POST")
	return &ginRoutesRouter{routes: r.engine.POST(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginEngine) PUT(relativePath string, handlers ...HandlerFunc) Router {
	r.pathRegistry.Add(relativePath, "PUT")
	return &ginRoutesRouter{routes: r.engine.PUT(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginEngine) DELETE(relativePath string, handlers ...HandlerFunc) Router {
	r.pathRegistry.Add(relativePath, "DELETE")
	return &ginRoutesRouter{routes: r.engine.DELETE(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginEngine) PATCH(relativePath string, handlers ...HandlerFunc) Router {
	r.pathRegistry.Add(relativePath, "PATCH")
	return &ginRoutesRouter{routes: r.engine.PATCH(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginEngine) OPTIONS(relativePath string, handlers ...HandlerFunc) Router {
	r.pathRegistry.Add(relativePath, "OPTIONS")
	return &ginRoutesRouter{routes: r.engine.OPTIONS(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginEngine) HEAD(relativePath string, handlers ...HandlerFunc) Router {
	r.pathRegistry.Add(relativePath, "HEAD")
	return &ginRoutesRouter{routes: r.engine.HEAD(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginEngine) NoMethod(handlers ...HandlerFunc) {
	r.engine.HandleMethodNotAllowed = true
	r.engine.NoMethod(MapMany(handlers, wrapHandler)...)
}

func (r *ginEngine) NoRoute(handlers ...HandlerFunc) {
	r.engine.NoRoute(MapMany(handlers, wrapHandler)...)
}

func (r *ginEngine) ExpandMethods() {
	expandMethods(r, r.pathRegistry)
	for _, g := range r.groups {
		g.expandMethods()
	}
}

func (r *ginEngine) Run(addr ...string) error {
	return r.engine.Run(addr...)
}

func (r *ginEngine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.engine.ServeHTTP(w, req)
}

type ginRoutesRouter struct {
	routes gin.IRoutes
}

func (r *ginRoutesRouter) Use(middleware ...HandlerFunc) Router {
	return &ginRoutesRouter{routes: r.routes.Use(MapMany(middleware, wrapHandler)...)}
}

func (r *ginRoutesRouter) Any(relativePath string, handlers ...HandlerFunc) Router {
	return &ginRoutesRouter{routes: r.routes.Any(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginRoutesRouter) GET(relativePath string, handlers ...HandlerFunc) Router {
	return &ginRoutesRouter{routes: r.routes.GET(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginRoutesRouter) POST(relativePath string, handlers ...HandlerFunc) Router {
	return &ginRoutesRouter{routes: r.routes.POST(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginRoutesRouter) PUT(relativePath string, handlers ...HandlerFunc) Router {
	return &ginRoutesRouter{routes: r.routes.PUT(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginRoutesRouter) DELETE(relativePath string, handlers ...HandlerFunc) Router {
	return &ginRoutesRouter{routes: r.routes.DELETE(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginRoutesRouter) PATCH(relativePath string, handlers ...HandlerFunc) Router {
	return &ginRoutesRouter{routes: r.routes.PATCH(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginRoutesRouter) OPTIONS(relativePath string, handlers ...HandlerFunc) Router {
	return &ginRoutesRouter{routes: r.routes.OPTIONS(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginRoutesRouter) HEAD(relativePath string, handlers ...HandlerFunc) Router {
	return &ginRoutesRouter{routes: r.routes.HEAD(relativePath, MapMany(handlers, wrapHandler)...)}
}

type ginRouterGroup struct {
	group        *gin.RouterGroup
	pathRegistry *PathRegistry
	groups       []*ginRouterGroup
}

func newGinRouterGroup(group *gin.RouterGroup) *ginRouterGroup {
	return &ginRouterGroup{
		group:        group,
		pathRegistry: NewPathRegistry(),
		groups:       make([]*ginRouterGroup, 0),
	}
}

func (r *ginRouterGroup) Use(middleware ...HandlerFunc) Router {
	return &ginRoutesRouter{routes: r.group.Use(MapMany(middleware, wrapHandler)...)}
}

func (r *ginRouterGroup) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	g := newGinRouterGroup(r.group.Group(relativePath, MapMany(handlers, wrapHandler)...))
	r.groups = append(r.groups, g)
	return g
}

func (r *ginRouterGroup) Any(relativePath string, handlers ...HandlerFunc) Router {
	r.pathRegistry.Add(relativePath, "GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD")
	return &ginRoutesRouter{routes: r.group.Any(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginRouterGroup) GET(relativePath string, handlers ...HandlerFunc) Router {
	r.pathRegistry.Add(relativePath, "GET")
	return &ginRoutesRouter{routes: r.group.GET(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginRouterGroup) POST(relativePath string, handlers ...HandlerFunc) Router {
	r.pathRegistry.Add(relativePath, "POST")
	return &ginRoutesRouter{routes: r.group.POST(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginRouterGroup) PUT(relativePath string, handlers ...HandlerFunc) Router {
	r.pathRegistry.Add(relativePath, "PUT")
	return &ginRoutesRouter{routes: r.group.PUT(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginRouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) Router {
	r.pathRegistry.Add(relativePath, "DELETE")
	return &ginRoutesRouter{routes: r.group.DELETE(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginRouterGroup) PATCH(relativePath string, handlers ...HandlerFunc) Router {
	r.pathRegistry.Add(relativePath, "PATCH")
	return &ginRoutesRouter{routes: r.group.PATCH(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginRouterGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) Router {
	r.pathRegistry.Add(relativePath, "OPTIONS")
	return &ginRoutesRouter{routes: r.group.OPTIONS(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginRouterGroup) HEAD(relativePath string, handlers ...HandlerFunc) Router {
	r.pathRegistry.Add(relativePath, "HEAD")
	return &ginRoutesRouter{routes: r.group.HEAD(relativePath, MapMany(handlers, wrapHandler)...)}
}

func (r *ginRouterGroup) expandMethods() {
	expandMethods(r, r.pathRegistry)
	for _, g := range r.groups {
		g.expandMethods()
	}
}

func expandMethods(router Router, registry *PathRegistry) {
	for _, p := range registry.Paths() {
		if !registry.Get(p, "GET") {
			router.GET(p, MethodNotAllowed)
		}
		if !registry.Get(p, "POST") {
			router.POST(p, MethodNotAllowed)
		}
		if !registry.Get(p, "PUT") {
			router.PUT(p, MethodNotAllowed)
		}
		if !registry.Get(p, "DELETE") {
			router.DELETE(p, MethodNotAllowed)
		}
		if !registry.Get(p, "PATCH") {
			router.PATCH(p, MethodNotAllowed)
		}
		if !registry.Get(p, "OPTIONS") {
			router.OPTIONS(p, MethodNotAllowed)
		}
		if !registry.Get(p, "HEAD") {
			router.HEAD(p, MethodNotAllowed)
		}
	}
}

type handlerFuncWrapper struct {
	f HandlerFunc
}

func wrapHandler(f HandlerFunc) gin.HandlerFunc {
	wrapper := &handlerFuncWrapper{
		f: f,
	}
	return wrapper.handle
}

func (w *handlerFuncWrapper) handle(c *gin.Context) {
	w.f(wrapContext(c))
}

type contextWrapper struct {
	c *gin.Context
}

func wrapContext(c *gin.Context) Context {
	return &contextWrapper{c: c}
}

func (w *contextWrapper) Get(name string) (any, bool) {
	return w.c.Get(name)
}

func (w *contextWrapper) MustGet(name string) any {
	return w.c.MustGet(name)
}

func (w *contextWrapper) Set(key string, value any) {
	w.c.Set(key, value)
}

func (w *contextWrapper) Query(key string) string {
	return w.c.Query(key)
}

func (w *contextWrapper) QueryArray(key string) []string {
	return w.c.QueryArray(key)
}

func (w *contextWrapper) IntQuery(key string) (int, error) {
	val := w.c.Query(key)
	if len(val) == 0 {
		return 0, nil
	}
	return strconv.Atoi(val)
}

func (w *contextWrapper) BoolQuery(key string) (bool, error) {
	val := w.c.Query(key)
	if len(val) == 0 {
		return false, nil
	}
	return strconv.ParseBool(val)
}

func (w *contextWrapper) Iso8601DateQuery(key string) (*time.Time, error) {
	val := w.c.Query(key)
	if len(val) == 0 {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", val)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (w *contextWrapper) Iso8601DateTimeQuery(key string) (*time.Time, error) {
	val := w.c.Query(key)
	if len(val) == 0 {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (w *contextWrapper) StringQuery(key string) (string, error) {
	val := w.c.Query(key)
	if len(val) == 0 {
		return "", nil
	}
	return url.QueryUnescape(val)
}

func (w *contextWrapper) DefaultQuery(key string, defaultValue string) string {
	return w.c.DefaultQuery(key, defaultValue)
}

func (w *contextWrapper) DefaultIntQuery(key string, defaultValue int) (int, error) {
	val := w.c.Query(key)
	if len(val) == 0 {
		return defaultValue, nil
	}
	return strconv.Atoi(val)
}

func (w *contextWrapper) DefaultBoolQuery(key string, defaultValue bool) (bool, error) {
	val := w.c.Query(key)
	if len(val) == 0 {
		return defaultValue, nil
	}
	return strconv.ParseBool(val)
}

func (w *contextWrapper) DefaultStringQuery(key string, defaultValue string) (string, error) {
	val := w.c.Query(key)
	if len(val) == 0 {
		return defaultValue, nil
	}
	return url.QueryUnescape(val)
}

func (w *contextWrapper) GetHeader(key string) string {
	return w.c.GetHeader(key)
}

func (w *contextWrapper) Param(key string) string {
	return w.c.Param(key)
}

func (w *contextWrapper) GetRawData() ([]byte, error) {
	return w.c.GetRawData()
}

func (w *contextWrapper) MayBindJSON(obj any) bool {
	return w.MayBindJSONV(obj, func() error {
		val, ok := obj.(Validatable)
		if ok {
			return val.Validate()
		}
		return nil
	})
}

func (w *contextWrapper) MayBindJSONV(obj any, validator func() error) bool {
	if err := w.c.ShouldBindJSON(obj); err != nil {
		if err == io.EOF {
			return true
		}
		w.RespondBadRequestE(err)
		return false
	}
	if err := validator(); err != nil {
		w.RespondBadRequestE(err)
		return false
	}
	return true
}

func (w *contextWrapper) MustBindJSON(obj any) bool {
	return w.MustBindJSONV(obj, func() error {
		val, ok := obj.(Validatable)
		if ok {
			return val.Validate()
		}
		return nil
	})
}

func (w *contextWrapper) MustBindJSONV(obj any, validator func() error) bool {
	if err := w.c.ShouldBindJSON(obj); err != nil {
		if err == io.EOF {
			w.RespondMissingRequestBody()
			return false
		}
		w.RespondBadRequestE(err)
		return false
	}
	if err := validator(); err != nil {
		w.RespondBadRequestE(err)
		return false
	}
	val, ok := obj.(Validatable)
	if ok {
		if err := val.Validate(); err != nil {
			w.RespondBadRequestE(err)
			return false
		}
	}
	return true
}

func (w *contextWrapper) Request() *http.Request {
	return w.c.Request
}

func (w *contextWrapper) Writer() http.ResponseWriter {
	return w.c.Writer
}

func (w *contextWrapper) ClientIP() string {
	return w.c.ClientIP()
}

func (w *contextWrapper) RemoteIP() string {
	return w.c.RemoteIP()
}

func (w *contextWrapper) Status(code int) Context {
	w.c.Status(code)
	return w
}

func (w *contextWrapper) String(code int, format string, values ...any) Context {
	w.c.String(code, format, values...)
	return w
}

func (w *contextWrapper) SetHeader(key string, value string) {
	w.c.Writer.Header().Set(key, value)
}

func (w *contextWrapper) SetContentType(value string) {
	w.SetHeader("Content-Type", value)
}

func (w *contextWrapper) Cookie(name string) (string, bool) {
	v, err := w.c.Cookie(name)
	if errors.Is(err, http.ErrNoCookie) {
		return "", false
	}
	return v, true
}

func (w *contextWrapper) SetCookie(name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool) {
	w.c.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}

func (w *contextWrapper) Stream(step func(w io.Writer) bool) bool {
	return w.c.Stream(step)
}

func (w *contextWrapper) SSEvent(name string, message any) {
	w.c.SSEvent(name, message)
}

func (w *contextWrapper) Data(code int, contentType string, data []byte) {
	w.c.Data(code, contentType, data)
}

func (w *contextWrapper) RespondOk(obj any) {
	w.respond(http.StatusOK, obj)
}

func (w *contextWrapper) RespondNoContent() {
	w.c.Status(http.StatusNoContent)
}

func (w *contextWrapper) RespondCreated(obj any) {
	w.respond(http.StatusCreated, obj)
}

func (w *contextWrapper) RespondForbidden(obj any) {
	w.respond(http.StatusForbidden, obj)
}

func (w *contextWrapper) RespondForbiddenE(err error) {
	w.respondE(http.StatusForbidden, err)
}

func (w *contextWrapper) RespondUnauthorized(obj any) {
	w.respond(http.StatusUnauthorized, obj)
}

func (w *contextWrapper) RespondUnauthorizedE(err error) {
	w.respondE(http.StatusUnauthorized, err)
}

func (w *contextWrapper) RespondBadRequest(obj any) {
	w.respond(http.StatusBadRequest, obj)
}

func (w *contextWrapper) RespondBadRequestE(err error) {
	w.respondE(http.StatusBadRequest, err)
}

func (w *contextWrapper) RespondNotFound(obj any) {
	w.respond(http.StatusNotFound, obj)
}

func (w *contextWrapper) RespondNotFoundE(err error) {
	w.respondE(http.StatusNotFound, err)
}

func (w *contextWrapper) RespondConflict(obj any) {
	w.respond(http.StatusConflict, obj)
}

func (w *contextWrapper) RespondConflictE(err error) {
	w.respondE(http.StatusConflict, err)
}

func (w *contextWrapper) RespondInternalServerError(obj any) {
	w.respond(http.StatusInternalServerError, obj)
}

func (w *contextWrapper) RespondInternalServerErrorE(err error) {
	w.respondE(http.StatusInternalServerError, err)
}

func (w *contextWrapper) RespondMissingRequestBody() {
	w.RespondBadRequestE(fmt.Errorf("request body is missing"))
}

func (w *contextWrapper) respond(status int, obj any) {
	if obj == nil {
		w.c.Status(status)
	} else {
		w.c.JSON(status, obj)
	}
}

func (w *contextWrapper) respondE(status int, err error) {
	w.c.JSON(status, gin.H{"error": err.Error()})
}

func (w *contextWrapper) Abort() {
	w.c.Abort()
}

func (w *contextWrapper) AbortWithError(code int, error error) {
	_ = w.c.AbortWithError(code, error)
}

func (w *contextWrapper) Next() {
	w.c.Next()
}

func (w *contextWrapper) HandleError(err error) {
	if e, ok := err.(*ResponseStatusError); ok {
		w.c.JSON(e.StatusCode, gin.H{"error": e.Message})
	} else {
		w.RespondInternalServerErrorE(err)
	}
}

func (w *contextWrapper) Deadline() (deadline time.Time, ok bool) {
	return w.c.Deadline()
}

func (w *contextWrapper) Done() <-chan struct{} {
	return w.c.Done()
}

func (w *contextWrapper) Err() error {
	return w.c.Err()
}

func (w *contextWrapper) Value(key any) any {
	return w.c.Value(key)
}
