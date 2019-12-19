package azurerm

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/blob/blobs"
)

func TestAccAzureRMStorageBlob_disappears(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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

func TestAccAzureRMStorageBlob_blockEmptyAccessTier(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockEmptyAccessTier(ri, rs, location, blobs.Cool),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "access_tier", "Cool"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "type"},
			},
			{
				Config: testAccAzureRMStorageBlob_blockEmptyAccessTier(ri, rs, location, blobs.Hot),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "access_tier", "Hot"),
				),
			},
			{
				Config: testAccAzureRMStorageBlob_blockEmptyAccessTier(ri, rs, location, blobs.Cool),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "access_tier", "Cool"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockFromInlineContent(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockFromInlineContent(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "source_content", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockFromPublicBlob(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "source_uri", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockFromPublicFile(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "source_uri", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockFromExistingBlob(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "source_uri", "type"},
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
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockFromLocalBlob(ri, rs, location, sourceBlob.Name()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
					testCheckAzureRMStorageBlobMatchesFile(resourceName, blobs.BlockBlob, sourceBlob.Name()),
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
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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

func TestAccAzureRMStorageBlob_contentTypePremium(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_contentTypePremium(ri, rs, location),
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
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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

func TestAccAzureRMStorageBlob_pageEmptyPremium(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_pageEmptyPremium(ri, rs, location),
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
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
				ImportStateVerifyIgnore: []string{"attempts", "parallelism", "size", "source_uri", "type"},
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
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_pageFromLocalBlob(ri, rs, location, sourceBlob.Name()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
					testCheckAzureRMStorageBlobMatchesFile(resourceName, blobs.PageBlob, sourceBlob.Name()),
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
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockFromPublicBlob(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMStorageBlob_requiresImport(ri, rs, location),
				ExpectError: acceptance.RequiresImportError("azurerm_storage_blob"),
			},
		},
	})
}

func TestAccAzureRMStorageBlob_update(t *testing.T) {
	resourceName := "azurerm_storage_blob.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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

		storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		name := rs.Primary.Attributes["name"]
		containerName := rs.Primary.Attributes["storage_container_name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for Blob %q (Container %q): %s", accountName, name, containerName, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
		}

		client, err := storageClient.BlobsClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building Blobs Client: %s", err)
		}

		input := blobs.GetPropertiesInput{}
		resp, err := client.GetProperties(ctx, accountName, containerName, name, input)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Blob %q (Container %q / Account %q / Resource Group %q) does not exist", name, containerName, accountName, account.ResourceGroup)
			}

			return fmt.Errorf("Bad: Get on BlobsClient: %+v", err)
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

		storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		name := rs.Primary.Attributes["name"]
		containerName := rs.Primary.Attributes["storage_container_name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for Blob %q (Container %q): %s", accountName, name, containerName, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
		}

		client, err := storageClient.BlobsClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building Blobs Client: %s", err)
		}

		input := blobs.DeleteInput{
			DeleteSnapshots: false,
		}
		if _, err := client.Delete(ctx, accountName, containerName, name, input); err != nil {
			return fmt.Errorf("Error deleting Blob %q (Container %q / Account %q / Resource Group %q): %s", name, containerName, accountName, account.ResourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMStorageBlobMatchesFile(resourceName string, kind blobs.BlobType, filePath string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		name := rs.Primary.Attributes["name"]
		containerName := rs.Primary.Attributes["storage_container_name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for Blob %q (Container %q): %s", accountName, name, containerName, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
		}

		client, err := storageClient.BlobsClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building Blobs Client: %s", err)
		}

		// first check the type
		getPropsInput := blobs.GetPropertiesInput{}
		props, err := client.GetProperties(ctx, accountName, containerName, name, getPropsInput)
		if err != nil {
			return fmt.Errorf("Error retrieving Properties for Blob %q (Container %q): %s", name, containerName, err)
		}

		if props.BlobType != kind {
			return fmt.Errorf("Bad: blob type %q does not match expected type %q", props.BlobType, kind)
		}

		// then compare the content itself
		getInput := blobs.GetInput{}
		actualProps, err := client.Get(ctx, accountName, containerName, name, getInput)
		if err != nil {
			return fmt.Errorf("Error retrieving Blob %q (Container %q): %s", name, containerName, err)
		}

		actualContents := actualProps.Contents

		// local file for comparison
		expectedContents, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}

		if string(actualContents) != string(expectedContents) {
			return fmt.Errorf("Bad: Storage Blob %q (storage container: %q) does not match contents", name, containerName)
		}

		return nil
	}
}

func testCheckAzureRMStorageBlobDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_blob" {
			continue
		}

		storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		name := rs.Primary.Attributes["name"]
		containerName := rs.Primary.Attributes["storage_container_name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for Blob %q (Container %q): %s", accountName, name, containerName, err)
		}
		if account == nil {
			return nil
		}

		client, err := storageClient.BlobsClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building Blobs Client: %s", err)
		}

		input := blobs.GetPropertiesInput{}
		props, err := client.GetProperties(ctx, accountName, containerName, name, input)
		if err != nil {
			if !utils.ResponseWasNotFound(props.Response) {
				return fmt.Errorf("Error retrieving Blob %q (Container %q / Account %q): %s", name, containerName, accountName, err)
			}
		}

		if utils.ResponseWasNotFound(props.Response) {
			return nil
		}

		return fmt.Errorf("Bad: Storage Blob %q (Storage Container: %q) still exists", name, containerName)
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

func testAccAzureRMStorageBlob_blockEmptyAccessTier(rInt int, rString string, location string, accessTier blobs.AccessTier) string {
	template := testAccAzureRMStorageBlob_templateBlockBlobStorage(rInt, rString, location, "private")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "block"
  access_tier            = "%s"
}
`, template, string(accessTier))
}

func testAccAzureRMStorageBlob_blockFromInlineContent(rInt int, rString, location string) string {
	template := testAccAzureRMStorageBlob_template(rInt, rString, location, "blob")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob" "test" {
  name                   = "rick.morty"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "block"
  source_content         = "Wubba Lubba Dub Dub"
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

func testAccAzureRMStorageBlob_contentTypePremium(rInt int, rString, location string) string {
	template := testAccAzureRMStorageBlob_templatePremium(rInt, rString, location, "private")
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

func testAccAzureRMStorageBlob_pageEmptyPremium(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageBlob_templatePremium(rInt, rString, location, "private")
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
	template := testAccAzureRMStorageBlob_blockFromPublicBlob(rInt, rString, location)
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
  metadata = {
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
  metadata = {
    hello = "world"
    panda = "pops"
  }
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

func testAccAzureRMStorageBlob_templateBlockBlobStorage(rInt int, rString, location, accessLevel string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_kind             = "StorageV2"
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

func testAccAzureRMStorageBlob_templatePremium(rInt int, rString, location, accessLevel string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Premium"
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
