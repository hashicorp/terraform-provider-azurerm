// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccEventgridNamespaceTopicIdAssociation_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_namespace_topic_id_association", "test")
	r := EventgridNamespaceTopicIdAssociationResource{}

	checkedFields := map[string]struct{}{
		"namespace_name":      {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_eventgrid_namespace_topic_id_association.test", checkedFields),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_eventgrid_namespace_topic_id_association.test", tfjsonpath.New("namespace_name"), tfjsonpath.New("eventgrid_namespace_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_eventgrid_namespace_topic_id_association.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("eventgrid_namespace_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_eventgrid_namespace_topic_id_association.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("eventgrid_namespace_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
