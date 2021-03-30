// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	imdsEndpoint = "http://169.254.169.254/metadata/identity/oauth2/token"
)

const (
	arcIMDSEndpoint    = "IMDS_ENDPOINT"
	identityEndpoint   = "IDENTITY_ENDPOINT"
	identityHeader     = "IDENTITY_HEADER"
	msiEndpoint        = "MSI_ENDPOINT"
	msiSecret          = "MSI_SECRET"
	imdsAPIVersion     = "2018-02-01"
	azureArcAPIVersion = "2019-08-15"
)

type msiType int

const (
	msiTypeUnknown             msiType = 0
	msiTypeIMDS                msiType = 1
	msiTypeAppServiceV20170901 msiType = 2
	msiTypeCloudShell          msiType = 3
	msiTypeUnavailable         msiType = 4
	msiTypeAppServiceV20190801 msiType = 5
	msiTypeAzureArc            msiType = 6
)

// managedIdentityClient provides the base for authenticating in managed identity environments
// This type includes an azcore.Pipeline and TokenCredentialOptions.
type managedIdentityClient struct {
	pipeline               azcore.Pipeline
	imdsAPIVersion         string
	imdsAvailableTimeoutMS time.Duration
	msiType                msiType
	endpoint               string
}

type wrappedNumber json.Number

func (n *wrappedNumber) UnmarshalJSON(b []byte) error {
	c := string(b)
	if c == "\"\"" {
		return nil
	}
	return json.Unmarshal(b, (*json.Number)(n))
}

// newManagedIdentityClient creates a new instance of the ManagedIdentityClient with the ManagedIdentityCredentialOptions
// that are passed into it along with a default pipeline.
// options: ManagedIdentityCredentialOptions configure policies for the pipeline and the authority host that
// will be used to retrieve tokens and authenticate
func newManagedIdentityClient(options *ManagedIdentityCredentialOptions) *managedIdentityClient {
	logEnvVars()
	return &managedIdentityClient{
		pipeline:               newDefaultMSIPipeline(*options), // a pipeline that includes the specific requirements for MSI authentication, such as custom retry policy options
		imdsAPIVersion:         imdsAPIVersion,                  // this field will be set to whatever value exists in the constant and is used when creating requests to IMDS
		imdsAvailableTimeoutMS: 500,                             // we allow a timeout of 500 ms since the endpoint might be slow to respond
		msiType:                msiTypeUnknown,                  // when creating a new managedIdentityClient, the current MSI type is unknown and will be tested for and replaced once authenticate() is called from GetToken on the credential side
	}
}

// authenticate creates an authentication request for a Managed Identity and returns the resulting Access Token if successful.
// ctx: The current context for controlling the request lifetime.
// clientID: The client (application) ID of the service principal.
// scopes: The scopes required for the token.
func (c *managedIdentityClient) authenticate(ctx context.Context, clientID string, scopes []string) (*azcore.AccessToken, error) {
	msg, err := c.createAuthRequest(ctx, clientID, scopes)
	if err != nil {
		return nil, err
	}

	resp, err := c.pipeline.Do(msg)
	if err != nil {
		return nil, err
	}

	if resp.HasStatusCode(successStatusCodes[:]...) {
		return c.createAccessToken(resp)
	}

	return nil, &AuthenticationFailedError{inner: newAADAuthenticationFailedError(resp)}
}

func (c *managedIdentityClient) createAccessToken(res *azcore.Response) (*azcore.AccessToken, error) {
	value := struct {
		// these are the only fields that we use
		Token        string        `json:"access_token,omitempty"`
		RefreshToken string        `json:"refresh_token,omitempty"`
		ExpiresIn    wrappedNumber `json:"expires_in,omitempty"` // this field should always return the number of seconds for which a token is valid
		ExpiresOn    string        `json:"expires_on,omitempty"` // the value returned in this field varies between a number and a date string
	}{}
	if err := res.UnmarshalAsJSON(&value); err != nil {
		return nil, fmt.Errorf("internal AccessToken: %w", err)
	}
	if value.ExpiresIn != "" {
		expiresIn, err := json.Number(value.ExpiresIn).Int64()
		if err != nil {
			return nil, err
		}
		return &azcore.AccessToken{Token: value.Token, ExpiresOn: time.Now().Add(time.Second * time.Duration(expiresIn)).UTC()}, nil
	}
	if expiresOn, err := strconv.Atoi(value.ExpiresOn); err == nil {
		return &azcore.AccessToken{Token: value.Token, ExpiresOn: time.Now().Add(time.Second * time.Duration(expiresOn)).UTC()}, nil
	}
	// this is the case when expires_on is a time string
	// this is the format of the string coming from the service
	if expiresOn, err := time.Parse("1/2/2006 15:04:05 PM +00:00", value.ExpiresOn); err == nil { // the date string specified is for Windows OS
		eo := expiresOn.UTC()
		return &azcore.AccessToken{Token: value.Token, ExpiresOn: eo}, nil
	} else if expiresOn, err := time.Parse("1/2/2006 15:04:05 +00:00", value.ExpiresOn); err == nil { // the date string specified is for Linux OS
		eo := expiresOn.UTC()
		return &azcore.AccessToken{Token: value.Token, ExpiresOn: eo}, nil
	} else {
		return nil, err
	}
}

