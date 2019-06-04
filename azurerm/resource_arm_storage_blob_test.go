package azurerm

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"strings"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMStorageBlob_basic(t *testing.T) {
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
				Config: testAccAzureRMStorageBlob_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "metadata.test", "value1"),
					resource.TestCheckResourceAttr(resourceName, "metadata.test2", "value2"),
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
func TestAccAzureRMStorageBlob_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
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
				Config: testAccAzureRMStorageBlob_basic(ri, rs, location),
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

func TestAccAzureRMStorageBlob_disappears(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMStorageBlob_basic(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
					testCheckAzureRMStorageBlobDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMStorageBlobBlock_source(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs1 := strings.ToLower(acctest.RandString(11))
	sourceBlob, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Failed to create local source blob file")
	}

	_, err = io.CopyN(sourceBlob, rand.Reader, 25*1024*1024)
	if err != nil {
		t.Fatalf("Failed to write random test to source blob")
	}

	err = sourceBlob.Close()
	if err != nil {
		t.Fatalf("Failed to close source blob")
	}

	config := testAccAzureRMStorageBlobBlock_source(ri, rs1, sourceBlob.Name(), testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobMatchesFile("azurerm_storage_blob.source", storage.BlobTypeBlock, sourceBlob.Name()),
				),
			},
		},
	})
}

func TestAccAzureRMStorageBlobPage_source(t *testing.T) {
	resourceName := "azurerm_storage_blob.source"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	sourceBlob, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Failed to create local source blob file")
	}

	err = sourceBlob.Truncate(25*1024*1024 + 512)
	if err != nil {
		t.Fatalf("Failed to truncate file to 25M")
	}

	for i := int64(0); i < 20; i = i + 2 {
		randomBytes := make([]byte, 1*1024*1024)
		_, err = rand.Read(randomBytes)
		if err != nil {
			t.Fatalf("Failed to read random bytes")
		}

		_, err = sourceBlob.WriteAt(randomBytes, i*1024*1024)
		if err != nil {
			t.Fatalf("Failed to write random bytes to file")
		}
	}

	randomBytes := make([]byte, 5*1024*1024)
	_, err = rand.Read(randomBytes)
	if err != nil {
		t.Fatalf("Failed to read random bytes")
	}

	_, err = sourceBlob.WriteAt(randomBytes, 20*1024*1024)
	if err != nil {
		t.Fatalf("Failed to write random bytes to file")
	}

	err = sourceBlob.Close()
	if err != nil {
		t.Fatalf("Failed to close source blob")
	}

	config := testAccAzureRMStorageBlobPage_source(ri, rs, sourceBlob.Name(), testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobMatchesFile(resourceName, storage.BlobTypePage, sourceBlob.Name()),
				),
			},
		},
	})
}

func TestAccAzureRMStorageBlob_source_uri(t *testing.T) {
	resourceName := "azurerm_storage_blob.destination"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	sourceBlob, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Failed to create local source blob file")
	}

	_, err = io.CopyN(sourceBlob, rand.Reader, 25*1024*1024)
	if err != nil {
		t.Fatalf("Failed to write random test to source blob")
	}

	err = sourceBlob.Close()
	if err != nil {
		t.Fatalf("Failed to close source blob")
	}

	config := testAccAzureRMStorageBlob_source_uri(ri, rs, sourceBlob.Name(), testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobMatchesFile(resourceName, storage.BlobTypeBlock, sourceBlob.Name()),
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

func TestAccAzureRMStorageBlobBlock_blockContentType(t *testing.T) {
	resourceName := "azurerm_storage_blob.source"
	ri := tf.AccRandTimeInt()
	rs1 := strings.ToLower(acctest.RandString(11))
	sourceBlob, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Failed to create local source blob file")
	}

	_, err = io.CopyN(sourceBlob, rand.Reader, 25*1024*1024)
	if err != nil {
		t.Fatalf("Failed to write random test to source blob")
	}

	err = sourceBlob.Close()
	if err != nil {
		t.Fatalf("Failed to close source blob")
	}

	config := testAccAzureRMStorageBlobPage_blockContentType(ri, rs1, testLocation(), sourceBlob.Name(), "text/plain")
	updateConfig := testAccAzureRMStorageBlobPage_blockContentType(ri, rs1, testLocation(), sourceBlob.Name(), "text/vnd.terraform.acctest.tmpfile")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "content_type", "text/plain"),
				),
			},
			{
				Config: updateConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "content_type", "text/vnd.terraform.acctest.tmpfile"),
				),
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
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for storage blob: %s", name)
		}

		armClient := testAccProvider.Meta().(*ArmClient)
		ctx := armClient.StopContext
		blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, resourceGroup, storageAccountName)
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
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for storage blob: %s", name)
		}

		armClient := testAccProvider.Meta().(*ArmClient)
		ctx := armClient.StopContext
		blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, resourceGroup, storageAccountName)
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
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for storage blob: %s", name)
		}

		armClient := testAccProvider.Meta().(*ArmClient)
		ctx := armClient.StopContext
		blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, resourceGroup, storageAccountName)
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
		blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, resourceGroup, storageAccountName)
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

