package logic

import "testing"

func Test_camelCase(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"1", args{"test"}, "Test"},
		{"2", args{"test_foo"}, "TestFoo"},
		{"2", args{"_test_foo"}, "TestFoo"},
		{"2", args{"test__foo"}, "TestFoo"},
		{"2", args{"test__foo_"}, "TestFoo"},
		{"2", args{"_"}, "_"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := camelCase(tt.args.name); got != tt.want {
				t.Errorf("camelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
