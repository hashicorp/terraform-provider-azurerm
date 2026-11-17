// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func AuthorizationRuleCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) error {
	listen, hasListen := d.GetOk("listen")
	send, hasSend := d.GetOk("send")
	manage, hasManage := d.GetOk("manage")

	if !hasListen && !hasSend && !hasManage {
		return errors.New("one of the `listen`, `send` or `manage` properties needs to be set")
	}

	if manage.(bool) && (!listen.(bool) || !send.(bool)) {
		return errors.New("if `manage` is set both `listen` and `send` must be set to true too")
	}

	return nil
}
