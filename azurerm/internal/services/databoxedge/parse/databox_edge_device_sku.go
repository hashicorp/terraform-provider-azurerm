package parse

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/databoxedge/mgmt/2019-08-01/databoxedge"
)

// DataboxEdgeDeviceSku type
type DataboxEdgeDeviceSku struct {
	Name string
	Tier string
}

// DataboxEdgeDeviceSkuName parses the input string into a DataboxEdgeDeviceSku type
func DataboxEdgeDeviceSkuName(input string) (*DataboxEdgeDeviceSku, error) {
	if len(strings.TrimSpace(input)) == 0 {
		return nil, fmt.Errorf("unable to parse Databox Edge Device 'sku_name' %q", input)
	}

	skuParts := strings.Split(input, "-")

	if strings.TrimSpace(skuParts[0]) == "" {
		return nil, fmt.Errorf("invalid Databox Edge Device 'sku_name' %q", input)
	}

	// There is only one possible Tier so always set value to Standard
	databoxEdgeDeviceSku := DataboxEdgeDeviceSku{
		Name: skuParts[0],
		Tier: string(databoxedge.Standard),
	}

	return &databoxEdgeDeviceSku, nil
}
