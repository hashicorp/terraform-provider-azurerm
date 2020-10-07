package validate

import "testing"

func TestCosmosAccountName(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			Value:  "foo-bar",
			Errors: 0,
		},
		{
			Value:  "foo",
			Errors: 0,
		},
		{
			Value:  "fu",
			Errors: 1,
		},
		{
			Value:  "foo_bar",
			Errors: 1,
		},
		{
			Value:  "fooB@r",
			Errors: 1,
		},
		{
			Value:  "foo-bar-foo-bar-foo-bar-foo-bar-foo-bar-foo-bar-foo-bar",
			Errors: 1,
		},
	}

	for _, tc := range cases {
		_, errors := CosmosAccountName(tc.Value, "throughput")
		if len(errors) != tc.Errors {
			t.Fatalf("Expected CosmosAccountName to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}

func TestCosmosEntityName(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			Value:  "",
			Errors: 1,
		},
		{
			Value:  "someEntityName",
			Errors: 0,
		},
		{
			Value:  "someEntityNamesomeEntityNamesomeEntityNamesomeEntityNamesomeEntityNamesomeEntityNamesomeEntityNamesomeEntityNamesomeEntityNamesomeEntityNamesomeEntityNamesomeEntityNamesomeEntityNamesomeEntityNamesomeEntityNamesomeEntityNamesomeEntityNamesomeEntityNamesomeEntityName",
			Errors: 1,
		},
	}

	for _, tc := range cases {
		_, errors := CosmosEntityName(tc.Value, "throughput")
		if len(errors) != tc.Errors {
			t.Fatalf("Expected CosmosEntityName to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}

func TestCosmosThroughput(t *testing.T) {
	cases := []struct {
		Value  int
		Errors int
	}{
		{
			Value:  400,
			Errors: 0,
		},
		{
			Value:  300,
			Errors: 1,
		},
		{
			Value:  450,
			Errors: 1,
		},
		{
			Value:  10000,
			Errors: 0,
		},
	}

	for _, tc := range cases {
		_, errors := CosmosThroughput(tc.Value, "throughput")
		if len(errors) != tc.Errors {
			t.Fatalf("Expected CosmosThroughput to trigger '%d' errors for '%d' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}

func TestCosmosMaxThroughput(t *testing.T) {
	cases := []struct {
		Value  interface{}
		Errors int
	}{
		{
			Value:  400,
			Errors: 2,
		},
		{
			Value:  1000,
			Errors: 1,
		},
		{
			Value:  4000,
			Errors: 0,
		},
		{
			Value:  4001,
			Errors: 1,
		},
		{
			Value:  10000,
			Errors: 0,
		},
		{
			Value:  54000,
			Errors: 0,
		},
		{
			Value:  1000000,
			Errors: 0,
		},
		{
			Value:  1100000,
			Errors: 1,
		},
		{
			Value:  "400",
			Errors: 1,
		},
	}

	for _, tc := range cases {
		_, errors := CosmosMaxThroughput(tc.Value, "throughput")
		if len(errors) != tc.Errors {
			t.Fatalf("Expected CosmosMaxThroughput to trigger '%d' errors for '%d' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}
