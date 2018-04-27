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

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/supress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSchedulerJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSchedulerJobCreateUpdate,
		Read:   resourceArmSchedulerJobRead,
		Update: resourceArmSchedulerJobCreateUpdate,
		Delete: resourceArmSchedulerJobDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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

			"resource_group_name": resourceGroupNameSchema(),

			"job_collection_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			//actions
			"action_web": resourceArmSchedulerJobActionWebSchema("action_web"),

			//each action can also be an error action
			"error_action_web": resourceArmSchedulerJobActionWebSchema("error_action_web"),

			//retry policy
			"retry": {
				Type:     schema.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						//silently fails if the duration is not in the correct format
						//todo validation
						"interval": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "00:00:30",
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
							DiffSuppressFunc: supress.CaseDifference,
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
							DiffSuppressFunc: supress.Rfc3339Time,
							ValidateFunc:     validate.Rfc3339Time,
						},

						"minutes": {
							Type:     schema.TypeSet,
							Optional: true,
							MinItems: 1,
							Elem:     &schema.Schema{Type: schema.TypeInt},
							Set:      set.HashInt,
						},

						"hours": {
							Type:     schema.TypeSet,
							Optional: true,
							MinItems: 1,
							Elem:     &schema.Schema{Type: schema.TypeInt},
							Set:      set.HashInt,
						},

						"week_days": { //used with weekly
							Type:          schema.TypeSet,
							Optional:      true,
							ConflictsWith: []string{"recurrence.0.month_days", "recurrence.0.monthly_occurrences"},
							MinItems:      1,
							Elem:          &schema.Schema{Type: schema.TypeString},
							//the constants are title cased but the API returns all lowercase
							//so lets ignore the case
							Set: set.HashStringIgnoreCase,
						},

						"month_days": { //used with monthly, -1, 1- must be between 1/31
							Type:          schema.TypeSet,
							Optional:      true,
							ConflictsWith: []string{"recurrence.0.week_days", "recurrence.0.monthly_occurrences"},
							MinItems:      1,
							Elem:          &schema.Schema{Type: schema.TypeInt},
							Set:           set.HashInt,
						},

						"monthly_occurrences": {
							Type:          schema.TypeSet,
							Optional:      true,
							ConflictsWith: []string{"recurrence.0.week_days", "recurrence.0.month_days"},
							MinItems:      1,
							Set:           resourceAzureRMSchedulerJobMonthlyOccurrenceHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"day": { // DatOfWeek (sunday monday)
										Type:     schema.TypeString,
										Required: true,
									},
									"occurrence": { //-5 - 5, not 0
										Type:     schema.TypeInt,
										Required: true,
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
				Default:          time.Now().Format(time.RFC3339), //default to now
				DiffSuppressFunc: supress.Rfc3339Time,
				ValidateFunc:     validate.Rfc3339Time, //times in the past just start immediately
			},

			"state": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: supress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(scheduler.JobStateEnabled),
					string(scheduler.JobStateDisabled),
					// JobStateFaulted & JobStateCompleted are also possible, but silly
				}, true),
			},

			//status
			"execution_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"failure_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"faulted_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_execution_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_execution_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: resourceArmSchedulerJobCustomizeDiff,
	}
}

