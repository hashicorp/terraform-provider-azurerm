package helpers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2022-03-01/containerapps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2022-03-01/daprcomponents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type Registry struct {
	PasswordSecretRef string `tfschema:"password_secret_reference"`
	Server            string `tfschema:"server"`
	UserName          string `tfschema:"username"`
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
					ValidateFunc: validation.StringIsNotEmpty, // TODO - Assuming this block is supported in Preview, can we validate here?
					Description:  "The hostname for the Container Registry.",
				},

				"username": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The username to use for this Container Registry.",
				},

				"password_secret_reference": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name of the Secret Reference containing the password value for this user on the Container Registry.",
				},
			},
		},
	}
}

func ExpandContainerAppRegistries(input []Registry) *[]containerapps.RegistryCredentials {
	if input == nil {
		return nil
	}

	registries := make([]containerapps.RegistryCredentials, 0)
	for _, v := range input {
		registries = append(registries, containerapps.RegistryCredentials{
			Server:            utils.String(v.Server),
			Username:          utils.String(v.UserName),
			PasswordSecretRef: utils.String(v.PasswordSecretRef),
		})
	}

	return &registries
}

func FlattenContainerAppRegistries(input *[]containerapps.RegistryCredentials) []Registry {
	if input == nil || len(*input) == 0 {
		return nil
	}

	result := make([]Registry, 0)
	for _, v := range *input {
		result = append(result, Registry{
			PasswordSecretRef: utils.NormalizeNilableString(v.PasswordSecretRef),
			Server:            utils.NormalizeNilableString(v.Server),
			UserName:          utils.NormalizeNilableString(v.Username),
		})
	}

	return result
}

type Ingress struct {
	AllowInsecure  bool            `tfschema:"allow_insecure_connections"`
	CustomDomains  []CustomDomain  `tfschema:"custom_domain"`
	IsExternal     bool            `tfschema:"is_external"`
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

				"custom_domain": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"is_external": {
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
					Type:             pluginsdk.TypeString,
					Optional:         true,
					Default:          string(containerapps.IngressTransportMethodAuto),
					ValidateFunc:     validation.StringInSlice(containerapps.PossibleValuesForIngressTransportMethod(), false),
					Description:      "The transport method for the Ingress. Possible values include `auto`, `http`, and `http2`. Defaults to `auto`",
					DiffSuppressFunc: suppress.CaseDifference,
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
		AllowInsecure: utils.Bool(ingress.AllowInsecure),
		CustomDomains: expandContainerAppIngressCustomDomain(ingress.CustomDomains),
		External:      utils.Bool(ingress.IsExternal),
		Fqdn:          utils.String(ingress.FQDN),
		TargetPort:    utils.Int64(int64(ingress.TargetPort)),
		Traffic:       expandContainerAppIngressTraffic(ingress.TrafficWeights, appName),
	}
	transport := containerapps.IngressTransportMethod(ingress.Transport)
	result.Transport = &transport

	return result
}

func FlattenContainerAppIngress(input *containerapps.Ingress, appName string) []Ingress {
	if input == nil {
		return nil
	}

	ingress := *input
	result := Ingress{
		AllowInsecure:  utils.NormaliseNilableBool(ingress.AllowInsecure),
		CustomDomains:  flattenContainerAppIngressCustomDomain(ingress.CustomDomains),
		IsExternal:     utils.NormaliseNilableBool(ingress.External),
		FQDN:           utils.NormalizeNilableString(ingress.Fqdn),
		TargetPort:     int(utils.NormaliseNilableInt64(ingress.TargetPort)),
		TrafficWeights: flattenContainerAppIngressTraffic(ingress.Traffic, appName),
	}

	if ingress.Transport != nil {
		result.Transport = string(*ingress.Transport)
	}

	return []Ingress{result}
}

type CustomDomain struct {
	CertBinding   string `tfschema:"certificate_binding_type"`
	CertificateId string `tfschema:"certificate_id"`
	Name          string `tfschema:"name"`
}

