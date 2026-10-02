// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

func TestHyperVHostTemplateUsesLocalAdminPasswordForAutoLogon(t *testing.T) {
	data := acceptance.TestData{
		RandomInteger: 123456789012345678,
		RandomString:  "abcde",
		Locations: acceptance.Regions{
			Primary: "eastus",
		},
	}

	config := HyperVHostTestResource{}.hyperVTemplate(data)

	autoLogonContent := `content = "<AutoLogon><Password><Value>${local.admin_password}</Value></Password><Enabled>true</Enabled><LogonCount>1</LogonCount><Username>${local.admin_name}</Username></AutoLogon>"`

	if !strings.Contains(config, `admin_password      = local.admin_password`) {
		t.Fatalf("expected generated config to reuse local.admin_password for the VM admin password:\n%s", config)
	}

	if !strings.Contains(config, autoLogonContent) {
		t.Fatalf("expected generated config to keep AutoLogon content on a single line using local.admin_password:\n%s", config)
	}
}
