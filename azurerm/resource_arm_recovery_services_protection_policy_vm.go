package azurerm

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2016-06-01/backup"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

//tim4ezones:
// https://main.recoveryservices.ext.azure.com/api/timeZones
//{"resultArray":[{"displayName":"(UTC-12:00) International Date Line West","id":"Dateline Standard Time","offsetHours":-12,"offsetMinutes":0},{"displayName":"(UTC-11:00) Coordinated Universal Time-11","id":"UTC-11","offsetHours":-11,"offsetMinutes":0},{"displayName":"(UTC-10:00) Aleutian Islands","id":"Aleutian Standard Time","offsetHours":-10,"offsetMinutes":0},{"displayName":"(UTC-10:00) Hawaii","id":"Hawaiian Standard Time","offsetHours":-10,"offsetMinutes":0},{"displayName":"(UTC-09:30) Marquesas Islands","id":"Marquesas Standard Time","offsetHours":-9,"offsetMinutes":-30},{"displayName":"(UTC-09:00) Alaska","id":"Alaskan Standard Time","offsetHours":-9,"offsetMinutes":0},{"displayName":"(UTC-09:00) Coordinated Universal Time-09","id":"UTC-09","offsetHours":-9,"offsetMinutes":0},{"displayName":"(UTC-08:00) Baja California","id":"Pacific Standard Time (Mexico)","offsetHours":-8,"offsetMinutes":0},{"displayName":"(UTC-08:00) Coordinated Universal Time-08","id":"UTC-08","offsetHours":-8,"offsetMinutes":0},{"displayName":"(UTC-08:00) Pacific Time (US & Canada)","id":"Pacific Standard Time","offsetHours":-8,"offsetMinutes":0},{"displayName":"(UTC-07:00) Arizona","id":"US Mountain Standard Time","offsetHours":-7,"offsetMinutes":0},{"displayName":"(UTC-07:00) Chihuahua, La Paz, Mazatlan","id":"Mountain Standard Time (Mexico)","offsetHours":-7,"offsetMinutes":0},{"displayName":"(UTC-07:00) Mountain Time (US & Canada)","id":"Mountain Standard Time","offsetHours":-7,"offsetMinutes":0},{"displayName":"(UTC-06:00) Central America","id":"Central America Standard Time","offsetHours":-6,"offsetMinutes":0},{"displayName":"(UTC-06:00) Central Time (US & Canada)","id":"Central Standard Time","offsetHours":-6,"offsetMinutes":0},{"displayName":"(UTC-06:00) Easter Island","id":"Easter Island Standard Time","offsetHours":-6,"offsetMinutes":0},{"displayName":"(UTC-06:00) Guadalajara, Mexico City, Monterrey","id":"Central Standard Time (Mexico)","offsetHours":-6,"offsetMinutes":0},{"displayName":"(UTC-06:00) Saskatchewan","id":"Canada Central Standard Time","offsetHours":-6,"offsetMinutes":0},{"displayName":"(UTC-05:00) Bogota, Lima, Quito, Rio Branco","id":"SA Pacific Standard Time","offsetHours":-5,"offsetMinutes":0},{"displayName":"(UTC-05:00) Chetumal","id":"Eastern Standard Time (Mexico)","offsetHours":-5,"offsetMinutes":0},{"displayName":"(UTC-05:00) Eastern Time (US & Canada)","id":"Eastern Standard Time","offsetHours":-5,"offsetMinutes":0},{"displayName":"(UTC-05:00) Haiti","id":"Haiti Standard Time","offsetHours":-5,"offsetMinutes":0},{"displayName":"(UTC-05:00) Havana","id":"Cuba Standard Time","offsetHours":-5,"offsetMinutes":0},{"displayName":"(UTC-05:00) Indiana (East)","id":"US Eastern Standard Time","offsetHours":-5,"offsetMinutes":0},{"displayName":"(UTC-05:00) Turks and Caicos","id":"Turks And Caicos Standard Time","offsetHours":-5,"offsetMinutes":0},{"displayName":"(UTC-04:00) Asuncion","id":"Paraguay Standard Time","offsetHours":-4,"offsetMinutes":0},{"displayName":"(UTC-04:00) Atlantic Time (Canada)","id":"Atlantic Standard Time","offsetHours":-4,"offsetMinutes":0},{"displayName":"(UTC-04:00) Caracas","id":"Venezuela Standard Time","offsetHours":-4,"offsetMinutes":0},{"displayName":"(UTC-04:00) Cuiaba","id":"Central Brazilian Standard Time","offsetHours":-4,"offsetMinutes":0},{"displayName":"(UTC-04:00) Georgetown, La Paz, Manaus, San Juan","id":"SA Western Standard Time","offsetHours":-4,"offsetMinutes":0},{"displayName":"(UTC-04:00) Santiago","id":"Pacific SA Standard Time","offsetHours":-4,"offsetMinutes":0},{"displayName":"(UTC-03:30) Newfoundland","id":"Newfoundland Standard Time","offsetHours":-3,"offsetMinutes":-30},{"displayName":"(UTC-03:00) Araguaina","id":"Tocantins Standard Time","offsetHours":-3,"offsetMinutes":0},{"displayName":"(UTC-03:00) Brasilia","id":"E. South America Standard Time","offsetHours":-3,"offsetMinutes":0},{"displayName":"(UTC-03:00) Cayenne, Fortaleza","id":"SA Eastern Standard Time","offsetHours":-3,"offsetMinutes":0},{"displayName":"(UTC-03:00) City of Buenos Aires","id":"Argentina Standard Time","offsetHours":-3,"offsetMinutes":0},{"displayName":"(UTC-03:00) Greenland","id":"Greenland Standard Time","offsetHours":-3,"offsetMinutes":0},{"displayName":"(UTC-03:00) Montevideo","id":"Montevideo Standard Time","offsetHours":-3,"offsetMinutes":0},{"displayName":"(UTC-03:00) Punta Arenas","id":"Magallanes Standard Time","offsetHours":-3,"offsetMinutes":0},{"displayName":"(UTC-03:00) Saint Pierre and Miquelon","id":"Saint Pierre Standard Time","offsetHours":-3,"offsetMinutes":0},{"displayName":"(UTC-03:00) Salvador","id":"Bahia Standard Time","offsetHours":-3,"offsetMinutes":0},{"displayName":"(UTC-02:00) Coordinated Universal Time-02","id":"UTC-02","offsetHours":-2,"offsetMinutes":0},{"displayName":"(UTC-02:00) Mid-Atlantic - Old","id":"Mid-Atlantic Standard Time","offsetHours":-2,"offsetMinutes":0},{"displayName":"(UTC-01:00) Azores","id":"Azores Standard Time","offsetHours":-1,"offsetMinutes":0},{"displayName":"(UTC-01:00) Cabo Verde Is.","id":"Cape Verde Standard Time","offsetHours":-1,"offsetMinutes":0},{"displayName":"(UTC) Coordinated Universal Time","id":"UTC","offsetHours":0,"offsetMinutes":0},{"displayName":"(UTC+00:00) Casablanca","id":"Morocco Standard Time","offsetHours":0,"offsetMinutes":0},{"displayName":"(UTC+00:00) Dublin, Edinburgh, Lisbon, London","id":"GMT Standard Time","offsetHours":0,"offsetMinutes":0},{"displayName":"(UTC+00:00) Monrovia, Reykjavik","id":"Greenwich Standard Time","offsetHours":0,"offsetMinutes":0},{"displayName":"(UTC+01:00) Amsterdam, Berlin, Bern, Rome, Stockholm, Vienna","id":"W. Europe Standard Time","offsetHours":1,"offsetMinutes":0},{"displayName":"(UTC+01:00) Belgrade, Bratislava, Budapest, Ljubljana, Prague","id":"Central Europe Standard Time","offsetHours":1,"offsetMinutes":0},{"displayName":"(UTC+01:00) Brussels, Copenhagen, Madrid, Paris","id":"Romance Standard Time","offsetHours":1,"offsetMinutes":0},{"displayName":"(UTC+01:00) Sao Tome","id":"Sao Tome Standard Time","offsetHours":1,"offsetMinutes":0},{"displayName":"(UTC+01:00) Sarajevo, Skopje, Warsaw, Zagreb","id":"Central European Standard Time","offsetHours":1,"offsetMinutes":0},{"displayName":"(UTC+01:00) West Central Africa","id":"W. Central Africa Standard Time","offsetHours":1,"offsetMinutes":0},{"displayName":"(UTC+02:00) Amman","id":"Jordan Standard Time","offsetHours":2,"offsetMinutes":0},{"displayName":"(UTC+02:00) Athens, Bucharest","id":"GTB Standard Time","offsetHours":2,"offsetMinutes":0},{"displayName":"(UTC+02:00) Beirut","id":"Middle East Standard Time","offsetHours":2,"offsetMinutes":0},{"displayName":"(UTC+02:00) Cairo","id":"Egypt Standard Time","offsetHours":2,"offsetMinutes":0},{"displayName":"(UTC+02:00) Chisinau","id":"E. Europe Standard Time","offsetHours":2,"offsetMinutes":0},{"displayName":"(UTC+02:00) Damascus","id":"Syria Standard Time","offsetHours":2,"offsetMinutes":0},{"displayName":"(UTC+02:00) Gaza, Hebron","id":"West Bank Standard Time","offsetHours":2,"offsetMinutes":0},{"displayName":"(UTC+02:00) Harare, Pretoria","id":"South Africa Standard Time","offsetHours":2,"offsetMinutes":0},{"displayName":"(UTC+02:00) Helsinki, Kyiv, Riga, Sofia, Tallinn, Vilnius","id":"FLE Standard Time","offsetHours":2,"offsetMinutes":0},{"displayName":"(UTC+02:00) Jerusalem","id":"Israel Standard Time","offsetHours":2,"offsetMinutes":0},{"displayName":"(UTC+02:00) Kaliningrad","id":"Kaliningrad Standard Time","offsetHours":2,"offsetMinutes":0},{"displayName":"(UTC+02:00) Khartoum","id":"Sudan Standard Time","offsetHours":2,"offsetMinutes":0},{"displayName":"(UTC+02:00) Tripoli","id":"Libya Standard Time","offsetHours":2,"offsetMinutes":0},{"displayName":"(UTC+02:00) Windhoek","id":"Namibia Standard Time","offsetHours":2,"offsetMinutes":0},{"displayName":"(UTC+03:00) Baghdad","id":"Arabic Standard Time","offsetHours":3,"offsetMinutes":0},{"displayName":"(UTC+03:00) Istanbul","id":"Turkey Standard Time","offsetHours":3,"offsetMinutes":0},{"displayName":"(UTC+03:00) Kuwait, Riyadh","id":"Arab Standard Time","offsetHours":3,"offsetMinutes":0},{"displayName":"(UTC+03:00) Minsk","id":"Belarus Standard Time","offsetHours":3,"offsetMinutes":0},{"displayName":"(UTC+03:00) Moscow, St. Petersburg, Volgograd","id":"Russian Standard Time","offsetHours":3,"offsetMinutes":0},{"displayName":"(UTC+03:00) Nairobi","id":"E. Africa Standard Time","offsetHours":3,"offsetMinutes":0},{"displayName":"(UTC+03:30) Tehran","id":"Iran Standard Time","offsetHours":3,"offsetMinutes":30},{"displayName":"(UTC+04:00) Abu Dhabi, Muscat","id":"Arabian Standard Time","offsetHours":4,"offsetMinutes":0},{"displayName":"(UTC+04:00) Astrakhan, Ulyanovsk","id":"Astrakhan Standard Time","offsetHours":4,"offsetMinutes":0},{"displayName":"(UTC+04:00) Baku","id":"Azerbaijan Standard Time","offsetHours":4,"offsetMinutes":0},{"displayName":"(UTC+04:00) Izhevsk, Samara","id":"Russia Time Zone 3","offsetHours":4,"offsetMinutes":0},{"displayName":"(UTC+04:00) Port Louis","id":"Mauritius Standard Time","offsetHours":4,"offsetMinutes":0},{"displayName":"(UTC+04:00) Saratov","id":"Saratov Standard Time","offsetHours":4,"offsetMinutes":0},{"displayName":"(UTC+04:00) Tbilisi","id":"Georgian Standard Time","offsetHours":4,"offsetMinutes":0},{"displayName":"(UTC+04:00) Yerevan","id":"Caucasus Standard Time","offsetHours":4,"offsetMinutes":0},{"displayName":"(UTC+04:30) Kabul","id":"Afghanistan Standard Time","offsetHours":4,"offsetMinutes":30},{"displayName":"(UTC+05:00) Ashgabat, Tashkent","id":"West Asia Standard Time","offsetHours":5,"offsetMinutes":0},{"displayName":"(UTC+05:00) Ekaterinburg","id":"Ekaterinburg Standard Time","offsetHours":5,"offsetMinutes":0},{"displayName":"(UTC+05:00) Islamabad, Karachi","id":"Pakistan Standard Time","offsetHours":5,"offsetMinutes":0},{"displayName":"(UTC+05:30) Chennai, Kolkata, Mumbai, New Delhi","id":"India Standard Time","offsetHours":5,"offsetMinutes":30},{"displayName":"(UTC+05:30) Sri Jayawardenepura","id":"Sri Lanka Standard Time","offsetHours":5,"offsetMinutes":30},{"displayName":"(UTC+05:45) Kathmandu","id":"Nepal Standard Time","offsetHours":5,"offsetMinutes":45},{"displayName":"(UTC+06:00) Astana","id":"Central Asia Standard Time","offsetHours":6,"offsetMinutes":0},{"displayName":"(UTC+06:00) Dhaka","id":"Bangladesh Standard Time","offsetHours":6,"offsetMinutes":0},{"displayName":"(UTC+06:00) Omsk","id":"Omsk Standard Time","offsetHours":6,"offsetMinutes":0},{"displayName":"(UTC+06:30) Yangon (Rangoon)","id":"Myanmar Standard Time","offsetHours":6,"offsetMinutes":30},{"displayName":"(UTC+07:00) Bangkok, Hanoi, Jakarta","id":"SE Asia Standard Time","offsetHours":7,"offsetMinutes":0},{"displayName":"(UTC+07:00) Barnaul, Gorno-Altaysk","id":"Altai Standard Time","offsetHours":7,"offsetMinutes":0},{"displayName":"(UTC+07:00) Hovd","id":"W. Mongolia Standard Time","offsetHours":7,"offsetMinutes":0},{"displayName":"(UTC+07:00) Krasnoyarsk","id":"North Asia Standard Time","offsetHours":7,"offsetMinutes":0},{"displayName":"(UTC+07:00) Novosibirsk","id":"N. Central Asia Standard Time","offsetHours":7,"offsetMinutes":0},{"displayName":"(UTC+07:00) Tomsk","id":"Tomsk Standard Time","offsetHours":7,"offsetMinutes":0},{"displayName":"(UTC+08:00) Beijing, Chongqing, Hong Kong, Urumqi","id":"China Standard Time","offsetHours":8,"offsetMinutes":0},{"displayName":"(UTC+08:00) Irkutsk","id":"North Asia East Standard Time","offsetHours":8,"offsetMinutes":0},{"displayName":"(UTC+08:00) Kuala Lumpur, Singapore","id":"Singapore Standard Time","offsetHours":8,"offsetMinutes":0},{"displayName":"(UTC+08:00) Perth","id":"W. Australia Standard Time","offsetHours":8,"offsetMinutes":0},{"displayName":"(UTC+08:00) Taipei","id":"Taipei Standard Time","offsetHours":8,"offsetMinutes":0},{"displayName":"(UTC+08:00) Ulaanbaatar","id":"Ulaanbaatar Standard Time","offsetHours":8,"offsetMinutes":0},{"displayName":"(UTC+08:30) Pyongyang","id":"North Korea Standard Time","offsetHours":8,"offsetMinutes":30},{"displayName":"(UTC+08:45) Eucla","id":"Aus Central W. Standard Time","offsetHours":8,"offsetMinutes":45},{"displayName":"(UTC+09:00) Chita","id":"Transbaikal Standard Time","offsetHours":9,"offsetMinutes":0},{"displayName":"(UTC+09:00) Osaka, Sapporo, Tokyo","id":"Tokyo Standard Time","offsetHours":9,"offsetMinutes":0},{"displayName":"(UTC+09:00) Seoul","id":"Korea Standard Time","offsetHours":9,"offsetMinutes":0},{"displayName":"(UTC+09:00) Yakutsk","id":"Yakutsk Standard Time","offsetHours":9,"offsetMinutes":0},{"displayName":"(UTC+09:30) Adelaide","id":"Cen. Australia Standard Time","offsetHours":9,"offsetMinutes":30},{"displayName":"(UTC+09:30) Darwin","id":"AUS Central Standard Time","offsetHours":9,"offsetMinutes":30},{"displayName":"(UTC+10:00) Brisbane","id":"E. Australia Standard Time","offsetHours":10,"offsetMinutes":0},{"displayName":"(UTC+10:00) Canberra, Melbourne, Sydney","id":"AUS Eastern Standard Time","offsetHours":10,"offsetMinutes":0},{"displayName":"(UTC+10:00) Guam, Port Moresby","id":"West Pacific Standard Time","offsetHours":10,"offsetMinutes":0},{"displayName":"(UTC+10:00) Hobart","id":"Tasmania Standard Time","offsetHours":10,"offsetMinutes":0},{"displayName":"(UTC+10:00) Vladivostok","id":"Vladivostok Standard Time","offsetHours":10,"offsetMinutes":0},{"displayName":"(UTC+10:30) Lord Howe Island","id":"Lord Howe Standard Time","offsetHours":10,"offsetMinutes":30},{"displayName":"(UTC+11:00) Bougainville Island","id":"Bougainville Standard Time","offsetHours":11,"offsetMinutes":0},{"displayName":"(UTC+11:00) Chokurdakh","id":"Russia Time Zone 10","offsetHours":11,"offsetMinutes":0},{"displayName":"(UTC+11:00) Magadan","id":"Magadan Standard Time","offsetHours":11,"offsetMinutes":0},{"displayName":"(UTC+11:00) Norfolk Island","id":"Norfolk Standard Time","offsetHours":11,"offsetMinutes":0},{"displayName":"(UTC+11:00) Sakhalin","id":"Sakhalin Standard Time","offsetHours":11,"offsetMinutes":0},{"displayName":"(UTC+11:00) Solomon Is., New Caledonia","id":"Central Pacific Standard Time","offsetHours":11,"offsetMinutes":0},{"displayName":"(UTC+12:00) Anadyr, Petropavlovsk-Kamchatsky","id":"Russia Time Zone 11","offsetHours":12,"offsetMinutes":0},{"displayName":"(UTC+12:00) Auckland, Wellington","id":"New Zealand Standard Time","offsetHours":12,"offsetMinutes":0},{"displayName":"(UTC+12:00) Coordinated Universal Time+12","id":"UTC+12","offsetHours":12,"offsetMinutes":0},{"displayName":"(UTC+12:00) Fiji","id":"Fiji Standard Time","offsetHours":12,"offsetMinutes":0},{"displayName":"(UTC+12:00) Petropavlovsk-Kamchatsky - Old","id":"Kamchatka Standard Time","offsetHours":12,"offsetMinutes":0},{"displayName":"(UTC+12:45) Chatham Islands","id":"Chatham Islands Standard Time","offsetHours":12,"offsetMinutes":45},{"displayName":"(UTC+13:00) Coordinated Universal Time+13","id":"UTC+13","offsetHours":13,"offsetMinutes":0},{"displayName":"(UTC+13:00) Nuku'alofa","id":"Tonga Standard Time","offsetHours":13,"offsetMinutes":0},{"displayName":"(UTC+13:00) Samoa","id":"Samoa Standard Time","offsetHours":13,"offsetMinutes":0},{"displayName":"(UTC+14:00) Kiritimati Island","id":"Line Islands Standard Time","offsetHours":14,"offsetMinutes":0}]}
func resourceArmRecoveryServicesProtectionPolicyVm() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRecoveryServicesProtectionPolicyVmCreateUpdate,
		Read:   resourceArmRecoveryServicesProtectionPolicyVmRead,
		Update: resourceArmRecoveryServicesProtectionPolicyVmCreateUpdate,
		Delete: resourceArmRecoveryServicesProtectionPolicyVmDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z][-_!a-zA-Z0-9]{2,149}$"),
					"Backup Policy name must be 3 - 150 characters long, start with a letter, contain only letters and numbers.",
				),
			},

			"resource_group_name": resourceGroupNameSchema(),

			"recovery_vault_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{1,49}$"),
					"Recovery Service Vault name must be 2 - 50 characters long, start with a letter, contain only letters, numbers and hyphens.",
				),
			},

			"backup": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"frequency": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(backup.ScheduleRunTypeDaily),
								string(backup.ScheduleRunTypeWeekly),
							}, true),
						},

						"time": { //applies to all backup schedules & retention times (they all must be the same)
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile("^([01][0-9]|[2][0-3]):([03][0])$"), //time must be on the hour or half past
								"Time of day must match the format HH:mm where HH is 00-23 and mm is 00 or 30",
							),
						},

						"weekdays": { //only for weekly
							Type:     schema.TypeSet,
							Optional: true,
							Set:      set.HashStringIgnoreCase,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								DiffSuppressFunc: suppress.CaseDifference,
								ValidateFunc:     validate.DayOfTheWeek(true),
							},
						},
					},
				},
			},

			"retention_daily": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 9999),
						},
					},
				},
			},

			"retention_weekly": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 9999),
						},

						"weekdays": {
							Type:     schema.TypeSet,
							Required: true,
							Set:      set.HashStringIgnoreCase,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								DiffSuppressFunc: suppress.CaseDifference,
								ValidateFunc:     validate.DayOfTheWeek(true),
							},
						},
					},
				},
			},

			"retention_monthly": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 9999),
						},

						"weeks": {
							Type:     schema.TypeSet,
							Required: true,
							Set:      set.HashStringIgnoreCase,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								DiffSuppressFunc: suppress.CaseDifference,
								ValidateFunc: validation.StringInSlice([]string{
									string(backup.First),
									string(backup.Second),
									string(backup.Third),
									string(backup.Fourth),
									string(backup.Last),
								}, true),
							},
						},

						"weekdays": {
							Type:     schema.TypeSet,
							Required: true,
							Set:      set.HashStringIgnoreCase,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								DiffSuppressFunc: suppress.CaseDifference,
								ValidateFunc:     validate.DayOfTheWeek(true),
							},
						},
					},
				},
			},

			"retention_yearly": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 9999),
						},

						"months": {
							Type:     schema.TypeSet,
							Required: true,
							Set:      set.HashStringIgnoreCase,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								DiffSuppressFunc: suppress.CaseDifference,
								ValidateFunc:     validate.Month(true),
							},
						},

						"weeks": {
							Type:     schema.TypeSet,
							Required: true,
							Set:      set.HashStringIgnoreCase,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								DiffSuppressFunc: suppress.CaseDifference,
								ValidateFunc: validation.StringInSlice([]string{
									string(backup.First),
									string(backup.Second),
									string(backup.Third),
									string(backup.Fourth),
									string(backup.Last),
								}, true),
							},
						},

						"weekdays": {
							Type:     schema.TypeSet,
							Required: true,
							Set:      set.HashStringIgnoreCase,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								DiffSuppressFunc: suppress.CaseDifference,
								ValidateFunc:     validate.DayOfTheWeek(true),
							},
						},
					},
				},
			},

			"tags": tagsSchema(),
		},

		//if daily, we need daily retention
		//if weekly daily cannot be set, and we need weekly
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {

			_, hasDaily := diff.GetOk("retention_daily")
			_, hasWeekly := diff.GetOk("retention_weekly")

			frequencyI, _ := diff.GetOk("backup.0.frequency")
			frequency := strings.ToLower(frequencyI.(string))
			if frequency == "daily" {
				if !hasDaily {
					return fmt.Errorf("`retention_daily` must be set when backup.0.frequency is daily")
				}

				if _, ok := diff.GetOk("backup.0.weekdays"); ok {
					return fmt.Errorf("`backup.0.weekdays` should be not set when backup.0.frequency is daily")
				}
			} else if frequency == "weekly" {
				if hasDaily {
					return fmt.Errorf("`retention_daily` must be not set when backup.0.frequency is weekly")
				}
				if !hasWeekly {
					return fmt.Errorf("`retention_weekly` must be set when backup.0.frequency is weekly")
				}
			} else {
				return fmt.Errorf("Unrecognized value for backup.0.frequency")
			}

			return nil
		},
	}
}

