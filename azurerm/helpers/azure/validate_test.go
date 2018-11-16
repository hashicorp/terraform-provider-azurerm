package azure

import "testing"

func TestHelper_AzureResourceID(t *testing.T) {
	cases := []struct {
		ID     string
		Errors int
	}{
		{
			ID:     "",
			Errors: 1,
		},
		{
			ID:     "nonsense",
			Errors: 1,
		},
		{
			ID:     "/slash",
			Errors: 1,
		},
		{
			ID:     "/path/to/nothing",
			Errors: 1,
		},
		{
			ID:     "/subscriptions",
			Errors: 1,
		},
		{
			ID:     "/providers",
			Errors: 1,
		},
		{
			ID:     "/subscriptions/not-a-guid",
			Errors: 0,
		},
		{
			ID:     "/providers/test",
			Errors: 0,
		},
		{
			ID:     "/subscriptions/00000000-0000-0000-0000-00000000000/",
			Errors: 0,
		},
		{
			ID:     "/providers/provider.name/",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.ID, func(t *testing.T) {
			_, errors := ValidateResourceID(tc.ID, "test")

			if len(errors) < tc.Errors {
				t.Fatalf("Expected ValidateResourceID to have %d not %d errors for %q", tc.Errors, len(errors), tc.ID)
			}
		})
	}
}

func TestAzureResourceIDOrEmpty(t *testing.T) {
	cases := []struct {
		ID     string
		Errors int
	}{
		{
			ID:     "",
			Errors: 0,
		},
		{
			ID:     "nonsense",
			Errors: 1,
		},
		//as this function just calls TestAzureResourceId lets not be as comprehensive
		{
			ID:     "/providers/provider.name/",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.ID, func(t *testing.T) {
			_, errors := ValidateResourceIDOrEmpty(tc.ID, "test")

			if len(errors) < tc.Errors {
				t.Fatalf("Expected TestAzureResourceIdOrEmpty to have %d not %d errors for %q", tc.Errors, len(errors), tc.ID)
			}
		})
	}
}

func TestAzureValidateMsSqlElasticPoolName(t *testing.T) {
	cases := []struct {
		ElasticPoolName string
		Errors          int
	}{
		{
			ElasticPoolName: "as",
			Errors:          2,
		},
		{
			ElasticPoolName: "Asd",
			Errors:          2,
		},
		{
			ElasticPoolName: "asd",
			Errors:          0,
		},
		{
			ElasticPoolName: "-asd",
			Errors:          1,
		},
		{
			ElasticPoolName: "asd-",
			Errors:          1,
		},
		{
			ElasticPoolName: "asd-1",
			Errors:          0,
		},
		{
			ElasticPoolName: "asd--1",
			Errors:          1,
		},
		{
			ElasticPoolName: "asd--1-",
			Errors:          2,
		},
		{
			ElasticPoolName: "asdfghjklzasdfghjklzasdfghjklzasdfghjklzasdfghjklz",
			Errors:          0,
		},
		{
			ElasticPoolName: "asdfghjklzasdfghjklzasdfghjklzasdfghjklzasdfghjklz1",
			Errors:          1,
		},
		{
			ElasticPoolName: "asdfghjklzasdfghjklzasdfghjklzasdfghjklzasdfghjklz1--",
			Errors:          3,
		},
		{
			ElasticPoolName: "Aasdfghjklzasdfghjklzasdfghjklzasdfghjklzasdfghjklz1--",
			Errors:          4,
		},
	}

	for _, tc := range cases {
		t.Run(tc.ElasticPoolName, func(t *testing.T) {
			_, errors := ValidateMsSqlElasticPoolName(tc.ElasticPoolName, "name")

			if len(errors) < tc.Errors {
				t.Fatalf("Expected TestAzureValidateMsSqlElasticPoolName to have %d not %d errors for %q", tc.Errors, len(errors), tc.ElasticPoolName)
			}
		})
	}
}
