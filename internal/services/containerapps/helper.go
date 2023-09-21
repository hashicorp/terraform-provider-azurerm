package containerapps

import (
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/jobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type TemplateModel struct {
	Containers     []ContainerModel     `tfschema:"containers"`
	Volumes        []VolumeModel        `tfschema:"volumes"`
	InitContainers []InitContainerModel `tfschema:"init_containers"`
}

type ContainerModel struct {
	Args         []string                  `tfschema:"args"`
	Command      []string                  `tfschema:"command"`
	Env          []EnvironmentVarModel     `tfschema:"env"`
	Image        string                    `tfschema:"image"`
	Name         string                    `tfschema:"name"`
	Probes       []ProbeModel              `tfschema:"probes"`
	Resources    []ContainerResourcesModel `tfschema:"resources"`
	VolumeMounts []VolumeMountModel        `tfschema:"volume_mounts"`
}

type InitContainerModel struct {
	Args         []string                  `tfschema:"args"`
	Command      []string                  `tfschema:"command"`
	Env          []EnvironmentVarModel     `tfschema:"env"`
	Image        string                    `tfschema:"image"`
	Name         string                    `tfschema:"name"`
	Resources    []ContainerResourcesModel `tfschema:"resources"`
	VolumeMounts []VolumeMountModel        `tfschema:"volume_mounts"`
}

type VolumeModel struct {
	MountOptions string              `tfschema:"mount_options"`
	Name         string              `tfschema:"name"`
	Secrets      []VolumeSecretModel `tfschema:"secrets"`
	StorageName  string              `tfschema:"storage_name"`
	StorageType  string              `tfschema:"storage_type"`
}

type EnvironmentVarModel struct {
	Name            string `tfschema:"name"`
	SecretReference string `tfschema:"secret_ref"`
	Value           string `tfschema:"value"`
}

type ProbeModel struct {
	FailureThreshold              int64                 `tfschema:"failure_threshold"`
	HttpGet                       []ProbeHttpGetModel   `tfschema:"http_get"`
	InitialDelaySeconds           int64                 `tfschema:"initial_delay_seconds"`
	PeriodSeconds                 int64                 `tfschema:"period_seconds"`
	SuccessThreshold              int64                 `tfschema:"success_threshold"`
	TcpSocket                     []ProbeTcpSocketModel `tfschema:"tcp_socket"`
	TerminationGracePeriodSeconds int64                 `tfschema:"termination_grace_period_seconds"`
	TimeoutSeconds                int64                 `tfschema:"timeout_seconds"`
	Type                          string                `tfschema:"type"`
}

type ContainerResourcesModel struct {
	CPU    float64 `tfschema:"cpu"`
	Memory string  `tfschema:"memory"`
}

type VolumeMountModel struct {
	MountPath  string `tfschema:"mount_path"`
	SubPath    string `tfschema:"sub_path"`
	VolumeName string `tfschema:"volume_name"`
}

type VolumeSecretModel struct {
	Path            string `tfschema:"path"`
	SecretReference string `tfschema:"secret_ref"`
}

type ProbeHttpGetModel struct {
	HttpHeaders []ProbeHttpGetHttpHeaderModel `tfschema:"http_headers"`
	Host        string                        `tfschema:"host"`
	Path        string                        `tfschema:"path"`
	Port        int64                         `tfschema:"port"`
	Scheme      string                        `tfschema:"scheme"`
}

type ProbeTcpSocketModel struct {
	Host string `tfschema:"host"`
	Port int64  `tfschema:"port"`
}

type ProbeHttpGetHttpHeaderModel struct {
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}

type ConfigurationModel struct {
	ReplicaRetryLimit     int64                          `tfschema:"replica_retry_limit"`
	ReplicaTimeout        int64                          `tfschema:"replica_timeout"`
	TriggerType           string                         `tfschema:"trigger_type"`
	Secret                []SecretModel                  `tfschema:"secret"`
	Registries            []RegistryModel                `tfschema:"registries"`
	EventTriggerConfig    []EventTriggerConfiguration    `tfschema:"event_trigger_config"`
	ScheduleTriggerConfig []ScheduleTriggerConfiguration `tfschema:"schedule_trigger_config"`
	ManualTriggerConfig   []ManualTriggerConfiguration   `tfschema:"manual_trigger_config"`
}

type SecretModel struct {
	Identity    string `tfschema:"identity"`
	KeyVaultURL string `tfschema:"key_vault_url"`
	Name        string `tfschema:"name"`
	Value       string `tfschema:"value"`
}

type RegistryModel struct {
	Identity                string `tfschema:"identity"`
	UserName                string `tfschema:"user_name"`
	PasswordSecretReference string `tfschema:"password_secret_ref"`
	Server                  string `tfschema:"server"`
}

type EventTriggerConfiguration struct {
	Parallelism            int64        `tfschema:"parallelism"`
	ReplicaCompletionCount int64        `tfschema:"replica_completion_count"`
	Scale                  []ScaleModel `tfschema:"scale"`
}

type ScheduleTriggerConfiguration struct {
	CronExpression         string `tfschema:"cron_expression"`
	Parallelism            int64  `tfschema:"parallelism"`
	ReplicaCompletionCount int64  `tfschema:"replica_completion_count"`
}

type ManualTriggerConfiguration struct {
	Parallelism            int64 `tfschema:"parallelism"`
	ReplicaCompletionCount int64 `tfschema:"replica_completion_count"`
}

type ScaleModel struct {
	MaxExecutions   int64       `tfschema:"max_executions"`
	MinExecutions   int64       `tfschema:"min_executions"`
	PollingInterval int64       `tfschema:"polling_interval"`
	Rules           []ScaleRule `tfschema:"rules"`
}

type ScaleRule struct {
	Auth     []ScaleRuleAuth `tfschema:"auth"`
	Metadata interface{}     `tfschema:"metadata"`
	Name     string          `tfschema:"name"`
	Type     string          `tfschema:"type"`
}

type ScaleRuleAuth struct {
	SecretReference  string `tfschema:"secret_ref"`
	TriggerParameter string `tfschema:"trigger_parameter"`
}

func containerSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"args": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"command": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"env": containerAppsJobsEnvSchema(),

				"image": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"probes": containerAppsJobsProbesSchema(),

				"resources": containerAppsJobsResourcesSchema(),

				"volume_mounts": containerAppsJobsVolumeMountsSchema(),
			},
		},
	}
}

func containerAppsJobsEnvSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"secret_ref": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"value": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func containerAppsJobsProbesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"failure_threshold": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  3,
				},

				"http_get": containerAppsJobsProbesHttpGetSchema(),

				"initial_delay_seconds": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},

				"period_seconds": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  10,
				},

				"success_threshold": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  1,
				},

				"tcp_socket": containerAppsJobsProbesTcpSocketSchema(),

				"termination_grace_period_seconds": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},

				"timeout_seconds": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},

				"type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(jobs.TypeLiveness),
						string(jobs.TypeReadiness),
						string(jobs.TypeStartup),
					}, false),
				},
			},
		},
	}
}

func containerAppsJobsResourcesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"cpu": {
					Type:         pluginsdk.TypeFloat,
					Optional:     true,
					Default:      0.5,
					ValidateFunc: validate.ContainerCpu,
				},

				"memory": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "1Gi",
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
				},
			},
		},
	}
}

func containerAppsJobsVolumeMountsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"mount_path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"sub_path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"volume_name": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
			},
		},
	}
}

func containerAppsJobsProbesHttpGetSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"http_headers": containerAppsJobsProbesHttpGetHttpHeadersSchema(),

				"host": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"port": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},

				"scheme": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(jobs.SchemeHTTP),
						string(jobs.SchemeHTTPS),
					}, false),
				},
			},
		},
	}
}

func containerAppsJobsProbesTcpSocketSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"host": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"port": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},
			},
		},
	}
}

func containerAppsJobsProbesHttpGetHttpHeadersSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"value": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func templateSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"containers": containerSchema(),

				"volumes": containerAppsJobsVolumesSchema(),

				"init_containers": containerSchema(),
			},
		},
	}
}

func containerAppsJobsVolumesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"mount_options": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"secret": containerAppsJobsVolumesSecretSchema(),

				"storage_name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"storage_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(jobs.StorageTypeAzureFile),
						string(jobs.StorageTypeEmptyDir),
						string(jobs.StorageTypeSecret),
					}, false),
				},
			},
		},
	}
}

func containerAppsJobsVolumesSecretSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"secret_ref": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func configurationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"replica_retry_limit": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},

				"replica_timeout": {
					Type:     pluginsdk.TypeInt,
					Required: true,
				},

				"trigger_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},

				"secret": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"identity": {
								Type:     pluginsdk.TypeString,
								Optional: true,
							},

							"key_vault_url": {
								Type:     pluginsdk.TypeString,
								Optional: true,
							},

							"name": {
								Type:     pluginsdk.TypeString,
								Optional: true,
							},

							"value": {
								Type:     pluginsdk.TypeString,
								Optional: true,
							},
						},
					},
				},

				"registries": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"identity": {
								Type:     pluginsdk.TypeString,
								Optional: true,
							},

							"user_name": {
								Type:     pluginsdk.TypeString,
								Optional: true,
							},

							"password_secret_ref": {
								Type:     pluginsdk.TypeString,
								Optional: true,
							},

							"server": {
								Type:     pluginsdk.TypeString,
								Optional: true,
							},
						},
					},
				},

				"event_trigger_config": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"parallelism": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},

							"replica_completion_count": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},

							"scale": containerAppsJobsScaleSchema(),
						},
					},
				},

				"schedule_trigger_config": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"cron_expression": {
								Type:     pluginsdk.TypeString,
								Optional: true,
							},

							"parallelism": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},

							"replica_completion_count": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},
						},
					},
				},

				"manual_trigger_config": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"parallelism": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},

							"replica_completion_count": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},
						},
					},
				},
			},
		},
	}
}

func containerAppsJobsScaleSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"max_executions": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},

				"min_executions": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},

				"polling_interval": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},

				"rules": containerAppsJobsScaleRulesSchema(),
			},
		},
	}
}

func containerAppsJobsScaleRulesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"auth": containerAppsJobsScaleRulesAuthSchema(),

				"metadata": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func containerAppsJobsScaleRulesAuthSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"secret_ref": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"trigger_parameter": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func expandContainerAppJobConfiguration(input []ConfigurationModel) (*jobs.JobConfiguration, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("configuration must be specified")
	}

	v := input[0]
	var configuration jobs.JobConfiguration

	if v.ReplicaRetryLimit != 0 {
		configuration.ReplicaRetryLimit = pointer.To(v.ReplicaRetryLimit)
	}

	if v.EventTriggerConfig != nil {
		eventTriggerConfig := expandContainerAppJobConfigurationEventTriggerConfig(v.EventTriggerConfig)
		configuration.EventTriggerConfig = eventTriggerConfig
	}

	if v.ManualTriggerConfig != nil {
		manualTriggerConfig := expandContainerAppJobConfigurationManualTriggerConfig(v.ManualTriggerConfig)
		configuration.ManualTriggerConfig = manualTriggerConfig
	}

	if v.ScheduleTriggerConfig != nil {
		scheduleTriggerConfig := expandContainerAppJobConfigurationScheduleTriggerConfig(v.ScheduleTriggerConfig)
		configuration.ScheduleTriggerConfig = scheduleTriggerConfig
	}

	if v.Registries != nil {
		registries := expandContainerAppJobConfigurationRegistries(v.Registries)
		configuration.Registries = registries
	}

	if v.Secret != nil {
		secrets := expandContainerAppJobConfigurationSecret(v.Secret)
		configuration.Secrets = secrets
	}

	if v.TriggerType != "" {
		configuration.TriggerType = jobs.TriggerType(v.TriggerType)
	}

	configuration.ReplicaTimeout = v.ReplicaTimeout

	return pointer.To(configuration), nil
}

func expandContainerAppJobConfigurationEventTriggerConfig(input []EventTriggerConfiguration) *jobs.JobConfigurationEventTriggerConfig {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	var eventTriggerConfig jobs.JobConfigurationEventTriggerConfig

	if v.Parallelism != 0 {
		eventTriggerConfig.Parallelism = pointer.To(v.Parallelism)
	}

	if v.ReplicaCompletionCount != 0 {
		eventTriggerConfig.ReplicaCompletionCount = pointer.To(v.ReplicaCompletionCount)
	}

	if v.Scale != nil {
		scale := expandContainerAppJobScale(v.Scale)
		eventTriggerConfig.Scale = scale
	}

	return pointer.To(eventTriggerConfig)
}

func expandContainerAppJobConfigurationManualTriggerConfig(input []ManualTriggerConfiguration) *jobs.JobConfigurationManualTriggerConfig {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	var manualTriggerConfig jobs.JobConfigurationManualTriggerConfig

	if v.Parallelism != 0 {
		manualTriggerConfig.Parallelism = pointer.To(v.Parallelism)
	}

	if v.ReplicaCompletionCount != 0 {
		manualTriggerConfig.ReplicaCompletionCount = pointer.To(v.ReplicaCompletionCount)
	}

	return pointer.To(manualTriggerConfig)
}

func expandContainerAppJobConfigurationScheduleTriggerConfig(input []ScheduleTriggerConfiguration) *jobs.JobConfigurationScheduleTriggerConfig {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	var scheduleTriggerConfig jobs.JobConfigurationScheduleTriggerConfig

	if v.CronExpression != "" {
		scheduleTriggerConfig.CronExpression = v.CronExpression
	}

	if v.Parallelism != 0 {
		scheduleTriggerConfig.Parallelism = pointer.To(v.Parallelism)
	}

	if v.ReplicaCompletionCount != 0 {
		scheduleTriggerConfig.ReplicaCompletionCount = pointer.To(v.ReplicaCompletionCount)
	}

	return pointer.To(scheduleTriggerConfig)
}

func expandContainerAppJobConfigurationRegistries(input []RegistryModel) *[]jobs.RegistryCredentials {
	var registries []jobs.RegistryCredentials
	for _, v := range input {
		var registry jobs.RegistryCredentials

		if v.Identity != "" {
			registry.Identity = pointer.To(v.Identity)
		}

		if v.PasswordSecretReference != "" {
			registry.PasswordSecretRef = pointer.To(v.PasswordSecretReference)
		}

		if v.Server != "" {
			registry.Server = pointer.To(v.Server)
		}

		if v.UserName != "" {
			registry.Username = pointer.To(v.UserName)
		}

		registries = append(registries, registry)
	}

	return pointer.To(registries)
}

func expandContainerAppJobConfigurationSecret(input []SecretModel) *[]jobs.Secret {
	var secrets []jobs.Secret
	for _, v := range input {
		var secret jobs.Secret

		if v.Identity != "" {
			secret.Identity = pointer.To(v.Identity)
		}

		if v.KeyVaultURL != "" {
			secret.KeyVaultUrl = pointer.To(v.KeyVaultURL)
		}

		if v.Name != "" {
			secret.Name = pointer.To(v.Name)
		}

		if v.Value != "" {
			secret.Value = pointer.To(v.Value)
		}

		secrets = append(secrets, secret)
	}

	return pointer.To(secrets)
}