func resourceArmRecoveryServicesProtectionPolicyVmCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).recoveryServicesProtectionPoliciesClient
	ctx := meta.(*ArmClient).StopContext

	policyName := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	log.Printf("[DEBUG] Creating/updating Recovery Service Protection Policy %s (resource group %q)", policyName, resourceGroup)

	//getting this ready now because its shared between *everything*, time is... complicated for this resource
	timeOfDay := d.Get("backup.0.time").(string)
	dateOfDay, err := time.Parse(time.RFC3339, fmt.Sprintf("2018-07-30T%s:00Z", timeOfDay))
	if err != nil {
		return fmt.Errorf("Error generating time from %q for policy %q (Resource Group %q): %+v", timeOfDay, policyName, resourceGroup, err)
	}
	times := append(make([]date.Time, 0), date.Time{Time: dateOfDay})

	policy := backup.ProtectionPolicyResource{
		Tags: expandTags(tags),
		Properties: &backup.AzureIaaSVMProtectionPolicy{
			BackupManagementType: backup.BackupManagementTypeAzureIaasVM,
			SchedulePolicy:       expandArmRecoveryServicesProtectionPolicySchedule(d, times),
			RetentionPolicy: &backup.LongTermRetentionPolicy{ //SimpleRetentionPolicy only has duration property ¯\_(ツ)_/¯
				RetentionPolicyType: backup.RetentionPolicyTypeLongTermRetentionPolicy,
				DailySchedule:       expandArmRecoveryServicesProtectionPolicyRetentionDaily(d, times),
				WeeklySchedule:      expandArmRecoveryServicesProtectionPolicyRetentionWeekly(d, times),
				MonthlySchedule:     expandArmRecoveryServicesProtectionPolicyRetentionMonthly(d, times),
				YearlySchedule:      expandArmRecoveryServicesProtectionPolicyRetentionYearly(d, times),
			},
		},
	}
	if _, err := client.CreateOrUpdate(ctx, vaultName, resourceGroup, policyName, policy); err != nil {
		return fmt.Errorf("Error creating/updating Recovery Service Protection Policy %q (Resource Group %q): %+v", policyName, resourceGroup, err)
	}

	resp, err := resourceArmRecoveryServicesProtectionPolicyWaitForState(client, ctx, true, vaultName, resourceGroup, policyName)
	if err != nil {
		return err
	}

	id := strings.Replace(*resp.ID, "Subscriptions", "subscriptions", 1)
	d.SetId(id)

	return resourceArmRecoveryServicesProtectionPolicyVmRead(d, meta)
}

