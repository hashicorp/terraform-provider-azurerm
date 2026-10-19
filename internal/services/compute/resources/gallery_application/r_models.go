// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package gallery_application

type GalleryApplicationModel struct {
	Name                string            `tfschema:"name"`
	GalleryId           string            `tfschema:"gallery_id"`
	Location            string            `tfschema:"location"`
	SupportedOSType     string            `tfschema:"supported_os_type"`
	Description         string            `tfschema:"description"`
	EndOfLifeDate       string            `tfschema:"end_of_life_date"`
	Eula                string            `tfschema:"eula"`
	PrivacyStatementURI string            `tfschema:"privacy_statement_uri"`
	ReleaseNoteURI      string            `tfschema:"release_note_uri"`
	Tags                map[string]string `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &GalleryApplicationModel{}
}
