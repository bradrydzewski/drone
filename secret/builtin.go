package secret

import (
	"github.com/drone/drone/model"
	"github.com/drone/drone/store"
)

type builtin struct {
	store.Store
}

// NewStore returns a new implementation of the secret store
// using a database connection.
func NewStore(s store.Store) Store {
	return &builtin{s}
}

func (s *builtin) GetSecret(repo *model.Repo, name string) (*Secret, error) {
	tmp, err := s.Store.GetSecret(repo, name)
	if err != nil {
		return nil, err
	}
	return fromRepoSecret(tmp), nil
}

func (s *builtin) GetSecretList(repo *model.Repo) ([]*Secret, error) {
	tmps, err := s.Store.GetSecretList(repo)
	if err != nil {
		return nil, err
	}
	secrets := []*Secret{}
	for _, tmp := range tmps {
		secrets = append(secrets, fromRepoSecret(tmp))
	}
	return secrets, nil
}

func (s *builtin) SetSecret(repo *model.Repo, secret *Secret) error {
	rsecret := toRepoSecret(repo, secret)
	err := s.Store.SetSecret(rsecret)
	if err != nil {
		return err
	}
	secret.ID = rsecret.ID
	return nil
}

func (s *builtin) DelSecret(repo *model.Repo, secret *Secret) error {
	return s.Store.DeleteSecret(
		toRepoSecret(repo, secret),
	)
}

func fromRepoSecret(secret *model.RepoSecret) *Secret {
	return &Secret{
		ID:         secret.ID,
		Name:       secret.Name,
		Value:      secret.Value,
		Events:     secret.Events,
		SkipVerify: secret.SkipVerify,
		Conceal:    secret.Conceal,
	}
}

func toRepoSecret(repo *model.Repo, secret *Secret) *model.RepoSecret {
	return &model.RepoSecret{
		ID:         secret.ID,
		Name:       secret.Name,
		Value:      secret.Value,
		Events:     secret.Events,
		SkipVerify: secret.SkipVerify,
		Conceal:    secret.Conceal,
		RepoID:     repo.ID,
	}
}
