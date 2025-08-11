// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
)

// Implementation handler for an UpgradeIdentity operation.
//
// This is used to encapsulate all upgrade logic from a prior identity to the
// current version when a Resource implements the
// ResourceWithUpgradeIdentity interface.
type IdentityUpgrader struct {
	// Schema information for the prior identity version. While not required,
	// setting this will populate the UpgradeIdentityRequest type Identity
	// field similar to other Resource data types. This allows for easier data
	// handling such as calling Get() or GetAttribute().
	//
	// If not set, prior identity data is available in the
	// UpgradeIdentityRequest type RawIdentity field.
	PriorSchema *identityschema.Schema

	// Provider defined logic for upgrading a resource identity from the prior
	// identity version to the current schema version.
	//
	// The context.Context parameter contains framework-defined loggers and
	// supports request cancellation.
	//
	// The UpgradeIdentityRequest parameter contains the prior identity data.
	// If PriorSchema was set, the Identity field will be available. Otherwise,
	// the RawIdentity must be used.
	//
	// The UpgradeIdentityResponse parameter should contain the upgraded
	// identity data and can be used to signal any logic warnings or errors.
	IdentityUpgrader func(context.Context, UpgradeIdentityRequest, *UpgradeIdentityResponse)
}
