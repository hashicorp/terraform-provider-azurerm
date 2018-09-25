package azurerm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/terraform/helper/schema"
)

type KeyManagemetResponse struct {
	Keys  *[]KeyManagementKey  `json:"keys"`
	Links *[]KeyManagementLink `json:"links"`
}

type KeyManagementKey struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type KeyManagementLink struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

func dataSourceArmFunctionAppKeyManagement() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmFuncionAppKeyManagementRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"function_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_name": resourceGroupNameDiffSuppressSchema(),
			"function_keys": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func getFunctionKeys(functionAppName string, functionName string, token string) (*[]KeyManagementKey, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", "https://"+functionAppName+".azurewebsites.net/admin/functions/"+functionName+"/keys", nil)
	if err != nil {
		return nil, err
	}
	authorization := "Bearer " + token
	req.Header.Add("Authorization", authorization)
	res, err := httpClient.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	var keyResponse KeyManagemetResponse
	err = json.Unmarshal(body, &keyResponse)
	if err != nil {
		return nil, err
	}
	return keyResponse.Keys, nil
}

func dataSourceArmFuncionAppKeyManagementRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesClient

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	functionName := d.Get("function_name").(string)

	ctx := meta.(*ArmClient).StopContext

	token, err := client.GetFunctionsAdminToken(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM Function App Admin Token %q: %+v", name, err)
	}

	keys, err := getFunctionKeys(name, functionName, to.String(token.Value))
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM Function App Key Management API %q: %+v", name, err)
	}

	functionKeys := make(map[string]string)
	for _, key := range *keys {
		functionKeys[key.Name] = key.Value
	}

	// d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("function_name", functionName)
	d.Set("resource_group_name", resourceGroup)
	d.Set("function_keys", functionKeys)

	return nil
}
