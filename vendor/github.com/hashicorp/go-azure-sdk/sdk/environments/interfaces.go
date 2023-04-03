// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package environments

type Api interface {
	AppId() (*string, bool)
	DomainSuffix() (*string, bool)
	Endpoint() (*string, bool)
	Name() string
	ResourceIdentifier() (*string, bool)
}
