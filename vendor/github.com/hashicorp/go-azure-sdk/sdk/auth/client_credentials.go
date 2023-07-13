// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package auth

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/go-uuid"
	"golang.org/x/oauth2"
)

type clientCredentialsType string

const (
	clientCredentialsAssertionType clientCredentialsType = "ClientCredentials"
	clientCredentialsSecretType    clientCredentialsType = "ClientSecret"
)

// clientCredentialsConfig is the configuration for using client credentials flow.
//
// For more information see:
// https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-client-creds-grant-flow#get-a-token
// https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-certificate-credentials
type clientCredentialsConfig struct {
	// Environment is the national cloud environment to use
	Environment environments.Environment

	// TenantID is the required tenant ID for the primary token
	TenantID string

	// AuxiliaryTenantIDs is an optional list of tenant IDs for which to obtain additional tokens
	AuxiliaryTenantIDs []string

	// ClientID is the application's ID.
	ClientID string

	// ClientSecret is the application's secret.
	ClientSecret string

	// PrivateKey contains the contents of an RSA private key or the
	// contents of a PEM file that contains a private key. The provided
	// private key is used to sign JWT assertions.
	// PEM containers with a passphrase are not supported.
	// Use the following command to convert a PKCS 12 file into a PEM.
	//
	//    $ openssl pkcs12 -in key.p12 -out key.pem -nodes
	//
	PrivateKey crypto.PrivateKey

	// Certificate contains the (optionally PEM encoded) X509 certificate registered
	// for the application with which you are authenticating. Used when FederatedAssertion is empty.
	Certificate *x509.Certificate

	// FederatedAssertion contains a JWT provided by a trusted third-party vendor
	// for obtaining an access token with a federated credential. When empty, an
	// assertion will be created and signed using the specified PrivateKey and Certificate
	FederatedAssertion string

	// Scopes specifies a list of requested permission scopes (used for v2 tokens)
	Scopes []string

	// TokenURL is the clientCredentialsToken endpoint, which overrides the default endpoint constructed from a tenant ID
	TokenURL string

	// Audience optionally specifies the intended audience of the
	// request.  If empty, the value of TokenURL is used as the
	// intended audience.
	Audience string
}

// TokenSource provides a source for obtaining access tokens using ClientAssertionAuthorizer or ClientSecretAuthorizer.
func (c *clientCredentialsConfig) TokenSource(_ context.Context, authType clientCredentialsType) (Authorizer, error) {
	switch authType {
	case clientCredentialsAssertionType:
		return NewCachedAuthorizer(&ClientAssertionAuthorizer{
			conf: c,
		})
	case clientCredentialsSecretType:
		return NewCachedAuthorizer(&ClientSecretAuthorizer{
			conf: c,
		})
	}
	return nil, fmt.Errorf("internal-error: unimplemented authType %q", string(authType))
}

type clientAssertionTokenHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
	KeyId     string `json:"kid"`
}

