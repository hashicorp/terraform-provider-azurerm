// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/blobs"
)

type StorageBlobResource struct{}

func TestAccStorageBlob_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.blockEmpty,
			TestResource: r,
		}),
	})
}

func TestAccStorageBlob_appendEmpty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appendEmpty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "type"),
	})
}

func TestAccStorageBlob_appendEmptyMetaData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appendEmptyMetaData(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "type"),
	})
}

func TestAccStorageBlob_blockEmpty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blockEmpty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "type"),
	})
}

func TestAccStorageBlob_blockEmptyAzureADAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blockEmptyAzureADAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "type"),
	})
}

func TestAccStorageBlob_blockEmptyMetaData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blockEmptyMetaData(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "type"),
	})
}

func TestAccStorageBlob_blockEmptyAccessTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blockEmptyAccessTier(data, blobs.Cool),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("access_tier").HasValue("Cool"),
			),
		},
		data.ImportStep("parallelism", "size", "type"),
		{
			Config: r.blockEmptyAccessTier(data, blobs.Hot),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("access_tier").HasValue("Hot"),
			),
		},
		{
			Config: r.blockEmptyAccessTier(data, blobs.Cool),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("access_tier").HasValue("Cool"),
			),
		},
	})
}

func TestAccStorageBlob_blockFromInlineContent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}
	content := "Wubba Lubba Dub Dubs"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blockFromInlineContent(data, content),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.blobMatchesContent(blobs.BlockBlob, []byte(content))),
			),
		},
		data.ImportStep("parallelism", "size", "source_content", "type"),
	})
}

func TestAccStorageBlob_blockFromInlineContentWithContentType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}
	content := `{"hello":"world"}`

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blockFromInlineContentWithContentType(data, content),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.blobMatchesContent(blobs.BlockBlob, []byte(content))),
			),
		},
		data.ImportStep("parallelism", "size", "source_content", "type"),
	})
}

func TestAccStorageBlob_blockFromPublicBlob(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blockFromPublicBlob(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "source_uri", "type"),
	})
}

func TestAccStorageBlob_blockFromPublicFile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blockFromPublicFile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "source_uri", "type"),
	})
}

func TestAccStorageBlob_blockFromExistingBlob(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blockFromExistingBlob(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "source_uri", "type"),
	})
}

func TestAccStorageBlob_blockFromLocalFile(t *testing.T) {
	sourceBlob, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Failed to create local source blob file")
	}

	if err := populateTempFile(sourceBlob); err != nil {
		t.Fatalf("Error populating temp file: %s", err)
	}
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blockFromLocalBlob(data, sourceBlob.Name()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.blobMatchesFile(blobs.BlockBlob, sourceBlob.Name())),
			),
		},
		data.ImportStep("parallelism", "size", "source", "type"),
	})
}

func TestAccStorageBlob_blockFromLocalFileWithContentMd5(t *testing.T) {
	sourceBlob, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Failed to create local source blob file")
	}

	if err := populateTempFile(sourceBlob); err != nil {
		t.Fatalf("Error populating temp file: %s", err)
	}
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.contentMd5ForLocalFile(data, sourceBlob.Name()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("example.vhd"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "source", sourceBlob.Name()),
			),
		},
		data.ImportStep("parallelism", "size", "source", "type"),
	})
}

func TestAccStorageBlob_cacheControl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cacheControl(data, "no-cache"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "type"),
		{
			Config: r.cacheControl(data, "max-age=3600"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "type"),
	})
}

func TestAccStorageBlob_contentType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.contentType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "type"),
		{
			Config: r.contentTypeUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "type"),
	})
}

func TestAccStorageBlob_contentTypePremium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.contentTypePremium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "type"),
	})
}

func TestAccStorageBlob_encryptionScope(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}
	content := "Wubba Lubba Dub Dubs"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.encryptionScope(data, content),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "source_content", "type"),
	})
}

func TestAccStorageBlob_encryptionScopeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}
	content := "Wubba Lubba Dub Dubs"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.encryptionScope(data, content),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "source_content", "type"),
		{
			Config: r.encryptionScopeUpdateMetadata(data, content),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "source_content", "type"),
		{
			Config: r.encryptionScopeUpdateProperties(data, content),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "source_content", "type"),
		{
			Config: r.encryptionScopeUpdateAccessTier(data, content),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "source_content", "type"),
	})
}

