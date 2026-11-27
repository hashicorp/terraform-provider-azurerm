// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package load_test

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/loadtests"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	loadtestserviceClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadtestservice/client"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := loadtestserviceClient.From(metadata.Client).V20221201.LoadTests

			id, err := loadtests.ParseLoadTestID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
