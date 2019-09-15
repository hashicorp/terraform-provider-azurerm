package authorizers

import (
	"encoding/base64"
	"net/http"
	"testing"
)

func TestBuildCanonicalizedHeader(t *testing.T) {
	testData := []struct {
		Input    http.Header
		Expected string
	}{
		{
			// no headers
			Expected: "",
			Input: map[string][]string{
				"": {""},
			},
		},
		{
			// no x-ms headers
			Expected: "",
			Input: map[string][]string{
				"panda": {"pops"},
			},
		},
		{
			// only a single x-ms header
			Expected: "x-ms-panda:nom",
			Input: map[string][]string{
				"x-ms-panda": {"nom"},
			},
		},
		{
			// multiple x-ms headers
			Expected: "x-ms-panda:nom\nx-ms-tiger:rawr",
			Input: map[string][]string{
				"x-ms-panda": {"nom"},
				"x-ms-tiger": {"rawr"},
			},
		},
		{
			// multiple x-ms headers, out of order
			Expected: "x-ms-panda:nom\nx-ms-tiger:rawr",
			Input: map[string][]string{
				"x-ms-tiger": {"rawr"},
				"x-ms-panda": {"nom"},
			},
		},
		{
			// mixed headers (some ms, some non-ms)
			Expected: "x-ms-panda:nom\nx-ms-tiger:rawr",
			Input: map[string][]string{
				"x-ms-tiger": {"rawr"},
				"panda":      {"pops"},
				"x-ms-panda": {"nom"},
			},
		},
		{
			// casing
			Expected: "x-ms-panda:nom\nx-ms-tiger:rawr",
			Input: map[string][]string{
				"X-Ms-Tiger": {"rawr"},
				"X-Ms-Panda": {"nom"},
			},
		},
	}

	for _, v := range testData {
		actual := buildCanonicalizedHeader(v.Input)
		if actual != v.Expected {
			t.Fatalf("Expected %q but got %q", v.Expected, actual)
		}
	}
}

func TestBuildCanonicalizedResource(t *testing.T) {
	testData := []struct {
		name          string
		accountName   string
		uri           string
		sharedKeyLite bool
		expected      string
		expectError   bool
	}{
		{
			name:          "invalid uri",
			accountName:   "example",
			uri:           "://example.com",
			sharedKeyLite: true,
			expected:      "",
			expectError:   true,
		},
		{
			name:          "storage emulator doesn't get prefix",
			accountName:   StorageEmulatorAccountName,
			uri:           "http://www.example.com/foo",
			sharedKeyLite: true,
			expected:      "/foo",
		},
		{
			name:          "non storage emulator gets prefix",
			accountName:   StorageEmulatorAccountName + "test",
			uri:           "http://www.example.com/foo",
			sharedKeyLite: true,
			expected:      "/" + StorageEmulatorAccountName + "test/foo",
		},
		{
			name:          "uri encoding",
			accountName:   "example",
			uri:           "<hello>",
			sharedKeyLite: true,
			expected:      "/example%3Chello%3E",
		},
		{
			name:          "comp-arg",
			accountName:   "example",
			uri:           "/endpoint?first=true&comp=bar&second=false&third=panda",
			sharedKeyLite: true,
			expected:      "/example/endpoint?comp=bar",
		},
		{
			name:          "arguments",
			accountName:   "example",
			uri:           "/endpoint?first=true&second=false&third=panda",
			sharedKeyLite: true,
			expected:      "/example/endpoint",
		},
		{
			name:          "arguments-sharedkey",
			accountName:   "example",
			uri:           "/endpoint?first=true&second=false&third=panda",
			sharedKeyLite: false,
			expected:      "/example/endpoint\nfirst:true\nsecond:false\nthird:panda",
			expectError:   false,
		},
		{
			name:          "arguments-sharedkey",
			accountName:   "example",
			uri:           "/endpoint?comp=strawberries&restype=pandas",
			sharedKeyLite: false,
			expected:      "/example/endpoint\ncomp:strawberries\nrestype:pandas",
			expectError:   false,
		},
		{
			name:          "extra-arguments-sharedkey",
			accountName:   "myaccount",
			uri:           "/mycontainer?restype=container&comp=list&include=snapshots&include=metadata&include=uncommittedblobs",
			expected:      "/myaccount/mycontainer\ncomp:list\ninclude:metadata,snapshots,uncommittedblobs\nrestype:container",
			sharedKeyLite: false,
			expectError:   false,
		},
	}

	for _, test := range testData {
		t.Logf("Test %q", test.name)
		actual, err := buildCanonicalizedResource(test.uri, test.accountName, test.sharedKeyLite)
		if err != nil {
			if test.expectError {
				continue
			}

			t.Fatalf("Error: %s", err)
		}

		if *actual != test.expected {
			t.Fatalf("Expected %q but got %q", test.expected, *actual)
		}
	}
}

