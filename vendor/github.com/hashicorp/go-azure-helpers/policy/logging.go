package policy

import (
	"log"
	"net/http/httputil"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func NewRequestLoggingPolicy(providerName string) azcore.Policy {
	return &requestLoggingPolicy{providerName: providerName}
}

type requestLoggingPolicy struct {
	providerName string
}

func (p requestLoggingPolicy) Do(r *azcore.Request) (*azcore.Response, error) {
	// strip the authorization header prior to printing
	authHeaderName := "Authorization"
	auth := r.Header.Get(authHeaderName)
	if auth != "" {
		r.Header.Del(authHeaderName)
	}

	// dump request to wire format
	if dump, err := httputil.DumpRequestOut(r.Request, true); err == nil {
		log.Printf("[DEBUG] %s Request: \n%s\n", p.providerName, dump)
	} else {
		// fallback to basic message
		log.Printf("[DEBUG] %s Request: %s to %s\n", p.providerName, r.Method, r.URL)
	}

	// add the auth header back
	if auth != "" {
		r.Header.Add(authHeaderName, auth)
	}

	return r.Next()
}
