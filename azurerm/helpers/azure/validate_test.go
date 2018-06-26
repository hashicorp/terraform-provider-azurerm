package azure

import "testing"

func TestHelper_Validate_AzureResourceId(t *testing.T) {
	cases := []struct {
		Id     string
		Errors int
	}{
		{
			Id:     "",
			Errors: 1,
		},
		{
			Id:     "nonsense",
			Errors: 1,
		},
		{
			Id:     "/slash",
			Errors: 1,
		},
		{
			Id:     "/path/to/nothing",
			Errors: 1,
		},
		{
			Id:     "/subscriptions",
			Errors: 1,
		},
		{
			Id:     "/providers",
			Errors: 1,
		},
		{
			Id:     "/subscriptions/not-a-guid",
			Errors: 0,
		},
		{
			Id:     "/providers/test",
			Errors: 0,
		},
		{
			Id:     "/subscriptions/00000000-0000-0000-0000-00000000000/",
			Errors: 0,
		},
		{
			Id:     "/providers/provider.name/",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		_, errors := ValidateResourceId(tc.Id, "test")

		if len(errors) < tc.Errors {
			t.Fatalf("Expected AzureResourceId to have an error for %q", tc.Id)
		}
	}
}
