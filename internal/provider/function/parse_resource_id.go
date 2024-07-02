package function

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ParseResourceIDFunction struct{}

var _ function.Function = ParseResourceIDFunction{}

var idParseResultTypes = map[string]attr.Type{
	"resource_name":       types.StringType,
	"resource_provider":   types.StringType,
	"resource_group_name": types.StringType,
	"resource_type":       types.StringType,
	"resource_scope":      types.StringType,
	"full_resource_type":  types.StringType,
	"subscription_id":     types.StringType,
	"parent_resources":    types.MapType{}.WithElementType(types.StringType),
}

var knownScopePrefixes = []string{
	"/subscriptions/",
	"/subscriptions//resourcegroups/",
	"/providers/Microsoft.Marketplace",
	"/providers/Microsoft.Subscription",
	"/providers/Microsoft.Management/managementGroups/",
	"/providers/Microsoft.Billing/billingAccounts//customers/",
	"/providers/Microsoft.Billing/billingAccounts//billingProfiles//invoiceSections/",
	"/providers/Microsoft.Billing/billingAccounts//enrollmentAccounts/",
}

func NewParseResourceIDFunction() function.Function {
	return &ParseResourceIDFunction{}
}

func (p ParseResourceIDFunction) Metadata(_ context.Context, _ function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "parse_resource_id"
}

func (p ParseResourceIDFunction) Definition(ctx context.Context, request function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Summary:             "parse_resource_id",
		Description:         "Parses an Azure Resource Manager ID and exposes the contained information",
		MarkdownDescription: "Parses an Azure Resource Manager ID and exposes the contained information",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "id",
				Description:         "Resource ID",
				MarkdownDescription: "Resource ID",
			},
		},
		Return: function.ObjectReturn{
			AttributeTypes: idParseResultTypes,
		},
	}
}

func (p ParseResourceIDFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var id string

	response.Error = function.ConcatFuncErrors(request.Arguments.Get(ctx, &id))

	if response.Error != nil {
		return
	}

	if len(id) == 0 {
		response.Error = function.NewFuncError("Got empty ID")
		return
	}

	// Try and fixup being passed a value missing the initial `/` as some APIs in Azure incorrectly omit it, but that's not the user's fault
	if !strings.HasPrefix(id, "/") {
		id = fmt.Sprintf("/%s", id)
	}

	idKey := ""

	segments := strings.Split(id, "/")
	if len(segments)%2 != 0 {
		for i := 1; len(segments) > i; i++ {
			if i%2 != 0 {
				key := segments[i]
				idKey = fmt.Sprintf("%s/%s/", idKey, key)
				if strings.EqualFold(key, "providers") && len(segments) >= i+2 {
					value := segments[i+1]
					idKey = fmt.Sprintf("%s%s", idKey, value)
				}
			}
		}
	}

	idKey = strings.ToLower(idKey)

	// These outputs should always have a value
	output := map[string]attr.Value{
		"resource_name":       types.StringValue(""),
		"resource_provider":   types.StringValue(""),
		"resource_group_name": types.StringValue(""),
		"resource_type":       types.StringValue(""),
		"resource_scope":      types.StringValue(""),
		"full_resource_type":  types.StringValue(""),
		"subscription_id":     types.StringValue(""),
		"parent_resources":    types.Map{},
	}
	idType := recaser.KnownResourceIds()[idKey]
	if idType == nil {
		// This might be a scoped ID type, lets try striping the common scope prefixes
		for _, v := range recaser.PotentialScopeValues() {
			if idType = recaser.KnownResourceIds()[strings.TrimPrefix(idKey, strings.TrimSuffix(v, "/"))]; idType != nil {
				break
			}
			if idType = recaser.KnownResourceIds()[strings.TrimPrefix(idKey, v)]; idType != nil {
				break
			}
		}
		if idType == nil {
			response.Error = function.NewFuncError("Unsupported Resource ID type")
			return
		}
	}
	parser := resourceids.NewParserFromResourceIdType(idType)
	parsed, err := parser.Parse(id, true)
	if err != nil {
		response.Error = function.NewFuncError(fmt.Sprintf("Parsing Resource ID Error: %s", err))
		return
	}
	err = idType.FromParseResult(*parsed)
	if err != nil {
		response.Error = function.NewFuncError(fmt.Sprintf("Expanding Parsed Resource ID Error: %s", err))
		return
	}

	s := idType.Segments()
	numSegments := len(s)
	pTemp := ""
	fullResourceType := ""
	parentMap := map[string]string{}
	for k, v := range s {
		switch v.Type {
		case resourceids.ResourceGroupSegmentType:
			output["resource_group_name"] = types.StringValue(parsed.Parsed[v.Name])

		case resourceids.ResourceProviderSegmentType:
			output["resource_provider"] = types.StringPointerValue(v.FixedValue)
			fullResourceType = pointer.From(v.FixedValue)

		case resourceids.SubscriptionIdSegmentType:
			output["subscription_id"] = types.StringValue(parsed.Parsed["subscriptionId"])

		case resourceids.StaticSegmentType:
			switch {
			case k == (numSegments - 2):
				{
					output["resource_type"] = types.StringPointerValue(v.FixedValue)
					fullResourceType = fmt.Sprintf("%s/%s", fullResourceType, pointer.From(v.FixedValue))
				}
			case v.FixedValue != nil && *v.FixedValue != "subscriptions" && *v.FixedValue != "resourceGroups" && *v.FixedValue != "providers":
				{
					pTemp = parsed.Parsed[v.Name]
					fullResourceType = fmt.Sprintf("%s/%s", fullResourceType, pointer.From(v.FixedValue))
				}
			}

		case resourceids.UserSpecifiedSegmentType:
			if k == (numSegments - 1) {
				output["resource_name"] = types.StringValue(parsed.Parsed[v.Name])
			} else {
				parentMap[pTemp] = parsed.Parsed[v.Name]
			}
		case resourceids.ScopeSegmentType:
			output["resource_scope"] = types.StringValue(parsed.Parsed[v.Name])
		}
	}
	parentMapValue, diags := types.MapValueFrom(ctx, types.StringType, parentMap)
	if diags.HasError() {
		response.Error = function.NewFuncError("failed to flatten parent resources")
		return
	}
	output["parent_resources"] = parentMapValue
	output["full_resource_type"] = types.StringValue(fullResourceType)

	result, diags := types.ObjectValue(idParseResultTypes, output)
	if diags.HasError() {
		response.Error = function.ConcatFuncErrors(response.Error, function.FuncErrorFromDiags(ctx, diags))
		return
	}

	response.Error = function.ConcatFuncErrors(response.Result.Set(ctx, result))
}
