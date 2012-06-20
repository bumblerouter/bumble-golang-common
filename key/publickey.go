package key

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
)

type pubKeyRegistrationPost struct {
	From   string `json:"f,omitempty"`
	PubKey string `json:"k,omitempty"`
}

func GetPublicKey(domain string, peerhash string) *rsa.PublicKey {
	url, err := GetPublicKeyURL(domain, peerhash)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	r, err := http.Get(url.String())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if r.StatusCode > 400 {
		//fmt.Printf("FAILED STATUS: %d\n", r.StatusCode)
		return nil
	}
	encap := make([]byte, r.ContentLength)
	r.Body.Read(encap)
	pubkey := PublicKeyFromPEM(string(encap))
	return pubkey
}

func GetPublicKeyURL(domain string, peerhash string) (*url.URL, error) {
	txts, err := net.LookupTXT(fmt.Sprintf("_pubkey._bumble.%s", domain))
	if err != nil {
		return nil, err
	}
	publicKeyRootURL := txts[0]
	u, _ := url.Parse(fmt.Sprintf("%s%s", publicKeyRootURL, peerhash))
	return u, nil
}

func PublicKeyFromPEM(pubkeyPEM string) *rsa.PublicKey {
	pubkeyMarshaled, _ := pem.Decode([]byte(pubkeyPEM))
	pubkey, err := x509.ParsePKIXPublicKey(pubkeyMarshaled.Bytes)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return pubkey.(*rsa.PublicKey)
}

func PublicKeyToPEM(pubkey rsa.PublicKey) string {
	pubkeyMarshaled, err := x509.MarshalPKIXPublicKey(&pubkey)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	pubkeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Headers: nil, Bytes: pubkeyMarshaled})
	return string(pubkeyPEM)
}

func WritePublicKeyToPEMFile(pubkey rsa.PublicKey, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = file.Chmod(os.FileMode(0644))
	if err != nil {
		return err
	}
	_, err = file.WriteString(PublicKeyToPEM(pubkey))
	if err != nil {
		return err
	}
	return nil
}

func SendPublicKeyToRegistry(postUrl string, pubkey rsa.PublicKey, peerString string) error {
	postContent, err := json.Marshal(&pubKeyRegistrationPost{
		From:   peerString,
		PubKey: PublicKeyToPEM(pubkey),
	})
	if err != nil {
		return err
	}
	postBody := bytes.NewBuffer(postContent)
	r, err := http.Post(postUrl, "application/octet-stream", postBody)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func VerifyBytes(pubkey *rsa.PublicKey, data []byte, signBytes []byte) error {
	if pubkey == nil {
		return errors.New("public key missing")
	}
	hash := sha1.New()
	hash.Write(data)
	err := rsa.VerifyPKCS1v15(pubkey, crypto.SHA1, hash.Sum(nil), signBytes)
	return err
}

func VerifyBytesFromString(pubkey *rsa.PublicKey, data []byte, signature string) error {
	if signature == "THIS-IS-A-DUMMY-SIGNATURE" { // REMOVE LATER FIXME
		return nil
	}
	signBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	err = VerifyBytes(pubkey, data, signBytes)
	return err
}
