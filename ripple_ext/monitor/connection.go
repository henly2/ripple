package monitor

import (
	"fmt"
	"github.com/rubblelabs/ripple/data"
	"github.com/rubblelabs/ripple/ripple_ext/logger"
	"github.com/rubblelabs/ripple/websockets"
	"gopkg.in/tomb.v1"
	"time"
)

type Connection struct {
	Ledgers chan *data.Ledger
	Err     error

	conn *websockets.Remote
	t1   tomb.Tomb
	t2   tomb.Tomb

	latestIndex   uint32
	ledgerIndexCh chan uint32

	logger logger.Logger
}

func NewConnection(logger logger.Logger, uri string, beginIndex uint32) (c *Connection, err error) {
	c = &Connection{
		logger:        logger,
		Ledgers:       make(chan *data.Ledger),
		ledgerIndexCh: make(chan uint32, 256),
	}

	// Connect to websocket server
	c.conn, err = websockets.NewRemote(uri)
	if err != nil {
		c.logger.Error("connect ws %s failed:%v", uri, err)
		return
	}

	// Subscribe to ledgers, and server messages
	var confirmation *websockets.SubscribeResult
	confirmation, err = c.conn.Subscribe(true, false, false, true)
	if err != nil {
		c.logger.Error("subscribe err: %v", err)
		return
	}
	c.latestIndex = confirmation.LedgerSequence
	c.logger.Info("subscribe ledger index:%d, basefee:%d", confirmation.LedgerSequence, confirmation.BaseFee)

	if beginIndex == 0 {
		beginIndex = c.latestIndex
	}

	go c.loop1()
	go c.loop2(beginIndex)

	return
}

func (c *Connection) Stop() {
	c.logger.Info("stop connection ...")

	c.Err = nil
	c.t1.Kill(nil)
	c.t2.Kill(nil)
}

func (c *Connection) Wait() {
	c.logger.Info("wait connection close ...")

	c.t1.Wait()
	c.t2.Wait()

	c.conn.Close()
	c.logger.Info("wait connection close fin...")
}

func (c *Connection) loop1() {
	defer c.logger.Info("connection loop1 exit...")
	defer c.t1.Done()
	defer close(c.ledgerIndexCh)

	for {
		// If the tomb is marked dying, exit cleanly
		select {
		case <-c.t1.Dying():
			return
		case msg, ok := <-c.conn.Incoming:
			if !ok {
				c.Err = fmt.Errorf("ws closed")
				c.t1.Kill(nil)
			} else {
				c.handleMessage(msg)
			}
		}
	}
}

func (c *Connection) handleMessage(msg interface{}) {
	switch msg := msg.(type) {
	case *websockets.LedgerStreamMsg:
		c.logger.Debug("LedgerStreamMsg, index=%d, tx count=%d", msg.LedgerSequence, msg.TxnCount)
		c.latestIndex = msg.LedgerSequence
		c.ledgerIndexCh <- msg.LedgerSequence
	case *websockets.TransactionStreamMsg:
		c.logger.Debug("TransactionStreamMsg, ledger index=%d, hash=%s", msg.LedgerSequence, msg.Transaction.GetBase().Hash)
	case *websockets.ServerStreamMsg:
		c.logger.Debug("ServerStreamMsg, txcost =%d, basefee=%d, status=%s\n", msg.TransactionCost(), msg.BaseFee, msg.Status)
	}
}

func (c *Connection) loop2(startIndex uint32) {
	defer c.logger.Info("connection loop2 quit...")
	defer c.t2.Done()
	defer close(c.Ledgers)

	needSleep := true
	index := startIndex
	for {
		needSleep = true

		if c.latestIndex+1 > index {
			cc, err := c.conn.Ledger(index, true)
			if err != nil {
				c.logger.Error("loopLedgerTx index:%d, err:%v", index, err)
			} else {
				if !cc.Ledger.Closed {
					c.logger.Warns("get unclosed ledger %d", cc.Ledger.LedgerSequence)
				} else {
					c.Ledgers <- &cc.Ledger

					index++
					needSleep = false
				}
			}
		}

		// If the tomb is marked dying, exit cleanly
		select {
		case <-c.t2.Dying():
			return
		case _, ok := <-c.ledgerIndexCh:
			if !ok {
				c.t2.Kill(nil)
				return
			}
			needSleep = false
		default:
		}

		// sleep
		if needSleep {
			c.logger.Info("connection loop2 sleep 5 seconds...")
			time.Sleep(time.Second * 5)
		}
	}
}
