// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/clusters"
	sdkclient "github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
)

func createOrUpdateKustoClusterThenPoll(ctx context.Context, client *clusters.ClustersClient, id commonids.KustoClusterId, input clusters.Cluster, federatedIdentityClientID *string, callback func() error) error {
	return updateKustoClusterRequestThenPoll(ctx, client, id, input, http.MethodPut, []int{http.StatusCreated, http.StatusOK}, clusters.CreateOrUpdateOperationOptions{}, federatedIdentityClientID, callback)
}

func updateKustoClusterThenPoll(ctx context.Context, client *clusters.ClustersClient, id commonids.KustoClusterId, input clusters.ClusterUpdate, federatedIdentityClientID *string, callback func() error) error {
	return updateKustoClusterRequestThenPoll(ctx, client, id, input, http.MethodPatch, []int{http.StatusAccepted, http.StatusCreated, http.StatusOK}, clusters.UpdateOperationOptions{}, federatedIdentityClientID, callback)
}

func updateKustoClusterRequestThenPoll(ctx context.Context, client *clusters.ClustersClient, id commonids.KustoClusterId, input any, method string, expectedStatusCodes []int, options sdkclient.Options, federatedIdentityClientID *string, callback func() error) error {
	body, err := marshalKustoClusterRequest(input, federatedIdentityClientID)
	if err != nil {
		return fmt.Errorf("marshalling request body: %+v", err)
	}

	req, err := client.Client.NewRequest(ctx, sdkclient.RequestOptions{
		ContentType:         "application/json; charset=utf-8",
		ExpectedStatusCodes: expectedStatusCodes,
		HttpMethod:          method,
		OptionsObject:       options,
		Path:                id.ID(),
	})
	if err != nil {
		return err
	}

	if err := req.Marshal(body); err != nil {
		return err
	}

	resp, err := req.Execute(ctx)
	if err != nil {
		return fmt.Errorf("performing %s: %+v", method, err)
	}

	poller, err := resourcemanager.PollerFromResponse(resp, client.Client)
	if err != nil {
		return err
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after %s: %+v", method, err)
	}

	return nil
}

// marshalKustoClusterRequest adds `federatedIdentityClientId` to the request body. The pinned SDK model for the
// generated client doesn't expose the field, so it's injected here until a 2025-02-14 SDK package is available.
func marshalKustoClusterRequest(input any, federatedIdentityClientID *string) (json.RawMessage, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	if federatedIdentityClientID == nil {
		return body, nil
	}

	var request map[string]any
	if err := json.Unmarshal(body, &request); err != nil {
		return nil, err
	}
	properties, ok := request["properties"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("request properties are missing")
	}
	keyVaultProperties, ok := properties["keyVaultProperties"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("request keyVaultProperties are missing")
	}
	keyVaultProperties["federatedIdentityClientId"] = *federatedIdentityClientID

	return json.Marshal(request)
}
