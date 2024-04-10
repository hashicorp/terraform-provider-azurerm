// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-06-02-preview/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-06-02-preview/managedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func validateKubernetesCluster(d *pluginsdk.ResourceData, cluster *managedclusters.ManagedCluster, resourceGroup, name string) error {
	if v, exists := d.GetOk("network_profile"); exists {
		rawProfiles := v.([]interface{})

		if len(rawProfiles) != 0 {
			// then ensure the conditionally-required fields are set
			profile := rawProfiles[0].(map[string]interface{})

			if networkPlugin := profile["network_plugin"].(string); networkPlugin != "" {
				dnsServiceIP := profile["dns_service_ip"].(string)
				serviceCidr := profile["service_cidr"].(string)
				podCidr := profile["pod_cidr"].(string)
				podCidrs := profile["pod_cidrs"].([]interface{})
				serviceCidrs := profile["service_cidrs"].([]interface{})
				networkPluginMode := profile["network_plugin_mode"].(string)
				isServiceCidrSet := serviceCidr != "" || len(serviceCidrs) != 0

				// Azure network plugin is not compatible with pod_cidr
				if podCidr != "" && strings.EqualFold(networkPlugin, "azure") && !strings.EqualFold(networkPluginMode, string(managedclusters.NetworkPluginModeOverlay)) {
					return fmt.Errorf("`pod_cidr` and `azure` cannot be set together unless specifying `network_plugin_mode` to `overlay`")
				}

				// if not All empty values or All set values.
				if !(dnsServiceIP == "" && !isServiceCidrSet) && !(dnsServiceIP != "" && isServiceCidrSet) {
					return fmt.Errorf("`dns_service_ip` and `service_cidr` should all be empty or both should be set")
				}

				ipVersions := profile["ip_versions"].([]interface{})
				if len(serviceCidrs) == 2 && len(ipVersions) != 2 {
					return fmt.Errorf("dual-stack networking must be enabled and `ip_versions` must be set to [\"IPv4\", \"IPv6\"] in order to specify multiple values in `service_cidrs`")
				}
				if len(podCidrs) == 2 && len(ipVersions) != 2 {
					return fmt.Errorf("dual-stack networking must be enabled and `ip_versions` must be set to [\"IPv4\", \"IPv6\"] in order to specify multiple values in `pod_cidrs`")
				}
			}
		}
	}

	// @tombuildsstuff: As of 2020-03-30 it's no longer possible to create a cluster using a Service Principal
	// for authentication (albeit this worked on 2020-03-27 via API version 2019-10-01 :shrug:). However it's
	// possible to rotate the Service Principal for an existing Cluster - so this needs to be supported via
	// update.
	//
	// For now we'll have to error out if attempting to create a new Cluster with an SP for auth - since otherwise
	// this gets silently converted to using MSI authentication.
	v, principalExists := d.GetOk("service_principal")
	if !principalExists {
		return nil
	}

	servicePrincipalsRaw, ok := v.([]interface{})
	if !ok || len(servicePrincipalsRaw) == 0 {
		// if it's an existing cluster, we need to check if there's currently a SP used on this cluster that isn't
		// defined locally, if so, we need to error out
		if cluster != nil {
			servicePrincipalExists := false
			if props := cluster.Properties; props != nil {
				if sp := props.ServicePrincipalProfile; sp != nil {
					if cid := sp.ClientId; cid != "" {
						// if it's MSI we ignore the block
						servicePrincipalExists = !strings.EqualFold(cid, "msi")
					}
				}
			}

			// a non-MI Service Principal exists on the cluster, but not locally
			if servicePrincipalExists {
				return existingClusterServicePrincipalRemovedErr
			}
		}

		return nil
	}

	// for a new cluster
	if cluster == nil {
		identityRaw, ok := d.GetOk("identity")
		if !ok {
			return nil
		}
		if vs := identityRaw.([]interface{}); len(vs) == 0 {
			return nil
		}

		// if we have both a Service Principal and an Identity Block defined
		return newClusterWithBothServicePrincipalAndMSIErr
	} else {
		// for an existing cluster
		servicePrincipalIsMsi := false
		if props := cluster.Properties; props != nil {
			if sp := props.ServicePrincipalProfile; sp != nil {
				if cid := sp.ClientId; cid != "" {
					servicePrincipalIsMsi = strings.EqualFold(cid, "msi")
				}
			}
		}

		// the user has a Service Principal block defined, but the Cluster's been upgraded to use MSI
		if servicePrincipalIsMsi {
			return existingClusterHasBeenUpgradedErr
		}

		hasIdentity := false
		if clusterIdentity := cluster.Identity; clusterIdentity != nil {
			hasIdentity = clusterIdentity.Type != identity.TypeNone
		}

		if hasIdentity {
			// there's a Service Principal block and an Identity block present - but it hasn't been upgraded
			// tell the user to update it
			return existingClusterRequiresUpgradingErr(resourceGroup, name)
		}
	}

	return nil
}

