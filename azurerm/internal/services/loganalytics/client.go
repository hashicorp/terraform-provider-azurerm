package loganalytics

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2015-11-01-preview/operationalinsights"
	"github.com/Azure/azure-sdk-for-go/services/preview/operationsmanagement/mgmt/2015-11-01-preview/operationsmanagement"
)

type Client struct {
	LinkedServicesClient operationalinsights.LinkedServicesClient
	SolutionsClient      operationsmanagement.SolutionsClient
	WorkspacesClient     operationalinsights.WorkspacesClient
}
