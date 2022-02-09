package websocket

import (
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"runtime"
	"sync/atomic"
	"time"
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
				t := time.Now().Unix()
				fmt.Println("TestRead begin", fd, ws.id, t)
				if err = ws.read(); err != nil {
					fmt.Printf("Testdelete1, index=%v, fd=%v, id=%v\n", r.index, fd, ws.id)
					r.epoller.delete(fd)
					fmt.Printf("Testdelete2, index=%v, fd=%v, id=%v\n", r.index, fd, ws.id)
					r.connections.deleteConn(ws.id, ws.fd)
					fmt.Printf("Testdelete3, index=%v, fd=%v, id=%v\n", r.index, fd, ws.id)
					r.addConnCount(-1)
					fmt.Printf("Testdelete4, index=%v, fd=%v, id=%v\n", r.index, fd, ws.id)
					ws.Lock()
					fmt.Printf("Testdelete5, index=%v, fd=%v, id=%v\n", r.index, fd, ws.id)
					ws.stop(err)
					fmt.Printf("Testdelete11, index=%v, fd=%v, id=%v\n", r.index, fd, ws.id)
					ws.Unlock()
				}
				fmt.Println("TestRead end", fd, ws.id, t)
			}
		}
		return err
	})
}

func (r *reactor) deleteAllConnections() {
	fmt.Println("delete All")
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
	c := conn.(*wsconn)
	if err := r.epoller.addRead(c.fd); err != nil {
		c.conn.Close()
		return err
	}
	fmt.Printf("Testadd, index=%v, fd=%v, id=%v\n", r.index, c.fd, c.id)
	c.server.clientOnConnect(c.id, c.fd, c)
	r.connections.registerConn(c.id, c.fd, c)
	r.addConnCount(1)
	return nil
}
