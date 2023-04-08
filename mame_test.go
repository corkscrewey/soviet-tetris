package main

import (
	"testing"
)

func Test_extract7zv2(t *testing.T) {
	dir := t.TempDir()
	type args struct {
		archive string
		offset  int64
		dstdir  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "extract7zv2",
			args: args{
				archive: "1/mame0253b_64bit.exe",
				offset:  205824,
				dstdir:  dir,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := extract7zv2(tt.args.archive, tt.args.offset, tt.args.dstdir); (err != nil) != tt.wantErr {
				t.Errorf("extract7zv2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
