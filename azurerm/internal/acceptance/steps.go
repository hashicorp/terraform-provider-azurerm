package acceptance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/helpers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/types"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

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
				client := buildClient()
				return helpers.ExistsInAzure(client, data.TestResource, td.ResourceName)(state)
			},
			func(state *terraform.State) error {
				client := buildClient()
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
	return resource.ComposeTestCheckFunc(
		func(state *terraform.State) error {
			rs, ok := state.RootModule().Resources[td.ResourceName]
			if !ok {
				return fmt.Errorf("Resource not found found: %s", td.ResourceName)
			}

			clients := buildClient()
			return check(clients.StopContext, clients, rs.Primary)
		},
	)
}

// ImportStep returns a Test Step which Imports the Resource, optionally
// ignoring any fields which may not be imported (for example, as they're
// not returned from the API)
func (td TestData) ImportStep(ignore ...string) resource.TestStep {
	step := resource.TestStep{
		ResourceName:      td.ResourceName,
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