func expandContainerAppJobScale(input []ScaleModel) *jobs.JobScale {
	v := input[0]
	var scale jobs.JobScale

	if v.MaxExecutions != 0 {
		scale.MaxExecutions = pointer.To(v.MaxExecutions)
	}

	if v.MinExecutions != 0 {
		scale.MinExecutions = pointer.To(v.MinExecutions)
	}

	if v.PollingInterval != 0 {
		scale.PollingInterval = pointer.To(v.PollingInterval)
	}

	if v.Rules != nil {
		rules := expandContainerAppJobScaleRules(v.Rules)
		scale.Rules = rules
	}

	return pointer.To(scale)
}

func expandContainerAppJobScaleRules(input []ScaleRule) *[]jobs.JobScaleRule {
	var rules []jobs.JobScaleRule
	for _, v := range input {
		var rule jobs.JobScaleRule

		if v.Auth != nil {
			auth := expandContainerAppJobScaleRulesAuth(v.Auth)
			rule.Auth = auth
		}

		if v.Metadata != "" {
			rule.Metadata = pointer.To(v.Metadata)
		}

		if v.Name != "" {
			rule.Name = pointer.To(v.Name)
		}

		if v.Type != "" {
			rule.Type = pointer.To(v.Type)
		}

		rules = append(rules, rule)
	}

	return pointer.To(rules)
}

func expandContainerAppJobScaleRulesAuth(input []ScaleRuleAuth) *[]jobs.ScaleRuleAuth {
	var auth []jobs.ScaleRuleAuth
	for _, v := range input {
		var ruleAuth jobs.ScaleRuleAuth

		if v.SecretReference != "" {
			ruleAuth.SecretRef = pointer.To(v.SecretReference)
		}

		if v.TriggerParameter != "" {
			ruleAuth.TriggerParameter = pointer.To(v.TriggerParameter)
		}

		auth = append(auth, ruleAuth)
	}

	return pointer.To(auth)
}

func expandContainerAppJobTemplate(input []TemplateModel) (*jobs.JobTemplate, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("template must be specified")
	}
	v := input[0]
	var template jobs.JobTemplate

	if len(v.Containers) != 0 {
		containers := expandContainerAppJobTemplateContainers(v.Containers)
		template.Containers = containers
	}

	if len(v.Volumes) != 0 {
		volumes := expandContainerAppJobTemplateVolumes(v.Volumes)
		template.Volumes = volumes
	}

	if len(v.InitContainers) != 0 {
		initContainers := expandContainerAppJobTemplateInitContainers(v.InitContainers)
		template.InitContainers = initContainers
	}

	return pointer.To(template), nil
}

func expandContainerAppJobTemplateInitContainers(input []InitContainerModel) *[]jobs.BaseContainer {
	var initContainers []jobs.BaseContainer
	for _, v := range input {
		var initContainer jobs.BaseContainer

		if len(v.Args) != 0 {
			initContainer.Args = pointer.To(v.Args)
		}

		if v.Name != "" {
			initContainer.Name = pointer.To(v.Name)
		}

		if len(v.Command) != 0 {
			initContainer.Command = pointer.To(v.Command)
		}

		if len(v.Env) != 0 {
			env := expandContainerAppJobTemplateContainersEnv(v.Env)
			initContainer.Env = env
		}

		if len(v.Resources) != 0 {
			resources := expandContainerAppJobTemplateContainersResources(v.Resources)
			initContainer.Resources = resources
		}

		if len(v.VolumeMounts) != 0 {
			volumeMounts := expandContainerAppJobTemplateContainersVolumeMounts(v.VolumeMounts)
			initContainer.VolumeMounts = volumeMounts
		}

		if v.Image != "" {
			initContainer.Image = pointer.To(v.Image)
		}

		initContainers = append(initContainers, initContainer)
	}

	return pointer.To(initContainers)
}

func expandContainerAppJobTemplateContainers(input []ContainerModel) *[]jobs.Container {
	var containers []jobs.Container
	for _, v := range input {
		var container jobs.Container

		if len(v.Args) != 0 {
			container.Args = pointer.To(v.Args)
		}

		if v.Name != "" {
			container.Name = pointer.To(v.Name)
		}

		if len(v.Command) != 0 {
			container.Command = pointer.To(v.Command)
		}

		if v.Image != "" {
			container.Image = pointer.To(v.Image)
		}

		if len(v.Env) != 0 {
			env := expandContainerAppJobTemplateContainersEnv(v.Env)
			container.Env = env
		}

		if len(v.Probes) != 0 {
			probes := expandContainerAppJobTemplateContainersProbes(v.Probes)
			container.Probes = probes
		}

		if len(v.Resources) != 0 {
			resources := expandContainerAppJobTemplateContainersResources(v.Resources)
			container.Resources = resources
		}

		if len(v.VolumeMounts) != 0 {
			volumeMounts := expandContainerAppJobTemplateContainersVolumeMounts(v.VolumeMounts)
			container.VolumeMounts = volumeMounts
		}
		containers = append(containers, container)
	}

	return pointer.To(containers)
}

func expandContainerAppJobTemplateContainersEnv(input []EnvironmentVarModel) *[]jobs.EnvironmentVar {
	var env []jobs.EnvironmentVar
	for _, v := range input {
		var environmentVar jobs.EnvironmentVar

		if v.Name != "" {
			environmentVar.Name = pointer.To(v.Name)
		}

		if v.SecretReference != "" {
			environmentVar.SecretRef = pointer.To(v.SecretReference)
		}

		if v.Value != "" {
			environmentVar.Value = pointer.To(v.Value)
		}
		env = append(env, environmentVar)
	}
	return pointer.To(env)
}

