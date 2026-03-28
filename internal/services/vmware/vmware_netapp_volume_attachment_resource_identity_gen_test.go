// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package vmware_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccVmwareNetappVolumeAttachment_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vmware_netapp_volume_attachment", "test")
	r := VmwareNetappVolumeAttachmentResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"cluster_name":        {},
		"private_cloud_name":  {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_vmware_netapp_volume_attachment.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_vmware_netapp_volume_attachment.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_vmware_netapp_volume_attachment.test", tfjsonpath.New("cluster_name"), tfjsonpath.New("vmware_cluster_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_vmware_netapp_volume_attachment.test", tfjsonpath.New("private_cloud_name"), tfjsonpath.New("vmware_cluster_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_vmware_netapp_volume_attachment.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("vmware_cluster_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_vmware_netapp_volume_attachment.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("vmware_cluster_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
