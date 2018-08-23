package ripple_api

import (
	"github.com/rubblelabs/ripple/data"
	"github.com/rubblelabs/ripple/websockets"
)

func BuildSetRegularKey(r *websockets.Remote, srcAddr, regularAddr string, fee int64, memo string) (data.Transaction, error) {
	tx := &data.SetRegularKey{}
	tx.TxBase.TransactionType = data.SET_REGULAR_KEY

	// src
	err := PrepareTx(tx, r, srcAddr, fee, memo)
	if err != nil {
		return nil, err
	}

	// set or reset
	if regularAddr != ""{
		// add regular
		dst, err := data.NewRegularKeyFromAddress(regularAddr)
		if err != nil {
			return nil, err
		}
		tx.RegularKey = dst
	} else {
		// remove all
		tx.RegularKey = nil
	}

	return tx, err
}

func SetRegularKey(r *websockets.Remote, srcAddr, regularAddr string, fee int64, memo string, seed string) (data.Transaction, *websockets.SubmitResult, error) {
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
