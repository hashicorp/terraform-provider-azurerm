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

			"profiles": {
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

						"rules": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							MaxItems: 10, // https://msdn.microsoft.com/en-us/library/azure/dn931928.aspx
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metricTrigger": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"metricName": {
													Type:     schema.TypeString,
													Required: true,
												},
												"metricResourceUri": {
													Type:     schema.TypeString,
													Required: true,
												},
												"timeGrain": {
													Type:     schema.TypeString,
													Required: true,
												},
												"statistic": {
													Type:     schema.TypeString,
													Required: true,
												},
												"timeWindow": {
													Type:     schema.TypeString,
													Required: true,
												},
												"timeAggregation": {
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
									"scaleAction": {
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

						"fixedDate": {
							Type:     schema.TypeSet,
							Optional: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"timeZone": {
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
												"timeZone": {
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
														Type: schema.TypeString,
													},
												},
												"minutes": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
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

			"notifications": {
				Type:     schema.TypeSet,
				Required: false,
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
									"sendToSubscriptionAdministrator": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"sendToSubscriptionCoAdministrator": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"customEmails": {
										Type:     schema.TypeList,
										Required: false,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"webhooks": {
							Type:     schema.TypeSet,
							Required: false,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"serviceUri": {
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

			"targetResourceUri": {
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

	parameters := insights.AutoscaleSettingResource{
		""

	}

	asClient.CreateOrUpdate(resourceGroupName, name, parameters)

	return nil
}

func resourceArmAutoscaleSettingsRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmAutoscaleSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