func TestAccStorageBlob_pageEmpty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.pageEmpty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "type"),
	})
}

func TestAccStorageBlob_pageEmptyPremium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.pageEmptyPremium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "type"),
	})
}

func TestAccStorageBlob_pageEmptyMetaData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.pageEmptyMetaData(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "type"),
	})
}

func TestAccStorageBlob_pageFromExistingBlob(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.pageFromExistingBlob(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "type", "source_uri"),
	})
}

func TestAccStorageBlob_pageFromLocalFile(t *testing.T) {
	sourceBlob, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Failed to create local source blob file")
	}

	if err := populateTempFile(sourceBlob); err != nil {
		t.Fatalf("Error populating temp file: %s", err)
	}
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.pageFromLocalBlob(data, sourceBlob.Name()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.blobMatchesFile(blobs.PageBlob, sourceBlob.Name())),
			),
		},
		data.ImportStep("parallelism", "size", "type", "source"),
	})
}

func TestAccStorageBlob_pageFromInlineContent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.pageFromInlineContent(data, 512),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "source_content", "type"),
		{
			Config:      r.pageFromInlineContent(data, 511),
			ExpectError: regexp.MustCompile(`"source" must be aligned to 512-byte boundary for "type" set to "Page"`),
		},
	})
}

func TestAccStorageBlob_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blockFromPublicBlob(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccStorageBlob_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "type"),
		{
			Config: r.updateUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "type"),
	})
}

func TestAccStorageBlob_archive(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")
	r := StorageBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.archive(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parallelism", "size", "type", "source_content"),
	})
}

func (r StorageBlobResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := blobs.ParseBlobID(state.ID, client.Storage.StorageDomainSuffix)
	if err != nil {
		return nil, err
	}
	account, err := client.Storage.FindAccount(ctx, client.Account.SubscriptionId, id.AccountId.AccountName)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, fmt.Errorf("unable to locate Account %q for Blob %q (Container %q)", id.AccountId.AccountName, id.BlobName, id.ContainerName)
	}
	blobsClient, err := client.Storage.BlobsDataPlaneClient(ctx, *account, client.Storage.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return nil, fmt.Errorf("building Blobs Client: %+v", err)
	}
	input := blobs.GetPropertiesInput{}
	resp, err := blobsClient.GetProperties(ctx, id.ContainerName, id.BlobName, input)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Blob %q (Container %q / Account %q): %+v", id.BlobName, id.ContainerName, id.AccountId.AccountName, err)
	}
	return utils.Bool(true), nil
}

func (r StorageBlobResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := blobs.ParseBlobID(state.ID, client.Storage.StorageDomainSuffix)
	if err != nil {
		return nil, err
	}
	account, err := client.Storage.FindAccount(ctx, client.Account.SubscriptionId, id.AccountId.AccountName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account %q for Blob %q (Container %q): %+v", id.AccountId.AccountName, id.BlobName, id.ContainerName, err)
	}
	blobsClient, err := client.Storage.BlobsDataPlaneClient(ctx, *account, client.Storage.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return nil, fmt.Errorf("building Blobs Client: %+v", err)
	}
	input := blobs.DeleteInput{
		DeleteSnapshots: false,
	}
	if _, err = blobsClient.Delete(ctx, id.ContainerName, id.BlobName, input); err != nil {
		return nil, fmt.Errorf("deleting Blob %q (Container %q / Account %q): %+v", id.BlobName, id.ContainerName, id.AccountId.AccountName, err)
	}
	return utils.Bool(true), nil
}

func (r StorageBlobResource) blobMatchesFile(kind blobs.BlobType, filePath string) func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
	expectedContents, err := os.ReadFile(filePath)
	if err != nil {
		return func(_ context.Context, _ *clients.Client, _ *pluginsdk.InstanceState) error {
			return fmt.Errorf("could not read local file for comparison: %v", err)
		}
	}

	return r.blobMatchesContent(kind, expectedContents)
}

