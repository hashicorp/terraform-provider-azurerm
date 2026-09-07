// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importDataConnectorTyped(expectKind dataconnectors.DataConnectorKind) func(ctx context.Context, metadata sdk.ResourceMetaData) error {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		return importSentinelDataConnector(expectKind)(ctx, metadata.ResourceData, metadata.Client)
	}
}

func importDataConnectorUntyped(expectKind dataconnectors.DataConnectorKind) *schema.ResourceImporter {
	return pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
		_, err := dataconnectors.ParseDataConnectorID(id)
		return err
	}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
		wrapped := sdk.NewPluginSdkResourceData(d)
		if err := importSentinelDataConnector(expectKind)(ctx, wrapped, meta); err != nil {
			return nil, err
		}
		return []*pluginsdk.ResourceData{d}, nil
	})
}

func importSentinelDataConnector(expectKind dataconnectors.DataConnectorKind) func(ctx context.Context, d sdk.ResourceData, meta interface{}) error {
	return func(ctx context.Context, d sdk.ResourceData, meta interface{}) error {
		id, err := dataconnectors.ParseDataConnectorID(d.Id())
		if err != nil {
			return err
		}

		client := meta.(*clients.Client).Sentinel.DataConnectorsClient

		resp, err := client.Get(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		return assertDataConnectorKind(resp.Model, expectKind)
	}
}

func assertDataConnectorKind(dc dataconnectors.DataConnector, expectKind dataconnectors.DataConnectorKind) error {
	var kind dataconnectors.DataConnectorKind
	switch dc.(type) {
	case dataconnectors.AADDataConnector:
		kind = dataconnectors.DataConnectorKindAzureActiveDirectory
	case dataconnectors.AATPDataConnector:
		kind = dataconnectors.DataConnectorKindAzureAdvancedThreatProtection
	case dataconnectors.ASCDataConnector:
		kind = dataconnectors.DataConnectorKindAzureSecurityCenter
	case dataconnectors.MCASDataConnector:
		kind = dataconnectors.DataConnectorKindMicrosoftCloudAppSecurity
	case dataconnectors.MTPDataConnector:
		kind = dataconnectors.DataConnectorKindMicrosoftThreatProtection
	case dataconnectors.IoTDataConnector:
		kind = dataconnectors.DataConnectorKindIOT
	case dataconnectors.Dynamics365DataConnector:
		kind = dataconnectors.DataConnectorKindDynamicsThreeSixFive
	case dataconnectors.Office365ProjectDataConnector:
		kind = dataconnectors.DataConnectorKindOfficeThreeSixFiveProject
	case dataconnectors.OfficeIRMDataConnector:
		kind = dataconnectors.DataConnectorKindOfficeIRM
	case dataconnectors.OfficeDataConnector:
		kind = dataconnectors.DataConnectorKindOfficeThreeSixFive
	case dataconnectors.OfficeATPDataConnector:
		kind = dataconnectors.DataConnectorKindOfficeATP
	case dataconnectors.OfficePowerBIDataConnector:
		kind = dataconnectors.DataConnectorKindOfficePowerBI
	case dataconnectors.AwsCloudTrailDataConnector:
		kind = dataconnectors.DataConnectorKindAmazonWebServicesCloudTrail
	case dataconnectors.MDATPDataConnector:
		kind = dataconnectors.DataConnectorKindMicrosoftDefenderAdvancedThreatProtection
	case dataconnectors.AwsS3DataConnector:
		kind = dataconnectors.DataConnectorKindAmazonWebServicesSThree
	case dataconnectors.TiTaxiiDataConnector:
		kind = dataconnectors.DataConnectorKindThreatIntelligenceTaxii
	case dataconnectors.TIDataConnector:
		kind = dataconnectors.DataConnectorKindThreatIntelligence
	}
	if expectKind != kind {
		return fmt.Errorf("'Sentinel Data Connector' has mismatched kind, expected: %q, got %q", expectKind, kind)
	}
	return nil
}
