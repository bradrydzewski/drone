package secret

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/drone/drone/model"
)

type remote struct {
	addr string
}

// NewRemoteStore returns a remote implementation of the secret store
// at the specified address.
func NewRemoteStore(addr string) Store {
	return &remote{addr}
}

// GetSecret gets the named secret from the remote endpoint for the
// given repository.
func (s *remote) GetSecret(repo *model.Repo, name string) (*Secret, error) {
	var (
		item = new(Secret)
		path = path.Join("/secrets", repo.Owner, repo.Name, name)
	)
	err := s.do("GET", path, nil, item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// GetSecretList makes a request to the remote endpoint to fetch
// the list of all registered secrets for the given repository.
func (s *remote) GetSecretList(repo *model.Repo) ([]*Secret, error) {
	var (
		list = []*Secret{}
		path = path.Join("/secrets", repo.Owner, repo.Name)
	)
	err := s.do("GET", path, nil, list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// SetSecret makes a request to the remote endpoint to create
// or update the secret for the given repository.
func (s *remote) SetSecret(repo *model.Repo, secret *Secret) error {
	path := path.Join("/secrets", repo.Owner, repo.Name)
	return s.do("POST", path, secret, nil)
}

// DelSecret makes a request to the remote endpoint to delete
// the given secret for the given repository.
func (s *remote) DelSecret(repo *model.Repo, secret *Secret) error {
	path := path.Join("/secrets", repo.Owner, repo.Name)
	return s.do("DELETE", path, nil, nil)
}

func (s *remote) do(method, path string, in, out interface{}) error {
	uri, err := url.Parse(path)
	if err != nil {
		return err
	}

	// if we are posting or putting data, we need to
	// write it to the body of the request.
	var buf io.ReadWriter
	if in != nil {
		buf = new(bytes.Buffer)
		jsonerr := json.NewEncoder(buf).Encode(in)
		if jsonerr != nil {
			return jsonerr
		}
	}

	// creates a new http request to bitbucket.
	req, err := http.NewRequest(method, uri.String(), buf)
	if err != nil {
		return err
	}
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// if an error is encountered, parse and return the
	// error response.
	if resp.StatusCode > http.StatusPartialContent {
		err := Error{}
		json.NewDecoder(resp.Body).Decode(&err)
		return &err
	}

	// if a json response is expected, parse and return
	// the json response.
	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}

	return nil
}

// Error represents a http error.
type Error struct {
	Message string `json:"message"`
}

// Error returns the error message in string format.
func (e *Error) Error() string {
	return e.Message
}
