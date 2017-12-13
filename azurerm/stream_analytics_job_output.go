package azurerm

import (
	"github.com/Azure/azure-sdk-for-go/arm/streamanalytics"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	DelimSpace       = " "
	DelimComma       = ","
	DelimTab         = "\t"
	DelimSemiColon   = ";"
	DelimVerticalBar = "|"
)

func streamAnalyticsOutputSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: false,
		Elem: map[string]*schema.Schema{
			"serialization": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(streamanalytics.TypeCsv),
								string(streamanalytics.TypeAvro),
								string(streamanalytics.TypeJSON),
							}, false),
						},
						"field_delimiter": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								DelimTab,
								DelimComma,
								DelimSemiColon,
								DelimVerticalBar,
								DelimSpace,
							}, false),
						},
						"encoding": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"format": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(streamanalytics.Array),
								string(streamanalytics.LineSeparated),
							}, false),
						},
					},
				},
			},
			"datasource": &schema.Schema{
				Type:     schema.TypeList,
				Optional: false,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blob":  &schema.Schema{},
						"table": &schema.Schema{},
					},
				},
			},
		},
	}
}
