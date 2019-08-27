package azurerm

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

// TODO: with the new SDK: changing the Tier of Blobs. Content type for Block blobs

var supportsNewStorageFeatures = false

func TestAccAzureRMStorageBlob_disappears(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockEmpty(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
					testCheckAzureRMStorageBlobDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMStorageBlob_appendEmpty(t *testing.T) {
	if !supportsNewStorageFeatures {
		t.Skip("Resource doesn't support Append Blobs Yet..")
	}

	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_appendEmpty(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_appendEmptyMetaData(t *testing.T) {
	if !supportsNewStorageFeatures {
		t.Skip("Resource doesn't support Append Blobs Yet..")
	}

	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_appendEmptyMetaData(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockEmpty(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockEmpty(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockEmptyMetaData(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockEmptyMetaData(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockFromPublicBlob(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockFromPublicBlob(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockFromPublicFile(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockFromPublicFile(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockFromExistingBlob(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockFromExistingBlob(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockFromLocalFile(t *testing.T) {
	sourceBlob, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Failed to create local source blob file")
	}

	if err := testAccAzureRMStorageBlob_populateTempFile(sourceBlob); err != nil {
		t.Fatalf("Error populating temp file: %s", err)
	}

	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockFromLocalBlob(ri, rs, location, sourceBlob.Name()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
					testCheckAzureRMStorageBlobMatchesFile(resourceName, storage.BlobTypeBlock, sourceBlob.Name()),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "source", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_contentType(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_contentType(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "type"},
			},
			{
				Config: testAccAzureRMStorageBlob_contentTypeUpdated(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_pageEmpty(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_pageEmpty(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_pageEmptyMetaData(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_pageEmptyMetaData(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_pageFromExistingBlob(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_pageFromExistingBlob(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_pageFromLocalFile(t *testing.T) {
	sourceBlob, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Failed to create local source blob file")
	}

	if err := testAccAzureRMStorageBlob_populateTempFile(sourceBlob); err != nil {
		t.Fatalf("Error populating temp file: %s", err)
	}

	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_pageFromLocalBlob(ri, rs, location, sourceBlob.Name()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
					testCheckAzureRMStorageBlobMatchesFile(resourceName, storage.BlobTypePage, sourceBlob.Name()),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "source", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_requiresImport(t *testing.T) {
	if features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockEmpty(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMStorageBlob_requiresImport(ri, rs, location),
				ExpectError: testRequiresImportError("azurerm_storage_blob"),
			},
		},
	})
}

func TestAccAzureRMStorageBlob_update(t *testing.T) {
	if !supportsNewStorageFeatures {
		t.Skip("Current implementation doesn't support updating the Content Type..")
	}

	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_update(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "type"},
			},
			{
				Config: testAccAzureRMStorageBlob_updateUpdated(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "type"},
			},
		},
	})
}

func testCheckAzureRMStorageBlobExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		storageAccountName := rs.Primary.Attributes["storage_account_name"]
		storageContainerName := rs.Primary.Attributes["storage_container_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		armClient := testAccProvider.Meta().(*ArmClient)
		ctx := armClient.StopContext
		blobClient, accountExists, err := armClient.storage.LegacyBlobClient(ctx, resourceGroup, storageAccountName)
		if err != nil {
			return err
		}
		if !accountExists {
			return fmt.Errorf("Bad: Storage Account %q does not exist", storageAccountName)
		}

		container := blobClient.GetContainerReference(storageContainerName)
		blob := container.GetBlobReference(name)
		exists, err := blob.Exists()
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("Bad: Storage Blob %q (storage container: %q) does not exist", name, storageContainerName)
		}

		return nil
	}
}

func testCheckAzureRMStorageBlobDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		storageAccountName := rs.Primary.Attributes["storage_account_name"]
		storageContainerName := rs.Primary.Attributes["storage_container_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		armClient := testAccProvider.Meta().(*ArmClient)
		ctx := armClient.StopContext
		blobClient, accountExists, err := armClient.storage.LegacyBlobClient(ctx, resourceGroup, storageAccountName)
		if err != nil {
			return err
		}
		if !accountExists {
			return fmt.Errorf("Bad: Storage Account %q does not exist", storageAccountName)
		}

		container := blobClient.GetContainerReference(storageContainerName)
		blob := container.GetBlobReference(name)
		options := &storage.DeleteBlobOptions{}
		_, err = blob.DeleteIfExists(options)
		return err
	}
}

func testCheckAzureRMStorageBlobMatchesFile(resourceName string, kind storage.BlobType, filePath string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		storageAccountName := rs.Primary.Attributes["storage_account_name"]
		storageContainerName := rs.Primary.Attributes["storage_container_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		armClient := testAccProvider.Meta().(*ArmClient)
		ctx := armClient.StopContext
		blobClient, accountExists, err := armClient.storage.LegacyBlobClient(ctx, resourceGroup, storageAccountName)
		if err != nil {
			return err
		}
		if !accountExists {
			return fmt.Errorf("Bad: Storage Account %q does not exist", storageAccountName)
		}

		containerReference := blobClient.GetContainerReference(storageContainerName)
		blobReference := containerReference.GetBlobReference(name)
		propertyOptions := &storage.GetBlobPropertiesOptions{}
		err = blobReference.GetProperties(propertyOptions)
		if err != nil {
			return err
		}

		properties := blobReference.Properties

		if properties.BlobType != kind {
			return fmt.Errorf("Bad: blob type %q does not match expected type %q", properties.BlobType, kind)
		}

		getOptions := &storage.GetBlobOptions{}
		blob, err := blobReference.Get(getOptions)
		if err != nil {
			return err
		}

		contents, err := ioutil.ReadAll(blob)
		if err != nil {
			return err
		}
		defer blob.Close()

		expectedContents, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}

		if string(contents) != string(expectedContents) {
			return fmt.Errorf("Bad: Storage Blob %q (storage container: %q) does not match contents", name, storageContainerName)
		}

		return nil
	}
}

func testCheckAzureRMStorageBlobDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_blob" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		storageAccountName := rs.Primary.Attributes["storage_account_name"]
		storageContainerName := rs.Primary.Attributes["storage_container_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for storage blob: %s", name)
		}

		armClient := testAccProvider.Meta().(*ArmClient)
		ctx := armClient.StopContext
		blobClient, accountExists, err := armClient.storage.LegacyBlobClient(ctx, resourceGroup, storageAccountName)
		if err != nil {
			return nil
		}
		if !accountExists {
			return nil
		}

		container := blobClient.GetContainerReference(storageContainerName)
		blob := container.GetBlobReference(name)
		exists, err := blob.Exists()
		if err != nil {
			return nil
		}

		if exists {
			return fmt.Errorf("Bad: Storage Blob %q (storage container: %q) still exists", name, storageContainerName)
		}
	}

	return nil
}

func testAccAzureRMStorageBlob_appendEmpty(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageBlob_template(rInt, rString, location, "private")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "append"
}
`, template)
}

func testAccAzureRMStorageBlob_appendEmptyMetaData(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageBlob_template(rInt, rString, location, "private")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "append"

  metadata = {
    hello = "world"
  }
}
`, template)
}

func testAccAzureRMStorageBlob_blockEmpty(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageBlob_template(rInt, rString, location, "private")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "block"
}
`, template)
}

func testAccAzureRMStorageBlob_blockEmptyMetaData(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageBlob_template(rInt, rString, location, "private")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "block"

  metadata = {
    hello = "world"
  }
}
`, template)
}

func testAccAzureRMStorageBlob_blockFromPublicBlob(rInt int, rString, location string) string {
	template := testAccAzureRMStorageBlob_template(rInt, rString, location, "blob")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "source" {
  name                   = "example.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "block"
  source_uri             = "http://releases.ubuntu.com/18.04.3/ubuntu-18.04.3-desktop-amd64.iso"
  content_type           = "application/x-iso9660-image"
}

resource "azurerm_storage_container" "second" {
  name                  = "second"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_storage_blob" "test" {
  name                   = "copied.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.second.name}"
  type                   = "block"
  source_uri             = "${azurerm_storage_blob.source.id}"
  content_type           = "${azurerm_storage_blob.source.content_type}"
}
`, template)
}

func testAccAzureRMStorageBlob_blockFromPublicFile(rInt int, rString, location string) string {
	template := testAccAzureRMStorageBlob_template(rInt, rString, location, "private")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "block"
  source_uri             = "http://releases.ubuntu.com/18.04.3/ubuntu-18.04.3-desktop-amd64.iso"
  content_type           = "application/x-iso9660-image"
}
`, template)
}

func testAccAzureRMStorageBlob_blockFromExistingBlob(rInt int, rString, location string) string {
	template := testAccAzureRMStorageBlob_template(rInt, rString, location, "private")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "source" {
  name                   = "example.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "block"
  source_uri             = "http://releases.ubuntu.com/18.04.3/ubuntu-18.04.3-desktop-amd64.iso"
  content_type           = "application/x-iso9660-image"
}

resource "azurerm_storage_blob" "test" {
  name                   = "copied.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "block"
  source_uri             = "${azurerm_storage_blob.source.id}"
  content_type           = "${azurerm_storage_blob.source.content_type}"
}
`, template)
}

func testAccAzureRMStorageBlob_blockFromLocalBlob(rInt int, rString, location, fileName string) string {
	template := testAccAzureRMStorageBlob_template(rInt, rString, location, "private")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "block"
  source                 = "%s"
}
`, template, fileName)
}

func testAccAzureRMStorageBlob_contentType(rInt int, rString, location string) string {
	template := testAccAzureRMStorageBlob_template(rInt, rString, location, "private")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.ext"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "page"
  size                   = 5120
  content_type           = "image/png"
}
`, template)
}

func testAccAzureRMStorageBlob_contentTypeUpdated(rInt int, rString, location string) string {
	template := testAccAzureRMStorageBlob_template(rInt, rString, location, "private")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.ext"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "page"
  size                   = 5120
  content_type           = "image/gif"
}
`, template)
}

func testAccAzureRMStorageBlob_pageEmpty(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageBlob_template(rInt, rString, location, "private")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "page"
  size                   = 5120
}
`, template)
}

func testAccAzureRMStorageBlob_pageEmptyMetaData(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageBlob_template(rInt, rString, location, "private")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "page"
  size                   = 5120

  metadata = {
    hello = "world"
  }
}
`, template)
}

func testAccAzureRMStorageBlob_pageFromExistingBlob(rInt int, rString, location string) string {
	template := testAccAzureRMStorageBlob_template(rInt, rString, location, "private")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "source" {
  name                   = "example.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "page"
  size                   = 5120
  content_type           = "application/x-iso9660-image"
}

resource "azurerm_storage_blob" "test" {
  name                   = "copied.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "page"
  source_uri             = "${azurerm_storage_blob.source.id}"
  content_type           = "${azurerm_storage_blob.source.content_type}"
}
`, template)
}

func testAccAzureRMStorageBlob_pageFromLocalBlob(rInt int, rString, location, fileName string) string {
	template := testAccAzureRMStorageBlob_template(rInt, rString, location, "private")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "page"
  source                 = "%s"
}
`, template, fileName)
}

func testAccAzureRMStorageBlob_populateTempFile(input *os.File) error {
	if err := input.Truncate(25*1024*1024 + 512); err != nil {
		return fmt.Errorf("Failed to truncate file to 25M")
	}

	for i := int64(0); i < 20; i = i + 2 {
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

func testAccAzureRMStorageBlob_requiresImport(rInt int, rString, location string) string {
	template := testAccAzureRMStorageBlob_blockEmpty(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "import" {
  name                   = "${azurerm_storage_blob.test.name}"
  resource_group_name    = "${azurerm_storage_blob.test.resource_group_name}"
  storage_account_name   = "${azurerm_storage_blob.test.storage_account_name}"
  storage_container_name = "${azurerm_storage_blob.test.storage_container_name}"
  type                   = "${azurerm_storage_blob.test.type}"
  size                   = "${azurerm_storage_blob.test.size}"
}
`, template)
}

func testAccAzureRMStorageBlob_template(rInt int, rString, location, accessLevel string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "%s"
}
`, rInt, location, rString, accessLevel)
}

func testAccAzureRMStorageBlob_update(rInt int, rString, location string) string {
	template := testAccAzureRMStorageBlob_template(rInt, rString, location, "private")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "block"
  size                   = 5120
  content_type           = "vnd/panda+pops"
  metadata               = {
    hello = "world"
  }
}
`, template)
}

func testAccAzureRMStorageBlob_updateUpdated(rInt int, rString, location string) string {
	template := testAccAzureRMStorageBlob_template(rInt, rString, location, "private")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "block"
  size                   = 5120
  content_type           = "vnd/mountain-mover-3000"
  metadata               = {
    hello = "world"
    panda = "pops"
  }
}
`, template)
}
