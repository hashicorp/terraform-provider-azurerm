// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package service_fabric_managed_cluster

import (
	"context"
	"fmt"
	"maps"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2024-04-01/managedcluster"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2024-04-01/nodetype"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (k Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clusterClient := metadata.Client.ServiceFabricManaged.ManagedClusterClient
			nodeTypeClient := metadata.Client.ServiceFabricManaged.NodeTypeClient

			id, err := managedcluster.ParseManagedClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			cluster, err := clusterClient.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(cluster.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("while reading data for cluster %q: %+v", id.ManagedClusterName, err)
			}

			clusterId := nodetype.NewManagedClusterID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName)
			nts, err := nodeTypeClient.ListByManagedClustersComplete(ctx, clusterId)
			if err != nil {
				return fmt.Errorf("while listing NodeTypes for cluster %q: +%v", id.ManagedClusterName, err)
			}

			model := flattenClusterProperties(cluster.Model)
			// Password is read-only
			model.Password = metadata.ResourceData.Get("password").(string)
			model.ResourceGroup = id.ResourceGroupName
			model.NodeTypes = make([]NodeType, 0)
			for _, nt := range nts.Items {
				provState := nt.Properties.ProvisioningState
				if provState == nil || *provState == nodetype.ManagedResourceProvisioningStateDeleted || *provState == nodetype.ManagedResourceProvisioningStateDeleting {
					continue
				}
				model.NodeTypes = append(model.NodeTypes, flattenNodetypeProperties(nt))
			}
			return metadata.Encode(model)
		},
		Timeout: 5 * time.Minute,
	}
}

func flattenClusterProperties(cluster *managedcluster.ManagedCluster) *ClusterResourceModel {
	model := &ClusterResourceModel{}
	if cluster == nil {
		return model
	}

	model.Name = pointer.From(cluster.Name)
	model.Location = cluster.Location
	model.Sku = cluster.Sku.Name

	properties := cluster.Properties
	if properties == nil {
		return model
	}

	model.DNSName = properties.DnsName
	model.SubnetId = pointer.From(properties.SubnetId)

	if features := properties.AddonFeatures; features != nil {
		for _, feature := range *features {
			switch feature {
			case managedcluster.ManagedClusterAddOnFeatureDnsService:
				model.DNSService = true
			case managedcluster.ManagedClusterAddOnFeatureBackupRestoreService:
				model.BackupRestoreService = true
			}
		}
	}
	model.Username = properties.AdminUserName

	if aad := properties.AzureActiveDirectory; aad != nil {
		model.Authentication = append(model.Authentication, Authentication{})
		adModels := make([]ADAuthentication, 0)

		adModel := ADAuthentication{}
		adModel.ClientApp = pointer.From(aad.ClientApplication)
		adModel.ClusterApp = pointer.From(aad.ClusterApplication)
		adModel.TenantId = pointer.From(aad.TenantId)

		adModels = append(adModels, adModel)
		model.Authentication[0].ADAuth = adModels
	}

	if clients := properties.Clients; clients != nil {
		certs := make([]ThumbprintAuth, len(*clients))
		for idx, client := range *clients {
			t := CertTypeReadOnly
			if client.IsAdmin {
				t = CertTypeAdmin
			}
			certs[idx] = ThumbprintAuth{
				CertificateType: t,
				CommonName:      pointer.From(client.CommonName),
				Thumbprint:      pointer.From(client.Thumbprint),
			}
		}
		if len(model.Authentication) == 0 {
			model.Authentication = append(model.Authentication, Authentication{})
		}
		model.Authentication[0].CertAuthentication = certs
	}

	if fss := properties.FabricSettings; fss != nil {
		cfs := make([]CustomFabricSetting, 0)
		for _, fs := range *fss {
			for _, param := range fs.Parameters {
				cfs = append(cfs, CustomFabricSetting{
					Section:   fs.Name,
					Parameter: param.Name,
					Value:     param.Value,
				})
			}
		}
		model.CustomFabricSettings = cfs
	}

	model.ClientConnectionPort = pointer.From(properties.ClientConnectionPort)
	model.HTTPGatewayPort = pointer.From(properties.HTTPGatewayConnectionPort)

	if lbrules := properties.LoadBalancingRules; lbrules != nil {
		model.LBRules = make([]LBRule, len(*lbrules))
		for idx, rule := range *lbrules {
			model.LBRules[idx] = LBRule{
				BackendPort:      rule.BackendPort,
				FrontendPort:     rule.FrontendPort,
				ProbeProtocol:    rule.ProbeProtocol,
				ProbeRequestPath: pointer.From(rule.ProbeRequestPath),
				Protocol:         rule.Protocol,
			}
		}
	}

	if upgradeWave := properties.ClusterUpgradeCadence; upgradeWave != nil {
		model.UpgradeWave = *upgradeWave
	}

	if t := cluster.Tags; t != nil {
		modelTags := make(map[string]any)
		for tag, value := range *t {
			// This tag is temporary and will be removed at a later date.
			// More info can be found here https://azure.microsoft.com/en-us/updates/default-outbound-access-for-vms-in-azure-will-be-retired-transition-to-a-new-method-of-internet-access/
			// In the meantime, we'll ignore it when setting tags into state
			if !strings.Contains(tag, "SFRP.DisableDefaultOutboundAccess") {
				modelTags[tag] = value
			}
		}
		model.Tags = modelTags
	}
	return model
}

