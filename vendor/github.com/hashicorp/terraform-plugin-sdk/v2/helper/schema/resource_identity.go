// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Implementation of a single identity schema version upgrade.
type IdentityUpgrader struct {
	// Version is the version schema that this Upgrader will handle, converting
	// it to Version+1.
	Version int64

	// Type describes the schema that this function can upgrade. Type is
	// required to decode the schema if the state was stored in a legacy
	// flatmap format.
	Type tftypes.Type

	// Upgrade takes the JSON encoded state and the provider meta value, and
	// upgrades the state one single schema version. The provided state is
	// decoded into the default json types using a map[string]interface{}. It
	// is up to the StateUpgradeFunc to ensure that the returned value can be
	// encoded using the new schema.
	Upgrade ResourceIdentityUpgradeFunc
}

type ResourceIdentity struct {
	// Version is the identity schema version.
	Version int64

	// SchemaFunc is the function that returns the schema for the
	// identity. Using a function for this field allows to prevent
	// storing all identity schema information in memory for the
	// lifecycle of a provider.
	// The types of the schema values are restricted to the types:
	//   - TypeBool
	//   - TypeFloat
	//   - TypeInt
	//   - TypeString
	//   - TypeList (of any of the above types)
	SchemaFunc func() map[string]*Schema

	// New struct, will be similar to (Resource).StateUpgraders
	IdentityUpgraders []IdentityUpgrader
}

// Function signature for an identity schema version upgrade handler.
//
// The Context parameter stores SDK information, such as loggers. It also
// is wired to receive any cancellation from Terraform such as a system or
// practitioner sending SIGINT (Ctrl-c).
//
// The map[string]interface{} parameter contains the previous identity schema
// version data for a managed resource instance. The keys are top level attribute
// names mapped to values that can be type asserted similar to
// fetching values using the ResourceData Get* methods:
//
//   - TypeBool: bool
//   - TypeFloat: float
//   - TypeInt: int
//   - TypeList: []interface{}
//   - TypeString: string
//
// In certain scenarios, the map may be nil, so checking for that condition
// upfront is recommended to prevent potential panics.
//
// The interface{} parameter is the result of the Provider type
// ConfigureFunc field execution. If the Provider does not define
// a ConfigureFunc, this will be nil. This parameter is conventionally
// used to store API clients and other provider instance specific data.
//
// The map[string]interface{} return parameter should contain the upgraded
// identity schema version data for a managed resource instance. Values must
// align to the typing mentioned above.
type ResourceIdentityUpgradeFunc func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error)

// SchemaMap returns the schema information for this resource identity
// defined via the SchemaFunc field.
func (ri *ResourceIdentity) SchemaMap() map[string]*Schema {
	if ri == nil || ri.SchemaFunc == nil {
		return nil
	}

	return ri.SchemaFunc()
}
