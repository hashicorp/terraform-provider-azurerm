// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_cache_rule

type ContainerRegistryCacheRuleModel struct {
	Name                string `tfschema:"name"`
	ContainerRegistryId string `tfschema:"container_registry_id"`
	CredentialSetId     string `tfschema:"credential_set_id"`
	SourceRepo          string `tfschema:"source_repo"`
	TargetRepo          string `tfschema:"target_repo"`
}

func (Resource) ModelObject() any {
	return &ContainerRegistryCacheRuleModel{}
}
