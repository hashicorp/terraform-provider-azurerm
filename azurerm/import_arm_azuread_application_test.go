package azurerm

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMActiveDirectoryApplication_importSimple(t *testing.T) {
	resourceName := "azurerm_azuread_application.test"

	id := uuid.New().String()
	config := testAccAzureRMActiveDirectoryApplication_simple(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMActiveDirectoryApplication_importAdvanced(t *testing.T) {
	resourceName := "azurerm_azuread_application.test"

	id := uuid.New().String()
	config := testAccAzureRMActiveDirectoryApplication_advanced(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMActiveDirectoryApplication_importKeyCredential(t *testing.T) {
	resourceName := "azurerm_azuread_application.test"

	id := uuid.New().String()
	keyId := uuid.New().String()
	config := testAccAzureRMActiveDirectoryApplication_keyCredential_single(id, keyId, "AsymmetricX509Cert")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProvidersWithTLS,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMActiveDirectoryApplication_importPasswordCredential(t *testing.T) {
	resourceName := "azurerm_azuread_application.test"

	id := uuid.New().String()
	keyId := uuid.New().String()
	timeStart := time.Now().UTC()
	timeEnd := timeStart.Add(time.Duration(1) * time.Hour)
	config := testAccAzureRMActiveDirectoryApplication_passwordCredential_single(id, keyId, timeStart, timeEnd)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
