// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

func importSentinelDataConnector(expectKind securityinsight.DataConnectorKind) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := parse.DataConnectorID(d.Id())
		if err != nil {
			return nil, err
		}

		client := azuresdkhacks.DataConnectorsClient{BaseClient: meta.(*clients.Client).Sentinel.DataConnectorsClient.BaseClient}
		resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
		if err != nil {
			return nil, fmt.Errorf("retrieving Sentinel Alert Rule %q: %+v", id, err)
		}

		if err := assertDataConnectorKind(resp.Value, expectKind); err != nil {
			return nil, err
		}
		return []*pluginsdk.ResourceData{d}, nil
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
	case securityinsight.TIDataConnector:
		kind = securityinsight.DataConnectorKindThreatIntelligence
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
	case azuresdkhacks.TiTaxiiDataConnector:
		kind = securityinsight.DataConnectorKindThreatIntelligenceTaxii
	case azuresdkhacks.TIDataConnector:
		kind = securityinsight.DataConnectorKindThreatIntelligence
	}
	if expectKind != kind {
		return fmt.Errorf("Sentinel Data Connector has mismatched kind, expected: %q, got %q", expectKind, kind)
	}
	return nil
}
