// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package service_fabric_managed_cluster

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2024-04-01/managedcluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (k Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"backup_service_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"dns_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z0-9]+(-*[a-z0-9])*$`), "The dns name of the cluster must have lowercase letters, numbers and hyphens. The first character must be a letter and the last character a letter or number"),
		},
		"dns_service_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"location": commonschema.Location(),
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(4, 23),
				validation.StringMatch(regexp.MustCompile(`^[a-z0-9]+(-*[a-z0-9])*$`), "The name of the cluster must have lowercase letters, numbers and hyphens. The first character must be a letter and the last character a letter or number"),
			),
		},
		"username": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 15),
				validation.StringMatch(regexp.MustCompile("^[^\\\\/\"\\[\\]:|<>+=;,?*$]{1,14}$"), "User names cannot contain special characters \\/\"\"[]:|<>+=;,$?*@"),
			),
		},
		"password": {
			Type:      pluginsdk.TypeString,
			Optional:  true,
			Sensitive: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(8, 123),
				validation.StringIsNotWhiteSpace,
			),
		},
		"resource_group_name": commonschema.ResourceGroupName(),

		"node_type":      nodeTypeSchema(),
		"authentication": authSchema(),
		"custom_fabric_setting": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"parameter": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotWhiteSpace,
					},
					"section": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotWhiteSpace,
					},
					"value": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotWhiteSpace,
					},
				},
			},
		},
		"client_connection_port": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1500, 65535),
		},
		"http_gateway_port": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1500, 65535),
		},
		"lb_rule": lbRulesSchema(),
		"sku": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "Basic",
			ValidateFunc: validation.StringInSlice(managedcluster.PossibleValuesForSkuName(), false),
		},
		"subnet_id": commonschema.ResourceIDReferenceOptionalForceNew(&commonids.SubnetId{}, "Subnet"),
		"tags":      commonschema.Tags(),
		"upgrade_wave": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      managedcluster.ClusterUpgradeCadenceWaveZero,
			ValidateFunc: validation.StringInSlice(managedcluster.PossibleValuesForClusterUpgradeCadence(), false),
		},
	}
}

func (k Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (k Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff
			sku := rd.Get("sku").(string)
			var primary bool
			for _, nti := range rd.Get("node_type").([]any) {
				nt := nti.(map[string]any)
				vmCount := nt["vm_instance_count"].(int)
				if sku == string(managedcluster.SkuNameBasic) && vmCount < 3 {
					return fmt.Errorf("basic SKU requires at least 3 instances in a node type")
				} else if sku == string(managedcluster.SkuNameStandard) && vmCount < 5 {
					return fmt.Errorf("standard SKU requires at least 5 instances in a node type")
				}
				isPrimary := nt["primary"].(bool)
				if isPrimary && !primary {
					primary = true
				} else if isPrimary {
					return fmt.Errorf("only one node type can be primary at any given time. ")
				}
			}

			for _, lbi := range rd.Get("lb_rule").([]any) {
				lb := lbi.(map[string]any)
				probeProto := lb["probe_protocol"].(string)
				if probeProto == string(managedcluster.ProbeProtocolHTTP) || probeProto == string(managedcluster.ProbeProtocolHTTPS) {
					probePath := lb["probe_request_path"]
					if probePath == nil || probePath.(string) == "" {
						return fmt.Errorf("probe_request_path needs to be set if probe protocol is %q", probeProto)
					}
				}
			}

			o, n := rd.GetChange("node_type")
			oi := o.([]any)
			ni := n.([]any)
			if len(oi) > 0 && !reflect.DeepEqual(oi, ni) {
				for idx := range oi {
					oNodeType := oi[idx].(map[string]any)
					for nIdx := range ni {
						newNodeType := ni[nIdx].(map[string]any)
						if oNodeType["name"].(string) != newNodeType["name"].(string) {
							continue
						}
						for _, k := range []string{"name", "vm_size", "primary", "stateless"} {
							attr := fmt.Sprintf("node_type.%d.%s", idx, k)
							if rd.HasChange(attr) {
								return fmt.Errorf("node type attribute %q cannot be changed once node type is created", k)
							}
						}
					}
				}
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
