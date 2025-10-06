package iotoperations

import (
    "fmt"
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataflowGraph_basic(t *testing.T) {
    resourceName := "azurerm_dataflow_graph.test"

    resource.ParallelTest(t, resource.TestCase{
        PreCheck:     func() { /* add pre-checks if needed */ },
        Providers:    testAccProviders,
        CheckDestroy: testAccCheckDataflowGraphDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccDataflowGraphConfig_basic(),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr(resourceName, "name", "test-dataflow-graph"),
                    resource.TestCheckResourceAttr(resourceName, "resource_group_name", "test-rg"),
                    resource.TestCheckResourceAttr(resourceName, "instance_name", "test-instance"),
                    resource.TestCheckResourceAttr(resourceName, "dataflow_profile_name", "test-profile")
                ),
            },
        },
    })
}

func testAccDataflowGraphConfig_basic() string {
    return fmt.Sprintf(`
resource "azurerm_dataflow_graph" "test" {
  name                = "test-dataflow-graph"
  resource_group_name = "test-rg"
  instance_name       = "test-instance"
  dataflow_profile_name = "test-profile"

  properties {
    mode = "Enabled"
    request_disk_persistence = "Enabled"
    nodes {
      type = "source"
      name = "temperature"
      # Add other node fields as needed
    }
    # Add connections and other properties as needed
  }

  extended_location {
    name = "qmbrfwcpwwhggszhrdjv"
    type = "CustomLocation"
  }
}
`)
}

func testAccCheckDataflowGraphDestroy(s *terraform.State) error {
    //verification to resource destruction
    return nil