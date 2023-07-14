// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginsdk

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