func (r StorageBlobResource) blobMatchesContent(kind blobs.BlobType, expectedContents []byte) func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithDeadline(ctx, time.Now().Add(10*time.Minute))
			defer cancel()
		}

		name := state.Attributes["name"]
		containerName := state.Attributes["storage_container_name"]
		accountName := state.Attributes["storage_account_name"]

		account, err := clients.Storage.FindAccount(ctx, clients.Account.SubscriptionId, accountName)
		if err != nil {
			return fmt.Errorf("retrieving Account %q for Blob %q (Container %q): %s", accountName, name, containerName, err)
		}
		if account == nil {
			return fmt.Errorf("unable to locate Storage Account %q!", accountName)
		}

		client, err := clients.Storage.BlobsDataPlaneClient(ctx, *account, clients.Storage.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return fmt.Errorf("building Blobs Client: %s", err)
		}

		// first check the type
		getPropsInput := blobs.GetPropertiesInput{}
		props, err := client.GetProperties(ctx, containerName, name, getPropsInput)
		if err != nil {
			return fmt.Errorf("retrieving Properties for Blob %q (Container %q): %s", name, containerName, err)
		}

		if props.BlobType != kind {
			return fmt.Errorf("bad: blob type %q does not match expected type %q", props.BlobType, kind)
		}

		// then compare the content itself
		getInput := blobs.GetInput{}
		actualProps, err := client.Get(ctx, containerName, name, getInput)
		if err != nil {
			return fmt.Errorf("retrieving Blob %q (Container %q): %s", name, containerName, err)
		}

		actualContents := actualProps.Contents

		if actualContents == nil {
			return fmt.Errorf("bad: Storage Blob %q (storage container: %q) returned nil contents", name, containerName)
		}

		// retrieved blob content may have arbitrary newline appended
		if strings.TrimSpace(string(*actualContents)) != strings.TrimSpace(string(expectedContents)) {
			return fmt.Errorf("bad: Storage Blob %q (storage container: %q) does not match contents", name, containerName)
		}

		return nil
	}
}

func (r StorageBlobResource) appendEmpty(data acceptance.TestData) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Append"
}
`, template)
}

func (r StorageBlobResource) appendEmptyMetaData(data acceptance.TestData) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Append"

  metadata = {
    hello = "world"
  }
}
`, template)
}

func (r StorageBlobResource) blockEmpty(data acceptance.TestData) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
}
`, template)
}

func (r StorageBlobResource) blockEmptyAzureADAuth(data acceptance.TestData) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  storage_use_azuread = true
}

%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
}
`, template)
}

func (r StorageBlobResource) blockEmptyMetaData(data acceptance.TestData) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"

  metadata = {
    hello = "world"
  }
}
`, template)
}

func (r StorageBlobResource) blockEmptyAccessTier(data acceptance.TestData, accessTier blobs.AccessTier) string {
	template := r.templateBlockBlobStorage(data, "private")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  access_tier            = "%s"
}
`, template, string(accessTier))
}

func (r StorageBlobResource) blockFromInlineContent(data acceptance.TestData, content string) string {
	template := r.template(data, "blob")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "test" {
  name                   = "rick.morty"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  source_content         = <<EOT
%s
EOT
}
`, template, content)
}

func (r StorageBlobResource) blockFromInlineContentWithContentType(data acceptance.TestData, content string) string {
	template := r.template(data, "blob")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "test" {
  name                   = "blob.json"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  content_type           = "application/json"
  source_content         = <<EOT
%s
EOT
}
`, template, content)
}

func (r StorageBlobResource) blockFromPublicBlob(data acceptance.TestData) string {
	template := r.template(data, "blob")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "source" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  source_uri             = "http://old-releases.ubuntu.com/releases/bionic/ubuntu-18.04-desktop-amd64.iso"
  content_type           = "application/x-iso9660-image"
}

resource "azurerm_storage_container" "second" {
  name                  = "second"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_storage_blob" "test" {
  name                   = "copied.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.second.name
  type                   = "Block"
  source_uri             = azurerm_storage_blob.source.id
  content_type           = azurerm_storage_blob.source.content_type
}
`, template)
}

func (r StorageBlobResource) blockFromPublicFile(data acceptance.TestData) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  source_uri             = "http://old-releases.ubuntu.com/releases/bionic/ubuntu-18.04-desktop-amd64.iso"
  content_type           = "application/x-iso9660-image"
}
`, template)
}

