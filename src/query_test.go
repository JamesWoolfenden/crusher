package crusher

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"

	"cloud.google.com/go/bigtable"
)

func TestCrusher_DeleteRows(t *testing.T) {
	type fields struct {
		TableID    *string
		ProjectID  *string
		InstanceID *string
		Filter     bigtable.Filter
	}
	type args struct {
		rows []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := &Crusher{
				TableID:    tt.fields.TableID,
				ProjectID:  tt.fields.ProjectID,
				InstanceID: tt.fields.InstanceID,
				Filter:     tt.fields.Filter,
			}
			if err := content.DeleteRows(tt.args.rows); (err != nil) != tt.wantErr {
				t.Errorf("DeleteRows() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCrusher_ReadWithFilter(t *testing.T) {
	type fields struct {
		TableID    *string
		ProjectID  *string
		InstanceID *string
		Filter     bigtable.Filter
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := &Crusher{
				TableID:    tt.fields.TableID,
				ProjectID:  tt.fields.ProjectID,
				InstanceID: tt.fields.InstanceID,
				Filter:     tt.fields.Filter,
			}
			got, err := content.ReadWithFilter()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadWithFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadWithFilter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrusher_btDelete(t *testing.T) {
	type fields struct {
		TableID    *string
		ProjectID  *string
		InstanceID *string
		Filter     bigtable.Filter
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := &Crusher{
				TableID:    tt.fields.TableID,
				ProjectID:  tt.fields.ProjectID,
				InstanceID: tt.fields.InstanceID,
				Filter:     tt.fields.Filter,
			}
			content.btDelete(tt.args.w, tt.args.r)
		})
	}
}

func TestCrusher_insertRows(t *testing.T) {
	type fields struct {
		TableID    *string
		ProjectID  *string
		InstanceID *string
		Filter     bigtable.Filter
	}
	type args struct {
		projectID  string
		instanceID string
		rows       []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := &Crusher{
				TableID:    tt.fields.TableID,
				ProjectID:  tt.fields.ProjectID,
				InstanceID: tt.fields.InstanceID,
				Filter:     tt.fields.Filter,
			}
			content.insertRows(tt.args.projectID, tt.args.instanceID, tt.args.rows)
		})
	}
}

func Test_printRow(t *testing.T) {
	type args struct {
		row bigtable.Row
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			printRow(w, tt.args.row)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("printRow() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
