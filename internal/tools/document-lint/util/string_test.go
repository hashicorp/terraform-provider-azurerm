// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"reflect"
	"testing"
)

func TestExtractCodeValue(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		wantRes []string
	}{
		{
			args: args{
				line: "",
			},
			wantRes: nil,
		},
		{
			args: args{
				line: "defauts to `def`.",
			},
			wantRes: []string{"def"},
		},
		{
			args: args{
				line: "`def`",
			},
			wantRes: []string{"def"},
		},
		{
			args: args{
				line: "abc is `def` `ghi` done",
			},
			wantRes: []string{"def", "ghi"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := ExtractCodeValue(tt.args.line); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("ExtractCodeValue() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestTimeoutIdx(t *testing.T) {
	line := "virtual_desktop_scaling_plan.html.markdown:* `create` - (Defaults to 1 hour) Used when creating the Virtual Desktop Scaling Plan."
	idxs := timeoutValueReg.FindStringSubmatchIndex(line)
	t.Logf("%v: %s", idxs, line[idxs[0]:idxs[1]])

	line = "* `delete` - (Defaults to 30 minutes) Used when deleting the EventGrid Event Subscription."
	idxs = timeoutValueReg.FindStringSubmatchIndex(line)
	t.Logf("%v: %s", idxs, line[idxs[0]:idxs[1]])

}
