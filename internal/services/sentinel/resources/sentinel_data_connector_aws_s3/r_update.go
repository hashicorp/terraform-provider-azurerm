// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_aws_s3

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := dataconnectors.ParseDataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var plan DataConnectorAwsS3Model
			if err := metadata.Decode(&plan); err != nil {
				return err
			}

			client := metadata.Client.Sentinel.DataConnectorsClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			params, ok := resp.Model.(dataconnectors.AwsS3DataConnector)
			if !ok {
				return fmt.Errorf("%s was not an AWS S3 Data Connector", id)
			}

			if props := params.Properties; props != nil {
				if metadata.ResourceData.HasChange("aws_role_arn") {
					props.RoleArn = plan.AwsRoleArm
				}
				if metadata.ResourceData.HasChange("destination_table") {
					props.DestinationTable = plan.DestinationTable
				}
				if metadata.ResourceData.HasChange("sqs_urls") {
					props.SqsURLs = plan.SqsUrls
				}
			}

			if _, err := client.CreateOrUpdate(ctx, *id, params); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
