package validate

import "testing"

func TestValidateStorageShareDirectoryName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    "abc123",
			Expected: true,
		},
		{
			Input:    "123abc",
			Expected: true,
		},
		{
			Input:    "abc-123",
			Expected: true,
		},
		{
			Input:    "123-abc",
			Expected: true,
		},
		{
			Input:    "-123-abc",
			Expected: true,
		},
		{
			Input:    "abc-123-",
			Expected: true,
		},
		{
			Input:    "abc--123",
			Expected: true,
		},
		{
			Input:    "hello/world",
			Expected: true,
		},
		{
			Input:    "hello/-world",
			Expected: true,
		},
		{
			Input:    "hello-/world",
			Expected: true,
		},
		{
			Input:    "abc-123/world-",
			Expected: true,
		},
		{
			Input:    "123-abc/hello-world",
			Expected: true,
		},
		{
			Input:    "abc-123/hello--world",
			Expected: true,
		},
		{
			Input:    "abc-123/hello/world",
			Expected: true,
		},
		{
			Input:    "abc-123/hello/world-",
			Expected: true,
		},
		{
			Input:    "hello/",
			Expected: false,
		},
		{
			Input:    "abc-test-123/hjg-345-test",
			Expected: true,
		},
		{
			Input:    "Abc123",
			Expected: true,
		},
		{
			Input:    "abc123A",
			Expected: true,
		},
		{
			Input:    "abC123",
			Expected: true,
		},
	}

	for _, v := range testCases {
		t.Logf("[DEBUG] Test Input %q", v.Input)

		warnings, errors := StorageShareDirectoryName(v.Input, "name")
		if len(warnings) != 0 {
			t.Fatalf("Expected no warnings but got %d", len(warnings))
		}

		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