func expandContainerAppJobTemplateContainersProbes(input []ProbeModel) *[]jobs.ContainerAppProbe {
	var probes []jobs.ContainerAppProbe
	for _, v := range input {
		var probe jobs.ContainerAppProbe

		probe.FailureThreshold = pointer.To(v.FailureThreshold)
		probe.InitialDelaySeconds = pointer.To(v.InitialDelaySeconds)
		probe.PeriodSeconds = pointer.To(v.PeriodSeconds)
		probe.SuccessThreshold = pointer.To(v.SuccessThreshold)

		if v.TerminationGracePeriodSeconds != 0 {
			probe.TimeoutSeconds = pointer.To(v.TimeoutSeconds)
		}

		if v.Type != "" {
			probe.Type = pointer.To(jobs.Type(v.Type))
		}

		if v.HttpGet != nil {
			httpGet := expandContainerAppJobTemplateContainersProbesHttpGet(v.HttpGet)
			probe.HTTPGet = httpGet
		}

		if v.TcpSocket != nil {
			tcpSocket := expandContainerAppJobTemplateContainersProbesTcpSocket(v.TcpSocket)
			probe.TcpSocket = tcpSocket
		}

		probes = append(probes, probe)
	}

	return pointer.To(probes)
}

func expandContainerAppJobTemplateContainersProbesTcpSocket(input []ProbeTcpSocketModel) *jobs.ContainerAppProbeTcpSocket {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	var tcpSocket jobs.ContainerAppProbeTcpSocket

	tcpSocket.Host = pointer.To(v.Host)
	tcpSocket.Port = v.Port

	return pointer.To(tcpSocket)
}

func expandContainerAppJobTemplateContainersProbesHttpGet(input []ProbeHttpGetModel) *jobs.ContainerAppProbeHTTPGet {
	v := input[0]
	var httpGet jobs.ContainerAppProbeHTTPGet

	httpGet.Host = pointer.To(v.Host)
	httpGet.Path = pointer.To(v.Path)
	httpGet.Port = v.Port
	httpGet.Scheme = pointer.To(jobs.Scheme(v.Scheme))

	if len(v.HttpHeaders) != 0 {
		httpHeaders := expandContainerAppJobTemplateContainersProbesHttpGetHttpHeaders(v.HttpHeaders)
		httpGet.HTTPHeaders = httpHeaders
	}

	return pointer.To(httpGet)
}

func expandContainerAppJobTemplateContainersProbesHttpGetHttpHeaders(input []ProbeHttpGetHttpHeaderModel) *[]jobs.ContainerAppProbeHTTPGetHTTPHeadersInlined {
	var httpHeaders []jobs.ContainerAppProbeHTTPGetHTTPHeadersInlined
	for _, v := range input {
		var httpHeader jobs.ContainerAppProbeHTTPGetHTTPHeadersInlined

		if v.Name != "" {
			httpHeader.Name = v.Name
		}

		if v.Value != "" {
			httpHeader.Value = v.Value
		}
		httpHeaders = append(httpHeaders, httpHeader)
	}

	return pointer.To(httpHeaders)
}

func expandContainerAppJobTemplateContainersResources(input []ContainerResourcesModel) *jobs.ContainerResources {
	v := input[0]
	var resources jobs.ContainerResources

	if v.CPU != 0 {
		resources.Cpu = pointer.To(v.CPU)
	}

	if v.Memory != "" {
		resources.Memory = pointer.To(v.Memory)
	}
	return pointer.To(resources)
}

func expandContainerAppJobTemplateContainersVolumeMounts(input []VolumeMountModel) *[]jobs.VolumeMount {
	var volumeMounts []jobs.VolumeMount
	for _, v := range input {
		var volumeMount jobs.VolumeMount

		if v.MountPath != "" {
			volumeMount.MountPath = pointer.To(v.MountPath)
		}

		if v.SubPath != "" {
			volumeMount.SubPath = pointer.To(v.SubPath)
		}

		if v.VolumeName != "" {
			volumeMount.VolumeName = pointer.To(v.VolumeName)
		}
		volumeMounts = append(volumeMounts, volumeMount)
	}

	return pointer.To(volumeMounts)
}

func expandContainerAppJobTemplateVolumes(input []VolumeModel) *[]jobs.Volume {
	var volumes []jobs.Volume
	for _, v := range input {
		var volume jobs.Volume

		if v.Name != "" {
			volume.Name = pointer.To(v.Name)
		}

		if v.MountOptions != "" {
			volume.MountOptions = pointer.To(v.MountOptions)
		}

		if v.StorageName != "" {
			volume.StorageName = pointer.To(v.StorageName)
		}

		if v.StorageType != "" {
			volume.StorageType = pointer.To(jobs.StorageType(v.StorageType))
		}

		if len(v.Secrets) != 0 {
			secrets := expandContainerAppJobTemplateVolumesSecrets(v.Secrets)
			volume.Secrets = secrets
		}
		volumes = append(volumes, volume)
	}

	return pointer.To(volumes)
}

