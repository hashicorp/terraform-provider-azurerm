package workspace

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

type WorkspaceContract struct {
	Id         *string                      `json:"id,omitempty"`
	Name       *string                      `json:"name,omitempty"`
	Properties *WorkspaceContractProperties `json:"properties,omitempty"`
	SystemData *systemdata.SystemData       `json:"systemData,omitempty"`
	Type       *string                      `json:"type,omitempty"`
}

type WorkspaceContractProperties struct {
	Description *string `json:"description,omitempty"`
	DisplayName string  `json:"displayName"`
}

type CreateOrUpdateOperationOptions struct {
	IfMatch *string
}

func DefaultCreateOrUpdateOperationOptions() CreateOrUpdateOperationOptions {
	return CreateOrUpdateOperationOptions{}
}

type DeleteOperationOptions struct {
	IfMatch *string
}

func DefaultDeleteOperationOptions() DeleteOperationOptions {
	return DeleteOperationOptions{}
}