func (c *managedIdentityClient) createAuthRequest(ctx context.Context, clientID string, scopes []string) (*azcore.Request, error) {
	switch c.msiType {
	case msiTypeIMDS:
		return c.createIMDSAuthRequest(ctx, scopes)
	case msiTypeAppServiceV20170901, msiTypeAppServiceV20190801:
		return c.createAppServiceAuthRequest(ctx, clientID, scopes)
	case msiTypeAzureArc:
		// need to perform preliminary request to retreive the secret key challenge provided by the HIMDS service
		key, err := c.getAzureArcSecretKey(ctx, scopes)
		if err != nil {
			return nil, &AuthenticationFailedError{inner: err, msg: "Failed to retreive secret key from the identity endpoint."}
		}
		return c.createAzureArcAuthRequest(ctx, key, scopes)
	case msiTypeCloudShell:
		return c.createCloudShellAuthRequest(ctx, clientID, scopes)
	default:
		errorMsg := ""
		switch c.msiType {
		case msiTypeUnavailable:
			errorMsg = "unavailable"
		default:
			errorMsg = "unknown"
		}
		return nil, &CredentialUnavailableError{credentialType: "Managed Identity Credential", message: "Make sure you are running in a valid Managed Identity Environment. Status: " + errorMsg}
	}
}

func (c *managedIdentityClient) createIMDSAuthRequest(ctx context.Context, scopes []string) (*azcore.Request, error) {
	request, err := azcore.NewRequest(ctx, http.MethodGet, c.endpoint)
	if err != nil {
		return nil, err
	}
	request.Header.Set(azcore.HeaderMetadata, "true")
	q := request.URL.Query()
	q.Add("api-version", c.imdsAPIVersion)
	q.Add("resource", strings.Join(scopes, " "))
	request.URL.RawQuery = q.Encode()
	return request, nil
}

func (c *managedIdentityClient) createAppServiceAuthRequest(ctx context.Context, clientID string, scopes []string) (*azcore.Request, error) {
	request, err := azcore.NewRequest(ctx, http.MethodGet, c.endpoint)
	if err != nil {
		return nil, err
	}
	q := request.URL.Query()
	if c.msiType == msiTypeAppServiceV20170901 {
		request.Header.Set("secret", os.Getenv(msiSecret))
		q.Add("api-version", "2017-09-01")
		q.Add("resource", strings.Join(scopes, " "))
		if clientID != "" {
			// the legacy 2017 API version specifically specifies "clientid" and not "client_id" as a query param
			q.Add("clientid", clientID)
		}
	} else if c.msiType == msiTypeAppServiceV20190801 {
		request.Header.Set("X-IDENTITY-HEADER", os.Getenv(identityHeader))
		q.Add("api-version", "2019-08-01")
		q.Add("resource", scopes[0])
		if clientID != "" {
			q.Add(qpClientID, clientID)
		}
	}

	request.URL.RawQuery = q.Encode()
	return request, nil
}

func (c *managedIdentityClient) getAzureArcSecretKey(ctx context.Context, resources []string) (string, error) {
	// create the request to retreive the secret key challenge provided by the HIMDS service
	request, err := azcore.NewRequest(ctx, http.MethodGet, c.endpoint)
	if err != nil {
		return "", err
	}
	request.Header.Set(azcore.HeaderMetadata, "true")
	q := request.URL.Query()
	q.Add("api-version", azureArcAPIVersion)
	q.Add("resource", strings.Join(resources, " "))
	request.URL.RawQuery = q.Encode()
	// send the initial request to get the short-lived secret key
	response, err := c.pipeline.Do(request)
	if err != nil {
		return "", err
	}
	// the endpoint is expected to return a 401 with the WWW-Authenticte header set to the location
	// of the secret key file. Any other status code indicates an error in the request.
	if response.StatusCode != 401 {
		return "", &AuthenticationFailedError{inner: newAADAuthenticationFailedError(response), msg: fmt.Sprintf("Expected a 401 Unauthorized response, received: %d", response.StatusCode)}
	}
	header := response.Header.Get("WWW-Authenticate")
	if len(header) == 0 {
		return "", errors.New("Did not receive a value from WWW-Authenticate header")
	}
	// the WWW-Authenticate header is expected in the following format: Basic realm=/some/file/path.key
	pos := strings.LastIndex(header, "=")
	if pos == -1 {
		return "", fmt.Errorf("Did not receive a correct value from WWW-Authenticate header: %s", header)
	}
	key, err := ioutil.ReadFile(header[pos+1:])
	if err != nil {
		return "", fmt.Errorf("Could not read file (%s) contents: %w", header[pos+1:], err)
	}
	return string(key), nil
}

