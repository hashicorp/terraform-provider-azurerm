package storagetasks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageTaskOperation struct {
	Name       StorageTaskOperationName `json:"name"`
	OnFailure  *OnFailure               `json:"onFailure,omitempty"`
	OnSuccess  *OnSuccess               `json:"onSuccess,omitempty"`
	Parameters *map[string]string       `json:"parameters,omitempty"`
}
