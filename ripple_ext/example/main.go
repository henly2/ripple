package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/rubblelabs/ripple/data"
	"github.com/rubblelabs/ripple/terminal"
	"github.com/rubblelabs/ripple/websockets"
	"time"
	"encoding/json"
	"github.com/rubblelabs/ripple/ripple_ext/ripple_api"
)

func checkErr(err error, quit bool) {
	if err != nil {
		terminal.Println(err.Error(), terminal.Default)
		if quit {
			os.Exit(1)
		}
	}
}

func myPrint(a interface{})  {
	b, _ := json.MarshalIndent(a, " ", "")
	fmt.Println(string(b))
}

var (
	localWs = "ws://127.0.0.1:6006"
	testWs = "wss://s.altnet.rippletest.net:51233"
	pubWs = "wss://s-east.ripple.com:443"

	host     = flag.String("host", "local", "websockets host to connect to")
	proposed = flag.Bool("proposed", false, "include proposed transactions")

	r *websockets.Remote

	local_genesis = ripple_api.RippledKeyPair{
		Address:"rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
		Seed:"snoPBrXtMeMyMHUVTgbuqAfg1SUTb",
	}

	local_addr1 = ripple_api.RippledKeyPair{
		Address:"rEiVPf31YhtgJtq8fvi1cpwNKDEaAp25Ud",
		Seed:"ssLUvvnEDgEjGTJorSB5aMscebBbt",
	}

	local_addr2_regular = ripple_api.RippledKeyPair{
		Address:"rBCSgao1iBP412bpFrQgfVnj7HGEeJUThr",
		Seed:"sh6JW3F9EtJ7iA8ZHLg85bcDVqFhy",
	}

	local_addr2 = ripple_api.RippledKeyPair{
		Address:"rntihVJZdoWTchfqgKRFGBK2byUUQ6iXze",
		Seed:"shEpDkHDQpPFPWsmgzCLMzaBWVF92",
	}

	local_addr3 = ripple_api.RippledKeyPair{
		Address:"r4cdWbeU1evVwFoVR3oz3ygJD5cK8QFVbm",
		Seed:"spFR1JxsaawVSuRpyD1Q84xhMHFoy",
	}

	test_addr = ripple_api.RippledKeyPair{
		Address:"rfkrt95xB744ghAv8yg9P8J2FB9Rg7JnpQ",
		Seed:"snR1mzgbAXNdZNYPdmiRhqdJBPSr5",
	}

	balancesIn = make(map[string]*data.Amount)
	balancesOut = make(map[string]*data.Amount)
)

func test4()  {
	ch, _ := ripple_api.StreamLedgers(r, 40930633, true)
	for {
		select {
		case ledger, ok := <- ch:
			if !ok {
				goto quit
			}
			parseLedger(ledger)
		}
	}
quit:
	fmt.Println("test4 quit")

	fmt.Println("in:")
	for k, v := range balancesIn {
		fmt.Println("%s-%v", k, v)
	}
	fmt.Println("out:")
	for k, v := range balancesOut {
		fmt.Println("%s-%v", k, v)
	}
}

func printGetBalance()  {
	fmt.Println("\nbalance...")
	a, err := ripple_api.GetAccountInfo(r, local_genesis.Address)
	fmt.Println(local_genesis.Address, ":")
	if err == nil {
		myPrint(a.AccountData.Balance.String())
	} else {
		fmt.Println(err)
	}

	a, err = ripple_api.GetAccountInfo(r, local_addr1.Address)
	fmt.Println(local_addr1.Address, ":")
	if err == nil {
		myPrint(a.AccountData.Balance.String())
	}else {
		fmt.Println(err)
	}

	a, err = ripple_api.GetAccountInfo(r, local_addr2.Address)
	fmt.Println(local_addr2.Address, ":")
	if err == nil {
		myPrint(a.AccountData.Balance.String())
	}else {
		fmt.Println(err)
	}

	a, err = ripple_api.GetAccountInfo(r, local_addr3.Address)
	fmt.Println(local_addr3.Address, ":")
	if err == nil {
		myPrint(a.AccountData.Balance.String())
	}else {
		fmt.Println(err)
	}

	a, err = ripple_api.GetAccountInfo(r, test_addr.Address)
	fmt.Println(test_addr.Address, ":")
	if err == nil {
		myPrint(a.AccountData.Balance.String())
	}else {
		fmt.Println(err)
	}
	fmt.Println("balance end...\n")
}


