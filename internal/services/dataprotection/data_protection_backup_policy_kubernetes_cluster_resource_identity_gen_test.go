// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

func TestAccDataProtectionBackupPolicyKubernetesCluster_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_kubernetes_cluster", "test")
	r := DataProtectionBackupPolicyKubernetesClusterResource{}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectIdentityValue("azurerm_data_protection_backup_policy_kubernetes_cluster.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_data_protection_backup_policy_kubernetes_cluster.test", tfjsonpath.New("backup_vault_name"), tfjsonpath.New("vault_name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_data_protection_backup_policy_kubernetes_cluster.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_data_protection_backup_policy_kubernetes_cluster.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("resource_group_name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
