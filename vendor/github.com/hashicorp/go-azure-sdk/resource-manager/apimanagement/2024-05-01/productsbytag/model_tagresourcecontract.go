package productsbytag

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagResourceContract struct {
	Api       *ApiTagResourceContractProperties       `json:"api,omitempty"`
	Operation *OperationTagResourceContractProperties `json:"operation,omitempty"`
	Product   *ProductTagResourceContractProperties   `json:"product,omitempty"`
	Tag       TagTagResourceContractProperties        `json:"tag"`
}