func main() {
	var err error

	flag.Parse()

	if *host == "test" {
		*host = testWs
	} else if *host == "pub" {
		*host = pubWs
	} else {
		*host = localWs
	}

	r, err = websockets.NewRemote(*host)
	checkErr(err, true)

	//test4()
	lr, err := r.Ledger(nil, false)
	if err != nil {
		fmt.Println("GetLedgerInfo err:", err)
	} else {
		fmt.Println("ledger:", lr.Ledger.LedgerSequence)
	}

	confirmation, err := r.Subscribe(true, true, false, true)
	checkErr(err, true)
	terminal.Println(fmt.Sprint("Subscribed at: ", confirmation.LedgerSequence), terminal.Default)

	printGetBalance()
	//return

	// Consume messages as they arrive
	go streamMessage()

	for {
		time.Sleep(time.Second*1)

		////_, re, err := ripple_api.SendPayment(r, local_genesis.Address, local_addr3.Address, 100000000, 10, "local3", local_addr2_regular.Seed)
		//tx, re, err := ripple_api.SendPayment(r, local_addr3.Address, local_addr2.Address, 100000000, 10, "local2", local_addr3.Seed)
		//if err != nil {
		//	fmt.Println("SendPayment err:", err)
		//	break
		//}
		//
		//fmt.Println("SendPayment result code:", re.EngineResultCode)
		//if re.EngineResultCode != 0 {
		//	fmt.Println("SendPayment result msg:", re.EngineResult.String(), re.EngineResultMessage)
		//
		//	_, re2, err2 := ripple_api.CancelTx(r, local_addr3.Address, tx.GetBase().Sequence, 10,"", local_addr3.Seed)
		//	if err2 != nil {
		//		fmt.Println("SendPayment err:", err2)
		//	} else {
		//		fmt.Println("SendPayment result code:", re2.EngineResultCode)
		//		if re2.EngineResultCode != 0 {
		//			fmt.Println("SendPayment result msg:", re2.EngineResult.String(), re2.EngineResultMessage)
		//		}
		//	}
		//
		//	break
		//}
		//sendPayment(local_genesis.Address, local_addr2.Address, 500000000, 10, local_addr2_regular.Seed)


		//setRegularKey(local_genesis.Address, local_genesis.Seed, local_addr2_regular.Address)
		break
	}

	fmt.Println("sleep 10 seconds...")
	time.Sleep(time.Second*100)
	//printGetBalance()

	fmt.Println("quit")
}

func streamMessage()  {
	for {
		msg, ok := <-r.Incoming
		if !ok {
			return
		}

		switch msg := msg.(type) {
		case *websockets.LedgerStreamMsg:
			//fmt.Println("LedgerStreamMsg===")
			//terminal.Println(msg, terminal.Default)
			//fmt.Println("LedgerStreamMsg===end")
		case *websockets.TransactionStreamMsg:
			parseTx(&msg.Transaction)
			//fmt.Println("TransactionStreamMsg===")
			//terminal.Println(&msg.Transaction, terminal.Indent)
			//
			//fmt.Println("TransactionStreamMsg-PathSet===")
			//for _, path := range msg.Transaction.PathSet() {
			//	terminal.Println(path, terminal.DoubleIndent)
			//}
			//
			//fmt.Println("TransactionStreamMsg-trades===")
			//trades, err := data.NewTradeSlice(&msg.Transaction)
			//checkErr(err, false)
			//for _, trade := range trades {
			//	terminal.Println(trade, terminal.DoubleIndent)
			//}
			//
			//fmt.Println("TransactionStreamMsg-balance===")
			//balances, err := msg.Transaction.Balances()
			//checkErr(err, false)
			//for _, balance := range balances {
			//	terminal.Println(balance, terminal.DoubleIndent)
			//
			//	myPrint(balance)
			//}
			//fmt.Println("TransactionStreamMsg===")
		case *websockets.ServerStreamMsg:
			//terminal.Println(msg, terminal.Default)
		}
	}
}

func parseLedger(ledger *websockets.LedgerResult)  {
	for _, t := range ledger.Ledger.Transactions {
		t.LedgerSequence = ledger.Ledger.LedgerSequence
		parseTx(t)
	}
}

func parseTx(v *data.TransactionWithMetaData)  {
	var (
		base   = v.GetBase()
		format = "sign:%s type:%-11s fee:%-8s memo:%s account:%-34s seq:%-9d "
		values = []interface{}{terminal.SignSymbol(v), base.GetType(), base.Fee, terminal.MemoSymbol(v), base.Account, base.Sequence}
	)

	ToAmount(balancesOut, base.Account.String(), base.Fee)

	switch tx := v.Transaction.(type) {
	case *data.Payment:
		format += "=> dst:%-34s amount:%-60s sendMax:%-60s"
		values = append(values, []interface{}{tx.Destination, tx.Amount, tx.SendMax}...)
		if v.MetaData.TransactionResult.Success(){
			ToAmount2(balancesOut, base.Account.String(), tx.Amount)
			ToAmount2(balancesIn, tx.Destination.String(), tx.Amount)
		}
		if tx.DestinationTag != nil {
			values = append(values, *tx.DestinationTag)
		} else {
			values = append(values, "nil")
		}
	case *data.SetRegularKey:
		format += "%-9d"
		values = append(values, tx.Sequence)
	case *data.AccountSet:
		format += "%-9d"
		values = append(values, tx.Sequence)
	case *data.TrustSet:
		format += "%-60s %d %d"
		values = append(values, tx.LimitAmount, tx.QualityIn, tx.QualityOut)
	default:
	}

	values = append(values, v.MetaData.TransactionResult.String())

	fmt.Println(values)

	for _, m := range base.Memos {
		fmt.Println(string(m.Memo.MemoType), string(m.Memo.MemoFormat), string(m.Memo.MemoData))
	}
}

func ToAmount(b map[string]*data.Amount, address string, v data.Value) {
	var err error

	am, ok := b[address];
	if !ok {
		b[address] = &data.Amount{
			Value: v.Clone(),
		}
		return
	}

	am1 := &data.Amount{
		Value: v.Clone(),
	}

	b[address], err = am.Add(am1)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func ToAmount2(b map[string]*data.Amount, address string, a data.Amount) {
	var err error

	am, ok := b[address];
	if !ok {
		b[address] = a.Clone()
		return
	}

	am1 := a.Clone()

	b[address], err = am.Add(am1)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

