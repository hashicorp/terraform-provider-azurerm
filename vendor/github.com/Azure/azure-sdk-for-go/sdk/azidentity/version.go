// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

const (
	// UserAgent is the string to be used in the user agent string when making requests.
	UserAgent = "azidentity/" + Version

	// Version is the semantic version (see http://semver.org) of this module.
	Version = "v0.6.1"
)
