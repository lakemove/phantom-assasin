package learn

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestECKeyGen(t *testing.T) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Equal(t, nil, err)
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	privBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	}
	buf := new(bytes.Buffer)
	pem.Encode(buf, privBlock)
	assert.Contains(t, buf.String(), "PRIVATE KEY")
	// decode
	decodedBlock, _ := pem.Decode(buf.Bytes())
	assert.True(t, decodedBlock != nil)
	decodedKey, err := x509.ParsePKCS8PrivateKey(decodedBlock.Bytes)
	assert.Equal(t, nil, err)
	assert.IsType(t, new(ecdsa.PrivateKey), decodedKey)
	assert.Equal(t, priv, decodedKey)
}

func TestRSAKeyGen(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Equal(t, nil, err)
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	assert.Equal(t, nil, err)
	privBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	}
	buf := new(bytes.Buffer)
	pem.Encode(buf, privBlock)
	assert.Contains(t, buf.String(), "PRIVATE KEY")
	// decode
	decodedBlock, _ := pem.Decode(buf.Bytes())
	assert.True(t, decodedBlock != nil)
	decodedKey, err := x509.ParsePKCS8PrivateKey(decodedBlock.Bytes)
	assert.Equal(t, nil, err)
	assert.IsType(t, new(rsa.PrivateKey), decodedKey)
	assert.Equal(t, priv, decodedKey)
}

func TestECSign(t *testing.T) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Equal(t, nil, err)
	signature, err := priv.Sign(rand.Reader, []byte("Hello Jay"), nil)
	sigHex := base64.StdEncoding.EncodeToString(signature)
	assert.True(t, len(sigHex) <= 96)
}

func TestRSASign(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Equal(t, nil, err)
	signature, err := priv.Sign(rand.Reader, []byte("Hello Jay"), nil)
	sigHex := base64.StdEncoding.EncodeToString(signature)
	assert.True(t, len(sigHex) <= 96)
}
