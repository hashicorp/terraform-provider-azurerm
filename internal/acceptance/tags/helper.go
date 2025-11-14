// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/types"
)

// ConfigGenerator is a function type that generates Terraform configuration for a resource
type ConfigGenerator func(data acceptance.TestData) string

// TestDefaultTagsProviderOnly tests that provider-only default tags are properly applied and visible in tags_all
func TestDefaultTagsProviderOnly(t *testing.T, resourceType string, testResource types.TestResource, configGenerator ConfigGenerator) {
	data := acceptance.BuildTestData(t, resourceType, "test")
	assert := check.That(data.ResourceName)

	data.ResourceTest(t, testResource, []acceptance.TestStep{
		{
			Config: configGenerator(data),
			Check: acceptance.ComposeTestCheckFunc(
				assert.ExistsInAzure(testResource),
				assert.Key("tags.%").HasValue("0"),
				assert.Key("tags_all.%").HasValue("2"),
				assert.Key("tags_all.managed_by").HasValue("terraform"),
				assert.Key("tags_all.environment").HasValue("test"),
			),
		},
	})
}

// TestDefaultTagsResourceOnly tests that resource-only tags are applied when no provider defaults exist
func TestDefaultTagsResourceOnly(t *testing.T, resourceType string, testResource types.TestResource, configGenerator ConfigGenerator) {
	data := acceptance.BuildTestData(t, resourceType, "test")
	assert := check.That(data.ResourceName)

	data.ResourceTest(t, testResource, []acceptance.TestStep{
		{
			Config: configGenerator(data),
			Check: acceptance.ComposeTestCheckFunc(
				assert.ExistsInAzure(testResource),
				assert.Key("tags.%").HasValue("2"),
				assert.Key("tags.cost_center").HasValue("Finance"),
				assert.Key("tags.team").HasValue("Backend"),
				assert.Key("tags_all.%").HasValue("2"),
				assert.Key("tags_all.cost_center").HasValue("Finance"),
				assert.Key("tags_all.team").HasValue("Backend"),
			),
		},
	})
}

// TestDefaultTagsProviderAndResourceNonOverlapping tests that provider and resource tags are merged when non-overlapping
func TestDefaultTagsProviderAndResourceNonOverlapping(t *testing.T, resourceType string, testResource types.TestResource, configGenerator ConfigGenerator) {
	data := acceptance.BuildTestData(t, resourceType, "test")
	assert := check.That(data.ResourceName)

	data.ResourceTest(t, testResource, []acceptance.TestStep{
		{
			Config: configGenerator(data),
			Check: acceptance.ComposeTestCheckFunc(
				assert.ExistsInAzure(testResource),
				assert.Key("tags.%").HasValue("2"),
				assert.Key("tags.cost_center").HasValue("Finance"),
				assert.Key("tags.team").HasValue("Backend"),
				assert.Key("tags_all.%").HasValue("4"),
				assert.Key("tags_all.managed_by").HasValue("terraform"),
				assert.Key("tags_all.environment").HasValue("test"),
				assert.Key("tags_all.cost_center").HasValue("Finance"),
				assert.Key("tags_all.team").HasValue("Backend"),
			),
		},
	})
}

// TestDefaultTagsProviderAndResourceOverlapping tests that resource tags take precedence when overlapping with provider defaults
func TestDefaultTagsProviderAndResourceOverlapping(t *testing.T, resourceType string, testResource types.TestResource, configGenerator ConfigGenerator) {
	data := acceptance.BuildTestData(t, resourceType, "test")
	assert := check.That(data.ResourceName)

	data.ResourceTest(t, testResource, []acceptance.TestStep{
		{
			Config: configGenerator(data),
			Check: acceptance.ComposeTestCheckFunc(
				assert.ExistsInAzure(testResource),
				assert.Key("tags.%").HasValue("2"),
				assert.Key("tags.environment").HasValue("production"),
				assert.Key("tags.team").HasValue("Backend"),
				assert.Key("tags_all.%").HasValue("3"),
				assert.Key("tags_all.environment").HasValue("production"),
				assert.Key("tags_all.managed_by").HasValue("terraform"),
				assert.Key("tags_all.team").HasValue("Backend"),
			),
		},
	})
}

// TestDefaultTagsUpdateProviderTags tests updating provider-level tags
func TestDefaultTagsUpdateProviderTags(t *testing.T, resourceType string, testResource types.TestResource, configGenerator1, configGenerator2 ConfigGenerator) {
	data := acceptance.BuildTestData(t, resourceType, "test")
	assert := check.That(data.ResourceName)

	data.ResourceTest(t, testResource, []acceptance.TestStep{
		{
			Config: configGenerator1(data),
			Check: acceptance.ComposeTestCheckFunc(
				assert.ExistsInAzure(testResource),
				assert.Key("tags_all.%").HasValue("2"),
				assert.Key("tags_all.managed_by").HasValue("terraform"),
				assert.Key("tags_all.environment").HasValue("test"),
			),
		},
		{
			Config: configGenerator2(data),
			Check: acceptance.ComposeTestCheckFunc(
				assert.ExistsInAzure(testResource),
				assert.Key("tags.%").HasValue("0"),
				assert.Key("tags_all.%").HasValue("3"),
				assert.Key("tags_all.managed_by").HasValue("terraform-updated"),
				assert.Key("tags_all.environment").HasValue("test"),
				assert.Key("tags_all.owner").HasValue("platform"),
			),
		},
	})
}

