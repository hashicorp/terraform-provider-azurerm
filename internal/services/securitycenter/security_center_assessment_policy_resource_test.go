// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SecurityCenterAssessmentPolicyResource struct{}

func testAccSecurityCenterAssessmentPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_assessment_policy", "test")
	r := SecurityCenterAssessmentPolicyResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccSecurityCenterAssessmentPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_assessment_policy", "test")
	r := SecurityCenterAssessmentPolicyResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccSecurityCenterAssessmentPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_assessment_policy", "test")
	r := SecurityCenterAssessmentPolicyResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SecurityCenterAssessmentPolicyResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r SecurityCenterAssessmentPolicyResource) basic() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_assessment_policy" "test" {
  display_name = "Test Display Name"
  severity     = "Medium"
  description  = "Test Description"
}
`
}

func (r SecurityCenterAssessmentPolicyResource) complete() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_assessment_policy" "test" {
  display_name            = "Test Display Name"
  severity                = "Low"
  description             = "Test Description"
  implementation_effort   = "Low"
  remediation_description = "Test Remediation Description"
  threats                 = ["DataExfiltration", "DataSpillage", "MaliciousInsider"]
  user_impact             = "Low"
  categories              = ["Data"]
}
`
}

func (r SecurityCenterAssessmentPolicyResource) update() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_assessment_policy" "test" {
  display_name            = "Updated Test Display Name"
  severity                = "Medium"
  description             = "Updated Test Description"
  implementation_effort   = "Moderate"
  remediation_description = "Updated Test Remediation Description"
  threats                 = ["DataExfiltration", "DataSpillage"]
  user_impact             = "Moderate"
}
`
}
