// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/containerapps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/daprcomponents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type Registry struct {
	PasswordSecretRef string `tfschema:"password_secret_name"`
	Server            string `tfschema:"server"`
	UserName          string `tfschema:"username"`
	Identity          string `tfschema:"identity"`
}

func ContainerAppRegistrySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MinItems: 1,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"server": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The hostname for the Container Registry.",
				},

				"username": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The username to use for this Container Registry.",
				},

				"password_secret_name": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The name of the Secret Reference containing the password value for this user on the Container Registry.",
				},

				"identity": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "ID of the System or User Managed Identity used to pull images from the Container Registry",
				},
			},
		},
	}
}

func ContainerAppRegistrySchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"server": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The hostname for the Container Registry.",
				},

				"username": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The username to use for this Container Registry.",
				},

				"password_secret_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The name of the Secret Reference containing the password value for this user on the Container Registry.",
				},

				"identity": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "ID of the System or User Managed Identity used to pull images from the Container Registry",
				},
			},
		},
	}
}

func ValidateContainerAppRegistry(r Registry) error {
	if r.Identity != "" && (r.UserName != "" || r.PasswordSecretRef != "") {
		return fmt.Errorf("identity and username/password_secret_name are mutually exclusive")
	}
	if r.Identity == "" && r.UserName == "" && r.PasswordSecretRef == "" {
		return fmt.Errorf("must supply either identity or username/password_secret_name")
	}
	if (r.UserName != "" && r.PasswordSecretRef == "") || (r.UserName == "" && r.PasswordSecretRef != "") {
		return fmt.Errorf("must supply both username and password_secret_name")
	}
	return nil
}

func ExpandContainerAppRegistries(input []Registry) (*[]containerapps.RegistryCredentials, error) {
	if input == nil {
		return nil, nil
	}

	registries := make([]containerapps.RegistryCredentials, 0)
	for _, v := range input {
		if err := ValidateContainerAppRegistry(v); err != nil {
			return nil, err
		}
		registries = append(registries, containerapps.RegistryCredentials{
			Server:            pointer.To(v.Server),
			Username:          pointer.To(v.UserName),
			PasswordSecretRef: pointer.To(v.PasswordSecretRef),
			Identity:          pointer.To(v.Identity),
		})
	}

	return &registries, nil
}

func FlattenContainerAppRegistries(input *[]containerapps.RegistryCredentials) []Registry {
	if input == nil || len(*input) == 0 {
		return []Registry{}
	}

	result := make([]Registry, 0)
	for _, v := range *input {
		result = append(result, Registry{
			PasswordSecretRef: pointer.From(v.PasswordSecretRef),
			Server:            pointer.From(v.Server),
			UserName:          pointer.From(v.Username),
			Identity:          pointer.From(v.Identity),
		})
	}

	return result
}

type Ingress struct {
	AllowInsecure          bool                    `tfschema:"allow_insecure_connections"`
	CustomDomains          []CustomDomain          `tfschema:"custom_domain"`
	IsExternal             bool                    `tfschema:"external_enabled"`
	FQDN                   string                  `tfschema:"fqdn"`
	TargetPort             int64                   `tfschema:"target_port"`
	ExposedPort            int64                   `tfschema:"exposed_port"`
	TrafficWeights         []TrafficWeight         `tfschema:"traffic_weight"`
	Transport              string                  `tfschema:"transport"`
	IpSecurityRestrictions []IpSecurityRestriction `tfschema:"ip_security_restriction"`
	ClientCertificateMode  string                  `tfschema:"client_certificate_mode"`
}

func ContainerAppIngressSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"allow_insecure_connections": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should this ingress allow insecure connections?",
				},

				"external_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Is this an external Ingress.",
				},

				"custom_domain": ContainerAppIngressCustomDomainSchemaComputed(),

				"fqdn": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The FQDN of the ingress.",
				},

				"ip_security_restriction": ContainerAppIngressIpSecurityRestriction(),

				"target_port": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
					Description:  "The target port on the container for the Ingress traffic.",
				},

				"exposed_port": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
					Description:  "The exposed port on the container for the Ingress traffic.",
				},

				"traffic_weight": ContainerAppIngressTrafficWeight(),

				"transport": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(containerapps.IngressTransportMethodAuto),
					ValidateFunc: validation.StringInSlice(containerapps.PossibleValuesForIngressTransportMethod(), false),
					Description:  "The transport method for the Ingress. Possible values include `auto`, `http`, and `http2`, `tcp`. Defaults to `auto`",
				},

				"client_certificate_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(containerapps.IngressClientCertificateModeAccept),
						string(containerapps.IngressClientCertificateModeRequire),
						string(containerapps.IngressClientCertificateModeIgnore),
					}, false),
					Description: "Client certificate mode for mTLS authentication. Ignore indicates server drops client certificate on forwarding. Accept indicates server forwards client certificate but does not require a client certificate. Require indicates server requires a client certificate.",
				},
			},
		},
	}
}

func ContainerAppIngressSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"allow_insecure_connections": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "Should this ingress allow insecure connections?",
				},

				"custom_domain": ContainerAppIngressCustomDomainSchemaComputed(),

				"external_enabled": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "Is this an external Ingress.",
				},

				"fqdn": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The FQDN of the ingress.",
				},

				"ip_security_restriction": ContainerAppIngressIpSecurityRestrictionComputed(),

				"target_port": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The target port on the container for the Ingress traffic.",
				},

				"exposed_port": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The exposed port on the container for the Ingress traffic.",
				},

				"traffic_weight": ContainerAppIngressTrafficWeightComputed(),

				"transport": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The transport method for the Ingress. Possible values include `auto`, `http`, and `http2`, `tcp`. Defaults to `auto`",
				},

				"client_certificate_mode": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "Client certificate mode for mTLS authentication. Ignore indicates server drops client certificate on forwarding. Accept indicates server forwards client certificate but does not require a client certificate. Require indicates server requires a client certificate.",
				},
			},
		},
	}
}

func ExpandContainerAppIngress(input []Ingress, appName string) *containerapps.Ingress {
	if len(input) == 0 {
		return nil
	}

	ingress := input[0]
	result := &containerapps.Ingress{
		AllowInsecure:          pointer.To(ingress.AllowInsecure),
		CustomDomains:          expandContainerAppIngressCustomDomain(ingress.CustomDomains),
		External:               pointer.To(ingress.IsExternal),
		Fqdn:                   pointer.To(ingress.FQDN),
		TargetPort:             pointer.To(ingress.TargetPort),
		ExposedPort:            pointer.To(ingress.ExposedPort),
		Traffic:                expandContainerAppIngressTraffic(ingress.TrafficWeights, appName),
		IPSecurityRestrictions: expandIpSecurityRestrictions(ingress.IpSecurityRestrictions),
	}
	transport := containerapps.IngressTransportMethod(ingress.Transport)
	result.Transport = &transport
	if ingress.ClientCertificateMode != "" {
		clientCertificateMode := containerapps.IngressClientCertificateMode(ingress.ClientCertificateMode)
		result.ClientCertificateMode = &clientCertificateMode
	}

	return result
}

func FlattenContainerAppIngress(input *containerapps.Ingress, appName string) []Ingress {
	if input == nil {
		return []Ingress{}
	}

	ingress := *input
	result := Ingress{
		AllowInsecure:          pointer.From(ingress.AllowInsecure),
		CustomDomains:          flattenContainerAppIngressCustomDomain(ingress.CustomDomains),
		IsExternal:             pointer.From(ingress.External),
		FQDN:                   pointer.From(ingress.Fqdn),
		TargetPort:             pointer.From(ingress.TargetPort),
		ExposedPort:            pointer.From(ingress.ExposedPort),
		TrafficWeights:         flattenContainerAppIngressTraffic(ingress.Traffic, appName),
		IpSecurityRestrictions: flattenContainerAppIngressIpSecurityRestrictions(ingress.IPSecurityRestrictions),
	}

	if ingress.Transport != nil {
		result.Transport = strings.ToLower(string(*ingress.Transport))
	}

	if ingress.ClientCertificateMode != nil {
		result.ClientCertificateMode = string(*ingress.ClientCertificateMode)
	}

	return []Ingress{result}
}

type CustomDomain struct {
	CertBinding   string `tfschema:"certificate_binding_type"`
	CertificateId string `tfschema:"certificate_id"`
	Name          string `tfschema:"name"`
}

func ContainerAppIngressCustomDomainSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"certificate_binding_type": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The Binding type. Possible values include `Disabled` and `SniEnabled`. Defaults to `Disabled`",
				},

				"certificate_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The ID of the Certificate.",
				},

				"name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The hostname of the Certificate. Must be the CN or a named SAN in the certificate.",
				},
			},
		},
	}
}

func expandContainerAppIngressCustomDomain(input []CustomDomain) *[]containerapps.CustomDomain {
	if len(input) == 0 {
		return nil
	}

	result := make([]containerapps.CustomDomain, 0)
	for _, v := range input {
		customDomain := containerapps.CustomDomain{
			Name:          v.Name,
			CertificateId: pointer.To(v.CertificateId),
		}
		bindingType := containerapps.BindingType(v.CertBinding)
		customDomain.BindingType = &bindingType

		result = append(result, customDomain)
	}

	return &result
}

func flattenContainerAppIngressCustomDomain(input *[]containerapps.CustomDomain) []CustomDomain {
	if input == nil {
		return []CustomDomain{}
	}

	result := make([]CustomDomain, 0)

	for _, v := range *input {
		customDomain := CustomDomain{
			Name: v.Name,
		}
		if v.BindingType != nil {
			customDomain.CertBinding = string(*v.BindingType)
		}
		if v.CertificateId != nil {
			customDomain.CertificateId = *v.CertificateId
		}
		result = append(result, customDomain)
	}

	return result
}

