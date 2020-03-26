package containers

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func validateKubernetesCluster(d *schema.ResourceData) error {
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

	v, principalExists := d.GetOk("service_principal")
	_, identityExists := d.GetOk("identity")

	if principalExists && identityExists {
		rawPrincipals := v.([]interface{})
		principal := rawPrincipals[0].(map[string]interface{})

		clientID := principal["client_id"].(string)

		if clientID != "msi" {
			return fmt.Errorf("`service_principal`(%q) and `identity` cannot both be set", clientID)
		}
	}

	if !principalExists && !identityExists {
		return fmt.Errorf("`service_principal` or `identity` must be set")
	}

	return nil
}
