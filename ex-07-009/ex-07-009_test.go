package main

import (
	"reflect"
	"testing"
)

func Test_parseSortQuery(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    order
		wantErr bool
	}{{
		name: "one arg",
		args: args{"title"},
		want: order{
			field: "title",
			order: asc,
		},
		wantErr: false,
	},
		{
			name: "empty",
			args: args{""},
			want: order{},
			wantErr: true,
		},
		{
			name: "many args",
			args: args{"title[asc]"},
			want: order{
				field: "title",
				order: asc,
			},
			wantErr: false,
		},
		{
			name: "many args asc default",
			args: args{"year"},
			want: order{
				field: "year",
				order: asc,
			},
			wantErr: false,
		},
		{
			name: "many args desc",
			args: args{"artist[desc]"},
			want: order{
				field: "artist",
				order: desc,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSortQuery(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSortQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSortQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
