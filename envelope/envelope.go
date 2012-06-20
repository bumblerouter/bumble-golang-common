package envelope

import (
	"bumbleserver.org/common/key"
	"bumbleserver.org/common/message"
	"bumbleserver.org/common/peer"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type Envelope struct {
	From      *peer.Peer `json:"f,omitempty"`
	To        *peer.Peer `json:"t,omitempty"`
	Message   string     `json:"m,omitempty"`
	Signature string     `json:"s,omitempty"`
	KeyBundle string     `json:"k,omitempty"`
}

func (e *Envelope) GetFrom() *peer.Peer {
	return e.From
}

func (e *Envelope) GetTo() *peer.Peer {
	return e.To
}

func (e *Envelope) GetMessageRaw() string {
	return e.Message
}

func (e *Envelope) GetMessage(prikey *rsa.PrivateKey) string {
	if len(e.KeyBundle) == 0 {
		fmt.Printf("MESSAGE FROM %s WAS NOT ENCRYPTED.\n", e.From)
		return e.Message
	}
	fmt.Printf("MESSAGE FROM %s WAS ENCRYPTED.\n", e.From)

	encryptedKeyIvBundle, err := base64.StdEncoding.DecodeString(e.KeyBundle)
	if err != nil {
		fmt.Printf("UNABLE TO DECODE KEYBUNDLE FROM BASE64 ENCODING: %s\n", err)
		return "{}"
	}

	keyIvBundle, err := rsa.DecryptPKCS1v15(rand.Reader, prikey, encryptedKeyIvBundle)
	if err != nil {
		fmt.Printf("UNABLE TO DECRYPT KEY FROM ENVELOPE: %s\n", err)
		return "{}"
	}

	iv := keyIvBundle[0:16]
	key := keyIvBundle[16:]

	encryptedMessage, err := base64.StdEncoding.DecodeString(e.Message)
	if err != nil {
		fmt.Printf("UNABLE TO DECODE MESSAGE FROM BASE64 ENCODING: %s\n", err)
		return "{}"
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("UNABLE TO USE DECRYPTED KEY FROM ENVELOPE TO PRODUCE CIPHER: %s\n", err)
	}

	cbc := cipher.NewCBCDecrypter(c, iv)
	decryptedMessage := make([]byte, len(encryptedMessage))
	cbc.CryptBlocks(decryptedMessage, encryptedMessage)
	padlen := decryptedMessage[len(decryptedMessage)-1]
	decryptedMessage = bytes.Replace(decryptedMessage, bytes.Repeat([]byte{padlen}, int(padlen)), nil, -1) // FIXME can this be better?
	return string(decryptedMessage)
}

func (e *Envelope) GetSignature() string {
	return e.Signature
}

func Package(msg message.Message, prikey *rsa.PrivateKey) (*Envelope, error) {
	env := new(Envelope)
	env.From = msg.GetFrom()
	env.To = msg.GetTo()
	if msg.GetDate().IsZero() {
		msg.SetDate(time.Now().UTC())
	}
	msgJson, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("PACKAGEENVELOPE ERROR: %s\n", err.Error())
		return nil, err
	}
	if prikey != nil {
		var pubkey *rsa.PublicKey
		if msg.GetTo() != nil {
			pubkey = msg.GetTo().PublicKey()
		}
		if pubkey == nil { // when is it okay to be nil?  probably never!  FIXME
			env.Message = string(msgJson)
		} else {
			keyIvBundle := make([]byte, 48)
			rand.Read(keyIvBundle)
			encryptedKeyIvBundle, err := rsa.EncryptPKCS1v15(rand.Reader, pubkey, keyIvBundle)
			if err != nil {
				fmt.Printf("UNABLE TO ENCRYPT KEY FOR ENVELOPE: %s\n", err)
			}

			iv := keyIvBundle[0:16]
			key := keyIvBundle[16:48]

			c, err := aes.NewCipher(key)
			if err != nil {
				fmt.Printf("UNABLE TO USE KEY TO PRODUCE CIPHER FOR ENVELOPE: %s\n", err)
			}

			cbc := cipher.NewCBCEncrypter(c, iv)
			padlen := c.BlockSize() - len(msgJson)%c.BlockSize()
			encryptedMessage := make([]byte, len(msgJson)+padlen)
			paddedMessage := bytes.Join([][]byte{msgJson, bytes.Repeat([]byte{byte(padlen)}, padlen)}, nil)
			cbc.CryptBlocks(encryptedMessage, paddedMessage)

			env.KeyBundle = base64.StdEncoding.EncodeToString(encryptedKeyIvBundle)
			env.Message = base64.StdEncoding.EncodeToString(encryptedMessage)
		}
		signature, err := key.SignBytesToString(prikey, []byte(env.Message))
		if err != nil {
			fmt.Printf("PACKAGEENVELOPE-SIGN ERROR: %s\n", err.Error())
			return nil, err
		} else {
			env.Signature = string(signature)
		}
	}

	out, _ := json.Marshal(env)
	fmt.Println(string(out))

	return env, nil
}
