package databoxedge

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databoxedge/validate"
)

func databoxEdgeCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, v interface{}) error {
	if value, ok := d.GetOk("shipment_address"); ok {
		shippingInfo := (value.([]interface{}))[0].(map[string]interface{})

		_, err := validate.DataboxEdgeStreetAddress(shippingInfo["address"].([]interface{}), "address")
		if err != nil {
			return fmt.Errorf("%s", err)
		}
	}

	return nil
}
