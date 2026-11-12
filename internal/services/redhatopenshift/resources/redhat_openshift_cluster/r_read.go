// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redhat_openshift_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2025-07-25/openshiftclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RedHatOpenShift.OpenShiftClustersClient

			id, err := openshiftclusters.ParseOpenShiftClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			var config RedHatOpenShiftClusterModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			return r.flatten(metadata, *id, resp.Model, config)
		},
	}
}

type ServicePrincipal struct {
	ClientId     string `tfschema:"client_id"`
	ClientSecret string `tfschema:"client_secret"`
}

type ClusterProfile struct {
	PullSecret               string `tfschema:"pull_secret"`
	Domain                   string `tfschema:"domain"`
	ManagedResourceGroupName string `tfschema:"managed_resource_group_name"`
	ResourceGroupId          string `tfschema:"resource_group_id"`
	Version                  string `tfschema:"version"`
	FipsEnabled              bool   `tfschema:"fips_enabled"`
}

type NetworkProfile struct {
	OutboundType                             string                `tfschema:"outbound_type"`
	PodCidr                                  string                `tfschema:"pod_cidr"`
	ServiceCidr                              string                `tfschema:"service_cidr"`
	PreconfiguredNetworkSecurityGroupEnabled bool                  `tfschema:"preconfigured_network_security_group_enabled"`
	LoadBalancerProfile                      []LoadBalancerProfile `tfschema:"load_balancer_profile"`
}

type LoadBalancerProfile struct {
	ManagedOutboundIpCount int64    `tfschema:"managed_outbound_ip_count"`
	EffectiveOutboundIps   []string `tfschema:"effective_outbound_ips"`
}

type PlatformWorkloadIdentityProfile struct {
	UpgradeableTo            string                     `tfschema:"upgradeable_to"`
	PlatformWorkloadIdentity []PlatformWorkloadIdentity `tfschema:"platform_workload_identity"`
}

type PlatformWorkloadIdentity struct {
	Name       string `tfschema:"name"`
	IdentityId string `tfschema:"identity_id"`
	ClientId   string `tfschema:"client_id"`
	ObjectId   string `tfschema:"object_id"`
}

type MainProfile struct {
	SubnetId                string `tfschema:"subnet_id"`
	VmSize                  string `tfschema:"vm_size"`
	DiskEncryptionSetId     string `tfschema:"disk_encryption_set_id"`
	EncryptionAtHostEnabled bool   `tfschema:"encryption_at_host_enabled"`
}

type WorkerProfile struct {
	VmSize                  string `tfschema:"vm_size"`
	SubnetId                string `tfschema:"subnet_id"`
	DiskEncryptionSetId     string `tfschema:"disk_encryption_set_id"`
	DiskSizeGb              int64  `tfschema:"disk_size_gb"`
	NodeCount               int64  `tfschema:"node_count"`
	EncryptionAtHostEnabled bool   `tfschema:"encryption_at_host_enabled"`
}

type IngressProfile struct {
	Visibility string `tfschema:"visibility"`
	IpAddress  string `tfschema:"ip_address"`
	Name       string `tfschema:"name"`
}

type ApiServerProfile struct {
	Visibility string `tfschema:"visibility"`
	IpAddress  string `tfschema:"ip_address"`
	Url        string `tfschema:"url"`
}
