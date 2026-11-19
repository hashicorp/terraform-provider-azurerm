// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_metadata

type MetadataModel struct {
	Name                     string                  `tfschema:"name"`
	WorkspaceId              string                  `tfschema:"workspace_id"`
	ContentId                string                  `tfschema:"content_id"`
	Kind                     string                  `tfschema:"kind"`
	ParentId                 string                  `tfschema:"parent_id"`
	Source                   []MetadataSourceModel   `tfschema:"source"`
	Author                   []MetadataAuthorModel   `tfschema:"author"`
	Support                  []MetadataSupportModel  `tfschema:"support"`
	Dependency               string                  `tfschema:"dependency"`
	Category                 []MetadataCategoryModel `tfschema:"category"`
	Providers                []string                `tfschema:"providers"`
	FirstPublishDate         string                  `tfschema:"first_publish_date"`
	LastPublishDate          string                  `tfschema:"last_publish_date"`
	ContentSchemaVersion     string                  `tfschema:"content_schema_version"`
	CustomVersion            string                  `tfschema:"custom_version"`
	IconId                   string                  `tfschema:"icon_id"`
	PreviewImages            []string                `tfschema:"preview_images"`
	PreviewImagesDark        []string                `tfschema:"preview_images_dark"`
	ThreatAnalysisTactics    []string                `tfschema:"threat_analysis_tactics"`
	ThreatAnalysisTechniques []string                `tfschema:"threat_analysis_techniques"`
	Version                  string                  `tfschema:"version"`
}

type MetadataSourceModel struct {
	Kind string `tfschema:"kind"`
	Name string `tfschema:"name"`
	Id   string `tfschema:"id"`
}

type MetadataAuthorModel struct {
	Name  string `tfschema:"name"`
	Email string `tfschema:"email"`
	Link  string `tfschema:"link"`
}

type MetadataSupportModel struct {
	Tier  string `tfschema:"tier"`
	Name  string `tfschema:"name"`
	Email string `tfschema:"email"`
	Link  string `tfschema:"link"`
}

type MetadataCategoryModel struct {
	Domains   []string `tfschema:"domains"`
	Verticals []string `tfschema:"verticals"`
}

func (a Resource) ModelObject() any {
	return &MetadataModel{}
}
