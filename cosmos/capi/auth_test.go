package capi

import (
	"net/http"
	"net/url"
	"testing"
)

func testRequest(method, path, date string) http.Request {
	r := http.Request{
		Method: method,
		URL: &url.URL{
			Path: path,
		},
	}

	r.Header = make(http.Header)
	r.Header.Set("x-ms-date", date)

	return r
}
func TestGenerateAuthorizationSignature(t *testing.T) {
	cases := []struct {
		Request  http.Request
		Key      string
		Expected string
	}{
		{
			Request:  testRequest("GET", "/dbs/SevenDayDB", "Thu, 28 Mar 2019 06:33:32 GMT"),
			Key:      "yvLnqDanONZn10a2ZgUge8cA3P9hkr3elsONP4yW6qADfj1RFkteABYBYEz627UAIUPGDRIeZNjaKqE4mBieqA==",
			Expected: "PaVILGLP+zDsS9kA3YQQNh0OEdZ+zGvaoIxEfZO9TZ8=",
		},
		{
			Request:  testRequest("POST", "/dbs", "Thu, 28 Mar 2019 06:23:48 GMT"),
			Key:      "yvLnqDanONZn10a2ZgUge8cA3P9hkr3elsONP4yW6qADfj1RFkteABYBYEz627UAIUPGDRIeZNjaKqE4mBieqA==",
			Expected: "bQ3HqSBMZk9mXEC8J33HWI6p8a9vrgW0agBIMsBs7Jk=",
		},
	}

	for _, tc := range cases {
		sig, err := GenerateAuthorizationSignature(&tc.Request, tc.Key)

		if err != nil {
			t.Fatalf("GenerateAuthorizationSignature encountered an error: %v", err)
		}

		if sig != tc.Expected {
			t.Fatalf("the signature returned by GenerateAuthorizationSignature did not match: `%s` != `%s`", sig, tc.Expected)
		}
	}
}
