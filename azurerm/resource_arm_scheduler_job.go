// nolint: megacheck
// entire automation SDK has been depreciated in v21.3 in favor of logic apps, an entirely different service.
package azurerm

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/scheduler/mgmt/2016-03-01/scheduler"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSchedulerJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSchedulerJobCreateUpdate,
		Read:   resourceArmSchedulerJobRead,
		Update: resourceArmSchedulerJobCreateUpdate,
		Delete: resourceArmSchedulerJobDelete,

		DeprecationMessage: "Scheduler Job's have been deprecated in favour of Logic Apps - more information can be found at https://docs.microsoft.com/en-us/azure/scheduler/migrate-from-scheduler-to-logic-apps",

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: resourceArmSchedulerJobCustomizeDiff,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z][-_a-zA-Z0-9].*$"),
					"Job Collection Name name must start with a letter and contain only letters, numbers, hyphens and underscores.",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"job_collection_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			//actions
			"action_web": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				Elem:          resourceArmSchedulerJobActionWebSchema("action_web"),
				ConflictsWith: []string{"action_storage_queue"},
			},

			"action_storage_queue": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				Elem:          resourceArmSchedulerJobActionStorageSchema(),
				ConflictsWith: []string{"action_web"},
			},

			//actions
			"error_action_web": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				Elem:          resourceArmSchedulerJobActionWebSchema("error_action_web"),
				ConflictsWith: []string{"error_action_storage_queue"},
			},

			"error_action_storage_queue": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				Elem:          resourceArmSchedulerJobActionStorageSchema(),
				ConflictsWith: []string{"error_action_web"},
			},

			//retry policy
			"retry": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						//silently fails if the duration is not in the correct format
						//todo validation
						"interval": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "00:00:30",
							ValidateFunc: validate.NoEmptyStrings,
						},

						"count": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      4,
							ValidateFunc: validation.IntBetween(1, 20),
						},
					},
				},
			},

			//recurrences (schedule in portal, recurrence in API)
			"recurrence": {
				Type:     schema.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"frequency": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(scheduler.Minute),
								string(scheduler.Hour),
								string(scheduler.Day),
								string(scheduler.Week),
								string(scheduler.Month),
							}, true),
						},

						"interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1, //defaults to 1 in the portal

							//maximum is dynamic:  1 min <= interval * frequency <= 500 days (bounded by JobCollection quotas)
							ValidateFunc: validation.IntAtLeast(1),
						},

						"count": {
							Type:     schema.TypeInt,
							Optional: true,
							//silently fails/produces odd results at >2147483647
							ValidateFunc: validation.IntBetween(1, 2147483647),
						},

						"end_time": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppress.RFC3339Time,
							ValidateFunc:     validate.RFC3339Time,
						},

						"minutes": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								ValidateFunc: validation.IntBetween(0, 59),
							},
							Set: set.HashInt,
						},

						"hours": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								ValidateFunc: validation.IntBetween(0, 23),
							},
							Set: set.HashInt,
						},

						"week_days": { //used with weekly
							Type:          schema.TypeSet,
							Optional:      true,
							ConflictsWith: []string{"recurrence.0.month_days", "recurrence.0.monthly_occurrences"},
							// the constants are title cased but the API returns all lowercase
							// so lets ignore the case
							Set: set.HashStringIgnoreCase,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								DiffSuppressFunc: suppress.CaseDifference,
								ValidateFunc:     validate.DayOfTheWeek(true),
							},
						},

						"month_days": { //used with monthly,
							Type:          schema.TypeSet,
							Optional:      true,
							ConflictsWith: []string{"recurrence.0.week_days", "recurrence.0.monthly_occurrences"},
							MinItems:      1,
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								ValidateFunc: validate.IntBetweenAndNot(-31, 31, 0),
							},
							Set: set.HashInt,
						},

						"monthly_occurrences": {
							Type:          schema.TypeSet,
							Optional:      true,
							ConflictsWith: []string{"recurrence.0.week_days", "recurrence.0.month_days"},
							MinItems:      1,
							Set:           resourceAzureRMSchedulerJobMonthlyOccurrenceHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"day": {
										Type:             schema.TypeString,
										Required:         true,
										DiffSuppressFunc: suppress.CaseDifference,
										ValidateFunc:     validate.DayOfTheWeek(true),
									},

									"occurrence": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validate.IntBetweenAndNot(-5, 5, 0),
									},
								},
							},
						},
					},
				},
			},

			"start_time": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true, //defaults to now in create function
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validate.RFC3339Time, //times in the past just start immediately
			},

			"state": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(scheduler.JobStateEnabled),
					string(scheduler.JobStateDisabled),
					// JobStateFaulted & JobStateCompleted are also possible, but silly
				}, true),
			},
		},
	}
}

