// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

func TestGeoReplicationPoller_Success(t *testing.T) {
	id := databases.NewDatabaseID("00000000-0000-0000-0000-000000000000", "my-rg", "amr1", "default")

	mockClient := &mockDatabasesClient{
		getResponse: &databases.GetOperationResponse{
			HttpResponse: &http.Response{StatusCode: 200},
			Model: &databases.Database{
				Properties: &databases.DatabaseCreateProperties{
					GeoReplication: &databases.DatabasePropertiesGeoReplication{
						LinkedDatabases: &[]databases.LinkedDatabase{
							{
								Id:    pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Cache/redisEnterprise/amr1/databases/default"),
								State: pointer.To(databases.LinkStateLinked),
							},
							{
								Id:    pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Cache/redisEnterprise/amr2/databases/default"),
								State: pointer.To(databases.LinkStateLinked),
							},
							{
								Id:    pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Cache/redisEnterprise/amr3/databases/default"),
								State: pointer.To(databases.LinkStateLinked),
							},
						},
					},
				},
			},
		},
	}

	pollerType := &geoReplicationPoller{
		client: mockClient,
		id:     id,
		toIds: []string{
			"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Cache/redisEnterprise/amr1/databases/default",
			"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Cache/redisEnterprise/amr2/databases/default",
			"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Cache/redisEnterprise/amr3/databases/default",
		},
	}

	result, err := pollerType.Poll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Status != pollers.PollingStatusSucceeded {
		t.Fatalf("expected status %s, got %s", pollers.PollingStatusSucceeded, result.Status)
	}
}

func TestGeoReplicationPoller_InProgress(t *testing.T) {
	id := databases.NewDatabaseID("00000000-0000-0000-0000-000000000000", "my-rg", "amr1", "default")

	mockClient := &mockDatabasesClient{
		getResponse: &databases.GetOperationResponse{
			HttpResponse: &http.Response{StatusCode: 200},
			Model: &databases.Database{
				Properties: &databases.DatabaseCreateProperties{
					GeoReplication: &databases.DatabasePropertiesGeoReplication{
						LinkedDatabases: &[]databases.LinkedDatabase{
							{
								Id:    pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Cache/redisEnterprise/amr1/databases/default"),
								State: pointer.To(databases.LinkStateLinked),
							},
							{
								Id:    pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Cache/redisEnterprise/amr2/databases/default"),
								State: pointer.To(databases.LinkStateLinked),
							},
						},
					},
				},
			},
		},
	}

	pollerType := &geoReplicationPoller{
		client: mockClient,
		id:     id,
		toIds: []string{
			"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Cache/redisEnterprise/amr1/databases/default",
			"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Cache/redisEnterprise/amr2/databases/default",
			"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Cache/redisEnterprise/amr3/databases/default",
		},
	}

	result, err := pollerType.Poll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Status != pollers.PollingStatusInProgress {
		t.Fatalf("expected status %s, got %s", pollers.PollingStatusInProgress, result.Status)
	}
}
