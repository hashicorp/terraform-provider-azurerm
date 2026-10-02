// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-10-01/agents"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SreAgentResource struct{}

// These local-draft acceptance definitions require Microsoft.App registration and
// an SRE Agent-enabled test subscription. Identity tests also require
// Microsoft.ManagedIdentity registration; complete tests create explicit role
// assignments. Registration is deliberately disabled in the provider config.
func TestAccSreAgent_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sre_agent", "test")
	r := SreAgentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctest-sre-%d", data.RandomIntOfLength(8))),
				check.That(data.ResourceName).Key("resource_group_name").MatchesOtherKey(check.That("azurerm_resource_group.test").Key("name")),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSreAgent_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sre_agent", "test")
	r := SreAgentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSreAgent_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sre_agent", "test")
	r := SreAgentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("2"),
				check.That(data.ResourceName).Key("action_configuration.0.identity_id").MatchesOtherKey(check.That("azurerm_user_assigned_identity.actions").Key("id")),
				check.That(data.ResourceName).Key("action_configuration.0.mode").HasValue("Review"),
				check.That(data.ResourceName).Key("action_configuration.0.access_level").HasValue("Low"),
				check.That(data.ResourceName).Key("resources_configuration.0.identity_id").MatchesOtherKey(check.That("azurerm_user_assigned_identity.resources").Key("id")),
				check.That(data.ResourceName).Key("resources_configuration.0.resource_ids.#").HasValue("1"),
				resource.TestCheckTypeSetElemAttrPair(data.ResourceName, "resources_configuration.0.resource_ids.*", "azurerm_resource_group.managed[0]", "id"),
				check.That(data.ResourceName).Key("tags.stage").HasValue("initial"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSreAgent_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sre_agent", "test")
	r := SreAgentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("action_configuration.0.access_level").HasValue("Low"),
				check.That(data.ResourceName).Key("resources_configuration.0.resource_ids.#").HasValue("1"),
				resource.TestCheckTypeSetElemAttrPair(data.ResourceName, "resources_configuration.0.resource_ids.*", "azurerm_resource_group.managed[0]", "id"),
				check.That(data.ResourceName).Key("tags.stage").HasValue("initial"),
			),
		},
		data.ImportStep(),
		{
			Config: r.configured(data, "High", "updated", true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("action_configuration.0.mode").HasValue("Review"),
				check.That(data.ResourceName).Key("action_configuration.0.access_level").HasValue("High"),
				check.That(data.ResourceName).Key("action_configuration.0.identity_id").MatchesOtherKey(check.That("azurerm_user_assigned_identity.actions").Key("id")),
				check.That(data.ResourceName).Key("resources_configuration.0.identity_id").MatchesOtherKey(check.That("azurerm_user_assigned_identity.resources").Key("id")),
				check.That(data.ResourceName).Key("resources_configuration.0.resource_ids.#").HasValue("2"),
				resource.TestCheckTypeSetElemAttrPair(data.ResourceName, "resources_configuration.0.resource_ids.*", "azurerm_resource_group.managed[0]", "id"),
				resource.TestCheckTypeSetElemAttrPair(data.ResourceName, "resources_configuration.0.resource_ids.*", "azurerm_resource_group.managed[1]", "id"),
				check.That(data.ResourceName).Key("tags.stage").HasValue("updated"),
			),
		},
		data.ImportStep(),
		{
			Config: r.configured(data, "Low", "", false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("action_configuration.0.access_level").HasValue("Low"),
				check.That(data.ResourceName).Key("resources_configuration.0.resource_ids.#").HasValue("1"),
				resource.TestCheckTypeSetElemAttrPair(data.ResourceName, "resources_configuration.0.resource_ids.*", "azurerm_resource_group.managed[0]", "id"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSreAgent_resourceScopeRoundTrip(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sre_agent", "test")
	r := SreAgentResource{}
	steps := make([]acceptance.TestStep, 0)
	for _, scenario := range []struct {
		members []int
		groups  int
	}{
		{[]int{0, 1, 2}, 3},
		{[]int{2, 0, 1}, 3},
		{[]int{1, 2, 3}, 4},
		{[]int{1, 2}, 4},
		{[]int{1, 2}, 3},
	} {
		checks := []acceptance.TestCheckFunc{
			check.That(data.ResourceName).ExistsInAzure(r),
			check.That(data.ResourceName).Key("resources_configuration.0.resource_ids.#").HasValue(fmt.Sprint(len(scenario.members))),
		}
		for _, index := range scenario.members {
			checks = append(checks, resource.TestCheckTypeSetElemAttrPair(data.ResourceName, "resources_configuration.0.resource_ids.*",
				fmt.Sprintf("azurerm_resource_group.managed[%d]", index), "id"))
		}
		steps = append(steps, acceptance.TestStep{
			Config: r.configuredScope(data, "Low", "scope", scenario.members, scenario.groups, false),
			Check:  acceptance.ComposeTestCheckFunc(checks...),
		}, data.ImportStep())
	}
	data.ResourceTest(t, r, steps)
}

func TestAccSreAgent_ignoreActionChanges(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sre_agent", "test")
	r := SreAgentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.configuredScope(data, "Low", "initial", []int{0}, 2, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("action_configuration.0.access_level").HasValue("Low"),
			),
		},
		data.ImportStep(),
		{
			Config: r.configuredScope(data, "High", "updated", []int{0}, 2, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("action_configuration.0.access_level").HasValue("Low"),
				check.That(data.ResourceName).Key("tags.stage").HasValue("updated"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSreAgent_identityTransitions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sre_agent", "test")
	r := SreAgentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
			),
		},
		data.ImportStep(),
		{
			Config: r.assignedIdentities(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.assignedIdentities(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				resource.TestCheckTypeSetElemAttrPair(data.ResourceName, "identity.0.identity_ids.*", "azurerm_user_assigned_identity.test[1]", "id"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func (SreAgentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := agents.ParseAgentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.SreAgent.Agents.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (SreAgentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  resource_provider_registrations = "none"
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-sre-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r SreAgentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sre_agent" "test" {
  name                = "acctest-sre-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomIntOfLength(8))
}

func (r SreAgentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sre_agent" "import" {
  name                = azurerm_sre_agent.test.name
  resource_group_name = azurerm_sre_agent.test.resource_group_name
  location            = azurerm_sre_agent.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, r.basic(data))
}

func (r SreAgentResource) assignedIdentities(data acceptance.TestData, both bool) string {
	ids := "azurerm_user_assigned_identity.test[1].id"
	if both {
		ids = "azurerm_user_assigned_identity.test[0].id, " + ids
	}
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  count               = 2
  name                = "acctest-identity-${count.index}-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_sre_agent" "test" {
  name                = "acctest-sre-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [%s]
  }
}
`, r.template(data), data.RandomIntOfLength(8), data.RandomIntOfLength(8), ids)
}

func (r SreAgentResource) complete(data acceptance.TestData) string {
	return r.configured(data, "Low", "initial", false)
}

func (r SreAgentResource) configured(data acceptance.TestData, accessLevel, stage string, expandedScope bool) string {
	members := []int{0}
	if expandedScope {
		members = append(members, 1)
	}
	return r.configuredScope(data, accessLevel, stage, members, 2, false)
}

func (r SreAgentResource) configuredScope(data acceptance.TestData, accessLevel, stage string, members []int, groups int, ignoreAction bool) string {
	resourceIDs := make([]string, len(members))
	for i, member := range members {
		resourceIDs[i] = fmt.Sprintf("azurerm_resource_group.managed[%d].id", member)
	}
	tags := ""
	if stage != "" {
		tags = fmt.Sprintf("tags = { stage = %q }", stage)
	}
	lifecycle := ""
	if ignoreAction {
		lifecycle = "lifecycle { ignore_changes = [action_configuration] }"
	}

	return fmt.Sprintf(`
%s

resource "azurerm_resource_group" "managed" {
  count = %d

  name     = "acctest-rg-sre-managed-${count.index}-%d"
  location = azurerm_resource_group.test.location
}

resource "azurerm_user_assigned_identity" "actions" {
  name                = "acctest-sre-actions-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_user_assigned_identity" "resources" {
  name                = "acctest-sre-resources-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "actions_reader" {
  count = %d

  scope                = azurerm_resource_group.managed[count.index].id
  role_definition_name = "Reader"
  principal_id         = azurerm_user_assigned_identity.actions.principal_id
}

resource "azurerm_role_assignment" "resources_reader" {
  count = %d

  scope                = azurerm_resource_group.managed[count.index].id
  role_definition_name = "Reader"
  principal_id         = azurerm_user_assigned_identity.resources.principal_id
}

resource "azurerm_sre_agent" "test" {
  name                = "acctest-sre-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.actions.id,
      azurerm_user_assigned_identity.resources.id,
    ]
  }

  action_configuration {
    identity_id  = azurerm_user_assigned_identity.actions.id
    mode         = "Review"
    access_level = "%s"
  }

  resources_configuration {
    identity_id  = azurerm_user_assigned_identity.resources.id
    resource_ids = [%s]
  }

  %s
  %s

  depends_on = [
    azurerm_role_assignment.actions_reader,
    azurerm_role_assignment.resources_reader,
  ]
}
`, r.template(data), groups, data.RandomInteger, data.RandomIntOfLength(8), data.RandomIntOfLength(8),
		groups, groups, data.RandomIntOfLength(8), accessLevel, strings.Join(resourceIDs, ", "), tags, lifecycle)
}