func flattenContainerAppIngressIpSecurityRestrictions(input *[]containerapps.IPSecurityRestrictionRule) []IpSecurityRestriction {
	if input == nil {
		return []IpSecurityRestriction{}
	}

	result := make([]IpSecurityRestriction, 0)
	for _, v := range *input {
		ipSecurityRestriction := IpSecurityRestriction{
			Description:    pointer.From(v.Description),
			IpAddressRange: v.IPAddressRange,
			Action:         string(v.Action),
			Name:           v.Name,
		}

		result = append(result, ipSecurityRestriction)
	}

	return result
}

type TrafficWeight struct {
	Label          string `tfschema:"label"`
	LatestRevision bool   `tfschema:"latest_revision"`
	RevisionSuffix string `tfschema:"revision_suffix"`
	Weight         int64  `tfschema:"percentage"`
}

type IpSecurityRestriction struct {
	Action         string `tfschema:"action"`
	Description    string `tfschema:"description"`
	IpAddressRange string `tfschema:"ip_address_range"`
	Name           string `tfschema:"name"`
}

func ContainerAppIngressIpSecurityRestriction() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"action": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(containerapps.PossibleValuesForAction(), false),
					Description:  "The action. Allow or Deny.",
				},

				"description": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "Describe the IP restriction rule that is being sent to the container-app.",
				},

				"ip_address_range": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.Any(validation.IsCIDR, validation.IsIPAddress),
					Description:  "The incoming IP address or range of IP addresses (in CIDR notation).",
				},

				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "Name for the IP restriction rule.",
				},
			},
		},
	}
}

func ContainerAppIngressIpSecurityRestrictionComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"action": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The action. Allow or Deny.",
				},

				"description": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "Describe the IP restriction rule that is being sent to the container-app.",
				},

				"ip_address_range": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "CIDR notation to match incoming IP address.",
				},

				"name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "Name for the IP restriction rule.",
				},
			},
		},
	}
}

func ContainerAppIngressTrafficWeight() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"label": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The label to apply to the revision as a name prefix for routing traffic.",
				},

				"revision_suffix": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The suffix string to append to the revision. This must be unique for the Container App's lifetime. A default hash created by the service will be used if this value is omitted.",
				},

				"latest_revision": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "This traffic Weight relates to the latest stable Container Revision.",
				},

				"percentage": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(0, 100),
					Description:  "The percentage of traffic to send to this revision.",
				},
			},
		},
	}
}

func ContainerAppIngressTrafficWeightComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"label": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The label to apply to the revision as a name prefix for routing traffic.",
				},

				"revision_suffix": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The suffix string to append to the revision. This must be unique for the Container App's lifetime. A default hash created by the service will be used if this value is omitted.",
				},

				"latest_revision": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "This traffic Weight relates to the latest stable Container Revision.",
				},

				"percentage": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The percentage of traffic to send to this revision.",
				},
			},
		},
	}
}

func expandContainerAppIngressTraffic(input []TrafficWeight, appName string) *[]containerapps.TrafficWeight {
	if len(input) == 0 {
		return nil
	}

	result := make([]containerapps.TrafficWeight, 0)

	for _, v := range input {
		traffic := containerapps.TrafficWeight{
			LatestRevision: pointer.To(v.LatestRevision),
			Weight:         pointer.To(v.Weight),
		}

		if !v.LatestRevision {
			traffic.RevisionName = pointer.To(fmt.Sprintf("%s--%s", appName, v.RevisionSuffix))
		}

		if v.Label != "" {
			traffic.Label = pointer.To(v.Label)
		}

		result = append(result, traffic)
	}

	return &result
}

func flattenContainerAppIngressTraffic(input *[]containerapps.TrafficWeight, appName string) []TrafficWeight {
	if input == nil {
		return []TrafficWeight{}
	}

	result := make([]TrafficWeight, 0)
	for _, v := range *input {
		prefix := fmt.Sprintf("%s--", appName)
		result = append(result, TrafficWeight{
			Label:          pointer.From(v.Label),
			LatestRevision: pointer.From(v.LatestRevision),
			RevisionSuffix: strings.TrimPrefix(pointer.From(v.RevisionName), prefix),
			Weight:         pointer.From(v.Weight),
		})
	}

	return result
}

func expandIpSecurityRestrictions(input []IpSecurityRestriction) *[]containerapps.IPSecurityRestrictionRule {
	if input == nil {
		return &[]containerapps.IPSecurityRestrictionRule{}
	}

	result := make([]containerapps.IPSecurityRestrictionRule, 0)
	for _, v := range input {
		ipSecurityRestrictionRule := containerapps.IPSecurityRestrictionRule{
			Action:         containerapps.Action(v.Action),
			Name:           v.Name,
			IPAddressRange: v.IpAddressRange,
			Description:    pointer.To(v.Description),
		}
		result = append(result, ipSecurityRestrictionRule)
	}

	return &result
}

type Dapr struct {
	AppId       string `tfschema:"app_id"`
	AppPort     int64  `tfschema:"app_port"`
	AppProtocol string `tfschema:"app_protocol"`
}

func ContainerDaprSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"app_id": {
					Type:        pluginsdk.TypeString,
					Required:    true,
					Description: "The Dapr Application Identifier.",
				},

				"app_port": {
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					Description: "The port which the application is listening on. This is the same as the `ingress` port.",
				},

				"app_protocol": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(containerapps.AppProtocolHTTP),
					ValidateFunc: validation.StringInSlice([]string{
						string(containerapps.AppProtocolHTTP),
						string(containerapps.AppProtocolGrpc),
					}, false),
					Description: "The protocol for the app. Possible values include `http` and `grpc`. Defaults to `http`.",
				},
			},
		},
	}
}

func ContainerDaprSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"app_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The Dapr Application Identifier.",
				},

				"app_port": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The port which the application is listening on. This is the same as the `ingress` port.",
				},

				"app_protocol": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The protocol for the app. Possible values include `http` and `grpc`. Defaults to `http`.",
				},
			},
		},
	}
}

func ExpandContainerAppDapr(input []Dapr) *containerapps.Dapr {
	if len(input) == 0 {
		return nil
	}

	dapr := input[0]
	if dapr.AppId == "" {
		return &containerapps.Dapr{
			Enabled: pointer.To(false),
		}
	}

	appProtocol := containerapps.AppProtocol(dapr.AppProtocol)

	return &containerapps.Dapr{
		AppId:       pointer.To(dapr.AppId),
		AppPort:     pointer.To(dapr.AppPort),
		AppProtocol: &appProtocol,
		Enabled:     pointer.To(true),
	}
}

func FlattenContainerAppDapr(input *containerapps.Dapr) []Dapr {
	if input == nil {
		return []Dapr{}
	}

	result := Dapr{
		AppId:   pointer.From(input.AppId),
		AppPort: pointer.From(input.AppPort),
	}
	if appProtocol := input.AppProtocol; appProtocol != nil {
		result.AppProtocol = string(*appProtocol)
	}

	return []Dapr{result}
}

type DaprMetadata struct {
	Name       string `tfschema:"name"`
	Value      string `tfschema:"value"`
	SecretName string `tfschema:"secret_name"`
}

func ContainerAppEnvironmentDaprMetadataSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name of the Metadata configuration item.",
				},

				"value": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The value for this metadata configuration item.",
				},

				"secret_name": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The name of a secret specified in the `secrets` block that contains the value for this metadata configuration item.",
				},
			},
		},
	}
}

type ContainerTemplate struct {
	Containers             []Container           `tfschema:"container"`
	InitContainers         []BaseContainer       `tfschema:"init_container"`
	Suffix                 string                `tfschema:"revision_suffix"`
	MinReplicas            int64                 `tfschema:"min_replicas"`
	MaxReplicas            int64                 `tfschema:"max_replicas"`
	AzureQueueScaleRules   []AzureQueueScaleRule `tfschema:"azure_queue_scale_rule"`
	CustomScaleRules       []CustomScaleRule     `tfschema:"custom_scale_rule"`
	HTTPScaleRules         []HTTPScaleRule       `tfschema:"http_scale_rule"`
	TCPScaleRules          []TCPScaleRule        `tfschema:"tcp_scale_rule"`
	Volumes                []ContainerVolume     `tfschema:"volume"`
	TerminationGracePeriod int64                 `tfschema:"termination_grace_period_seconds"`
}

func ContainerTemplateSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Required: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"container": ContainerAppContainerSchema(),

				"init_container": InitContainerAppContainerSchema(),

				"min_replicas": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      0,
					ValidateFunc: validation.IntBetween(0, 300),
					Description:  "The minimum number of replicas for this container.",
				},

				"max_replicas": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      10,
					ValidateFunc: validation.IntBetween(1, 300),
					Description:  "The maximum number of replicas for this container.",
				},

				"azure_queue_scale_rule": AzureQueueScaleRuleSchema(),

				"custom_scale_rule": CustomScaleRuleSchema(),

				"http_scale_rule": HTTPScaleRuleSchema(),

				"tcp_scale_rule": TCPScaleRuleSchema(),

				"volume": ContainerVolumeSchema(),

				"revision_suffix": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Computed:    true, // Note: O+C This value is always present and non-zero but if not user specified, then the service will generate a value.
					Description: "The suffix for the revision. This value must be unique for the lifetime of the Resource. If omitted the service will use a hash function to create one.",
				},

				"termination_grace_period_seconds": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      0,
					ValidateFunc: validation.IntBetween(0, 600),
					Description:  "The time in seconds after the container is sent the termination signal before the process if forcibly killed.",
				},
			},
		},
	}
}

func ContainerTemplateSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"container": ContainerAppContainerSchemaComputed(),

				"init_container": InitContainerAppContainerSchemaComputed(),

				"min_replicas": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The minimum number of replicas for this container.",
				},

				"max_replicas": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The maximum number of replicas for this container.",
				},

				"azure_queue_scale_rule": AzureQueueScaleRuleSchemaComputed(),

				"custom_scale_rule": CustomScaleRuleSchemaComputed(),

				"http_scale_rule": HTTPScaleRuleSchemaComputed(),

				"tcp_scale_rule": TCPScaleRuleSchemaComputed(),

				"volume": ContainerVolumeSchemaComputed(),

				"revision_suffix": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"termination_grace_period_seconds": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func ExpandContainerAppTemplate(input []ContainerTemplate, metadata sdk.ResourceMetaData) *containerapps.Template {
	if len(input) != 1 {
		return nil
	}

	config := input[0]
	template := &containerapps.Template{
		Containers:     expandContainerAppContainers(config.Containers),
		InitContainers: expandInitContainerAppContainers(config.InitContainers),
		Volumes:        expandContainerAppVolumes(config.Volumes),
	}

	if config.MaxReplicas != 0 {
		if template.Scale == nil {
			template.Scale = &containerapps.Scale{}
		}
		template.Scale.MaxReplicas = pointer.To(config.MaxReplicas)
	}

	if config.TerminationGracePeriod != 0 {
		template.TerminationGracePeriodSeconds = pointer.To(config.TerminationGracePeriod)
	}

	if config.MinReplicas != 0 {
		if template.Scale == nil {
			template.Scale = &containerapps.Scale{}
		}
		template.Scale.MinReplicas = pointer.To(config.MinReplicas)
	}

	if rules := config.expandContainerAppScaleRules(); len(rules) != 0 {
		if template.Scale == nil {
			template.Scale = &containerapps.Scale{}
		}

		template.Scale.Rules = pointer.To(rules)
	}

	if config.Suffix != "" {
		if metadata.ResourceData.HasChange("template.0.revision_suffix") {
			template.RevisionSuffix = pointer.To(config.Suffix)
		}
	}

	return template
}

func FlattenContainerAppTemplate(input *containerapps.Template) []ContainerTemplate {
	if input == nil {
		return []ContainerTemplate{}
	}
	result := ContainerTemplate{
		Containers:             flattenContainerAppContainers(input.Containers),
		InitContainers:         flattenInitContainerAppContainers(input.InitContainers),
		Suffix:                 pointer.From(input.RevisionSuffix),
		TerminationGracePeriod: pointer.From(input.TerminationGracePeriodSeconds),
		Volumes:                flattenContainerAppVolumes(input.Volumes),
	}

	if scale := input.Scale; scale != nil {
		result.MaxReplicas = pointer.From(scale.MaxReplicas)
		result.MinReplicas = pointer.From(scale.MinReplicas)
		result.flattenContainerAppScaleRules(scale.Rules)
	}

	return []ContainerTemplate{result}
}

type Container struct {
	Name             string                       `tfschema:"name"`
	Image            string                       `tfschema:"image"`
	CPU              float64                      `tfschema:"cpu"`
	Memory           string                       `tfschema:"memory"`
	EphemeralStorage string                       `tfschema:"ephemeral_storage"`
	Env              []ContainerEnvVar            `tfschema:"env"`
	Args             []string                     `tfschema:"args"`
	Command          []string                     `tfschema:"command"`
	LivenessProbe    []ContainerAppLivenessProbe  `tfschema:"liveness_probe"`
	ReadinessProbe   []ContainerAppReadinessProbe `tfschema:"readiness_probe"`
	StartupProbe     []ContainerAppStartupProbe   `tfschema:"startup_probe"`
	VolumeMounts     []ContainerVolumeMount       `tfschema:"volume_mounts"`
}

func ContainerAppContainerSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validate.ContainerAppContainerName,
					Description:  "The name of the container.",
				},

				"image": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The image to use to create the container.",
				},

				"cpu": {
					Type:         pluginsdk.TypeFloat,
					Required:     true,
					ValidateFunc: validation.FloatAtLeast(0.1),
					Description:  "The amount of vCPU to allocate to the container. Possible values include `0.25`, `0.5`, `0.75`, `1.0`, `1.25`, `1.5`, `1.75`, and `2.0`. **NOTE:** `cpu` and `memory` must be specified in `0.25'/'0.5Gi` combination increments. e.g. `1.0` / `2.0` or `0.5` / `1.0`. When there's a workload profile specified, there's no such constraint.",
				},

				"memory": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The amount of memory to allocate to the container. Possible values include `0.5Gi`, `1.0Gi`, `1.5Gi`, `2.0Gi`, `2.5Gi`, `3.0Gi`, `3.5Gi`, and `4.0Gi`. **NOTE:** `cpu` and `memory` must be specified in `0.25'/'0.5Gi` combination increments. e.g. `1.25` / `2.5Gi` or `0.75` / `1.5Gi`. When there's a workload profile specified, there's no such constraint.",
				},

				"ephemeral_storage": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The amount of ephemeral storage available to the Container App.",
				},

				"env": ContainerEnvVarSchema(),

				"args": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "A list of args to pass to the container.",
				},

				"command": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "A command to pass to the container to override the default. This is provided as a list of command line elements without spaces.",
				},

				"liveness_probe": ContainerAppLivenessProbeSchema(),

				"readiness_probe": ContainerAppReadinessProbeSchema(),

				"startup_probe": ContainerAppStartupProbeSchema(),

				"volume_mounts": ContainerVolumeMountSchema(),
			},
		},
	}
}

func ContainerAppContainerSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The name of the container.",
				},

				"image": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The image to use to create the container.",
				},

				"cpu": {
					Type:        pluginsdk.TypeFloat,
					Computed:    true,
					Description: "The amount of vCPU to allocate to the container. Possible values include `0.25`, `0.5`, `0.75`, `1.0`, `1.25`, `1.5`, `1.75`, and `2.0`. **NOTE:** `cpu` and `memory` must be specified in `0.25'/'0.5Gi` combination increments. e.g. `1.0` / `2.0` or `0.5` / `1.0`",
				},

				"memory": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The amount of memory to allocate to the container. Possible values include `0.5Gi`, `1.0Gi`, `1.5Gi`, `2.0Gi`, `2.5Gi`, `3.0Gi`, `3.5Gi`, and `4.0Gi`. **NOTE:** `cpu` and `memory` must be specified in `0.25'/'0.5Gi` combination increments. e.g. `1.25` / `2.5Gi` or `0.75` / `1.5Gi`",
				},

				"ephemeral_storage": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The amount of ephemeral storage available to the Container App.",
				},

				"env": ContainerEnvVarSchemaComputed(),

				"args": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "A list of args to pass to the container.",
				},

				"command": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "A command to pass to the container to override the default. This is provided as a list of command line elements without spaces.",
				},

				"liveness_probe": ContainerAppLivenessProbeSchemaComputed(),

				"readiness_probe": ContainerAppReadinessProbeSchemaComputed(),

				"startup_probe": ContainerAppStartupProbeSchemaComputed(),

				"volume_mounts": ContainerVolumeMountSchemaComputed(),
			},
		},
	}
}

type BaseContainer struct {
	Name             string                 `tfschema:"name"`
	Image            string                 `tfschema:"image"`
	CPU              float64                `tfschema:"cpu"`
	Memory           string                 `tfschema:"memory"`
	EphemeralStorage string                 `tfschema:"ephemeral_storage"`
	Env              []ContainerEnvVar      `tfschema:"env"`
	Args             []string               `tfschema:"args"`
	Command          []string               `tfschema:"command"`
	VolumeMounts     []ContainerVolumeMount `tfschema:"volume_mounts"`
}

func InitContainerAppContainerSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validate.ContainerAppContainerName,
					Description:  "The name of the container.",
				},

				"image": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The image to use to create the container.",
				},

				"cpu": {
					Type:         pluginsdk.TypeFloat,
					Optional:     true,
					ValidateFunc: validation.FloatAtLeast(0.1),
					Description:  "The amount of vCPU to allocate to the container. Possible values include `0.25`, `0.5`, `0.75`, `1.0`, `1.25`, `1.5`, `1.75`, and `2.0`. **NOTE:** `cpu` and `memory` must be specified in `0.25'/'0.5Gi` combination increments. e.g. `1.0` / `2.0` or `0.5` / `1.0`. When there's a workload profile specified, there's no such constraint.",
				},

				"memory": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The amount of memory to allocate to the container. Possible values include `0.5Gi`, `1.0Gi`, `1.5Gi`, `2.0Gi`, `2.5Gi`, `3.0Gi`, `3.5Gi`, and `4.0Gi`. **NOTE:** `cpu` and `memory` must be specified in `0.25'/'0.5Gi` combination increments. e.g. `1.25` / `2.5Gi` or `0.75` / `1.5Gi`. When there's a workload profile specified, there's no such constraint.",
				},

				"ephemeral_storage": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The amount of ephemeral storage available to the Container App.",
				},

				"env": ContainerEnvVarSchema(),

				"args": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "A list of args to pass to the container.",
				},

				"command": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "A command to pass to the container to override the default. This is provided as a list of command line elements without spaces.",
				},

				"volume_mounts": ContainerVolumeMountSchema(),
			},
		},
	}
}

func InitContainerAppContainerSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The name of the container.",
				},

				"image": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The image to use to create the container.",
				},

				"cpu": {
					Type:        pluginsdk.TypeFloat,
					Computed:    true,
					Description: "The amount of vCPU to allocate to the container. Possible values include `0.25`, `0.5`, `0.75`, `1.0`, `1.25`, `1.5`, `1.75`, and `2.0`. **NOTE:** `cpu` and `memory` must be specified in `0.25'/'0.5Gi` combination increments. e.g. `1.0` / `2.0` or `0.5` / `1.0`",
				},

				"memory": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The amount of memory to allocate to the container. Possible values include `0.5Gi`, `1.0Gi`, `1.5Gi`, `2.0Gi`, `2.5Gi`, `3.0Gi`, `3.5Gi`, and `4.0Gi`. **NOTE:** `cpu` and `memory` must be specified in `0.25'/'0.5Gi` combination increments. e.g. `1.25` / `2.5Gi` or `0.75` / `1.5Gi`",
				},

				"ephemeral_storage": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The amount of ephemeral storage available to the Container App.",
				},

				"env": ContainerEnvVarSchemaComputed(),

				"args": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "A list of args to pass to the container.",
				},

				"command": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "A command to pass to the container to override the default. This is provided as a list of command line elements without spaces.",
				},

				"volume_mounts": ContainerVolumeMountSchemaComputed(),
			},
		},
	}
}

