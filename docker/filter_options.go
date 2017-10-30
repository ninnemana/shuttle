package docker

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// FilterOptions defines the available set of query parameters
// that can be used to query existing docker images.
type FilterOptions struct {
	All      *bool   `json:"all"`
	NameLike *string `json:"nameLike"`
	Name     *string `json:"name"`
	Before   *string `json:"before"`
	After    *string `json:"after"`
	Dangling *bool   `json:"dangling"`
	Tag      *string `json:"tag"`
}

// Validate determines if the provided FilterOptions are suitable for listing
// available docker images.
func (f FilterOptions) Validate() error {
	return validation.ValidateStruct(&f, nil)
}

// Build converts the FilterOptions into types.ImageListOptions that can be
// used by the Docker SDK.
func (f FilterOptions) Build() types.ImageListOptions {

	opts := filters.NewArgs()
	reference := ""
	switch {
	case f.NameLike == nil:
		reference = "*"
	case *f.NameLike == "":
		reference = "*"
	default:
		reference = fmt.Sprintf("*%s*", *f.NameLike)
	}

	if f.Name != nil && *f.Name != "" {
		reference = *f.Name
	}
	if f.Tag != nil && *f.Tag != "" {
		reference = fmt.Sprintf("%s:%s", reference, *f.Tag)
	}

	opts.Add("reference", reference)

	args := types.ImageListOptions{
		Filters: opts,
	}

	if f.All != nil && *f.All == true {
		args.All = true
	}

	return args
}
