package authorizers

import (
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
)

// TODO: switch to using the version from github.com/Azure/go-autorest
// once https://github.com/Azure/go-autorest/pull/416 has been merged

// SharedKeyLiteAuthorizer implements an authorization for Shared Key Lite
// this can be used for interaction with Blob, File and Queue Storage Endpoints
type SharedKeyLiteAuthorizer struct {
	storageAccountName string
	storageAccountKey  string
}

// NewSharedKeyLiteAuthorizer creates a SharedKeyLiteAuthorizer using the given credentials
func NewSharedKeyLiteAuthorizer(accountName, accountKey string) *SharedKeyLiteAuthorizer {
	return &SharedKeyLiteAuthorizer{
		storageAccountName: accountName,
		storageAccountKey:  accountKey,
	}
}

// WithAuthorization returns a PrepareDecorator that adds an HTTP Authorization header whose
// value is "SharedKeyLite " followed by the computed key.
// This can be used for the Blob, Queue, and File Services
//
// from: https://docs.microsoft.com/en-us/rest/api/storageservices/authorize-with-shared-key
// You may use Shared Key Lite authorization to authorize a request made against the
// 2009-09-19 version and later of the Blob and Queue services,
// and version 2014-02-14 and later of the File services.
func (skl *SharedKeyLiteAuthorizer) WithAuthorization() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r, err := p.Prepare(r)
			if err != nil {
				return r, err
			}

			key, err := buildSharedKeyLite(skl.storageAccountName, skl.storageAccountKey, r)
			if err != nil {
				return r, err
			}

			sharedKeyHeader := formatSharedKeyLiteAuthorizationHeader(skl.storageAccountName, *key)
			return autorest.Prepare(r, autorest.WithHeader(HeaderAuthorization, sharedKeyHeader))
		})
	}
}
func buildSharedKeyLite(accountName, storageAccountKey string, r *http.Request) (*string, error) {
	// first ensure the relevant headers are configured
	prepareHeadersForRequest(r)

	sharedKey, err := computeSharedKeyLite(r.Method, r.URL.String(), accountName, r.Header)
	if err != nil {
		return nil, err
	}

	// we then need to HMAC that value
	hmacdValue := hmacValue(storageAccountKey, *sharedKey)
	return &hmacdValue, nil
}

// computeSharedKeyLite computes the Shared Key Lite required for Storage Authentication
// NOTE: this function assumes that the `x-ms-date` field is set
func computeSharedKeyLite(verb, url string, accountName string, headers http.Header) (*string, error) {
	canonicalizedResource, err := buildCanonicalizedResource(url, accountName, true)
	if err != nil {
		return nil, err
	}

	canonicalizedHeaders := buildCanonicalizedHeader(headers)
	canonicalizedString := buildCanonicalizedStringForSharedKeyLite(verb, headers, canonicalizedHeaders, *canonicalizedResource)
	return &canonicalizedString, nil
}

func buildCanonicalizedStringForSharedKeyLite(verb string, headers http.Header, canonicalizedHeaders, canonicalizedResource string) string {
	return strings.Join([]string{
		verb,
		headers.Get(HeaderContentMD5),
		headers.Get(HeaderContentType),
		"", // date should be nil, apparently :shrug:
		canonicalizedHeaders,
		canonicalizedResource,
	}, "\n")
}
