// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armcore

const (
	// UserAgent is the string to be used in the user agent string when making requests.
	UserAgent = "armcore/" + Version

	// Version is the semantic version (see http://semver.org) of this module.
	Version = "v0.5.1"
)
