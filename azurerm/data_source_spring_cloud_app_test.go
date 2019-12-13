package azurerm

//import (
//	"fmt"
//	"testing"
//
//	"github.com/hashicorp/terraform/helper/resource"
//)
//
//func TestAccDataSourceAzureRMSpringCloudApp_basic(t *testing.T) {
//	dataSourceName := "data.azurerm_spring_cloud_app.test"
//
//	resource.ParallelTest(t, resource.TestCase{
//		PreCheck:  func() { testAccPreCheck(t) },
//		Providers: testAccProviders,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccDataSourceSpringCloudApp_basic(),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(dataSourceName, "test", "test"),
//				),
//			},
//		},
//	})
//}
//
//func testAccDataSourceSpringCloudApp_basic() string {
//	config := testAccAzureRMSpringCloudApp_basic()
//	return fmt.Sprintf(`
//%s
//
//data "azurerm_spring_cloud_app" "test" {
//  resource_group = "${azurerm_spring_cloud_app.test.resource_group}"
//  service_name   = "${azurerm_spring_cloud_app.test.service_name}"
//  name           = "${azurerm_spring_cloud_app.test.name}"
//}
//`, config)
//}