func expandInitContainerAppContainers(input []BaseContainer) *[]containerapps.BaseContainer {
	if input == nil {
		return nil
	}

	result := make([]containerapps.BaseContainer, 0)
	for _, v := range input {
		container := containerapps.BaseContainer{
			Env:   expandInitContainerEnvVar(v),
			Image: pointer.To(v.Image),
			Name:  pointer.To(v.Name),
			Resources: &containerapps.ContainerResources{
				Cpu:              pointer.To(v.CPU),
				EphemeralStorage: pointer.To(v.EphemeralStorage),
				Memory:           pointer.To(v.Memory),
			},
			VolumeMounts: expandContainerVolumeMounts(v.VolumeMounts),
		}
		if len(v.Args) != 0 {
			container.Args = pointer.To(v.Args)
		}
		if len(v.Command) != 0 {
			container.Command = pointer.To(v.Command)
		}

		result = append(result, container)
	}

	return &result
}

func flattenInitContainerAppContainers(input *[]containerapps.BaseContainer) []BaseContainer {
	if input == nil || len(*input) == 0 {
		return []BaseContainer{}
	}
	result := make([]BaseContainer, 0)
	for _, v := range *input {
		container := BaseContainer{
			Name:         pointer.From(v.Name),
			Image:        pointer.From(v.Image),
			Args:         pointer.From(v.Args),
			Command:      pointer.From(v.Command),
			Env:          flattenContainerEnvVar(v.Env),
			VolumeMounts: flattenContainerVolumeMounts(v.VolumeMounts),
		}

		if resources := v.Resources; resources != nil {
			container.CPU = pointer.From(resources.Cpu)
			container.Memory = pointer.From(resources.Memory)
			container.EphemeralStorage = pointer.From(resources.EphemeralStorage)
		}

		result = append(result, container)
	}
	return result
}

func expandContainerAppContainers(input []Container) *[]containerapps.Container {
	if input == nil {
		return nil
	}

	result := make([]containerapps.Container, 0)
	for _, v := range input {
		container := containerapps.Container{
			Env:    expandContainerEnvVar(v),
			Image:  pointer.To(v.Image),
			Name:   pointer.To(v.Name),
			Probes: expandContainerProbes(v),
			Resources: &containerapps.ContainerResources{
				Cpu:              pointer.To(v.CPU),
				EphemeralStorage: pointer.To(v.EphemeralStorage),
				Memory:           pointer.To(v.Memory),
			},
			VolumeMounts: expandContainerVolumeMounts(v.VolumeMounts),
		}
		if len(v.Args) != 0 {
			container.Args = pointer.To(v.Args)
		}
		if len(v.Command) != 0 {
			container.Command = pointer.To(v.Command)
		}

		result = append(result, container)
	}

	return &result
}

func flattenContainerAppContainers(input *[]containerapps.Container) []Container {
	if input == nil || len(*input) == 0 {
		return []Container{}
	}
	result := make([]Container, 0)
	for _, v := range *input {
		container := Container{
			Name:         pointer.From(v.Name),
			Image:        pointer.From(v.Image),
			Args:         pointer.From(v.Args),
			Command:      pointer.From(v.Command),
			Env:          flattenContainerEnvVar(v.Env),
			VolumeMounts: flattenContainerVolumeMounts(v.VolumeMounts),
		}
		if v.Probes != nil {
			for _, p := range *v.Probes {
				switch *p.Type {
				case containerapps.TypeLiveness:
					container.LivenessProbe = flattenContainerAppLivenessProbe(p)
				case containerapps.TypeReadiness:
					container.ReadinessProbe = flattenContainerAppReadinessProbe(p)
				case containerapps.TypeStartup:
					container.StartupProbe = flattenContainerAppStartupProbe(p)
				}
			}
		}

		if resources := v.Resources; resources != nil {
			container.CPU = pointer.From(resources.Cpu)
			container.Memory = pointer.From(resources.Memory)
			container.EphemeralStorage = pointer.From(resources.EphemeralStorage)
		}

		result = append(result, container)
	}
	return result
}

type ContainerVolume struct {
	Name         string `tfschema:"name"`
	StorageName  string `tfschema:"storage_name"`
	StorageType  string `tfschema:"storage_type"`
	MountOptions string `tfschema:"mount_options"`
}

func ContainerVolumeSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name of the volume.",
				},

				"storage_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "EmptyDir",
					ValidateFunc: validation.StringInSlice(
						containerapps.PossibleValuesForStorageType(),
						false),
					Description: "The type of storage volume. Possible values include `AzureFile` and `EmptyDir`. Defaults to `EmptyDir`.",
				},

				"storage_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validate.ManagedEnvironmentStorageName,
					Description:  "The name of the `AzureFile` storage. Required when `storage_type` is `AzureFile`",
				},

				"mount_options": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "Mount options used while mounting the AzureFile. Must be a comma-separated string.",
				},
			},
		},
	}
}

func ContainerVolumeSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"storage_type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"storage_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"mount_options": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func expandContainerAppVolumes(input []ContainerVolume) *[]containerapps.Volume {
	if input == nil {
		return nil
	}

	volumes := make([]containerapps.Volume, 0)

	for _, v := range input {
		volume := containerapps.Volume{
			Name: pointer.To(v.Name),
		}
		if v.StorageName != "" {
			volume.StorageName = pointer.To(v.StorageName)
		}
		if v.StorageType != "" {
			storageType := containerapps.StorageType(v.StorageType)
			volume.StorageType = &storageType
		}
		if v.MountOptions != "" {
			volume.MountOptions = pointer.To(v.MountOptions)
		}
		volumes = append(volumes, volume)
	}

	return &volumes
}

func flattenContainerAppVolumes(input *[]containerapps.Volume) []ContainerVolume {
	if input == nil || len(*input) == 0 {
		return []ContainerVolume{}
	}

	result := make([]ContainerVolume, 0)
	for _, v := range *input {
		containerVolume := ContainerVolume{
			Name:        pointer.From(v.Name),
			StorageName: pointer.From(v.StorageName),
		}
		if v.StorageType != nil {
			containerVolume.StorageType = string(*v.StorageType)
		}
		if v.MountOptions != nil {
			containerVolume.MountOptions = pointer.From(v.MountOptions)
		}

		result = append(result, containerVolume)
	}

	return result
}

type ContainerVolumeMount struct {
	Name    string `tfschema:"name"`
	Path    string `tfschema:"path"`
	SubPath string `tfschema:"sub_path"`
}

func ContainerVolumeMountSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name of the Volume to be mounted in the container.",
				},

				"path": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The path in the container at which to mount this volume.",
				},

				"sub_path": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The sub path of the volume to be mounted in the container.",
				},
			},
		},
	}
}

func ContainerVolumeMountSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The name of the Volume to be mounted in the container.",
				},

				"path": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The path in the container at which to mount this volume.",
				},

				"sub_path": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The sub path of the volume to be mounted in the container.",
				},
			},
		},
	}
}

func expandContainerVolumeMounts(input []ContainerVolumeMount) *[]containerapps.VolumeMount {
	if input == nil {
		return nil
	}

	volumeMounts := make([]containerapps.VolumeMount, 0)
	for _, v := range input {
		volumeMounts = append(volumeMounts, containerapps.VolumeMount{
			MountPath:  pointer.To(v.Path),
			VolumeName: pointer.To(v.Name),
			SubPath:    pointer.To(v.SubPath),
		})
	}

	return &volumeMounts
}

func flattenContainerVolumeMounts(input *[]containerapps.VolumeMount) []ContainerVolumeMount {
	if input == nil || len(*input) == 0 {
		return []ContainerVolumeMount{}
	}

	result := make([]ContainerVolumeMount, 0)
	for _, v := range *input {
		result = append(result, ContainerVolumeMount{
			Name:    pointer.From(v.VolumeName),
			Path:    pointer.From(v.MountPath),
			SubPath: pointer.From(v.SubPath),
		})
	}

	return result
}

type ContainerEnvVar struct {
	Name            string `tfschema:"name"`
	Value           string `tfschema:"value"`
	SecretReference string `tfschema:"secret_name"`
}

func ContainerEnvVarSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MinItems: 1,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name of the environment variable for the container.",
				},

				"value": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The value for this environment variable. **NOTE:** This value is ignored if `secret_name` is used",
				},

				"secret_name": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The name of the secret that contains the value for this environment variable.",
				},
			},
		},
	}
}

func ContainerEnvVarSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The name of the environment variable for the container.",
				},

				"value": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The value for this environment variable. **NOTE:** This value is ignored if `secret_name` is used",
				},

				"secret_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The name of the secret that contains the value for this environment variable.",
				},
			},
		},
	}
}

func expandInitContainerEnvVar(input BaseContainer) *[]containerapps.EnvironmentVar {
	envs := make([]containerapps.EnvironmentVar, 0)
	if len(input.Env) == 0 {
		return &envs
	}

	for _, v := range input.Env {
		env := containerapps.EnvironmentVar{
			Name: pointer.To(v.Name),
		}
		if v.SecretReference != "" {
			env.SecretRef = pointer.To(v.SecretReference)
		} else {
			env.Value = pointer.To(v.Value)
		}

		envs = append(envs, env)
	}

	return &envs
}

