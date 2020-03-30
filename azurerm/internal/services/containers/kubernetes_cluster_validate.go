package containers

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-11-01/containerservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func validateKubernetesCluster(d *schema.ResourceData, cluster *containerservice.ManagedClusterProperties) error {
	if v, exists := d.GetOk("network_profile"); exists {
		rawProfiles := v.([]interface{})

		if len(rawProfiles) != 0 {
			// then ensure the conditionally-required fields are set
			profile := rawProfiles[0].(map[string]interface{})

			if networkPlugin := profile["network_plugin"].(string); networkPlugin != "" {
				dockerBridgeCidr := profile["docker_bridge_cidr"].(string)
				dnsServiceIP := profile["dns_service_ip"].(string)
				serviceCidr := profile["service_cidr"].(string)
				podCidr := profile["pod_cidr"].(string)

				// Azure network plugin is not compatible with pod_cidr
				if podCidr != "" && networkPlugin == "azure" {
					return fmt.Errorf("`pod_cidr` and `azure` cannot be set together")
				}

				// if not All empty values or All set values.
				if !(dockerBridgeCidr == "" && dnsServiceIP == "" && serviceCidr == "") && !(dockerBridgeCidr != "" && dnsServiceIP != "" && serviceCidr != "") {
					return fmt.Errorf("`docker_bridge_cidr`, `dns_service_ip` and `service_cidr` should all be empty or all should be set")
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
		return nil
	}

	// create
	if cluster == nil {
		identityRaw, ok := d.GetOk("identity")
		if !ok {
			return nil
		}
		if vs := identityRaw.([]interface{}); len(vs) == 0 {
			return nil
		}

		// if we have both a Service Principal and an Identity Block defined
		return fmt.Errorf(`
Due to a breaking change in the AKS API - it's no longer possible to provision
new clusters using a Service Principal for cluster authentication - however
existing clusters can have their Service Principals updated.

Instead use the 'identity' block to create an AKS Cluster using a Managed Identity
for cluster authentication rather than a Service Principal.

Unfortunately at this time there's no migration path from using a Service Principal
for cluster authentication to using a Managed Identity for authentication without rebuilding
the Kubernetes Cluster.

The AKS feature request for this is: https://github.com/Azure/AKS/issues/1520
`)
	} else {
		servicePrincipalIsMsi := false
		if props := cluster.ServicePrincipalProfile; props != nil {
			if cid := props.ClientID; cid != nil {
				servicePrincipalIsMsi = strings.EqualFold(*cid, "msi")
			}
		}

		if servicePrincipalIsMsi {
			return fmt.Errorf(`
When the Kubernetes Cluster is using a Managed Identity for Cluster Authentication, it's not possible
to specify a Service Principal for the Cluster.

In addition, due to a breaking change in the AKS API - it's no longer possible to create a Kubernetes
Cluster using a Service Principal for Cluster Authentication - or to migrate between them. Instead you
must use a Managed Identity for Cluster Authentication for new Clusters (however existing clusters can
continue using Service Principals if they were created with them).
`)
		}
	}

	return nil
}
