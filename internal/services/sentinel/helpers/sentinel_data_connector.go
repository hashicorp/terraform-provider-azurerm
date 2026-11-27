// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	sentinelClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func ImportDataConnectorTyped(expectKind dataconnectors.DataConnectorKind) func(ctx context.Context, metadata sdk.ResourceMetaData) error {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		return ImportSentinelDataConnector(expectKind)(ctx, metadata.ResourceData, metadata.Client)
	}
}

func ImportDataConnectorUntyped(expectKind dataconnectors.DataConnectorKind) *schema.ResourceImporter {
	return pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
		_, err := dataconnectors.ParseDataConnectorID(id)
		return err
	}, func(ctx context.Context, d *pluginsdk.ResourceData, meta any) ([]*pluginsdk.ResourceData, error) {
		wrapped := sdk.NewPluginSdkResourceData(d)
		if err := ImportSentinelDataConnector(expectKind)(ctx, wrapped, meta); err != nil {
			return nil, err
		}
		return []*pluginsdk.ResourceData{d}, nil
	})
}

func ImportSentinelDataConnector(expectKind dataconnectors.DataConnectorKind) func(ctx context.Context, d sdk.ResourceData, meta any) error {
	return func(ctx context.Context, d sdk.ResourceData, meta any) error {
		id, err := dataconnectors.ParseDataConnectorID(d.Id())
		if err != nil {
			return err
		}

		client := sentinelClient.From(meta).DataConnectorsClient

		resp, err := client.Get(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		return AssertDataConnectorKind(resp.Model, expectKind)
	}
}

func AssertDataConnectorKind(dc dataconnectors.DataConnector, expectKind dataconnectors.DataConnectorKind) error {
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
