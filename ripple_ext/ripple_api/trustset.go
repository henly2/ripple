package ripple_api

import (
	"github.com/rubblelabs/ripple/data"
	"github.com/rubblelabs/ripple/websockets"
)

func BuildTrustSet(r *websockets.Remote, srcAddr string, currency, issuer string, fee int64, memo string) (data.Transaction, error) {
	tx := &data.TrustSet{}
	tx.TxBase.TransactionType = data.TRUST_SET

	// src
	err := PrepareTx(tx, r, srcAddr, fee, memo)
	if err != nil {
		return nil, err
	}

	la, err := data.NewAmount("")
	if err != nil {
		return nil, err
	}
	tx.LimitAmount = *la

	return tx, err
}

func TrustSet(r *websockets.Remote, srcAddr, regularAddr string, fee int64, memo string, seed string) (data.Transaction, *websockets.SubmitResult, error) {
	// build
	tx, err := BuildSetRegularKey(r, srcAddr, regularAddr, fee, memo)
	if err != nil {
		return nil, nil, err
	}

	// sign
	err = SignTx(tx, seed)
	if err != nil {
		return nil, nil, err
	}

	// submit
	result, err := SubmitTx(r, tx)
	if err != nil {
		return nil, nil, err
	}

	return tx, result, nil
}
