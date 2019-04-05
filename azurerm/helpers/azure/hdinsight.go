package azure

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/hdinsight/mgmt/2018-06-01-preview/hdinsight"
	"github.com/hashicorp/go-getter/helper/url"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SchemaHDInsightName() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validateHDInsightName,
	}
}

func SchemaHDInsightDataSourceName() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validateHDInsightName,
	}
}

func validateHDInsightName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// The name must be 59 characters or less and can contain letters, numbers, and hyphens (but the first and last character must be a letter or number).
	if matched := regexp.MustCompile(`(^[a-zA-Z0-9])([a-zA-Z0-9-]{1,57})([a-zA-Z0-9]$)`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q must be 59 characters or less and can contain letters, numbers, and hyphens (but the first and last character must be a letter or number).", k))
	}

	return warnings, errors
}

func SchemaHDInsightTier() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		ValidateFunc: validation.StringInSlice([]string{
			string(hdinsight.Standard),
			string(hdinsight.Premium),
		}, false),
		// TODO: file a bug about this
		DiffSuppressFunc: SuppressLocationDiff,
	}
}

func SchemaHDInsightClusterVersion() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// TODO: tests
			// `3.6` gets converted to `3.6.1000.67`; so let's just compare major/minor if possible
			o := strings.Split(old, ".")
			n := strings.Split(new, ".")

			if len(o) >= 2 && len(n) >= 2 {
				oldMajor := o[0]
				oldMinor := o[1]
				newMajor := n[0]
				newMinor := n[1]

				return oldMajor == newMajor && oldMinor == newMinor
			}

			return false
		},
	}
}

func SchemaHDInsightsGateway() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:     schema.TypeBool,
					Required: true,
					ForceNew: true,
				},
				// NOTE: these are Required since if these aren't present you get a `500 bad request`
				"username": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"password": {
					Type:      schema.TypeString,
					Required:  true,
					ForceNew:  true,
					Sensitive: true,
				},
			},
		},
	}
}

func ExpandHDInsightsConfigurations(input []interface{}) map[string]interface{} {
	vs := input[0].(map[string]interface{})

	// NOTE: Admin username must be different from SSH Username
	enabled := vs["enabled"].(bool)
	username := vs["username"].(string)
	password := vs["password"].(string)

	return map[string]interface{}{
		"gateway": map[string]interface{}{
			"restAuthCredential.isEnabled": enabled,
			"restAuthCredential.username":  username,
			"restAuthCredential.password":  password,
		},
	}
}

func FlattenHDInsightsConfigurations(input map[string]*string) []interface{} {
	enabled := false
	if v, exists := input["restAuthCredential.isEnabled"]; exists && v != nil {
		e, err := strconv.ParseBool(*v)
		if err == nil {
			enabled = e
		}
	}

	username := ""
	if v, exists := input["restAuthCredential.username"]; exists && v != nil {
		username = *v
	}

	password := ""
	if v, exists := input["restAuthCredential.password"]; exists && v != nil {
		password = *v
	}

	return []interface{}{
		map[string]interface{}{
			"enabled":  enabled,
			"username": username,
			"password": password,
		},
	}
}

func SchemaHDInsightsStorageAccounts() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"storage_account_key": {
					Type:      schema.TypeString,
					Required:  true,
					ForceNew:  true,
					Sensitive: true,
				},
				"storage_container_id": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"is_default": {
					Type:     schema.TypeBool,
					Required: true,
					ForceNew: true,
				},
			},
		},
	}
}

func ExpandHDInsightsStorageAccounts(input []interface{}) (*[]hdinsight.StorageAccount, error) {
	results := make([]hdinsight.StorageAccount, 0)

	for _, vs := range input {
		v := vs.(map[string]interface{})

		storageAccountKey := v["storage_account_key"].(string)
		storageContainerId := v["storage_container_id"].(string)
		isDefault := v["is_default"].(bool)

		// https://foo.blob.core.windows.net/example
		uri, err := url.Parse(storageContainerId)
		if err != nil {
			return nil, fmt.Errorf("Error parsing %q: %s", storageContainerId, err)
		}

		result := hdinsight.StorageAccount{
			Name:      utils.String(uri.Host),
			Container: utils.String(strings.TrimPrefix(uri.Path, "/")),
			Key:       utils.String(storageAccountKey),
			IsDefault: utils.Bool(isDefault),
		}
		results = append(results, result)
	}

	return &results, nil
}

