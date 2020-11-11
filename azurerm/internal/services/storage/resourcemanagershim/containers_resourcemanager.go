package resourcemanagershim

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/containers"
)

type ResourceManagerStorageContainerWrapper struct {
	client *storage.BlobContainersClient
}

func NewResourceManagerStorageContainerWrapper(client *storage.BlobContainersClient) ResourceManagerStorageContainerWrapper {
	return ResourceManagerStorageContainerWrapper{
		client: client,
	}
}

func (w ResourceManagerStorageContainerWrapper) Create(ctx context.Context, resourceGroup, accountName, containerName string, input containers.CreateInput) error {
	rmInput := storage.BlobContainer{
		ContainerProperties: &storage.ContainerProperties{
			PublicAccess: w.mapAccessLevel(input.AccessLevel),
			Metadata:     w.mapDataPlaneMetaData(input.MetaData),
		},
	}
	_, err := w.client.Create(ctx, resourceGroup, accountName, containerName, rmInput)
	return err
}

func (w ResourceManagerStorageContainerWrapper) Delete(ctx context.Context, resourceGroup, accountName, containerName string) error {
	resp, err := w.client.Delete(ctx, resourceGroup, accountName, containerName)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return err
	}
	return nil
}

func (w ResourceManagerStorageContainerWrapper) Exists(ctx context.Context, resourceGroup, accountName, containerName string) (*bool, error) {
	container, err := w.client.Get(ctx, resourceGroup, accountName, containerName)
	if err != nil {
		if utils.ResponseWasNotFound(container.Response) {
			return utils.Bool(false), nil
		}

		return nil, err
	}

	return utils.Bool(container.ContainerProperties != nil), nil
}

func (w ResourceManagerStorageContainerWrapper) Get(ctx context.Context, resourceGroup, accountName, containerName string) (*StorageContainerProperties, error) {
	container, err := w.client.Get(ctx, resourceGroup, accountName, containerName)
	if err != nil {
		if utils.ResponseWasNotFound(container.Response) {
			return nil, nil
		}

		return nil, err
	}

	if container.ContainerProperties == nil {
		return nil, fmt.Errorf("`properties` is null in the API response")
	}

	var nullableBoolToBool = func(input *bool) bool {
		if input == nil {
			return false
		}

		return *input
	}

	output := StorageContainerProperties{
		AccessLevel:           w.mapPublicAccess(container.ContainerProperties.PublicAccess),
		HasLegalHold:          nullableBoolToBool(container.ContainerProperties.HasLegalHold),
		HasImmutabilityPolicy: nullableBoolToBool(container.ContainerProperties.HasImmutabilityPolicy),
		MetaData:              w.mapResourceManagerMetaData(container.ContainerProperties.Metadata),
	}
	return &output, nil
}

func (w ResourceManagerStorageContainerWrapper) UpdateAccessLevel(ctx context.Context, resourceGroup, accountName, containerName string, level containers.AccessLevel) error {
	rmInput := storage.BlobContainer{
		ContainerProperties: &storage.ContainerProperties{
			PublicAccess: w.mapAccessLevel(level),
		},
	}
	_, err := w.client.Update(ctx, resourceGroup, accountName, containerName, rmInput)
	return err
}

func (w ResourceManagerStorageContainerWrapper) UpdateMetaData(ctx context.Context, resourceGroup, accountName, containerName string, metaData map[string]string) error {
	rmInput := storage.BlobContainer{
		ContainerProperties: &storage.ContainerProperties{
			Metadata: w.mapDataPlaneMetaData(metaData),
		},
	}
	_, err := w.client.Update(ctx, resourceGroup, accountName, containerName, rmInput)
	return err
}

func (w ResourceManagerStorageContainerWrapper) mapAccessLevel(input containers.AccessLevel) storage.PublicAccess {
	switch input {
	case containers.Blob:
		return storage.PublicAccessBlob

	case containers.Container:
		return storage.PublicAccessContainer
	}

	// this is an empty string, or _could_ be a value going forwards - it's easier to default this
	return storage.PublicAccessNone
}

func (w ResourceManagerStorageContainerWrapper) mapDataPlaneMetaData(input map[string]string) map[string]*string {
	output := make(map[string]*string, len(input))

	for k, v := range input {
		output[k] = &v
	}

	return output
}

func (w ResourceManagerStorageContainerWrapper) mapPublicAccess(input storage.PublicAccess) containers.AccessLevel {
	switch input {
	case storage.PublicAccessBlob:
		return containers.Blob

	case storage.PublicAccessContainer:
		return containers.Container
	}

	return containers.Private
}

func (w ResourceManagerStorageContainerWrapper) mapResourceManagerMetaData(input map[string]*string) map[string]string {
	output := make(map[string]string, len(input))

	for k, v := range input {
		if v == nil {
			continue
		}

		output[k] = *v
	}

	return output
}
