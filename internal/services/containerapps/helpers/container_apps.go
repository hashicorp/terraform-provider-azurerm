// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/containerapps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/daprcomponents"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
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
	AllowInsecure  bool            `tfschema:"allow_insecure_connections"`
	CustomDomains  []CustomDomain  `tfschema:"custom_domain"`
	IsExternal     bool            `tfschema:"external_enabled"`
	FQDN           string          `tfschema:"fqdn"`
	TargetPort     int             `tfschema:"target_port"`
	TrafficWeights []TrafficWeight `tfschema:"traffic_weight"`
	Transport      string          `tfschema:"transport"`
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

				"custom_domain": ContainerAppIngressCustomDomainSchema(),

				"external_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Is this an external Ingress.",
				},

				"fqdn": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The FQDN of the ingress.",
				},

				"target_port": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
					Description:  "The target port on the container for the Ingress traffic.",
				},

				"traffic_weight": ContainerAppIngressTrafficWeight(),

				"transport": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(containerapps.IngressTransportMethodAuto),
					ValidateFunc: validation.StringInSlice(containerapps.PossibleValuesForIngressTransportMethod(), false),
					Description:  "The transport method for the Ingress. Possible values include `auto`, `http`, and `http2`. Defaults to `auto`",
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

				"target_port": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The target port on the container for the Ingress traffic.",
				},

				"traffic_weight": ContainerAppIngressTrafficWeightComputed(),

				"transport": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The transport method for the Ingress. Possible values include `auto`, `http`, and `http2`. Defaults to `auto`",
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
		AllowInsecure: pointer.To(ingress.AllowInsecure),
		CustomDomains: expandContainerAppIngressCustomDomain(ingress.CustomDomains),
		External:      pointer.To(ingress.IsExternal),
		Fqdn:          pointer.To(ingress.FQDN),
		TargetPort:    pointer.To(int64(ingress.TargetPort)),
		Traffic:       expandContainerAppIngressTraffic(ingress.TrafficWeights, appName),
	}
	transport := containerapps.IngressTransportMethod(ingress.Transport)
	result.Transport = &transport

	return result
}

func FlattenContainerAppIngress(input *containerapps.Ingress, appName string) []Ingress {
	if input == nil {
		return []Ingress{}
	}

	ingress := *input
	result := Ingress{
		AllowInsecure:  pointer.From(ingress.AllowInsecure),
		CustomDomains:  flattenContainerAppIngressCustomDomain(ingress.CustomDomains),
		IsExternal:     pointer.From(ingress.External),
		FQDN:           pointer.From(ingress.Fqdn),
		TargetPort:     int(pointer.From(ingress.TargetPort)),
		TrafficWeights: flattenContainerAppIngressTraffic(ingress.Traffic, appName),
	}

	if ingress.Transport != nil {
		result.Transport = strings.ToLower(string(*ingress.Transport))
	}

	return []Ingress{result}
}

type CustomDomain struct {
	CertBinding   string `tfschema:"certificate_binding_type"`
	CertificateId string `tfschema:"certificate_id"`
	Name          string `tfschema:"name"`
}

func ContainerAppIngressCustomDomainSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"certificate_binding_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      containerapps.BindingTypeDisabled,
					ValidateFunc: validation.StringInSlice(containerapps.PossibleValuesForBindingType(), false),
					Description:  "The Binding type. Possible values include `Disabled` and `SniEnabled`. Defaults to `Disabled`",
				},

				"certificate_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: managedenvironments.ValidateCertificateID,
				},

				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The hostname of the Certificate. Must be the CN or a named SAN in the certificate.",
				},
			},
		},
	}
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

type TrafficWeight struct {
	Label          string `tfschema:"label"`
	LatestRevision bool   `tfschema:"latest_revision"`
	RevisionSuffix string `tfschema:"revision_suffix"`
	Weight         int    `tfschema:"percentage"`
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
			Weight:         pointer.To(int64(v.Weight)),
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
			Weight:         int(pointer.From(v.Weight)),
		})
	}

	return result
}

type Dapr struct {
	AppId       string `tfschema:"app_id"`
	AppPort     int    `tfschema:"app_port"`
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
		AppPort:     pointer.To(int64(dapr.AppPort)),
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
		AppPort: int(pointer.From(input.AppPort)),
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
	Containers  []Container       `tfschema:"container"`
	Suffix      string            `tfschema:"revision_suffix"`
	MinReplicas int               `tfschema:"min_replicas"`
	MaxReplicas int               `tfschema:"max_replicas"`
	Volumes     []ContainerVolume `tfschema:"volume"`
}

func ContainerTemplateSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Required: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"container": ContainerAppContainerSchema(),

				"min_replicas": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
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

				"volume": ContainerVolumeSchema(),

				"revision_suffix": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "The suffix for the revision. This value must be unique for the lifetime of the Resource. If omitted the service will use a hash function to create one.",
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

				"volume": ContainerVolumeSchema(),

				"revision_suffix": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The suffix for the revision. This value must be unique for the lifetime of the Resource. If omitted the service will use a hash function to create one.",
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
		Containers: expandContainerAppContainers(config.Containers),
		Volumes:    expandContainerAppVolumes(config.Volumes),
	}

	if config.MaxReplicas != 0 {
		if template.Scale == nil {
			template.Scale = &containerapps.Scale{}
		}
		template.Scale.MaxReplicas = pointer.To(int64(config.MaxReplicas))
	}

	if config.MinReplicas != 0 {
		if template.Scale == nil {
			template.Scale = &containerapps.Scale{}
		}
		template.Scale.MinReplicas = pointer.To(int64(config.MinReplicas))
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
		Containers: flattenContainerAppContainers(input.Containers),
		Suffix:     pointer.From(input.RevisionSuffix),
		Volumes:    flattenContainerAppVolumes(input.Volumes),
	}

	if scale := input.Scale; scale != nil {
		result.MaxReplicas = int(pointer.From(scale.MaxReplicas))
		result.MinReplicas = int(pointer.From(scale.MinReplicas))
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
					ValidateFunc: validate.ContainerCpu,
					Description:  "The amount of vCPU to allocate to the container. Possible values include `0.25`, `0.5`, `0.75`, `1.0`, `1.25`, `1.5`, `1.75`, and `2.0`. **NOTE:** `cpu` and `memory` must be specified in `0.25'/'0.5Gi` combination increments. e.g. `1.0` / `2.0` or `0.5` / `1.0`",
				},

				"memory": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						"0.5Gi",
						"1Gi",
						"1.5Gi",
						"2Gi",
						"2.5Gi",
						"3Gi",
						"3.5Gi",
						"4Gi",
					}, false),
					Description: "The amount of memory to allocate to the container. Possible values include `0.5Gi`, `1.0Gi`, `1.5Gi`, `2.0Gi`, `2.5Gi`, `3.0Gi`, `3.5Gi`, and `4.0Gi`. **NOTE:** `cpu` and `memory` must be specified in `0.25'/'0.5Gi` combination increments. e.g. `1.25` / `2.5Gi` or `0.75` / `1.5Gi`",
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
			container.Args = &v.Args
		}
		if len(v.Command) != 0 {
			container.Command = &v.Command
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
	Name        string `tfschema:"name"`
	StorageName string `tfschema:"storage_name"`
	StorageType string `tfschema:"storage_type"`
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

		result = append(result, containerVolume)
	}

	return result
}

type ContainerVolumeMount struct {
	Name string `tfschema:"name"`
	Path string `tfschema:"path"`
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
			Name: pointer.From(v.VolumeName),
			Path: pointer.From(v.MountPath),
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

func expandContainerEnvVar(input Container) *[]containerapps.EnvironmentVar {
	envs := make([]containerapps.EnvironmentVar, 0)
	if input.Env == nil || len(input.Env) == 0 {
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
	Port             int          `tfschema:"port"`
	Path             string       `tfschema:"path"`
	Headers          []HttpHeader `tfschema:"header"`
	Interval         int          `tfschema:"interval_seconds"`
	Timeout          int          `tfschema:"timeout"`
	FailureThreshold int          `tfschema:"failure_count_threshold"`
	SuccessThreshold int          `tfschema:"success_count_threshold"`
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
					Computed:    true,
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
					ValidateFunc: validation.IntBetween(1, 10),
					Description:  "The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.",
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
					Description: "The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.",
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
		Type:             &probeType,
		PeriodSeconds:    pointer.To(int64(input.Interval)),
		TimeoutSeconds:   pointer.To(int64(input.Timeout)),
		FailureThreshold: pointer.To(int64(input.FailureThreshold)),
		SuccessThreshold: pointer.To(int64(input.SuccessThreshold)),
	}

	switch p := strings.ToUpper(input.Transport); p {
	case "HTTP", "HTTPS":
		scheme := containerapps.Scheme(p)
		result.HTTPGet = &containerapps.ContainerAppProbeHTTPGet{
			Host:   pointer.To(input.Host),
			Path:   pointer.To(input.Path),
			Port:   int64(input.Port),
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
			Port: int64(input.Port),
		}
	}

	return result
}

func flattenContainerAppReadinessProbe(input containerapps.ContainerAppProbe) []ContainerAppReadinessProbe {
	result := make([]ContainerAppReadinessProbe, 0)
	probe := ContainerAppReadinessProbe{
		Interval:         int(pointer.From(input.PeriodSeconds)),
		Timeout:          int(pointer.From(input.TimeoutSeconds)),
		FailureThreshold: int(pointer.From(input.FailureThreshold)),
		SuccessThreshold: int(pointer.From(input.SuccessThreshold)),
	}

	if httpGet := input.HTTPGet; httpGet != nil {
		if httpGet.Scheme != nil {
			probe.Transport = string(*httpGet.Scheme)
		}
		probe.Host = pointer.From(httpGet.Host)
		probe.Port = int(httpGet.Port)
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
		probe.Port = int(tcpSocket.Port)
	}

	result = append(result, probe)

	return result
}

type ContainerAppLivenessProbe struct {
	Transport              string       `tfschema:"transport"`
	Host                   string       `tfschema:"host"`
	Port                   int          `tfschema:"port"`
	Path                   string       `tfschema:"path"`
	Headers                []HttpHeader `tfschema:"header"`
	InitialDelay           int          `tfschema:"initial_delay"`
	Interval               int          `tfschema:"interval_seconds"`
	Timeout                int          `tfschema:"timeout"`
	FailureThreshold       int          `tfschema:"failure_count_threshold"`
	TerminationGracePeriod int          `tfschema:"termination_grace_period_seconds"`
}

func ContainerAppLivenessProbeSchema() *pluginsdk.Schema {
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
					Computed:    true,
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
					ValidateFunc: validation.IntBetween(1, 60),
					Description:  "The time in seconds to wait after the container has started before the probe is started.",
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
					ValidateFunc: validation.IntBetween(1, 10),
					Description:  "The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.",
				},

				"termination_grace_period_seconds": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The time in seconds after the container is sent the termination signal before the process if forcibly killed.",
				},
			},
		},
	}
}