func testAccAzureRMStorageBlob_basic(rInt int, rString string, location string) string {
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

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_storage_blob" "test" {
  name = "herpderp1.vhd"

  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"

  type = "page"
  size = 5120
  
  metadata = {
    test = "value1"
    test2 = "value2"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageBlob_requiresImport(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageBlob_basic(rInt, rString, location)
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

func testAccAzureRMStorageBlobBlock_source(rInt int, rString string, sourceBlobName string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "source" {
  name                     = "acctestacc%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "source" {
  name                  = "source"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.source.name}"
  container_access_type = "blob"
}

resource "azurerm_storage_blob" "source" {
  name = "source.vhd"

  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.source.name}"
  storage_container_name = "${azurerm_storage_container.source.name}"

  type        = "block"
  source      = "%s"
  parallelism = 4
  attempts    = 2
}
`, rInt, location, rString, sourceBlobName)
}

func testAccAzureRMStorageBlobPage_source(rInt int, rString string, sourceBlobName string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "source" {
  name                     = "acctestacc%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "source" {
  name                  = "source"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.source.name}"
  container_access_type = "blob"
}

resource "azurerm_storage_blob" "source" {
  name = "source.vhd"

  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.source.name}"
  storage_container_name = "${azurerm_storage_container.source.name}"

  type        = "page"
  source      = "%s"
  parallelism = 3
  attempts    = 3
}
`, rInt, location, rString, sourceBlobName)
}

func testAccAzureRMStorageBlob_source_uri(rInt int, rString string, sourceBlobName string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "source" {
  name                     = "acctestacc%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "source" {
  name                  = "source"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.source.name}"
  container_access_type = "blob"
}

resource "azurerm_storage_blob" "source" {
  name = "source.vhd"

  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.source.name}"
  storage_container_name = "${azurerm_storage_container.source.name}"

  type        = "block"
  source      = "%s"
  parallelism = 4
  attempts    = 2
}

resource "azurerm_storage_blob" "destination" {
  name                   = "destination.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.source.name}"
  storage_container_name = "${azurerm_storage_container.source.name}"
  source_uri             = "${azurerm_storage_blob.source.url}"
  type                   = "block"
}
`, rInt, location, rString, sourceBlobName)
}

func testAccAzureRMStorageBlobPage_blockContentType(rInt int, rString, location string, sourceBlobName, contentType string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "source" {
  name                     = "acctestacc%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "source" {
  name                  = "source"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.source.name}"
  container_access_type = "blob"
}

resource "azurerm_storage_blob" "source" {
  name = "source.vhd"

  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.source.name}"
  storage_container_name = "${azurerm_storage_container.source.name}"

  type         = "page"
  source       = "%s"
  content_type = "%s"
  parallelism  = 3
  attempts     = 3
}
`, rInt, location, rString, sourceBlobName, contentType)
}
