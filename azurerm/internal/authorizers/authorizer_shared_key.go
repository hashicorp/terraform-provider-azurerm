package authorizers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
)

// TODO: switch to using the version from github.com/Azure/go-autorest
// once https://github.com/Azure/go-autorest/pull/416 has been merged

// SharedKeyAuthorizer implements an authorization for Shared Key
// this can be used for interaction with Blob, File and Queue Storage Endpoints
type SharedKeyAuthorizer struct {
	storageAccountName string
	storageAccountKey  string
}

// NewSharedKeyAuthorizer creates a SharedKeyAuthorizer using the given credentials
func NewSharedKeyAuthorizer(accountName, accountKey string) *SharedKeyAuthorizer {
	return &SharedKeyAuthorizer{
		storageAccountName: accountName,
		storageAccountKey:  accountKey,
	}
}

// WithAuthorization returns a PrepareDecorator that adds an HTTP Authorization header whose
// value is "SharedKey " followed by the computed key.
// This can be used for the Blob, Queue, and File Services
//
// from: https://docs.microsoft.com/en-us/rest/api/storageservices/authorize-with-shared-key
// You may use Shared Key Lite authorization to authorize a request made against the
// 2009-09-19 version and later of the Blob and Queue services,
// and version 2014-02-14 and later of the File services.
func (skl *SharedKeyAuthorizer) WithAuthorization() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r, err := p.Prepare(r)
			if err != nil {
				return r, err
			}

			key, err := buildSharedKey(skl.storageAccountName, skl.storageAccountKey, r)
			if err != nil {
				return r, err
			}

			sharedKeyHeader := formatSharedKeyAuthorizationHeader(skl.storageAccountName, *key)
			return autorest.Prepare(r, autorest.WithHeader(HeaderAuthorization, sharedKeyHeader))
		})
	}
}
func buildSharedKey(accountName, storageAccountKey string, r *http.Request) (*string, error) {
	// first ensure the relevant headers are configured
	prepareHeadersForRequest(r)

	sharedKey, err := computeSharedKey(r.Method, r.URL.String(), accountName, r.Header)
	if err != nil {
		return nil, err
	}

	// we then need to HMAC that value
	hmacdValue := hmacValue(storageAccountKey, *sharedKey)
	return &hmacdValue, nil
}

// computeSharedKeyLite computes the Shared Key Lite required for Storage Authentication
// NOTE: this function assumes that the `x-ms-date` field is set
func computeSharedKey(verb, url string, accountName string, headers http.Header) (*string, error) {
	canonicalizedResource, err := buildCanonicalizedResource(url, accountName, false)
	if err != nil {
		return nil, err
	}

	canonicalizedHeaders := buildCanonicalizedHeader(headers)
	canonicalizedString := buildCanonicalizedStringForSharedKey(verb, headers, canonicalizedHeaders, *canonicalizedResource)
	return &canonicalizedString, nil
}

func buildCanonicalizedStringForSharedKey(verb string, headers http.Header, canonicalizedHeaders, canonicalizedResource string) string {
	lengthString := headers.Get(HeaderContentLength)

	// empty string when zero
	if lengthString == "0" {
		lengthString = ""
	}

	return strings.Join([]string{
		verb,                                 // HTTP Verb
		headers.Get(HeaderContentEncoding),   // Content-Encoding
		headers.Get(HeaderContentLanguage),   // Content-Language
		lengthString,                         // Content-Length (empty string when zero)
		headers.Get(HeaderContentMD5),        // Content-MD5
		headers.Get(HeaderContentType),       // Content-Type
		"",                                   // date should be nil, apparently :shrug:
		headers.Get(HeaderIfModifiedSince),   // If-Modified-Since
		headers.Get(HeaderIfMatch),           // If-Match
		headers.Get(HeaderIfNoneMatch),       // If-None-Match
		headers.Get(HeaderIfUnmodifiedSince), // If-Unmodified-Since
		headers.Get(HeaderRange),             // Range
		canonicalizedHeaders,
		canonicalizedResource,
	}, "\n")
}

func formatSharedKeyAuthorizationHeader(accountName, key string) string {
	canonicalizedAccountName := primaryStorageAccountName(accountName)
	return fmt.Sprintf("SharedKey %s:%s", canonicalizedAccountName, key)
}