func resourceArmSchedulerJobActionWebSchema(propertyName string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MinItems: 1,
		MaxItems: 1,
		Optional: true,
		//ConflictsWith: conflictsWith,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{

				// we can determine the type (HTTP/HTTPS) from the url
				// but we need to make sure the url starts with http/https
				// both so we can determine the type and as azure requires it
				"url": {
					Type:             schema.TypeString,
					Optional:         true,
					DiffSuppressFunc: supress.CaseDifference,
					ValidateFunc:     validate.Url,
				},

				"method": {
					Type:     schema.TypeString,
					Optional: true,
					//DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
					Default: "Get", //todo have a default or force user to pick?
					ValidateFunc: validation.StringInSlice([]string{
						"Get", "Put", "Post", "Delete",
					}, true),
				},

				//only valid/used when action type is put
				"body": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"headers": {
					Type:     schema.TypeMap,
					Optional: true,
				},

				//authentication requires HTTPS
				"authentication_basic": {
					Type:     schema.TypeList,
					MinItems: 1,
					MaxItems: 1,
					Optional: true,
					ConflictsWith: []string{
						fmt.Sprintf("%s.0.authentication_certificate", propertyName),
						fmt.Sprintf("%s.0.authentication_active_directory", propertyName),
					},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"username": {
								Type:     schema.TypeString,
								Required: true,
							},

							"password": {
								Type:      schema.TypeString,
								Required:  true,
								Sensitive: true,
							},
						},
					},
				},

				"authentication_certificate": {
					Type:     schema.TypeList,
					Optional: true,
					MinItems: 1,
					MaxItems: 1,
					ConflictsWith: []string{
						fmt.Sprintf("%s.0.authentication_basic", propertyName),
						fmt.Sprintf("%s.0.authentication_active_directory", propertyName),
					},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"pfx": {
								Type:      schema.TypeString,
								Required:  true,
								Sensitive: true, //sensitive & shortens diff
							},

							"password": {
								Type:      schema.TypeString,
								Required:  true,
								Sensitive: true,
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
					MinItems: 1,
					MaxItems: 1,
					ConflictsWith: []string{
						fmt.Sprintf("%s.0.authentication_basic", propertyName),
						fmt.Sprintf("%s.0.authentication_certificate", propertyName),
					},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"tenant_id": {
								Type:     schema.TypeString,
								Required: true,
							},

							"client_id": {
								Type:     schema.TypeString,
								Required: true,
							},

							"secret": {
								Type:      schema.TypeString,
								Required:  true,
								Sensitive: true,
							},

							"audience": {
								Type:     schema.TypeString,
								Optional: true,
								Default:  "https://management.core.windows.net/",
							},
						},
					},
				},
			},
		},
	}
}

func resourceArmSchedulerJobCustomizeDiff(diff *schema.ResourceDiff, v interface{}) error {

	_, hasWeb := diff.GetOk("action_web")
	if !hasWeb {
		return fmt.Errorf("One of `action_web`, `action_servicebus` or `action_storagequeue` must be set")
	}

	if b, ok := diff.GetOk("recurrence"); ok {
		if recurrence, ok := b.([]interface{})[0].(map[string]interface{}); ok {

			//if neither count nor end time is set the API will silently fail
			_, hasCount := recurrence["count"]
			_, hasEnd := recurrence["end_time"]
			if !hasCount && !hasEnd {
				return fmt.Errorf("One of `count` or `end_time` must be set for the 'recurrence' block.")
			}

			if v, ok := recurrence["minutes"].(*schema.Set); ok {
				for _, e := range v.List() {
					//leverage existing function, validates type and value
					if _, errors := validation.IntBetween(0, 59)(e, "minutes"); len(errors) > 0 {
						return errors[0]
					}
				}
			}

			if v, ok := recurrence["hours"].(*schema.Set); ok {
				for _, e := range v.List() {
					//leverage existing function, validates type and value
					if _, errors := validation.IntBetween(0, 23)(e, "hours"); len(errors) > 0 {
						return errors[0]
					}
				}
			}

			if v, ok := recurrence["week_days"].(*schema.Set); ok {
				for _, e := range v.List() {
					//leverage existing function, validates type and value
					if _, errors := validation.StringInSlice([]string{
						string(scheduler.Monday),
						string(scheduler.Tuesday),
						string(scheduler.Wednesday),
						string(scheduler.Thursday),
						string(scheduler.Friday),
						string(scheduler.Saturday),
						string(scheduler.Sunday),
					}, true)(e, "week_days"); len(errors) > 0 {
						return errors[0] //string in slice can only return one
					}
				}
			}

			if v, ok := recurrence["month_days"].(*schema.Set); ok {
				for _, e := range v.List() {
					v := e.(int)
					if (-31 < v && v > 31) && v != 0 {
						return fmt.Errorf("expected 'month_days' to be in the range (-31 - 31) excluding 0, got %d", v)
					}
				}
			}
		}
	}

	return nil
}

func resourceArmSchedulerJobCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).schedulerJobsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	jobCollection := d.Get("job_collection_name").(string)

	job := scheduler.JobDefinition{
		Properties: &scheduler.JobProperties{
			Action: &scheduler.JobAction{},
		},
	}

	log.Printf("[DEBUG] Creating/updating Scheduler Job %q (resource group %q)", name, resourceGroup)

	//action
	if b, ok := d.GetOk("action_web"); ok {
		job.Properties.Action.Request, job.Properties.Action.Type = expandAzureArmSchedulerJobActionRequest(b)
	}

	//error action
	if b, ok := d.GetOk("error_action_web"); ok {
		job.Properties.Action.ErrorAction = &scheduler.JobErrorAction{}
		job.Properties.Action.ErrorAction.Request, job.Properties.Action.ErrorAction.Type = expandAzureArmSchedulerJobActionRequest(b)
	}

	//retry policy
	if b, ok := d.GetOk("retry"); ok {
		job.Properties.Action.RetryPolicy = expandAzureArmSchedulerJobActionRetry(b)
	} else {
		job.Properties.Action.RetryPolicy = &scheduler.RetryPolicy{
			RetryType: scheduler.None,
		}
	}

	//schedule (recurrence)
	if b, ok := d.GetOk("recurrence"); ok {
		job.Properties.Recurrence = expandAzureArmSchedulerJobRecurrence(b)
	}

	//start time
	startTime, err := time.Parse(time.RFC3339, d.Get("start_time").(string))
	if err != nil {
		return fmt.Errorf("Error parsing start time (%s) for job %q (Resource Group %q): %+v", d.Get("start_time"), name, resourceGroup, err)
	}
	job.Properties.StartTime = &date.Time{startTime}

	//state
	if state, ok := d.GetOk("state"); ok {
		job.Properties.State = scheduler.JobState(state.(string))
	}

	resp, err := client.CreateOrUpdate(ctx, resourceGroup, jobCollection, name, job)
	if err != nil {
		return fmt.Errorf("Error creating/updating Scheduler Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmSchedulerJobPopulate(d, resourceGroup, jobCollection, &resp)
}

func resourceArmSchedulerJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).schedulerJobsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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

	return resourceArmSchedulerJobPopulate(d, resourceGroup, jobCollection, &job)
}

