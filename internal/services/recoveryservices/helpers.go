package recoveryservices

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotectioncontainers"
)

// This code is a workaround for this bug https://github.com/Azure/azure-sdk-for-go/issues/2824
func handleAzureSdkForGoBug2824(id string) string {
	return strings.Replace(id, "/Subscriptions/", "/subscriptions/", 1)
}

func wasBadRequestWithNotExist(resp *http.Response, err error) bool {
	e, ok := err.(autorest.DetailedError)
	if !ok {
		return false
	}

	r, ok := e.Original.(*azure.RequestError)
	if !ok {
		return false
	}

	if r.ServiceError == nil || len(r.ServiceError.Details) == 0 {
		return false
	}

	sc, ok := r.ServiceError.Details[0]["code"]
	if !ok {
		return false
	}

	return response.WasBadRequest(resp) && sc == "SubscriptionIdNotRegisteredWithSrs"
}

func fetchHyperVContainerIdByFabricId(ctx context.Context, containerClient *replicationprotectioncontainers.ReplicationProtectionContainersClient, fabricId replicationfabrics.ReplicationFabricId) (string, error) {
	id, err := replicationprotectioncontainers.ParseReplicationFabricID(fabricId.ID())
	if err != nil {
		return "", fmt.Errorf("parsing %s: %+v", fabricId.ID(), err)
	}

	resp, err := containerClient.ListByReplicationFabricsComplete(ctx, *id)
	if err != nil {
		return "", fmt.Errorf("listing containers: %+v", err)
	}

	if len(resp.Items) == 0 || len(resp.Items) > 1 {
		return "", fmt.Errorf("expected one container but got %d", len(resp.Items))
	}

	if resp.Items[0].Id == nil {
		return "", fmt.Errorf("container id is nil")
	}

	return handleAzureSdkForGoBug2824(*resp.Items[0].Id), nil
}
