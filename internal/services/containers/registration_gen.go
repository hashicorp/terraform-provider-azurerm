package containers

// NOTE: this file is generated - manual changes will be overwritten.

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/kubernetes_cluster_trusted_access_role_binding"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/kubernetes_fleet_member"
)

var _ sdk.TypedServiceRegistration = autoRegistration{}

type autoRegistration struct{}

func (autoRegistration) Name() string {
	return "ContainerService"
}

func (autoRegistration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (autoRegistration) Resources() []sdk.Resource {
	return []sdk.Resource{
		kubernetes_cluster_trusted_access_role_binding.Resource{},
		kubernetes_fleet_member.Resource{},
	}
}

func (autoRegistration) WebsiteCategories() []string {
	return []string{
		"Container",
	}
}
