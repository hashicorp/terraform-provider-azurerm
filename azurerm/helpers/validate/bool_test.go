package validate

import "testing"

func TestBoolIsTrue(t *testing.T) {
	testCases := []struct {
		Value           bool
		ShouldHaveError bool
	}{
		{
			Value:           true,
			ShouldHaveError: false,
		}, {
			Value:           false,
			ShouldHaveError: true,
		},
	}

	t.Run("TestBoolIsTrue", func(t *testing.T) {
		for _, value := range testCases {
			_, errors := BoolIsTrue()(value.Value, "dummy")
			hasErrors := len(errors) > 0

			if value.ShouldHaveError && !hasErrors {
				t.Fatalf("Expected an error but didn't get one for %t", value.Value)
				return
			}

			if !value.ShouldHaveError && hasErrors {
				t.Fatalf("Expected %t to return no errors, but got some %+v", value.Value, errors)
				return
			}
		}
	})
}
