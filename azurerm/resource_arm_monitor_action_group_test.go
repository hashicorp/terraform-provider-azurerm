package azurerm

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMMonitorActionGroup_basic(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "short_name", "acctestag"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMMonitorActionGroup_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_app_service_custom_hostname_binding"),
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_emailReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMonitorActionGroup_emailReceiver(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.0.email_address", "admin@contoso.com"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_itsmReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMonitorActionGroup_itsmReceiver(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.0.workspace_id", "6eee3a18-aac3-40e4-b98e-1f309f329816"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_azureAppPushReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMonitorActionGroup_azureAppPushReceiver(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.0.email_address", "admin@contoso.com"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_smsReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMonitorActionGroup_smsReceiver(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.0.country_code", "1"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.0.phone_number", "1231231234"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_webhookReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMonitorActionGroup_webhookReceiver(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.0.service_uri", "http://example.com/alert"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_automationRunbookReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMonitorActionGroup_automationRunbookReceiver(ri, testLocation())

	aaName := fmt.Sprintf("acctestAA-%d", ri)
	webhookName := "webhook_alert"
	resGroup := fmt.Sprintf("acctestRG-%d", ri)
	aaWebhookResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/AutomationAccounts/%s/webhooks/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, aaName, webhookName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.webhook_resource_id", aaWebhookResourceID),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_voiceReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMonitorActionGroup_voiceReceiver(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.0.country_code", "1"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.0.phone_number", "1231231234"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_logicAppReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMonitorActionGroup_logicAppReceiver(ri, testLocation())

	laName := fmt.Sprintf("acctestLA-%d", ri)
	resGroup := fmt.Sprintf("acctestRG-%d", ri)
	laResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/workflows/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, laName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.resource_id", laResourceID),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_azureFunctionReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMonitorActionGroup_azureFunctionReceiver(ri, testLocation())

	faName := fmt.Sprintf("acctestFA-%d", ri)
	resGroup := fmt.Sprintf("acctestRG-%d", ri)
	faResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, faName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.function_app_resource_id", faResourceID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_complete(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMonitorActionGroup_complete(ri, testLocation())

	aaName := fmt.Sprintf("acctestAA-%d", ri)
	faName := fmt.Sprintf("acctestFA-%d", ri)
	laName := fmt.Sprintf("acctestLA-%d", ri)
	webhookName := "webhook_alert"
	resGroup := fmt.Sprintf("acctestRG-%d", ri)
	aaWebhookResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/AutomationAccounts/%s/webhooks/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, aaName, webhookName)
	faResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, faName)
	laResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/workflows/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, laName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.0.email_address", "admin@contoso.com"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.1.email_address", "devops@contoso.com"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.0.workspace_id", "6eee3a18-aac3-40e4-b98e-1f309f329816"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.0.connection_id", "53de6956-42b4-41ba-be3c-b154cdf17b13"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.0.ticket_configuration", "{}"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.0.region", "eastus"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.0.email_address", "admin@contoso.com"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.0.country_code", "1"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.0.phone_number", "1231231234"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.1.country_code", "86"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.1.phone_number", "13888888888"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.0.service_uri", "http://example.com/alert"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.1.service_uri", "https://backup.example.com/warning"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.webhook_resource_id", aaWebhookResourceID),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.0.country_code", "1"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.0.phone_number", "1231231234"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.resource_id", laResourceID),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.function_app_resource_id", faResourceID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_disabledUpdate(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMMonitorActionGroup_disabledBasic(ri, location)
	postConfig := testAccAzureRMMonitorActionGroup_basic(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_singleReceiverUpdate(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	emailConfig := testAccAzureRMMonitorActionGroup_emailReceiver(ri, location)
	itsmConfig := testAccAzureRMMonitorActionGroup_itsmReceiver(ri, location)
	azureAppPushConfig := testAccAzureRMMonitorActionGroup_azureAppPushReceiver(ri, location)
	smsConfig := testAccAzureRMMonitorActionGroup_smsReceiver(ri, location)
	webhookConfig := testAccAzureRMMonitorActionGroup_webhookReceiver(ri, location)
	automationRunbookConfig := testAccAzureRMMonitorActionGroup_automationRunbookReceiver(ri, location)
	voiceConfig := testAccAzureRMMonitorActionGroup_voiceReceiver(ri, location)
	logicAppConfig := testAccAzureRMMonitorActionGroup_logicAppReceiver(ri, location)
	azureFunctionConfig := testAccAzureRMMonitorActionGroup_azureFunctionReceiver(ri, location)

	aaName := fmt.Sprintf("acctestAA-%d", ri)
	faName := fmt.Sprintf("acctestFA-%d", ri)
	laName := fmt.Sprintf("acctestLA-%d", ri)
	webhookName := "webhook_alert"
	resGroup := fmt.Sprintf("acctestRG-%d", ri)
	aaResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/AutomationAccounts/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, aaName)
	aaWebhookResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/AutomationAccounts/%s/webhooks/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, aaName, webhookName)
	faResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, faName)
	laResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/workflows/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, laName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: emailConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.0.email_address", "admin@contoso.com"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
			{
				Config: itsmConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.0.workspace_id", "6eee3a18-aac3-40e4-b98e-1f309f329816"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.0.connection_id", "53de6956-42b4-41ba-be3c-b154cdf17b13"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.0.ticket_configuration", "{}"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.0.region", "eastus"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
			{
				Config: azureAppPushConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.0.email_address", "admin@contoso.com"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
			{
				Config: smsConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.0.country_code", "1"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.0.phone_number", "1231231234"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
			{
				Config: webhookConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.0.service_uri", "http://example.com/alert"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
			{
				Config: automationRunbookConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.automation_account_id", aaResourceID),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.runbook_name", webhookName),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.webhook_resource_id", aaWebhookResourceID),
					resource.TestCheckResourceAttrSet(resourceName, "automation_runbook_receiver.is_global_runbook"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.service_uri", "https://s13events.azure-automation.net/webhooks?token=randomtoken"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
			{
				Config: voiceConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.0.phone_number", "1231231234"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
			{
				Config: logicAppConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.resource_id", laResourceID),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.callback_url", "http://test-host:100/workflows/fb9c8d79b15f41ce9b12861862f43546/versions/08587100027316071865/triggers/manualTrigger/paths/invoke?api-version=2015-08-01-preview&sp=%2Fversions%2F08587100027316071865%2Ftriggers%2FmanualTrigger%2Frun&sv=1.0&sig=IxEQ_ygZf6WNEQCbjV0Vs6p6Y4DyNEJVAa86U5B4xhk"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
			{
				Config: azureFunctionConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.function_app_resource_id", faResourceID),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.function_name", "myfunc"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.http_trigger_url", "https://example.com/trigger"),
				),
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_multipleReceiversUpdate(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	basicConfig := testAccAzureRMMonitorActionGroup_basic(ri, location)
	completeConfig := testAccAzureRMMonitorActionGroup_complete(ri, location)

	aaName := fmt.Sprintf("acctestAA-%d", ri)
	faName := fmt.Sprintf("acctestFA-%d", ri)
	laName := fmt.Sprintf("acctestLA-%d", ri)
	webhookName := "webhook_alert"
	resGroup := fmt.Sprintf("acctestRG-%d", ri)
	aaResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/AutomationAccounts/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, aaName)
	aaWebhookResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/AutomationAccounts/%s/webhooks/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, aaName, webhookName)
	faResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, faName)
	laResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/workflows/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, laName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: basicConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
			{
				Config: completeConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.0.email_address", "admin@contoso.com"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.1.email_address", "devops@contoso.com"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.0.workspace_id", "6eee3a18-aac3-40e4-b98e-1f309f329816"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.0.connection_id", "53de6956-42b4-41ba-be3c-b154cdf17b13"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.0.ticket_configuration", "{}"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.0.region", "eastus"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.0.email_address", "admin@contoso.com"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.0.country_code", "1"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.0.phone_number", "1231231234"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.1.country_code", "86"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.1.phone_number", "13888888888"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.0.service_uri", "http://example.com/alert"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.1.service_uri", "https://backup.example.com/warning"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.automation_account_id", aaResourceID),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.runbook_name", webhookName),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.webhook_resource_id", aaWebhookResourceID),
					resource.TestCheckResourceAttrSet(resourceName, "automation_runbook_receiver.is_global_runbook"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.service_uri", "https://s13events.azure-automation.net/webhooks?token=randomtoken"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.0.country_code", "1"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.0.phone_number", "1231231234"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.1.country_code", "86"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.1.phone_number", "13888888888"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.resource_id", laResourceID),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.callback_url", "http://test-host:100/workflows/fb9c8d79b15f41ce9b12861862f43546/versions/08587100027316071865/triggers/manualTrigger/paths/invoke?api-version=2015-08-01-preview&sp=%2Fversions%2F08587100027316071865%2Ftriggers%2FmanualTrigger%2Frun&sv=1.0&sig=IxEQ_ygZf6WNEQCbjV0Vs6p6Y4DyNEJVAa86U5B4xhk"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.function_app_resource_id", faResourceID),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.function_name", "myfunc"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.http_trigger_url", "https://example.com/trigger"),
				),
			},
			{
				Config: basicConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_function_receiver.#", "0"),
				),
			},
		},
	})
}