type HDInsightNodeDefinition struct {
	CanSpecifyInstanceCount bool
	MinInstanceCount        int
	MaxInstanceCount        int
	ValidVmSizes            []string
	CanSpecifyDisks         bool
	MaxNumberOfDisksPerNode *int

	// TODO: pull these in for the Expand
	FixedMinInstanceCount    *int32
	FixedTargetInstanceCount *int32
}

func SchemaHDInsightNodeDefinition(schemaLocation string, definition HDInsightNodeDefinition) *schema.Schema {
	result := map[string]*schema.Schema{
		"vm_size": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			// TODO: re-enable me
			//ValidateFunc: validation.StringInSlice(definition.ValidVmSizes, false),
		},
		"username": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			ForceNew:  true,
			Sensitive: true,
			// TODO: validation
			// The password must be at least 10 characters in length and must contain at least one digit, one uppercase and one lower case letter, one non-alphanumeric character (except characters ' " ` \).
		},
		"ssh_keys": {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validate.NoEmptyStrings,
			},
			Set: schema.HashString,
			ConflictsWith: []string{
				fmt.Sprintf("%s.0.password", schemaLocation),
			},
		},

		"subnet_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: ValidateResourceIDOrEmpty,
		},

		"virtual_network_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: ValidateResourceIDOrEmpty,
		},
	}

	if definition.CanSpecifyInstanceCount {
		// TODO: should we make this validate func optional?
		result["min_instance_count"] = &schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(definition.MinInstanceCount, definition.MaxInstanceCount),
		}
		result["target_instance_count"] = &schema.Schema{
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(definition.MinInstanceCount, definition.MaxInstanceCount),
		}
	}

	if definition.CanSpecifyDisks {
		result["number_of_disks_per_node"] = &schema.Schema{
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, *definition.MaxNumberOfDisksPerNode),
		}
	}

	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: result,
		},
	}
}

func ExpandHDInsightNodeDefinition(name string, input []interface{}, definition HDInsightNodeDefinition) (*hdinsight.Role, error) {
	v := input[0].(map[string]interface{})
	vmSize := v["vm_size"].(string)
	username := v["username"].(string)
	password := v["password"].(string)
	virtualNetworkId := v["virtual_network_id"].(string)
	subnetId := v["subnet_id"].(string)

	role := hdinsight.Role{
		Name: utils.String(name),
		HardwareProfile: &hdinsight.HardwareProfile{
			VMSize: utils.String(vmSize),
		},
		OsProfile: &hdinsight.OsProfile{
			LinuxOperatingSystemProfile: &hdinsight.LinuxOperatingSystemProfile{
				Username: utils.String(username),
			},
		},
	}

	virtualNetworkSpecified := virtualNetworkId != ""
	subnetSpecified := subnetId != ""
	if virtualNetworkSpecified && subnetSpecified {
		role.VirtualNetworkProfile = &hdinsight.VirtualNetworkProfile{
			ID:     utils.String(virtualNetworkId),
			Subnet: utils.String(subnetId),
		}
	} else if (virtualNetworkSpecified && !subnetSpecified) || (subnetSpecified && !virtualNetworkSpecified) {
		return nil, fmt.Errorf("`virtual_network_id` and `subnet_id` must both either be set or empty!")
	}

	if password != "" {
		role.OsProfile.LinuxOperatingSystemProfile.Password = utils.String(password)
	} else {
		sshKeysRaw := v["ssh_keys"].(*schema.Set).List()
		sshKeys := make([]hdinsight.SSHPublicKey, 0)
		for _, v := range sshKeysRaw {
			sshKeys = append(sshKeys, hdinsight.SSHPublicKey{
				CertificateData: utils.String(v.(string)),
			})
		}

		if len(sshKeys) == 0 {
			return nil, fmt.Errorf("Either a `password` or `ssh_key` must be specified!")
		}

		role.OsProfile.LinuxOperatingSystemProfile.SSHProfile = &hdinsight.SSHProfile{
			PublicKeys: &sshKeys,
		}
	}

	if definition.CanSpecifyInstanceCount {
		minInstanceCount := v["min_instance_count"].(int)
		if minInstanceCount > 0 {
			role.MinInstanceCount = utils.Int32(int32(minInstanceCount))
		}

		targetInstanceCount := v["target_instance_count"].(int)
		role.TargetInstanceCount = utils.Int32(int32(targetInstanceCount))
	} else {
		role.MinInstanceCount = definition.FixedMinInstanceCount
		role.TargetInstanceCount = definition.FixedTargetInstanceCount
	}

	if definition.CanSpecifyDisks {
		numberOfDisksPerNode := v["number_of_disks_per_node"].(int)
		if numberOfDisksPerNode > 0 {
			role.DataDisksGroups = &[]hdinsight.DataDisksGroups{
				{
					DisksPerNode: utils.Int32(int32(numberOfDisksPerNode)),
				},
			}
		}
	}

	return &role, nil
}