func expandContainerAppJobTemplateVolumesSecrets(input []VolumeSecretModel) *[]jobs.SecretVolumeItem {
	var secrets []jobs.SecretVolumeItem
	for _, v := range input {
		var secret jobs.SecretVolumeItem

		if v.Path != "" {
			secret.Path = pointer.To(v.Path)
		}

		if v.SecretReference != "" {
			secret.SecretRef = pointer.To(v.SecretReference)
		}
		secrets = append(secrets, secret)
	}

	return pointer.To(secrets)
}

func flattenContainerAppJobTemplate(input *jobs.JobTemplate) []TemplateModel {
	if input == nil {
		return nil
	}

	var template TemplateModel
	if input.Containers != nil {
		template.Containers = flattenContainerAppJobTemplateContainers(input.Containers)
	}

	if input.Volumes != nil {
		template.Volumes = flattenContainerAppJobTemplateVolumes(input.Volumes)
	}

	if input.InitContainers != nil {
		template.InitContainers = flattenContainerAppJobTemplateInitContainers(input.InitContainers)
	}

	return []TemplateModel{
		{
			Containers:     template.Containers,
			Volumes:        template.Volumes,
			InitContainers: template.InitContainers,
		},
	}
}

func flattenContainerAppJobTemplateContainers(input *[]jobs.Container) []ContainerModel {
	var containers []ContainerModel
	for _, v := range *input {
		var container ContainerModel

		if v.Args != nil {
			container.Args = *v.Args
		}

		if v.Command != nil {
			container.Command = *v.Command
		}

		if v.Env != nil {
			env := flattenContainerAppJobTemplateContainersEnv(v.Env)
			container.Env = env
		}

		if v.Image != nil {
			container.Image = *v.Image
		}

		if v.Name != nil {
			container.Name = *v.Name
		}

		if v.Probes != nil {
			probes := flattenContainerAppJobTemplateContainersProbes(v.Probes)
			container.Probes = probes
		}

		if v.Resources != nil {
			resources := flattenContainerAppJobTemplateContainersResources(v.Resources)
			container.Resources = resources
		}

		if v.VolumeMounts != nil {
			volumeMounts := flattenContainerAppJobTemplateContainersVolumeMounts(v.VolumeMounts)
			container.VolumeMounts = volumeMounts
		}

		containers = append(containers, container)
	}

	return containers
}

func flattenContainerAppJobTemplateVolumes(input *[]jobs.Volume) []VolumeModel {
	var volumes []VolumeModel
	for _, v := range *input {
		var volume VolumeModel

		if v.Name != nil {
			volume.Name = *v.Name
		}

		if v.MountOptions != nil {
			volume.MountOptions = *v.MountOptions
		}

		if v.StorageName != nil {
			volume.StorageName = *v.StorageName
		}

		if v.StorageType != nil {
			volume.StorageType = string(*v.StorageType)
		}

		if v.Secrets != nil {
			secrets := flattenContainerAppJobTemplateVolumesSecrets(v.Secrets)
			volume.Secrets = secrets
		}

		volumes = append(volumes, volume)
	}

	return volumes
}

func flattenContainerAppJobTemplateInitContainers(input *[]jobs.BaseContainer) []InitContainerModel {
	var initContainers []InitContainerModel
	for _, v := range *input {
		var initContainer InitContainerModel

		if v.Args != nil {
			initContainer.Args = *v.Args
		}

		if v.Command != nil {
			initContainer.Command = *v.Command
		}

		if v.Env != nil {
			env := flattenContainerAppJobTemplateContainersEnv(v.Env)
			initContainer.Env = env
		}

		if v.Image != nil {
			initContainer.Image = *v.Image
		}

		if v.Name != nil {
			initContainer.Name = *v.Name
		}

		if v.Resources != nil {
			resources := flattenContainerAppJobTemplateContainersResources(v.Resources)
			initContainer.Resources = resources
		}

		if v.VolumeMounts != nil {
			volumeMounts := flattenContainerAppJobTemplateContainersVolumeMounts(v.VolumeMounts)
			initContainer.VolumeMounts = volumeMounts
		}

		initContainers = append(initContainers, initContainer)
	}

	return initContainers
}

func flattenContainerAppJobTemplateContainersEnv(input *[]jobs.EnvironmentVar) []EnvironmentVarModel {
	var env []EnvironmentVarModel
	for _, v := range *input {
		var environmentVar EnvironmentVarModel

		if v.Name != nil {
			environmentVar.Name = *v.Name
		}

		if v.SecretRef != nil {
			environmentVar.SecretReference = *v.SecretRef
		}

		if v.Value != nil {
			environmentVar.Value = *v.Value
		}

		env = append(env, environmentVar)
	}

	return env
}

func flattenContainerAppJobTemplateVolumesSecrets(input *[]jobs.SecretVolumeItem) []VolumeSecretModel {
	var secrets []VolumeSecretModel
	for _, v := range *input {
		var secret VolumeSecretModel

		if v.Path != nil {
			secret.Path = *v.Path
		}

		if v.SecretRef != nil {
			secret.SecretReference = *v.SecretRef
		}

		secrets = append(secrets, secret)
	}

	return secrets
}

