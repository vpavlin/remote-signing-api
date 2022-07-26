package signer

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
)

type SignedMessage struct {
	Address   string
	Hash      []byte
	Signature []byte
}

func IsValidSignature(signer string, hash []byte, signature []byte) (bool, error) {
	sigPubKey, err := crypto.Ecrecover(hash, signature)
	if err != nil {
		return false, err
	}

	logrus.Debugf("Recovered: %s", sigPubKey)

	hashAddress := crypto.Keccak256Hash(sigPubKey[1:])
	recoveredAddrBytes := hashAddress.Bytes()[12:]

	recoveredAddr := common.HexToAddress(hex.EncodeToString(recoveredAddrBytes))

	logrus.Debugf("Recovered address: %s, used address: %s, match: %t", recoveredAddr, signer, recoveredAddr == common.HexToAddress(signer))

	return recoveredAddr == common.HexToAddress(signer), nil
}

func BytesToString(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func StringToBytes(data string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}
