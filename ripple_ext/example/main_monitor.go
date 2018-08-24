package main

import (
	"github.com/rubblelabs/ripple/ripple_ext/monitor"
	"fmt"
	"flag"
	"os"
	"os/signal"
	"syscall"
	l4g "github.com/alecthomas/log4go"
)

var uris []string = []string{
	//"wss://s-west.ripple.com:443",
	//"wss://s-east.ripple.com:443",
	//"wss://s1.ripple.com:443",
	//"wss://s.altnet.rippletest.net:51233",
	"ws://127.0.0.1:6006",
}
func loop(m *monitor.Monitor) {
	for {
		ledger := <-m.Ledgers()
		fmt.Printf("Ledger %d with %d transactions:\n", ledger.LedgerSequence, len(ledger.Transactions))
		for _, txn := range ledger.Transactions {
			fmt.Printf("  %s %s\n", txn.GetBase().Hash, txn.GetBase().TransactionType)
		}
	}
}

type CCLog struct {
}

func (self *CCLog) Info(arg0 interface{}, args ...interface{}) {
	l4g.Global.Info(arg0, args...)
}

func (self *CCLog) Trace(arg0 interface{}, args ...interface{}) {
	l4g.Global.Trace(arg0, args...)
}

func (self *CCLog) Warn(arg0 interface{}, args ...interface{}) {
	l4g.Global.Warn(arg0, args...)
}

func (self *CCLog) Error(arg0 interface{}, args ...interface{}) {
	l4g.Global.Error(arg0, args...)
}

func (self *CCLog) Debug(arg0 interface{}, args ...interface{}) {
	l4g.Global.Debug(arg0, args...)
}

func main()  {
	flag.Parse()

	l4g.LoadConfiguration("log.xml")
	defer l4g.Close()

	m := monitor.NewMonitor(&CCLog{}, uris, 0)
	go loop(m)

	fmt.Println("Press Ctrl+c to quit...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	l4g.Debug("stop...")
	m.Stop()
	l4g.Debug("stop... ok")
	return
}