package validate

import (
	s "strings"
	"testing"
)

func TestAzureRMApiManagementServiceName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "a",
			ErrCount: 0,
		},
		{
			Value:    "abc",
			ErrCount: 0,
		},
		{
			Value:    "api1",
			ErrCount: 0,
		},
		{
			Value:    "company-api",
			ErrCount: 0,
		},
		{
			Value:    "hello_world",
			ErrCount: 1,
		},
		{
			Value:    "helloworld21!",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := ApiManagementServiceName(tc.Value, "azurerm_api_management")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Api Management Service Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}

func TestAzureRMApiManagementPublisherName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "",
			ErrCount: 1,
		},
		{
			Value:    "a",
			ErrCount: 0,
		},
		{
			Value:    "abc",
			ErrCount: 0,
		},
		{
			Value:    "api1",
			ErrCount: 0,
		},
		{
			Value:    "company-api",
			ErrCount: 0,
		},
		{
			Value:    "hello_world",
			ErrCount: 0,
		},
		{
			Value:    "helloworld21!",
			ErrCount: 0,
		},
		{
			Value:    "company api",
			ErrCount: 0,
		},
		{
			Value:    "alsdkjflasjkdflajsdlfjkalsdfjkalskdjflajksdflkjasdlfkjasldkfjalksdjflakjsdfljkasdlkfjalskdjfalksdjfdd",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := ApiManagementServicePublisherName(tc.Value, "azurerm_api_management")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Api Management Service Publisher Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}

func TestAzureRMApiManagementApiPath_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "",
			ErrCount: 0,
		},
		{
			Value:    "/",
			ErrCount: 1,
		},
		{
			Value:    "/abc",
			ErrCount: 1,
		},
		{
			Value:    "api1",
			ErrCount: 0,
		},
		{
			Value:    "api1/",
			ErrCount: 1,
		},
		{
			Value:    "api1/sub",
			ErrCount: 0,
		},
		{
			Value:    ".well-known",
			ErrCount: 0,
		},
		{
			Value:    s.Repeat("x", 401),
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := ApiManagementApiPath(tc.Value, "azurerm_api_management_api")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Api Management Api Path to trigger a validation error for '%s'", tc.Value)
		}
	}
}

func TestAzureRMApiManagementApiName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "",
			ErrCount: 1,
		},
		{
			Value:    "asdf+",
			ErrCount: 1,
		},
		{
			Value:    "adsf&",
			ErrCount: 1,
		},
		{
			Value:    "asdfasdf#",
			ErrCount: 1,
		},
		{
			Value:    "asdf*",
			ErrCount: 1,
		},
		{
			Value:    "alksdjl asdlfj laskdjflkjasdlfj lasdf",
			ErrCount: 0,
		},
		{
			Value:    "ddlfj laskdjflkjasdlfj lasdf alksdjflka sdlfjalsdjflajdsflkjasd alsdkjflaksjd flajksdl fjasldkjf lasjdflkajs dfljas ldfjj aljds fljasldkf jalsdjf lakjsdf ljasldkfjalskdjf lakjsd flkajs dlfkja lsdkjf laksdjf lkasjdf lkajsdlfk jasldkfj asldkjfal ksdjf laksjdf",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := ApiManagementApiName(tc.Value, "azurerm_api_management_api")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Api Management Api Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}

func TestApiManagementChildName(t *testing.T) {
	cases := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "Empty",
			input: "",
			valid: false,
		},
		{
			name:  "1",
			input: "1",
			valid: true,
		},
		{
			name:  "v-",
			input: "v-",
			valid: false,
		},
		{
			name:  "v-1",
			input: "v-1",
			valid: true,
		},
		{
			name:  "V1",
			input: "V1",
			valid: true,
		},
		{
			name:  "v_1",
			input: "v_1",
			valid: false,
		},
		{
			name:  "v.1",
			input: "v.1",
			valid: false,
		},
		{
			name:  "v1-",
			input: "v1-",
			valid: false,
		},
		{
			name:  "-v1",
			input: "-v1",
			valid: false,
		},
	}

	for _, tt := range cases {
		_, err := ApiManagementChildName(tt.input, "azurerm_api_management_api_version_set")
		valid := err == nil
		if valid != tt.valid {
			t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
		}
	}
}
