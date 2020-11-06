package tests

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/blobs"
)

type testChangeableSetPropertiesInput struct {
	CacheControl       string
	ContentType        string
	ContentEncoding    string
	ContentLanguage    string
	ContentDisposition string
}

func TestAccAzureRMStorageBlob_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockEmpty(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
					testCheckAzureRMStorageBlobDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMStorageBlob_appendEmpty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_appendEmpty(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_appendEmptyMetaData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_appendEmptyMetaData(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockEmpty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockEmpty(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockEmptyAzureADAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockEmptyAzureADAuth(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockEmptyMetaData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockEmptyMetaData(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockEmptyAccessTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockEmptyAccessTier(data, blobs.Cool),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "access_tier", "Cool"),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
			{
				Config: testAccAzureRMStorageBlob_blockEmptyAccessTier(data, blobs.Hot),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "access_tier", "Hot"),
				),
			},
			{
				Config: testAccAzureRMStorageBlob_blockEmptyAccessTier(data, blobs.Cool),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "access_tier", "Cool"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockFromInlineContent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockFromInlineContent(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "source_content", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockFromPublicBlob(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockFromPublicBlob(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "source_uri", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockFromPublicFile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockFromPublicFile(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "source_uri", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockFromExistingBlob(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockFromExistingBlob(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "source_uri", "type"},
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
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockFromLocalBlob(data, sourceBlob.Name()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
					testCheckAzureRMStorageBlobMatchesFile(data.ResourceName, blobs.BlockBlob, sourceBlob.Name()),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "source", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blockFromLocalFileWithContentMd5(t *testing.T) {
	sourceBlob, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Failed to create local source blob file")
	}

	if err := testAccAzureRMStorageBlob_populateTempFile(sourceBlob); err != nil {
		t.Fatalf("Error populating temp file: %s", err)
	}
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_contentMd5ForLocalFile(data, sourceBlob.Name()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "example.vhd"),
					resource.TestCheckResourceAttr(data.ResourceName, "source", sourceBlob.Name()),
				),
			},
			data.ImportStep("parallelism", "size", "source", "type"),
		},
	})
}

func TestAccAzureRMStorageBlob_contentType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_contentType(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
			{
				Config: testAccAzureRMStorageBlob_contentTypeUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_contentTypePremium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_contentTypePremium(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blobContentPropertiesCacheControl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blobContentPropertiesUpdated(data, testChangeableSetPropertiesInput{
					CacheControl:       "max-age=5",
					ContentType:        "application/octet-stream",
					ContentEncoding:    "identity",
					ContentLanguage:    "en",
					ContentDisposition: "attachment",
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
			{
				Config: testAccAzureRMStorageBlob_blobContentPropertiesUpdated(data, testChangeableSetPropertiesInput{
					CacheControl:       "max-age=500",
					ContentType:        "application/octet-stream",
					ContentEncoding:    "identity",
					ContentLanguage:    "en",
					ContentDisposition: "attachment",
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blobContentPropertiesContentType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blobContentPropertiesUpdated(data, testChangeableSetPropertiesInput{
					CacheControl:       "max-age=5",
					ContentType:        "application/octet-stream",
					ContentEncoding:    "identity",
					ContentLanguage:    "en",
					ContentDisposition: "attachment",
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
			{
				Config: testAccAzureRMStorageBlob_blobContentPropertiesUpdated(data, testChangeableSetPropertiesInput{
					CacheControl:       "max-age=5",
					ContentType:        "text/plain",
					ContentEncoding:    "deflate",
					ContentLanguage:    "en",
					ContentDisposition: "attachment",
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blobContentPropertiesContentEncoding(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blobContentPropertiesUpdated(data, testChangeableSetPropertiesInput{
					CacheControl:       "max-age=5",
					ContentType:        "application/octet-stream",
					ContentEncoding:    "identity",
					ContentLanguage:    "en",
					ContentDisposition: "attachment",
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
			{
				Config: testAccAzureRMStorageBlob_blobContentPropertiesUpdated(data, testChangeableSetPropertiesInput{
					CacheControl:       "max-age=5",
					ContentType:        "application/octet-stream",
					ContentEncoding:    "deflate",
					ContentLanguage:    "en",
					ContentDisposition: "attachment",
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blobContentPropertiesContentLanguage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blobContentPropertiesUpdated(data, testChangeableSetPropertiesInput{
					CacheControl:       "max-age=5",
					ContentType:        "application/octet-stream",
					ContentEncoding:    "identity",
					ContentLanguage:    "en",
					ContentDisposition: "attachment",
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
			{
				Config: testAccAzureRMStorageBlob_blobContentPropertiesUpdated(data, testChangeableSetPropertiesInput{
					CacheControl:       "max-age=5",
					ContentType:        "application/octet-stream",
					ContentEncoding:    "identity",
					ContentLanguage:    "fr",
					ContentDisposition: "attachment",
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_blobContentPropertiesContentDisposition(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blobContentPropertiesUpdated(data, testChangeableSetPropertiesInput{
					CacheControl:       "max-age=5",
					ContentType:        "application/octet-stream",
					ContentEncoding:    "identity",
					ContentLanguage:    "en",
					ContentDisposition: "attachment",
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
			{
				Config: testAccAzureRMStorageBlob_blobContentPropertiesUpdated(data, testChangeableSetPropertiesInput{
					CacheControl:       "max-age=5",
					ContentType:        "application/octet-stream",
					ContentEncoding:    "identity",
					ContentLanguage:    "en",
					ContentDisposition: "inline",
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_pageEmpty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_pageEmpty(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_pageEmptyPremium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_pageEmptyPremium(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_pageEmptyMetaData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_pageEmptyMetaData(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_pageFromExistingBlob(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_pageFromExistingBlob(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "source_uri", "type"},
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
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_pageFromLocalBlob(data, sourceBlob.Name()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
					testCheckAzureRMStorageBlobMatchesFile(data.ResourceName, blobs.PageBlob, sourceBlob.Name()),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "source", "type"},
			},
		},
	})
}

func TestAccAzureRMStorageBlob_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_blockFromPublicBlob(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMStorageBlob_requiresImport),
		},
	})
}

func TestAccAzureRMStorageBlob_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageBlob_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
			{
				Config: testAccAzureRMStorageBlob_updateUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageBlobExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parallelism", "size", "type"},
			},
		},
	})
}

func testCheckAzureRMStorageBlobExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

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
		storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

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
		storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

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
	storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_blob" {
			continue
		}

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

func testAccAzureRMStorageBlob_appendEmpty(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_template(data, "private")
	return fmt.Sprintf(`
provider "azurerm" {}

%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Append"
}
`, template)
}

func testAccAzureRMStorageBlob_appendEmptyMetaData(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_template(data, "private")
	return fmt.Sprintf(`
provider "azurerm" {}

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

func testAccAzureRMStorageBlob_blockEmpty(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_template(data, "private")
	return fmt.Sprintf(`
provider "azurerm" {}

%s

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
}
`, template)
}

func testAccAzureRMStorageBlob_blockEmptyAzureADAuth(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_template(data, "private")
	return fmt.Sprintf(`
provider "azurerm" {
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

func testAccAzureRMStorageBlob_blockEmptyMetaData(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_template(data, "private")
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

func testAccAzureRMStorageBlob_blockEmptyAccessTier(data acceptance.TestData, accessTier blobs.AccessTier) string {
	template := testAccAzureRMStorageBlob_templateBlockBlobStorage(data, "private")
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

func testAccAzureRMStorageBlob_blockFromInlineContent(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_template(data, "blob")
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
  source_content         = "Wubba Lubba Dub Dub"
}
`, template)
}

func testAccAzureRMStorageBlob_blockFromPublicBlob(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_template(data, "blob")
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

func testAccAzureRMStorageBlob_blockFromPublicFile(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_template(data, "private")
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

func testAccAzureRMStorageBlob_blockFromExistingBlob(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_template(data, "private")
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

func testAccAzureRMStorageBlob_blockFromLocalBlob(data acceptance.TestData, fileName string) string {
	template := testAccAzureRMStorageBlob_template(data, "private")
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

func testAccAzureRMStorageBlob_contentMd5ForLocalFile(data acceptance.TestData, fileName string) string {
	template := testAccAzureRMStorageBlob_template(data, "blob")
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

func testAccAzureRMStorageBlob_contentType(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_template(data, "private")
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

func testAccAzureRMStorageBlob_contentTypePremium(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_templatePremium(data, "private")
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

func testAccAzureRMStorageBlob_contentTypeUpdated(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_template(data, "private")
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

func testAccAzureRMStorageBlob_blobContentPropertiesUpdated(data acceptance.TestData, changeable testChangeableSetPropertiesInput) string {
	template := testAccAzureRMStorageBlob_template(data, "private")
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
  size                   = 1024
  cache_control          = "%s"
  content_type           = "%s"
  content_encoding       = "%s"
  content_language       = "%s"
  content_disposition    = "%s"
}
`, template, changeable.CacheControl, changeable.ContentType, changeable.ContentEncoding, changeable.ContentLanguage, changeable.ContentDisposition)
}

func testAccAzureRMStorageBlob_pageEmpty(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_template(data, "private")
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

func testAccAzureRMStorageBlob_pageEmptyPremium(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_templatePremium(data, "private")
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

func testAccAzureRMStorageBlob_pageEmptyMetaData(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_template(data, "private")
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

func testAccAzureRMStorageBlob_pageFromExistingBlob(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_template(data, "private")
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

func testAccAzureRMStorageBlob_pageFromLocalBlob(data acceptance.TestData, fileName string) string {
	template := testAccAzureRMStorageBlob_template(data, "private")
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

func testAccAzureRMStorageBlob_populateTempFile(input *os.File) error {
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

func testAccAzureRMStorageBlob_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_blockFromPublicBlob(data)
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

func testAccAzureRMStorageBlob_update(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_template(data, "private")
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

func testAccAzureRMStorageBlob_updateUpdated(data acceptance.TestData) string {
	template := testAccAzureRMStorageBlob_template(data, "private")
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

func testAccAzureRMStorageBlob_template(data acceptance.TestData, accessLevel string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  allow_blob_public_access = true
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, accessLevel)
}

func testAccAzureRMStorageBlob_templateBlockBlobStorage(data acceptance.TestData, accessLevel string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  allow_blob_public_access = true
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, accessLevel)
}

func testAccAzureRMStorageBlob_templatePremium(data acceptance.TestData, accessLevel string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Premium"
  account_replication_type = "LRS"
  allow_blob_public_access = true
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, accessLevel)
}
