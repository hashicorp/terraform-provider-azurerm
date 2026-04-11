package skillsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchIndexerIndexProjections struct {
	Parameters *SearchIndexerIndexProjectionsParameters `json:"parameters,omitempty"`
	Selectors  []SearchIndexerIndexProjectionSelector   `json:"selectors"`
}
