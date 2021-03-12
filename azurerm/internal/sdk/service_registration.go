package sdk

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

// TypedServiceRegistration is a Service Registration using Types
// meaning that we can abstract on top of the Plugin SDK and use
// Native Types where possible
type TypedServiceRegistration interface {
	// Name is the name of this Service
	Name() string

	// PackagePath is the relative path to this package
	PackagePath() string

	// DataSources returns a list of Data Sources supported by this Service
	DataSources() []DataSource

	// Resources returns a list of Resources supported by this Service
	Resources() []Resource

	// WebsiteCategories returns a list of categories which can be used for the sidebar
	WebsiteCategories() []string
}

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
