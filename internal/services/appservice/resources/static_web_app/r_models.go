// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package static_web_app

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type StaticWebAppResourceModel struct {
	Name                string                                     `tfschema:"name"`
	ResourceGroupName   string                                     `tfschema:"resource_group_name"`
	Location            string                                     `tfschema:"location"`
	AppSettings         map[string]string                          `tfschema:"app_settings"`
	BasicAuth           []BasicAuth                                `tfschema:"basic_auth"`
	ConfigFileChanges   bool                                       `tfschema:"configuration_file_changes_enabled"`
	Identity            []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	PreviewEnvironments bool                                       `tfschema:"preview_environments_enabled"`
	PublicNetworkAccess bool                                       `tfschema:"public_network_access_enabled"`
	SkuTier             string                                     `tfschema:"sku_tier"`
	SkuSize             string                                     `tfschema:"sku_size"`
	Tags                map[string]string                          `tfschema:"tags"`

	ApiKey          string `tfschema:"api_key"`
	DefaultHostName string `tfschema:"default_host_name"`

	RepositoryUrl    string `tfschema:"repository_url"`
	RepositoryToken  string `tfschema:"repository_token"`
	RepositoryBranch string `tfschema:"repository_branch"`
}

func (r Resource) ModelObject() any {
	return &StaticWebAppResourceModel{}
}