func resourceArmRecoveryServicesProtectionPolicyVmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).recoveryServicesProtectionPoliciesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	policyName := id.Path["backupPolicies"]
	vaultName := id.Path["vaults"]
	resourceGroup := id.ResourceGroup

	log.Printf("[DEBUG] Reading Recovery Service Protection Policy %q (resource group %q)", policyName, resourceGroup)

	resp, err := client.Get(ctx, vaultName, resourceGroup, policyName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Recovery Service Protection Policy %q (Resource Group %q): %+v", policyName, resourceGroup, err)
	}

	d.Set("name", policyName)
	d.Set("resource_group_name", resourceGroup)
	d.Set("recovery_vault_name", vaultName)

	if properties, ok := resp.Properties.AsAzureIaaSVMProtectionPolicy(); ok && properties != nil {
		if schedule, ok := properties.SchedulePolicy.AsSimpleSchedulePolicy(); ok && schedule != nil {

			if err := d.Set("backup", flattenArmRecoveryServicesProtectionPolicySchedule(schedule)); err != nil {
				return fmt.Errorf("Error setting `backup`: %+v", err)
			}
		}

		if retention, ok := properties.RetentionPolicy.AsLongTermRetentionPolicy(); ok && retention != nil {
			if s := retention.DailySchedule; s != nil {
				if err := d.Set("retention_daily", flattenArmRecoveryServicesProtectionPolicyRetentionDaily(s)); err != nil {
					return fmt.Errorf("Error setting `retention_daily`: %+v", err)
				}
			} else {
				d.Set("retention_daily", nil)
			}

			if s := retention.WeeklySchedule; s != nil {
				if err := d.Set("retention_weekly", flattenArmRecoveryServicesProtectionPolicyRetentionWeekly(s)); err != nil {
					return fmt.Errorf("Error setting `retention_weekly`: %+v", err)
				}
			} else {
				d.Set("retention_weekly", nil)
			}

			if s := retention.MonthlySchedule; s != nil {
				if err := d.Set("retention_monthly", flattenArmRecoveryServicesProtectionPolicyRetentionMonthly(s)); err != nil {
					return fmt.Errorf("Error setting `retention_monthly`: %+v", err)
				}
			} else {
				d.Set("retention_monthly", nil)
			}

			if s := retention.YearlySchedule; s != nil {
				if err := d.Set("retention_yearly", flattenArmRecoveryServicesProtectionPolicyRetentionYearly(s)); err != nil {
					return fmt.Errorf("Error setting `retention_yearly`: %+v", err)
				}
			} else {
				d.Set("retention_yearly", nil)
			}

		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmRecoveryServicesProtectionPolicyVmDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).recoveryServicesProtectionPoliciesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	policyName := id.Path["backupPolicies"]
	resourceGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]

	log.Printf("[DEBUG] Deleting Recovery Service Protected Item %q (resource group %q)", policyName, resourceGroup)

	resp, err := client.Delete(ctx, vaultName, resourceGroup, policyName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error issuing delete request for Recovery Service Protection Policy %q (Resource Group %q): %+v", policyName, resourceGroup, err)
		}
	}

	if _, err := resourceArmRecoveryServicesProtectionPolicyWaitForState(client, ctx, false, vaultName, resourceGroup, policyName); err != nil {
		return err
	}

	return nil
}