func expandContainerAppIngressCustomDomain(input []CustomDomain) *[]containerapps.CustomDomain {
	if input == nil || len(input) == 0 {
		return nil
	}

	result := make([]containerapps.CustomDomain, 0)
	for _, v := range input {
		customDomain := containerapps.CustomDomain{
			Name:          v.Name,
			CertificateId: v.CertificateId,
		}
		bindingType := containerapps.BindingType(v.CertBinding)
		customDomain.BindingType = &bindingType

		result = append(result, customDomain)
	}

	return &result
}

func flattenContainerAppIngressCustomDomain(input *[]containerapps.CustomDomain) []CustomDomain {
	if input == nil {
		return nil
	}

	result := make([]CustomDomain, 0)

	for _, v := range *input {
		customDomain := CustomDomain{
			CertificateId: v.CertificateId,
			Name:          v.Name,
		}
		if v.BindingType != nil {
			customDomain.CertBinding = string(*v.BindingType)
		}
		result = append(result, customDomain)
	}

	return result
}

type TrafficWeight struct {
	Label          string `tfschema:"label"`
	LatestRevision bool   `tfschema:"latest_revision"`
	RevisionSuffix string `tfschema:"revision_suffix"`
	Weight         int    `tfschema:"weight"`
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

				"weight": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(0, 100), // TODO - this may just be multiples of 10?
					Description:  "The weight (%) of traffic to send to this revision.",
				},
			},
		},
	}
}

func expandContainerAppIngressTraffic(input []TrafficWeight, appName string) *[]containerapps.TrafficWeight {
	if input == nil || len(input) == 0 {
		return nil
	}

	result := make([]containerapps.TrafficWeight, 0)

	for _, v := range input {
		traffic := containerapps.TrafficWeight{
			LatestRevision: utils.Bool(v.LatestRevision),
			Weight:         utils.Int64(int64(v.Weight)),
		}

		if !v.LatestRevision {
			traffic.RevisionName = utils.String(fmt.Sprintf("%s--%s", appName, v.RevisionSuffix))
			// traffic.RevisionName = utils.String(v.RevisionName)
		}

		if v.Label != "" {
			traffic.Label = utils.String(v.Label)
		}

		result = append(result, traffic)
	}

	return &result
}

func flattenContainerAppIngressTraffic(input *[]containerapps.TrafficWeight, appName string) []TrafficWeight {
	if input == nil {
		return nil
	}

	result := make([]TrafficWeight, 0)
	for _, v := range *input {
		prefix := fmt.Sprintf("%s--", appName)
		result = append(result, TrafficWeight{
			Label:          utils.NormalizeNilableString(v.Label),
			LatestRevision: utils.NormaliseNilableBool(v.LatestRevision),
			RevisionSuffix: strings.TrimPrefix(utils.NormalizeNilableString(v.RevisionName), prefix),
			Weight:         int(utils.NormaliseNilableInt64(v.Weight)),
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
					Required:    true,
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

func ExpandContainerAppDapr(input []Dapr) *containerapps.Dapr {
	if len(input) == 0 {
		return nil
	}

	dapr := input[0]
	if dapr.AppId == "" {
		return &containerapps.Dapr{
			Enabled: utils.Bool(false),
		}
	}

	appProtocol := containerapps.AppProtocol(dapr.AppProtocol)

	return &containerapps.Dapr{
		AppId:       utils.String(dapr.AppId),
		AppPort:     utils.Int64(int64(dapr.AppPort)),
		AppProtocol: &appProtocol,
		Enabled:     utils.Bool(true),
	}
}

func FlattenContainerAppDapr(input *containerapps.Dapr) []Dapr {
	if input == nil {
		return nil
	}

	result := Dapr{
		AppId:   utils.NormalizeNilableString(input.AppId),
		AppPort: int(utils.NormaliseNilableInt64(input.AppPort)),
	}
	if appProtocol := input.AppProtocol; appProtocol != nil {
		result.AppProtocol = string(*appProtocol)
	}

	return []Dapr{result}
}

type DaprMetadata struct {
	Name      string `tfschema:"name"`
	Value     string `tfschema:"value"`
	SecretRef string `tfschema:"secret_reference"`
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

				"secret_reference": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The name of a secret specified in the `secrets` block that contains the value for this metadata configuration item.",
				},
			},
		},
	}
}

