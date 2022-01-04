package consumption

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmConsumptionBudgetSubscriptionDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: resourceArmConsumptionBudgetSubscriptionDataSourceRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ConsumptionBudgetName(),
			},

			"subscription_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"amount": {
				Type:     pluginsdk.TypeFloat,
				Computed: true,
			},

			"filter": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"dimension": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"operator": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"values": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},
						"tag": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"operator": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"values": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},
						"not": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"dimension": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"name": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"operator": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"values": {
													Type:     pluginsdk.TypeList,
													Computed: true,
													Elem: &pluginsdk.Schema{
														Type: pluginsdk.TypeString,
													},
												},
											},
										},
									},
									"tag": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"name": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"operator": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"values": {
													Type:     pluginsdk.TypeList,
													Computed: true,
													Elem: &pluginsdk.Schema{
														Type: pluginsdk.TypeString,
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
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"threshold": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
						"threshold_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"operator": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"contact_emails": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"contact_groups": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"contact_roles": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"time_grain": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"time_period": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"start_date": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"end_date": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmConsumptionBudgetSubscriptionDataSourceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Consumption.BudgetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewConsumptionBudgetId(d.Get("subscription_id").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.Scope, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.Name)
	d.Set("subscription_id", id.Scope)
	if resp.Amount != nil {
		amount, _ := resp.Amount.Float64()
		d.Set("amount", amount)
	}
	d.Set("time_grain", string(resp.TimeGrain))
	d.Set("time_period", flattenConsumptionBudgetTimePeriod(resp.TimePeriod))
	d.Set("notification", flattenConsumptionBudgetNotifications(resp.Notifications, id.Scope))
	d.Set("filter", flattenConsumptionBudgetFilter(resp.Filter))

	return nil
}