func flattenNodetypeProperties(nt nodetype.NodeType) NodeType {
	props := nt.Properties
	if props == nil {
		return NodeType{Name: pointer.From(nt.Name)}
	}

	out := NodeType{
		DataDiskSize:     pointer.From(nt.Properties.DataDiskSizeGB),
		Name:             pointer.From(nt.Name),
		Primary:          props.IsPrimary,
		VmImageOffer:     pointer.From(props.VMImageOffer),
		VmImagePublisher: pointer.From(props.VMImagePublisher),
		VmImageSku:       pointer.From(props.VMImageSku),
		VmImageVersion:   pointer.From(props.VMImageVersion),
		VmInstanceCount:  props.VMInstanceCount,
		VmSize:           pointer.From(props.VMSize),
		Id:               pointer.From(nt.Id),
	}

	if appPorts := props.ApplicationPorts; appPorts != nil {
		out.ApplicationPorts = fmt.Sprintf("%d-%d", appPorts.StartPort, appPorts.EndPort)
	}

	if ephemeralPorts := props.EphemeralPorts; ephemeralPorts != nil {
		out.EphemeralPorts = fmt.Sprintf("%d-%d", ephemeralPorts.StartPort, ephemeralPorts.EndPort)
	}

	if mpg := props.MultiplePlacementGroups; mpg != nil {
		out.MultiplePlacementGroupsEnabled = *mpg
	}

	if stateless := props.IsStateless; stateless != nil {
		out.Stateless = *stateless
	}

	if capacities := props.Capacities; capacities != nil {
		caps := make(map[string]string)
		maps.Copy(caps, *capacities)
		out.Capacities = caps
	}

	if diskType := props.DataDiskType; diskType != nil {
		out.DataDiskType = *diskType
	}

	if placementProps := props.PlacementProperties; placementProps != nil {
		placements := make(map[string]string)
		maps.Copy(placements, *placementProps)
		out.PlacementProperties = placements
	}

	if secrets := props.VMSecrets; secrets != nil {
		secs := make([]VmSecrets, len(*secrets))
		for idx, sec := range *secrets {
			certs := make([]VaultCertificates, len(sec.VaultCertificates))
			for idx, cert := range sec.VaultCertificates {
				certs[idx] = VaultCertificates{
					Store: cert.CertificateStore,
					Url:   cert.CertificateURL,
				}
			}
			secs[idx] = VmSecrets{
				SourceVault:  pointer.From(sec.SourceVault.Id),
				Certificates: certs,
			}
		}
		out.VmSecrets = secs
	}
	return out
}

type CustomFabricSetting struct {
	Parameter string `tfschema:"parameter"`
	Section   string `tfschema:"section"`
	Value     string `tfschema:"value"`
}

type LBRule struct {
	BackendPort      int64                        `tfschema:"backend_port"`
	FrontendPort     int64                        `tfschema:"frontend_port"`
	ProbeProtocol    managedcluster.ProbeProtocol `tfschema:"probe_protocol"`
	ProbeRequestPath string                       `tfschema:"probe_request_path"`
	Protocol         managedcluster.Protocol      `tfschema:"protocol"`
}

type ThumbprintAuth struct {
	CertificateType CertType `tfschema:"type"`
	CommonName      string   `tfschema:"common_name"`
	Thumbprint      string   `tfschema:"thumbprint"`
}

type ADAuthentication struct {
	ClientApp  string `tfschema:"client_application_id"`
	ClusterApp string `tfschema:"cluster_application_id"`
	TenantId   string `tfschema:"tenant_id"`
}

type Authentication struct {
	ADAuth             []ADAuthentication `tfschema:"active_directory"`
	CertAuthentication []ThumbprintAuth   `tfschema:"certificate"`
}

