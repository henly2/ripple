package monitor

import (
	"github.com/rubblelabs/ripple/data"
	"github.com/rubblelabs/ripple/ripple_ext/logger"
	"gopkg.in/tomb.v1"
	"time"
	"github.com/rubblelabs/ripple/websockets"
	"sync"
)

type Monitor struct {
	t       tomb.Tomb
	ledgers chan *data.Ledger

	uris          []string
	scanLedgerIdx uint32

	logger logger.Logger

	ledgerStream websockets.LedgerStreamMsg
	locker sync.Mutex
}

func NewMonitor(logger logger.Logger, uris []string, startLedgerIdx uint32) *Monitor {
	m := &Monitor{
		logger:        logger,
		ledgers:       make(chan *data.Ledger),
		scanLedgerIdx: startLedgerIdx,
		uris:          uris,
	}

	go m.loop()
	return m
}

func (m *Monitor) Stop() {
	m.logger.Info("monitor stop")

	m.t.Kill(nil)
	m.t.Wait()

	m.logger.Info("monitor stop fin")
}

func (m *Monitor) Ledgers() chan *data.Ledger {
	return m.ledgers
}

func (m *Monitor) LedgerStream() websockets.LedgerStreamMsg {
	m.locker.Lock()
	defer m.locker.Unlock()
	return m.ledgerStream
}

func (m *Monitor) SetLedgerStream(msg *websockets.LedgerStreamMsg) {
	m.locker.Lock()
	defer m.locker.Unlock()
	if msg != nil {
		m.ledgerStream = *msg
	}
}

func (m *Monitor) loop() {
	uriIndex := 0
	defer m.logger.Info("monitor loop quit...")
	defer m.t.Done()

	for {
		m.logger.Info("monitor connecting: %s", m.uris[uriIndex])
		err := m.handleConnection(m.uris[uriIndex])
		if err != nil {
			m.logger.Error("monitor connection failed: %v", err)
		}

		uriIndex = (uriIndex + 1) % len(m.uris)

		select {
		case <-m.t.Dying():
			// If the tomb is marked dying, exit cleanly
			return
		default:
			// Wait a bit before reconnecting to another server
			time.Sleep(5 * time.Second)
		}
	}
}

func (m *Monitor) handleConnection(uri string) (err error) {
	var (
		c    *Connection
		stop bool
	)

	// Open a new connection to ripple
	c, err = NewConnection(m.logger, uri, m.scanLedgerIdx)
	if err != nil {
		return err
	}
	defer c.Wait()

	m.SetLedgerStream(c.SubResult().LedgerStreamMsg)

	stop = false
	for {
		select {
		case <-m.t.Dying():
			// We are exiting cleanly
			if stop == false {
				stop = true
				c.Stop()
			}
		case d, ok := <-c.Data:
			if !ok {
				// Ripple connection is dead
				return c.Err
			}
			if ledger, ok := d.(*data.Ledger); ok {
				m.ledgers <- ledger
				m.scanLedgerIdx = ledger.LedgerSequence
			} else if msg, ok := d.(*websockets.LedgerStreamMsg); ok {
				m.SetLedgerStream(msg)
			}
		}
	}
}
