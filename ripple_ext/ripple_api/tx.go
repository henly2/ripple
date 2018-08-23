package ripple_api

import (
	"github.com/rubblelabs/ripple/data"
	"fmt"
	"github.com/rubblelabs/ripple/websockets"
)

func SignTx(tx data.Signer, seed string) error {
	// get key
	key, err := GetKeyBySeed(seed)
	if err != nil {
		return err
	}

	// sign
	var sequence uint32
	err = data.Sign(tx, key, &sequence)
	if err != nil {
		return err
	}

	// verify
	ok, err := data.CheckSignature(tx)
	if err != nil {
		return err
	}
	if !ok {
		err = fmt.Errorf("check sig is false")
		return err
	}

	return nil
}

func SubmitTx(r *websockets.Remote, tx data.Transaction) (*websockets.SubmitResult, error) {
	result, err := r.Submit(tx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func PrepareTx(tx data.Transaction, r *websockets.Remote, srcAddr string, fee int64, memo string) error {
	base := tx.GetBase()

	// src
	src, err := data.NewAccountFromAddress(srcAddr)
	if err != nil {
		return err
	}
	base.Account = *src

	//seq
	srcAccount, err := GetAccountInfo(r, srcAddr)
	if err != nil {
		return err
	}
	base.Sequence = *srcAccount.AccountData.Sequence

	//LastLedgerSequence
	//base.LastLedgerSequence = new(uint32)
	//*base.LastLedgerSequence = srcAccount.LedgerSequence + 2

	// fee
	f, err := data.NewNativeValue(fee)
	if err != nil {
		return err
	}
	base.Fee = *f

	// memo
	if memo != ""{
		mm := data.Memo{
			Memo: struct {
				MemoType   data.VariableLength
				MemoData   data.VariableLength
				MemoFormat data.VariableLength
			}{MemoType: data.VariableLength("test"), MemoData: data.VariableLength(memo), MemoFormat: data.VariableLength("plain/text")},
		}
		base.Memos = append(base.Memos, mm)
	}

	return nil
}
