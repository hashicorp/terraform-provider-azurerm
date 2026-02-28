// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagemover_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StorageMoverMultiCloudConnectorEndpointTestResource struct{}

func TestAccStorageMoverMultiCloudConnectorEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_multi_cloud_connector_endpoint", "test")
	r := StorageMoverMultiCloudConnectorEndpointTestResource{}

	multiCloudConnectorId := os.Getenv("ARM_TEST_MULTI_CLOUD_CONNECTOR_ID")
	awsS3BucketId := os.Getenv("ARM_TEST_AWS_S3_BUCKET_ID")
	if multiCloudConnectorId == "" || awsS3BucketId == "" {
		t.Skip("Skipping as ARM_TEST_MULTI_CLOUD_CONNECTOR_ID and/or ARM_TEST_AWS_S3_BUCKET_ID are not set")
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, multiCloudConnectorId, awsS3BucketId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageMoverMultiCloudConnectorEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_multi_cloud_connector_endpoint", "test")
	r := StorageMoverMultiCloudConnectorEndpointTestResource{}

	multiCloudConnectorId := os.Getenv("ARM_TEST_MULTI_CLOUD_CONNECTOR_ID")
	awsS3BucketId := os.Getenv("ARM_TEST_AWS_S3_BUCKET_ID")
	if multiCloudConnectorId == "" || awsS3BucketId == "" {
		t.Skip("Skipping as ARM_TEST_MULTI_CLOUD_CONNECTOR_ID and/or ARM_TEST_AWS_S3_BUCKET_ID are not set")
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, multiCloudConnectorId, awsS3BucketId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImport(data, multiCloudConnectorId, awsS3BucketId)
		}),
	})
}

func TestAccStorageMoverMultiCloudConnectorEndpoint_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_multi_cloud_connector_endpoint", "test")
	r := StorageMoverMultiCloudConnectorEndpointTestResource{}

	multiCloudConnectorId := os.Getenv("ARM_TEST_MULTI_CLOUD_CONNECTOR_ID")
	awsS3BucketId := os.Getenv("ARM_TEST_AWS_S3_BUCKET_ID")
	if multiCloudConnectorId == "" || awsS3BucketId == "" {
		t.Skip("Skipping as ARM_TEST_MULTI_CLOUD_CONNECTOR_ID and/or ARM_TEST_AWS_S3_BUCKET_ID are not set")
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, multiCloudConnectorId, awsS3BucketId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageMoverMultiCloudConnectorEndpoint_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_multi_cloud_connector_endpoint", "test")
	r := StorageMoverMultiCloudConnectorEndpointTestResource{}

	multiCloudConnectorId := os.Getenv("ARM_TEST_MULTI_CLOUD_CONNECTOR_ID")
	awsS3BucketId := os.Getenv("ARM_TEST_AWS_S3_BUCKET_ID")
	if multiCloudConnectorId == "" || awsS3BucketId == "" {
		t.Skip("Skipping as ARM_TEST_MULTI_CLOUD_CONNECTOR_ID and/or ARM_TEST_AWS_S3_BUCKET_ID are not set")
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, multiCloudConnectorId, awsS3BucketId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, multiCloudConnectorId, awsS3BucketId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageMoverMultiCloudConnectorEndpointTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := endpoints.ParseEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.StorageMover.EndpointsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r StorageMoverMultiCloudConnectorEndpointTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_storage_mover" "test" {
  name                = "acctest-ssm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r StorageMoverMultiCloudConnectorEndpointTestResource) basic(data acceptance.TestData, multiCloudConnectorId, awsS3BucketId string) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_mover_multi_cloud_connector_endpoint" "test" {
  name                     = "acctest-smmcce-%d"
  storage_mover_id         = azurerm_storage_mover.test.id
  multi_cloud_connector_id = "%s"
  aws_s3_bucket_id         = "%s"
}
`, template, data.RandomInteger, multiCloudConnectorId, awsS3BucketId)
}

func (r StorageMoverMultiCloudConnectorEndpointTestResource) requiresImport(data acceptance.TestData, multiCloudConnectorId, awsS3BucketId string) string {
	config := r.basic(data, multiCloudConnectorId, awsS3BucketId)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_mover_multi_cloud_connector_endpoint" "import" {
  name                     = azurerm_storage_mover_multi_cloud_connector_endpoint.test.name
  storage_mover_id         = azurerm_storage_mover.test.id
  multi_cloud_connector_id = azurerm_storage_mover_multi_cloud_connector_endpoint.test.multi_cloud_connector_id
  aws_s3_bucket_id         = azurerm_storage_mover_multi_cloud_connector_endpoint.test.aws_s3_bucket_id
}
`, config)
}

func (r StorageMoverMultiCloudConnectorEndpointTestResource) complete(data acceptance.TestData, multiCloudConnectorId, awsS3BucketId string) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_mover_multi_cloud_connector_endpoint" "test" {
  name                     = "acctest-smmcce-%d"
  storage_mover_id         = azurerm_storage_mover.test.id
  multi_cloud_connector_id = "%s"
  aws_s3_bucket_id         = "%s"
  description              = "Example Multi-Cloud Connector Endpoint Description"
}
`, template, data.RandomInteger, multiCloudConnectorId, awsS3BucketId)
}

func (r StorageMoverMultiCloudConnectorEndpointTestResource) update(data acceptance.TestData, multiCloudConnectorId, awsS3BucketId string) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_mover_multi_cloud_connector_endpoint" "test" {
  name                     = "acctest-smmcce-%d"
  storage_mover_id         = azurerm_storage_mover.test.id
  multi_cloud_connector_id = "%s"
  aws_s3_bucket_id         = "%s"
  description              = "Updated Multi-Cloud Connector Endpoint Description"
}
`, template, data.RandomInteger, multiCloudConnectorId, awsS3BucketId)
}
