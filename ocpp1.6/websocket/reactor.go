package websocket

import (
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"runtime"
	"sync/atomic"
)

var (
	ErrReactorShutdown = errors.New("reactor is going to be shutdown")
)

const (
	ErrEvents = unix.EPOLLERR | unix.EPOLLHUP | unix.EPOLLRDHUP
	OutEvents = ErrEvents | unix.EPOLLOUT
	InEvents  = ErrEvents | unix.EPOLLIN | unix.EPOLLPRI
)

type reactor struct {
	server       *Server
	index        int
	connections  *wsconns
	connectCount int32
	epoller      *epoller
}

func (r *reactor) activity() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer func() {
		r.stopSignal()
		r.deleteAllConnections()
	}()
	r.epoller.reactor(func(fd int, ev uint32) error {
		var err error
		if ws, ok := r.connections.getConnByFD(fd); ok {
			if ev&OutEvents != 0 {
				//TODO() support epollout
			}
			if ev&InEvents != 0 {
				if err = ws.read(); err != nil {
					r.epoller.delete(fd)
					r.connections.deleteConn(ws.id, ws.fd)
					r.addConnCount(-1)
					ws.Lock()
					ws.stop(err)
					ws.Unlock()
				}
			}
		}
		return err
	})
}

func (r *reactor) deleteAllConnections() {
	r.connections.Lock()
	for fd, wsconn := range r.connections.wsfdmap {
		r.epoller.delete(fd)
		wsconn.stop(fmt.Errorf("reactor%d shut down", r.index))
		delete(r.connections.wsmap, wsconn.id)
		delete(r.connections.wsfdmap, wsconn.fd)
		r.addConnCount(-1)
	}
	r.connections.Unlock()
	r.epoller.close()
}

func (r *reactor) stopSignal() {
	r.server.once.Do(func() {
		r.server.cond.L.Lock()
		r.server.cond.Signal()
		r.server.cond.L.Unlock()
	})
}

func (r *reactor) addConnCount(delta int32) {
	atomic.AddInt32(&r.connectCount, delta)
}

func (r *reactor) registerConn(conn interface{}) error {
	c := conn.(*Wsconn)
	if err := r.epoller.addRead(c.fd); err != nil {
		c.conn.Close()
		return err
	}
	c.server.clientOnConnect(c)
	r.connections.registerConn(c.id, c.fd, c)
	r.addConnCount(1)
	go c.writedump()
	return nil
}