func expandArmRecoveryServicesProtectionPolicySchedule(d *schema.ResourceData, times []date.Time) *backup.SimpleSchedulePolicy {
	if bb, ok := d.Get("backup").([]interface{}); ok && len(bb) > 0 {
		block := bb[0].(map[string]interface{})

		schedule := backup.SimpleSchedulePolicy{ //LongTermSchedulePolicy has no properties
			SchedulePolicyType: backup.SchedulePolicyTypeSimpleSchedulePolicy,
			ScheduleRunTimes:   &times,
		}

		if v, ok := block["frequency"].(string); ok {
			schedule.ScheduleRunFrequency = backup.ScheduleRunType(v)
		}

		if v, ok := block["weekdays"].(*schema.Set); ok {
			days := make([]backup.DayOfWeek, 0)
			for _, day := range v.List() {
				days = append(days, backup.DayOfWeek(day.(string)))
			}
			schedule.ScheduleRunDays = &days
		}

		return &schedule
	}

	return nil
}

func expandArmRecoveryServicesProtectionPolicyRetentionDaily(d *schema.ResourceData, times []date.Time) *backup.DailyRetentionSchedule {
	if rb, ok := d.Get("retention_daily").([]interface{}); ok && len(rb) > 0 {
		block := rb[0].(map[string]interface{})

		return &backup.DailyRetentionSchedule{
			RetentionTimes: &times,
			RetentionDuration: &backup.RetentionDuration{
				Count:        utils.Int32(int32(block["count"].(int))),
				DurationType: backup.RetentionDurationType(backup.RetentionDurationTypeDays),
			},
		}
	}

	return nil
}

