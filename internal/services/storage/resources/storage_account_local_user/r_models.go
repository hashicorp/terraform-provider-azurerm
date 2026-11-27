// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_local_user

type PermissionsModel struct {
	Create bool `tfschema:"create"`
	Delete bool `tfschema:"delete"`
	List   bool `tfschema:"list"`
	Read   bool `tfschema:"read"`
	Write  bool `tfschema:"write"`
}

type PermissionScopeModel struct {
	Permissions  []PermissionsModel `tfschema:"permissions"`
	ResourceName string             `tfschema:"resource_name"`
	Service      string             `tfschema:"service"`
}

type SshAuthorizedKeyModel struct {
	Description string `tfschema:"description"`
	Key         string `tfschema:"key"`
}

type LocalUserModel struct {
	HomeDirectory      string                  `tfschema:"home_directory"`
	Name               string                  `tfschema:"name"`
	Password           string                  `tfschema:"password"`
	PermissionScope    []PermissionScopeModel  `tfschema:"permission_scope"`
	Sid                string                  `tfschema:"sid"`
	SshAuthorizedKey   []SshAuthorizedKeyModel `tfschema:"ssh_authorized_key"`
	SshKeyEnabled      bool                    `tfschema:"ssh_key_enabled"`
	SshPasswordEnabled bool                    `tfschema:"ssh_password_enabled"`
	StorageAccountId   string                  `tfschema:"storage_account_id"`
}

func (r Resource) ModelObject() any {
	return &LocalUserModel{}
}
