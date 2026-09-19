package documents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutocompleteItem struct {
	QueryPlusText string `json:"queryPlusText"`
	Text          string `json:"text"`
}
