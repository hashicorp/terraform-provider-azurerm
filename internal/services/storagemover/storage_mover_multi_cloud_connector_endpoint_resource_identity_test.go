// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagemover_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccStorageMoverMultiCloudConnectorEndpoint_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_multi_cloud_connector_endpoint", "test")
	r := StorageMoverMultiCloudConnectorEndpointTestResource{}

	multiCloudConnectorId := os.Getenv("ARM_TEST_MULTI_CLOUD_CONNECTOR_ID")
	awsS3BucketId := os.Getenv("ARM_TEST_AWS_S3_BUCKET_ID")
	if multiCloudConnectorId == "" || awsS3BucketId == "" {
		t.Skip("Skipping as ARM_TEST_MULTI_CLOUD_CONNECTOR_ID and/or ARM_TEST_AWS_S3_BUCKET_ID are not set")
	}

	checkedFields := map[string]struct{}{
		"name":                {},
		"resource_group_name": {},
		"storage_mover_name":  {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data, multiCloudConnectorId, awsS3BucketId),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_storage_mover_multi_cloud_connector_endpoint.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_storage_mover_multi_cloud_connector_endpoint.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_storage_mover_multi_cloud_connector_endpoint.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("storage_mover_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_storage_mover_multi_cloud_connector_endpoint.test", tfjsonpath.New("storage_mover_name"), tfjsonpath.New("storage_mover_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_storage_mover_multi_cloud_connector_endpoint.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("storage_mover_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
