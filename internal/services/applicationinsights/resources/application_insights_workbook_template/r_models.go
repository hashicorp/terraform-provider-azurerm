// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_workbook_template

type ApplicationInsightsWorkbookTemplateModel struct {
	Name              string                         `tfschema:"name"`
	ResourceGroupName string                         `tfschema:"resource_group_name"`
	Author            string                         `tfschema:"author"`
	Galleries         []WorkbookTemplateGalleryModel `tfschema:"galleries"`
	Localized         string                         `tfschema:"localized"`
	Location          string                         `tfschema:"location"`
	Priority          int64                          `tfschema:"priority"`
	Tags              map[string]string              `tfschema:"tags"`
	TemplateData      string                         `tfschema:"template_data"`
}

type WorkbookTemplateGalleryModel struct {
	Name         string `tfschema:"name"`
	Category     string `tfschema:"category"`
	Order        int64  `tfschema:"order"`
	ResourceType string `tfschema:"resource_type"`
	Type         string `tfschema:"type"`
}

func (r Resource) ModelObject() any {
	return &ApplicationInsightsWorkbookTemplateModel{}
}
