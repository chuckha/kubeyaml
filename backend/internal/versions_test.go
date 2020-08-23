package internal

import (
	"reflect"
	"testing"
)

func TestSortVersionsNoErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "simple test",
			args: []string{"1.15", "1.17", "1.18", "1.16"},
			want: []string{"1.15", "1.16", "1.17", "1.18"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SortVersions(tt.args...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortVersions() = %v, want %v", got, tt.want)
			}
			if err != nil {
				t.Errorf("err should be nil but is %v", err)
			}
		})
	}
}
