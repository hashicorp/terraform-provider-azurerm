package sdk

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"

// SetID uses the specified ID Formatter to set the Resource ID
func (rmd ResourceMetaData) SetID(formatter resourceid.Formatter) {
	rmd.ResourceData.SetId(formatter.ID())
}
