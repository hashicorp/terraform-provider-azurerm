// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package consumption_budget_resource_group

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	consumptionClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Consumption Budget for a specific Resource Group.",
		Read:        resourceArmConsumptionBudgetResourceGroupDataSourceRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of this Consumption Budget.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ConsumptionBudgetName(),
			},

			"resource_group_id": {
				Description: "The ID of the subscription.",
				Type:        pluginsdk.TypeString,
				Required:    true,
			},

			"amount": {
				Description: "The total amount of cost to track with the budget.",
				Type:        pluginsdk.TypeFloat,
				Computed:    true,
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
										Description: "The name of the column to use for the filter.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"operator": {
										Description: "The operator to use for comparison.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
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
										Description: "The name of the tag used for the filter.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"operator": {
										Description: "The operator used for comparison.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"values": {
										Description: "A list of values for the tag.",
										Type:        pluginsdk.TypeList,
										Computed:    true,
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
													Description: "The name of the column to use for the filter.",
													Type:        pluginsdk.TypeString,
													Computed:    true,
												},
												"operator": {
													Description: "The operator to use for comparison.",
													Type:        pluginsdk.TypeString,
													Computed:    true,
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
													Description: "The name of the tag used for the filter.",
													Type:        pluginsdk.TypeString,
													Computed:    true,
												},
												"operator": {
													Description: "The operator used for comparison.",
													Type:        pluginsdk.TypeString,
													Computed:    true,
												},
												"values": {
													Description: "A list of values for the tag.",
													Type:        pluginsdk.TypeList,
													Computed:    true,
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
							Description: "Whether the notification is enabled.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},
						"threshold": {
							Description: "Threshold value associated with the notification.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},
						"threshold_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"operator": {
							Description: "The comparison operator for the notification.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"contact_emails": {
							Description: "A list of email addresses to send the budget notification to when the threshold is exceeded.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"contact_groups": {
							Description: "A list of Action Group IDs to send the budget notification to when the threshold is exceeded.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"contact_roles": {
							Description: "A list of contact roles to send the budget notification to when the threshold is exceeded.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"time_grain": {
				Description: "The time covered by a budget.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"time_period": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"start_date": {
							Description: "The start date for the budget.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"end_date": {
							Description: "The end date for the budget.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceArmConsumptionBudgetResourceGroupDataSourceRead(d *pluginsdk.ResourceData, meta any) error {
	client := consumptionClient.From(meta).BudgetsClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id := budgets.NewScopedBudgetID(d.Get("resource_group_id").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making read request on %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.BudgetName)
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("amount", props.Amount)
			d.Set("time_grain", string(props.TimeGrain))
			d.Set("time_period", helpers.FlattenConsumptionBudgetTimePeriod(&props.TimePeriod))
			d.Set("notification", helpers.FlattenConsumptionBudgetNotifications(props.Notifications, id.Scope))
			d.Set("filter", helpers.FlattenConsumptionBudgetFilter(props.Filter))
		}
	}

	return nil
}
