// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_threat_intelligence_taxii

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient

			id, err := dataconnectors.ParseDataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var plan DataConnectorThreatIntelligenceTAXIIModel
			if err := metadata.Decode(&plan); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			dc, ok := existing.Model.(dataconnectors.TiTaxiiDataConnector)
			if !ok {
				return fmt.Errorf("%s was not an Threat Intelligence TAXII Data Connector", id)
			}

			if props := dc.Properties; props != nil {
				if metadata.ResourceData.HasChange("display_name") {
					props.FriendlyName = &plan.DisplayName
				}
				if metadata.ResourceData.HasChange("api_root_url") {
					props.TaxiiServer = &plan.APIRootURL
				}
				if metadata.ResourceData.HasChange("collection_id") {
					props.CollectionId = &plan.CollectionID
				}
				if metadata.ResourceData.HasChange("user_name") {
					props.UserName = &plan.UserName
				}
				if metadata.ResourceData.HasChange("password") {
					props.Password = &plan.Password
				}
				if metadata.ResourceData.HasChange("polling_frequency") {
					props.PollingFrequency = dataconnectors.PollingFrequency(plan.PollingFrequency)
				}
				if metadata.ResourceData.HasChange("lookback_date") {
					// Format is guaranteed by schema validation
					lookbackDate, _ := time.Parse(time.RFC3339, plan.LookbackDate)
					props.TaxiiLookbackPeriod = pointer.To(lookbackDate.Format(time.RFC3339))
				}

				// Setting the user name and password if non empty in plan, which are required by the API.
				if plan.UserName != "" {
					props.UserName = &plan.UserName
				}
				if plan.Password != "" {
					props.Password = &plan.Password
				}
			}

			if _, err := client.CreateOrUpdate(ctx, *id, dc); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