func testAccAzureRMMonitorActionGroup_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_requiresImport(rInt int, location string) string {
	template := testAccAzureRMMonitorActionGroup_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_group" "import" {
  name                = "${azurerm_monitor_action_group.test.name}"
  resource_group_name = "${azurerm_monitor_action_group.test.resource_group_name}"
  short_name          = "${azurerm_monitor_action_group.test.short_name}"
}
`, template)
}

func testAccAzureRMMonitorActionGroup_emailReceiver(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"

  email_receiver {
    name          = "sendtoadmin"
    email_address = "admin@contoso.com"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_itsmReceiver(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"

  itsm_receiver {
		name          = "createorupdateticket"
		workspace_id = "6eee3a18-aac3-40e4-b98e-1f309f329816"
		connection_id = "53de6956-42b4-41ba-be3c-b154cdf17b13"
		ticket_configuration = "{}"
		region = "eastus"
	}
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_azureAppPushReceiver(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
	short_name          = "acctestag"

	azure_app_push_receiver {
		name          = "pushtoadmin"
		email_address = "admin@contoso.com"
	}
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_smsReceiver(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"

  sms_receiver {
    name         = "oncallmsg"
    country_code = "1"
    phone_number = "1231231234"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_webhookReceiver(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"

  webhook_receiver {
    name        = "callmyapiaswell"
    service_uri = "http://example.com/alert"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_automationRunbookReceiver(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"

  automation_runbook_receiver {
		name = "action_name_1"
		automation_account_id = "${azurerm_automation_account.test.id}"
  	runbook_name = "${azurerm_automation_runbook.test.name}"
		webhook_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg-runbooks/providers/microsoft.automation/automationaccounts/aaa001/webhooks/webhook_alert"
		is_global_runbook = true
		service_uri = "https://s13events.azure-automation.net/webhooks?token=randomtoken"
	}
}

resource "azurerm_automation_account" "test" {
	name                = "acctestAA-%d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"

	sku {
		name = "Basic"
	}
}

resource "azurerm_automation_runbook" "test" {
	name                = "Get-AzureVMTutorial"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	account_name        = "${azurerm_automation_account.test.name}"
	log_verbose         = "true"
	log_progress        = "true"
	description         = "This is an test runbook"
	runbook_type        = "PowerShellWorkflow"

	publish_content_link {
		uri = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/101-automation-runbook-getvms/Runbooks/Get-AzureVMTutorial.ps1"
	}
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMMonitorActionGroup_voiceReceiver(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"

	voice_receiver {
		name         = "oncallmsg"
		country_code = "1"
		phone_number = "1231231234"
	}
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_logicAppReceiver(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
	short_name          = "acctestag"

	logic_app_receiver {
		name = "logicappaction"
		resource_id = "${azurerm_logic_app_workflow.test.id}"
		callback_url = "http://test-host:100/workflows/fb9c8d79b15f41ce9b12861862f43546/versions/08587100027316071865/triggers/manualTrigger/paths/invoke?api-version=2015-08-01-preview&sp=%%2Fversions%%2F08587100027316071865%%2Ftriggers%%2FmanualTrigger%%2Frun&sv=1.0&sig=IxEQ_ygZf6WNEQCbjV0Vs6p6Y4DyNEJVAa86U5B4xhk"
	}
}

resource "azurerm_logic_app_workflow" "test" {
	name                = "acctestLA-%d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_logic_app_trigger_http_request" "test" {
	name         = "some-http-trigger"
	logic_app_id = "${azurerm_logic_app_workflow.test.id}"

schema = <<SCHEMA
{
	"type": "object",
	"properties": {
		"hello": {
			"type": "string"
		}
	}
}
SCHEMA
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMMonitorActionGroup_azureFunctionReceiver(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"

	azure_function_receiver {
		name = "funcaction"
		function_app_resource_id = "${azurerm_function_app.test.id}"
		function_name = "myfunc"
		http_trigger_url = "https://example.com/trigger"
	}
}

resource "azurerm_storage_account" "test" {
	name                     = "acctestSA-%d"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
	name                = "acctestSP-%d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"

	sku {
		tier = "Standard"
		size = "S1"
	}
}

resource "azurerm_function_app" "test" {
	name                      = "acctestFA-%d"
	location                  = "${azurerm_resource_group.test.location}"
	resource_group_name       = "${azurerm_resource_group.test.name}"
	app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
	storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
}
`, rInt, location, rInt, rInt/1e12, rInt, rInt)
}