func expandArmRecoveryServicesProtectionPolicyRetentionWeekly(d *schema.ResourceData, times []date.Time) *backup.WeeklyRetentionSchedule {
	if rb, ok := d.Get("retention_weekly").([]interface{}); ok && len(rb) > 0 {
		block := rb[0].(map[string]interface{})

		retention := backup.WeeklyRetentionSchedule{
			RetentionTimes: &times,
			RetentionDuration: &backup.RetentionDuration{
				Count:        utils.Int32(int32(block["count"].(int))),
				DurationType: backup.RetentionDurationType(backup.RetentionDurationTypeWeeks),
			},
		}

		if v, ok := block["weekdays"].(*schema.Set); ok {
			days := make([]backup.DayOfWeek, 0)
			for _, day := range v.List() {
				days = append(days, backup.DayOfWeek(day.(string)))
			}
			retention.DaysOfTheWeek = &days
		}

		return &retention
	}

	return nil
}

func expandArmRecoveryServicesProtectionPolicyRetentionMonthly(d *schema.ResourceData, times []date.Time) *backup.MonthlyRetentionSchedule {
	if rb, ok := d.Get("retention_monthly").([]interface{}); ok && len(rb) > 0 {
		block := rb[0].(map[string]interface{})

		retention := backup.MonthlyRetentionSchedule{
			RetentionScheduleFormatType: backup.RetentionScheduleFormatWeekly, //this is always weekly ¯\_(ツ)_/¯
			RetentionScheduleDaily:      nil,                                  //and this is always nil..
			RetentionScheduleWeekly:     expandArmRecoveryServicesProtectionPolicyRetentionWeeklyFormat(block),
			RetentionTimes:              &times,
			RetentionDuration: &backup.RetentionDuration{
				Count:        utils.Int32(int32(block["count"].(int))),
				DurationType: backup.RetentionDurationType(backup.RetentionDurationTypeMonths),
			},
		}

		return &retention
	}

	return nil
}

