package crusher

import "cloud.google.com/go/bigtable"

type Crusher struct {
	TableID    string
	ProjectID  string
	InstanceID string
	KeyFilter  string
	DryRun     bool
	Filter     bigtable.Filter
}
