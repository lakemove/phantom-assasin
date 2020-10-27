package learn

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestECKeyGen(t *testing.T) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(t, err)
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
	assert.Nil(t, err)
	assert.IsType(t, new(ecdsa.PrivateKey), decodedKey)
	assert.Equal(t, priv, decodedKey)
}

func TestRSAKeyGen(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(t, err)
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	assert.Nil(t, err)
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
	assert.Nil(t, err)
	assert.IsType(t, new(rsa.PrivateKey), decodedKey)
	assert.Equal(t, priv, decodedKey)
}

func TestECSign(t *testing.T) {
	//SHA256 with ECDSA
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(t, err)
	message := "Hello Jay"
	digest := sha256.Sum256([]byte(message))
	signature, err := priv.Sign(rand.Reader, digest[:], crypto.SHA256)
	assert.Nil(t, err, "sign ecdsa")
	sigHex := base64.StdEncoding.EncodeToString(signature)
	assert.True(t, len(sigHex) <= 96)
}

func TestRSASign(t *testing.T) {
	//SHA256 with RSA
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(t, err)
	message := "Hello Jay"
	digest := sha256.Sum256([]byte(message))
	signature, err := priv.Sign(rand.Reader, digest[:], crypto.SHA256)
	assert.Nil(t, err, "sign rsa")
	sigHex := base64.StdEncoding.EncodeToString(signature)
	t.Log(sigHex)
	assert.True(t, len(sigHex) <= 96)
}