func expandArmRecoveryServicesProtectionPolicyRetentionYearly(d *schema.ResourceData, times []date.Time) *backup.YearlyRetentionSchedule {
	if rb, ok := d.Get("retention_yearly").([]interface{}); ok && len(rb) > 0 {
		block := rb[0].(map[string]interface{})

		retention := backup.YearlyRetentionSchedule{
			RetentionScheduleFormatType: backup.RetentionScheduleFormatWeekly, //this is always weekly ¯\_(ツ)_/¯
			RetentionScheduleDaily:      nil,                                  //and this is always nil..
			RetentionScheduleWeekly:     expandArmRecoveryServicesProtectionPolicyRetentionWeeklyFormat(block),
			RetentionTimes:              &times,
			RetentionDuration: &backup.RetentionDuration{
				Count:        utils.Int32(int32(block["count"].(int))),
				DurationType: backup.RetentionDurationType(backup.RetentionDurationTypeYears),
			},
		}

		if v, ok := block["months"].(*schema.Set); ok {
			months := make([]backup.MonthOfYear, 0)
			for _, month := range v.List() {
				months = append(months, backup.MonthOfYear(month.(string)))
			}
			retention.MonthsOfYear = &months
		}

		return &retention
	}

	return nil
}

func expandArmRecoveryServicesProtectionPolicyRetentionWeeklyFormat(block map[string]interface{}) *backup.WeeklyRetentionFormat {
	weekly := backup.WeeklyRetentionFormat{}

	if v, ok := block["weekdays"].(*schema.Set); ok {
		days := make([]backup.DayOfWeek, 0)
		for _, day := range v.List() {
			days = append(days, backup.DayOfWeek(day.(string)))
		}
		weekly.DaysOfTheWeek = &days
	}

	if v, ok := block["weeks"].(*schema.Set); ok {
		weeks := make([]backup.WeekOfMonth, 0)
		for _, week := range v.List() {
			weeks = append(weeks, backup.WeekOfMonth(week.(string)))
		}
		weekly.WeeksOfTheMonth = &weeks
	}

	return &weekly
}