func (r StorageBlobResource) blockFromExistingBlob(data acceptance.TestData) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "source" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  source_uri             = "http://old-releases.ubuntu.com/releases/bionic/ubuntu-18.04-desktop-amd64.iso"
  content_type           = "application/x-iso9660-image"
}

resource "azurerm_storage_blob" "test" {
  name                   = "copied.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  source_uri             = azurerm_storage_blob.source.id
  content_type           = azurerm_storage_blob.source.content_type
}
`, template)
}

func (r StorageBlobResource) blockFromLocalBlob(data acceptance.TestData, fileName string) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  source                 = "%s"
}
`, template, fileName)
}

func (r StorageBlobResource) contentMd5ForLocalFile(data acceptance.TestData, fileName string) string {
	template := r.template(data, "blob")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  source                 = "%s"
  content_md5            = "${filemd5("%s")}"
}
`, template, fileName, fileName)
}

func (r StorageBlobResource) contentType(data acceptance.TestData) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.ext"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  size                   = 5120
  content_type           = "image/png"
}
`, template)
}

func (r StorageBlobResource) contentTypePremium(data acceptance.TestData) string {
	template := r.templatePremium(data, "private")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.ext"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  size                   = 5120
  content_type           = "image/png"
}
`, template)
}

func (r StorageBlobResource) contentTypeUpdated(data acceptance.TestData) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.ext"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  size                   = 5120
  content_type           = "image/gif"
}
`, template)
}

func (r StorageBlobResource) encryptionScope(data acceptance.TestData, content string) string {
	template := r.template(data, "blob")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_storage_encryption_scope" "test" {
  name               = "acctestEScontainer%[2]d"
  storage_account_id = azurerm_storage_account.test.id
  source             = "Microsoft.Storage"
}

resource "azurerm_storage_blob" "test" {
  name                   = "rick.morty"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  encryption_scope       = azurerm_storage_encryption_scope.test.name
  source_content         = <<EOT
%[3]s
EOT
}
`, template, data.RandomInteger, content)
}

func (r StorageBlobResource) encryptionScopeUpdateMetadata(data acceptance.TestData, content string) string {
	template := r.template(data, "blob")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_storage_encryption_scope" "test" {
  name               = "acctestEScontainer%[2]d"
  storage_account_id = azurerm_storage_account.test.id
  source             = "Microsoft.Storage"
}

resource "azurerm_storage_blob" "test" {
  name                   = "rick.morty"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  encryption_scope       = azurerm_storage_encryption_scope.test.name
  source_content         = <<EOT
%[3]s
EOT

  metadata = {
    hello = "world"
  }
}
`, template, data.RandomInteger, content)
}

func (r StorageBlobResource) encryptionScopeUpdateProperties(data acceptance.TestData, content string) string {
	template := r.template(data, "blob")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_storage_encryption_scope" "test" {
  name               = "acctestEScontainer%[2]d"
  storage_account_id = azurerm_storage_account.test.id
  source             = "Microsoft.Storage"
}

resource "azurerm_storage_blob" "test" {
  name                   = "rick.morty"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  encryption_scope       = azurerm_storage_encryption_scope.test.name
  source_content         = <<EOT
%[3]s
EOT

  metadata = {
    hello = "world"
  }

  content_type = "text/plain"
}
`, template, data.RandomInteger, content)
}

func (r StorageBlobResource) encryptionScopeUpdateAccessTier(data acceptance.TestData, content string) string {
	template := r.template(data, "blob")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_storage_encryption_scope" "test" {
  name               = "acctestEScontainer%[2]d"
  storage_account_id = azurerm_storage_account.test.id
  source             = "Microsoft.Storage"
}

resource "azurerm_storage_blob" "test" {
  name                   = "rick.morty"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  encryption_scope       = azurerm_storage_encryption_scope.test.name
  source_content         = <<EOT
%[3]s
EOT

  metadata = {
    hello = "world"
  }

  content_type = "text/plain"
  access_tier  = "Hot"
}
`, template, data.RandomInteger, content)
}

