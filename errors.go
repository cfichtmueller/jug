// Copyright 2023 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jug

import "net/http"

type ResponseStatusError struct {
	StatusCode int
	Message    string
}

func NewResponseStatusError(statusCode int, message string) *ResponseStatusError {
	return &ResponseStatusError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func NewBadRequestError(message string) *ResponseStatusError {
	return NewResponseStatusError(http.StatusBadRequest, message)
}

func NewUnauthorizedError(message string) *ResponseStatusError {
	return NewResponseStatusError(http.StatusUnauthorized, message)
}

func NewForbiddenError(message string) *ResponseStatusError {
	return NewResponseStatusError(http.StatusForbidden, message)
}

func NewNotFoundError(message string) *ResponseStatusError {
	return NewResponseStatusError(http.StatusNotFound, message)
}

func NewConflictError(message string) *ResponseStatusError {
	return NewResponseStatusError(http.StatusConflict, message)
}

func (e *ResponseStatusError) Error() string {
	return e.Message
}
