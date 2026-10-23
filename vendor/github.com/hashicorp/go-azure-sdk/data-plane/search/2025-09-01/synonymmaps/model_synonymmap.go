package synonymmaps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SynonymMap struct {
	EncryptionKey *SearchResourceEncryptionKey `json:"encryptionKey,omitempty"`
	Format        SynonymMapFormat             `json:"format"`
	Name          string                       `json:"name"`
	OdataEtag     *string                      `json:"@odata.etag,omitempty"`
	Synonyms      string                       `json:"synonyms"`
}
