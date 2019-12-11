package msi

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// TODO: we should probably rename this Identity, or move into Authorization

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Managed Service Identities"
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}
