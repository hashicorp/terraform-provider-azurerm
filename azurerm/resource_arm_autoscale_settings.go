package azurerm

import (
	"github.com/Azure/azure-sdk-for-go/arm/insights"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmAutoscaleSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutoscaleSettingsCreateOrUpdate,
		Read:   resourceArmAutoscaleSettingsRead,
		Update: resourceArmAutoscaleSettingsCreateOrUpdate,
		Delete: resourceArmAutoscaleSettingsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"profile": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 20, // https://msdn.microsoft.com/en-us/library/azure/dn931928.aspx
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"capacity": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"minimum": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"maximum": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"default": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},

						"rule": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							MaxItems: 10, // https://msdn.microsoft.com/en-us/library/azure/dn931928.aspx
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_trigger": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"metric_name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"metric_resource_uri": {
													Type:     schema.TypeString,
													Required: true,
												},
												"time_grain": {
													Type:     schema.TypeString,
													Required: true,
												},
												"statistic": {
													Type:     schema.TypeString,
													Required: true,
												},
												"time_window": {
													Type:     schema.TypeString,
													Required: true,
												},
												"time_aggregation": {
													Type:     schema.TypeString,
													Required: true,
												},
												"operator": {
													Type:     schema.TypeString,
													Required: true,
												},
												"threshold": {
													Type:     schema.TypeInt,
													Required: true,
												},
											},
										},
									},
									"scale_action": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"direction": {
													Type:     schema.TypeString,
													Required: true,
												},
												"type": {
													Type:     schema.TypeString,
													Required: true,
												},
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
												"cooldown": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},

						"fixed_date": {
							Type:     schema.TypeSet,
							Optional: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time_zone": {
										Type:     schema.TypeString,
										Required: true,
									},
									"start": {
										Type:     schema.TypeString,
										Required: true,
									},
									"end": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},

						"recurrence": {
							Type:     schema.TypeSet,
							Optional: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"frequency": {
										Type:     schema.TypeString,
										Required: true,
									},
									"schedule": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"time_zone": {
													Type:     schema.TypeString,
													Required: true,
												},
												"days": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"hours": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"minutes": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"notification": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operation": {
							Type:     schema.TypeString,
							Required: true,
						},
						"email": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"send_to_subscription_administrator": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"send_to_subscription_co_administrator": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"custom_emails": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"webhook": {
							Type:     schema.TypeSet,
							Optional: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_uri": {
										Type:     schema.TypeString,
										Required: true,
									},
									"properties": {
										Type:     schema.TypeMap,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"target_resource_uri": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmAutoscaleSettingsCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	asClient := armClient.autoscaleSettingsClient

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)
	resourceType := "Microsoft.Insights/autoscaleSettings"
	location := d.Get("location").(string)
	enabled := d.Get("enabled").(bool)
	targetResourceURI := d.Get("target_resource_uri").(string)
	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)
	autoscaleSettings := insights.AutoscaleSetting{
		Name:              &name,
		Enabled:           &enabled,
		TargetResourceURI: &targetResourceURI,
		// Profiles:
		// Notifications:
	}

	parameters := insights.AutoscaleSettingResource{
		Name:             &name,
		Type:             &resourceType,
		Location:         &location,
		Tags:             expandedTags,
		AutoscaleSetting: &autoscaleSettings,
	}

	result, err := asClient.CreateOrUpdate(resourceGroupName, name, parameters)
	if err != nil {
		return err
	}

	d.SetId(*result.ID)

	return resourceArmAutoscaleSettingsRead(d, meta)
}

func resourceArmAutoscaleSettingsRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmAutoscaleSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	asClient := armClient.autoscaleSettingsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroupName := id.ResourceGroup
	autoscaleSettingName := id.Path["autoscalesettings"]

	_, err = asClient.Delete(resourceGroupName, autoscaleSettingName)
	return err
}

// func expandAzureRmAutoscaleSettingProfile(d *schema.ResourceData) (*[]insights.AutoscaleProfile, error) {
// 	d.Get()
// }
