// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_zone

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.Zones

			id, err := zones.ParseDnsZoneID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, zones.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