var existingClusterCommonErr = `
Azure Kubernetes Service has recently made several breaking changes to Cluster Authentication as
the Managed Identity Preview has concluded and entered General Availability.

During the Preview it was possible to create a Kubernetes Cluster using Mixed Mode Authentication,
meaning that you could specify both a Service Principal and a Managed Identity. Now that this has
entered General Availability, Managed Identity is used for all cluster related authentication.

This means that it's no longer possible to create a Mixed-Mode cluster - as such the 'identity'
and 'service_principal' block cannot be specified together.

Existing clusters using Mixed-Mode authentication will be updated to use only Managed Identity for
authentication when any change is made to the Cluster (but _not_ the Node Pool) - for example when
a tag is added or removed from the Cluster.
`

// existing cluster which has been switched to using MSI - user config needs updating
var existingClusterHasBeenUpgradedErr = fmt.Errorf(`
%s

It appears that this Kubernetes Cluster has been updated to use Managed Identity - as such it is
no longer possible to specify both a 'service_principal' and 'identity' block for this cluster.

To be able to continue managing this cluster in Terraform - please remove the 'service_principal'
block from the resource - which will match the changes made by Azure (where this cluster is now
only using Managed Identity for Cluster Authentication).
`, existingClusterCommonErr)

// existing cluster requires updating to using MI - then the user config needs updating
var existingClusterRequiresUpgradingErr = func(resourceGroup, name string) error {
	return fmt.Errorf(`
%s

This Kubernetes Cluster requires upgrading to use a Managed Identity rather than Mixed-Mode
authentication. Whilst Terraform could attempt to do this automatically, unfortunately this wouldn't
work for all users - and as such this needs to be performed out-of-band.

You can do this by making any change to the Cluster (not the Node Pool/Resource Group) - for example
adding a Tag, which can be done via the Azure CLI:

$ az resource update \
	--resource-type "Microsoft.ContainerService/managedClusters"
	--resource-group "%s"\
	--name "%s"\
	--set "tags.Foo='Bar'"

Which will set a tag of 'Foo' with the value 'Bar' on this Kubernetes Cluster.

Once the Cluster has finished updating, you can confirm Managed Identity is being used by running the
following Azure CLI command:

$ az resource show\
  --resource-type "Microsoft.ContainerService/managedClusters"\
  --q "properties.servicePrincipalProfile"\
  --resource-group "%s"\
  --name "%s"

.. which if successful should show:

{
  "clientId": "msi"
}

meaning that the cluster is using only a Managed Identity for Cluster Authentication.

---

Now that the Cluster has been updated - to continue using this Cluster in Terraform, remove the
'service_principal' block from your Terraform Configuration (since this is no longer required), at
which point this Cluster can be managed in Terraform as before.
`, existingClusterCommonErr, resourceGroup, name, resourceGroup, name)
}

// an existing cluster exists with an SP, but it's not defined in the config
var existingClusterServicePrincipalRemovedErr = fmt.Errorf(`
A Service Principal exists for this Kubernetes Cluster but has not been defined in the Terraform
Configuration.

At this time it's not possible to migrate from using a Service Principal for Cluster Authentication
to using a Managed Identity for Cluster Authentication - although support for this is being tracked
in this Github issue: https://github.com/Azure/AKS/issues/1520

To be able to continue managing this Kubernetes Cluster in Terraform, please re-introduce the
'service_principal' block. Alternatively you can re-create this Kubernetes Cluster by Tainting the
resource in Terraform, which will cause all Pods running on this Kubernetes Cluster to be recreated.
`)

