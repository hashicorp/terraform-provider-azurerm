package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SecurityCenterAssessmentMetadataResource struct{}

func TestAccSecurityCenterAssessmentMetadata_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_assessment_metadata", "test")
	r := SecurityCenterAssessmentMetadataResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityCenterAssessmentMetadata_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_assessment_metadata", "test")
	r := SecurityCenterAssessmentMetadataResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityCenterAssessmentMetadata_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_assessment_metadata", "test")
	r := SecurityCenterAssessmentMetadataResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SecurityCenterAssessmentMetadataResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	assessmentMetadataClient := client.SecurityCenter.AssessmentsMetadataClient
	id, err := parse.AssessmentMetadataID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := assessmentMetadataClient.GetInSubscription(ctx, id.AssessmentMetadataName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving Azure Security Center Assessment Metadata %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.AssessmentMetadataProperties != nil), nil
}

func (r SecurityCenterAssessmentMetadataResource) basic() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_assessment_metadata" "test" {
  display_name    = "Test Display Name"
  assessment_type = "CustomerManaged"
  severity        = "Medium"
  description     = "Test Description"
}
`
}

func (r SecurityCenterAssessmentMetadataResource) complete() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_assessment_metadata" "test" {
  display_name            = "Test Display Name"
  assessment_type         = "CustomerManaged"
  severity                = "Low"
  description             = "Test Description"
  implementation_effort   = "Low"
  remediation_description = "Test Remediation Description"
  threats                 = ["DataExfiltration", "DataSpillage", "MaliciousInsider"]
  user_impact             = "Low"
}
`
}

func (r SecurityCenterAssessmentMetadataResource) update() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_assessment_metadata" "test" {
  display_name            = "Updated Test Display Name"
  assessment_type         = "CustomerManaged"
  severity                = "Medium"
  description             = "Updated Test Description"
  implementation_effort   = "Moderate"
  remediation_description = "Updated Test Remediation Description"
  threats                 = ["DataExfiltration", "DataSpillage"]
  user_impact             = "Moderate"
}
`
}
