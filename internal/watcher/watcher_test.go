package watcher

import (
	"reflect"
	"testing"

	"github.com/mgla96/ssh-watcher/internal/notifier"
)

func TestParseLogLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want notifier.LogLine
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseLogLine(tt.args.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseLogLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
