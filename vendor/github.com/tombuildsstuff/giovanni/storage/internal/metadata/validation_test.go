package metadata

import "testing"

func TestValidationCSharpKeywords(t *testing.T) {
	for key := range cSharpKeywords {
		t.Logf("[DEBUG] Testing %q", key)

		err := Validate(map[string]string{
			key: "value",
		})
		if err == nil {
			t.Fatalf("Expected an error but didn't get one for %q", key)
		}
	}
}

func TestValidation(t *testing.T) {
	testData := []struct {
		Input         string
		ShouldBeValid bool
	}{
		{
			Input:         "",
			ShouldBeValid: false,
		},
		{
			Input:         "abc123",
			ShouldBeValid: true,
		},
		{
			Input:         "_abc123",
			ShouldBeValid: true,
		},
		{
			Input:         "123abc",
			ShouldBeValid: false,
		},
		{
			Input:         "a_123abc",
			ShouldBeValid: true,
		},
		{
			Input:         "abc_123",
			ShouldBeValid: true,
		},
		{
			Input:         "abc123_",
			ShouldBeValid: true,
		},
		{
			Input:         "ABC123",
			ShouldBeValid: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		err := Validate(map[string]string{
			v.Input: "value",
		})
		actual := err == nil
		if v.ShouldBeValid != actual {
			t.Fatalf("Expected %t but got %t for %q", v.ShouldBeValid, actual, v.Input)
		}
	}
}
