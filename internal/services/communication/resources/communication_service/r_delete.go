// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package communication_service

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/communicationservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.ServiceClient

			id, err := communicationservices.ParseCommunicationServiceID(metadata.ResourceData.Id())
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
