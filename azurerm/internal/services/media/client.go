package media

import "github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2018-07-01/media"

type Client struct {
	ServicesClient media.MediaservicesClient
}