func ContainerAppEnvironmentDaprMetadataDataSourceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The name of the Metadata configuration item.",
				},

				"value": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The value for this metadata configuration item.",
				},

				"secret_reference": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
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
					ValidateFunc: validation.IntBetween(1, 10), // TODO - Double check this against API
					Description:  "The minimum number of replicas for this container.",
				},

				"max_replicas": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(1, 10), // TODO - Double check this against API
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
		template.Scale.MaxReplicas = utils.Int64(int64(config.MaxReplicas))
	}

	if config.MinReplicas != 0 {
		if template.Scale == nil {
			template.Scale = &containerapps.Scale{}
		}
		template.Scale.MinReplicas = utils.Int64(int64(config.MinReplicas))
	}

	if config.Suffix != "" {
		if metadata.ResourceData.HasChange("template.0.revision_suffix") {
			template.RevisionSuffix = utils.String(config.Suffix)
		}
	}

	return template
}

func FlattenContainerAppTemplate(input *containerapps.Template) []ContainerTemplate {
	result := ContainerTemplate{
		Containers: flattenContainerAppContainers(input.Containers),
		Suffix:     utils.NormalizeNilableString(input.RevisionSuffix),
		Volumes:    flattenContainerAppVolumes(input.Volumes),
	}

	if scale := input.Scale; scale != nil {
		result.MaxReplicas = int(utils.NormaliseNilableInt64(scale.MaxReplicas))
		result.MinReplicas = int(utils.NormaliseNilableInt64(scale.MinReplicas))
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
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty, // TODO - Check boundaries / regex
					Description:  "The name of the container",
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
					ValidateFunc: ValidateContainerCpu,
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
					Optional:    true,
					Description: "The amount of ephemeral storage available to the Container App.",
				}, // Not supported?

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
		// MaxItems: 1, // Only 1 container per app supported atm?
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
			Image:  utils.String(v.Image),
			Name:   utils.String(v.Name),
			Probes: expandContainerProbes(v),
			Resources: &containerapps.ContainerResources{
				Cpu:              utils.Float(v.CPU),
				EphemeralStorage: utils.String(v.EphemeralStorage),
				Memory:           utils.String(v.Memory),
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
		return nil
	}
	result := make([]Container, 0)
	for _, v := range *input {
		container := Container{
			Name:         utils.NormalizeNilableString(v.Name),
			Image:        utils.NormalizeNilableString(v.Image),
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

		if args := v.Args; args != nil && len(*args) != 0 {
			container.Args = *args
		}

		if command := v.Command; command != nil && len(*command) != 0 {
			container.Command = *command
		}

		if resources := v.Resources; resources != nil {
			container.CPU = utils.NormaliseNilableFloat64(resources.Cpu)
			container.Memory = utils.NormalizeNilableString(resources.Memory)
			container.EphemeralStorage = utils.NormalizeNilableString(resources.EphemeralStorage)
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
					ValidateFunc: validation.StringIsNotEmpty, // TODO - Boundary / character checks
					Description:  "The name of the volume.",
				},

				"storage_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "EmptyDir",
					ValidateFunc: validation.StringInSlice([]string{
						"EmptyDir",
						"AzureFile",
					}, false),
					Description: "The type of storage volume. Possible values include `AzureFile` and `EmptyDir`. Defaults to `EmptyDir`.",
				},

				"storage_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name of the `AzureFile` storage.",
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
			Name:        utils.String(v.Name),
			StorageName: utils.String(v.StorageName),
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
		return nil
	}

	result := make([]ContainerVolume, 0)
	for _, v := range *input {
		containerVolume := ContainerVolume{
			Name:        utils.NormalizeNilableString(v.Name),
			StorageName: utils.NormalizeNilableString(v.StorageName),
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

func expandContainerVolumeMounts(input []ContainerVolumeMount) *[]containerapps.VolumeMount {
	if input == nil {
		return nil
	}

	volumeMounts := make([]containerapps.VolumeMount, 0)
	for _, v := range input {
		volumeMounts = append(volumeMounts, containerapps.VolumeMount{
			MountPath:  utils.String(v.Path),
			VolumeName: utils.String(v.Name),
		})
	}

	return &volumeMounts
}

func flattenContainerVolumeMounts(input *[]containerapps.VolumeMount) []ContainerVolumeMount {
	if input == nil || len(*input) == 0 {
		return nil
	}

	result := make([]ContainerVolumeMount, 0)
	for _, v := range *input {
		result = append(result, ContainerVolumeMount{
			Name: utils.NormalizeNilableString(v.VolumeName),
			Path: utils.NormalizeNilableString(v.MountPath),
		})
	}

	return result
}

type ContainerEnvVar struct {
	Name            string `tfschema:"name"`
	Value           string `tfschema:"value"`
	SecretReference string `tfschema:"secret_reference"`
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
					Description: "The value for this environment variable. **NOTE:** This value is ignored if `secret_reference` is used",
				},

				"secret_reference": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
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
			Name: utils.String(v.Name),
		}
		if v.SecretReference != "" {
			env.SecretRef = utils.String(v.SecretReference)
		} else {
			env.Value = utils.String(v.Value)
		}

		envs = append(envs, env)
	}

	return &envs
}

