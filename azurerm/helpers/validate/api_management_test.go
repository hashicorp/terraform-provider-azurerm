package validate

import (
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
