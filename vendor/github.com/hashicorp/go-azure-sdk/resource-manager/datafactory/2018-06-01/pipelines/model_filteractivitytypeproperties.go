package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FilterActivityTypeProperties struct {
	Condition Expression `json:"condition"`
	Items     Expression `json:"items"`
}