func flattenContainerEnvVar(input *[]containerapps.EnvironmentVar) []ContainerEnvVar {
	if input == nil || len(*input) == 0 {
		return nil
	}

	result := make([]ContainerEnvVar, 0)

	for _, v := range *input {
		result = append(result, ContainerEnvVar{
			Name:            utils.NormalizeNilableString(v.Name),
			SecretReference: utils.NormalizeNilableString(v.SecretRef),
			Value:           utils.NormalizeNilableString(v.Value),
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
	Interval         int          `tfschema:"interval"`
	Timeout          int          `tfschema:"timeout"`
	FailureThreshold int          `tfschema:"failure_threshold"`
	SuccessThreshold int          `tfschema:"success_threshold"`
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
						"tcp",
						"http",
						"https",
					}, true),
					Description:      "Type of probe. Possible values are `tcp`, `http`, and `https`.",
					DiffSuppressFunc: suppress.CaseDifference,
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
					Description: "The URI to use for http type probes. Not valid for `tcp` type probes. Defaults to `/`.",
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

				"interval": {
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

				"failure_threshold": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      3,
					ValidateFunc: validation.IntBetween(1, 10),
					Description:  "The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.",
				},

				"success_threshold": {
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

func expandContainerAppReadinessProbe(input ContainerAppReadinessProbe) containerapps.ContainerAppProbe {
	probeType := containerapps.TypeReadiness
	result := containerapps.ContainerAppProbe{
		Type:             &probeType,
		PeriodSeconds:    utils.Int64(int64(input.Interval)),
		TimeoutSeconds:   utils.Int64(int64(input.Timeout)),
		FailureThreshold: utils.Int64(int64(input.FailureThreshold)),
		SuccessThreshold: utils.Int64(int64(input.SuccessThreshold)),
	}

	switch p := strings.ToUpper(input.Transport); p {
	case "HTTP", "HTTPS":
		scheme := containerapps.Scheme(p)
		result.HTTPGet = &containerapps.ContainerAppProbeHTTPGet{
			Host:   utils.String(input.Host),
			Path:   utils.String(input.Path),
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
			Host: utils.String(input.Host),
			Port: int64(input.Port),
		}
	}

	return result
}

func flattenContainerAppReadinessProbe(input containerapps.ContainerAppProbe) []ContainerAppReadinessProbe {
	result := make([]ContainerAppReadinessProbe, 0)
	probe := ContainerAppReadinessProbe{
		Interval:         int(utils.NormaliseNilableInt64(input.PeriodSeconds)),
		Timeout:          int(utils.NormaliseNilableInt64(input.TimeoutSeconds)),
		FailureThreshold: int(utils.NormaliseNilableInt64(input.FailureThreshold)),
		SuccessThreshold: int(utils.NormaliseNilableInt64(input.SuccessThreshold)),
	}

	if httpGet := input.HTTPGet; httpGet != nil {
		if httpGet.Scheme != nil {
			probe.Transport = string(*httpGet.Scheme)
		}
		probe.Host = utils.NormalizeNilableString(httpGet.Host)
		probe.Port = int(httpGet.Port)
		probe.Path = utils.NormalizeNilableString(httpGet.Path)

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
		probe.Transport = "tcp"
		probe.Host = utils.NormalizeNilableString(tcpSocket.Host)
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
	Interval               int          `tfschema:"interval"`
	Timeout                int          `tfschema:"timeout"`
	FailureThreshold       int          `tfschema:"failure_threshold"`
	TerminationGracePeriod int          `tfschema:"termination_grace_period"` // Alpha feature requiring `ProbeTerminationGracePeriod` to be enabled on the subscription?
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
						"tcp",
						"http",
						"https",
					}, true),
					Description:      "Type of probe. Possible values are `tcp`, `http`, and `https`.",
					DiffSuppressFunc: suppress.CaseDifference,
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
					Description: "The URI to use with the `host` for http type probes. Not valid for `tcp` type probes. Defaults to `/`.",
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

				"interval": {
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

				"failure_threshold": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      3,
					ValidateFunc: validation.IntBetween(1, 10),
					Description:  "The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.",
				},

				"termination_grace_period": {
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
		InitialDelaySeconds: utils.Int64(int64(input.InitialDelay)),
		PeriodSeconds:       utils.Int64(int64(input.Interval)),
		TimeoutSeconds:      utils.Int64(int64(input.Timeout)),
		FailureThreshold:    utils.Int64(int64(input.FailureThreshold)),
	}

	switch p := strings.ToUpper(input.Transport); p {
	case "HTTP", "HTTPS":
		scheme := containerapps.Scheme(p)
		result.HTTPGet = &containerapps.ContainerAppProbeHTTPGet{
			Host:   utils.String(input.Host),
			Path:   utils.String(input.Path),
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
			Host: utils.String(input.Host),
			Port: int64(input.Port),
		}
	}

	return result
}

func flattenContainerAppLivenessProbe(input containerapps.ContainerAppProbe) []ContainerAppLivenessProbe {
	result := make([]ContainerAppLivenessProbe, 0)
	probe := ContainerAppLivenessProbe{
		InitialDelay:           int(utils.NormaliseNilableInt64(input.InitialDelaySeconds)),
		Interval:               int(utils.NormaliseNilableInt64(input.PeriodSeconds)),
		Timeout:                int(utils.NormaliseNilableInt64(input.TimeoutSeconds)),
		FailureThreshold:       int(utils.NormaliseNilableInt64(input.FailureThreshold)),
		TerminationGracePeriod: int(utils.NormaliseNilableInt64(input.TerminationGracePeriodSeconds)),
	}
	if httpGet := input.HTTPGet; httpGet != nil {
		if httpGet.Scheme != nil {
			probe.Transport = string(*httpGet.Scheme)
		}
		probe.Host = utils.NormalizeNilableString(httpGet.Host)
		probe.Port = int(httpGet.Port)
		probe.Path = utils.NormalizeNilableString(httpGet.Path)

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
		probe.Transport = "tcp"
		probe.Host = utils.NormalizeNilableString(tcpSocket.Host)
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
	Interval               int          `tfschema:"interval"`
	Timeout                int          `tfschema:"timeout"`
	FailureThreshold       int          `tfschema:"failure_threshold"`
	TerminationGracePeriod int          `tfschema:"termination_grace_period"` // Alpha feature requiring `ProbeTerminationGracePeriod` to be enabled on the subscription?
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
						"tcp",
						"http",
						"https",
					}, true),
					Description:      "Type of probe. Possible values are `tcp`, `http`, and `https`.",
					DiffSuppressFunc: suppress.CaseDifference,
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
					Description: "The URI to use with the `host` for http type probes. Not valid for `tcp` type probes. Defaults to `/`.",
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

				"interval": {
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

				"failure_threshold": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      3,
					ValidateFunc: validation.IntBetween(1, 10),
					Description:  "The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.",
				},

				"termination_grace_period": {
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
		PeriodSeconds:    utils.Int64(int64(input.Interval)),
		TimeoutSeconds:   utils.Int64(int64(input.Timeout)),
		FailureThreshold: utils.Int64(int64(input.FailureThreshold)),
	}

	switch p := strings.ToUpper(input.Transport); p {
	case "HTTP", "HTTPS":
		scheme := containerapps.Scheme(p)
		result.HTTPGet = &containerapps.ContainerAppProbeHTTPGet{
			Host:   utils.String(input.Host),
			Path:   utils.String(input.Path),
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
			Host: utils.String(input.Host),
			Port: int64(input.Port),
		}
	}

	return result
}

func flattenContainerAppStartupProbe(input containerapps.ContainerAppProbe) []ContainerAppStartupProbe {
	result := make([]ContainerAppStartupProbe, 0)
	probe := ContainerAppStartupProbe{
		Interval:               int(utils.NormaliseNilableInt64(input.PeriodSeconds)),
		Timeout:                int(utils.NormaliseNilableInt64(input.TimeoutSeconds)),
		FailureThreshold:       int(utils.NormaliseNilableInt64(input.FailureThreshold)),
		TerminationGracePeriod: int(utils.NormaliseNilableInt64(input.TerminationGracePeriodSeconds)),
	}

	if httpGet := input.HTTPGet; httpGet != nil {
		if httpGet.Scheme != nil {
			probe.Transport = string(*httpGet.Scheme)
		}
		probe.Host = utils.NormalizeNilableString(httpGet.Host)
		probe.Port = int(httpGet.Port)
		probe.Path = utils.NormalizeNilableString(httpGet.Path)

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
		probe.Transport = "tcp"
		probe.Host = utils.NormalizeNilableString(tcpSocket.Host)
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
		Type:         pluginsdk.TypeList,
		Optional:     true,
		Sensitive:    true,
		ValidateFunc: nil,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: ValidateSecretName,
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
			Name:  utils.String(v.Name),
			Value: utils.String(v.Value),
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
					Name:  utils.String(v["name"].(string)),
					Value: utils.String(v["value"].(string)),
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
		result = append(result, containerapps.Secret{
			Name:  v.Name,
			Value: v.Value,
		})
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

func ExpandDaprSecrets(input []Secret) *[]daprcomponents.Secret {
	if len(input) == 0 {
		return nil
	}

	result := make([]daprcomponents.Secret, 0)

	for _, v := range input {
		result = append(result, daprcomponents.Secret{
			Name:  utils.String(v.Name),
			Value: utils.String(v.Value),
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
		return nil
	}
	result := make([]Secret, 0)
	for _, v := range input.Value {
		result = append(result, Secret{
			Name:  utils.NormalizeNilableString(v.Name),
			Value: utils.NormalizeNilableString(v.Value),
		})
	}

	return result
}

func FlattenContainerAppDaprSecrets(input *daprcomponents.DaprSecretsCollection) []Secret {
	if input == nil || input.Value == nil {
		return nil
	}
	result := make([]Secret, 0)
	for _, v := range input.Value {
		result = append(result, Secret{
			Name:  utils.NormalizeNilableString(v.Name),
			Value: utils.NormalizeNilableString(v.Value),
		})
	}

	return result
}

func ParseContainerAppLatestRevision(containerApp *containerapps.ContainerApp) (string, error) {
	if containerApp == nil || containerApp.Properties == nil || containerApp.Properties.LatestRevisionName == nil {
		return "", fmt.Errorf("could not read latest revision from API")
	}

	return *containerApp.Properties.LatestRevisionName, nil
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
