package shim

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/containers"
)

type DataPlaneStorageContainerWrapper struct {
	client *containers.Client
}

func NewDataPlaneStorageContainerWrapper(client *containers.Client) StorageContainerWrapper {
	return DataPlaneStorageContainerWrapper{
		client: client,
	}
}

func (w DataPlaneStorageContainerWrapper) Create(ctx context.Context, _, accountName, containerName string, input containers.CreateInput) error {
	timeout, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context is missing a timeout")
	}

	if resp, err := w.client.Create(ctx, accountName, containerName, input); err != nil {
		// If we fail due to previous delete still in progress, then we can retry
		if utils.ResponseWasConflict(resp.Response) && strings.Contains(err.Error(), "ContainerBeingDeleted") {
			stateConf := &pluginsdk.StateChangeConf{
				Pending:        []string{"waitingOnDelete"},
				Target:         []string{"succeeded"},
				Refresh:        w.createRefreshFunc(ctx, accountName, containerName, input),
				PollInterval:   10 * time.Second,
				NotFoundChecks: 180,
				Timeout:        time.Until(timeout),
			}

			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("failed creating container: %+v", err)
			}
		} else {
			return fmt.Errorf("failed creating container: %+v", err)
		}
	}
	return nil
}

func (w DataPlaneStorageContainerWrapper) Delete(ctx context.Context, _, accountName, containerName string) error {
	resp, err := w.client.Delete(ctx, accountName, containerName)
	if utils.ResponseWasNotFound(resp) {
		return nil
	}

	return err
}

func (w DataPlaneStorageContainerWrapper) Exists(ctx context.Context, _, accountName, containerName string) (*bool, error) {
	existing, err := w.client.GetProperties(ctx, accountName, containerName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return nil, err
		}
	}

	exists := !utils.ResponseWasNotFound(existing.Response)
	return &exists, nil
}

func (w DataPlaneStorageContainerWrapper) Get(ctx context.Context, _, accountName, containerName string) (*StorageContainerProperties, error) {
	props, err := w.client.GetProperties(ctx, accountName, containerName)
	if err != nil {
		if utils.ResponseWasNotFound(props.Response) {
			return nil, nil
		}

		return nil, err
	}

	return &StorageContainerProperties{
		AccessLevel:           props.AccessLevel,
		MetaData:              props.MetaData,
		HasImmutabilityPolicy: props.HasImmutabilityPolicy,
		HasLegalHold:          props.HasLegalHold,
	}, nil
}

func (w DataPlaneStorageContainerWrapper) UpdateAccessLevel(ctx context.Context, _, accountName, containerName string, level containers.AccessLevel) error {
	_, err := w.client.SetAccessControl(ctx, accountName, containerName, level)
	return err
}

func (w DataPlaneStorageContainerWrapper) UpdateMetaData(ctx context.Context, _, accountName, containerName string, metaData map[string]string) error {
	_, err := w.client.SetMetaData(ctx, accountName, containerName, metaData)
	return err
}

func (w DataPlaneStorageContainerWrapper) createRefreshFunc(ctx context.Context, accountName string, containerName string, input containers.CreateInput) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := w.client.Create(ctx, accountName, containerName, input)
		if err != nil {
			if !utils.ResponseWasConflict(resp.Response) {
				return nil, "", err
			}

			if utils.ResponseWasConflict(resp.Response) && strings.Contains(err.Error(), "ContainerBeingDeleted") {
				return nil, "waitingOnDelete", nil
			}
		}

		return "succeeded", "succeeded", nil
	}
}
