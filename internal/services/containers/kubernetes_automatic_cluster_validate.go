// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/managedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
)

func validateKubernetesAutomaticClusterTyped(model *KubernetesAutomaticClusterModel, cluster *managedclusters.ManagedCluster) error {
	identityExists := len(model.Identity) > 0

	if cluster == nil {
		if !identityExists {
			return fmt.Errorf("either an `identity` or `service_principal` block must be specified for cluster authentication")
		}
	} else {
		servicePrincipalExistsOnCluster := false
		if props := cluster.Properties; props != nil {
			if sp := props.ServicePrincipalProfile; sp != nil {
				if cid := sp.ClientId; cid != "" {
					// if it's MSI we ignore the block
					servicePrincipalExistsOnCluster = !strings.EqualFold(cid, "msi")
				}
			}
		}

		// a non-MI Service Principal exists on the cluster, but not locally
		if servicePrincipalExistsOnCluster {
			return fmt.Errorf("the Service Principal block cannot be removed once it has been set")
		}
		// Check if the user has a Service Principal block defined, but the Cluster's been upgraded to use MSI
	}
	return nil
}

// returned when the Control Plane for the AKS Cluster must be upgraded in order to deploy this version to the Node Pool
var errAutomaticClusterControlPlaneMustBeUpgraded = func(resourceGroup, clusterName, nodePoolName string, clusterVersion *string, desiredNodePoolVersion string, availableVersions []string) error {
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
details can be found at https://aka.ms/version-skew-policy`, desiredNodePoolVersion, nodePoolName, clusterName, resourceGroup, clusterVersionDetails, versionsList)
}

func validateNodePoolAutomaticSupportsVersion(ctx context.Context, client *client.Client, currentNodePoolVersion string, defaultNodePoolId agentpools.AgentPoolId, desiredNodePoolVersion string) error {
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

		return errAutomaticClusterControlPlaneMustBeUpgraded(defaultNodePoolId.ResourceGroupName, defaultNodePoolId.ManagedClusterName, defaultNodePoolId.AgentPoolName, clusterVersion, desiredNodePoolVersion, supportedVersions)
	}

	return nil
}