func testAccAzureRMMonitorActionGroup_complete(rInt int, location string) string {
	return fmt.Sprintf(`
	resource "azurerm_resource_group" "test" {
		name     = "acctestRG-%d"
		location = "%s"
	}

	resource "azurerm_monitor_action_group" "test" {
		name                = "acctestActionGroup-%d"
		resource_group_name = "${azurerm_resource_group.test.name}"
		short_name          = "acctestag"

		email_receiver {
			name          = "sendtoadmin"
			email_address = "admin@contoso.com"
		}

		email_receiver {
			name          = "sendtodevops"
			email_address = "devops@contoso.com"
		}

		itsm_receiver {
			name          = "createorupdateticket"
			workspace_id = "6eee3a18-aac3-40e4-b98e-1f309f329816"
			connection_id = "53de6956-42b4-41ba-be3c-b154cdf17b13"
			ticket_configuration = "{}"
			region = "eastus"
		}

		azure_app_push_receiver {
			name          = "pushtoadmin"
			email_address = "admin@contoso.com"
		}

		sms_receiver {
			name         = "oncallmsg"
			country_code = "1"
			phone_number = "1231231234"
		}

		sms_receiver {
			name         = "remotesupport"
			country_code = "86"
			phone_number = "13888888888"
		}

		webhook_receiver {
			name        = "callmyapiaswell"
			service_uri = "http://example.com/alert"
		}

		webhook_receiver {
			name        = "callmybackupapi"
			service_uri = "https://backup.example.com/warning"
		}

		automation_runbook_receiver {
			name = "action_name_1"
			automation_account_id = "${azurerm_automation_account.test.id}"
    	runbook_name = "${azurerm_automation_runbook.test.name}"
			webhook_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg-runbooks/providers/microsoft.automation/automationaccounts/aaa001/webhooks/webhook_alert"
			is_global_runbook = true
			service_uri = "https://s13events.azure-automation.net/webhooks?token=randomtoken"
		}

		voice_receiver {
			name         = "oncallmsg"
			country_code = "1"
			phone_number = "1231231234"
		}

		logic_app_receiver {
			name = "logicappaction"
			resource_id = "${azurerm_logic_app_workflow.test.id}"
			callback_url = "http://test-host:100/workflows/fb9c8d79b15f41ce9b12861862f43546/versions/08587100027316071865/triggers/manualTrigger/paths/invoke?api-version=2015-08-01-preview&sp=%%2Fversions%%2F08587100027316071865%%2Ftriggers%%2FmanualTrigger%%2Frun&sv=1.0&sig=IxEQ_ygZf6WNEQCbjV0Vs6p6Y4DyNEJVAa86U5B4xhk"
		}

		azure_function_receiver {
			name = "funcaction"
			function_app_resource_id = "${azurerm_function_app.test.id}"
			function_name = "myfunc"
			http_trigger_url = "https://example.com/trigger"
		}

	}

	resource "azurerm_automation_account" "test" {
		name                = "acctestAA-%d"
		location            = "${azurerm_resource_group.test.location}"
		resource_group_name = "${azurerm_resource_group.test.name}"

		sku {
		  name = "Basic"
		}
	}

	resource "azurerm_automation_runbook" "test" {
		name                = "Get-AzureVMTutorial"
		location            = "${azurerm_resource_group.test.location}"
		resource_group_name = "${azurerm_resource_group.test.name}"
		account_name        = "${azurerm_automation_account.test.name}"
		log_verbose         = "true"
		log_progress        = "true"
		description         = "This is an test runbook"
		runbook_type        = "PowerShellWorkflow"

		publish_content_link {
		  uri = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/101-automation-runbook-getvms/Runbooks/Get-AzureVMTutorial.ps1"
		}
	}

	resource "azurerm_logic_app_workflow" "test" {
		name                = "acctestLA-%d"
		location            = "${azurerm_resource_group.test.location}"
		resource_group_name = "${azurerm_resource_group.test.name}"
	}

	resource "azurerm_logic_app_trigger_http_request" "test" {
		name         = "some-http-trigger"
		logic_app_id = "${azurerm_logic_app_workflow.test.id}"

	schema = <<SCHEMA
	{
		"type": "object",
		"properties": {
			"hello": {
				"type": "string"
			}
		}
	}
	SCHEMA
	}

	resource "azurerm_storage_account" "test" {
		name                     = "acctestSA-%d"
		resource_group_name      = "${azurerm_resource_group.test.name}"
		location                 = "${azurerm_resource_group.test.location}"
		account_tier             = "Standard"
		account_replication_type = "LRS"
	}

	resource "azurerm_app_service_plan" "test" {
		name                = "acctestSP-%d"
		location            = "${azurerm_resource_group.test.location}"
		resource_group_name = "${azurerm_resource_group.test.name}"

		sku {
			tier = "Standard"
			size = "S1"
		}
	}

	resource "azurerm_function_app" "test" {
		name                      = "acctestFA-%d"
		location                  = "${azurerm_resource_group.test.location}"
		resource_group_name       = "${azurerm_resource_group.test.name}"
		app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
		storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
	}
`, rInt, location, rInt, rInt, rInt, rInt/1e12, rInt, rInt)
}

func testAccAzureRMMonitorActionGroup_disabledBasic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"
  enabled             = false
}
`, rInt, location, rInt)
}

func testCheckAzureRMMonitorActionGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).monitor.ActionGroupsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_action_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Action Group still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMMonitorActionGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Action Group Instance: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).monitor.ActionGroupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on monitorActionGroupsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Action Group Instance %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}
