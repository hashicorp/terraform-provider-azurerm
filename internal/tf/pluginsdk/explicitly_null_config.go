// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package pluginsdk

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/ctyhelpers"
)

// IsExplicitlyNullInConfig determines whether the specified 'configFieldName'
// exists in the configuration file or not.
//
// Returns 'true' if the 'configFieldName' is not found in the configuration file or 'false' if
// the 'configFieldName' is found in the configuration file.
func IsExplicitlyNullInConfig(d *ResourceData, configFieldName string) bool {
	var isNull bool

	v := d.GetRawConfig().AsValueMap()[configFieldName]

	if v.IsNull() {
		isNull = true
	}

	return isNull
}

// rawConfigReader is a ResourceData or a ResourceDiff, both of which can read the configuration as written.
type rawConfigReader interface {
	GetRawConfigAt(cty.Path) (cty.Value, diag.Diagnostics)
}

// IsExplicitlyNullInConfigAt is IsExplicitlyNullInConfig for a property inside a block, given by its path, e.g. `node_type.0.is_primary`. A path that is
// not in the configuration counts as null.
func IsExplicitlyNullInConfigAt(d rawConfigReader, path string) bool {
	v, diags := d.GetRawConfigAt(ctyhelpers.ConstructCtyPath(path))

	return diags.HasError() || v.IsNull()
}
