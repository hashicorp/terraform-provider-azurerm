package azurerm

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMAdApplication_importSimple(t *testing.T) {
	resourceName := "azurerm_ad_application.test"

	id := uuid.New().String()
	config := testAccAzureRMAdApplication_simple(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAdApplicationDestroy,
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

func TestAccAzureRMAdApplication_importAdvanced(t *testing.T) {
	resourceName := "azurerm_ad_application.test"

	id := uuid.New().String()
	config := testAccAzureRMAdApplication_advanced(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAdApplicationDestroy,
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

func TestAccAzureRMAdApplication_importKeyCredential(t *testing.T) {
	resourceName := "azurerm_ad_application.test"

	id := uuid.New().String()
	keyId := uuid.New().String()
	config := testAccAzureRMAdApplication_keyCredential(id, keyId)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProvidersWithTLS,
		CheckDestroy: testCheckAzureRMAdApplicationDestroy,
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

func TestAccAzureRMAdApplication_importPasswordCredential(t *testing.T) {
	resourceName := "azurerm_ad_application.test"

	id := uuid.New().String()
	keyId := uuid.New().String()
	timeStart := string(time.Now().UTC().Format(time.RFC3339))
	timeEnd := string(time.Now().UTC().Add(time.Duration(1) * time.Hour).Format(time.RFC3339))
	config := testAccAzureRMAdApplication_passwordCredential(id, keyId, timeStart, timeEnd)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAdApplicationDestroy,
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
