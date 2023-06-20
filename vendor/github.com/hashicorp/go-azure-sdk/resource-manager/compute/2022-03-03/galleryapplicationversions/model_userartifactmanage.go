package galleryapplicationversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserArtifactManage struct {
	Install string  `json:"install"`
	Remove  string  `json:"remove"`
	Update  *string `json:"update,omitempty"`
}
