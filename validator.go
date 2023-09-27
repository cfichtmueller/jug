// Copyright 2023 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jug

import (
	"fmt"
	"regexp"
	"strings"
)

type Validator struct {
	errors strings.Builder
}

func NewValidator() *Validator {
	return &Validator{
		errors: strings.Builder{},
	}
}

// V invokes a validation function on the validator
func (v *Validator) V(fun func(*Validator)) *Validator {
	fun(v)
	return v
}

// Require requires a condition to be truthy
func (v *Validator) Require(condition bool, message string) *Validator {
	if !condition {
		v.append(message)
	}
	return v
}

// RequireEnum requires a value to be found in a given enum
func (v *Validator) RequireEnum(s string, message string, values ...string) *Validator {
	if len(s) == 0 {
		return v
	}
	for _, value := range values {
		if s == value {
			return v
		}
	}
	v.append(message)
	return v
}

// RequireStringSliceMinLength requires the given slice to have at least min elements
func (v *Validator) RequireStringSliceMinLength(s []string, min int, message string) *Validator {
	if len(s) < min {
		v.append(message)
	}
	return v
}

// RequireStringSliceNotEmpty requires the given slice not to be empty
func (v *Validator) RequireStringSliceNotEmpty(s []string, message string) *Validator {
	if len(s) == 0 {
		v.append(message)
	}
	return v
}

// RequireStringSliceEnum requires the given slice to only contain elements from values...
func (v *Validator) RequireStringSliceEnum(s []string, message string, values ...string) *Validator {
	if len(s) == 0 {
		return v
	}
	m := make(map[string]bool)
	for _, v := range values {
		m[v] = true
	}
	for _, i := range s {
		_, ok := m[i]
		if !ok {
			v.append(message)
			return v
		}
	}
	return v
}

// RequireMatchesRegex requires a value to match a given regular expression
func (v *Validator) RequireMatchesRegex(s string, regex *regexp.Regexp, message string) *Validator {
	if len(s) > 0 && !regex.MatchString(s) {
		v.append(message)
	}
	return v
}

// RequireStringMinLength requires a value to have a given minimum length
func (v *Validator) RequireStringMinLength(s string, min int, message string) *Validator {
	return v.Require(len(s) >= min, message)
}

// RequireStringMaxLength requires a value to have a given maximum length
func (v *Validator) RequireStringMaxLength(s string, max int, message string) *Validator {
	return v.Require(len(s) < max, message)
}

// RequireStringNotEmpty requires a string not to be empty
func (v *Validator) RequireStringNotEmpty(s string, message string) *Validator {
	return v.Require(len(s) > 0, message)
}

// RequireStringLengthBetween requires a string TODO
func (v *Validator) RequireStringLengthBetween(s string, min int, max int, message string) *Validator {
	return v.Require(len(s) >= min && len(s) < max, message)
}

// Validate performs the validation
func (v *Validator) Validate() error {
	if v.errors.Len() > 0 {
		return fmt.Errorf(v.errors.String())
	}
	return nil
}

func (v *Validator) append(msg string) {
	if v.errors.Len() > 0 {
		v.errors.WriteString(", ")
	}
	v.errors.WriteString(msg)
}

// ValidateSub performs validation on a sub item.
func ValidateSub[T Validatable](v *Validator, key string, items []T) *Validator {
	for i, item := range items {
		if err := item.Validate(); err != nil {
			v.append(fmt.Sprintf("%s[%d]: %s", key, i, err.Error()))
		}
	}
	return v
}
