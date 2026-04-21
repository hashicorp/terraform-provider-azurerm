// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/applicationgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	commonValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	containerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ACIConnectorLinuxModel struct {
	SubnetName        string                   `tfschema:"subnet_name"`
	ConnectorIdentity []ConnectorIdentityModel `tfschema:"connector_identity"`
}

type ConnectorIdentityModel struct {
	ClientID               string `tfschema:"client_id"`
	ObjectID               string `tfschema:"object_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

type ConfidentialComputingModel struct {
	SGXQuoteHelperEnabled bool `tfschema:"sgx_quote_helper_enabled"`
}

type IngressApplicationGatewayModel struct {
	GatewayID                         string                                   `tfschema:"gateway_id"`
	GatewayName                       string                                   `tfschema:"gateway_name"`
	SubnetCIDR                        string                                   `tfschema:"subnet_cidr"`
	SubnetID                          string                                   `tfschema:"subnet_id"`
	EffectiveGatewayID                string                                   `tfschema:"effective_gateway_id"`
	IngressApplicationGatewayIdentity []IngressApplicationGatewayIdentityModel `tfschema:"ingress_application_gateway_identity"`
}

type IngressApplicationGatewayIdentityModel struct {
	ClientID               string `tfschema:"client_id"`
	ObjectID               string `tfschema:"object_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

type KeyVaultSecretsProviderModel struct {
	// SecretRotationEnabled  bool                  `tfschema:"secret_rotation_enabled"`
	SecretRotationInterval string                `tfschema:"secret_rotation_interval"`
	SecretIdentity         []SecretIdentityModel `tfschema:"secret_identity"`
}

type SecretIdentityModel struct {
	ClientID               string `tfschema:"client_id"`
	ObjectID               string `tfschema:"object_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

type OMSAgentModel struct {
	LogAnalyticsWorkspaceID     string                  `tfschema:"log_analytics_workspace_id"`
	MSIAuthForMonitoringEnabled *bool                   `tfschema:"msi_auth_for_monitoring_enabled"`
	OMSAgentIdentity            []OMSAgentIdentityModel `tfschema:"oms_agent_identity"`
}

type OMSAgentIdentityModel struct {
	ClientID               string `tfschema:"client_id"`
	ObjectID               string `tfschema:"object_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

// Schema Definition for Automatic Cluster Addons

func schemaKubernetesAutomaticClusterAddOnsTyped() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"open_service_mesh_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
		"aci_connector_linux": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"connector_identity": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"object_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"user_assigned_identity_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
		// "azure_policy_enabled": {
		// 	Type:     pluginsdk.TypeBool,
		// 	Optional: true,
		//	Default:  true,
		// },
		"confidential_computing": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"sgx_quote_helper_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
		},
		"http_application_routing_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
		"oms_agent": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"log_analytics_workspace_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: workspaces.ValidateWorkspaceID,
					},
					"msi_auth_for_monitoring_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"oms_agent_identity": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"object_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"user_assigned_identity_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
		"ingress_application_gateway": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"gateway_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ConflictsWith: []string{
							"ingress_application_gateway.0.subnet_cidr",
							"ingress_application_gateway.0.subnet_id",
						},
						AtLeastOneOf: []string{
							"ingress_application_gateway.0.gateway_id",
							"ingress_application_gateway.0.subnet_cidr",
							"ingress_application_gateway.0.subnet_id",
						},
						ValidateFunc: applicationgateways.ValidateApplicationGatewayID,
					},
					"gateway_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"subnet_cidr": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ConflictsWith: []string{
							"ingress_application_gateway.0.gateway_id",
							"ingress_application_gateway.0.subnet_id",
						},
						AtLeastOneOf: []string{
							"ingress_application_gateway.0.gateway_id",
							"ingress_application_gateway.0.subnet_cidr",
							"ingress_application_gateway.0.subnet_id",
						},
						ValidateFunc: commonValidate.CIDR,
					},
					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ConflictsWith: []string{
							"ingress_application_gateway.0.gateway_id",
							"ingress_application_gateway.0.subnet_cidr",
						},
						AtLeastOneOf: []string{
							"ingress_application_gateway.0.gateway_id",
							"ingress_application_gateway.0.subnet_cidr",
							"ingress_application_gateway.0.subnet_id",
						},
						ValidateFunc: commonids.ValidateSubnetID,
					},
					"effective_gateway_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"ingress_application_gateway_identity": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"object_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"user_assigned_identity_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
		"key_vault_secrets_provider": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// "secret_rotation_enabled": {
					// 	Type:     pluginsdk.TypeBool,
					// 	Default:  false,
					// 	Optional: true,
					// 	AtLeastOneOf: []string{
					// 		"key_vault_secrets_provider.0.secret_rotation_enabled",
					// 		"key_vault_secrets_provider.0.secret_rotation_interval",
					// 	},
					// },
					"secret_rotation_interval": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "2m",
						// AtLeastOneOf: []string{
						// 	"key_vault_secrets_provider.0.secret_rotation_enabled",
						// 	"key_vault_secrets_provider.0.secret_rotation_interval",
						// },
						ValidateFunc: containerValidate.Duration,
					},
					"secret_identity": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"object_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"user_assigned_identity_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func expandKubernetesAddOnsTyped(input *KubernetesAutomaticClusterModel, env environments.Environment) (*map[string]managedclusters.ManagedClusterAddonProfile, error) {
	addonProfiles := map[string]managedclusters.ManagedClusterAddonProfile{}

	addonProfiles[openServiceMeshKey] = managedclusters.ManagedClusterAddonProfile{
		Enabled: input.OpenServiceMeshEnabled,
	}

	addonProfiles[confidentialComputingKey] = managedclusters.ManagedClusterAddonProfile{
		Enabled: false,
	}
	if len(input.ConfidentialComputing) > 0 {
		cc := input.ConfidentialComputing[0]
		config := make(map[string]string)
		quoteHelperEnabled := "false"
		if cc.SGXQuoteHelperEnabled {
			quoteHelperEnabled = "true"
		}
		config["ACCSGXQuoteHelperEnabled"] = quoteHelperEnabled
		addonProfiles[confidentialComputingKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: true,
			Config:  &config,
		}
	}

	if input.HTTPApplicationRoutingEnabled {
		addonProfiles[httpApplicationRoutingKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: input.HTTPApplicationRoutingEnabled,
		}
	}

	addonProfiles[omsAgentKey] = managedclusters.ManagedClusterAddonProfile{
		Enabled: false,
	}
	if len(input.OMSAgent) > 0 {
		oms := input.OMSAgent[0]
		config := make(map[string]string)

		if oms.LogAnalyticsWorkspaceID != "" {
			lawid, err := workspaces.ParseWorkspaceIDInsensitively(oms.LogAnalyticsWorkspaceID)
			if err != nil {
				return nil, fmt.Errorf("parsing Log Analytics Workspace ID: %+v", err)
			}
			config["logAnalyticsWorkspaceResourceID"] = lawid.ID()
		}

		if oms.MSIAuthForMonitoringEnabled != nil {
			config["useAADAuth"] = fmt.Sprintf("%t", *oms.MSIAuthForMonitoringEnabled)
		}

		addonProfiles[omsAgentKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: true,
			Config:  &config,
		}
	}

	addonProfiles[aciConnectorKey] = managedclusters.ManagedClusterAddonProfile{
		Enabled: false,
	}
	if len(input.ACIConnectorLinux) > 0 {
		aci := input.ACIConnectorLinux[0]
		config := make(map[string]string)

		if aci.SubnetName != "" {
			config["SubnetName"] = aci.SubnetName
		}

		addonProfiles[aciConnectorKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: true,
			Config:  &config,
		}
	}

	// addonProfiles[azurePolicyKey] = managedclusters.ManagedClusterAddonProfile{
	// 	Enabled: input.AzurePolicyEnabled,
	// 	Config: pointer.To(map[string]string{
	//		"version": "v2",
	//	}),
	//}

	addonProfiles[ingressApplicationGatewayKey] = managedclusters.ManagedClusterAddonProfile{
		Enabled: false,
	}
	if len(input.IngressApplicationGateway) > 0 {
		iag := input.IngressApplicationGateway[0]
		config := make(map[string]string)

		if iag.GatewayID != "" {
			config["applicationGatewayId"] = iag.GatewayID
		}

		if iag.GatewayName != "" {
			config["applicationGatewayName"] = iag.GatewayName
		}

		if iag.SubnetCIDR != "" {
			config["subnetCIDR"] = iag.SubnetCIDR
		}

		if iag.SubnetID != "" {
			config["subnetId"] = iag.SubnetID
		}

		addonProfiles[ingressApplicationGatewayKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: true,
			Config:  &config,
		}
	}

	if len(input.KeyVaultSecretsProvider) > 0 {
		kvsp := input.KeyVaultSecretsProvider[0]
		config := make(map[string]string)

		// enableSecretRotation := fmt.Sprintf("%t", kvsp.SecretRotationEnabled)
		config["enableSecretRotation"] = fmt.Sprintf("%t", true)
		config["rotationPollInterval"] = kvsp.SecretRotationInterval

		addonProfiles[azureKeyvaultSecretsProviderKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: true,
			Config:  &config,
		}
	}

	return filterUnsupportedKubernetesAddOns(addonProfiles, env)
}

