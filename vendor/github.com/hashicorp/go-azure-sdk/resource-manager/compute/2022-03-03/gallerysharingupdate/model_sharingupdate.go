package gallerysharingupdate

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SharingUpdate struct {
	Groups        *[]SharingProfileGroup      `json:"groups,omitempty"`
	OperationType SharingUpdateOperationTypes `json:"operationType"`
}
