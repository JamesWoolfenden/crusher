package main

import (
	"testing"
)

func TestValidateInput(t *testing.T) {
	tests := []struct {
		name       string
		projectID  string
		instanceID string
		tableID    string
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "valid inputs",
			projectID:  "my-project",
			instanceID: "my-instance",
			tableID:    "my-table",
			wantErr:    false,
		},
		{
			name:       "empty project ID",
			projectID:  "",
			instanceID: "my-instance",
			tableID:    "my-table",
			wantErr:    true,
			errMsg:     "project ID cannot be empty",
		},
		{
			name:       "whitespace project ID",
			projectID:  "   ",
			instanceID: "my-instance",
			tableID:    "my-table",
			wantErr:    true,
			errMsg:     "project ID cannot be empty",
		},
		{
			name:       "empty instance ID",
			projectID:  "my-project",
			instanceID: "",
			tableID:    "my-table",
			wantErr:    true,
			errMsg:     "instance ID cannot be empty",
		},
		{
			name:       "whitespace instance ID",
			projectID:  "my-project",
			instanceID: "   ",
			tableID:    "my-table",
			wantErr:    true,
			errMsg:     "instance ID cannot be empty",
		},
		{
			name:       "empty table ID",
			projectID:  "my-project",
			instanceID: "my-instance",
			tableID:    "",
			wantErr:    true,
			errMsg:     "table ID cannot be empty",
		},
		{
			name:       "whitespace table ID",
			projectID:  "my-project",
			instanceID: "my-instance",
			tableID:    "   ",
			wantErr:    true,
			errMsg:     "table ID cannot be empty",
		},
		{
			name:       "all empty",
			projectID:  "",
			instanceID: "",
			tableID:    "",
			wantErr:    true,
			errMsg:     "project ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInput(tt.projectID, tt.instanceID, tt.tableID)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("validateInput() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}