// users trying to create a new cluster with an SP & MSI
var newClusterWithBothServicePrincipalAndMSIErr = fmt.Errorf(`
Azure Kubernetes Service has recently made several breaking changes to Cluster Authentication as
the Managed Identity Preview has concluded and entered General Availability.

During the Preview it was possible to create a Kubernetes Cluster using Mixed Mode Authentication,
meaning that you could specify both a Service Principal and a Managed Identity. Now that this has
entered General Availability, Managed Identity is used for all cluster related authentication.

This means that it's no longer possible to create a Mixed-Mode cluster - as such the 'identity'
and 'service_principal' block cannot be specified together. Instead you can either use a Service
Principal or Managed Identity for Cluster Authentication - but not both.

In order to create this Kubernetes Cluster, please remove either the 'identity' block or the
'service_principal' block.
`)

// returned when the Control Plane for the AKS Cluster must be upgraded in order to deploy this version to the Node Pool
var clusterControlPlaneMustBeUpgradedError = func(resourceGroup, clusterName, nodePoolName string, clusterVersion *string, desiredNodePoolVersion string, availableVersions []string) error {
	versions := make([]string, 0)
	for _, version := range availableVersions {
		versions = append(versions, fmt.Sprintf(" * %s", version))
	}
	versionsList := strings.Join(versions, "\n")
	clusterVersionDetails := "We were unable to determine the version of Kubernetes available on the Kubernetes Cluster."
	if clusterVersion != nil {
		clusterVersionDetails = fmt.Sprintf("The Kubernetes Cluster is running version %q.", *clusterVersion)
	}

	return fmt.Errorf(`
The Kubernetes/Orchestrator Version %q is not available for Node Pool %q.

Please confirm that this version is supported by the Kubernetes Cluster %q
(Resource Group %q) - which may need to be upgraded first.

%s

The supported Orchestrator Versions for this Node Pool/supported by this Kubernetes Cluster are:
%s

Node Pools cannot use a version of Kubernetes that is not supported on the Control Plane. More
details can be found at https://aka.ms/version-skew-policy.
`, desiredNodePoolVersion, nodePoolName, clusterName, resourceGroup, clusterVersionDetails, versionsList)
}

func validateNodePoolSupportsVersion(ctx context.Context, client *client.Client, currentNodePoolVersion string, defaultNodePoolId agentpools.AgentPoolId, desiredNodePoolVersion string) error {
	// confirm the version being used is >= the version of the control plane
	clusterId := commonids.NewKubernetesClusterID(defaultNodePoolId.SubscriptionId, defaultNodePoolId.ResourceGroupName, defaultNodePoolId.ManagedClusterName)
	resp, err := client.AgentPoolsClient.GetAvailableAgentPoolVersions(ctx, clusterId)
	if err != nil {
		return fmt.Errorf("retrieving Available Agent Pool Versions for %s: %+v", defaultNodePoolId, err)
	}
	versionExists := false
	supportedVersions := make([]string, 0)

	// when updating a cluster running a deprecated version of k8s then the validation should pass
	if currentNodePoolVersion == desiredNodePoolVersion {
		versionExists = true
	}

	// when creating a new cluster or upgrading the desired version should be supported
	if versions := resp.Model; !versionExists && versions != nil && versions.Properties.AgentPoolVersions != nil {
		for _, version := range *versions.Properties.AgentPoolVersions {
			if version.KubernetesVersion == nil {
				continue
			}

			v := *version.KubernetesVersion
			supportedVersions = append(supportedVersions, v)
			// alias versions (major.minor) are also fine as the latest supported GA patch version is chosen automatically in this case
			if v == desiredNodePoolVersion || v[:strings.LastIndex(v, ".")] == desiredNodePoolVersion {
				versionExists = true
			}
		}
	}

	if !versionExists {
		clusterId := commonids.NewKubernetesClusterID(defaultNodePoolId.SubscriptionId, defaultNodePoolId.ResourceGroupName, defaultNodePoolId.ManagedClusterName)
		cluster, err := client.KubernetesClustersClient.Get(ctx, clusterId)
		if err != nil {
			if !response.WasStatusCode(cluster.HttpResponse, http.StatusUnauthorized) {
				return fmt.Errorf("retrieving %s: %+v", clusterId, err)
			}
		}

		// nilable since a user may not necessarily have access, and this is trying to be helpful
		var clusterVersion *string
		if clusterModel := cluster.Model; clusterModel != nil && clusterModel.Properties != nil {
			clusterVersion = clusterModel.Properties.CurrentKubernetesVersion
		}

		return clusterControlPlaneMustBeUpgradedError(defaultNodePoolId.ResourceGroupName, defaultNodePoolId.ManagedClusterName, defaultNodePoolId.AgentPoolName, clusterVersion, desiredNodePoolVersion, supportedVersions)
	}

	return nil
}
