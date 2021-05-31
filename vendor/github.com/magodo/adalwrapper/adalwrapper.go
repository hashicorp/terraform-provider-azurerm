package adalwrapper

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/go-autorest/autorest"
)

// TokenCredential is a wrapper on the autorest.Authorizer, which implements azcore.TokenCredential.
type TokenCredential struct {
	autorest.Authorizer
}

// NewTokenCredential wraps an autorest.Authorizer into a TokenCredential.
func NewTokenCredential(auth autorest.Authorizer) *TokenCredential{
	return &TokenCredential{auth}
}

// AuthenticationPolicy returns a policy that requests the credential and applies it to the HTTP request.
func (cred *TokenCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	// The input options currently only contains a TokenRequestOptions, which in turns only contains the "scope", that is not
	// used in the "old" MS auth endpoint. The similar functionality is covered by the "resource", which is meant to specify
	// during constructing the autorest.BearerAuthorizer.
	return newPolicy(cred)
}

// GetToken requests an access token for the specified set of scopes.
// NOTE: currently this method is not functional, calling it will return error.
func (cred *TokenCredential) GetToken(ctx context.Context, options azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	return nil, fmt.Errorf("GetToken is not supported due to the access token is refreshed internally in adal")
}

// policy implements the azcore.Policy that leverage the internal TokenCredential (actually the autorest.Authorizer internally) to
// requests the credential and applies it to the HTTP request.
type policy struct {
	cred *TokenCredential
}

func newPolicy(cred *TokenCredential) *policy {
	return &policy{cred}
}

func (policy *policy) Do(req *azcore.Request) (*azcore.Response, error) {
	r, err := autorest.Prepare(req.Request, policy.cred.WithAuthorization())
	if err != nil {
		return nil, err
	}
	req.Request = r

	return req.Next()
}
