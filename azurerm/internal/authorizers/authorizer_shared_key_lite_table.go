package authorizers

import (
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
)

// TODO: switch to using the version from github.com/Azure/go-autorest
// once https://github.com/Azure/go-autorest/pull/416 has been merged

// SharedKeyLiteTableAuthorizer implements an authorization for Shared Key Lite
// this can be used for interaction with Table Storage Endpoints
type SharedKeyLiteTableAuthorizer struct {
	storageAccountName string
	storageAccountKey  string
}

// NewSharedKeyLiteAuthorizer crates a SharedKeyLiteAuthorizer using the given credentials
func NewSharedKeyLiteTableAuthorizer(accountName, accountKey string) *SharedKeyLiteTableAuthorizer {
	return &SharedKeyLiteTableAuthorizer{
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
func (skl *SharedKeyLiteTableAuthorizer) WithAuthorization() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r, err := p.Prepare(r)
			if err != nil {
				return r, err
			}

			key, err := buildSharedKeyLiteTable(skl.storageAccountName, skl.storageAccountKey, r)
			if err != nil {
				return r, err
			}

			sharedKeyHeader := formatSharedKeyLiteAuthorizationHeader(skl.storageAccountName, *key)
			return autorest.Prepare(r, autorest.WithHeader(HeaderAuthorization, sharedKeyHeader))
		})
	}
}

func buildSharedKeyLiteTable(accountName, storageAccountKey string, r *http.Request) (*string, error) {
	// first ensure the relevant headers are configured
	prepareHeadersForRequest(r)

	sharedKey, err := computeSharedKeyLiteTable(r.URL.String(), accountName, r.Header)
	if err != nil {
		return nil, err
	}

	// we then need to HMAC that value
	hmacdValue := hmacValue(storageAccountKey, *sharedKey)
	return &hmacdValue, nil
}

// computeSharedKeyLite computes the Shared Key Lite required for Storage Authentication
// NOTE: this function assumes that the `x-ms-date` field is set
func computeSharedKeyLiteTable(url string, accountName string, headers http.Header) (*string, error) {
	dateHeader := headers.Get("x-ms-date")
	canonicalizedResource, err := buildCanonicalizedResource(url, accountName, true)
	if err != nil {
		return nil, err
	}

	canonicalizedString := buildCanonicalizedStringForSharedKeyLiteTable(*canonicalizedResource, dateHeader)
	return &canonicalizedString, nil
}

func buildCanonicalizedStringForSharedKeyLiteTable(canonicalizedResource, dateHeader string) string {
	return strings.Join([]string{
		dateHeader,
		canonicalizedResource,
	}, "\n")
}