func resourceArmSchedulerJobActionWebSchema(propertyName string) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// we can determine the type (HTTP/HTTPS) from the url
			// but we need to make sure the url starts with http/https
			// both so we can determine the type and as azure requires it
			"url": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validate.URLIsHTTPOrHTTPS,
			},

			"method": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Get", "Put", "Post", "Delete",
				}, true),
			},

			//only valid/used when action type is put
			"body": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"headers": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			//authentication requires HTTPS
			"authentication_basic": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				ConflictsWith: []string{
					fmt.Sprintf("%s.0.authentication_certificate", propertyName),
					fmt.Sprintf("%s.0.authentication_active_directory", propertyName),
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"password": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"authentication_certificate": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ConflictsWith: []string{
					fmt.Sprintf("%s.0.authentication_basic", propertyName),
					fmt.Sprintf("%s.0.authentication_active_directory", propertyName),
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pfx": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true, //sensitive & shortens diff
							ValidateFunc: validate.NoEmptyStrings,
						},

						"password": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"thumbprint": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"expiration": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"subject_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"authentication_active_directory": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ConflictsWith: []string{
					fmt.Sprintf("%s.0.authentication_basic", propertyName),
					fmt.Sprintf("%s.0.authentication_certificate", propertyName),
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"client_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"secret": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"audience": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true, //is defaulted to the ServiceManagementEndpoint in create
						},
					},
				},
			},
		},
	}
}

func resourceArmSchedulerJobActionStorageSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{

			"storage_account_name": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validateArmStorageAccountName,
			},

			"storage_queue_name": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validateArmStorageQueueName,
			},

			"sas_token": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"message": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
		},
	}
}

func resourceArmSchedulerJobCustomizeDiff(diff *schema.ResourceDiff, _ interface{}) error {
	_, hasWeb := diff.GetOk("action_web")
	_, hasStorage := diff.GetOk("action_storage_queue")
	if !hasWeb && !hasStorage {
		return fmt.Errorf("One of `action_web` or `action_storage_queue` must be set")
	}

	if b, ok := diff.GetOk("recurrence"); ok {
		if recurrence, ok := b.([]interface{})[0].(map[string]interface{}); ok {
			//if neither count nor end time is set the API will silently fail
			_, hasCount := recurrence["count"]
			_, hasEnd := recurrence["end_time"]
			if !hasCount && !hasEnd {
				return fmt.Errorf("One of `count` or `end_time` must be set for the 'recurrence' block.")
			}
		}
	}

	return nil
}

func resourceArmSchedulerJobCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Scheduler.JobsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	jobCollection := d.Get("job_collection_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, jobCollection, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Scheduler Job %q (resource group %q)", name, resourceGroup)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_scheduler_job", *existing.ID)
		}
	}

	job := scheduler.JobDefinition{
		Properties: &scheduler.JobProperties{
			Action: expandAzureArmSchedulerJobAction(d, meta),
		},
	}

	log.Printf("[DEBUG] Creating/updating Scheduler Job %q (resource group %q)", name, resourceGroup)

	//schedule (recurrence)
	if b, ok := d.GetOk("recurrence"); ok {
		job.Properties.Recurrence = expandAzureArmSchedulerJobRecurrence(b)
	}

	//start time, should be validated by schema, also defaults to now if not set
	if v, ok := d.GetOk("start_time"); ok {
		startTime, _ := time.Parse(time.RFC3339, v.(string))
		job.Properties.StartTime = &date.Time{Time: startTime}
	} else {
		job.Properties.StartTime = &date.Time{Time: time.Now()}
	}

	//state
	if state, ok := d.GetOk("state"); ok {
		job.Properties.State = scheduler.JobState(state.(string))
	}

	resp, err := client.CreateOrUpdate(ctx, resourceGroup, jobCollection, name, job)
	if err != nil {
		return fmt.Errorf("Error creating/updating Scheduler Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmSchedulerJobRead(d, meta)
}

func resourceArmSchedulerJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Scheduler.JobsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["jobs"]
	resourceGroup := id.ResourceGroup
	jobCollection := id.Path["jobCollections"]

	log.Printf("[DEBUG] Reading Scheduler Job %q (resource group %q)", name, resourceGroup)

	job, err := client.Get(ctx, resourceGroup, jobCollection, name)
	if err != nil {
		if utils.ResponseWasNotFound(job.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Scheduler Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	//standard properties
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("job_collection_name", jobCollection)

	//check & get properties
	properties := job.Properties
	if properties != nil {
		//action
		action := properties.Action
		if action != nil {
			actionType := strings.ToLower(string(action.Type))

			if strings.EqualFold(actionType, string(scheduler.HTTP)) || strings.EqualFold(actionType, string(scheduler.HTTPS)) {
				if err := d.Set("action_web", flattenAzureArmSchedulerJobActionRequest(d, "action_web", action.Request)); err != nil {
					return err
				}
			} else if strings.EqualFold(actionType, string(scheduler.StorageQueue)) {
				if err := d.Set("action_storage_queue", flattenAzureArmSchedulerJobActionQueueMessage(d, "action_storage_queue", action.QueueMessage)); err != nil {
					return err
				}
			}

			//error action
			if errorAction := action.ErrorAction; errorAction != nil {
				errorActionType := strings.ToLower(string(errorAction.Type))

				if strings.EqualFold(errorActionType, string(scheduler.HTTP)) || strings.EqualFold(errorActionType, string(scheduler.HTTPS)) {
					if err := d.Set("error_action_web", flattenAzureArmSchedulerJobActionRequest(d, "error_action_web", errorAction.Request)); err != nil {
						return err
					}
				} else if strings.EqualFold(errorActionType, string(scheduler.StorageQueue)) {
					if err := d.Set("error_action_storage_queue", flattenAzureArmSchedulerJobActionQueueMessage(d, "error_action_storage_queue", errorAction.QueueMessage)); err != nil {
						return err
					}
				}
			}

			//retry
			if retry := action.RetryPolicy; retry != nil {
				//if its not fixed we should not have a retry block
				//api returns whatever casing it gets so do a case insensitive comparison
				if strings.EqualFold(string(retry.RetryType), string(scheduler.Fixed)) {
					if err := d.Set("retry", flattenAzureArmSchedulerJobActionRetry(retry)); err != nil {
						return err
					}
				}
			}
		}

		//schedule
		if recurrence := properties.Recurrence; recurrence != nil {
			if err := d.Set("recurrence", flattenAzureArmSchedulerJobSchedule(recurrence)); err != nil {
				return err
			}
		}

		if v := properties.StartTime; v != nil {
			d.Set("start_time", (*v).Format(time.RFC3339))
		}

		//status && state
		d.Set("state", properties.State)
	}

	return nil
}

func resourceArmSchedulerJobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Scheduler.JobsClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["jobs"]
	resourceGroup := id.ResourceGroup
	jobCollection := id.Path["jobCollections"]

	log.Printf("[DEBUG] Deleting Scheduler Job %q (resource group %q)", name, resourceGroup)

	resp, err := client.Delete(ctx, resourceGroup, jobCollection, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error issuing delete request for Scheduler Job %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandAzureArmSchedulerJobAction(d *schema.ResourceData, meta interface{}) *scheduler.JobAction {
	action := scheduler.JobAction{}

	//action
	if b, ok := d.GetOk("action_web"); ok {
		action.Request, action.Type = expandAzureArmSchedulerJobActionRequest(b, meta)
	} else if b, ok := d.GetOk("action_storage_queue"); ok {
		action.QueueMessage = expandAzureArmSchedulerJobActionStorage(b)
		action.Type = scheduler.StorageQueue
	}

	//error action
	if b, ok := d.GetOk("error_action_web"); ok {
		action.ErrorAction = &scheduler.JobErrorAction{}
		action.ErrorAction.Request, action.ErrorAction.Type = expandAzureArmSchedulerJobActionRequest(b, meta)
	} else if b, ok := d.GetOk("error_action_storage_queue"); ok {
		action.ErrorAction = &scheduler.JobErrorAction{}
		action.ErrorAction.QueueMessage = expandAzureArmSchedulerJobActionStorage(b)
		action.ErrorAction.Type = scheduler.StorageQueue
	}

	//retry policy
	if b, ok := d.GetOk("retry"); ok {
		action.RetryPolicy = expandAzureArmSchedulerJobActionRetry(b)
	} else {
		action.RetryPolicy = &scheduler.RetryPolicy{
			RetryType: scheduler.None,
		}
	}

	return &action
}

func expandAzureArmSchedulerJobActionRequest(b interface{}, meta interface{}) (*scheduler.HTTPRequest, scheduler.JobActionType) {
	block := b.([]interface{})[0].(map[string]interface{})

	url := block["url"].(string)

	request := scheduler.HTTPRequest{
		URI:     &url,
		Method:  utils.String(block["method"].(string)),
		Headers: map[string]*string{},
	}

	// determine type from the url, the property validation must ensure this
	// otherwise we need to worry about what happens if neither is true
	var jobType scheduler.JobActionType
	if strings.HasPrefix(strings.ToLower(url), "https://") {
		jobType = scheduler.HTTPS
	} else {
		jobType = scheduler.HTTP
	}

	//load headers
	for k, v := range block["headers"].(map[string]interface{}) {
		request.Headers[k] = utils.String(v.(string))
	}

	//only valid for a set
	if v, ok := block["body"].(string); ok && v != "" {
		request.Body = utils.String(block["body"].(string))
	}

	//authentications
	if v, ok := block["authentication_basic"].([]interface{}); ok && len(v) > 0 {
		b := v[0].(map[string]interface{})
		request.Authentication = &scheduler.BasicAuthentication{
			Type:     scheduler.TypeBasic,
			Username: utils.String(b["username"].(string)),
			Password: utils.String(b["password"].(string)),
		}
	}

	if v, ok := block["authentication_certificate"].([]interface{}); ok && len(v) > 0 {
		b := v[0].(map[string]interface{})
		request.Authentication = &scheduler.ClientCertAuthentication{
			Type:     scheduler.TypeClientCertificate,
			Pfx:      utils.String(b["pfx"].(string)),
			Password: utils.String(b["password"].(string)),
		}
	}

	if v, ok := block["authentication_active_directory"].([]interface{}); ok && len(v) > 0 {
		b := v[0].(map[string]interface{})
		oauth := &scheduler.OAuthAuthentication{
			Type:     scheduler.TypeActiveDirectoryOAuth,
			Tenant:   utils.String(b["tenant_id"].(string)),
			ClientID: utils.String(b["client_id"].(string)),
			Secret:   utils.String(b["secret"].(string)),
		}

		//default to the service Management Endpoint
		if v, ok := b["audience"].(string); ok {
			oauth.Audience = utils.String(v)
		} else {
			oauth.Audience = utils.String(meta.(*ArmClient).Account.Environment.ServiceManagementEndpoint)
		}

		request.Authentication = oauth
	}

	return &request, jobType
}

func expandAzureArmSchedulerJobActionStorage(b interface{}) *scheduler.StorageQueueMessage {
	block := b.([]interface{})[0].(map[string]interface{})

	message := scheduler.StorageQueueMessage{
		StorageAccount: utils.String(block["storage_account_name"].(string)),
		QueueName:      utils.String(block["storage_queue_name"].(string)),
		SasToken:       utils.String(block["sas_token"].(string)),
		Message:        utils.String(block["message"].(string)),
	}

	return &message
}

func expandAzureArmSchedulerJobActionRetry(b interface{}) *scheduler.RetryPolicy {
	block := b.([]interface{})[0].(map[string]interface{})
	retry := scheduler.RetryPolicy{
		RetryType: scheduler.Fixed,
	}

	if v, ok := block["interval"].(string); ok && v != "" {
		retry.RetryInterval = utils.String(v)
	}
	if v, ok := block["count"].(int); ok {
		retry.RetryCount = utils.Int32(int32(v))
	}

	return &retry
}

func expandAzureArmSchedulerJobRecurrence(b interface{}) *scheduler.JobRecurrence {
	block := b.([]interface{})[0].(map[string]interface{})

	recurrence := scheduler.JobRecurrence{
		Frequency: scheduler.RecurrenceFrequency(block["frequency"].(string)),
		Interval:  utils.Int32(int32(block["interval"].(int))),
	}
	if v, ok := block["count"].(int); ok {
		recurrence.Count = utils.Int32(int32(v))
	}
	if v, ok := block["end_time"].(string); ok && v != "" {
		endTime, _ := time.Parse(time.RFC3339, v) //validated by schema
		recurrence.EndTime = &date.Time{Time: endTime}
	}

	schedule := scheduler.JobRecurrenceSchedule{}
	if s, ok := block["minutes"].(*schema.Set); ok && s.Len() > 0 {
		schedule.Minutes = set.ToSliceInt32P(s)
	}
	if s, ok := block["hours"].(*schema.Set); ok && s.Len() > 0 {
		schedule.Hours = set.ToSliceInt32P(s)
	}

	if s, ok := block["week_days"].(*schema.Set); ok && s.Len() > 0 {
		var slice []scheduler.DayOfWeek
		for _, m := range s.List() {
			slice = append(slice, scheduler.DayOfWeek(m.(string)))
		}
		schedule.WeekDays = &slice
	}

	if s, ok := block["month_days"].(*schema.Set); ok && s.Len() > 0 {
		schedule.MonthDays = set.ToSliceInt32P(s)
	}
	if s, ok := block["monthly_occurrences"].(*schema.Set); ok && s.Len() > 0 {
		var slice []scheduler.JobRecurrenceScheduleMonthlyOccurrence
		for _, e := range s.List() {
			b := e.(map[string]interface{})
			slice = append(slice, scheduler.JobRecurrenceScheduleMonthlyOccurrence{
				Day:        scheduler.JobScheduleDay(b["day"].(string)),
				Occurrence: utils.Int32(int32(b["occurrence"].(int))),
			})
		}
		schedule.MonthlyOccurrences = &slice
	}

	// if non of these are set and we try and send out a empty JobRecurrenceSchedule block
	// the API will not respond so kindly
	if schedule.Minutes != nil ||
		schedule.Hours != nil ||
		schedule.WeekDays != nil ||
		schedule.MonthDays != nil ||
		schedule.MonthlyOccurrences != nil {
		recurrence.Schedule = &schedule
	}
	return &recurrence
}

// flatten (API --> terraform)

func flattenAzureArmSchedulerJobActionRequest(d *schema.ResourceData, blockName string, request *scheduler.HTTPRequest) []interface{} {
	block := map[string]interface{}{}

	if v := request.URI; v != nil {
		block["url"] = *v
	}
	if v := request.Method; v != nil {
		block["method"] = *v
	}
	if v := request.Body; v != nil {
		block["body"] = *v
	}

	if v := request.Headers; v != nil {
		headers := map[string]interface{}{}
		for k, v := range v {
			headers[k] = *v
		}

		block["headers"] = headers
	}

	if auth := request.Authentication; auth != nil {
		authBlock := map[string]interface{}{}

		if basic, ok := auth.AsBasicAuthentication(); ok {
			block["authentication_basic"] = []interface{}{authBlock}

			if v := basic.Username; v != nil {
				authBlock["username"] = *v
			}

			//password is not returned so lets fetch it
			if v, ok := d.GetOk(fmt.Sprintf("%s.0.authentication_basic.0.password", blockName)); ok {
				authBlock["password"] = v.(string)
			}
		} else if cert, ok := auth.AsClientCertAuthentication(); ok {
			block["authentication_certificate"] = []interface{}{authBlock}

			if v := cert.CertificateThumbprint; v != nil {
				authBlock["thumbprint"] = *v
			}
			if v := cert.CertificateExpirationDate; v != nil {
				authBlock["expiration"] = (*v).Format(time.RFC3339)
			}
			if v := cert.CertificateSubjectName; v != nil {
				authBlock["subject_name"] = *v
			}

			//these properties not returned, so lets grab them
			if v, ok := d.GetOk(fmt.Sprintf("%s.0.authentication_certificate.0.pfx", blockName)); ok {
				authBlock["pfx"] = v.(string)
			}
			if v, ok := d.GetOk(fmt.Sprintf("%s.0.authentication_certificate.0.password", blockName)); ok {
				authBlock["password"] = v.(string)
			}
		} else if oauth, ok := auth.AsOAuthAuthentication(); ok {
			block["authentication_active_directory"] = []interface{}{authBlock}

			if v := oauth.Audience; v != nil {
				authBlock["audience"] = *v
			}
			if v := oauth.ClientID; v != nil {
				authBlock["client_id"] = *v
			}
			if v := oauth.Tenant; v != nil {
				authBlock["tenant_id"] = *v
			}

			//secret is not returned
			if v, ok := d.GetOk(fmt.Sprintf("%s.0.authentication_active_directory.0.secret", blockName)); ok {
				authBlock["secret"] = v.(string)
			}
		}
	}

	return []interface{}{block}
}

func flattenAzureArmSchedulerJobActionQueueMessage(d *schema.ResourceData, blockName string, qm *scheduler.StorageQueueMessage) []interface{} {
	block := map[string]interface{}{}

	if v := qm.StorageAccount; v != nil {
		block["storage_account_name"] = *v
	}
	if v := qm.QueueName; v != nil {
		block["storage_queue_name"] = *v
	}
	if v := qm.Message; v != nil {
		block["message"] = *v
	}

	//sas_token is not returned by the API
	if v, ok := d.GetOk(fmt.Sprintf("%s.0.sas_token", blockName)); ok {
		block["sas_token"] = v.(string)
	}

	return []interface{}{block}
}

func flattenAzureArmSchedulerJobActionRetry(retry *scheduler.RetryPolicy) []interface{} {
	block := map[string]interface{}{}

	if v := retry.RetryInterval; v != nil {
		block["interval"] = *v
	}
	if v := retry.RetryCount; v != nil {
		block["count"] = *v
	}

	return []interface{}{block}
}

func flattenAzureArmSchedulerJobSchedule(recurrence *scheduler.JobRecurrence) []interface{} {
	block := map[string]interface{}{}

	block["frequency"] = string(recurrence.Frequency)

	if v := recurrence.Interval; v != nil {
		block["interval"] = *v
	}
	if v := recurrence.Count; v != nil {
		block["count"] = *v
	}
	if v := recurrence.EndTime; v != nil {
		block["end_time"] = (*v).Format(time.RFC3339)
	}

	if schedule := recurrence.Schedule; schedule != nil {
		if v := schedule.Minutes; v != nil {
			block["minutes"] = set.FromInt32Slice(*v)
		}
		if v := schedule.Hours; v != nil {
			block["hours"] = set.FromInt32Slice(*v)
		}

		if v := schedule.WeekDays; v != nil {
			s := &schema.Set{F: schema.HashString}
			for _, v := range *v {
				s.Add(string(v))
			}
			block["week_days"] = s
		}
		if v := schedule.MonthDays; v != nil {
			block["month_days"] = set.FromInt32Slice(*v)
		}

		if monthly := schedule.MonthlyOccurrences; monthly != nil {
			s := &schema.Set{F: resourceAzureRMSchedulerJobMonthlyOccurrenceHash}
			for _, e := range *monthly {
				m := map[string]interface{}{
					"day": string(e.Day),
				}

				if v := e.Occurrence; v != nil {
					m["occurrence"] = int(*v)
				}

				s.Add(m)
			}
			block["monthly_occurrences"] = s
		}
	}

	return []interface{}{block}
}

func resourceAzureRMSchedulerJobMonthlyOccurrenceHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		//day returned by azure is in a different case then the API constants
		buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["day"].(string))))
		buf.WriteString(fmt.Sprintf("%d-", m["occurrence"].(int)))
	}

	return hashcode.String(buf.String())
}
