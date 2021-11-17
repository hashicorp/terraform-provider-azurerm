package generatedetailedcostreportoperationstatus

import "github.com/Azure/go-autorest/autorest"

type GenerateDetailedCostReportOperationStatusClient struct {
	Client  autorest.Client
	baseUri string
}

func NewGenerateDetailedCostReportOperationStatusClientWithBaseURI(endpoint string) GenerateDetailedCostReportOperationStatusClient {
	return GenerateDetailedCostReportOperationStatusClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