func flattenContainerAppJobTemplateContainersProbes(input *[]jobs.ContainerAppProbe) []ProbeModel {
	var probes []ProbeModel
	for _, v := range *input {
		var probe ProbeModel

		if v.FailureThreshold != nil {
			probe.FailureThreshold = *v.FailureThreshold
		}

		if v.InitialDelaySeconds != nil {
			probe.InitialDelaySeconds = *v.InitialDelaySeconds
		}

		if v.Type != nil {
			probe.Type = string(*v.Type)
		}

		if v.PeriodSeconds != nil {
			probe.PeriodSeconds = *v.PeriodSeconds
		}

		if v.SuccessThreshold != nil {
			probe.SuccessThreshold = *v.SuccessThreshold
		}

		if v.TerminationGracePeriodSeconds != nil {
			probe.TerminationGracePeriodSeconds = *v.TerminationGracePeriodSeconds
		}

		if v.TimeoutSeconds != nil {
			probe.TimeoutSeconds = *v.TimeoutSeconds
		}

		if v.TcpSocket != nil {
			tcpSocket := flattenContainerAppJobTemplateContainersProbesTcpSocket(v.TcpSocket)
			probe.TcpSocket = tcpSocket
		}

		if v.HTTPGet != nil {
			httpGet := flattenContainerAppJobTemplateContainersProbesHttpGet(v.HTTPGet)
			probe.HttpGet = httpGet
		}
	}

	return probes
}

func flattenContainerAppJobTemplateContainersResources(input *jobs.ContainerResources) []ContainerResourcesModel {
	var resources ContainerResourcesModel

	if input.Cpu != nil {
		resources.CPU = pointer.From(input.Cpu)
	}

	if input.Memory != nil {
		resources.Memory = pointer.From(input.Memory)
	}

	return []ContainerResourcesModel{resources}
}

func flattenContainerAppJobTemplateContainersVolumeMounts(input *[]jobs.VolumeMount) []VolumeMountModel {
	var volumeMounts []VolumeMountModel
	for _, v := range *input {
		var volumeMount VolumeMountModel

		if v.MountPath != nil {
			volumeMount.MountPath = *v.MountPath
		}

		if v.SubPath != nil {
			volumeMount.SubPath = *v.SubPath
		}

		if v.VolumeName != nil {
			volumeMount.VolumeName = *v.VolumeName
		}

		volumeMounts = append(volumeMounts, volumeMount)
	}

	return volumeMounts
}

func flattenContainerAppJobTemplateContainersProbesTcpSocket(input *jobs.ContainerAppProbeTcpSocket) []ProbeTcpSocketModel {
	var tcpSocket ProbeTcpSocketModel

	if input.Host != nil {
		tcpSocket.Host = *input.Host
	}

	tcpSocket.Port = input.Port

	return []ProbeTcpSocketModel{tcpSocket}
}

func flattenContainerAppJobTemplateContainersProbesHttpGet(input *jobs.ContainerAppProbeHTTPGet) []ProbeHttpGetModel {
	var httpGet ProbeHttpGetModel

	if input.Host != nil {
		httpGet.Host = *input.Host
	}

	if input.Path != nil {
		httpGet.Path = *input.Path
	}

	httpGet.Port = input.Port

	if input.Scheme != nil {
		httpGet.Scheme = string(*input.Scheme)
	}

	if input.HTTPHeaders != nil {
		httpHeaders := flattenContainerAppJobTemplateContainersProbesHttpGetHttpHeaders(input.HTTPHeaders)
		httpGet.HttpHeaders = httpHeaders
	}

	return []ProbeHttpGetModel{httpGet}
}

func flattenContainerAppJobTemplateContainersProbesHttpGetHttpHeaders(input *[]jobs.ContainerAppProbeHTTPGetHTTPHeadersInlined) []ProbeHttpGetHttpHeaderModel {
	var httpHeaders []ProbeHttpGetHttpHeaderModel
	for _, v := range *input {
		var httpHeader ProbeHttpGetHttpHeaderModel

		if v.Name != "" {
			httpHeader.Name = v.Name
		}

		if v.Value != "" {
			httpHeader.Value = v.Value
		}

		httpHeaders = append(httpHeaders, httpHeader)
	}

	return httpHeaders
}

func flattenContainerAppJobConfiguration(input *jobs.JobConfiguration) []ConfigurationModel {
	var config ConfigurationModel

	if input.ReplicaRetryLimit != nil {
		config.ReplicaRetryLimit = *input.ReplicaRetryLimit
	}

	if input.Registries != nil {
		registries := flattenContainerAppJobConfigurationRegistries(*input.Registries)
		config.Registries = registries
	}

	if input.EventTriggerConfig != nil {
		eventTriggerConfig := flattenContainerAppJobConfigurationEventTriggerConfig(input.EventTriggerConfig)
		config.EventTriggerConfig = eventTriggerConfig
	}

	if input.ScheduleTriggerConfig != nil {
		scheduleTriggerConfig := flattenContainerAppJobConfigurationScheduleTriggerConfig(input.ScheduleTriggerConfig)
		config.ScheduleTriggerConfig = scheduleTriggerConfig
	}

	if input.ManualTriggerConfig != nil {
		manualTriggerConfig := flattenContainerAppJobConfigurationManualTriggerConfig(input.ManualTriggerConfig)
		config.ManualTriggerConfig = manualTriggerConfig
	}

	if input.Secrets != nil {
		secrets := flattenContainerAppJobConfigurationSecret(*input.Secrets)
		config.Secret = secrets
	}

	config.TriggerType = string(input.TriggerType)

	config.ReplicaTimeout = input.ReplicaTimeout

	return []ConfigurationModel{config}
}

