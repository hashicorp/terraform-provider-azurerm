// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nginx_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	nginxClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/nginx/client"
)

func (m Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := nginxconfiguration.ParseConfigurationID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := nginxClient.From(meta.Client).NginxConfiguration

			if err := client.ConfigurationsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}

			return nil
		},
	}
}