func (r StorageBlobResource) pageEmpty(data acceptance.TestData) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  size                   = 5120
}
`, template)
}

func (r StorageBlobResource) pageEmptyPremium(data acceptance.TestData) string {
	template := r.templatePremium(data, "private")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  size                   = 5120
}
`, template)
}

func (r StorageBlobResource) pageEmptyMetaData(data acceptance.TestData) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  size                   = 5120

  metadata = {
    hello = "world"
  }
}
`, template)
}

func (r StorageBlobResource) pageFromExistingBlob(data acceptance.TestData) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "source" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  size                   = 5120
  content_type           = "application/x-iso9660-image"
}

resource "azurerm_storage_blob" "test" {
  name                   = "copied.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  source_uri             = azurerm_storage_blob.source.id
  content_type           = azurerm_storage_blob.source.content_type
}
`, template)
}

func (r StorageBlobResource) pageFromLocalBlob(data acceptance.TestData, fileName string) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  source                 = "%s"
}
`, template, fileName)
}

func (r StorageBlobResource) pageFromInlineContent(data acceptance.TestData, length int) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "test" {
  name                   = "rick.morty"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  source_content         = join("", [for i in range(0, %d) : "a"])
}
`, template, length)
}

func (r StorageBlobResource) requiresImport(data acceptance.TestData) string {
	template := r.blockFromPublicBlob(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "import" {
  name                   = azurerm_storage_blob.test.name
  storage_account_name   = azurerm_storage_blob.test.storage_account_name
  storage_container_name = azurerm_storage_blob.test.storage_container_name
  type                   = azurerm_storage_blob.test.type
  size                   = azurerm_storage_blob.test.size
}
`, template)
}

func (r StorageBlobResource) update(data acceptance.TestData) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  size                   = 5120
  content_type           = "vnd/panda+pops"
  metadata = {
    hello = "world"
  }
}
`, template)
}

func (r StorageBlobResource) updateUpdated(data acceptance.TestData) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  size                   = 5120
  content_type           = "vnd/mountain-mover-3000"
  metadata = {
    hello = "world"
    panda = "pops"
  }
}
`, template)
}

func (r StorageBlobResource) cacheControl(data acceptance.TestData, cacheControl string) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.ext"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  size                   = 5120
  content_type           = "image/png"
  cache_control          = "%s"
}
`, template, cacheControl)
}

func (r StorageBlobResource) template(data acceptance.TestData, accessLevel string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                            = "acctestacc%s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, accessLevel)
}

func (r StorageBlobResource) templateBlockBlobStorage(data acceptance.TestData, accessLevel string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                            = "acctestacc%s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_kind                    = "StorageV2"
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, accessLevel)
}

func (r StorageBlobResource) templatePremium(data acceptance.TestData, accessLevel string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                            = "acctestacc%s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_tier                    = "Premium"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, accessLevel)
}

func (r StorageBlobResource) archive(data acceptance.TestData) string {
	template := r.template(data, "private")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_blob" "test" {
  name                   = "rick.morty"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  source_content         = "Wubba Lubba Dub Dub"
  access_tier            = "Archive"
}
`, template)
}

func populateTempFile(input *os.File) error {
	if err := input.Truncate(25*1024*1024 + 512); err != nil {
		return fmt.Errorf("Failed to truncate file to 25M")
	}

	for i := int64(0); i < 20; i += 2 {
		randomBytes := make([]byte, 1*1024*1024)
		if _, err := rand.Read(randomBytes); err != nil {
			return fmt.Errorf("Failed to read random bytes")
		}

		if _, err := input.WriteAt(randomBytes, i*1024*1024); err != nil {
			return fmt.Errorf("Failed to write random bytes to file")
		}
	}

	randomBytes := make([]byte, 5*1024*1024)
	if _, err := rand.Read(randomBytes); err != nil {
		return fmt.Errorf("Failed to read random bytes")
	}

	if _, err := input.WriteAt(randomBytes, 20*1024*1024); err != nil {
		return fmt.Errorf("Failed to write random bytes to file")
	}

	if err := input.Close(); err != nil {
		return fmt.Errorf("Failed to close source blob")
	}

	return nil
}