func FlattenHDInsightNodeDefinition(input *hdinsight.Role, existing []interface{}, definition HDInsightNodeDefinition) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := map[string]interface{}{
		"vm_size":            "",
		"username":           "",
		"password":           "",
		"ssh_keys":           schema.NewSet(schema.HashString, []interface{}{}),
		"subnet_id":          "",
		"virtual_network_id": "",
	}

	if profile := input.OsProfile; profile != nil {
		if osProfile := profile.LinuxOperatingSystemProfile; osProfile != nil {
			if username := osProfile.Username; username != nil {
				output["username"] = *username
			}
		}
	}

	// neither Password / SSH Keys are returned from the API, so we need to look them up to not force a diff
	if len(existing) > 0 {
		existingV := existing[0].(map[string]interface{})
		output["password"] = existingV["password"].(string)

		sshKeys := existingV["ssh_keys"].(*schema.Set).List()
		output["ssh_keys"] = schema.NewSet(schema.HashString, sshKeys)

		// whilst the VMSize can be returned from `input.HardwareProfile.VMSize` - it can be malformed
		// for example, `small`, `medium`, `large` and `extralarge` can be returned inside of actual VM Size
		// after extensive experimentation it appears multiple instance sizes fit `extralarge`, as such
		// unfortunately we can't transform these; since it can't be changed
		// we should be "safe" to try and pull it from the state instead, but clearly this isn't ideal
		vmSize := existingV["vm_size"].(string)
		output["vm_size"] = vmSize
	}

	if profile := input.VirtualNetworkProfile; profile != nil {
		if profile.ID != nil {
			output["virtual_network_id"] = *profile.ID
		}
		if profile.Subnet != nil {
			output["subnet_id"] = *profile.Subnet
		}
	}

	if definition.CanSpecifyInstanceCount {
		output["min_instance_count"] = 0
		output["target_instance_count"] = 0

		if input.MinInstanceCount != nil {
			output["min_instance_count"] = int(*input.MinInstanceCount)
		}

		if input.TargetInstanceCount != nil {
			output["target_instance_count"] = int(*input.TargetInstanceCount)
		}
	}

	if definition.CanSpecifyDisks {
		output["number_of_disks_per_node"] = 0
		if input.DataDisksGroups != nil && len(*input.DataDisksGroups) > 0 {
			group := (*input.DataDisksGroups)[0]
			if group.DisksPerNode != nil {
				output["number_of_disks_per_node"] = int(*group.DisksPerNode)
			}
		}
	}

	return []interface{}{output}
}

func FindHDInsightRole(input *[]hdinsight.Role, name string) *hdinsight.Role {
	if input == nil {
		return nil
	}

	for _, v := range *input {
		if v.Name == nil {
			continue
		}

		actualName := *v.Name
		if strings.EqualFold(name, actualName) {
			return &v
		}
	}

	return nil
}

func FindHDInsightConnectivityEndpoint(name string, input *[]hdinsight.ConnectivityEndpoint) string {
	if input == nil {
		return ""
	}

	for _, v := range *input {
		if v.Name == nil || v.Location == nil {
			continue
		}

		if strings.EqualFold(*v.Name, name) {
			return *v.Location
		}
	}

	return ""
}
