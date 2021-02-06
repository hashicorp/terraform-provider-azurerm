package consumption

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/consumption/parse"

	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2019-01-01/consumption"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/shopspring/decimal"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmConsumptionBudgetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Consumption.BudgetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ConsumptionBudgetID(d.Id())
	if err != nil {
		return err
	}

	scope := id.Scope
	name := id.Name
	resp, err := client.Get(ctx, scope, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error making read request on Azure Consumption Budget %q for scope %q: %+v", name, scope, err)
	}

	d.Set("name", resp.Name)
	d.Set("category", string(resp.Category))
	amount, _ := resp.Amount.Float64()
	d.Set("amount", amount)
	d.Set("time_grain", string(resp.TimeGrain))
	d.Set("time_period", FlattenConsumptionBudgetTimePeriod(resp.TimePeriod))
	d.Set("notification", schema.NewSet(schema.HashResource(SchemaAzureConsumptionBudgetNotificationElement()), FlattenConsumptionBudgetNotifications(resp.Notifications)))
	d.Set("filter", FlattenConsumptionBudgetFilter(resp.Filters))

	return nil
}

func resourceArmConsumptionBudgetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Consumption.BudgetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ConsumptionBudgetID(d.Id())
	if err != nil {
		return err
	}

	scope := id.Scope
	name := id.Name

	resp, err := client.Delete(ctx, scope, name)

	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("error issuing delete request on Azure Consumption Budget %q for scope %q: %+v", name, scope, err)
		}
	}

	return nil
}

func resourceArmConsumptionBudgetCreateUpdate(d *schema.ResourceData, meta interface{}, resourceName, scope string) error {
	client := meta.(*clients.Client).Consumption.BudgetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, scope, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("error checking for presence of existing Consumption Budget %q for scope %q: %s", name, scope, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError(resourceName, *existing.ID)
		}
	}

	amount := decimal.NewFromFloat(d.Get("amount").(float64))
	timePeriod, err := ExpandConsumptionBudgetTimePeriod(d.Get("time_period").([]interface{}))
	if err != nil {
		return fmt.Errorf("error in expanding`time_period`: %+v", err)
	}

	parameters := consumption.Budget{
		Name: utils.String(name),
		BudgetProperties: &consumption.BudgetProperties{
			Amount:        &amount,
			Category:      consumption.CategoryType(d.Get("category").(string)),
			Filters:       ExpandConsumptionBudgetFilter(d.Get("filter").([]interface{})),
			Notifications: ExpandConsumptionBudgetNotifications(d.Get("notification").(*schema.Set).List()),
			TimeGrain:     consumption.TimeGrainType(d.Get("time_grain").(string)),
			TimePeriod:    timePeriod,
		},
	}

	read, err := client.CreateOrUpdate(ctx, scope, name, parameters)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("cannot read Azure Consumption Budget %q for scope %q", name, scope)
	}

	d.SetId(*read.ID)

	return nil
}
