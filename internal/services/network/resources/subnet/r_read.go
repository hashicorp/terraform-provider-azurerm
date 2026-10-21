// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subnet

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubnetRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.Subnets

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSubnetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if err := resourceSubnetFlatten(d, *id, resp.Model); err != nil {
		return fmt.Errorf("encoding %s: %+v", id, err)
	}

	return nil
}

var subnetDelegationServiceNames = []string{
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

// subnetServiceEndpointDiffSuppress is a TEMPORARY solution to suppress diffs for a service endpoint list where the
// individual entries have no change, and only the order of endpoints is different.
// The Azure API will not persist an order change, which results in an infinite diff loop if order was the only change.
// This function should be removed when the SDK v2 supports this feature for TypeList natively.
func subnetServiceEndpointDiffSuppress(_, _, _ string, d *schema.ResourceData) bool {
	if d == nil {
		return false
	}

	listKey := "service_endpoint"

	oldRaw, newRaw := d.GetChange(listKey)
	if oldRaw == nil || newRaw == nil {
		return false
	}

	oldList, okOld := oldRaw.([]any)
	newList, okNew := newRaw.([]any)
	if !okOld || !okNew {
		return false
	}

	if len(oldList) != len(newList) {
		return false
	}

	type endpointItem struct {
		service           string
		networkIdentifier string
	}

	canonicalize := func(rawList []any) []endpointItem {
		items := make([]endpointItem, 0, len(rawList))
		for _, raw := range rawList {
			m, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			item := endpointItem{}
			if s, ok := m["service"].(string); ok {
				item.service = s
			}
			if ni, ok := m["network_identifier"].(string); ok {
				item.networkIdentifier = ni
			}
			items = append(items, item)
		}

		sort.Slice(items, func(i, j int) bool {
			if items[i].service != items[j].service {
				return items[i].service < items[j].service
			}
			return items[i].networkIdentifier < items[j].networkIdentifier
		})

		return items
	}

	canonicalOld := canonicalize(oldList)
	canonicalNew := canonicalize(newList)

	return reflect.DeepEqual(canonicalOld, canonicalNew)
}
