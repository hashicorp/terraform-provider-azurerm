package logic

import (
	"go/types"
	"strings"
	"testing"
)

func Test_newSelectorExpr(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{"", "a.b", "a.b"},
		{"", "model.Field", "model.Field"},
		{"", "model", ""},
		{"", "", ""},
		{"", "a.b.c.d.e.f.g.h.a.b.c.s", "a.b.c.d.e.f.g.h.a.b.c.s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newSelectorExpr(strings.Split(tt.args, ".")...)
			if got == nil {
				if tt.want != "" {
					t.Errorf("newSelectorExpr() = nil, want %v", tt.want)
				}
				return
			}
			if types.ExprString(got) != tt.want {
				t.Errorf("newSelectorExpr() = %v, want %v", types.ExprString(got), tt.want)
			}
		})
	}
}
