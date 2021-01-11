package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
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
			Config: r.basic(uuid.New().String()),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityCenterAssessmentMetadata_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_assessment_metadata", "test")
	r := SecurityCenterAssessmentMetadataResource{}
	uuid := uuid.New().String()

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(uuid),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImport(uuid)
		}),
	})
}

func TestAccSecurityCenterAssessmentMetadata_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_assessment_metadata", "test")
	r := SecurityCenterAssessmentMetadataResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(uuid.New().String()),
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
	uuid := uuid.New().String()

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(uuid),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(uuid),
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

func (r SecurityCenterAssessmentMetadataResource) requiresImport(uuid string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_assessment_metadata" "import" {
  name            = azurerm_security_center_assessment_metadata.test.name
  display_name    = azurerm_security_center_assessment_metadata.test.display_name
  assessment_type = azurerm_security_center_assessment_metadata.test.assessment_type
  severity        = azurerm_security_center_assessment_metadata.test.severity
  description     = azurerm_security_center_assessment_metadata.test.description
}
`, r.basic(uuid))
}

func (r SecurityCenterAssessmentMetadataResource) basic(uuid string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_assessment_metadata" "test" {
  name            = "%s"
  display_name    = "Test Display Name"
  assessment_type = "CustomerManaged"
  severity        = "Medium"
  description     = "Test Description"
}
`, uuid)
}

func (r SecurityCenterAssessmentMetadataResource) complete(uuid string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_assessment_metadata" "test" {
  name                    = "%s"
  display_name            = "Test Display Name"
  assessment_type         = "CustomerManaged"
  severity                = "Low"
  description             = "Test Description"
  implementation_effort   = "Low"
  is_preview              = false
  remediation_description = "Test Remediation Description"
  threats                 = ["DataExfiltration", "DataSpillage", "MaliciousInsider"]
  user_impact             = "Low"
}
`, uuid)
}

func (r SecurityCenterAssessmentMetadataResource) update(uuid string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_assessment_metadata" "test" {
  name                    = "%s"
  display_name            = "Updated Test Display Name"
  assessment_type         = "CustomerManaged"
  severity                = "Medium"
  description             = "Updated Test Description"
  implementation_effort   = "Moderate"
  is_preview              = true
  remediation_description = "Updated Test Remediation Description"
  threats                 = ["DataExfiltration", "DataSpillage"]
  user_impact             = "Moderate"
}
`, uuid)
}