func flattenArmRecoveryServicesProtectionPolicySchedule(schedule *backup.SimpleSchedulePolicy) []interface{} {
	block := map[string]interface{}{}

	block["frequency"] = string(schedule.ScheduleRunFrequency)

	if times := schedule.ScheduleRunTimes; times != nil && len(*times) > 0 {
		time := (*times)[0]
		block["time"] = time.Format("15:04")
	}

	if days := schedule.ScheduleRunDays; days != nil {
		weekdays := make([]interface{}, 0)
		for _, d := range *days {
			weekdays = append(weekdays, string(d))
		}
		block["weekdays"] = schema.NewSet(schema.HashString, weekdays)
	}

	return []interface{}{block}
}

func flattenArmRecoveryServicesProtectionPolicyRetentionDaily(daily *backup.DailyRetentionSchedule) []interface{} {
	block := map[string]interface{}{}

	if duration := daily.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	return []interface{}{block}
}

func flattenArmRecoveryServicesProtectionPolicyRetentionWeekly(weekly *backup.WeeklyRetentionSchedule) []interface{} {
	block := map[string]interface{}{}

	if duration := weekly.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	if days := weekly.DaysOfTheWeek; days != nil {
		weekdays := make([]interface{}, 0)
		for _, d := range *days {
			weekdays = append(weekdays, string(d))
		}
		block["weekdays"] = schema.NewSet(schema.HashString, weekdays)
	}

	return []interface{}{block}
}

