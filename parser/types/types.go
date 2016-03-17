package types

import (
	"strings"
)

// MapEqualSlice represents a map of key value items in string format
// which can alternatively be declared as a slice of string items
// in key=value format.
type MapEqualSlice struct {
	parts map[string]string
}

// MarshalYAML implements the Marshaller interface.
func (s MapEqualSlice) MarshalYAML() (interface{}, error) {
	return s.parts, nil
}

// UnmarshalYAML implements the Unmarshaller interface.
func (s *MapEqualSlice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	s.parts = map[string]string{}
	err := unmarshal(&s.parts)
	if err == nil {
		return nil
	}

	var sliceType []string

	err = unmarshal(&sliceType)
	if err != nil {
		return err
	}

	for _, v := range sliceType {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			val := parts[1]
			s.parts[key] = val
		}
	}

	return nil
}

// Map gets the parts of the MapEqualSlice as a map of key value
// pairs in string format.
func (s *MapEqualSlice) Map() map[string]string {
	return s.parts
}

// Len returns the number of parts of the MapEqualSlice.
func (s *MapEqualSlice) Len() int {
	if s == nil {
		return 0
	}
	return len(s.parts)
}

// Stringorslice represents a string or an array of strings.
// TODO use docker/docker/pkg/stringutils.StrSlice once 1.9.x is released.
type Stringorslice struct {
	parts []string
}

// MarshalYAML implements the Marshaller interface.
func (s Stringorslice) MarshalYAML() (interface{}, error) {
	return s.parts, nil
}

// UnmarshalYAML implements the Unmarshaller interface.
func (s *Stringorslice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var sliceType []string
	err := unmarshal(&sliceType)
	if err == nil {
		s.parts = sliceType
		return nil
	}

	var stringType string
	err = unmarshal(&stringType)
	if err == nil {
		sliceType = make([]string, 0, 1)
		s.parts = append(sliceType, string(stringType))
		return nil
	}
	return err
}

// Len returns the number of parts of the Stringorslice.
func (s *Stringorslice) Len() int {
	if s == nil {
		return 0
	}
	return len(s.parts)
}

// Slice gets the parts of the StrSlice as a Slice of string.
func (s *Stringorslice) Slice() []string {
	if s == nil {
		return nil
	}
	return s.parts
}
