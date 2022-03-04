package arckubernetes_test

import (
	"context"
	"fmt"
	"testing"

	arckubernetes "github.com/hashicorp/terraform-provider-azurerm/internal/services/arckubernetes/sdk/2021-10-01/hybridkubernetes"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ArcKubernetesClusterResource struct{}

func TestAccArcKubernetesCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_cluster", "test")
	r := ArcKubernetesClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccArcKubernetesCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_cluster", "test")
	r := ArcKubernetesClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccArcKubernetesCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_cluster", "test")
	r := ArcKubernetesClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccArcKubernetesCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_cluster", "test")
	r := ArcKubernetesClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ArcKubernetesClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := arckubernetes.ParseConnectedClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.ArcKubernetes.ArcKubernetesClient
	resp, err := client.ConnectedClusterGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r ArcKubernetesClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ArcKubernetesClusterResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_arc_kubernetes_cluster" "test" {
  name                         = "acctest-akcc-%d"
  resource_group_name          = azurerm_resource_group.test.name
  agent_public_key_certificate = "MIICCgKCAgEAsSpALlON3394ysLQdRSy6cCBwL08NgZp7c1xsy0kQH/wHuixfoCwtL1OZ0a5kqod9vE6L8ICsXAE+iEdU1OspcJxL9J/gSyiOCMYPUabbYRXFy5x258RRLtn60NoaqcaDW+Z80HLwJOMECdJ/yDkuuNbnL0J2cyR8/WXjoeee8cG52QmDuxB6a4ROOushroIE2NS3FuhJh3b3Ddj+NU3gfbgIjjCMJQDpJWqUi/68lB33K1FQwuY1CCT8rKC5dDRKnlTgEWtnhBBly1D+C9GPTh7l9NTxcAEcBKo1ZIkMcixcS+gTcsTtDRRwuTtt8kybohfRMvmBA/S9bkD6cCaGJMe8YerOyJsDd4zSUHz/qN9iTt8FDdFcCIjhEdsMKl350wMj7+UNPvchlRAAve3oCIntW9063fDlQHrsaPfhCplTlKLAt1jcUkGSaeuylnRi8te+hmDCIhmo8wDqMv8Yvy7BTX4bYMg+6j0EGiIKPMRf5NHt6bXbiWEcO5LWEwgCbvpTjf7XdrU/xJ+eB+uAP1etvE0tCYAEFisfFnqNZxMQCFjovI0ZNxizrZOBznk15fWTr4KTNsPUTkEvbJfZPQqc/QyC5yqIAFHAAj+jusd4tK9f19nCsoi78xeKGH2s1zrD0AmBNGUrPLPVzgfdueadJKbVbiPteYec9qSG50CAwEAAQ=="
  identity {
    type = "SystemAssigned"
  }
  location = "%s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r ArcKubernetesClusterResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_arc_kubernetes_cluster" "import" {
  name                         = azurerm_arc_kubernetes_cluster.test.name
  resource_group_name          = azurerm_arc_kubernetes_cluster.test.resource_group_name
  agent_public_key_certificate = azurerm_arc_kubernetes_cluster.test.agent_public_key_certificate
  identity {
    type = "SystemAssigned"
  }
  location = azurerm_arc_kubernetes_cluster.test.location
}
`, config)
}

func (r ArcKubernetesClusterResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_arc_kubernetes_cluster" "test" {
  name                         = "acctest-akcc-%d"
  resource_group_name          = azurerm_resource_group.test.name
  agent_public_key_certificate = "MIICCgKCAgEAsSpALlON3394ysLQdRSy6cCBwL08NgZp7c1xsy0kQH/wHuixfoCwtL1OZ0a5kqod9vE6L8ICsXAE+iEdU1OspcJxL9J/gSyiOCMYPUabbYRXFy5x258RRLtn60NoaqcaDW+Z80HLwJOMECdJ/yDkuuNbnL0J2cyR8/WXjoeee8cG52QmDuxB6a4ROOushroIE2NS3FuhJh3b3Ddj+NU3gfbgIjjCMJQDpJWqUi/68lB33K1FQwuY1CCT8rKC5dDRKnlTgEWtnhBBly1D+C9GPTh7l9NTxcAEcBKo1ZIkMcixcS+gTcsTtDRRwuTtt8kybohfRMvmBA/S9bkD6cCaGJMe8YerOyJsDd4zSUHz/qN9iTt8FDdFcCIjhEdsMKl350wMj7+UNPvchlRAAve3oCIntW9063fDlQHrsaPfhCplTlKLAt1jcUkGSaeuylnRi8te+hmDCIhmo8wDqMv8Yvy7BTX4bYMg+6j0EGiIKPMRf5NHt6bXbiWEcO5LWEwgCbvpTjf7XdrU/xJ+eB+uAP1etvE0tCYAEFisfFnqNZxMQCFjovI0ZNxizrZOBznk15fWTr4KTNsPUTkEvbJfZPQqc/QyC5yqIAFHAAj+jusd4tK9f19nCsoi78xeKGH2s1zrD0AmBNGUrPLPVzgfdueadJKbVbiPteYec9qSG50CAwEAAQ=="
  distribution                 = "kind"

  identity {
    type = "SystemAssigned"
  }
  infrastructure = "generic"
  location       = "%s"

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r ArcKubernetesClusterResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_arc_kubernetes_cluster" "test" {
  name                         = "acctest-akcc-%d"
  resource_group_name          = azurerm_resource_group.test.name
  agent_public_key_certificate = "MIICCgKCAgEAsSpALlON3394ysLQdRSy6cCBwL08NgZp7c1xsy0kQH/wHuixfoCwtL1OZ0a5kqod9vE6L8ICsXAE+iEdU1OspcJxL9J/gSyiOCMYPUabbYRXFy5x258RRLtn60NoaqcaDW+Z80HLwJOMECdJ/yDkuuNbnL0J2cyR8/WXjoeee8cG52QmDuxB6a4ROOushroIE2NS3FuhJh3b3Ddj+NU3gfbgIjjCMJQDpJWqUi/68lB33K1FQwuY1CCT8rKC5dDRKnlTgEWtnhBBly1D+C9GPTh7l9NTxcAEcBKo1ZIkMcixcS+gTcsTtDRRwuTtt8kybohfRMvmBA/S9bkD6cCaGJMe8YerOyJsDd4zSUHz/qN9iTt8FDdFcCIjhEdsMKl350wMj7+UNPvchlRAAve3oCIntW9063fDlQHrsaPfhCplTlKLAt1jcUkGSaeuylnRi8te+hmDCIhmo8wDqMv8Yvy7BTX4bYMg+6j0EGiIKPMRf5NHt6bXbiWEcO5LWEwgCbvpTjf7XdrU/xJ+eB+uAP1etvE0tCYAEFisfFnqNZxMQCFjovI0ZNxizrZOBznk15fWTr4KTNsPUTkEvbJfZPQqc/QyC5yqIAFHAAj+jusd4tK9f19nCsoi78xeKGH2s1zrD0AmBNGUrPLPVzgfdueadJKbVbiPteYec9qSG50CAwEAAQ=="
  distribution                 = "kind"

  identity {
    type = "SystemAssigned"
  }
  infrastructure = "generic"
  location       = "%s"

  tags = {
    ENV = "TestUpdate"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}
