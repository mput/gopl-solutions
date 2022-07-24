package main

import (
	"testing"
)

func Test_parseDollars(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    dollars
		wantErr bool
	}{
		{
			name:    "simple",
			args:    args{s: "2"},
			want:    200,
			wantErr: false,
		}, {
			name:    "simple",
			args:    args{s: "2.01"},
			want:    201,
			wantErr: false,
		}, {
			name:    "simple",
			args:    args{s: "0.01"},
			want:    1,
			wantErr: false,
		}, {
			name:    "simple",
			args:    args{s: "2."},
			want:    0,
			wantErr: true,
		},{
			name:    "simple",
			args:    args{s: "2.2"},
			want:    220,
			wantErr: false,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDollars(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDollars() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseDollars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dollars_String(t *testing.T) {
	tests := []struct {
		name string
		d    dollars
		want string
	}{
		{
			name: "simple",
			d:    dollars(1),
			want: "0.01"},
		{
			name: "simple",
			d:    dollars(10),
			want: "0.10"},
		{
			name: "simple",
			d:    dollars(101),
			want: "1.01"},
		{
			name: "simple",
			d:    dollars(1010),
			want: "10.10"},
		{
			name: "simple",
			d:    dollars(10075),
			want: "100.75"},
		{
			name: "simple",
			d:    dollars(0),
			want: "0.00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("dollars.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
