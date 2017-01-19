package secret

import "github.com/drone/drone/model"

// Secret represents a secret variable.
type Secret struct {
	// the id for this secret.
	ID int64 `json:"id"`

	// the name of the secret which will be used as the environment variable
	// name at runtime.
	Name string `json:"name"`

	// the value of the secret which will be provided to the runtime environment
	// as a named environment variable.
	Value string `json:"value"`

	// the secret is restricted to this list of events.
	Events []string `json:"event"`

	// whether the secret requires verification
	SkipVerify bool `json:"skip_verify"`

	// whether the secret should be concealed in the build log
	Conceal bool `json:"conceal"`
}

// Store defines a storage driver for getting and setting secrets.
type Store interface {
	GetSecret(*model.Repo, string) (*Secret, error)
	GetSecretList(*model.Repo) ([]*Secret, error)
	SetSecret(*model.Repo, *Secret) error
	DelSecret(*model.Repo, *Secret) error
}
