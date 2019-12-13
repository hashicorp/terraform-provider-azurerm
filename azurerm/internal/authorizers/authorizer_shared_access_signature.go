package authorizers

import (
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
)

// SharedAccessSignatureAuthorizer implements an authorization for Shared Access Signature
// this can be used for interaction with Blob, File and Queue Storage Endpoints
type SharedAccessSignatureAuthorizer struct {
	sasQuery map[string]interface{}
}

// NewSharedAccessSignatureAuthorizer creates a SharedAccessSignatureAuthorizer using sasToken
func NewSharedAccessSignatureAuthorizer(sasToken string) *SharedAccessSignatureAuthorizer {
	m, _ := url.ParseQuery(sasToken)
	query := make(map[string]interface{}, len(m))
	for key, value := range m {
		for i, v := range value {
			value[i] = url.QueryEscape(v)
		}
		query[key] = value
	}
	return &SharedAccessSignatureAuthorizer{
		sasQuery: query,
	}
}

// WithAuthorization returns a PrepareDecorator that adds an HTTP Authorization header whose
// value is "SharedKey " followed by the computed key.
// This can be used for the Blob, Queue, and File Services
//
// from: https://docs.microsoft.com/en-us/rest/api/storageservices/delegate-access-with-shared-access-signature
func (skl *SharedAccessSignatureAuthorizer) WithAuthorization() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r, err := p.Prepare(r)
			if err != nil {
				return r, err
			}

			return autorest.Prepare(r, autorest.WithQueryParameters(skl.sasQuery))
		})
	}
}
