// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var SubnetResourceName = "azurerm_subnet"

var SubnetDelegationServiceNames = []string{
	"GitHub.Network/networkSettings",
	"Informatica.DataManagement/organizations",
	"Microsoft.ApiManagement/service",
	"Microsoft.Apollo/npu",
	"Microsoft.App/environments",
	"Microsoft.App/testClients",
	"Microsoft.AVS/PrivateClouds",
	"Microsoft.AzureCosmosDB/clusters",
	"Microsoft.BareMetal/AzureHostedService",
	"Microsoft.BareMetal/AzureHPC",
	"Microsoft.BareMetal/AzurePaymentHSM",
	"Microsoft.BareMetal/AzureVMware",
	"Microsoft.BareMetal/CrayServers",
	"Microsoft.BareMetal/MonitoringServers",
	"Microsoft.Batch/batchAccounts",
	"Microsoft.CloudTest/hostedpools",
	"Microsoft.CloudTest/images",
	"Microsoft.CloudTest/pools",
	"Microsoft.Codespaces/plans",
	"Microsoft.ContainerInstance/containerGroups",
	"Microsoft.ContainerService/managedClusters",
	"Microsoft.ContainerService/TestClients",
	"Microsoft.Databricks/workspaces",
	"Microsoft.DBforMySQL/flexibleServers",
	"Microsoft.DBforMySQL/servers",
	"Microsoft.DBforMySQL/serversv2",
	"Microsoft.DBforPostgreSQL/flexibleServers",
	"Microsoft.DBforPostgreSQL/serversv2",
	"Microsoft.DBforPostgreSQL/singleServers",
	"Microsoft.DelegatedNetwork/controller",
	"Microsoft.DevCenter/networkConnection",
	"Microsoft.DevOpsInfrastructure/pools",
	"Microsoft.DocumentDB/cassandraClusters",
	"Microsoft.Fidalgo/networkSettings",
	"Microsoft.HardwareSecurityModules/dedicatedHSMs",
	"Microsoft.Kusto/clusters",
	"Microsoft.LabServices/labplans",
	"Microsoft.Logic/integrationServiceEnvironments",
	"Microsoft.MachineLearningServices/workspaces",
	"Microsoft.MessagingConnectors/connectors",
	"Microsoft.Netapp/volumes",
	"Microsoft.Network/applicationGateways",
	"Microsoft.Network/dnsResolvers",
	"Microsoft.Network/managedResolvers",
	"Microsoft.Network/fpgaNetworkInterfaces",
	"Microsoft.Network/networkWatchers.",
	"Microsoft.Network/virtualNetworkGateways",
	"Microsoft.Orbital/orbitalGateways",
	"Microsoft.PowerAutomate/hostedRpa",
	"Microsoft.PowerPlatform/enterprisePolicies",
	"Microsoft.PowerPlatform/vnetaccesslinks",
	"Microsoft.ServiceFabricMesh/networks",
	"Microsoft.ServiceNetworking/trafficControllers",
	"Microsoft.Singularity/accounts/networks",
	"Microsoft.Singularity/accounts/npu",
	"Microsoft.Sql/managedInstances",
	"Microsoft.Sql/managedInstancesOnebox",
	"Microsoft.Sql/managedInstancesStage",
	"Microsoft.Sql/managedInstancesTest",
	"Microsoft.Sql/servers",
	"Microsoft.StoragePool/diskPools",
	"Microsoft.StreamAnalytics/streamingJobs",
	"Microsoft.Synapse/workspaces",
	"Microsoft.Web/hostingEnvironments",
	"Microsoft.Web/serverFarms",
	"NGINX.NGINXPLUS/nginxDeployments",
	"Oracle.Database/networkAttachments",
	"PaloAltoNetworks.Cloudngfw/firewalls",
	"Qumulo.Storage/fileSystems",
	"PureStorage.Block/storagePools",
}

func SubnetProvisioningStateRefreshFunc(ctx context.Context, client *subnets.SubnetsClient, id commonids.SubnetId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.Get(ctx, id, subnets.DefaultGetOperationOptions())
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id.String(), err)
		}

		if res.Model != nil && res.Model.Properties != nil && res.Model.Properties.ProvisioningState != nil {
			return res, string(*res.Model.Properties.ProvisioningState), nil
		}
		return nil, "", errors.New("unable to read provisioning state")
	}
}
