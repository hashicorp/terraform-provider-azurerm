package azuresdkhacks

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2020-02-02/insights"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/parse"
)

type WebTestsClient struct {
	client *insights.WebTestsClient
}

func NewWebTestsClient(client insights.WebTestsClient) WebTestsClient {
	return WebTestsClient{
		client: &client,
	}
}

// CreateOrUpdate is a workaround to handle that the Azure API can return either a 200 / 201
// rather than the 200 documented in the Swagger - this is a workaround until the upstream PR is merged
// TF issue: https://github.com/hashicorp/terraform-provider-azurerm/issues/16805
// Swagger PR: https://github.com/Azure/azure-rest-api-specs/pull/19104
func (c WebTestsClient) CreateOrUpdate(ctx context.Context, id parse.WebTestId, webTestDefinition insights.WebTest) (result insights.WebTest, err error) {
	req, err := c.client.CreateOrUpdatePreparer(ctx, id.ResourceGroup, id.Name, webTestDefinition)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.WebTestsClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	resp, err := c.client.CreateOrUpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "insights.WebTestsClient", "CreateOrUpdate", resp, "Failure sending request")
		return
	}

	result, err = c.createOrUpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.WebTestsClient", "CreateOrUpdate", resp, "Failure responding to request")
		return
	}

	return
}

// createOrUpdateResponder is a workaround to handle that the Azure API can return either a 200 / 201
// rather than the 200 documented in the Swagger - this is a workaround until the upstream PR is merged
// TF issue: https://github.com/hashicorp/terraform-provider-azurerm/issues/16805
// Swagger PR: https://github.com/Azure/azure-rest-api-specs/pull/19104
func (c WebTestsClient) createOrUpdateResponder(resp *http.Response) (result insights.WebTest, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

func (c WebTestsClient) Delete(ctx context.Context, id parse.WebTestId) (result autorest.Response, err error) {
	return c.client.Delete(ctx, id.ResourceGroup, id.Name)
}

func (c WebTestsClient) Get(ctx context.Context, id parse.WebTestId) (result insights.WebTest, err error) {
	return c.client.Get(ctx, id.ResourceGroup, id.Name)
}
