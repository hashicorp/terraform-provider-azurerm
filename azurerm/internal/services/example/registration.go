package example

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Example"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Example",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	resources := []Resource{
		ExampleResource{},
	}
	out, err := r.magicGlueCode(resources)
	if err != nil {
		// TODO: handle errors the Go Native way
		panic(err)
	}
	return *out
}

func (r Registration) magicGlueCode(input []Resource) (*map[string]*schema.Resource, error) {
	out := make(map[string]*schema.Resource, 0)

	for _, v := range input {
		wrapper := NewResourceWrapper(v)
		resource, err := wrapper.Resource()
		if err != nil {
			return nil, err
		}

		name := v.ResourceType()
		if _, existing := out[name]; existing {
			return nil, fmt.Errorf("resource already exists with the type %q", name)
		}

		out[name] = resource
	}

	return &out, nil
}