func flattenKubernetesAddOnsTyped(profile map[string]managedclusters.ManagedClusterAddonProfile) (
	aciConnectorLinux []ACIConnectorLinuxModel,
	// azurePolicyEnabled bool,
	confidentialComputing []ConfidentialComputingModel,
	httpApplicationRoutingEnabled bool,
	httpApplicationRoutingZoneName string,
	ingressApplicationGateway []IngressApplicationGatewayModel,
	keyVaultSecretsProvider []KeyVaultSecretsProviderModel,
	omsAgent []OMSAgentModel,
	openServiceMeshEnabled bool,
) {
	aciConnector := kubernetesAddonProfileLocate(profile, aciConnectorKey)
	if aciConnector.Enabled {
		subnetName := ""
		if v := aciConnector.Config; v != nil && (*v)["SubnetName"] != "" {
			subnetName = (*v)["SubnetName"]
		}

		identity := flattenKubernetesClusterAddOnIdentityProfileTyped(aciConnector.Identity)

		aciConnectorLinux = []ACIConnectorLinuxModel{{
			SubnetName:        subnetName,
			ConnectorIdentity: identity,
		}}
	}

	// azurePolicy := kubernetesAddonProfileLocate(profile, azurePolicyKey)
	// azurePolicyEnabled = azurePolicy.Enabled

	confidentialComputingProfile := kubernetesAddonProfileLocate(profile, confidentialComputingKey)
	if confidentialComputingProfile.Enabled {
		quoteHelperEnabled := false
		if v := kubernetesAddonProfilelocateInConfig(confidentialComputingProfile.Config, "ACCSGXQuoteHelperEnabled"); v != "" && v != "false" {
			quoteHelperEnabled = true
		}
		confidentialComputing = []ConfidentialComputingModel{{
			SGXQuoteHelperEnabled: quoteHelperEnabled,
		}}
	}

	httpApplicationRouting := kubernetesAddonProfileLocate(profile, httpApplicationRoutingKey)
	httpApplicationRoutingEnabled = httpApplicationRouting.Enabled

	if v := kubernetesAddonProfilelocateInConfig(httpApplicationRouting.Config, "HTTPApplicationRoutingZoneName"); v != "" {
		httpApplicationRoutingZoneName = v
	}

	omsAgentProfile := kubernetesAddonProfileLocate(profile, omsAgentKey)
	if omsAgentProfile.Enabled {
		workspaceID := ""
		useAADAuth := false

		if v := kubernetesAddonProfilelocateInConfig(omsAgentProfile.Config, "logAnalyticsWorkspaceResourceID"); v != "" {
			if lawid, err := workspaces.ParseWorkspaceIDInsensitively(v); err == nil {
				workspaceID = lawid.ID()
			}
		}

		if v := kubernetesAddonProfilelocateInConfig(omsAgentProfile.Config, "useAADAuth"); v != "false" && v != "" {
			useAADAuth = true
		}

		omsAgentIdentity := flattenKubernetesClusterAddOnIdentityProfileTyped(omsAgentProfile.Identity)

		omsAgent = []OMSAgentModel{{
			LogAnalyticsWorkspaceID:     workspaceID,
			MSIAuthForMonitoringEnabled: pointer.To(useAADAuth),
			OMSAgentIdentity:            flattenOMSAgentIdentityTyped(omsAgentIdentity),
		}}
	}

	ingressApplicationGatewayProfile := kubernetesAddonProfileLocate(profile, ingressApplicationGatewayKey)
	if ingressApplicationGatewayProfile.Enabled {
		gatewayId := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGatewayProfile.Config, "applicationGatewayId"); v != "" {
			gatewayId = v
		}

		gatewayName := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGatewayProfile.Config, "applicationGatewayName"); v != "" {
			gatewayName = v
		}

		effectiveGatewayId := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGatewayProfile.Config, "effectiveApplicationGatewayId"); v != "" {
			effectiveGatewayId = v
		}

		subnetCIDR := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGatewayProfile.Config, "subnetCIDR"); v != "" {
			subnetCIDR = v
		}

		subnetId := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGatewayProfile.Config, "subnetId"); v != "" {
			subnetId = v
		}

		ingressApplicationGatewayIdentity := flattenKubernetesClusterAddOnIdentityProfileTyped(ingressApplicationGatewayProfile.Identity)

		ingressApplicationGateway = []IngressApplicationGatewayModel{{
			GatewayID:                         gatewayId,
			GatewayName:                       gatewayName,
			EffectiveGatewayID:                effectiveGatewayId,
			SubnetCIDR:                        subnetCIDR,
			SubnetID:                          subnetId,
			IngressApplicationGatewayIdentity: flattenIngressApplicationGatewayIdentityTyped(ingressApplicationGatewayIdentity),
		}}
	}

	openServiceMesh := kubernetesAddonProfileLocate(profile, openServiceMeshKey)
	openServiceMeshEnabled = openServiceMesh.Enabled

	azureKeyVaultSecretsProviderProfile := kubernetesAddonProfileLocate(profile, azureKeyvaultSecretsProviderKey)
	if azureKeyVaultSecretsProviderProfile.Enabled {
		// enableSecretRotation := false
		rotationPollInterval := ""

		// if v := kubernetesAddonProfilelocateInConfig(azureKeyVaultSecretsProviderProfile.Config, "enableSecretRotation"); v != "false" {
		// 	enableSecretRotation = true
		// }

		if v := kubernetesAddonProfilelocateInConfig(azureKeyVaultSecretsProviderProfile.Config, "rotationPollInterval"); v != "" {
			rotationPollInterval = v
		}

		azureKeyvaultSecretsProviderIdentity := flattenKubernetesClusterAddOnIdentityProfileTyped(azureKeyVaultSecretsProviderProfile.Identity)

		keyVaultSecretsProvider = []KeyVaultSecretsProviderModel{{
			// SecretRotationEnabled:  enableSecretRotation,
			SecretRotationInterval: rotationPollInterval,
			SecretIdentity:         flattenSecretIdentityTyped(azureKeyvaultSecretsProviderIdentity),
		}}
	}

	return
}

