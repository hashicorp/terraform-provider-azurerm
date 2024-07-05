// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

func importDataConnectorTyped(expectKind securityinsight.DataConnectorKind) func(ctx context.Context, metadata sdk.ResourceMetaData) error {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		return importSentinelDataConnector(expectKind)(ctx, metadata.ResourceData, metadata.Client)
	}
}

func importDataConnectorUntyped(expectKind securityinsight.DataConnectorKind) *schema.ResourceImporter {
	return pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
		_, err := parse.DataConnectorID(id)
		return err
	}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
		wrapped := sdk.NewPluginSdkResourceData(d)
		if err := importSentinelDataConnector(expectKind)(ctx, wrapped, meta); err != nil {
			return nil, err
		}
		return []*pluginsdk.ResourceData{d}, nil
	})
}

func importSentinelDataConnector(expectKind securityinsight.DataConnectorKind) func(ctx context.Context, d sdk.ResourceData, meta interface{}) error {
	return func(ctx context.Context, d sdk.ResourceData, meta interface{}) error {
		id, err := parse.DataConnectorID(d.Id())
		if err != nil {
			return err
		}

		client := meta.(*clients.Client).Sentinel.DataConnectorsClient

		resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Alert Rule %q: %+v", id, err)
		}

		return assertDataConnectorKind(resp.Value, expectKind)
	}
}

func assertDataConnectorKind(dc securityinsight.BasicDataConnector, expectKind securityinsight.DataConnectorKind) error {
	var kind securityinsight.DataConnectorKind
	switch dc.(type) {
	case securityinsight.AADDataConnector:
		kind = securityinsight.DataConnectorKindAzureActiveDirectory
	case securityinsight.AATPDataConnector:
		kind = securityinsight.DataConnectorKindAzureAdvancedThreatProtection
	case securityinsight.ASCDataConnector:
		kind = securityinsight.DataConnectorKindAzureSecurityCenter
	case securityinsight.MCASDataConnector:
		kind = securityinsight.DataConnectorKindMicrosoftCloudAppSecurity
	case securityinsight.MTPDataConnector:
		kind = securityinsight.DataConnectorKindMicrosoftThreatProtection
	case securityinsight.IoTDataConnector:
		kind = securityinsight.DataConnectorKindIOT
	case securityinsight.Dynamics365DataConnector:
		kind = securityinsight.DataConnectorKindDynamics365
	case securityinsight.Office365ProjectDataConnector:
		kind = securityinsight.DataConnectorKindOffice365Project
	case securityinsight.OfficeIRMDataConnector:
		kind = securityinsight.DataConnectorKindOfficeIRM
	case securityinsight.OfficeDataConnector:
		kind = securityinsight.DataConnectorKindOffice365
	case securityinsight.OfficeATPDataConnector:
		kind = securityinsight.DataConnectorKindOfficeATP
	case securityinsight.OfficePowerBIDataConnector:
		kind = securityinsight.DataConnectorKindOfficePowerBI
	case securityinsight.AwsCloudTrailDataConnector:
		kind = securityinsight.DataConnectorKindAmazonWebServicesCloudTrail
	case securityinsight.MDATPDataConnector:
		kind = securityinsight.DataConnectorKindMicrosoftDefenderAdvancedThreatProtection
	case securityinsight.AwsS3DataConnector:
		kind = securityinsight.DataConnectorKindAmazonWebServicesS3
	case securityinsight.TiTaxiiDataConnector:
		kind = securityinsight.DataConnectorKindThreatIntelligenceTaxii
	case securityinsight.TIDataConnector:
		kind = securityinsight.DataConnectorKindThreatIntelligence
	}
	if expectKind != kind {
		return fmt.Errorf("Sentinel Data Connector has mismatched kind, expected: %q, got %q", expectKind, kind)
	}
	return nil
}
