package provider

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/sdk"
)

func TestTypedDataSourcesContainValidModelObjects(t *testing.T) {
	for _, service := range SupportedTypedServices() {
		t.Logf("Service %q..", service.Name())
		for _, resource := range service.SupportedDataSources() {
			t.Logf("- DataSources %q..", resource.ResourceType())
			obj := resource.ModelObject()
			if err := sdk.ValidateModelObject(&obj); err != nil {
				t.Fatalf("validating model: %+v", err)
			}
		}
	}
}

func TestTypedResourcesContainValidModelObjects(t *testing.T) {
	for _, service := range SupportedTypedServices() {
		t.Logf("Service %q..", service.Name())
		for _, resource := range service.SupportedResources() {
			t.Logf("- Resource %q..", resource.ResourceType())
			obj := resource.ModelObject()
			if err := sdk.ValidateModelObject(&obj); err != nil {
				t.Fatalf("validating model: %+v", err)
			}
		}
	}
}
