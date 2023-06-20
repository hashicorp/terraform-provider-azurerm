package frontdoors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RulesEngineProperties struct {
	ResourceState *FrontDoorResourceState `json:"resourceState,omitempty"`
	Rules         *[]RulesEngineRule      `json:"rules,omitempty"`
}
