package planmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

// TODO: rename this file and modifier to something more descriptive

type Modifier struct {
	privateStateKey     string
	privateStateValue   string
	privateStateDefault string
}

var _ planmodifier.String = &Modifier{}

// TODO: rename
func SuppressIfPrivateStateNilOrDoesNotEqual(privateStateKey, privateStateValue, privateStateDefault string) Modifier {
	return Modifier{
		privateStateKey:     privateStateKey,
		privateStateValue:   privateStateValue,
		privateStateDefault: privateStateDefault,
	}
}

func (m Modifier) Description(ctx context.Context) string {
	return "TODO"
}

func (m Modifier) MarkdownDescription(ctx context.Context) string {
	return "TODO"
}

// Sets a private state property used for tracking whether users have set something in config
// to then be used in the Read() funcs
func (m Modifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Do not check on resource creation
	if req.State.Raw.IsNull() {
		return
	}

	// Do not check on resource destroy
	if req.Plan.Raw.IsNull() {
		return
	}

	// If values are equal, ignore
	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	//private, diags := req.Private.GetKey(ctx, m.privateStateKey)
	//if diags.HasError() {
	//	resp.Diagnostics.Append(diags...)
	//	return
	//}

	resp.Diagnostics.AddWarning("plan modifier run", "")

	// If raw is null value, set to default, meaning the whatever property this is tracking in priv state
	// was not set by the user.
	if req.ConfigValue.IsNull() {
		resp.Private.SetKey(ctx, m.privateStateKey, sdk.NewPrivateStateValue(m.privateStateDefault).Bytes())
	} else {
		resp.Private.SetKey(ctx, m.privateStateKey, sdk.NewPrivateStateValue(m.privateStateValue).Bytes())
	}
}
