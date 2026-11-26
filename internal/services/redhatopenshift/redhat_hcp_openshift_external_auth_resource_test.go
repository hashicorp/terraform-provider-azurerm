// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redhatopenshift_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2026-10-01/externalauths"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type RedhatHcpOpenshiftExternalAuthResource struct{}

func TestAccRedhatHcpOpenshiftExternalAuth_requiredPrefix(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_external_auth", "test")
	resource := RedhatHcpOpenshiftExternalAuthResource{}

	data.ResourceTest(t, resource, []acceptance.TestStep{
		{
			Config:      resource.complete(data, ""),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(regexp.QuoteMeta("`claim.0.mappings.0.username.0.prefix` is required when `claim.0.mappings.0.username.0.prefix_policy` is `Prefix`")),
		},
	})
}

func TestAccRedhatHcpOpenshiftExternalAuth_requiredClaim(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_external_auth", "test")
	resource := RedhatHcpOpenshiftExternalAuthResource{}

	data.ResourceTest(t, resource, []acceptance.TestStep{
		{
			Config: resource.basicWithValidationRules(data, `
    validation_rule {
      type = "RequiredClaim"
    }
`),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(regexp.QuoteMeta("`claim.0.validation_rule.0.required_claim` is required when `claim.0.validation_rule.0.type` is `RequiredClaim`")),
		},
		{
			Config: resource.basicWithValidationRules(data, `
    validation_rule {
      type = "RequiredClaim"
      required_claim {
        claim          = "tid"
        required_value = data.azurerm_client_config.current.tenant_id
      }
    }
    validation_rule {
      type = "RequiredClaim"
    }
`),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(regexp.QuoteMeta("`claim.0.validation_rule.1.required_claim` is required when `claim.0.validation_rule.1.type` is `RequiredClaim`")),
		},
	})
}

