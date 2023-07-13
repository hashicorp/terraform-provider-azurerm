// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func TestTypedDataSourcesContainValidModelObjects(t *testing.T) {
	for _, service := range SupportedTypedServices() {
		t.Logf("Service %q..", service.Name())
		for _, resource := range service.DataSources() {
			t.Logf("- DataSources %q..", resource.ResourceType())
			obj := resource.ModelObject()
			if err := sdk.ValidateModelObject(obj); err != nil {
				t.Fatalf("validating model: %+v", err)
			}
		}
	}
}

func TestTypedResourcesContainValidModelObjects(t *testing.T) {
	for _, service := range SupportedTypedServices() {
		t.Logf("Service %q..", service.Name())
		for _, resource := range service.Resources() {
			t.Logf("- Resource %q..", resource.ResourceType())
			obj := resource.ModelObject()
			if err := sdk.ValidateModelObject(obj); err != nil {
				t.Fatalf("validating model: %+v", err)
			}
		}
	}
}

func TestTypedResourcesContainValidIDParsers(t *testing.T) {
	// This test confirms that all of the Typed Resources return an ID Validation method
	// which is used to ensure that each of the resources will validate the Resource ID
	// during import time. Whilst this may seem unnecessary as it's an interface method
	// since we could return nil, this test is double-checking.
	//
	// Untyped Resources are checked via TestUntypedResourcesContainImporters
	for _, service := range SupportedTypedServices() {
		t.Logf("Service %q..", service.Name())
		for _, resource := range service.Resources() {
			t.Logf("- Resource %q..", resource.ResourceType())
			obj := resource.IDValidationFunc()
			if obj == nil {
				t.Fatalf("IDValidationFunc returns nil - all resources must return an ID Validation function")
			}
		}
	}
}

func TestUntypedResourcesContainImporters(t *testing.T) {
	// Typed Resources are checked via TestTypedResourcesContainValidIDParsers
	// as if an ID Parser is returned it's automatically used (and it's a required
	// method on the sdk.Resource interface)
	for _, service := range SupportedUntypedServices() {
		deprecatedResourcesWhichDontSupportImport := map[string]struct{}{
			// @tombuildsstuff: new resources shouldn't be added to this list - instead add an Import function
			// to the resource, for example:
			//
			//		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			//			_, err := ParseTheId(id)
			//			return err
			//		})
			"azurerm_security_center_server_vulnerability_assessment": {},
			"azurerm_template_deployment":                             {},
		}
		for k, v := range service.SupportedResources() {
			if _, ok := deprecatedResourcesWhichDontSupportImport[k]; ok {
				t.Logf("the resource %q doesn't support import but it's deprecated so we're skipping..", k)
				continue
			}

			if v.Importer == nil {
				t.Fatalf("all resources must support import, however the resource %q does not support import", k)
			}
		}
	}
}
