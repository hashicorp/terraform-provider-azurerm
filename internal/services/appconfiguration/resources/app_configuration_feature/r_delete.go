// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_configuration_feature

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/parse"
	"github.com/jackofallops/kermit/sdk/appconfiguration/1.0/appconfiguration"
)

func (k Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			nestedItemId, err := parse.ParseNestedItemID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClientWithEndpoint(nestedItemId.ConfigurationStoreEndpoint)
			if err != nil {
				return err
			}

			kv, err := client.GetKeyValues(ctx, nestedItemId.Key, nestedItemId.Label, "", "", []appconfiguration.KeyValueFields{})
			if err != nil {
				return fmt.Errorf("while checking for feature's %q existence: %+v", nestedItemId.Key, err)
			}
			keysFound := kv.Values()
			if len(keysFound) == 0 {
				return nil
			}

			if _, err = client.DeleteLock(ctx, nestedItemId.Key, nestedItemId.Label, "", ""); err != nil {
				return fmt.Errorf("while unlocking key %q: %+v", *nestedItemId, err)
			}

			if _, err = client.DeleteKeyValue(ctx, nestedItemId.Key, nestedItemId.Label, ""); err != nil {
				return fmt.Errorf("while removing key %q: %+v", *nestedItemId, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
