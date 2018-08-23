package ripple_api

import (
	"github.com/rubblelabs/ripple/websockets"
	"github.com/rubblelabs/ripple/data"
)

func BuildCancelTx(r * websockets.Remote, srcAddr string, seqNum uint32, fee int64, memo string) (data.Transaction, error) {
	tx := &data.AccountSet{}
	tx.TxBase.TransactionType = data.ACCOUNT_SET

	err := PrepareTx(tx, r, srcAddr, fee, memo)
	if err != nil {
		return nil, err
	}

	// seq
	tx.TxBase.Sequence = seqNum

	// flag
	clearFlag := uint32(5)
	tx.SetFlag = &clearFlag

	return tx, nil
}

func CancelTx(r * websockets.Remote, srcAddr string, seqNum uint32, fee int64, memo string, seed string) (data.Transaction, *websockets.SubmitResult, error) {
	tx, err := BuildCancelTx(r, srcAddr, seqNum, fee, memo)
	if err != nil {
		return nil, nil, err
	}

	// sign
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
