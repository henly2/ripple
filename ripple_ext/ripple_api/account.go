package ripple_api

import (
	"crypto/rand"
	"github.com/rubblelabs/ripple/crypto"
	"github.com/rubblelabs/ripple/data"
	"github.com/rubblelabs/ripple/websockets"
	"fmt"
)

type RippledKeyPair struct {
	Address string
	Seed 	string
}

func NewAccount() (*RippledKeyPair, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	s, err := crypto.NewFamilySeed(b)
	if err != nil {
		return nil, err
	}

	addr, err := GetAddressBySeed(s.String())
	if err != nil {
		return nil, err
	}

	kp := &RippledKeyPair{
		Address: addr,
		Seed:s.String(),
	}

	return kp, nil
}

func GetKeyBySeed(seed string) (crypto.Key, error) {
	s, err := data.NewSeedFromAddress(seed)
	if err != nil {
		return nil, nil
	}

	return s.Key(KeyTypeECDSA), nil
}

func GetAddressBySeed(seed string) (string, error) {
	s, err := data.NewSeedFromAddress(seed)
	if err != nil {
		return "", nil
	}

	var sequence uint32
	a := s.AccountId(KeyTypeECDSA, &sequence)
	return a.String(), nil
}

func GetAccountInfo(r *websockets.Remote, address string) (*websockets.AccountInfoResult, error) {
	src, err := data.NewAccountFromAddress(address)
	if err != nil {
		return nil, err
	}

	ai, err := r.AccountInfo(*src)
	if err != nil {
		return nil, err
	}

	return ai, nil
}

func GetAccountSeqNum(r *websockets.Remote, address string) (uint32, error) {
	a, err := GetAccountInfo(r, address)
	if err != nil {
		return 0, err
	}

	if a.AccountData.Sequence == nil {
		return 0, fmt.Errorf("sequence field is nil")
	}

	return *a.AccountData.Sequence, nil
}

func IsValidAccount(addressOrSeed string) (bool) {
	// seed?
	_, err := data.NewSeedFromAddress(addressOrSeed)
	if err == nil {
		return true
	}

	// address
	_, err = data.NewAccountFromAddress(addressOrSeed)
	return err == nil
}