func (c *managedIdentityClient) createAzureArcAuthRequest(ctx context.Context, key string, resources []string) (*azcore.Request, error) {
	request, err := azcore.NewRequest(ctx, http.MethodGet, c.endpoint)
	if err != nil {
		return nil, err
	}
	request.Header.Set(azcore.HeaderMetadata, "true")
	request.Header.Set(azcore.HeaderAuthorization, fmt.Sprintf("Basic %s", key))
	q := request.URL.Query()
	q.Add("api-version", azureArcAPIVersion)
	q.Add("resource", strings.Join(resources, " "))
	request.URL.RawQuery = q.Encode()
	return request, nil
}

func (c *managedIdentityClient) createCloudShellAuthRequest(ctx context.Context, clientID string, scopes []string) (*azcore.Request, error) {
	request, err := azcore.NewRequest(ctx, http.MethodPost, c.endpoint)
	if err != nil {
		return nil, err
	}
	request.Header.Set(azcore.HeaderMetadata, "true")
	data := url.Values{}
	data.Set("resource", strings.Join(scopes, " "))
	if clientID != "" {
		data.Set("client_id", clientID)
	}
	dataEncoded := data.Encode()
	body := azcore.NopCloser(strings.NewReader(dataEncoded))
	if err := request.SetBody(body, azcore.HeaderURLEncoded); err != nil {
		return nil, err
	}
	return request, nil
}

func (c *managedIdentityClient) getMSIType() (msiType, error) {
	if c.msiType == msiTypeUnknown { // if we haven't already determined the msiType
		if endpointEnvVar := os.Getenv(msiEndpoint); endpointEnvVar != "" { // if the env var MSI_ENDPOINT is set
			c.endpoint = endpointEnvVar
			if secretEnvVar := os.Getenv(msiSecret); secretEnvVar != "" { // if BOTH the env vars MSI_ENDPOINT and MSI_SECRET are set the msiType is AppService
				c.msiType = msiTypeAppServiceV20170901
			} else { // if ONLY the env var MSI_ENDPOINT is set the msiType is CloudShell
				c.msiType = msiTypeCloudShell
			}
		} else if endpointEnvVar := os.Getenv(identityEndpoint); endpointEnvVar != "" { // check for IDENTITY_ENDPOINT
			c.endpoint = endpointEnvVar
			if header := os.Getenv(identityHeader); header != "" { // if BOTH the env vars IDENTITY_ENDPOINT and IDENTITY_HEADER are set the msiType is AppService
				c.msiType = msiTypeAppServiceV20190801
			} else if arcIMDS := os.Getenv(arcIMDSEndpoint); arcIMDS != "" {
				c.msiType = msiTypeAzureArc
			} else {
				c.msiType = msiTypeUnavailable
				return c.msiType, &CredentialUnavailableError{credentialType: "Managed Identity Credential", message: "This Managed Identity Environment is not supported yet"}
			}
		} else if c.imdsAvailable() { // if MSI_ENDPOINT is NOT set AND the IMDS endpoint is available the msiType is IMDS. This will timeout after 500 milliseconds
			c.endpoint = imdsEndpoint
			c.msiType = msiTypeIMDS
		} else { // if MSI_ENDPOINT is NOT set and IMDS endpoint is not available Managed Identity is not available
			c.msiType = msiTypeUnavailable
			return c.msiType, &CredentialUnavailableError{credentialType: "Managed Identity Credential", message: "Make sure you are running in a valid Managed Identity Environment"}
		}
	}
	return c.msiType, nil
}

// performs an I/O request that has a timeout of 500 milliseconds
func (c *managedIdentityClient) imdsAvailable() bool {
	tempCtx, cancel := context.WithTimeout(context.Background(), c.imdsAvailableTimeoutMS*time.Millisecond)
	defer cancel()
	// this should never fail
	request, _ := azcore.NewRequest(tempCtx, http.MethodGet, imdsEndpoint)
	q := request.URL.Query()
	q.Add("api-version", c.imdsAPIVersion)
	request.URL.RawQuery = q.Encode()
	resp, err := c.pipeline.Do(request)
	if err == nil {
		resp.Drain()
	}
	return err == nil
}
