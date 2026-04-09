// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package auth

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

var _ CachingAuthorizer = &CachedAuthorizer{}

// CachedAuthorizer caches a token until it expires, then acquires a new token from Source
type CachedAuthorizer struct {
	// Source contains the underlying Authorizer for obtaining tokens
	Source Authorizer

	mutex     sync.RWMutex
	token     *oauth2.Token
	auxTokens []*oauth2.Token
}

// Token returns the current token if it's still valid, else will acquire a new token
func (c *CachedAuthorizer) Token(ctx context.Context, req *http.Request) (*oauth2.Token, error) {
	c.mutex.RLock()
	dueForRenewal := tokenDueForRenewal(c.token)
	c.mutex.RUnlock()

	if dueForRenewal {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		var err error
		c.token, err = c.Source.Token(ctx, req)
		if err != nil {
			return nil, err
		}
	}

	return c.token, nil
}

// AuxiliaryTokens returns additional tokens for auxiliary tenant IDs, for use in multi-tenant scenarios
func (c *CachedAuthorizer) AuxiliaryTokens(ctx context.Context, req *http.Request) ([]*oauth2.Token, error) {
	c.mutex.RLock()
	var dueForRenewal bool
	for _, token := range c.auxTokens {
		if dueForRenewal = tokenDueForRenewal(token); dueForRenewal {
			break
		}
	}
	c.mutex.RUnlock()

	if dueForRenewal || len(c.auxTokens) == 0 {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		var err error
		c.auxTokens, err = c.Source.AuxiliaryTokens(ctx, req)
		if err != nil {
			return nil, err
		}
	}

	return c.auxTokens, nil
}

// InvalidateCachedTokens expires the currently cached token and auxTokens, forcing new
// tokens to be acquired when Token() or AuxiliaryTokens() are next called
func (c *CachedAuthorizer) InvalidateCachedTokens() error {
	if c.token == nil {
		return nil
	}
	c.token.Expiry = time.Now()
	for i := range c.auxTokens {
		c.auxTokens[i].Expiry = time.Now()
	}
	return nil
}

// NewCachedAuthorizer returns an Authorizer that caches an access token for the duration of its validity.
// If the cached token expires, a new one is acquired and cached.
func NewCachedAuthorizer(src Authorizer) (CachingAuthorizer, error) {
	if _, ok := src.(*SharedKeyAuthorizer); ok {
		return nil, fmt.Errorf("internal-error: SharedKeyAuthorizer cannot be cached")
	}
	return &CachedAuthorizer{
		Source: src,
	}, nil
}
