package resourcehealth

import "github.com/Azure/go-autorest/autorest"

type ResourceHealthMetadataResource struct {
	autorest.Response `json:"-"`
	Value             []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Name       string `json:"name"`
		Properties struct {
			DisplayName         string   `json:"displayName"`
			ApplicableScenarios []string `json:"applicableScenarios,omitempty"`
			SupportedValues     []struct {
				ID            string   `json:"id"`
				DisplayName   string   `json:"displayName,omitempty"`
				ResourceTypes []string `json:"resourceTypes,omitempty"`
			} `json:"supportedValues"`
		} `json:"properties,omitempty"`
	} `json:"value"`
}
