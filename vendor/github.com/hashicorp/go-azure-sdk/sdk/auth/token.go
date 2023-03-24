// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"golang.org/x/oauth2"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/claims"
)

const tokenExpiryDelta = 20 * time.Minute

// tokenExpiresSoon returns true if the token expires within 10 minutes, or if more than 50% of its validity period has elapsed (if this can be determined), whichever is later
func tokenDueForRenewal(token *oauth2.Token) bool {
	if token == nil {
		return true
	}

	// Some tokens may never expire
	if token.Expiry.IsZero() {
		return false
	}

	expiry := token.Expiry.Round(0)
	delta := tokenExpiryDelta
	now := time.Now()
	expiresWithinTenMinutes := expiry.Add(-delta).Before(now)

	// Try to parse the token claims to retrieve the issuedAt time
	if claims, err := claims.ParseClaims(token); err == nil {
		if claims.IssuedAt > 0 {
			issued := time.Unix(claims.IssuedAt, 0)
			validity := expiry.Sub(issued)

			// If the validity period is less than double the expiry delta, then instead
			// determine whether >50% of the validity period has elapsed
			if validity < delta*2 {
				halfValidityHasElapsed := issued.Add(validity / 2).Before(now)
				return halfValidityHasElapsed
			}
		}
	}

	return expiresWithinTenMinutes
}