func (h *clientAssertionTokenHeader) encode() (string, error) {
	b, err := json.Marshal(h)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

type clientAssertionTokenClaims struct {
	Audience  string `json:"aud"`
	Expiry    int64  `json:"exp"`
	Issuer    string `json:"iss"`
	JwtId     string `json:"jti"`
	NotBefore int64  `json:"nbf"`
	Subject   string `json:"sub"`
}

func (c *clientAssertionTokenClaims) encode() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

type clientAssertionToken struct {
	header clientAssertionTokenHeader
	claims clientAssertionTokenClaims
}

func (c *clientAssertionToken) encode(key crypto.PrivateKey) (*string, error) {
	var err error

	c.claims.NotBefore = time.Now().Unix()
	c.claims.Expiry = time.Now().Add(time.Hour).Unix()
	c.claims.JwtId, err = uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	var hash = crypto.SHA256
	var sign func([]byte, []byte) ([]byte, error)

	// determine algorithm and signing function, fail for unsupported keys
	if k, ok := key.(*ecdsa.PrivateKey); ok {
		c.header.Algorithm = "ES256"
		sign = func(data []byte, sum []byte) ([]byte, error) {
			return ecdsa.SignASN1(rand.Reader, k, sum)
		}
	} else if k, ok := key.(*rsa.PrivateKey); ok {
		c.header.Algorithm = "RS256"
		sign = func(data []byte, sum []byte) ([]byte, error) {
			return rsa.SignPKCS1v15(rand.Reader, k, hash, sum)
		}
	} else {
		return nil, fmt.Errorf("unrecognized/unsupported key type: %T", key)
	}

	// encode the header
	hs, err := c.header.encode()
	if err != nil {
		return nil, err
	}

	// encode the claims
	cs, err := c.claims.encode()
	if err != nil {
		return nil, err
	}

	// sign the token
	ss := fmt.Sprintf("%s.%s", hs, cs)
	h := hash.New()
	h.Write([]byte(ss))
	sig, err := sign([]byte(ss), h.Sum(nil))
	if err != nil {
		return nil, err
	}

	ret := fmt.Sprintf("%s.%s", ss, base64.RawURLEncoding.EncodeToString(sig))
	return &ret, nil
}

var _ Authorizer = &ClientAssertionAuthorizer{}

type ClientAssertionAuthorizer struct {
	conf *clientCredentialsConfig
}

func (a *ClientAssertionAuthorizer) assertion(tokenUrl string) (*string, error) {
	if a.conf == nil {
		return nil, fmt.Errorf("internal-error: ClientAssertionAuthorizer not configured")
	}

	if a.conf.Certificate == nil {
		return nil, fmt.Errorf("internal-error: ClientAssertionAuthorizer misconfigured; Certificate was nil")
	}

	keySig := sha1.Sum(a.conf.Certificate.Raw)
	keyId := base64.URLEncoding.EncodeToString(keySig[:])

	audience := a.conf.Audience
	if audience == "" {
		audience = tokenUrl
	}

	t := clientAssertionToken{
		header: clientAssertionTokenHeader{
			Type:  "JWT",
			KeyId: keyId,
		},
		claims: clientAssertionTokenClaims{
			Audience: audience,
			Issuer:   a.conf.ClientID,
			Subject:  a.conf.ClientID,
		},
	}

	assertion, err := t.encode(a.conf.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("ClientAssertionAuthorizer: failed to encode and sign JWT assertion: %v", err)
	}

	return assertion, nil
}

func (a *ClientAssertionAuthorizer) token(ctx context.Context, tokenUrl string) (*oauth2.Token, error) {
	if a.conf == nil {
		return nil, fmt.Errorf("internal-error: ClientAssertionAuthorizer not configured")
	}

	assertion := a.conf.FederatedAssertion
	if assertion == "" {
		a, err := a.assertion(tokenUrl)
		if err != nil {
			return nil, err
		}
		if a == nil {
			return nil, fmt.Errorf("ClientAssertionAuthorizer: assertion was nil")
		}
		assertion = *a
	}

	v := url.Values{
		"client_assertion":      {assertion},
		"client_assertion_type": {"urn:ietf:params:oauth:client-assertion-type:jwt-bearer"},
		"client_id":             {a.conf.ClientID},
		"grant_type":            {"client_credentials"},
		// NOTE: we intentionally only support v2 (MSAL) Tokens at this time since v1 (ADAL) is EOL
		"scope": []string{
			strings.Join(a.conf.Scopes, " "),
		},
	}

	return clientCredentialsToken(ctx, tokenUrl, &v)
}

func (a *ClientAssertionAuthorizer) Token(ctx context.Context, _ *http.Request) (*oauth2.Token, error) {
	if a.conf == nil {
		return nil, fmt.Errorf("could not request token: conf is nil")
	}

	tokenUrl := a.conf.TokenURL
	if tokenUrl == "" {
		if a.conf.Environment.Authorization == nil {
			return nil, fmt.Errorf("no `authorization` configuration was found for this environment")
		}
		tokenUrl = tokenEndpoint(*a.conf.Environment.Authorization, a.conf.TenantID)
	}

	return a.token(ctx, tokenUrl)
}

// AuxiliaryTokens returns additional tokens for auxiliary tenant IDs, for use in multi-tenant scenarios
func (a *ClientAssertionAuthorizer) AuxiliaryTokens(ctx context.Context, _ *http.Request) ([]*oauth2.Token, error) {
	if a.conf == nil {
		return nil, fmt.Errorf("could not request token: conf is nil")
	}

	tokens := make([]*oauth2.Token, 0)

	if len(a.conf.AuxiliaryTenantIDs) == 0 {
		return tokens, nil
	}

	for _, tenantId := range a.conf.AuxiliaryTenantIDs {
		tokenUrl := a.conf.TokenURL
		if tokenUrl == "" {
			if a.conf.Environment.Authorization == nil {
				return nil, fmt.Errorf("no `authorization` configuration was found for this environment")
			}
			tokenUrl = tokenEndpoint(*a.conf.Environment.Authorization, tenantId)
		}

		token, err := a.token(ctx, tokenUrl)
		if err != nil {
			return tokens, err
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

func clientCredentialsToken(ctx context.Context, endpoint string, params *url.Values) (*oauth2.Token, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer([]byte(params.Encode())))
	if err != nil {
		return nil, fmt.Errorf("clientCredentialsToken: failed to build request: %+v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("clientCredentialsToken: cannot request token: %v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("clientCredentialsToken: cannot parse response: %v", err)
	}

	if c := resp.StatusCode; c < 200 || c > 299 {
		return nil, fmt.Errorf("clientCredentialsToken: received HTTP status %d with response: %s", resp.StatusCode, body)
	}

	// clientCredentialsToken response can arrive with numeric values as integers or strings :(
	var tokenRes struct {
		AccessToken string      `json:"access_token"`
		TokenType   string      `json:"token_type"`
		IDToken     string      `json:"id_token"`
		Resource    string      `json:"resource"`
		Scope       string      `json:"scope"`
		ExpiresIn   interface{} `json:"expires_in"` // relative seconds from now
		ExpiresOn   interface{} `json:"expires_on"` // timestamp
	}
	if err := json.Unmarshal(body, &tokenRes); err != nil {
		return nil, fmt.Errorf("clientCredentialsToken: cannot unmarshal response: %v", err)
	}

	token := &oauth2.Token{
		AccessToken: tokenRes.AccessToken,
		TokenType:   tokenRes.TokenType,
	}

	var secs time.Duration
	if exp, ok := tokenRes.ExpiresIn.(string); ok && exp != "" {
		if v, err := strconv.Atoi(exp); err == nil {
			secs = time.Duration(v)
		}
	} else if exp, ok := tokenRes.ExpiresIn.(int64); ok {
		secs = time.Duration(exp)
	} else if exp, ok := tokenRes.ExpiresIn.(float64); ok {
		secs = time.Duration(exp)
	}
	if secs > 0 {
		token.Expiry = time.Now().Add(secs * time.Second)
	}

	return token, nil
}

func tokenEndpoint(endpoint environments.Authorization, tenant string) string {
	if tenant == "" {
		tenant = "common"
	}
	return fmt.Sprintf("%s/%s/oauth2/v2.0/token", endpoint.LoginEndpoint, tenant)
}
