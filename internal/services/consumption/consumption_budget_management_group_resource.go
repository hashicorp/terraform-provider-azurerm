package consumption

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2019-10-01/consumption"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/parse"
	managementGroupParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmConsumptionBudgetManagementGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmConsumptionBudgetManagementGroupCreateUpdate,
		Read:   resourceArmConsumptionBudgetManagementGroupRead,
		Update: resourceArmConsumptionBudgetManagementGroupCreateUpdate,
		Delete: resourceArmConsumptionBudgetManagementGroupDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ConsumptionBudgetManagementGroupID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: SchemaConsumptionBudgetManagementGroupResource(),
	}
}

func resourceArmConsumptionBudgetManagementGroupCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	managementGroupId, err := managementGroupParse.ManagementGroupID(d.Get("management_group_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewConsumptionBudgetManagementGroupID(managementGroupId.Name, d.Get("name").(string))

	err = resourceArmConsumptionBudgetCreateUpdate(d, meta, consumptionBudgetManagementGroupName, managementGroupId.ID())
	if err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceArmConsumptionBudgetManagementGroupRead(d, meta)
}

func resourceArmConsumptionBudgetManagementGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	consumptionBudgetId, err := parse.ConsumptionBudgetManagementGroupID(d.Id())
	if err != nil {
		return err
	}

	managementGroupId := managementGroupParse.NewManagementGroupId(consumptionBudgetId.ManagementGroupName)

	err = resourceArmConsumptionBudgetRead(d, meta, managementGroupId.ID(), consumptionBudgetId.BudgetName, SchemaConsumptionBudgetNotificationManagementGroupElement, flattenConsumptionBudgetManagementGroupNotifications)
	if err != nil {
		return err
	}

	d.Set("management_group_id", managementGroupId.ID())

	return nil
}

func resourceArmConsumptionBudgetManagementGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	consumptionBudgetId, err := parse.ConsumptionBudgetManagementGroupID(d.Id())
	if err != nil {
		return err
	}

	managementGroupId := managementGroupParse.NewManagementGroupId(consumptionBudgetId.ManagementGroupName)

	return resourceArmConsumptionBudgetDelete(d, meta, managementGroupId.ID())
}

func flattenConsumptionBudgetManagementGroupNotifications(input map[string]*consumption.Notification) []interface{} {
	notifications := make([]interface{}, 0)

	if input == nil {
		return notifications
	}

	for _, v := range input {
		if v != nil {
			notificationBlock := make(map[string]interface{})

			notificationBlock["enabled"] = *v.Enabled
			notificationBlock["operator"] = string(v.Operator)
			threshold, _ := v.Threshold.Float64()
			notificationBlock["threshold"] = int(threshold)
			notificationBlock["threshold_type"] = string(v.ThresholdType)
			notificationBlock["contact_emails"] = utils.FlattenStringSlice(v.ContactEmails)

			notifications = append(notifications, notificationBlock)
		}
	}

	return notifications
}
