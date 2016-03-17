package secret

import "github.com/drone/drone/parser/types"

type File struct {
	Checksum string      `json:"checksum"`
	Runtime  []*Secret   `json:"runtime"`
	Registry []*Registry `json:"registry"`
}

type Secret struct {
	Image []string          `json:"image"`
	Event []string          `json:"event"`
	Data  map[string]string `json:"data"`
}

type Registry struct {
	Hostname string `json:"hostname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

//
// yaml representations of the above structures. we unmarshal
// to these intermediary types and then convert to something
// that is more Go friendly.
//

type file struct {
	Checksum string               `yaml:"checksum"`
	Runtime  map[string]*runtime  `yaml:"runtime"`
	Registry map[string]*registry `yaml:"registry"`
}

type runtime struct {
	Image types.Stringorslice `yaml:"image"`
	Event types.Stringorslice `yaml:"event"`
	Data  map[string]string   `yaml:",inline"`
}

type registry struct {
	Hostname string `yaml:"hostname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Email    string `yaml:"email"`
}
