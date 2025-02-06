// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package environments

import (
	"strings"
)

type Api interface {
	// AppId is a GUID that identifies the application/API in the cloud environment
	AppId() (*string, bool)

	// Available returns whether the application/API is supported in the cloud environment
	Available() bool

	// DomainSuffix is the specific domain suffix for constructing endpoints for a data plane API in the cloud environment
	DomainSuffix() (*string, bool)

	// Endpoint is the common endpoint for the application/API in the cloud environment
	Endpoint() (*string, bool)

	// Name returns the name of the application/API
	Name() string

	// ResourceIdentifier is a URI that identifies the application/API in the cloud environment and
	// is used for constructing scopes/roles when authorizing connections.
	ResourceIdentifier() (*string, bool)

	// WithResourceIdentifier overrides the default resource ID for the API and is useful for APIs that offer multiple authorization scopes
	WithResourceIdentifier(string) Api
}

// ApiIsKnownPublished determines whether the provided Api represents the specified known API as published in PublishedApis
func ApiIsKnownPublished(api Api, apiName string) bool {
	appId, ok := api.AppId()
	if !ok || appId == nil {
		return false
	}
	knownApiAppId, ok := PublishedApis[apiName]
	if !ok {
		return false
	}
	if !strings.EqualFold(*appId, knownApiAppId) {
		return false
	}
	return true
}
