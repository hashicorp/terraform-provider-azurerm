package encodings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Deinterlace struct {
	Mode   *DeinterlaceMode   `json:"mode,omitempty"`
	Parity *DeinterlaceParity `json:"parity,omitempty"`
}
