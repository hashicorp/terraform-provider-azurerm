package common

import "github.com/hashicorp/terraform/helper/schema"

// NOTE: unfortunately this has to live in it's own package to avoid a circular reference
// since the Services will need to access ArmClient which is moving to `internal/common`

type ServiceRegistration interface {
	// Name is the name of this Service
	Name() string

	// SupportedDataSources returns the supported Data Sources supported by this Service
	SupportedDataSources() map[string]*schema.Resource

	// SupportedResources returns the supported Resources supported by this Service
	SupportedResources() map[string]*schema.Resource
}