func TestFormatSharedKeyLiteAuthorizationHeader(t *testing.T) {
	testData := []struct {
		name        string
		accountName string
		accountKey  string
		expected    string
	}{
		{
			name:        "primary",
			accountName: "account1",
			accountKey:  "examplekey",
			expected:    "SharedKeyLite account1:examplekey",
		},
		{
			name:        "secondary",
			accountName: "account1-secondary",
			accountKey:  "examplekey",
			expected:    "SharedKeyLite account1:examplekey",
		},
	}

	for _, test := range testData {
		t.Logf("Test: %q", test.name)
		actual := formatSharedKeyLiteAuthorizationHeader(test.accountName, test.accountKey)

		if actual != test.expected {
			t.Fatalf("Expected %q but got %q", test.expected, actual)
		}
	}
}

func TestHMAC(t *testing.T) {
	testData := []struct {
		Expected            string
		StorageAccountKey   string
		CanonicalizedString string
	}{
		{
			// When Storage Key isn't base-64 encoded
			Expected:            "",
			StorageAccountKey:   "bar",
			CanonicalizedString: "foobarzoo",
		},
		{
			// Valid
			Expected:            "h5U0ATVX6SpbFX1H6GNuxIMeXXCILLoIvhflPtuQZ30=",
			StorageAccountKey:   base64.StdEncoding.EncodeToString([]byte("bar")),
			CanonicalizedString: "foobarzoo",
		},
	}

	for _, v := range testData {
		actual := hmacValue(v.StorageAccountKey, v.CanonicalizedString)
		if actual != v.Expected {
			t.Fatalf("Expected %q but got %q", v.Expected, actual)
		}
	}
}

func TestTestPrepareHeadersForRequest(t *testing.T) {
	request, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	headers := []string{
		"Date",
		"X-Ms-Date",
	}

	for _, header := range headers {
		existingVal := request.Header.Get(header)
		if existingVal != "" {
			t.Fatalf("%q had a value prior to being set: %q", header, existingVal)
		}
	}

	prepareHeadersForRequest(request)

	for _, header := range headers {
		updatedVal := request.Header.Get(header)
		if updatedVal == "" {
			t.Fatalf("%q didn't have a value after being set: %q", header, updatedVal)
		}
	}
}

func TestPrepareHeadersForRequestWithNoneConfigured(t *testing.T) {
	request, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	request.Header = nil
	prepareHeadersForRequest(request)

	if request.Header == nil {
		t.Fatalf("Expected `request.Header` to not be nil, but it was!")
	}
}

func TestPrimaryStorageAccountName(t *testing.T) {
	testData := []struct {
		Expected string
		Input    string
	}{
		{
			// Empty
			Expected: "",
			Input:    "",
		},
		{
			// Primary
			Expected: "bar",
			Input:    "bar",
		},
		{
			// Secondary
			Expected: "bar",
			Input:    "bar-secondary",
		},
	}

	for _, v := range testData {
		actual := primaryStorageAccountName(v.Input)
		if actual != v.Expected {
			t.Fatalf("Expected %q but got %q", v.Expected, actual)
		}
	}
}