func flattenKubernetesClusterAddOnIdentityProfileTyped(profile *managedclusters.UserAssignedIdentity) []ConnectorIdentityModel {
	if profile == nil {
		return []ConnectorIdentityModel{}
	}

	clientID := ""
	if profile.ClientId != nil {
		clientID = *profile.ClientId
	}

	objectID := ""
	if profile.ObjectId != nil {
		objectID = *profile.ObjectId
	}

	userAssignedIdentityID := ""
	if profile.ResourceId != nil {
		userAssignedIdentityID = *profile.ResourceId
	}

	return []ConnectorIdentityModel{{
		ClientID:               clientID,
		ObjectID:               objectID,
		UserAssignedIdentityID: userAssignedIdentityID,
	}}
}

func flattenOMSAgentIdentityTyped(input []ConnectorIdentityModel) []OMSAgentIdentityModel {
	if len(input) == 0 {
		return []OMSAgentIdentityModel{}
	}
	return []OMSAgentIdentityModel{{
		ClientID:               input[0].ClientID,
		ObjectID:               input[0].ObjectID,
		UserAssignedIdentityID: input[0].UserAssignedIdentityID,
	}}
}

func flattenIngressApplicationGatewayIdentityTyped(input []ConnectorIdentityModel) []IngressApplicationGatewayIdentityModel {
	if len(input) == 0 {
		return []IngressApplicationGatewayIdentityModel{}
	}
	return []IngressApplicationGatewayIdentityModel{{
		ClientID:               input[0].ClientID,
		ObjectID:               input[0].ObjectID,
		UserAssignedIdentityID: input[0].UserAssignedIdentityID,
	}}
}

func flattenSecretIdentityTyped(input []ConnectorIdentityModel) []SecretIdentityModel {
	if len(input) == 0 {
		return []SecretIdentityModel{}
	}
	return []SecretIdentityModel{{
		ClientID:               input[0].ClientID,
		ObjectID:               input[0].ObjectID,
		UserAssignedIdentityID: input[0].UserAssignedIdentityID,
	}}
}