func expandContainerEnvVar(input Container) *[]containerapps.EnvironmentVar {
	envs := make([]containerapps.EnvironmentVar, 0)
	if len(input.Env) == 0 {
		return &envs
	}

	for _, v := range input.Env {
		env := containerapps.EnvironmentVar{
			Name: pointer.To(v.Name),
		}
		if v.SecretReference != "" {
			env.SecretRef = pointer.To(v.SecretReference)
		} else {
			env.Value = pointer.To(v.Value)
		}

		envs = append(envs, env)
	}

	return &envs
}

func flattenContainerEnvVar(input *[]containerapps.EnvironmentVar) []ContainerEnvVar {
	if input == nil || len(*input) == 0 {
		return []ContainerEnvVar{}
	}

	result := make([]ContainerEnvVar, 0)

	for _, v := range *input {
		result = append(result, ContainerEnvVar{
			Name:            pointer.From(v.Name),
			SecretReference: pointer.From(v.SecretRef),
			Value:           pointer.From(v.Value),
		})
	}

	return result
}

type ContainerAppReadinessProbe struct {
	Transport        string       `tfschema:"transport"`
	Host             string       `tfschema:"host"`
	Port             int64        `tfschema:"port"`
	Path             string       `tfschema:"path"`
	Headers          []HttpHeader `tfschema:"header"`
	InitialDelay     int64        `tfschema:"initial_delay"`
	Interval         int64        `tfschema:"interval_seconds"`
	Timeout          int64        `tfschema:"timeout"`
	FailureThreshold int64        `tfschema:"failure_count_threshold"`
	SuccessThreshold int64        `tfschema:"success_count_threshold"`
}

func ContainerAppReadinessProbeSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"transport": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						"TCP",
						string(containerapps.SchemeHTTP),
						string(containerapps.SchemeHTTPS),
					}, true),
					Description: "Type of probe. Possible values are `TCP`, `HTTP`, and `HTTPS`.",
				},

				"port": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
					Description:  "The port number on which to connect. Possible values are between `1` and `65535`.",
				},

				"host": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The probe hostname. Defaults to the pod IP address. Setting a value for `Host` in `headers` can be used to override this for `http` and `https` type probes.",
				},

				"path": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Computed:    true, // Note: O+C Needs to remain computed as this has a variable default and since it is part of a list we cannot diffsuppress it.
					Description: "The URI to use for http type probes. Not valid for `TCP` type probes. Defaults to `/`.",
				},

				"header": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
								Description:  "The HTTP Header Name.",
							},

							"value": {
								Type:        pluginsdk.TypeString,
								Required:    true,
								Description: "The HTTP Header value.",
							},
						},
					},
				},

				"initial_delay": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      0,
					ValidateFunc: validation.IntBetween(0, 60),
					Description:  "The number of seconds elapsed after the container has started before the probe is initiated. Possible values are between `0` and `60`. Defaults to `0` seconds.",
				},

				"interval_seconds": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      10,
					ValidateFunc: validation.IntBetween(1, 240),
					Description:  "How often, in seconds, the probe should run. Possible values are between `1` and `240`. Defaults to `10`",
				},

				"timeout": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      1,
					ValidateFunc: validation.IntBetween(1, 240),
					Description:  "Time in seconds after which the probe times out. Possible values are between `1` an `240`. Defaults to `1`.",
				},

				"failure_count_threshold": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      3,
					ValidateFunc: validation.IntBetween(1, 30),
					Description:  "The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `30`. Defaults to `3`.",
				},

				"success_count_threshold": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      3,
					ValidateFunc: validation.IntBetween(1, 10),
					Description:  "The number of consecutive successful responses required to consider this probe as successful. Possible values are between `1` and `10`. Defaults to `3`.",
				},
			},
		},
	}
}

func ContainerAppReadinessProbeSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"transport": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "Type of probe. Possible values are `TCP`, `HTTP`, and `HTTPS`.",
				},

				"port": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The port number on which to connect. Possible values are between `1` and `65535`.",
				},

				"host": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The probe hostname. Defaults to the pod IP address. Setting a value for `Host` in `headers` can be used to override this for `http` and `https` type probes.",
				},

				"path": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The URI to use for http type probes. Not valid for `TCP` type probes. Defaults to `/`.",
				},

				"header": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Type:        pluginsdk.TypeString,
								Computed:    true,
								Description: "The HTTP Header Name.",
							},

							"value": {
								Type:        pluginsdk.TypeString,
								Computed:    true,
								Description: "The HTTP Header value.",
							},
						},
					},
				},

				"initial_delay": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The number of seconds elapsed after the container has started before the probe is initiated. Possible values are between `0` and `60`. Defaults to `0` seconds.",
				},

				"interval_seconds": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "How often, in seconds, the probe should run. Possible values are between `1` and `240`. Defaults to `10`",
				},

				"timeout": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "Time in seconds after which the probe times out. Possible values are between `1` an `240`. Defaults to `1`.",
				},

				"failure_count_threshold": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `30`. Defaults to `3`.",
				},

				"success_count_threshold": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The number of consecutive successful responses required to consider this probe as successful. Possible values are between `1` and `10`. Defaults to `3`.",
				},
			},
		},
	}
}

func expandContainerAppReadinessProbe(input ContainerAppReadinessProbe) containerapps.ContainerAppProbe {
	probeType := containerapps.TypeReadiness
	result := containerapps.ContainerAppProbe{
		Type:                &probeType,
		InitialDelaySeconds: pointer.To(input.InitialDelay),
		PeriodSeconds:       pointer.To(input.Interval),
		TimeoutSeconds:      pointer.To(input.Timeout),
		FailureThreshold:    pointer.To(input.FailureThreshold),
		SuccessThreshold:    pointer.To(input.SuccessThreshold),
	}

	switch p := strings.ToUpper(input.Transport); p {
	case "HTTP", "HTTPS":
		scheme := containerapps.Scheme(p)
		result.HTTPGet = &containerapps.ContainerAppProbeHTTPGet{
			Host:   pointer.To(input.Host),
			Path:   pointer.To(input.Path),
			Port:   input.Port,
			Scheme: &scheme,
		}
		if input.Headers != nil {
			headers := make([]containerapps.ContainerAppProbeHTTPGetHTTPHeadersInlined, 0)

			for _, h := range input.Headers {
				headers = append(headers, containerapps.ContainerAppProbeHTTPGetHTTPHeadersInlined{
					Name:  h.Name,
					Value: h.Value,
				})
			}
			result.HTTPGet.HTTPHeaders = &headers
		}

	default:
		result.TcpSocket = &containerapps.ContainerAppProbeTcpSocket{
			Host: pointer.To(input.Host),
			Port: input.Port,
		}
	}

	return result
}

func flattenContainerAppReadinessProbe(input containerapps.ContainerAppProbe) []ContainerAppReadinessProbe {
	result := make([]ContainerAppReadinessProbe, 0)
	probe := ContainerAppReadinessProbe{
		InitialDelay:     pointer.From(input.InitialDelaySeconds),
		Interval:         pointer.From(input.PeriodSeconds),
		Timeout:          pointer.From(input.TimeoutSeconds),
		FailureThreshold: pointer.From(input.FailureThreshold),
		SuccessThreshold: pointer.From(input.SuccessThreshold),
	}

	if httpGet := input.HTTPGet; httpGet != nil {
		if httpGet.Scheme != nil {
			probe.Transport = string(*httpGet.Scheme)
		}
		probe.Host = pointer.From(httpGet.Host)
		probe.Port = httpGet.Port
		probe.Path = pointer.From(httpGet.Path)

		if httpGet.HTTPHeaders != nil {
			headers := make([]HttpHeader, 0)
			for _, h := range *httpGet.HTTPHeaders {
				headers = append(headers, HttpHeader{
					Name:  h.Name,
					Value: h.Value,
				})
			}
			probe.Headers = headers
		}
	}

	if tcpSocket := input.TcpSocket; tcpSocket != nil {
		probe.Transport = "TCP"
		probe.Host = pointer.From(tcpSocket.Host)
		probe.Port = tcpSocket.Port
	}

	result = append(result, probe)

	return result
}

type ContainerAppLivenessProbe struct {
	Transport              string       `tfschema:"transport"`
	Host                   string       `tfschema:"host"`
	Port                   int64        `tfschema:"port"`
	Path                   string       `tfschema:"path"`
	Headers                []HttpHeader `tfschema:"header"`
	InitialDelay           int64        `tfschema:"initial_delay"`
	Interval               int64        `tfschema:"interval_seconds"`
	Timeout                int64        `tfschema:"timeout"`
	FailureThreshold       int64        `tfschema:"failure_count_threshold"`
	TerminationGracePeriod int64        `tfschema:"termination_grace_period_seconds,removedInNextMajorVersion"`
}

