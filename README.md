# Jug

A small shim for building opinionated REST services in go.

[User Guide](#user-guide)

## TL;DR

This library puts a thin layer on top of the [gin](https://github.com/gin-gonic/gin) library.
It's aimed at writing REST based web services that mainly use JSON representations.

## Quick Start

Install dependencies.

```go
go get github.com/cfichtmueller/jug
```

Write your server.

```go
type Message struct {
	M string `json:"message"`
}

// Create a new router
router := jug.New()

// Set up routes

router.GET("/api/hello", func(c jug.Context) {
	c.RespondOk(Message{M: "Hello World"})
})

router.GET("/not-found", func(c jug.Context) {
	c.RespondNotFound(nil) // no response body
})

// Set up method not allowed handlers for all unset verbs on configured routes
router.ExpandMethods()

// Start the server
if err := router.Run("0.0.0.0:8000"); err != nil {
	log.Fatal(err)
}
```

## User Guide

- [Setting up Routes](#setting-up-routes)
- [Expand Methods](#expand-methods)
- [Reading Path Parameters](#reading-path-parameters)
- [Reading Query Parameters](#reading-query-parameters)
- [Reading Request Headers](#reading-headers)
- [Reading Request Body](#reading-request-body)
- [Validating Input](#validating-input)
- [Simple Responses](#simple-responses)
- [Streaming Responses](#streaming-responses)
- [Server Sent Events](#server-sent-events)
- [Cookies](#cookies)
- [Using Middleware](#using-middleware)
- [Using the Context](#using-the-context)
- [Handling Errors](#handling-errors)
- [Debug Mode](#debug-mode)

### Setting up Routes

```go
// set up a router
router := jug.New()

// handle GET /api/users
router.GET("/api/users", func(c jug.Context) {
	c.RespondOk(nil)
})

// create a sub group
projects := router.Group("/api/projects")

// handle GET /api/projects
projects.GET("", handler)
// handle POST /api/projects
projects.POST("", handler)
// handle GET /api/projects/:id
projects.GET("/:id", handler)
```

### Expand Methods

`ExpandMethods` sets up 405 Method Not Allowed handlers for methods on routes that don't have a handler yet.

```go
router := jug.New()
router.GET("/api/users", ...)
router.POST("/api/users", ...)
router.ExpandMethods()
```
```
GET /api/users    -> 200
POST /api/users   -> 201
PUT /api/users    -> 405
DELETE /api/users -> 405
GET /foo/bar      -> 404
```

### Reading Path Parameters

```go
router := jug.New()

router.GET("/api/users/:userId", func(c jug.Context) {
	userId := c.Param("userId")
	c.String(http.StatusOK, "user: %s", userId)
})
```

### Reading Query Parameters

```go
func query(c jug.Context) {
	
	rawValue := c.Query(key)
	
	sliceValue := c.QueryArray(key)
	
	intValue, err := c.IntQuery(key)
	
	boolValue, err := c.BoolQuery(key)
	
	dateValue, err := c.Iso8601DateQuery(key)
	
	dateTimeValue, err := c.Iso8601DateTimeQuery(key)
	
	stringValue, err := c.StringQuery(key)
	
	valueOrDefault := c.DefaultQuery(key, defaultValue)
	
	intValueOrDefault, err := c.DefaultIntQuery(key, defaultValue)
	
	boolValueOrDefault, err := c.DefaultBoolQuery(key, defaultValue)
	
	stringValueOrDefault, err := c.DefaultStringQuery(key, defaultValue)
}
```

### Reading Headers

```go
func headers(c jug.Context) {
    headerValue := c.GetHeader(headerName)
}
```

### Reading Request Body

```go
func rawData(c jug.Context) {
    // get the raw request body
    data, err := c.GetRawData()
}

func mayBindJSON(c jug.Context) {
    var query Query
    if c.MayBindJSON(&query) {
        c.RespondOk(query)
        return
    }
    c.String(http.StatusOk, "No request body given")
}

func mustBindJSON(c jug.Context) {
    var query Query

    // MustBindJSON responds 400 if the binding fails
    if !c.MustBindJSON(&query) {
        return
    }
    c.RespondOk(query)
}

```

### Validating Input

Use the `Validator` to validate data.

```go
message := "Hello World"

err := jug.NewValidator().
	RequireStringnotEmpty(message, "message is required").
	Validate()

if err != nil {
	log.Println("invalid message", err)
}
```

Types that implement the `Validatable` interface are automatically validated when bound from JSON.

```go
type CreateUserRequest struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

func (r CreateUserRequest) Validate() error {
	return jug.NewValidator().
		RequireStringNotEmpty(r.Name, "name is required").
		RequireStringMinLength(r.Name, 2, "name needs at least 2 characters").
		Validate()
}

handler := func(c jug.Context) {
	var req CreateUserRequest
	if !c.MustBindJSON(&req) {
		// if binding or validation fails an appropriate response is generated
		return
    }
	c.RespondCreated(req)
}
```

If you need more control over the validation process you can use the *V functions and provide a validator function your own.
This is useful if you need to access contextual information during the validation.

```go
type CreateUserRequest struct {
	Name string `json:"name"`
	Age int `json:"age"`
	Role string `json:"role""`
}

handler := func(c jug.Context) {
	var req CreateUserRequest
	if !c.MustBindJSONV(&req, func()error {
        return jug.NewValidator().
        RequireEnum(r.Role, loadRolesFromDatabase(), "invalid role")
        Validate()
    }) {
		// if binding or validation fails an appropriate response is generated
		return
    }
	c.RespondCreated(req)
}
```

### Simple Responses

Response methods that take a response body argument marshal the given object to JSON unless specified otherwise.

```go
var c jug.Context

c.Data(statusCode int, contentType string, data []byte)

c.String(statusCode int, format, args...)

c.RespondOk(responseBody any)

c.RespondNoContent()

c.RespondCreated(responseBody any)

c.RespondForbidden(responseBody any)

c.RespondUnauthorized(responseBody any)
c.RespondUnauthorizedE(err error)

c.RespondBadRequest(responseBody any)
c.RespondBadRequestE(err error)

c.RespondNotFound(responseBody any)
c.RedpondNotFoundE(err error)

c.RespondConflict(responseBody any)
c.RespondConflictE(err error)

c.RespondInternalServerError(responseBody any)
c.RespondInternalServerErrorE(err error)

c.RespondMissingRequestBody()
```

### Streaming Responses

Create streaming responses using the `Stream` method.
It takes a step function that is provided with a writer.
Use the writer to write to the response.
Return `true` to keep the response open. The step function will be called again.
Return `false` to end the response and close the connection.

```go
router.GET("/api/stream", func(c jug.Context) {
	i := 10
	c.Stream(func(w io.Writer) bool {
        if i > 0 {
            _, _ = fmt.Fprintf(w, "Remaining: %d\n", i)
            i = i -1
            time.Sleep(time.Second)
            return true
		}
	    return false
    })
})
```

### Server Sent Events

To emit server sent events, use `SSEvent` inside a `Stream` step function.

```go
router.GET("/api/sse", func(c jug.Context) {
    clientChan := make(chan string)
    go func() {
        data := fetchDataWhichTakesSomeTime()
        clientChan <- data
        close(clientChan)
    }()

    c.Stream(func(w io.Writer) bool {
        if msg, ok := <-clientChan; ok {
            c.SSEvent("message", msg)
            return true
        }
        return false
    })
})
```

### Cookies

Getting cookies:

```go
func getCookie(c jug.Context) {
    cookie, err := c.Cookie("my-cookie")
    if err != nil {
        if errors.is(err, http.ErrNoCookie) {
            log.Println("no cookie in request")
            return
        }
        log.Fatal(err)
    }
    log.Println("cookie value", cookie)
}
```

Setting cookies:

```go
func setCookie(c jug.Context) {
	maxAge := 0
	path := ""
	domain := ""
	secure := false
	httOnly := true
	c.SetCookie("my-cookie", "cookies be great", maxAge, path, domain, secure, httpOnly)
}
```

### Using Middleware

Middleware are global handlers that will be executed for every single request.
Typical use cases are loggers, authentication filters, error handlers etc.

A middleware can be registered at the root level or for a route.
It will be executed for all child routes too.

```go
func middleware(c jug.Context) {
	log.Println("hello from the middleware")
}

router := jug.New()
router.Use(middleware)
router.GET("/api/users", handler)
router.GET("/api/projects", handler)

GET /api/users ->  200, hello from the middleware
GET /api/projects -> 200, hello from the middleware
GET /foo/bar -> 404, hello from the middleware

router := jug.New()

usersGroup := router.Group("/api/users", middleware)
usersGroup.GET("", handler)
usersGroup.GET("/:id", handler)
router.GET("/api/projects", handler)

GET / -> 404
GET /api/users -> 200, hello from the middleware
GET /api/users/3 -> 200, hello from the middleware
GET /api/useres/3/foo -> 404, hello from the middleware
GET /api/projects -> 200
```

### Using the Context

Each client request has its own context. Handlers can set and get data to and from the context.

```go
middleware := func(c jug.Context) {
    requestId := c.GetHeader("RequestId")
    if len(requestId) > 0 {
        c.Set("requestId", requestId)
    }
}

handler := func(c jug.Context) {
    requestId, ok := c.Get("requestId")
    if ok {
        c.String(http.StatusOK, "Request id: %s", requestId.(string))
    } else {
        c.String(http.StatusOK, "No Request id given")
    }
}
```

### Handling Errors

Use `HandleError` to inspect an error and to write an appropriate response.

```go
func (c jug.Context) {
	data, err := loadDataFromDatabase()
	if err != nil {
		c.HandleError(err)
		return
    }
	c.RespondOk(data)
}
```

`HandleError` sets appropriate status codes and response messages for:

```go
err := jug.NewResponseStatusError(statusCode, message)

err := jug.NewBadRequestError(message)

err := jug.NewUnauthorizedError(message)

err := jug.NewForbiddenError(message)

err := jug.NewConflictError(message)
```

Unsupported errors lead to an HTTP 500 response.

### Debug Mode

Enables the gin debug mode.

```go
router := jug.New()
router.EnableDebugMode()
```