func resourceArmSchedulerJobPopulate(d *schema.ResourceData, resourceGroup string, jobCollection string, job *scheduler.JobDefinition) error {

	//standard properties
	name := strings.Split(*job.Name, "/")[1] //job.Name is actually "{job_collection_name}/{job_name}
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("job_collection_name", jobCollection)

	//check & get properties
	properties := job.Properties
	if properties == nil {
		return fmt.Errorf("job properties is nil")
	}

	//action
	action := properties.Action
	if action == nil {
		return fmt.Errorf("job action is nil")
	}
	actionType := strings.ToLower(string(action.Type))
	if strings.EqualFold(actionType, string(scheduler.HTTP)) || strings.EqualFold(actionType, string(scheduler.HTTPS)) {
		d.Set("action_web", flattenAzureArmSchedulerJobActionRequest(action.Request, d.Get("action_web")))
	} else {
		return fmt.Errorf("Unknown job type %q for scheduler job %q action (Resource Group %q)", action.Type, name, resourceGroup)
	}

	//error action
	if errorAction := action.ErrorAction; errorAction != nil {
		if strings.EqualFold(actionType, string(scheduler.HTTP)) || strings.EqualFold(actionType, string(scheduler.HTTPS)) {
			d.Set("error_action_web", flattenAzureArmSchedulerJobActionRequest(errorAction.Request, d.Get("error_action_web")))
		} else {
			return fmt.Errorf("Unknown job type %q for scheduler job %q error action (Resource Group %q)", errorAction.Type, name, resourceGroup)
		}
	}

	//retry
	if retry := properties.Action.RetryPolicy; retry != nil {
		//if its not fixed we should not have a retry block
		if retry.RetryType == scheduler.Fixed {
			d.Set("retry", flattenAzureArmSchedulerJobActionRetry(retry))
		}
	}

	//schedule
	if recurrence := properties.Recurrence; recurrence != nil {
		d.Set("recurrence", flattenAzureArmSchedulerJobSchedule(recurrence))
	}

	d.Set("start_time", properties.StartTime.Format(time.RFC3339))
	d.Set("state", properties.State)

	//status
	status := properties.Status
	if status != nil {
		if v := status.ExecutionCount; v != nil {
			d.Set("execution_count", *v)
		}
		if v := status.FailureCount; v != nil {
			d.Set("failure_count", *v)
		}
		if v := status.FaultedCount; v != nil {
			d.Set("faulted_count", *v)
		}

		//these can be nil, if so set to empty so any outputs referencing them won't explode
		if v := status.LastExecutionTime; v != nil {
			d.Set("last_execution_time", (*v).Format(time.RFC3339))
		} else {
			d.Set("last_execution_time", "")
		}
		if v := status.NextExecutionTime; v != nil {
			d.Set("next_execution_time", (*v).Format(time.RFC3339))
		} else {
			d.Set("next_execution_time", "")
		}
	}

	return nil
}

func resourceArmSchedulerJobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).schedulerJobsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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

