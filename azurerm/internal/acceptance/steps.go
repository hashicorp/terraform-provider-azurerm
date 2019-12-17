package acceptance

import "github.com/hashicorp/terraform-plugin-sdk/helper/resource"

type DisappearsStepData struct {
	// Config is a function which returns the Terraform Configuration which should be used for this step
	Config func(data TestData) string

	// CheckExists is a function which confirms that the given Resource in
	// the state exists
	CheckExists func(resourceName string) resource.TestCheckFunc

	// Destroy is a function which looks up the given Resource in the State
	// and then ensures that it's deleted
	Destroy func(resourceName string) resource.TestCheckFunc
}

// DisappearsStep returns a Test Step which first confirms the resource exists
// then destroys it, and expects that the plan at the end of this should show
// that the resource needs to be created (since it's been destroyed)
func (td TestData) DisappearsStep(data DisappearsStepData) resource.TestStep {
	config := data.Config(td)
	return resource.TestStep{
		Config: config,
		Check: resource.ComposeTestCheckFunc(
			data.CheckExists(td.ResourceName),
			data.Destroy(td.ResourceName),
		),
		ExpectNonEmptyPlan: true,
	}
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
		ExpectError: RequiresImportError(td.resourceType),
	}
}
