package crusher

import "cloud.google.com/go/bigtable"

type Crusher struct {
	TableID    *string
	ProjectID  *string
	InstanceID *string
	Filter     bigtable.Filter
}
