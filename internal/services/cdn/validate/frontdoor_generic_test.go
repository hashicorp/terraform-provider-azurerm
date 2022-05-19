package validate

import "testing"

func TestCdnFrontdoorEndpointName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Empty
			Input: "",
			Valid: false,
		},

		{
			// Max Len
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIE",
			Valid: true,
		},

		{
			// Too Long
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIE",
			Valid: false,
		},

		{
			// Invalid Character
			Input: "1%A",
			Valid: false,
		},

		{
			// Start With Hyphen
			Input: "-1",
			Valid: false,
		},

		{
			// End With Hyphen
			Input: "1-",
			Valid: false,
		},

		{
			// Too Short
			Input: "1",
			Valid: false,
		},

		{
			// Start With Letter, End With Letter
			Input: "AA",
			Valid: true,
		},

		{
			// Start With Number, End With Number
			Input: "11",
			Valid: true,
		},

		{
			// Start With Letter, End With Number
			Input: "A1",
			Valid: true,
		},

		{
			// Start With Number, End With Letter
			Input: "1A",
			Valid: true,
		},

		{
			// Start With Letter, End With Letter and Hyphen Separator
			Input: "A-A",
			Valid: true,
		},

		{
			// Start With Number, End With Number and Hyphen Separator
			Input: "1-1",
			Valid: true,
		},

		{
			// Start With Letter, End With Number and Hyphen Separator
			Input: "A-1",
			Valid: true,
		},

		{
			// Start With Number, End With Letter and Hyphen Separator
			Input: "1-A",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := CdnFrontdoorEndpointName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

func TestCdnFrontdoorName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Empty
			Input: "",
			Valid: false,
		},

		{
			// Max Len
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: true,
		},

		{
			// Too Long
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: false,
		},

		{
			// Invalid Character
			Input: "1%A",
			Valid: false,
		},

		{
			// Start With Hyphen
			Input: "-1",
			Valid: false,
		},

		{
			// End With Hyphen
			Input: "1-",
			Valid: false,
		},

		{
			// Too Short
			Input: "1",
			Valid: false,
		},

		{
			// Start With Letter, End With Letter
			Input: "AA",
			Valid: true,
		},

		{
			// Start With Number, End With Number
			Input: "11",
			Valid: true,
		},

		{
			// Start With Letter, End With Number
			Input: "A1",
			Valid: true,
		},

		{
			// Start With Number, End With Letter
			Input: "1A",
			Valid: true,
		},

		{
			// Start With Letter, End With Letter and Hyphen Separator
			Input: "A-A",
			Valid: true,
		},

		{
			// Start With Number, End With Number and Hyphen Separator
			Input: "1-1",
			Valid: true,
		},

		{
			// Start With Letter, End With Number and Hyphen Separator
			Input: "A-1",
			Valid: true,
		},

		{
			// Start With Number, End With Letter and Hyphen Separator
			Input: "1-A",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := CdnFrontdoorName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

func TestCdnFrontdoorOriginGroupName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Empty
			Input: "",
			Valid: false,
		},

		{
			// Max Len
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: true,
		},

		{
			// Too Long
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: false,
		},

		{
			// Invalid Character
			Input: "1%A",
			Valid: false,
		},

		{
			// Start With Hyphen
			Input: "-1",
			Valid: false,
		},

		{
			// End With Hyphen
			Input: "1-",
			Valid: false,
		},

		{
			// Too Short
			Input: "1",
			Valid: false,
		},

		{
			// Start With Letter, End With Letter
			Input: "AA",
			Valid: true,
		},

		{
			// Start With Number, End With Number
			Input: "11",
			Valid: true,
		},

		{
			// Start With Letter, End With Number
			Input: "A1",
			Valid: true,
		},

		{
			// Start With Number, End With Letter
			Input: "1A",
			Valid: true,
		},

		{
			// Start With Letter, End With Letter and Hyphen Separator
			Input: "A-A",
			Valid: true,
		},

		{
			// Start With Number, End With Number and Hyphen Separator
			Input: "1-1",
			Valid: true,
		},

		{
			// Start With Letter, End With Number and Hyphen Separator
			Input: "A-1",
			Valid: true,
		},

		{
			// Start With Number, End With Letter and Hyphen Separator
			Input: "1-A",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := CdnFrontdoorOriginGroupName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

func TestCdnFrontdoorOriginName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Empty
			Input: "",
			Valid: false,
		},

		{
			// Max Len
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: true,
		},

		{
			// Too Long
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: false,
		},

		{
			// Invalid Character
			Input: "1%A",
			Valid: false,
		},

		{
			// Start With Hyphen
			Input: "-1",
			Valid: false,
		},

		{
			// End With Hyphen
			Input: "1-",
			Valid: false,
		},

		{
			// Too Short
			Input: "1",
			Valid: false,
		},

		{
			// Start With Letter, End With Letter
			Input: "AA",
			Valid: true,
		},

		{
			// Start With Number, End With Number
			Input: "11",
			Valid: true,
		},

		{
			// Start With Letter, End With Number
			Input: "A1",
			Valid: true,
		},

		{
			// Start With Number, End With Letter
			Input: "1A",
			Valid: true,
		},

		{
			// Start With Letter, End With Letter and Hyphen Separator
			Input: "A-A",
			Valid: true,
		},

		{
			// Start With Number, End With Number and Hyphen Separator
			Input: "1-1",
			Valid: true,
		},

		{
			// Start With Letter, End With Number and Hyphen Separator
			Input: "A-1",
			Valid: true,
		},

		{
			// Start With Number, End With Letter and Hyphen Separator
			Input: "1-A",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := CdnFrontdoorOriginName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

func TestCdnFrontdoorRouteName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Empty
			Input: "",
			Valid: false,
		},

		{
			// Max Len
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: true,
		},

		{
			// Too Long
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: false,
		},

		{
			// Invalid Character
			Input: "1%A",
			Valid: false,
		},

		{
			// Start With Hyphen
			Input: "-1",
			Valid: false,
		},

		{
			// End With Hyphen
			Input: "1-",
			Valid: false,
		},

		{
			// Too Short
			Input: "1",
			Valid: false,
		},

		{
			// Start With Letter, End With Letter
			Input: "AA",
			Valid: true,
		},

		{
			// Start With Number, End With Number
			Input: "11",
			Valid: true,
		},

		{
			// Start With Letter, End With Number
			Input: "A1",
			Valid: true,
		},

		{
			// Start With Number, End With Letter
			Input: "1A",
			Valid: true,
		},

		{
			// Start With Letter, End With Letter and Hyphen Separator
			Input: "A-A",
			Valid: true,
		},

		{
			// Start With Number, End With Number and Hyphen Separator
			Input: "1-1",
			Valid: true,
		},

		{
			// Start With Letter, End With Number and Hyphen Separator
			Input: "A-1",
			Valid: true,
		},

		{
			// Start With Number, End With Letter and Hyphen Separator
			Input: "1-A",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := CdnFrontdoorRouteName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

func TestCdnFrontdoorCacheDuration(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Empty
			Input: "",
			Valid: false,
		},

		{
			// Not Enough Days
			Input: "0.00:00:00",
			Valid: false,
		},

		{
			// Too Many Days
			Input: "366.00:00:00",
			Valid: false,
		},

		{
			// Max Days
			Input: "365.00:00:00",
			Valid: true,
		},

		{
			// Min Days
			Input: "1.00:00:00",
			Valid: true,
		},

		{
			// Max Duration
			Input: "365.23:59:59",
			Valid: true,
		},

		{
			// Min Duration
			Input: "00:00:01",
			Valid: true,
		},

		{
			// Too Short Duration
			Input: "00:00:00",
			Valid: false,
		},

		{
			// Max Hours
			Input: "23:00:00",
			Valid: true,
		},

		{
			// Too Many Hours
			Input: "24:00:00",
			Valid: false,
		},

		{
			// Max Minutes
			Input: "00:59:00",
			Valid: true,
		},

		{
			// Too Many Minutes
			Input: "00:60:00",
			Valid: false,
		},

		{
			// Max Seconds
			Input: "00:00:59",
			Valid: true,
		},

		{
			// Too Many Seconds
			Input: "00:00:60",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := CdnFrontdoorCacheDuration(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

func TestCdnFrontdoorUrlPathConditionMatchValue(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Empty
			Input: "",
			Valid: true,
		},

		{
			// Starts with slash
			Input: "/foo",
			Valid: false,
		},

		{
			// Does not start with slash
			Input: "foo",
			Valid: true,
		},

		{
			// Has embedded path slash
			Input: "foo/bar",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := CdnFrontdoorUrlPathConditionMatchValue(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

func TestCdnFrontdoorCustomDomainName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Empty
			Input: "",
			Valid: false,
		},

		{
			// Starts with hyphen
			Input: "-foo",
			Valid: false,
		},

		{
			// Ends with hyphen
			Input: "foo-",
			Valid: false,
		},

		{
			// Starts with number
			Input: "1foo",
			Valid: true,
		},

		{
			// Ends with number
			Input: "foo1",
			Valid: true,
		},

		{
			// Has embedded hyphen
			Input: "foo-bar",
			Valid: true,
		},

		{
			// Too short
			Input: "f",
			Valid: false,
		},

		{
			// Min Len
			Input: "fo",
			Valid: true,
		},

		{
			// Max Len
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: true,
		},

		{
			// Too Long
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIII-EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := CdnFrontdoorCustomDomainName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

func TestCdnFrontdoorSecretName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Empty
			Input: "",
			Valid: false,
		},

		{
			// Starts with hyphen
			Input: "-foo",
			Valid: false,
		},

		{
			// Ends with hyphen
			Input: "foo-",
			Valid: false,
		},

		{
			// Starts with number
			Input: "1foo",
			Valid: true,
		},

		{
			// Ends with number
			Input: "foo1",
			Valid: true,
		},

		{
			// Has embedded hyphen
			Input: "foo-bar",
			Valid: true,
		},

		{
			// Too short
			Input: "f",
			Valid: false,
		},

		{
			// Min Len
			Input: "fo",
			Valid: true,
		},

		{
			// Max Len
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: true,
		},

		{
			// Too Long
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIII-EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := CdnFrontdoorSecretName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

func TestLegacyFrontdoorWAFName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Empty
			Input: "",
			Valid: false,
		},

		{
			// Starts with hyphen
			Input: "-foo",
			Valid: false,
		},

		{
			// Ends with hyphen
			Input: "foo-",
			Valid: false,
		},

		{
			// Starts with number
			Input: "1foo",
			Valid: false,
		},

		{
			// Ends with number
			Input: "foo1",
			Valid: true,
		},

		{
			// Has embedded hyphen
			Input: "foo-bar",
			Valid: false,
		},

		{
			// Min Len
			Input: "f",
			Valid: true,
		},

		{
			// Max Len
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: true,
		},

		{
			// Too Long
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE1",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := LegacyFrontdoorWAFName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

func TestLegacyCustomBlockResponseBody(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Empty
			Input: "",
			Valid: false,
		},

		{
			// Valid Base64 String
			Input: "V2VsbCwgc2hlIHR1cm5lZCBtZSBpbnRvIGEgbmV3dC4=",
			Valid: true,
		},

		{
			// Invalid Base64 String
			Input: "QSBuZXd0Pw=",
			Valid: false,
		},

		{
			// Valid Base64 String(punchline)
			Input: "SSBnb3QgYmV0dGVy",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := LegacyCustomBlockResponseBody(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

func TestCdnFrontdoorRuleName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Empty
			Input: "",
			Valid: false,
		},

		{
			// Starts invalid character
			Input: "-foo",
			Valid: false,
		},

		{
			// Ends with invalid character
			Input: "foo-",
			Valid: false,
		},

		{
			// Has embedded invalid character
			Input: "foo-bar",
			Valid: false,
		},

		{
			// Starts with number
			Input: "1foo",
			Valid: false,
		},

		{
			// Ends with number
			Input: "foo1",
			Valid: true,
		},

		{
			// Min Len
			Input: "f",
			Valid: true,
		},

		{
			// Max Len
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: true,
		},

		{
			// Too Long
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := CdnFrontdoorRuleName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

func TestCdnFrontdoorRuleSetName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Empty
			Input: "",
			Valid: false,
		},

		{
			// Starts invalid character
			Input: "-foo",
			Valid: false,
		},

		{
			// Ends with invalid character
			Input: "foo-",
			Valid: false,
		},

		{
			// Has embedded invalid character
			Input: "foo-bar",
			Valid: false,
		},

		{
			// Starts with number
			Input: "1foo",
			Valid: false,
		},

		{
			// Ends with number
			Input: "foo1",
			Valid: true,
		},

		{
			// Min Len
			Input: "f",
			Valid: true,
		},

		{
			// Max Len
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEE",
			Valid: true,
		},

		{
			// Too Long
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEE",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := CdnFrontdoorRuleSetName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

func TestCdnFrontdoorUrlRedirectActionQueryString(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Empty
			Input: "",
			Valid: false,
		},

		{
			// Invalid
			Input: "a",
			Valid: false,
		},

		{
			// Prefix Ampersand
			Input: "&a=b",
			Valid: false,
		},

		{
			// Suffix Ampersand
			Input: "a=b&",
			Valid: false,
		},

		{
			// Prefix Question mark
			Input: "?a=b",
			Valid: false,
		},

		{
			// Suffix Question mark
			Input: "a=b?",
			Valid: false,
		},

		{
			// Expected format
			Input: "a=b",
			Valid: true,
		},

		{
			// Prepend additional query string
			Input: "a=b&c=d",
			Valid: false,
		},
		{
			// Use action server variable
			Input: "a={http_version}",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := CdnFrontdoorUrlRedirectActionQueryString(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

func TestCdnFrontdoorUrlRedirectActionDestinationPath(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Empty
			Input: "",
			Valid: true,
		},

		{
			// Starts with slash
			Input: "/foo",
			Valid: true,
		},

		{
			// Does not start with slash
			Input: "foo",
			Valid: false,
		},

		{
			// Does not start with slash with embedded path slash
			Input: "foo/bar",
			Valid: false,
		},

		{
			// Starts with slash with embedded path slash
			Input: "/foo/bar",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := CdnFrontdoorUrlRedirectActionDestinationPath(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

func TestCdnFrontdoorResourceID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Valid: false,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Valid: false,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Valid: false,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Valid: false,
		},

		{
			// missing ProfileName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/",
			Valid: false,
		},

		{
			// missing value for ProfileName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/",
			Valid: false,
		},

		{
			// missing CustomDomainName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/",
			Valid: false,
		},

		{
			// missing value for CustomDomainName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/customDomains/",
			Valid: false,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/customDomains/customDomain1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.CDN/PROFILES/PROFILE1/CUSTOMDOMAINS/CUSTOMDOMAIN1",
			Valid: false,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Valid: false,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Valid: false,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Valid: false,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Valid: false,
		},

		{
			// missing ProfileName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/",
			Valid: false,
		},

		{
			// missing value for ProfileName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/",
			Valid: false,
		},

		{
			// missing AfdEndpointName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/",
			Valid: false,
		},

		{
			// missing value for AfdEndpointName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/",
			Valid: false,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.CDN/PROFILES/PROFILE1/AFDENDPOINTS/ENDPOINT1",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := CdnFrontdoorResourceID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
