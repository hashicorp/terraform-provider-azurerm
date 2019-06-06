package cognitive

import "github.com/Azure/azure-sdk-for-go/services/cognitiveservices/mgmt/2017-04-18/cognitiveservices"

type Client struct {
	AccountsClient cognitiveservices.AccountsClient
}
