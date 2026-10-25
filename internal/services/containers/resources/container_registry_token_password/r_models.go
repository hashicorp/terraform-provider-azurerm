// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_token_password

type ContainerRegistryTokenPasswordModel struct {
	TokenId   string                           `tfschema:"container_registry_token_id"`
	Password1 []ContainerRegistryTokenPassword `tfschema:"password1"`
	Password2 []ContainerRegistryTokenPassword `tfschema:"password2"`
}

func (r Resource) ModelObject() any {
	return &ContainerRegistryTokenPasswordModel{}
}
