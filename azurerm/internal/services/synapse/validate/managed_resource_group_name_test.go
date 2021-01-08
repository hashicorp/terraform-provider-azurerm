package validate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestManagedResourceGroupName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "a",
			Valid: true,
		},
		{
			Input: "1",
			Valid: true,
		},
		{
			Input: "hello",
			Valid: true,
		},
		{
			Input: "Hello",
			Valid: true,
		},
		{
			Input: "hello-world",
			Valid: true,
		},
		{
			Input: "Hello_World",
			Valid: true,
		},
		{
			Input: "HelloWithNumbers12345",
			Valid: true,
		},
		{
			Input: "(Did)You(Know)That(Brackets)Are(Allowed)In(Resource)Group(Names)",
			Valid: true,
		},
		{
			Input: "EndingWithAPeriod.",
			Valid: false,
		},
		{
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			Valid: false,
		},
		{
			Input: acctest.RandString(90),
			Valid: true,
		},
		{
			Input: acctest.RandString(91),
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ManagedResourceGroupName()(tc.Input, "name")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
