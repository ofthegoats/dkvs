package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"io"
	"log"
	"time"

	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/rep"
	"go.nanomsg.org/mangos/v3/protocol/req"
)

// Reasonable default subject to change
const (
	TIMEOUT = 1 * time.Second
)

func EncryptAES(key, plaintext []byte) []byte {
	c, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(c)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce) // fill nonce with cryptographically secure random sequence
	return gcm.Seal(nonce, nonce, plaintext, nil)
}

func DecryptAES(key, ciphertext []byte) []byte {
	c, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(c)
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, _ := gcm.Open(nil, nonce, ciphertext, nil)
	return plaintext
}

// Set up a listening port and take the first message that reaches it, close after this.
func Listen(socket string, messages chan<- Rumour, key []byte) error {
	lSocket, err := rep.NewSocket() // listenSocket
	defer lSocket.Close()           // Close when finished
	if err != nil {                 // If the node fails to establish a socket to listen on, this is fatal
		log.Fatalln(err)
	}
	lSocket.Listen(socket)
	encrypted_bs, err := lSocket.Recv()
	if err != nil { // failed to read bytes, maybe not be fatal, so don't panic
		log.Println(err)
	}
	bs := DecryptAES(key, encrypted_bs)
	// Decode gob structure and send it on messages channel
	reader := bytes.NewReader(bs)
	decoder := gob.NewDecoder(reader)
	var msg Rumour
	err = decoder.Decode(&msg)
	if err != nil { // if decoding failed, decryption must be wrong, implying message was not for this node
		return err
	}
	messages <- msg
	return nil
}

// Send a Rumour to a socket (recipient), AES/GCM encrypted with key
func Send(recipient string, key []byte, rumour Rumour) error {
	sSocket, err := req.NewSocket() // sendSocket
	if err != nil {                 // failed to establish a socket
		log.Println(err) // not fatal for whole program, but does mean this method failed.
		return err
	}
	sSocket.SetOption(mangos.OptionSendDeadline, TIMEOUT)
	sSocket.Dial(recipient)
	// turn the message into bytes to be sent over the network as gob data
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err = encoder.Encode(rumour)
	if err != nil { // failed to encode, do not send
		log.Println(err)
		return err
	}
	encrypted_buffer := EncryptAES(key, buffer.Bytes())
	err = sSocket.Send(encrypted_buffer)
	if err != nil { // sending the bytes failed
		log.Println(err) // not fatal for whole program, but does mean this method failed.
		return err
	}
	return nil
}
