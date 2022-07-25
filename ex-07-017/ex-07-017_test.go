package main

import (
	"encoding/xml"
	"reflect"
	"testing"
)

func Test_parseSelector(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    Selector
		wantErr bool
	}{
		{"simple", args{[]string{"help"}}, Selector{{Name: xml.Name{Local: "help"}}}, false},
		{"with attr",
			args{[]string{"help[me=you]"}},
			Selector{{
				Name: xml.Name{Local: "help"},
				Attr: []xml.Attr{{Name: xml.Name{Local: "me"}, Value: "you"}}}},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSelector(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSelector() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSelector() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
