package mssql_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccMssqlVirtualMachine_listByServerID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_machine", "testlist1")
	r := MsSqlVirtualMachineResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_mssql_virtual_machine.list", 1), // expect at least the 1 we created
					querycheck.ExpectIdentity(
						"azurerm_mssql_virtual_machine.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func TestAccMssqlVirtualMachine_listByResourceGroupID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_VirtualMachine", "testlist1")
	r := MsSqlVirtualMachineResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_mssql_virtual_machine.list", 1), // expect at least the 1 we created
					querycheck.ExpectIdentity(
						"azurerm_mssql_virtual_machine.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r MsSqlVirtualMachineResource) basicQuery() string {
	return `
list "azurerm_mssql_server" "list" {
  provider = azurerm
  config {}
}
`
}

func (r MsSqlVirtualMachineResource) basicQueryByResourceGroupName() string {
	return fmt.Sprintf(`
list "azurerm_mssql_server" "list" {
  provider = azurerm
  config {
    resource_group_name = "{azurerm_resource_group.test.id}"
  }
}
`)
}
