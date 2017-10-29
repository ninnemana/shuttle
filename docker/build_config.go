package docker

import (
	"strings"

	"github.com/docker/docker/api/types"
	validation "github.com/go-ozzo/ozzo-validation"
)

// BuildConfig describes some meta parameters
// that we will listen for when building a docker image.
type BuildConfig struct {
	Directory  *string
	Dockerfile *string
	Remove     *bool
	Tags       []string
}

const (
	// DefaultDockerfile is what we will fallback on for the name of
	// the docker build definition.
	DefaultDockerfile = "Dockerfile"
)

// Validate checks the validation requirements for fields on BuildConfig.
func (bc BuildConfig) Validate() error {
	return validation.ValidateStruct(&bc,
		validation.Field(&bc.Directory, validation.Required),
	)
}

// ToOptions converts our BuildConfig to types.ImageBuildOptions
// which can be used by the Docker service.
func (bc BuildConfig) ToOptions() types.ImageBuildOptions {
	var dockerfile string
	var remove bool

	switch {
	case bc.Dockerfile == nil:
		dockerfile = DefaultDockerfile
	case strings.TrimSpace(*bc.Dockerfile) == "":
		dockerfile = DefaultDockerfile
	default:
		dockerfile = *bc.Dockerfile
	}

	if bc.Remove != nil && *bc.Remove != false {
		remove = true
	}

	return types.ImageBuildOptions{
		Tags:           bc.Tags,
		SuppressOutput: false,
		Remove:         remove,
		ForceRemove:    remove,
		PullParent:     true,
		Dockerfile:     dockerfile,
	}
}
