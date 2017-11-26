package azurerm

import (
	"github.com/Azure/azure-sdk-for-go/arm/streamanalytics"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func streamAnalyticsInputSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(streamanalytics.TypeReference),
						string(streamanalytics.TypeStream),
					}, false),
				},
				"serialization": &schema.Schema{
					Type:     schema.TypeList,
					Required: true,
					MaxItems: 1,
					MinItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"type": &schema.Schema{
								Type:     schema.TypeString,
								Required: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(streamanalytics.TypeAvro),
									string(streamanalytics.TypeCsv),
									string(streamanalytics.TypeJSON),
								}, false),
							},
							"field_delimiter": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"encoding": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
								Default:  string(streamanalytics.UTF8),
								ValidateFunc: validation.StringInSlice([]string{
									string(streamanalytics.UTF8),
								}, false),
							},
						},
					},
				},
				"datasource": &schema.Schema{
					Type:     schema.TypeList,
					MaxItems: 1,
					Required: true,
					MinItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"blob": &schema.Schema{
								Type:     schema.TypeList,
								MaxItems: 1,
								Optional: true,
								MinItems: 1,
								ConflictsWith: []string{
									"event_hub",
									"iot_hub",
								},
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"storage_account_name": &schema.Schema{
											Type:     schema.TypeString,
											Required: true,
										},
										"storage_account_key": &schema.Schema{
											Type:     schema.TypeString,
											Required: true,
										},
										"container": &schema.Schema{
											Type:     schema.TypeString,
											Required: true,
										},
										"path_pattern": &schema.Schema{
											Type:     schema.TypeString,
											Required: true,
										},
										"date_format": &schema.Schema{
											Type:     schema.TypeString,
											Optional: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