//expand (terraform -> API)
func expandAzureArmSchedulerJobActionRequest(b interface{}) (*scheduler.HTTPRequest, scheduler.JobActionType) {

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
		//} else if strings.HasPrefix(strings.ToLower(url), "http://") {
	} else {
		jobType = scheduler.HTTP
	}

	//load headers
	//if v, ok := block["headers"].(map[string]interface{}); ok { //check doesn't seem to be needed
	for k, v := range block["headers"].(map[string]interface{}) {
		(request.Headers)[k] = utils.String(v.(string))
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
		request.Authentication = &scheduler.OAuthAuthentication{
			Type:     scheduler.TypeActiveDirectoryOAuth,
			Tenant:   utils.String(b["tenant_id"].(string)),
			ClientID: utils.String(b["client_id"].(string)),
			Audience: utils.String(b["audience"].(string)),
			Secret:   utils.String(b["secret"].(string)),
		}
	}

	return &request, jobType
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
	recurrence := scheduler.JobRecurrence{}
	schedule := scheduler.JobRecurrenceSchedule{}

	if v, ok := block["frequency"].(string); ok && v != "" {
		recurrence.Frequency = scheduler.RecurrenceFrequency(v)
	}
	if v, ok := block["interval"].(int); ok {
		recurrence.Interval = utils.Int32(int32(v))
	}
	if v, ok := block["count"].(int); ok {
		recurrence.Count = utils.Int32(int32(v))
	}
	if v, ok := block["end_time"].(string); ok && v != "" {
		endTime, _ := time.Parse(time.RFC3339, v)
		recurrence.EndTime = &date.Time{Time: endTime}
	}

	if s, ok := block["minutes"].(*schema.Set); ok && s.Len() > 0 {
		schedule.Minutes = setToSliceInt32P(s)
	}
	if s, ok := block["hours"].(*schema.Set); ok && s.Len() > 0 {
		schedule.Hours = setToSliceInt32P(s)
	}

	if s, ok := block["week_days"].(*schema.Set); ok && s.Len() > 0 {
		var slice []scheduler.DayOfWeek
		for _, m := range s.List() {
			slice = append(slice, scheduler.DayOfWeek(m.(string)))
		}
		schedule.WeekDays = &slice
	}

	if s, ok := block["month_days"].(*schema.Set); ok && s.Len() > 0 {
		schedule.MonthDays = setToSliceInt32P(s)
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

func flattenAzureArmSchedulerJobActionRequest(request *scheduler.HTTPRequest, ob interface{}) []interface{} {
	oldBlock := map[string]interface{}{}

	if v, ok := ob.([]interface{}); ok && len(v) > 0 {
		oldBlock = v[0].(map[string]interface{})
	}

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

			//password is always blank, so preserve state
			if v, ok := oldBlock["authentication_basic"].([]interface{}); ok && len(v) > 0 {
				oab := v[0].(map[string]interface{})
				authBlock["password"] = oab["password"]
			}

		} else if cert, ok := auth.AsClientCertAuthentication(); ok {
			block["authentication_certificate"] = []interface{}{authBlock}

			//pfx and password are always empty, so preserve state
			if v, ok := oldBlock["authentication_certificate"].([]interface{}); ok && len(v) > 0 {
				oab := v[0].(map[string]interface{})
				authBlock["pfx"] = oab["pfx"]
				authBlock["password"] = oab["password"]
			}

			if v := cert.CertificateThumbprint; v != nil {
				authBlock["thumbprint"] = *v
			}
			if v := cert.CertificateExpirationDate; v != nil {
				authBlock["expiration"] = (*v).Format(time.RFC3339)
			}
			if v := cert.CertificateSubjectName; v != nil {
				authBlock["subject_name"] = *v
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

			//secret is always empty, so preserve state
			if v, ok := oldBlock["authentication_active_directory"].([]interface{}); ok && len(v) > 0 {
				oab := v[0].(map[string]interface{})
				authBlock["secret"] = oab["secret"]
			}
		}
	}

	return []interface{}{block}
}

func flattenAzureArmSchedulerJobActionRetry(retry *scheduler.RetryPolicy) []interface{} {
	block := map[string]interface{}{}

	block["type"] = string(retry.RetryType)
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
			block["minutes"] = sliceToSetInt32(*v)
		}
		if v := schedule.Hours; v != nil {
			block["hours"] = sliceToSetInt32(*v)
		}

		if v := schedule.WeekDays; v != nil {
			set := &schema.Set{F: schema.HashString}
			for _, v := range *v {
				set.Add(string(v))
			}
			block["week_days"] = set
		}
		if v := schedule.MonthDays; v != nil {
			block["month_days"] = sliceToSetInt32(*v)
		}

		if monthly := schedule.MonthlyOccurrences; monthly != nil {
			set := &schema.Set{F: resourceAzureRMSchedulerJobMonthlyOccurrenceHash}
			for _, e := range *monthly {

				m := map[string]interface{}{
					"day": string(e.Day),
				}

				if v := e.Occurrence; v != nil {
					m["occurrence"] = int(*v)
				}

				set.Add(m)
			}
			block["monthly_occurrences"] = set
		}
	}

	return []interface{}{block}
}

func resourceAzureRMSchedulerJobMonthlyOccurrenceHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	//day returned by azure is in a different case then the API constants
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["day"].(string))))
	buf.WriteString(fmt.Sprintf("%d-", m["occurrence"].(int)))

	return hashcode.String(buf.String())
}

func sliceToSetInt32(slice []int32) *schema.Set {
	set := &schema.Set{F: set.HashInt}
	for _, v := range slice {
		set.Add(int(v))
	}
	return set
}

func setToSliceInt32P(set *schema.Set) *[]int32 {
	var slice []int32
	for _, m := range set.List() {
		slice = append(slice, int32(m.(int)))
	}
	return &slice
}
