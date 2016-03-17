package secret

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/square/go-jose"
	"gopkg.in/yaml.v2"
)

// Parse parses and returns the secure section of the
// yaml file as plaintext parameters.
func Parse(in, privKey string) (*File, error) {
	// unarmshal the private key from PEM
	rsaPrivKey, err := decodePrivateKey(privKey)
	if err != nil {
		return nil, err
	}

	// decrypt the Yaml file
	plain, err := decrypt(in, rsaPrivKey)
	if err != nil {
		return nil, err
	}

	// unmarshal the yaml
	out := &file{}
	err = yaml.Unmarshal(plain, out)
	if err != nil {
		return nil, err
	}

	// convert the yaml structure to a structure
	// that is a bit easier to work with in Go.
	file := &File{}
	file.Checksum = out.Checksum
	for k, v := range out.Runtime {
		secret := &Secret{
			Image: v.Image.Slice(),
			Event: v.Event.Slice(),
			Data:  v.Data,
		}
		if v.Image.Len() == 0 {
			secret.Image = []string{k}
		}
		file.Runtime = append(file.Runtime, secret)
	}
	for k, v := range out.Registry {
		registry := &Registry{
			Hostname: v.Hostname,
			Username: v.Username,
			Password: v.Password,
			Email:    v.Email,
		}
		if registry.Hostname == "" {
			registry.Hostname = k
		}
		file.Registry = append(file.Registry, registry)
	}
	return file, nil
}

// decrypt decrypts a JOSE string and returns the
// plaintext value.
func decrypt(secret string, privKey *rsa.PrivateKey) ([]byte, error) {
	object, err := jose.ParseEncrypted(secret)
	if err != nil {
		return nil, err
	}
	return object.Decrypt(privKey)
}

// decodePrivateKey is a helper function that unmarshals a PEM
// bytes to an RSA Private Key
func decodePrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	derBlock, _ := pem.Decode([]byte(privateKey))
	return x509.ParsePKCS1PrivateKey(derBlock.Bytes)
}

// encodePrivateKey is a helper function that marshals an RSA
// Private Key to a PEM encoded file.
func encodePrivateKey(privkey *rsa.PrivateKey) string {
	privateKeyMarshaled := x509.MarshalPKCS1PrivateKey(privkey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PRIVATE KEY",
		Bytes:   privateKeyMarshaled,
		Headers: nil,
	})
	return string(privateKeyPEM)
}

// encrypt encrypts a plaintext variable using JOSE with
// RSA_OAEP and A128GCM algorithms.
func encrypt(text string, pubKey *rsa.PublicKey) (string, error) {
	var encrypted string
	var plaintext = []byte(text)

	// Creates a new encrypter using defaults
	encrypter, err := jose.NewEncrypter(jose.RSA_OAEP, jose.A128GCM, pubKey)
	if err != nil {
		return encrypted, err
	}
	// Encrypts the plaintext value and serializes
	// as a JOSE string.
	object, err := encrypter.Encrypt(plaintext)
	if err != nil {
		return encrypted, err
	}
	return object.CompactSerialize()
}
