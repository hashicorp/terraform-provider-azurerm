package expressroutecircuitroutestablesummary

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCircuitRoutesTableSummary struct {
	As          *int64  `json:"as,omitempty"`
	Neighbor    *string `json:"neighbor,omitempty"`
	StatePfxRcd *string `json:"statePfxRcd,omitempty"`
	UpDown      *string `json:"upDown,omitempty"`
	V           *int64  `json:"v,omitempty"`
}
