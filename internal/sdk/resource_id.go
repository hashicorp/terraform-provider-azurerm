// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import "github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"

// SetID uses the specified ID Formatter to set the Resource ID
func (rmd ResourceMetaData) SetID(formatter resourceids.Id) {
	rmd.ResourceData.SetId(formatter.ID())
}