// TestDefaultTagsUpdateResourceTags tests updating resource-level tags
func TestDefaultTagsUpdateResourceTags(t *testing.T, resourceType string, testResource types.TestResource, configGenerator1, configGenerator2 ConfigGenerator) {
	data := acceptance.BuildTestData(t, resourceType, "test")
	assert := check.That(data.ResourceName)

	data.ResourceTest(t, testResource, []acceptance.TestStep{
		{
			Config: configGenerator1(data),
			Check: acceptance.ComposeTestCheckFunc(
				assert.ExistsInAzure(testResource),
				assert.Key("tags.%").HasValue("2"),
				assert.Key("tags_all.%").HasValue("4"),
			),
		},
		{
			Config: configGenerator2(data),
			Check: acceptance.ComposeTestCheckFunc(
				assert.ExistsInAzure(testResource),
				assert.Key("tags.%").HasValue("3"),
				assert.Key("tags.cost_center").HasValue("Finance"),
				assert.Key("tags.team").HasValue("Backend"),
				assert.Key("tags.project").HasValue("Project-X"),
				assert.Key("tags_all.%").HasValue("5"),
				assert.Key("tags_all.managed_by").HasValue("terraform"),
				assert.Key("tags_all.environment").HasValue("test"),
				assert.Key("tags_all.cost_center").HasValue("Finance"),
				assert.Key("tags_all.team").HasValue("Backend"),
				assert.Key("tags_all.project").HasValue("Project-X"),
			),
		},
	})
}

// TestDefaultTagsUpdateToProviderOnly tests removing resource tags to go back to provider-only tags
func TestDefaultTagsUpdateToProviderOnly(t *testing.T, resourceType string, testResource types.TestResource, configGenerator1, configGenerator2 ConfigGenerator) {
	data := acceptance.BuildTestData(t, resourceType, "test")
	assert := check.That(data.ResourceName)

	data.ResourceTest(t, testResource, []acceptance.TestStep{
		{
			Config: configGenerator1(data),
			Check: acceptance.ComposeTestCheckFunc(
				assert.ExistsInAzure(testResource),
				assert.Key("tags.%").HasValue("2"),
				assert.Key("tags.cost_center").HasValue("Finance"),
				assert.Key("tags.team").HasValue("Backend"),
				assert.Key("tags_all.%").HasValue("4"),
				assert.Key("tags_all.managed_by").HasValue("terraform"),
				assert.Key("tags_all.environment").HasValue("test"),
			),
		},
		{
			Config: configGenerator2(data),
			Check: acceptance.ComposeTestCheckFunc(
				assert.ExistsInAzure(testResource),
				assert.Key("tags.%").HasValue("0"),
				assert.Key("tags_all.%").HasValue("2"),
				assert.Key("tags_all.managed_by").HasValue("terraform"),
				assert.Key("tags_all.environment").HasValue("test"),
			),
		},
	})
}

// TestDefaultTagsUpdateToResourceOnly tests removing provider tags to go back to resource-only tags
func TestDefaultTagsUpdateToResourceOnly(t *testing.T, resourceType string, testResource types.TestResource, configGenerator1, configGenerator2 ConfigGenerator) {
	data := acceptance.BuildTestData(t, resourceType, "test")
	assert := check.That(data.ResourceName)

	data.ResourceTest(t, testResource, []acceptance.TestStep{
		{
			Config: configGenerator1(data),
			Check: acceptance.ComposeTestCheckFunc(
				assert.ExistsInAzure(testResource),
				assert.Key("tags.%").HasValue("0"),
				assert.Key("tags_all.%").HasValue("2"),
				assert.Key("tags_all.managed_by").HasValue("terraform"),
				assert.Key("tags_all.environment").HasValue("test"),
			),
		},
		{
			Config: configGenerator2(data),
			Check: acceptance.ComposeTestCheckFunc(
				assert.ExistsInAzure(testResource),
				assert.Key("tags.%").HasValue("2"),
				assert.Key("tags.cost_center").HasValue("Finance"),
				assert.Key("tags.team").HasValue("Backend"),
				assert.Key("tags_all.%").HasValue("4"),
				assert.Key("tags_all.managed_by").HasValue("terraform"),
				assert.Key("tags_all.environment").HasValue("test"),
				assert.Key("tags_all.cost_center").HasValue("Finance"),
				assert.Key("tags_all.team").HasValue("Backend"),
			),
		},
	})
}

// TestDefaultTagsProviderAndResourceDuplicateTag tests that resource tags take precedence when same key exists in both provider and resource
func TestDefaultTagsProviderAndResourceDuplicateTag(t *testing.T, resourceType string, testResource types.TestResource, configGenerator ConfigGenerator) {
	data := acceptance.BuildTestData(t, resourceType, "test")
	assert := check.That(data.ResourceName)

	data.ResourceTest(t, testResource, []acceptance.TestStep{
		{
			Config: configGenerator(data),
			Check: acceptance.ComposeTestCheckFunc(
				assert.ExistsInAzure(testResource),
				// When same key exists in both provider and resource, resource value takes precedence
				assert.Key("tags.%").HasValue("2"),
				assert.Key("tags.environment").HasValue("production"),
				assert.Key("tags.team").HasValue("Backend"),
				assert.Key("tags_all.%").HasValue("3"),
				assert.Key("tags_all.environment").HasValue("production"),
				assert.Key("tags_all.managed_by").HasValue("terraform"),
				assert.Key("tags_all.team").HasValue("Backend"),
			),
		},
	})
}
