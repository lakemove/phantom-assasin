package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBindXML(t *testing.T) {
	x0 := `<MxML><deliverable><currency>LUS</currency></deliverable></MxML>`
	var v struct {
		Currency string `xml:"deliverable>currency"`
	}
	xml.Unmarshal([]byte(x0), &v)
	assert.Equal(t, "LUS", v.Currency)
}

func TestECKeyGen(t *testing.T) {
	curve := elliptic.P256()
	priv, x, y, err := elliptic.GenerateKey(curve, rand.Reader)
	fmt.Println(priv, x, y)
	assert.Equal(t, nil, err)
	assert.Equal(t, "LUS", "LUS")
}