func ContainerAppLivenessProbeSchemaComputed() *pluginsdk.Schema {
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
					Description: "The time in seconds to wait after the container has started before the probe is started.",
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
					Description: "The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.",
				},

				"termination_grace_period_seconds": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The time in seconds after the container is sent the termination signal before the process if forcibly killed.",
				},
			},
		},
	}
}

func expandContainerAppLivenessProbe(input ContainerAppLivenessProbe) containerapps.ContainerAppProbe {
	probeType := containerapps.TypeLiveness
	result := containerapps.ContainerAppProbe{
		Type:                &probeType,
		InitialDelaySeconds: pointer.To(int64(input.InitialDelay)),
		PeriodSeconds:       pointer.To(int64(input.Interval)),
		TimeoutSeconds:      pointer.To(int64(input.Timeout)),
		FailureThreshold:    pointer.To(int64(input.FailureThreshold)),
	}

	switch p := strings.ToUpper(input.Transport); p {
	case "HTTP", "HTTPS":
		scheme := containerapps.Scheme(p)
		result.HTTPGet = &containerapps.ContainerAppProbeHTTPGet{
			Host:   pointer.To(input.Host),
			Path:   pointer.To(input.Path),
			Port:   int64(input.Port),
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
			Port: int64(input.Port),
		}
	}

	return result
}

func flattenContainerAppLivenessProbe(input containerapps.ContainerAppProbe) []ContainerAppLivenessProbe {
	result := make([]ContainerAppLivenessProbe, 0)
	probe := ContainerAppLivenessProbe{
		InitialDelay:           int(pointer.From(input.InitialDelaySeconds)),
		Interval:               int(pointer.From(input.PeriodSeconds)),
		Timeout:                int(pointer.From(input.TimeoutSeconds)),
		FailureThreshold:       int(pointer.From(input.FailureThreshold)),
		TerminationGracePeriod: int(pointer.From(input.TerminationGracePeriodSeconds)),
	}
	if httpGet := input.HTTPGet; httpGet != nil {
		if httpGet.Scheme != nil {
			probe.Transport = string(*httpGet.Scheme)
		}
		probe.Host = pointer.From(httpGet.Host)
		probe.Port = int(httpGet.Port)
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
		probe.Port = int(tcpSocket.Port)
	}

	result = append(result, probe)

	return result
}

type ContainerAppStartupProbe struct {
	Transport              string       `tfschema:"transport"`
	Host                   string       `tfschema:"host"`
	Port                   int          `tfschema:"port"`
	Path                   string       `tfschema:"path"`
	Headers                []HttpHeader `tfschema:"header"`
	Interval               int          `tfschema:"interval_seconds"`
	Timeout                int          `tfschema:"timeout"`
	FailureThreshold       int          `tfschema:"failure_count_threshold"`
	TerminationGracePeriod int          `tfschema:"termination_grace_period_seconds"`
}

func ContainerAppStartupProbeSchema() *pluginsdk.Schema {
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
					Computed:    true,
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
					ValidateFunc: validation.IntBetween(1, 10),
					Description:  "The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.",
				},

				"termination_grace_period_seconds": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The time in seconds after the container is sent the termination signal before the process if forcibly killed.",
				},
			},
		},
	}
}

func ContainerAppStartupProbeSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
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
					Description: "The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.",
				},

				"termination_grace_period_seconds": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The time in seconds after the container is sent the termination signal before the process if forcibly killed.",
				},
			},
		},
	}
}

func expandContainerAppStartupProbe(input ContainerAppStartupProbe) containerapps.ContainerAppProbe {
	probeType := containerapps.TypeStartup
	result := containerapps.ContainerAppProbe{
		Type:             &probeType,
		PeriodSeconds:    pointer.To(int64(input.Interval)),
		TimeoutSeconds:   pointer.To(int64(input.Timeout)),
		FailureThreshold: pointer.To(int64(input.FailureThreshold)),
	}

	switch p := strings.ToUpper(input.Transport); p {
	case "HTTP", "HTTPS":
		scheme := containerapps.Scheme(p)
		result.HTTPGet = &containerapps.ContainerAppProbeHTTPGet{
			Host:   pointer.To(input.Host),
			Path:   pointer.To(input.Path),
			Port:   int64(input.Port),
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
			Port: int64(input.Port),
		}
	}

	return result
}

func flattenContainerAppStartupProbe(input containerapps.ContainerAppProbe) []ContainerAppStartupProbe {
	result := make([]ContainerAppStartupProbe, 0)
	probe := ContainerAppStartupProbe{
		Interval:               int(pointer.From(input.PeriodSeconds)),
		Timeout:                int(pointer.From(input.TimeoutSeconds)),
		FailureThreshold:       int(pointer.From(input.FailureThreshold)),
		TerminationGracePeriod: int(pointer.From(input.TerminationGracePeriodSeconds)),
	}

	if httpGet := input.HTTPGet; httpGet != nil {
		if httpGet.Scheme != nil {
			probe.Transport = string(*httpGet.Scheme)
		}
		probe.Host = pointer.From(httpGet.Host)
		probe.Port = int(httpGet.Port)
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
		probe.Port = int(tcpSocket.Port)
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
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}

func SecretsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:      pluginsdk.TypeSet,
		Optional:  true,
		Sensitive: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validate.SecretName,
					Sensitive:    true,
					Description:  "The Secret name.",
				},

				"value": {
					Type:        pluginsdk.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The value for this secret.",
				},
			},
		},
	}
}

func SecretsDataSourceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:      pluginsdk.TypeList,
		Computed:  true,
		Sensitive: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The Secret name.",
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

func ExpandContainerSecrets(input []Secret) *[]containerapps.Secret {
	if len(input) == 0 {
		return nil
	}

	result := make([]containerapps.Secret, 0)

	for _, v := range input {
		result = append(result, containerapps.Secret{
			Name:  pointer.To(v.Name),
			Value: pointer.To(v.Value),
		})
	}

	return &result
}

func ExpandFormerContainerSecrets(metadata sdk.ResourceMetaData) *[]containerapps.Secret {
	secretsRaw, _ := metadata.ResourceData.GetChange("secret")
	result := make([]containerapps.Secret, 0)
	if secrets, ok := secretsRaw.([]interface{}); ok {

		for _, secret := range secrets {
			if v, ok := secret.(map[string]interface{}); ok {
				result = append(result, containerapps.Secret{
					Name:  pointer.To(v["name"].(string)),
					Value: pointer.To(v["value"].(string)),
				})
			}
		}
	}

	return &result
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
			// TODO: add support for Identity & KeyVaultUrl
			Name:  v.Name,
			Value: v.Value,
		})
	}

	return &result
}

func ExpandDaprSecrets(input []Secret) *[]daprcomponents.Secret {
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

func FlattenSecrets(input []interface{}) []Secret {
	secrets := make([]Secret, 0)
	for _, s := range input {
		secret := s.(map[string]interface{})
		name, ok := secret["name"].(string)
		if !ok {
			continue
		}
		value := ""
		if val, ok := secret["value"].(string); ok {
			value = val
		}
		secrets = append(secrets, Secret{
			Name:  name,
			Value: value,
		})
	}

	return secrets
}

func FlattenContainerAppSecrets(input *containerapps.SecretsCollection) []Secret {
	if input == nil || input.Value == nil {
		return []Secret{}
	}
	result := make([]Secret, 0)
	for _, v := range input.Value {
		result = append(result, Secret{
			Name:  pointer.From(v.Name),
			Value: pointer.From(v.Value),
		})
	}

	return result
}

func FlattenContainerAppDaprSecrets(input *daprcomponents.DaprSecretsCollection) []Secret {
	if input == nil || input.Value == nil {
		return []Secret{}
	}
	result := make([]Secret, 0)
	for _, v := range input.Value {
		result = append(result, Secret{
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
