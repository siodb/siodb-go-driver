// Copyright (C) 2019-2020 Siodb GmbH. All rights reserved.
// Use of this source code is governed by a license that can be found
// in the LICENSE file.

package siodb

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"io/ioutil"
)

func (sc *siodbConn) authenticate() (err error) {

	var buf [binary.MaxVarintLen32]byte
	var encodedLength int

	// Begin Session Request
	beginSessionRequest := &BeginSessionRequest{
		UserName: sc.cfg.user,
	}
	sc.debug("authenticate | %v", beginSessionRequest)
	encodedLength = binary.PutUvarint(buf[:], uint64(5))
	sc.netConn.Write(buf[:encodedLength])
	writeMessage(sc.netConn, beginSessionRequest)

	// Get Session Response
	var beginSessionResponse BeginSessionResponse
	if _, err = sc.ReadMessage(6, &beginSessionResponse); err != nil {
		return err
	}

	// Check if Siodb has started the session
	if !beginSessionResponse.GetSessionStarted() {
		return &siodbDriverError{"Starting session failed: " + beginSessionResponse.GetMessage().GetText()}
	}
	sc.debug("authenticate | beginSessionResponse | %v", sc.cfg.privateKey)

	// Sign challenge sha512 digest
	sha512 := sha512.New()
	sha512.Write(beginSessionResponse.GetChallenge())
	signature, err := rsa.SignPKCS1v15(nil, sc.cfg.privateKey, crypto.SHA512, sha512.Sum(nil))
	sc.debug("authenticate | signature | %v", err)

	// Begin Session Request
	clientAuthenticationRequest := &ClientAuthenticationRequest{
		Signature: signature,
	}
	sc.debug("authenticate | clientAuthenticationRequest | %v", clientAuthenticationRequest)
	encodedLength = binary.PutUvarint(buf[:], uint64(7))
	sc.netConn.Write(buf[:encodedLength])
	writeMessage(sc.netConn, clientAuthenticationRequest)

	// Get Session Response
	var clientAuthenticationResponse ClientAuthenticationResponse
	if _, err = sc.ReadMessage(8, &clientAuthenticationResponse); err != nil {
		return err
	}
	sc.debug("authenticate | clientAuthenticationResponse | %v", clientAuthenticationResponse)

	// Check if Siodb has authenticated the session
	if !clientAuthenticationResponse.GetAuthenticated() {
		return &siodbDriverError{"Authentication failed: " + clientAuthenticationResponse.GetMessage().GetText()}
	}
	sc.debug("authenticate | %v", clientAuthenticationResponse)

	// Setup session Id
	sc.sessionID = clientAuthenticationResponse.GetSessionId()

	return nil
}

func loadPrivateKey(rsaPKeyPath string, rsaPKeyPwd string) (pk *rsa.PrivateKey, err error) {

	priv, err := ioutil.ReadFile(rsaPKeyPath)
	if err != nil {
		return pk, &siodbDriverError{"Paring URI: Indentity file '" + rsaPKeyPath + "' not found."}
	}

	privatePem, _ := pem.Decode(priv)
	var privatePemBytes []byte
	if privatePem.Type != "RSA PRIVATE KEY" {
		return pk, &siodbDriverError{"RSA private key is of the wrong type."}
	}

	if rsaPKeyPwd != "" {
		privatePemBytes, err = x509.DecryptPEMBlock(privatePem, []byte(rsaPKeyPwd))
	} else {
		privatePemBytes = privatePem.Bytes
	}

	var parsedPrivateKey interface{}
	if parsedPrivateKey, err = x509.ParsePKCS1PrivateKey(privatePemBytes); err != nil {
		if parsedPrivateKey, err = x509.ParsePKCS8PrivateKey(privatePemBytes); err != nil {
			return pk, &siodbDriverError{"Unable to parse RSA private key."}
		}
	}

	var ok bool
	pk, ok = parsedPrivateKey.(*rsa.PrivateKey)
	if !ok {
		return pk, &siodbDriverError{"Unable to parse RSA private key"}
	}

	return pk, nil
}
