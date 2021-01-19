package sentinel

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2019-01-01-preview/securityinsight"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sentinel/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataConnectorID(dc securityinsight.BasicDataConnector) *string {
	if dc == nil {
		return nil
	}
	switch dc := dc.(type) {
	case securityinsight.AADDataConnector:
		return dc.ID
	case securityinsight.ASCDataConnector:
		return dc.ID
	case securityinsight.MCASDataConnector:
		return dc.ID
	case securityinsight.TIDataConnector:
		return dc.ID
	case securityinsight.TiTaxiiDataConnector:
		return dc.ID
	case securityinsight.OfficeDataConnector:
		return dc.ID
	case securityinsight.OfficeATPDataConnector:
		return dc.ID
	case securityinsight.AwsCloudTrailDataConnector:
		return dc.ID
	case securityinsight.AATPDataConnector:
		return dc.ID
	case securityinsight.MDATPDataConnector:
		return dc.ID
	default:
		return nil
	}
}

func importSentinelDataConnector(expectKind securityinsight.DataConnectorKind) func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
	return func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
		id, err := parse.DataConnectorID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).Sentinel.DataConnectorsClient
		ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
		defer cancel()

		resp, err := client.Get(ctx, id.ResourceGroup, operationalInsightsResourceProvider, id.WorkspaceName, id.Name)
		if err != nil {
			return nil, fmt.Errorf("retrieving Sentinel Alert Rule %q: %+v", id, err)
		}

		if err := assertDataConnectorKind(resp.Value, expectKind); err != nil {
			return nil, err
		}
		return []*schema.ResourceData{d}, nil
	}
}

func assertDataConnectorKind(dc securityinsight.BasicDataConnector, expectKind securityinsight.DataConnectorKind) error {
	var kind securityinsight.DataConnectorKind
	switch dc.(type) {
	case securityinsight.AADDataConnector:
		kind = securityinsight.DataConnectorKindAzureActiveDirectory
	case securityinsight.ASCDataConnector:
		kind = securityinsight.DataConnectorKindAzureSecurityCenter
	case securityinsight.MCASDataConnector:
		kind = securityinsight.DataConnectorKindMicrosoftCloudAppSecurity
	case securityinsight.TIDataConnector:
		kind = securityinsight.DataConnectorKindThreatIntelligence
	case securityinsight.TiTaxiiDataConnector:
		kind = securityinsight.DataConnectorKindThreatIntelligenceTaxii
	case securityinsight.OfficeDataConnector:
		kind = securityinsight.DataConnectorKindOffice365
	case securityinsight.OfficeATPDataConnector:
		kind = securityinsight.DataConnectorKindOfficeATP
	case securityinsight.AwsCloudTrailDataConnector:
		kind = securityinsight.DataConnectorKindAmazonWebServicesCloudTrail
	case securityinsight.AATPDataConnector:
		kind = securityinsight.DataConnectorKindAzureAdvancedThreatProtection
	case securityinsight.MDATPDataConnector:
		kind = securityinsight.DataConnectorKindMicrosoftDefenderAdvancedThreatProtection
	}
	if expectKind != kind {
		return fmt.Errorf("Sentinel Data Connector has mismatched kind, expected: %q, got %q", expectKind, kind)
	}
	return nil
}
