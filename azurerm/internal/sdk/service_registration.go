package sdk

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

// UntypedServiceRegistration is the interface used for untyped/raw Plugin SDK resources
// in the future this'll be superseded by the TypedServiceRegistration which allows for
// stronger Typed resources to be used.
type UntypedServiceRegistration interface {
	// Name is the name of this Service
	Name() string

	// WebsiteCategories returns a list of categories which can be used for the sidebar
	WebsiteCategories() []string

	// SupportedDataSources returns the supported Data Sources supported by this Service
	SupportedDataSources() map[string]*schema.Resource

	// SupportedResources returns the supported Resources supported by this Service
	SupportedResources() map[string]*schema.Resource
}
