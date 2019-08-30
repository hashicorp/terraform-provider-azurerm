package authorizers

import (
	"net/http"
	"testing"
)

func TestBuildCanonicalizedStringForSharedKeyLite(t *testing.T) {
	testData := []struct {
		name                  string
		headers               map[string][]string
		canonicalizedHeaders  string
		canonicalizedResource string
		verb                  string
		expected              string
	}{
		{
			name: "completed",
			verb: "NOM",
			headers: map[string][]string{
				"Content-MD5":  {"abc123"},
				"Content-Type": {"vnd/panda-pops+v1"},
			},
			canonicalizedHeaders:  "all-the-headers",
			canonicalizedResource: "all-the-resources",
			expected:              "NOM\n\nvnd/panda-pops+v1\n\nall-the-headers\nall-the-resources",
		},
	}

	for _, test := range testData {
		t.Logf("Test: %q", test.name)
		actual := buildCanonicalizedStringForSharedKeyLite(test.verb, test.headers, test.canonicalizedHeaders, test.canonicalizedResource)
		if actual != test.expected {
			t.Fatalf("Expected %q but got %q", test.expected, actual)
		}
	}
}

func TestComputeSharedKey(t *testing.T) {
	testData := []struct {
		name        string
		accountName string
		method      string
		url         string
		headers     map[string]string
		expected    string
	}{
		{
			name:        "No Path",
			accountName: "unlikely23exst2acct23wi",
			method:      "GET",
			url:         "https://unlikely23exst2acct23wi.queue.core.windows.net?comp=properties&restype=service",
			headers: map[string]string{
				"Content-Type": "application/xml; charset=utf-8",
				"X-Ms-Date":    "Wed, 21 Aug 2019 11:00:25 GMT",
				"X-Ms-Version": "2018-11-09",
			},
			expected: `GET

application/xml; charset=utf-8

x-ms-date:Wed, 21 Aug 2019 11:00:25 GMT
x-ms-version:2018-11-09
/unlikely23exst2acct23wi?comp=properties`,
		},
		{
			name:        "With Path",
			accountName: "unlikely23exst2accti1t0",
			method:      "GET",
			url:         "https://unlikely23exst2accti1t0.queue.core.windows.net/?comp=properties&restype=service",
			headers: map[string]string{
				"Content-Type": "application/xml; charset=utf-8",
				"X-Ms-Date":    "Wed, 21 Aug 2019 11:53:48 GMT",
				"X-Ms-Version": "2018-11-09",
			},
			expected: `GET

application/xml; charset=utf-8

x-ms-date:Wed, 21 Aug 2019 11:53:48 GMT
x-ms-version:2018-11-09
/unlikely23exst2accti1t0/?comp=properties`,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Test %q", v.name)

		headers := http.Header{}
		for hn, hv := range v.headers {
			headers.Add(hn, hv)
		}

		actual, err := computeSharedKeyLite(v.method, v.url, v.accountName, headers)
		if err != nil {
			t.Fatalf("Error computing shared key: %s", err)
		}

		if *actual != v.expected {
			t.Fatalf("Expected %q but got %q", v.expected, *actual)
		}
	}
}