func flattenContainerAppJobConfigurationRegistries(input []jobs.RegistryCredentials) []RegistryModel {
	var registries []RegistryModel
	for _, v := range input {
		var registry RegistryModel

		if v.Identity != nil {
			registry.Identity = *v.Identity
		}

		if v.PasswordSecretRef != nil {
			registry.PasswordSecretReference = *v.PasswordSecretRef
		}

		if v.Server != nil {
			registry.Server = *v.Server
		}

		if v.Username != nil {
			registry.UserName = *v.Username
		}

		registries = append(registries, registry)
	}

	return registries
}

func flattenContainerAppJobConfigurationEventTriggerConfig(input *jobs.JobConfigurationEventTriggerConfig) []EventTriggerConfiguration {
	var eventTriggerConfig EventTriggerConfiguration

	if input.Parallelism != nil {
		eventTriggerConfig.Parallelism = *input.Parallelism
	}

	if input.ReplicaCompletionCount != nil {
		eventTriggerConfig.ReplicaCompletionCount = *input.ReplicaCompletionCount
	}

	if input.Scale != nil {
		scale := flattenContainerAppJobScale(input.Scale)
		eventTriggerConfig.Scale = scale
	}

	return []EventTriggerConfiguration{eventTriggerConfig}
}

func flattenContainerAppJobScale(input *jobs.JobScale) []ScaleModel {
	var scale ScaleModel

	if input.MaxExecutions != nil {
		scale.MaxExecutions = *input.MaxExecutions
	}

	if input.MinExecutions != nil {
		scale.MinExecutions = *input.MinExecutions
	}

	if input.PollingInterval != nil {
		scale.PollingInterval = *input.PollingInterval
	}

	if input.Rules != nil {
		rules := flattenContainerAppJobScaleRules(input.Rules)
		scale.Rules = rules
	}

	return []ScaleModel{scale}
}

func flattenContainerAppJobScaleRules(input *[]jobs.JobScaleRule) []ScaleRule {
	var rules []ScaleRule
	for _, v := range *input {
		var rule ScaleRule

		if v.Auth != nil {
			auth := flattenContainerAppJobScaleRulesAuth(v.Auth)
			rule.Auth = auth
		}

		if v.Metadata != nil {
			rule.Metadata = *v.Metadata
		}

		if v.Name != nil {
			rule.Name = *v.Name
		}

		if v.Type != nil {
			rule.Type = *v.Type
		}

		rules = append(rules, rule)
	}

	return rules
}

func flattenContainerAppJobConfigurationScheduleTriggerConfig(input *jobs.JobConfigurationScheduleTriggerConfig) []ScheduleTriggerConfiguration {
	var scheduleTriggerConfig ScheduleTriggerConfiguration

	if input.CronExpression != "" {
		scheduleTriggerConfig.CronExpression = input.CronExpression
	}

	if input.Parallelism != nil {
		scheduleTriggerConfig.Parallelism = *input.Parallelism
	}

	if input.ReplicaCompletionCount != nil {
		scheduleTriggerConfig.ReplicaCompletionCount = *input.ReplicaCompletionCount
	}

	return []ScheduleTriggerConfiguration{scheduleTriggerConfig}
}

func flattenContainerAppJobConfigurationManualTriggerConfig(input *jobs.JobConfigurationManualTriggerConfig) []ManualTriggerConfiguration {
	var manualTriggerConfig ManualTriggerConfiguration

	if input.Parallelism != nil {
		manualTriggerConfig.Parallelism = *input.Parallelism
	}

	if input.ReplicaCompletionCount != nil {
		manualTriggerConfig.ReplicaCompletionCount = *input.ReplicaCompletionCount
	}

	return []ManualTriggerConfiguration{manualTriggerConfig}
}

func flattenContainerAppJobConfigurationSecret(input []jobs.Secret) []SecretModel {
	var secrets []SecretModel
	for _, v := range input {
		var secret SecretModel

		if v.Identity != nil {
			secret.Identity = *v.Identity
		}

		if v.KeyVaultUrl != nil {
			secret.KeyVaultURL = *v.KeyVaultUrl
		}

		if v.Name != nil {
			secret.Name = *v.Name
		}

		if v.Value != nil {
			secret.Value = *v.Value
		}

		secrets = append(secrets, secret)
	}

	return secrets
}

func flattenContainerAppJobScaleRulesAuth(input *[]jobs.ScaleRuleAuth) []ScaleRuleAuth {
	var auth []ScaleRuleAuth
	for _, v := range *input {
		var ruleAuth ScaleRuleAuth

		if v.SecretRef != nil {
			ruleAuth.SecretReference = *v.SecretRef
		}

		if v.TriggerParameter != nil {
			ruleAuth.TriggerParameter = *v.TriggerParameter
		}

		auth = append(auth, ruleAuth)
	}

	return auth
}
