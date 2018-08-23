package ripple_api

import (
	"github.com/rubblelabs/ripple/data"
	"github.com/rubblelabs/ripple/websockets"
)

func BuildPayment(r *websockets.Remote, srcAddr, dstAddr string, amount interface{}, fee int64, memo string) (data.Transaction, error) {
	tx := &data.Payment{}
	tx.TxBase.TransactionType = data.PAYMENT

	// src
	err := PrepareTx(tx, r, srcAddr, fee, memo)
	if err != nil {
		return nil, err
	}

	// dst
	dst, err := data.NewAccountFromAddress(dstAddr)
	if err != nil {
		return nil, err
	}
	tx.Destination = *dst

	// amount
	am, err := data.NewAmount(amount)
	if err != nil {
		return nil, err
	}
	tx.Amount = *am

	// dst tag
	// we put memo to DestinationTag field
	tag, err := BuildTag(memo)
	if err != nil {
		return nil, err
	}
	tx.DestinationTag = tag

	return tx, nil
}

func SendPayment(r *websockets.Remote, srcAddr, dstAddr string, amount interface{}, fee int64, memo string, seed string) (data.Transaction, *websockets.SubmitResult, error) {
	tx, err := BuildPayment(r, srcAddr, dstAddr, amount, fee, memo)
	if err != nil {
		return nil, nil, err
	}

	err = SignTx(tx, seed)
	if err != nil {
		return nil, nil, err
	}

	result, err := SubmitTx(r, tx)
	if err != nil {
		return nil, nil, err
	}

	return tx, result, nil
}
