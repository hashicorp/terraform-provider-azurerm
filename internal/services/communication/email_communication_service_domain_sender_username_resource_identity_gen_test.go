// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package communication_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccEmailCommunicationServiceDomainSenderUsername_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_email_communication_service_domain_sender_username", "test")
	r := EmailCommunicationServiceDomainSenderUsernameResource{}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectIdentityValue("azurerm_email_communication_service_domain_sender_username.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_email_communication_service_domain_sender_username.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_email_communication_service_domain_sender_username.test", tfjsonpath.New("domain_name"), tfjsonpath.New("email_service_domain_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_email_communication_service_domain_sender_username.test", tfjsonpath.New("email_service_name"), tfjsonpath.New("email_service_domain_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_email_communication_service_domain_sender_username.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("email_service_domain_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
