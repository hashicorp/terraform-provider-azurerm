// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package acceptance

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/testclient"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

// ApplyStep returns a Test Step which applies a Configuration and then check that the
// resource exists. This doesn't do any other assertions since it's expected that an
// ImportStep will be called afterwards to validate that.
func (td TestData) ApplyStep(config func(data TestData) string, testResource types.TestResource) resource.TestStep {
	return resource.TestStep{
		Config: config(td),
		Check: ComposeTestCheckFunc(
			check.That(td.ResourceName).ExistsInAzure(testResource),
		),
	}
}

type DisappearsStepData struct {
	// Config is a function which returns the Terraform Configuration which should be used for this step
	Config func(data TestData) string

	// TestResource is a reference to a TestResource which can destroy the resource
	// to enable a Disappears step
	TestResource types.TestResourceVerifyingRemoved
}

// DisappearsStep returns a Test Step which first confirms the resource exists
// then destroys it, and expects that the plan at the end of this should show
// that the resource needs to be created (since it's been destroyed)
func (td TestData) DisappearsStep(data DisappearsStepData) resource.TestStep {
	config := data.Config(td)
	return resource.TestStep{
		Config: config,
		Check: resource.ComposeTestCheckFunc(
			func(state *terraform.State) error {
				client, err := testclient.Build()
				if err != nil {
					return fmt.Errorf("building client: %+v", err)
				}
				return helpers.ExistsInAzure(client, data.TestResource, td.ResourceName)(state)
			},
			func(state *terraform.State) error {
				client, err := testclient.Build()
				if err != nil {
					return fmt.Errorf("building client: %+v", err)
				}
				return helpers.DeleteResourceFunc(client, data.TestResource, td.ResourceName)(state)
			},
		),
		ExpectNonEmptyPlan: true,
	}
}

type ClientCheckFunc func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error

// CheckWithClient returns a TestCheckFunc which will call a ClientCheckFunc
// with the provider context and clients
func (td TestData) CheckWithClient(check ClientCheckFunc) resource.TestCheckFunc {
	return td.CheckWithClientForResource(check, td.ResourceName)
}

// CheckWithClientForResource returns a TestCheckFunc which will call a ClientCheckFunc
// with the provider context and clients for the named resource
func (td TestData) CheckWithClientForResource(check ClientCheckFunc, resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		func(state *terraform.State) error {
			rs, ok := state.RootModule().Resources[resourceName]
			if !ok {
				return fmt.Errorf("Resource not found: %s", resourceName)
			}

			client, err := testclient.Build()
			if err != nil {
				return fmt.Errorf("building client: %+v", err)
			}
			return check(client.StopContext, client, rs.Primary)
		},
	)
}

// CheckWithClientWithoutResource returns a TestCheckFunc which will call a ClientCheckFunc
// with the provider context and clients to find if a resource exists when the resource that created it is destroyed.
func (td TestData) CheckWithClientWithoutResource(check ClientCheckFunc) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		func(state *terraform.State) error {
			client, err := testclient.Build()
			if err != nil {
				return fmt.Errorf("building client: %+v", err)
			}
			return check(client.StopContext, client, nil)
		},
	)
}

// ImportStep returns a Test Step which Imports the Resource, optionally
// ignoring any fields which may not be imported (for example, as they're
// not returned from the API)
func (td TestData) ImportStep(ignore ...string) resource.TestStep {
	return td.ImportStepFor(td.ResourceName, ignore...)
}

// ImportStepFor returns a Test Step which Imports a given resource by name,
// optionally ignoring any fields which may not be imported (for example, as they're
// not returned from the API)
func (td TestData) ImportStepFor(resourceName string, ignore ...string) resource.TestStep {
	if strings.HasPrefix(resourceName, "data.") {
		return resource.TestStep{
			ResourceName: resourceName,
			SkipFunc: func() (bool, error) {
				return false, fmt.Errorf("Data Sources (%q) do not support import - remove the ImportStep / ImportStepFor`", resourceName)
			},
		}
	}

	step := resource.TestStep{
		ResourceName:      resourceName,
		ImportState:       true,
		ImportStateVerify: true,
	}

	if len(ignore) > 0 {
		step.ImportStateVerifyIgnore = ignore
	}

	return step
}

// RequiresImportErrorStep returns a Test Step which expects a Requires Import
// error to be returned when running this step
func (td TestData) RequiresImportErrorStep(configBuilder func(data TestData) string) resource.TestStep {
	config := configBuilder(td)
	return resource.TestStep{
		Config:      config,
		ExpectError: RequiresImportError(td.ResourceType),
	}
}

// RequiresImportAssociationErrorStep returns a Test Step which expects a Requires Import
// error for an association resource to be returned when running this step
func (td TestData) RequiresImportAssociationErrorStep(configBuilder func(data TestData) string) resource.TestStep {
	config := configBuilder(td)
	return resource.TestStep{
		Config:      config,
		ExpectError: RequiresImportAssociationError(td.ResourceType),
	}
}
