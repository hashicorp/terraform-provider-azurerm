// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-01-01/buckets"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

type NetAppVolumeBucketResource struct{}

func TestAccNetAppVolumeBucket_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_bucket", "test")
	r := NetAppVolumeBucketResource{}

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

func TestAccNetAppVolumeBucket_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_bucket", "test")
	r := NetAppVolumeBucketResource{}

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

func TestAccNetAppVolumeBucket_path(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_bucket", "test")
	r := NetAppVolumeBucketResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withPath(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("path").HasValue("/data"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolumeBucket_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_bucket", "test")
	r := NetAppVolumeBucketResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolumeBucket_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_bucket", "test")
	r := NetAppVolumeBucketResource{}

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

func (t NetAppVolumeBucketResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := buckets.ParseBucketID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NetApp.BucketsClient.Get(ctx, *id)
	if err != nil {
		if resp.HttpResponse != nil && resp.HttpResponse.StatusCode == http.StatusNotFound {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(true), nil
}

// firstBucket provisions the first bucket on the shared volume using the
// `azurerm_netapp_volume_bucket_with_server` resource. The first bucket on a
// volume must establish the bucket server (FQDN + certificate); subsequent
// buckets are created with the server-less `azurerm_netapp_volume_bucket`
// resource and reuse that server configuration.
func (NetAppVolumeBucketResource) firstBucket(data acceptance.TestData) string {
	template := NetAppVolumeBucketWithServerResource{}.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_bucket_with_server" "first" {
  name      = "acctest-bucket-first-%[2]d"
  volume_id = azurerm_netapp_volume.test.id

  file_system_nfs_user {
    group_id = 1000
    user_id  = 1000
  }

  server {
    fqdn            = local.bucket_fqdn
    certificate_pem = base64encode("${tls_self_signed_cert.test.cert_pem}${tls_private_key.test.private_key_pem}")
  }
}
`, template, data.RandomInteger)
}

func (r NetAppVolumeBucketResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_bucket" "test" {
  name      = "acctest-bucket-%[2]d"
  volume_id = azurerm_netapp_volume.test.id

  file_system_nfs_user {
    group_id = 1000
    user_id  = 1000
  }

  depends_on = [azurerm_netapp_volume_bucket_with_server.first]
}
`, r.firstBucket(data), data.RandomInteger)
}

func (r NetAppVolumeBucketResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_bucket" "test" {
  name        = "acctest-bucket-%[2]d"
  volume_id   = azurerm_netapp_volume.test.id
  permissions = "ReadWrite"

  file_system_nfs_user {
    group_id = 2000
    user_id  = 2000
  }

  depends_on = [azurerm_netapp_volume_bucket_with_server.first]
}
`, r.firstBucket(data), data.RandomInteger)
}

func (r NetAppVolumeBucketResource) withPath(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_bucket" "test" {
  name      = "acctest-bucket-%[2]d"
  volume_id = azurerm_netapp_volume.test.id
  path      = "/data"

  file_system_nfs_user {
    group_id = 1000
    user_id  = 1000
  }

  depends_on = [azurerm_netapp_volume_bucket_with_server.first]
}
`, r.firstBucket(data), data.RandomInteger)
}

func (r NetAppVolumeBucketResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_bucket" "import" {
  name      = azurerm_netapp_volume_bucket.test.name
  volume_id = azurerm_netapp_volume_bucket.test.volume_id

  file_system_nfs_user {
    group_id = 1000
    user_id  = 1000
  }
}
`, r.basic(data))
}