func flattenArmRecoveryServicesProtectionPolicyRetentionMonthly(monthly *backup.MonthlyRetentionSchedule) []interface{} {
	block := map[string]interface{}{}

	if duration := monthly.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	if weekly := monthly.RetentionScheduleWeekly; weekly != nil {
		block["weekdays"], block["weeks"] = flattenArmRecoveryServicesProtectionPolicyRetentionWeeklyFormat(weekly)
	}

	return []interface{}{block}
}

func flattenArmRecoveryServicesProtectionPolicyRetentionYearly(yearly *backup.YearlyRetentionSchedule) []interface{} {
	block := map[string]interface{}{}

	if duration := yearly.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	if weekly := yearly.RetentionScheduleWeekly; weekly != nil {
		block["weekdays"], block["weeks"] = flattenArmRecoveryServicesProtectionPolicyRetentionWeeklyFormat(weekly)
	}

	if months := yearly.MonthsOfYear; months != nil {
		slice := make([]interface{}, 0)
		for _, d := range *months {
			slice = append(slice, string(d))
		}
		block["months"] = schema.NewSet(schema.HashString, slice)
	}

	return []interface{}{block}
}

func flattenArmRecoveryServicesProtectionPolicyRetentionWeeklyFormat(retention *backup.WeeklyRetentionFormat) (weekdays, weeks *schema.Set) {
	if days := retention.DaysOfTheWeek; days != nil {
		slice := make([]interface{}, 0)
		for _, d := range *days {
			slice = append(slice, string(d))
		}
		weekdays = schema.NewSet(schema.HashString, slice)
	}

	if days := retention.WeeksOfTheMonth; days != nil {
		slice := make([]interface{}, 0)
		for _, d := range *days {
			slice = append(slice, string(d))
		}
		weeks = schema.NewSet(schema.HashString, slice)
	}

	return weekdays, weeks
}

func resourceArmRecoveryServicesProtectionPolicyWaitForState(client backup.ProtectionPoliciesClient, ctx context.Context, found bool, vaultName, resourceGroup, policyName string) (backup.ProtectionPolicyResource, error) {
	state := &resource.StateChangeConf{
		Timeout:    30 * time.Minute,
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Refresh: func() (interface{}, string, error) {

			resp, err := client.Get(ctx, vaultName, resourceGroup, policyName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return resp, "NotFound", nil
				}

				return resp, "Error", fmt.Errorf("Error making Read request on Recovery Service Protection Policy %q (Resource Group %q): %+v", policyName, resourceGroup, err)
			}

			return resp, "Found", nil
		},
	}

	if found {
		state.Pending = []string{"NotFound"}
		state.Target = []string{"Found"}
	} else {
		state.Pending = []string{"Found"}
		state.Target = []string{"NotFound"}
	}

	resp, err := state.WaitForState()
	if err != nil {
		return resp.(backup.ProtectionPolicyResource), fmt.Errorf("Error waiting for the Recovery Service Protection Policy %q to be %t (Resource Group %q) to provision: %+v", policyName, found, resourceGroup, err)
	}

	return resp.(backup.ProtectionPolicyResource), nil
}
