// Copyright 2023 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jug

import "testing"

func TestNewValidator(t *testing.T) {
	v := NewValidator()
	if v == nil {
		t.Fatal("Expected to get a validator")
	}
}

func TestValidator_Require(t *testing.T) {
	e1 := NewValidator().Require(false, "message").Validate()
	if e1 == nil {
		t.Fatal("Require(false) should return an error")
	}
	if e1.Error() != "message" {
		t.Fatal("error should contain the provided message, got", e1.Error())
	}

	e2 := NewValidator().Require(true, "m").Validate()
	if e2 != nil {
		t.Fatal("Require(true) should not return an error, got", e2)
	}
}

func TestValidator_RequireEnum(t *testing.T) {
	if err := NewValidator().RequireEnum("a", "message", "a", "b").Validate(); err != nil {
		t.Fatal("RequireEnum() should not fail when given an enum value, got", err)
	}

	err := NewValidator().RequireEnum("c", "message", "a", "b").Validate()
	if err == nil {
		t.Fatal("RequireEnum() should fail when given a non enum value")
	}
	if err.Error() != "message" {
		t.Fatal("error should contain the provided message, got", err.Error())
	}
}
