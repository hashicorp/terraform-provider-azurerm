package contact

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SpacecraftsListAvailableContactsOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AvailableContacts
}

type SpacecraftsListAvailableContactsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AvailableContacts
}

type SpacecraftsListAvailableContactsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SpacecraftsListAvailableContactsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SpacecraftsListAvailableContacts ...
func (c ContactClient) SpacecraftsListAvailableContacts(ctx context.Context, id SpacecraftId, input ContactParameters) (result SpacecraftsListAvailableContactsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &SpacecraftsListAvailableContactsCustomPager{},
		Path:       fmt.Sprintf("%s/listAvailableContacts", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// SpacecraftsListAvailableContactsThenPoll performs SpacecraftsListAvailableContacts then polls until it's completed
func (c ContactClient) SpacecraftsListAvailableContactsThenPoll(ctx context.Context, id SpacecraftId, input ContactParameters) error {
	result, err := c.SpacecraftsListAvailableContacts(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SpacecraftsListAvailableContacts: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after SpacecraftsListAvailableContacts: %+v", err)
	}

	return nil
}