type VaultCertificates struct {
	Store string `tfschema:"store"`
	Url   string `tfschema:"url"`
}

type VmSecrets struct {
	SourceVault  string              `tfschema:"vault_id"`
	Certificates []VaultCertificates `tfschema:"certificates"`
}

type NodeType struct {
	DataDiskSize                   int64  `tfschema:"data_disk_size_gb"`
	Id                             string `tfschema:"id"`
	MultiplePlacementGroupsEnabled bool   `tfschema:"multiple_placement_groups_enabled"`
	Name                           string `tfschema:"name"`
	Primary                        bool   `tfschema:"primary"`
	Stateless                      bool   `tfschema:"stateless"`
	VmImageOffer                   string `tfschema:"vm_image_offer"`
	VmImagePublisher               string `tfschema:"vm_image_publisher"`
	VmImageSku                     string `tfschema:"vm_image_sku"`
	VmImageVersion                 string `tfschema:"vm_image_version"`
	VmInstanceCount                int64  `tfschema:"vm_instance_count"`
	VmSize                         string `tfschema:"vm_size"`

	ApplicationPorts    string            `tfschema:"application_port_range"`
	Capacities          map[string]string `tfschema:"capacities"`
	DataDiskType        nodetype.DiskType `tfschema:"data_disk_type"`
	EphemeralPorts      string            `tfschema:"ephemeral_port_range"`
	PlacementProperties map[string]string `tfschema:"placement_properties"`
	VmSecrets           []VmSecrets       `tfschema:"vm_secrets"`
}

type CertType string

const (
	CertTypeAdmin    CertType = "AdminClient"
	CertTypeReadOnly CertType = "ReadOnlyClient"
)

func parsePortRange(input string) (int64, int64, error) {
	if len(input) == 0 {
		return 0, 0, fmt.Errorf("port range is an empty string")
	}
	toks := strings.Split(input, "-")
	if len(toks) != 2 {
		return 0, 0, fmt.Errorf("invalid port range format in string %q", input)
	}
	from, err := strconv.ParseInt(toks[0], 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("while parsing %q as integer: %s", toks[0], err)
	}

	to, err := strconv.ParseInt(toks[1], 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("while parsing %q as integer: %s", toks[1], err)
	}
	return from, to, nil
}

func nodeTypeSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"data_disk_size_gb": {
					Type:     pluginsdk.TypeInt,
					Required: true,
				},
				"multiple_placement_groups_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
				"name": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"primary": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
				"stateless": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
				"vm_image_offer": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"vm_image_publisher": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"vm_image_sku": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"vm_image_version": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"vm_instance_count": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(3, 100),
				},
				"vm_size": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"application_port_range": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: func(i any, s string) ([]string, []error) {
						errors := make([]error, 0)
						if _, _, err := parsePortRange(i.(string)); err != nil {
							errors = append(errors, err)
						}
						return nil, errors
					},
				},
				"capacities": {
					Type:     pluginsdk.TypeMap,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
				"data_disk_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(nodetype.DiskTypeStandardLRS),
					ValidateFunc: validation.StringInSlice(nodetype.PossibleValuesForDiskType(), false),
				},
				"ephemeral_port_range": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"placement_properties": {
					Type:     pluginsdk.TypeMap,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
				"vm_secrets": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"vault_id": {
								Type:     pluginsdk.TypeString,
								Required: true,
							},
							"certificates": {
								Type:     pluginsdk.TypeList,
								Required: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"store": {
											Type:     pluginsdk.TypeString,
											Required: true,
										},
										"url": {
											Type:     pluginsdk.TypeString,
											Required: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func authSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"active_directory": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"client_application_id": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.IsUUID,
							},
							"cluster_application_id": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.IsUUID,
							},
							"tenant_id": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.IsUUID,
							},
						},
					},
				},
				"certificate": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"common_name": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotWhiteSpace,
							},
							"thumbprint": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotWhiteSpace,
							},
							"type": {
								Type:     pluginsdk.TypeString,
								Required: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(CertTypeAdmin),
									string(CertTypeReadOnly),
								}, false),
							},
						},
					},
				},
			},
		},
	}
}

func lbRulesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"backend_port": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
				},
				"frontend_port": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
				},
				"probe_protocol": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(managedcluster.PossibleValuesForProbeProtocol(), false),
				},
				"probe_request_path": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
				"protocol": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(managedcluster.PossibleValuesForProtocol(), false),
				},
			},
		},
	}
}
