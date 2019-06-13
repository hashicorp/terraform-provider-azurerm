package notificationhub

import (
	"github.com/Azure/azure-sdk-for-go/services/notificationhubs/mgmt/2017-04-01/notificationhubs"
)

type Client struct {
	HubsClient       notificationhubs.Client
	NamespacesClient notificationhubs.NamespacesClient
}
