// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_task

type ContainerRegistryTaskModel struct {
	Name                string               `tfschema:"name"`
	ContainerRegistryId string               `tfschema:"container_registry_id"`
	AgentConfig         []AgentConfig        `tfschema:"agent_setting"`
	AgentPoolName       string               `tfschema:"agent_pool_name"`
	IsSystemTask        bool                 `tfschema:"is_system_task"`
	LogTemplate         string               `tfschema:"log_template"`
	Platform            []Platform           `tfschema:"platform"`
	Enabled             bool                 `tfschema:"enabled"`
	TimeoutInSec        int64                `tfschema:"timeout_in_seconds"`
	DockerStep          []DockerStep         `tfschema:"docker_step"`
	FileTaskStep        []FileTaskStep       `tfschema:"file_step"`
	EncodedTaskStep     []EncodedTaskStep    `tfschema:"encoded_step"`
	BaseImageTrigger    []BaseImageTrigger   `tfschema:"base_image_trigger"`
	SourceTrigger       []SourceTrigger      `tfschema:"source_trigger"`
	TimerTrigger        []TimerTrigger       `tfschema:"timer_trigger"`
	RegistryCredential  []RegistryCredential `tfschema:"registry_credential"`
	Tags                map[string]string    `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &ContainerRegistryTaskModel{}
}
