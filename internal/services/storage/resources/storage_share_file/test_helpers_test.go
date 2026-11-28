// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_share_file_test

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type StorageShareDirectoryResource struct{}

func (r StorageShareDirectoryResource) nestedWithBackslashes(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "c" {
  name              = "c"
  storage_share_url = azurerm_storage_share.test.url
}

resource "azurerm_storage_share_directory" "dos" {
  name              = "c\\dos"
  storage_share_url = azurerm_storage_share.test.url
  depends_on        = [azurerm_storage_share_directory.c]
}
`, template)
}

func (r StorageShareDirectoryResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "test" {
  name               = "fileshare"
  storage_account_id = azurerm_storage_account.test.id
  quota              = 50
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func populateTempFile(input *os.File) error {
	if err := input.Truncate(25*1024*1024 + 512); err != nil {
		return errors.New("Failed to truncate file to 25M")
	}

	for i := int64(0); i < 20; i += 2 {
		randomBytes := make([]byte, 1*1024*1024)
		if _, err := rand.Read(randomBytes); err != nil {
			return errors.New("Failed to read random bytes")
		}

		if _, err := input.WriteAt(randomBytes, i*1024*1024); err != nil {
			return errors.New("Failed to write random bytes to file")
		}
	}

	randomBytes := make([]byte, 5*1024*1024)
	if _, err := rand.Read(randomBytes); err != nil {
		return errors.New("Failed to read random bytes")
	}

	if _, err := input.WriteAt(randomBytes, 20*1024*1024); err != nil {
		return errors.New("Failed to write random bytes to file")
	}

	if err := input.Close(); err != nil {
		return errors.New("Failed to close source blob")
	}

	return nil
}
