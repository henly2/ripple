package ripple_api

import (
	"github.com/rubblelabs/ripple/websockets"
)

func StreamLedgers(r *websockets.Remote, ledgerStartIndex uint32, transactions bool) (chan *websockets.LedgerResult, error) {
	ch := make(chan *websockets.LedgerResult)

	//confirmation, err := r.Subscribe(true, true, false, true)
	//if err != nil {
	//	return nil, err
	//}
	//
	//for {
	//	msg, ok := <-r.Incoming
	//	if !ok {
	//		return
	//	}
	//
	//	switch msg := msg.(type) {
	//	case *websockets.LedgerStreamMsg:
	//		// do nothing now
	//	case *websockets.TransactionStreamMsg:
	//		//fmt.Println("TransactionStreamMsg===")
	//		//terminal.Println(&msg.Transaction, terminal.Indent)
	//		//
	//		//fmt.Println("TransactionStreamMsg-PathSet===")
	//		//for _, path := range msg.Transaction.PathSet() {
	//		//	terminal.Println(path, terminal.DoubleIndent)
	//		//}
	//		//
	//		//fmt.Println("TransactionStreamMsg-trades===")
	//		//trades, err := data.NewTradeSlice(&msg.Transaction)
	//		//checkErr(err, false)
	//		//for _, trade := range trades {
	//		//	terminal.Println(trade, terminal.DoubleIndent)
	//		//}
	//		//
	//		//fmt.Println("TransactionStreamMsg-balance===")
	//		//balances, err := msg.Transaction.Balances()
	//		//checkErr(err, false)
	//		//for _, balance := range balances {
	//		//	terminal.Println(balance, terminal.DoubleIndent)
	//		//
	//		//	myPrint(balance)
	//		//}
	//		//fmt.Println("TransactionStreamMsg===")
	//	case *websockets.ServerStreamMsg:
	//		// do nothing now
	//	}
	//}
	//
	//// get old ledger
	//go func(ledgerStartIndex uint32, transactions bool, ch chan *websockets.LedgerResult) {
	//	defer close(ch)
	//
	//	var ledgerNextIndex = ledgerStartIndex
	//
	//	for {
	//		cc, err := r.Ledger(ledgerNextIndex, transactions)
	//		if err != nil {
	//			return
	//		}
	//
	//		ch <- cc
	//		if !cc.Ledger.Closed {
	//			return
	//		}
	//
	//		ledgerNextIndex = cc.Ledger.Ledger() + 1
	//	}
	//}(ledgerStartIndex, transactions, ch)

	// then subscriber new

	return ch, nil
}