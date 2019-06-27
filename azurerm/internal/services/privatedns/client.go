package privatedns

import (
	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
)

type Client struct {
	PrivateZonesClient privatedns.PrivateZonesClient
}
