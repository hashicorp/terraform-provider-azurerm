package eventhub_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type EventHubSharedAccessSignatureDataSource struct {
}

func TestAccEventHubSharedAccessSignatureDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_sas", "test")
	r := EventHubSharedAccessSignatureDataSource{}
	utcNow := time.Now().UTC()
	endDate := utcNow.Add(time.Hour * 24).Format(time.RFC3339)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data, endDate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sas").Exists(),
			),
		},
	})
}

func (EventHubSharedAccessSignatureDataSource) basic(data acceptance.TestData, enddate string) string {
	return fmt.Sprintf(`

data "azurerm_eventhub_sas" "test" {
  connection_string = "Endpoint=sb://acctesteventhubnamespace-test01.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=IUSvXLiPZ3uAQcso/cL7vTiL4zsc/EMtcUzNCC2dhaM="
  expiry            = "%s"
}
`, enddate)
}