func ContainerAppLivenessProbeSchema() *pluginsdk.Schema {
	schema := &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"transport": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						"TCP",
						string(containerapps.SchemeHTTP),
						string(containerapps.SchemeHTTPS),
					}, false),
					Description: "Type of probe. Possible values are `TCP`, `HTTP`, and `HTTPS`.",
				},

				"port": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
					Description:  "The port number on which to connect. Possible values are between `1` and `65535`.",
				},

				"host": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The probe hostname. Defaults to the pod IP address. Setting a value for `Host` in `headers` can be used to override this for `http` and `https` type probes.",
				},

				"path": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Computed:    true, // Note: O+C Needs to remain computed as this has a variable default and since it is part of a list we cannot diffsuppress it.
					Description: "The URI to use with the `host` for http type probes. Not valid for `TCP` type probes. Defaults to `/`.",
				},

				"header": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
								Description:  "The HTTP Header Name.",
							},

							"value": {
								Type:        pluginsdk.TypeString,
								Required:    true,
								Description: "The HTTP Header value.",
							},
						},
					},
				},

				"initial_delay": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      1,
					ValidateFunc: validation.IntBetween(0, 60),
					Description:  "The number of seconds elapsed after the container has started before the probe is initiated. Possible values are between `0` and `60`. Defaults to `1` seconds.",
				},

				"interval_seconds": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      10,
					ValidateFunc: validation.IntBetween(1, 240),
					Description:  "How often, in seconds, the probe should run. Possible values are between `1` and `240`. Defaults to `10`",
				},

				"timeout": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      1,
					ValidateFunc: validation.IntBetween(1, 240),
					Description:  "Time in seconds after which the probe times out. Possible values are between `1` an `240`. Defaults to `1`.",
				},

				"failure_count_threshold": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      3,
					ValidateFunc: validation.IntBetween(1, 30),
					Description:  "The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `30`. Defaults to `3`.",
				},
			},
		},
	}

	if !features.FivePointOh() {
		schema.Elem.(*pluginsdk.Resource).Schema["termination_grace_period_seconds"] = &pluginsdk.Schema{
			Type:        pluginsdk.TypeInt,
			Computed:    true,
			Description: "The time in seconds after the container is sent the termination signal before the process if forcibly killed.",
		}
	}

	return schema
}

func ContainerAppLivenessProbeSchemaComputed() *pluginsdk.Schema {
	schema := &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"transport": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "Type of probe. Possible values are `TCP`, `HTTP`, and `HTTPS`.",
				},

				"port": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The port number on which to connect. Possible values are between `1` and `65535`.",
				},

				"host": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The probe hostname. Defaults to the pod IP address. Setting a value for `Host` in `headers` can be used to override this for `http` and `https` type probes.",
				},

				"path": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The URI to use with the `host` for http type probes. Not valid for `TCP` type probes. Defaults to `/`.",
				},

				"header": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Type:        pluginsdk.TypeString,
								Computed:    true,
								Description: "The HTTP Header Name.",
							},

							"value": {
								Type:        pluginsdk.TypeString,
								Computed:    true,
								Description: "The HTTP Header value.",
							},
						},
					},
				},

				"initial_delay": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The number of seconds elapsed after the container has started before the probe is initiated. Possible values are between `0` and `60`. Defaults to `1` seconds.",
				},

				"interval_seconds": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "How often, in seconds, the probe should run. Possible values are between `1` and `240`. Defaults to `10`",
				},

				"timeout": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "Time in seconds after which the probe times out. Possible values are between `1` an `240`. Defaults to `1`.",
				},

				"failure_count_threshold": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `30`. Defaults to `3`.",
				},
			},
		},
	}

	if !features.FivePointOh() {
		schema.Elem.(*pluginsdk.Resource).Schema["termination_grace_period_seconds"] = &pluginsdk.Schema{
			Type:        pluginsdk.TypeInt,
			Computed:    true,
			Description: "The time in seconds after the container is sent the termination signal before the process if forcibly killed.",
		}
	}

	return schema
}

func expandContainerAppLivenessProbe(input ContainerAppLivenessProbe) containerapps.ContainerAppProbe {
	probeType := containerapps.TypeLiveness
	result := containerapps.ContainerAppProbe{
		Type:                &probeType,
		InitialDelaySeconds: pointer.To(input.InitialDelay),
		PeriodSeconds:       pointer.To(input.Interval),
		TimeoutSeconds:      pointer.To(input.Timeout),
		FailureThreshold:    pointer.To(input.FailureThreshold),
	}

	switch p := strings.ToUpper(input.Transport); p {
	case "HTTP", "HTTPS":
		scheme := containerapps.Scheme(p)
		result.HTTPGet = &containerapps.ContainerAppProbeHTTPGet{
			Host:   pointer.To(input.Host),
			Path:   pointer.To(input.Path),
			Port:   input.Port,
			Scheme: &scheme,
		}
		if input.Headers != nil {
			headers := make([]containerapps.ContainerAppProbeHTTPGetHTTPHeadersInlined, 0)

			for _, h := range input.Headers {
				headers = append(headers, containerapps.ContainerAppProbeHTTPGetHTTPHeadersInlined{
					Name:  h.Name,
					Value: h.Value,
				})
			}
			result.HTTPGet.HTTPHeaders = &headers
		}

	default:
		result.TcpSocket = &containerapps.ContainerAppProbeTcpSocket{
			Host: pointer.To(input.Host),
			Port: input.Port,
		}
	}

	return result
}

func flattenContainerAppLivenessProbe(input containerapps.ContainerAppProbe) []ContainerAppLivenessProbe {
	result := make([]ContainerAppLivenessProbe, 0)
	probe := ContainerAppLivenessProbe{
		InitialDelay:           pointer.From(input.InitialDelaySeconds),
		Interval:               pointer.From(input.PeriodSeconds),
		Timeout:                pointer.From(input.TimeoutSeconds),
		FailureThreshold:       pointer.From(input.FailureThreshold),
		TerminationGracePeriod: pointer.From(input.TerminationGracePeriodSeconds),
	}
	if httpGet := input.HTTPGet; httpGet != nil {
		if httpGet.Scheme != nil {
			probe.Transport = string(*httpGet.Scheme)
		}
		probe.Host = pointer.From(httpGet.Host)
		probe.Port = httpGet.Port
		probe.Path = pointer.From(httpGet.Path)

		if httpGet.HTTPHeaders != nil {
			headers := make([]HttpHeader, 0)
			for _, h := range *httpGet.HTTPHeaders {
				headers = append(headers, HttpHeader{
					Name:  h.Name,
					Value: h.Value,
				})
			}
			probe.Headers = headers
		}
	}

	if tcpSocket := input.TcpSocket; tcpSocket != nil {
		probe.Transport = "TCP"
		probe.Host = pointer.From(tcpSocket.Host)
		probe.Port = tcpSocket.Port
	}

	result = append(result, probe)

	return result
}

type ContainerAppStartupProbe struct {
	Transport              string       `tfschema:"transport"`
	Host                   string       `tfschema:"host"`
	Port                   int64        `tfschema:"port"`
	Path                   string       `tfschema:"path"`
	Headers                []HttpHeader `tfschema:"header"`
	InitialDelay           int64        `tfschema:"initial_delay"`
	Interval               int64        `tfschema:"interval_seconds"`
	Timeout                int64        `tfschema:"timeout"`
	FailureThreshold       int64        `tfschema:"failure_count_threshold"`
	TerminationGracePeriod int64        `tfschema:"termination_grace_period_seconds,removedInNextMajorVersion"`
}

func ContainerAppStartupProbeSchema() *pluginsdk.Schema {
	schema := &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"transport": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						"TCP",
						string(containerapps.SchemeHTTP),
						string(containerapps.SchemeHTTPS),
					}, false),
					Description: "Type of probe. Possible values are `TCP`, `HTTP`, and `HTTPS`.",
				},

				"port": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
					Description:  "The port number on which to connect. Possible values are between `1` and `65535`.",
				},

				"host": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The probe hostname. Defaults to the pod IP address. Setting a value for `Host` in `headers` can be used to override this for `http` and `https` type probes.",
				},

				"path": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Computed:    true, // Note: O+C Needs to remain computed as this has a variable default and since it is part of a list we cannot diffsuppress it.
					Description: "The URI to use with the `host` for http type probes. Not valid for `TCP` type probes. Defaults to `/`.",
				},

				"header": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
								Description:  "The HTTP Header Name.",
							},

							"value": {
								Type:        pluginsdk.TypeString,
								Required:    true,
								Description: "The HTTP Header value.",
							},
						},
					},
				},

				"initial_delay": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      0,
					ValidateFunc: validation.IntBetween(0, 60),
					Description:  "The number of seconds elapsed after the container has started before the probe is initiated. Possible values are between `0` and `60`. Defaults to `0` seconds.",
				},

				"interval_seconds": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      10,
					ValidateFunc: validation.IntBetween(1, 240),
					Description:  "How often, in seconds, the probe should run. Possible values are between `1` and `240`. Defaults to `10`",
				},

				"timeout": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      1,
					ValidateFunc: validation.IntBetween(1, 240),
					Description:  "Time in seconds after which the probe times out. Possible values are between `1` an `240`. Defaults to `1`.",
				},

				"failure_count_threshold": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      3,
					ValidateFunc: validation.IntBetween(1, 30),
					Description:  "The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `30`. Defaults to `3`.",
				},
			},
		},
	}

	if !features.FivePointOh() {
		schema.Elem.(*pluginsdk.Resource).Schema["termination_grace_period_seconds"] = &pluginsdk.Schema{
			Type:        pluginsdk.TypeInt,
			Computed:    true,
			Description: "The time in seconds after the container is sent the termination signal before the process if forcibly killed.",
		}
	}

	return schema
}

func ContainerAppStartupProbeSchemaComputed() *pluginsdk.Schema {
	schema := &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"transport": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"port": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"host": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"path": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"header": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Type:        pluginsdk.TypeString,
								Computed:    true,
								Description: "The HTTP Header Name.",
							},

							"value": {
								Type:        pluginsdk.TypeString,
								Computed:    true,
								Description: "The HTTP Header value.",
							},
						},
					},
				},

				"initial_delay": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The number of seconds elapsed after the container has started before the probe is initiated. Possible values are between `0` and `60`. Defaults to `0` seconds.",
				},

				"interval_seconds": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "How often, in seconds, the probe should run. Possible values are between `1` and `240`. Defaults to `10`",
				},

				"timeout": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "Time in seconds after which the probe times out. Possible values are between `1` an `240`. Defaults to `1`.",
				},

				"failure_count_threshold": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `30`. Defaults to `3`.",
				},
			},
		},
	}

	if !features.FivePointOh() {
		schema.Elem.(*pluginsdk.Resource).Schema["termination_grace_period_seconds"] = &pluginsdk.Schema{
			Type:        pluginsdk.TypeInt,
			Computed:    true,
			Description: "The time in seconds after the container is sent the termination signal before the process if forcibly killed.",
		}
	}

	return schema
}

