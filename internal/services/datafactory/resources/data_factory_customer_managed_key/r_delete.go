// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_customer_managed_key

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			log.Printf(`[INFO] The Customer Managed Key cannot be removed from the Data Factory once added. To remove the Customer Managed Key delete and recreate the parent Data Factory.`)
			return nil
		},
	}
}
