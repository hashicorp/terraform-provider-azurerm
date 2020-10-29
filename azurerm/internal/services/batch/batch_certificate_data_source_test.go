package batch_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccBatchCertificateDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_batch_certificate", "test")

	//lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchCertificateDataSource_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", "sha1-42c107874fd0e4a9583292a2f1098e8fe4b2edda"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_name", fmt.Sprintf("testaccbatch%s", data.RandomString)),
					resource.TestCheckResourceAttr(data.ResourceName, "format", "Pfx"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_data", "MIIFqzCCA5OgAwIBAgIJAMs4jwMPq7T1MA0GCSqGSIb3DQEBCwUAMGwxCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApTb21lLVN0YXRlMRgwFgYDVQQKDA9UZXJyYWZvcm0gVGVzdHMxDjAMBgNVBAsMBUF6dXJlMR4wHAYDVQQDDBVUZXJyYWZvcm0gQXBwIEdhdGV3YXkwHhcNMTYxMTAxMTcxOTEyWhcNMjYxMDMwMTcxOTEyWjBsMQswCQYDVQQGEwJVUzETMBEGA1UECAwKU29tZS1TdGF0ZTEYMBYGA1UECgwPVGVycmFmb3JtIFRlc3RzMQ4wDAYDVQQLDAVBenVyZTEeMBwGA1UEAwwVVGVycmFmb3JtIEFwcCBHYXRld2F5MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA49HW2pYIlW/mlaadLA1AsXiV48xVhXAvGVk3DEl1ffjp5bN8rap5WV1D83uMg1Ii7CJM8yNHkRkvN8n5WXFng4R5V1jPxGOTAj+xLybvEASi++GZelWdpOuMk8/nAoKPMbQ5NyKFy5WzlOduMldR7Awt2pwdId3akqm1i9ITG9Js+4P4nYXM8vfJCajILqi4YfhEoCNvS1EUgvlpSFE7pfNhc2W+zsfUWxWmB2SpWwX9MgQ1D4OmdKp+Eo+b6vzst3XArKMHMadPTUAk8H+ZgAnlX9yO+3vQ6z86vma/WgrG2LH6GCGXBjmKlhxVCPMLA5LeRUwEGc/Q7X/ClitGWY9umPN1XVj5e5Di1K2M082Y14mgbTTRTpv/nx7Xlph+MHnVhEWvaGMpqCHuM1W1y7wIS1IREYQ2q+K54xxZSPKYJMSnmj6A0hR/LBV0rL1uVhedEpdviduuO76qCyZrGG4HwBlW4hnIaahLzgqlvlmbDUQonAVPDgi3brVdXJgLv2zi7/ZHFW3IHgDylUVIdig0ccbzxKymlkGQ0RsLBjWOyxak2J8bN5JNVyxSwX43NZqxJ8yOv5xjB+rVMri9SX3Dl5NbFzOjynov601Pmwvb7zYnyttG2Hl5EKrkahjijGRjGy3EWEiBiArLkdTKCDHBlHxykTEvY6ZH5B9waP0CAwEAAaNQME4wHQYDVR0OBBYEFD2/Hq3IivZ5RMOKrPsM7ijIFHmMMB8GA1UdIwQYMBaAFD2/Hq3IivZ5RMOKrPsM7ijIFHmMMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQELBQADggIBAKxHWO/Q4labjnCVxYi+kaMRCPJUdHj7lga8yi8EGHaL+CbwynkaiyTfPvtmcqiuaZM9BaXsuNMRcHMtXM0EHBsjViwAHk6SrqLXd/opFvMI2QbG93koFUCpczrpyO9GvnRN4iOIYbSPXAdGOB6bkpMbm/XajORoDrua+/ET/X/1FP0GZBTmEFwojuCfOI/VuJXj0OW8XzkLmsXiLpOiakjU1obBup/1lz9DtOEBsiB9Ury+f5gZ+FnZuqhgQxeDxlZ69P6YYAfkzhcfbf7HO+nMKhppAj1BFeR4SBb+F/fLchCGO5yohwkxWz3i2q9gTDhBgo31416viyCKFWSVW3Vn7jbsjZ+Q9MK1jVSOSxC7qoQkRoNy9SKpqylunXZb+K6F3HfBkDQvn3OwsxYiSOcX9JaWpQAInNIZVg+WrJ1PXm8PFIaVPJfMgP3GOdm9vRAMjOM5Bc9iqGr2spimFd5h0GmgLvh35B3jHHWF4i3NupJQ6hUvHQZtYZOxfwxnY0/LVBTyLTVlniFA7dGSI+5Uexm+Pjh7IMGI532jTONlfNm9Bz/jdf1o0FlOclzG6Eif22gml3GM3xCUVlaElylYNAjO2lfvZuRVo5GKdMwtV9acNl0OwSx+0zbMYY2Ni3jQCI4kOL5Csctryf0rHXTlCCvnzBYVDPKmFJPna61T"),
					resource.TestCheckResourceAttr(data.ResourceName, "thumbprint", "42c107874fd0e4a9583292a2f1098e8fe4b2edda"),
					resource.TestCheckResourceAttr(data.ResourceName, "thumbprint_algorithm", "sha1"), // api now always returns this as lowercase
				),
			},
		},
	})
}

func testAccBatchCertificateDataSource_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
}

resource "azurerm_batch_certificate" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  account_name         = azurerm_batch_account.test.name
  certificate          = filebase64("testdata/batch_certificate.pfx")
  format               = "Pfx"
  password             = "terraform"
  thumbprint           = "42c107874fd0e4a9583292a2f1098e8fe4b2edda"
  thumbprint_algorithm = "SHA1"
}

data "azurerm_batch_certificate" "test" {
  name                = azurerm_batch_certificate.test.name
  account_name        = azurerm_batch_account.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