func expandContainerAppStartupProbe(input ContainerAppStartupProbe) containerapps.ContainerAppProbe {
	probeType := containerapps.TypeStartup
	result := containerapps.ContainerAppProbe{
		Type:                &probeType,
		InitialDelaySeconds: pointer.To(input.InitialDelay),
		PeriodSeconds:       pointer.To(input.Interval),
		TimeoutSeconds:      pointer.To(input.Timeout),
		FailureThreshold:    pointer.To(input.FailureThreshold),
	}

	switch p := strings.ToUpper(input.Transport); p {
	case "HTTP", "HTTPS":
		scheme := containerapps.Scheme(p)
		result.HTTPGet = &containerapps.ContainerAppProbeHTTPGet{
			Host:   pointer.To(input.Host),
			Path:   pointer.To(input.Path),
			Port:   input.Port,
			Scheme: &scheme,
		}
		if input.Headers != nil {
			headers := make([]containerapps.ContainerAppProbeHTTPGetHTTPHeadersInlined, 0)

			for _, h := range input.Headers {
				headers = append(headers, containerapps.ContainerAppProbeHTTPGetHTTPHeadersInlined{
					Name:  h.Name,
					Value: h.Value,
				})
			}
			result.HTTPGet.HTTPHeaders = &headers
		}

	default:
		result.TcpSocket = &containerapps.ContainerAppProbeTcpSocket{
			Host: pointer.To(input.Host),
			Port: input.Port,
		}
	}

	return result
}

func flattenContainerAppStartupProbe(input containerapps.ContainerAppProbe) []ContainerAppStartupProbe {
	result := make([]ContainerAppStartupProbe, 0)
	probe := ContainerAppStartupProbe{
		InitialDelay:           pointer.From(input.InitialDelaySeconds),
		Interval:               pointer.From(input.PeriodSeconds),
		Timeout:                pointer.From(input.TimeoutSeconds),
		FailureThreshold:       pointer.From(input.FailureThreshold),
		TerminationGracePeriod: pointer.From(input.TerminationGracePeriodSeconds),
	}

	if httpGet := input.HTTPGet; httpGet != nil {
		if httpGet.Scheme != nil {
			probe.Transport = string(*httpGet.Scheme)
		}
		probe.Host = pointer.From(httpGet.Host)
		probe.Port = httpGet.Port
		probe.Path = pointer.From(httpGet.Path)

		if httpGet.HTTPHeaders != nil {
			headers := make([]HttpHeader, 0)
			for _, h := range *httpGet.HTTPHeaders {
				headers = append(headers, HttpHeader{
					Name:  h.Name,
					Value: h.Value,
				})
			}
			probe.Headers = headers
		}
	}

	if tcpSocket := input.TcpSocket; tcpSocket != nil {
		probe.Transport = "TCP"
		probe.Host = pointer.From(tcpSocket.Host)
		probe.Port = tcpSocket.Port
	}

	result = append(result, probe)

	return result
}

type HttpHeader struct {
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}

func expandContainerProbes(input Container) *[]containerapps.ContainerAppProbe {
	probes := make([]containerapps.ContainerAppProbe, 0)

	if len(input.LivenessProbe) == 1 {
		probes = append(probes, expandContainerAppLivenessProbe(input.LivenessProbe[0]))
	}

	if len(input.ReadinessProbe) == 1 {
		probes = append(probes, expandContainerAppReadinessProbe(input.ReadinessProbe[0]))
	}

	if len(input.StartupProbe) == 1 {
		probes = append(probes, expandContainerAppStartupProbe(input.StartupProbe[0]))
	}

	return &probes
}

type Secret struct {
	Identity         string `tfschema:"identity"`
	KeyVaultSecretId string `tfschema:"key_vault_secret_id"`
	Name             string `tfschema:"name"`
	Value            string `tfschema:"value"`
}

func SecretsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:      pluginsdk.TypeSet,
		Optional:  true,
		Sensitive: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"identity": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.Any(
						commonids.ValidateUserAssignedIdentityID,
						validation.StringInSlice([]string{"System"}, false),
					),
					Description: "The identity to use for accessing key vault reference.",
				},

				"key_vault_secret_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
					Description:  "The Key Vault Secret ID. Could be either one of `id` or `versionless_id`.",
				},

				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validate.SecretName,
					Description:  "The secret name.",
				},

				"value": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Sensitive:   true,
					Description: "The value for this secret.",
				},
			},
		},
	}
}

func SecretsDataSourceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:      pluginsdk.TypeSet,
		Computed:  true,
		Sensitive: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"identity": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The identity to use for accessing key vault reference.",
				},

				"key_vault_secret_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The id of the key vault secret.",
				},

				"name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The secret name.",
				},

				"value": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Sensitive:   true,
					Description: "The value for this secret.",
				},
			},
		},
	}
}

func ExpandContainerSecrets(input []Secret) (*[]containerapps.Secret, error) {
	if len(input) == 0 {
		return nil, nil
	}

	result := make([]containerapps.Secret, 0)

	for _, v := range input {
		result = append(result, containerapps.Secret{
			Identity:    pointer.To(v.Identity),
			KeyVaultURL: pointer.To(v.KeyVaultSecretId),
			Name:        pointer.To(v.Name),
			Value:       pointer.To(v.Value),
		})
	}

	return &result, nil
}

func UnpackContainerSecretsCollection(input *containerapps.SecretsCollection) *[]containerapps.Secret {
	if input == nil || len(input.Value) == 0 {
		return nil
	}

	result := make([]containerapps.Secret, 0)
	for _, v := range input.Value {
		result = append(result, containerapps.Secret(v))
	}

	return &result
}

func UnpackContainerDaprSecretsCollection(input *daprcomponents.DaprSecretsCollection) *[]daprcomponents.Secret {
	if input == nil || len(input.Value) == 0 {
		return nil
	}

	result := make([]daprcomponents.Secret, 0)
	for _, v := range input.Value {
		result = append(result, daprcomponents.Secret{
			Name:  v.Name,
			Value: v.Value,
		})
	}

	return &result
}

type DaprSecret struct {
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}

func ExpandDaprSecrets(input []DaprSecret) *[]daprcomponents.Secret {
	if len(input) == 0 {
		return nil
	}

	result := make([]daprcomponents.Secret, 0)

	for _, v := range input {
		result = append(result, daprcomponents.Secret{
			Name:  pointer.To(v.Name),
			Value: pointer.To(v.Value),
		})
	}

	return &result
}

func FlattenContainerAppSecrets(input *containerapps.SecretsCollection) []Secret {
	if input == nil || input.Value == nil {
		return []Secret{}
	}
	result := make([]Secret, 0)
	for _, v := range input.Value {
		secret := Secret{
			Identity:         pointer.From(v.Identity),
			KeyVaultSecretId: pointer.From(v.KeyVaultURL),
			Name:             pointer.From(v.Name),
		}
		if v.KeyVaultURL == nil {
			secret.Value = pointer.From(v.Value)
		}
		result = append(result, secret)
	}

	return result
}

func FlattenContainerAppDaprSecrets(input *daprcomponents.DaprSecretsCollection) []DaprSecret {
	if input == nil || input.Value == nil {
		return []DaprSecret{}
	}
	result := make([]DaprSecret, 0)
	for _, v := range input.Value {
		result = append(result, DaprSecret{
			Name:  pointer.From(v.Name),
			Value: pointer.From(v.Value),
		})
	}

	return result
}

func ContainerAppProbesRemoved(metadata sdk.ResourceMetaData) bool {
	var hasLiveness, hasReadiness, hasStartup bool

	if metadata.ResourceData.HasChange("template.0.container.0.liveness_probe") {
		_, newLivenessRaw := metadata.ResourceData.GetChange("template.0.container.0.liveness_probe")
		if newLiveness, ok := newLivenessRaw.([]interface{}); ok && len(newLiveness) > 0 {
			hasLiveness = true
		}
	}

	if metadata.ResourceData.HasChange("template.0.container.0.readiness_probe") {
		_, newReadinessRaw := metadata.ResourceData.GetChange("template.0.container.0.readiness_probe")
		if newReadiness, ok := newReadinessRaw.([]interface{}); ok && len(newReadiness) > 0 {
			hasReadiness = true
		}
	}

	if metadata.ResourceData.HasChange("template.0.container.0.startup_probe") {
		_, newStartupRaw := metadata.ResourceData.GetChange("template.0.container.0.startup_probe")
		if newStartup, ok := newStartupRaw.([]interface{}); ok && len(newStartup) > 0 {
			hasStartup = true
		}
	}

	return !(hasLiveness || hasReadiness || hasStartup)
}

type AzureQueueScaleRule struct {
	Name            string                    `tfschema:"name"`
	QueueLength     int64                     `tfschema:"queue_length"`
	QueueName       string                    `tfschema:"queue_name"`
	Authentications []ScaleRuleAuthentication `tfschema:"authentication"`
}

func AzureQueueScaleRuleSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"queue_length": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(1),
				},

				"queue_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"authentication": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MinItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"secret_name": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validate.SecretName,
							},

							"trigger_parameter": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
		},
	}
}

func AzureQueueScaleRuleSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"queue_length": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"queue_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"authentication": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"secret_name": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"trigger_parameter": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

type CustomScaleRule struct {
	Name            string                    `tfschema:"name"`
	Metadata        map[string]string         `tfschema:"metadata"`
	CustomRuleType  string                    `tfschema:"custom_rule_type"`
	Authentications []ScaleRuleAuthentication `tfschema:"authentication"`
}

func CustomScaleRuleSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validate.LowerCaseAlphaNumericWithHyphensAndPeriods,
				},

				"metadata": {
					Type:     pluginsdk.TypeMap,
					Required: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"custom_rule_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						"activemq", "artemis-queue", "kafka", "pulsar", "aws-cloudwatch",
						"aws-dynamodb", "aws-dynamodb-streams", "aws-kinesis-stream", "aws-sqs-queue",
						"azure-app-insights", "azure-blob", "azure-data-explorer", "azure-eventhub",
						"azure-log-analytics", "azure-monitor", "azure-pipelines", "azure-servicebus",
						"azure-queue", "cassandra", "cpu", "cron", "datadog", "elasticsearch", "external",
						"external-push", "gcp-stackdriver", "gcp-storage", "gcp-pubsub", "graphite", "http",
						"huawei-cloudeye", "ibmmq", "influxdb", "kubernetes-workload", "liiklus", "memory",
						"metrics-api", "mongodb", "mssql", "mysql", "nats-jetstream", "stan", "tcp", "new-relic",
						"openstack-metric", "openstack-swift", "postgresql", "predictkube", "prometheus",
						"rabbitmq", "redis", "redis-cluster", "redis-sentinel", "redis-streams",
						"redis-cluster-streams", "redis-sentinel-streams", "selenium-grid",
						"solace-event-queue", "github-runner",
					}, false), // Note - this can be any KEDA compatible source in a user's environment
				},

				"authentication": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MinItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"secret_name": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validate.SecretName,
							},

							"trigger_parameter": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
		},
	}
}

func CustomScaleRuleSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"metadata": {
					Type:     pluginsdk.TypeMap,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"custom_rule_type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"authentication": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"secret_name": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"trigger_parameter": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

type HTTPScaleRule struct {
	Name               string                    `tfschema:"name"`
	ConcurrentRequests string                    `tfschema:"concurrent_requests"`
	Authentications    []ScaleRuleAuthentication `tfschema:"authentication"`
}

func HTTPScaleRuleSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"concurrent_requests": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validate.ContainerAppScaleRuleConcurrentRequests,
				},

				"authentication": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MinItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"secret_name": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validate.SecretName,
							},

							"trigger_parameter": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
		},
	}
}

func HTTPScaleRuleSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"concurrent_requests": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"authentication": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"secret_name": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"trigger_parameter": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

type TCPScaleRule struct {
	Name               string                    `tfschema:"name"`
	ConcurrentRequests string                    `tfschema:"concurrent_requests"`
	Authentications    []ScaleRuleAuthentication `tfschema:"authentication"`
}

func TCPScaleRuleSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"concurrent_requests": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validate.ContainerAppScaleRuleConcurrentRequests,
				},

				"authentication": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MinItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"secret_name": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validate.SecretName,
							},

							"trigger_parameter": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
		},
	}
}

func TCPScaleRuleSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"concurrent_requests": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"authentication": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"secret_name": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"trigger_parameter": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

type ScaleRuleAuthentication struct {
	SecretRef    string `tfschema:"secret_name"`
	TriggerParam string `tfschema:"trigger_parameter"`
}

func (c *ContainerTemplate) expandContainerAppScaleRules() []containerapps.ScaleRule {
	if len(c.AzureQueueScaleRules) == 0 && len(c.CustomScaleRules) == 0 && len(c.HTTPScaleRules) == 0 && len(c.TCPScaleRules) == 0 {
		return nil
	}
	result := make([]containerapps.ScaleRule, 0)
	for _, v := range c.AzureQueueScaleRules {
		r := containerapps.ScaleRule{
			Name: pointer.To(v.Name),
			AzureQueue: &containerapps.QueueScaleRule{
				QueueLength: pointer.To(v.QueueLength),
				QueueName:   pointer.To(v.QueueName),
			},
		}

		auths := make([]containerapps.ScaleRuleAuth, 0)
		for _, a := range v.Authentications {
			auth := containerapps.ScaleRuleAuth{
				TriggerParameter: pointer.To(a.TriggerParam),
				SecretRef:        pointer.To(a.SecretRef),
			}
			auths = append(auths, auth)
		}

		r.AzureQueue.Auth = pointer.To(auths)

		result = append(result, r)
	}

	for _, v := range c.CustomScaleRules {
		r := containerapps.ScaleRule{
			Name: pointer.To(v.Name),
			Custom: &containerapps.CustomScaleRule{
				Metadata: pointer.To(v.Metadata),
				Type:     pointer.To(v.CustomRuleType),
			},
		}

		auths := make([]containerapps.ScaleRuleAuth, 0)
		for _, a := range v.Authentications {
			auth := containerapps.ScaleRuleAuth{
				TriggerParameter: pointer.To(a.TriggerParam),
				SecretRef:        pointer.To(a.SecretRef),
			}
			auths = append(auths, auth)
		}

		r.Custom.Auth = pointer.To(auths)

		result = append(result, r)
	}

	for _, v := range c.HTTPScaleRules {
		metaData := make(map[string]string, 0)
		metaData["concurrentRequests"] = v.ConcurrentRequests
		r := containerapps.ScaleRule{
			Name: pointer.To(v.Name),
			HTTP: &containerapps.HTTPScaleRule{
				Metadata: pointer.To(metaData),
			},
		}

		auths := make([]containerapps.ScaleRuleAuth, 0)
		for _, a := range v.Authentications {
			auth := containerapps.ScaleRuleAuth{
				TriggerParameter: pointer.To(a.TriggerParam),
				SecretRef:        pointer.To(a.SecretRef),
			}
			auths = append(auths, auth)
		}

		r.HTTP.Auth = pointer.To(auths)

		result = append(result, r)
	}

	for _, v := range c.TCPScaleRules {
		metaData := make(map[string]string, 0)
		metaData["concurrentRequests"] = v.ConcurrentRequests
		r := containerapps.ScaleRule{
			Name: pointer.To(v.Name),
			Tcp: &containerapps.TcpScaleRule{
				Metadata: pointer.To(metaData),
			},
		}

		auths := make([]containerapps.ScaleRuleAuth, 0)
		for _, a := range v.Authentications {
			auth := containerapps.ScaleRuleAuth{
				TriggerParameter: pointer.To(a.TriggerParam),
				SecretRef:        pointer.To(a.SecretRef),
			}
			auths = append(auths, auth)
		}

		r.Tcp.Auth = pointer.To(auths)

		result = append(result, r)
	}

	return result
}

func (c *ContainerTemplate) flattenContainerAppScaleRules(input *[]containerapps.ScaleRule) {
	if input != nil && len(*input) != 0 {
		rules := *input
		azureQueueScaleRules := make([]AzureQueueScaleRule, 0)
		customScaleRules := make([]CustomScaleRule, 0)
		httpScaleRules := make([]HTTPScaleRule, 0)
		tcpScaleRules := make([]TCPScaleRule, 0)
		for _, v := range rules {
			if q := v.AzureQueue; q != nil {
				rule := AzureQueueScaleRule{
					Name:        pointer.From(v.Name),
					QueueLength: pointer.From(q.QueueLength),
					QueueName:   pointer.From(q.QueueName),
				}

				authentications := make([]ScaleRuleAuthentication, 0)
				if auths := q.Auth; auths != nil {
					for _, a := range *auths {
						authentications = append(authentications, ScaleRuleAuthentication{
							SecretRef:    pointer.From(a.SecretRef),
							TriggerParam: pointer.From(a.TriggerParameter),
						})
					}
				}

				rule.Authentications = authentications

				azureQueueScaleRules = append(azureQueueScaleRules, rule)
				continue
			}

			if r := v.Custom; r != nil {
				rule := CustomScaleRule{
					Name:           pointer.From(v.Name),
					Metadata:       pointer.From(r.Metadata),
					CustomRuleType: pointer.From(r.Type),
				}

				authentications := make([]ScaleRuleAuthentication, 0)
				if auths := r.Auth; auths != nil {
					for _, a := range *auths {
						authentications = append(authentications, ScaleRuleAuthentication{
							SecretRef:    pointer.From(a.SecretRef),
							TriggerParam: pointer.From(a.TriggerParameter),
						})
					}
				}
				rule.Authentications = authentications

				customScaleRules = append(customScaleRules, rule)
				continue
			}

			if r := v.HTTP; r != nil {
				metaData := pointer.From(r.Metadata)
				concurrentReqs := ""

				if m, ok := metaData["concurrentRequests"]; ok {
					concurrentReqs = m
				}

				rule := HTTPScaleRule{
					Name:               pointer.From(v.Name),
					ConcurrentRequests: concurrentReqs,
				}

				authentications := make([]ScaleRuleAuthentication, 0)
				if auths := r.Auth; auths != nil {
					for _, a := range *auths {
						authentications = append(authentications, ScaleRuleAuthentication{
							SecretRef:    pointer.From(a.SecretRef),
							TriggerParam: pointer.From(a.TriggerParameter),
						})
					}
				}

				rule.Authentications = authentications

				httpScaleRules = append(httpScaleRules, rule)
				continue
			}

			if r := v.Tcp; r != nil {
				metaData := pointer.From(r.Metadata)
				concurrentReqs := ""

				if m, ok := metaData["concurrentRequests"]; ok {
					concurrentReqs = m
				}

				rule := TCPScaleRule{
					Name:               pointer.From(v.Name),
					ConcurrentRequests: concurrentReqs,
				}

				authentications := make([]ScaleRuleAuthentication, 0)
				if auths := r.Auth; auths != nil {
					for _, a := range *auths {
						authentications = append(authentications, ScaleRuleAuthentication{
							SecretRef:    pointer.From(a.SecretRef),
							TriggerParam: pointer.From(a.TriggerParameter),
						})
					}
				}
				rule.Authentications = authentications

				tcpScaleRules = append(tcpScaleRules, rule)
				continue
			}
		}

		c.AzureQueueScaleRules = azureQueueScaleRules
		c.CustomScaleRules = customScaleRules
		c.HTTPScaleRules = httpScaleRules
		c.TCPScaleRules = tcpScaleRules
	}
}
