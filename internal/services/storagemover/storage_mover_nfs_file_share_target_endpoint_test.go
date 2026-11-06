// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storagemover

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/endpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/storagemovers"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func TestStorageMoverNfsFileShareTargetEndpointFlatten(t *testing.T) {
	id := endpoints.NewEndpointID("00000000-0000-0000-0000-000000000000", "example-rg", "example-mover", "example-endpoint")
	storageAccountID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Storage/storageAccounts/example"
	r := StorageMoverNfsFileShareTargetEndpointResource{}
	metadata := sdk.NewResourceMetaData(nil, r)
	metadata.SetID(id)

	for _, description := range []*string{pointer.To("description"), nil} {
		model := endpoints.Endpoint{
			Properties: endpoints.AzureStorageNfsFileShareEndpointProperties{
				FileShareName:            "example-share",
				StorageAccountResourceId: storageAccountID,
				Description:              description,
			},
		}
		if err := r.flatten(metadata, &id, &model); err != nil {
			t.Fatalf("flattening endpoint: %v", err)
		}

		var state StorageMoverNfsFileShareTargetEndpointModel
		if err := metadata.Decode(&state); err != nil {
			t.Fatalf("decoding endpoint state: %v", err)
		}
		expected := StorageMoverNfsFileShareTargetEndpointModel{
			Name:             id.EndpointName,
			StorageMoverId:   storagemovers.NewStorageMoverID(id.SubscriptionId, id.ResourceGroupName, id.StorageMoverName).ID(),
			FileShareName:    "example-share",
			StorageAccountId: storageAccountID,
			Description:      pointer.From(description),
		}
		if !reflect.DeepEqual(state, expected) {
			t.Fatalf("expected state %#v, got %#v", expected, state)
		}
	}

	identity, err := metadata.ResourceData.Identity()
	if err != nil {
		t.Fatalf("reading resource identity: %v", err)
	}
	for name, expected := range map[string]string{
		"name":                id.EndpointName,
		"resource_group_name": id.ResourceGroupName,
		"storage_mover_name":  id.StorageMoverName,
		"subscription_id":     id.SubscriptionId,
	} {
		if actual := identity.Get(name); actual != expected {
			t.Errorf("expected identity %s = %q, got %q", name, expected, actual)
		}
	}
}

func TestStorageMoverNfsFileShareTargetEndpointListRegistration(t *testing.T) {
	ctx := context.Background()
	resourceType := StorageMoverNfsFileShareTargetEndpointResource{}.ResourceType()
	for _, registered := range (Registration{}).ListResources() {
		var metadata resource.MetadataResponse
		registered.Metadata(ctx, resource.MetadataRequest{}, &metadata)
		if metadata.TypeName != resourceType {
			continue
		}

		configurable, ok := registered.(sdk.FrameworkListWrappedResourceWithConfig)
		if !ok {
			t.Fatal("endpoint list resource must accept a storage_mover_id")
		}
		var schema list.ListResourceSchemaResponse
		configurable.ListResourceConfigSchema(ctx, list.ListResourceSchemaRequest{}, &schema)
		if attribute, ok := schema.Schema.Attributes["storage_mover_id"]; !ok || !attribute.IsRequired() {
			t.Fatal("storage_mover_id must be a required list argument")
		}
		if registered.ResourceFunc().Identity == nil {
			t.Fatal("endpoint list resource must expose resource identity")
		}
		return
	}
	t.Fatalf("list resource %q is not registered", resourceType)
}