func TestAccRedhatHcpOpenshiftExternalAuth_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_external_auth", "test")
	resource := RedhatHcpOpenshiftExternalAuthResource{}

	data.ResourceTest(t, resource, []acceptance.TestStep{
		{
			Config: resource.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(resource),
				check.That(data.ResourceName).Key("hcp_openshift_cluster_id").MatchesOtherKey(check.That("azurerm_redhat_hcp_openshift_cluster.test").Key("id")),
				check.That(data.ResourceName).Key("issuer.0.audiences.0").MatchesOtherKey(check.That("data.azurerm_client_config.current").Key("client_id")),
				check.That(data.ResourceName).Key("claim.0.mappings.0.username.0.claim").HasValue("sub"),
				check.That(data.ResourceName).Key("claim.0.mappings.0.username.0.prefix_policy").HasValue("NoPrefix"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedhatHcpOpenshiftExternalAuth_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_external_auth", "test")
	resource := RedhatHcpOpenshiftExternalAuthResource{}

	data.ResourceTest(t, resource, []acceptance.TestStep{
		{
			Config: resource.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(resource),
		},
		data.RequiresImportErrorStep(resource.requiresImport),
	})
}

func TestAccRedhatHcpOpenshiftExternalAuth_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_external_auth", "test")
	resource := RedhatHcpOpenshiftExternalAuthResource{}

	data.ResourceTest(t, resource, []acceptance.TestStep{
		{
			Config: resource.complete(data, "oidc:"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(resource),
				check.That(data.ResourceName).Key("claim.0.mappings.0.username.0.claim").HasValue("preferred_username"),
				check.That(data.ResourceName).Key("claim.0.mappings.0.username.0.prefix_policy").HasValue("Prefix"),
				check.That(data.ResourceName).Key("claim.0.mappings.0.username.0.prefix").HasValue("oidc:"),
				check.That(data.ResourceName).Key("claim.0.mappings.0.groups.0.claim").HasValue("groups"),
				check.That(data.ResourceName).Key("claim.0.mappings.0.groups.0.prefix").HasValue("oidc:"),
				check.That(data.ResourceName).Key("claim.0.validation_rule.0.type").HasValue("RequiredClaim"),
				check.That(data.ResourceName).Key("claim.0.validation_rule.0.required_claim.0.claim").HasValue("tid"),
				check.That(data.ResourceName).Key("claim.0.validation_rule.0.required_claim.0.required_value").MatchesOtherKey(check.That("data.azurerm_client_config.current").Key("tenant_id")),
				check.That(data.ResourceName).Key("client.0.client_id").MatchesOtherKey(check.That("data.azurerm_client_config.current").Key("client_id")),
				check.That(data.ResourceName).Key("client.0.type").HasValue("Confidential"),
				check.That(data.ResourceName).Key("client.0.extra_scopes.#").HasValue("2"),
				check.That(data.ResourceName).Key("client.0.component.0.name").HasValue("console"),
				check.That(data.ResourceName).Key("client.0.component.0.auth_client_namespace").HasValue("openshift-console"),
			),
		},
		data.ImportStep(),
		{
			Config: RedhatHcpOpenshiftNodePoolResource{}.removeFromState(RedhatHcpOpenshiftClusterResource{}.basic(data)),
		},
	})
}

func TestAccRedhatHcpOpenshiftExternalAuth_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_external_auth", "test")
	resource := RedhatHcpOpenshiftExternalAuthResource{}

	data.ResourceTest(t, resource, []acceptance.TestStep{
		{
			Config: resource.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(resource),
				check.That(data.ResourceName).Key("claim.0.mappings.0.username.0.prefix_policy").HasValue("NoPrefix"),
				check.That(data.ResourceName).Key("claim.0.mappings.0.username.0.prefix").HasValue(""),
			),
		},
		data.ImportStep(),
		{
			Config: resource.basicWithPrefix(data, "oidc:"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(resource),
				check.That(data.ResourceName).Key("claim.0.mappings.0.username.0.prefix_policy").HasValue("Prefix"),
				check.That(data.ResourceName).Key("claim.0.mappings.0.username.0.prefix").HasValue("oidc:"),
			),
		},
		data.ImportStep(),
		{
			Config: resource.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(resource),
				check.That(data.ResourceName).Key("claim.0.mappings.0.username.0.prefix_policy").HasValue("NoPrefix"),
				check.That(data.ResourceName).Key("claim.0.mappings.0.username.0.prefix").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func (RedhatHcpOpenshiftExternalAuthResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := externalauths.ParseExternalAuthID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.RedHatOpenShift.ExternalAuthsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (resource RedhatHcpOpenshiftExternalAuthResource) basic(data acceptance.TestData) string {
	return resource.basicWithValidationRules(data, "")
}

func (resource RedhatHcpOpenshiftExternalAuthResource) basicWithPrefix(data acceptance.TestData, prefix string) string {
	return strings.Replace(resource.basic(data), `prefix_policy = "NoPrefix"`, fmt.Sprintf(`prefix_policy = "Prefix"
        prefix        = %q`, prefix), 1)
}

func (RedhatHcpOpenshiftExternalAuthResource) basicWithValidationRules(data acceptance.TestData, validationRules string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_hcp_openshift_external_auth" "test" {
  name                     = "acctestea%[2]s"
  hcp_openshift_cluster_id = azurerm_redhat_hcp_openshift_cluster.test.id

  issuer {
    url       = "https://login.microsoftonline.com/${data.azurerm_client_config.current.tenant_id}/v2.0"
    audiences = [data.azurerm_client_config.current.client_id]
  }

  claim {
    mappings {
      username {
        claim         = "sub"
        prefix_policy = "NoPrefix"
      }
    }
%[3]s
  }
}
`, RedhatHcpOpenshiftClusterResource{}.basic(data), data.RandomString, validationRules)
}

func (RedhatHcpOpenshiftExternalAuthResource) complete(data acceptance.TestData, prefix string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_hcp_openshift_external_auth" "test" {
  name                     = "acctestea%[2]s"
  hcp_openshift_cluster_id = azurerm_redhat_hcp_openshift_cluster.test.id

  issuer {
    url       = "https://login.microsoftonline.com/${data.azurerm_client_config.current.tenant_id}/v2.0"
    audiences = [data.azurerm_client_config.current.client_id]
  }

  claim {
    mappings {
      username {
        claim         = "preferred_username"
        prefix_policy = "Prefix"
        prefix        = "%[3]s"
      }

      groups {
        claim  = "groups"
        prefix = "%[3]s"
      }
    }

    validation_rule {
      type = "RequiredClaim"

      required_claim {
        claim          = "tid"
        required_value = data.azurerm_client_config.current.tenant_id
      }
    }
  }

  client {
    client_id    = data.azurerm_client_config.current.client_id
    type         = "Confidential"
    extra_scopes = ["email", "profile"]

    component {
      name                  = "console"
      auth_client_namespace = "openshift-console"
    }
  }

  depends_on = [azurerm_redhat_hcp_openshift_node_pool.test]
}
`, RedhatHcpOpenshiftNodePoolResource{}.basic(data), data.RandomString, prefix)
}

func (resource RedhatHcpOpenshiftExternalAuthResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_redhat_hcp_openshift_external_auth" "import" {
  name                     = azurerm_redhat_hcp_openshift_external_auth.test.name
  hcp_openshift_cluster_id = azurerm_redhat_hcp_openshift_external_auth.test.hcp_openshift_cluster_id

  issuer {
    url       = azurerm_redhat_hcp_openshift_external_auth.test.issuer[0].url
    audiences = azurerm_redhat_hcp_openshift_external_auth.test.issuer[0].audiences
  }

  claim {
    mappings {
      username {
        claim         = azurerm_redhat_hcp_openshift_external_auth.test.claim[0].mappings[0].username[0].claim
        prefix_policy = azurerm_redhat_hcp_openshift_external_auth.test.claim[0].mappings[0].username[0].prefix_policy
      }
    }
  }
}
`, resource.basic(data))
}
