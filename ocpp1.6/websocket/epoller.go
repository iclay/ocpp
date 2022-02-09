package websocket

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"runtime"
	"sync/atomic"
	"unsafe"
)

const (
	readEvents      = unix.EPOLLPRI | unix.EPOLLIN
	writeEvents     = unix.EPOLLOUT
	readWriteEvents = readEvents | writeEvents
)

const (
	defauleEpollEventCap = 2 >> 1
)

type epollevent struct {
	size   int
	events []unix.EpollEvent
}

func initEpollEvents(size int) *epollevent {
	return &epollevent{
		size:   size,
		events: make([]unix.EpollEvent, size),
	}
}

type epoller struct {
	epollfd    int
	eventfd    int
	eventfdBuf []byte
	wakeSignal int32
	asyncQueue Queue
}

var (
	u uint64 = 1
	b        = (*(*[8]byte)(unsafe.Pointer(&u)))[:]
)

func createEpoller() (*epoller, error) {
	var err error
	epoller := new(epoller)
	if epoller.eventfd, err = unix.Eventfd(0, unix.EFD_NONBLOCK|unix.EFD_CLOEXEC); err != nil {
		return nil, os.NewSyscallError("eventfd", err)
	}
	defer func() {
		if err != nil {
			_ = os.NewSyscallError("close eventfd", unix.Close(epoller.eventfd))
		}
	}()
	// if epoller.epollfd, err = unix.EpollCreate1(unix.EPOLL_CLOEXEC); err != nil {
	if epoller.epollfd, err = unix.EpollCreate(1024); err != nil {
		return nil, os.NewSyscallError("epoll_create1", err)
	}
	defer func() {
		if err != nil {
			_ = os.NewSyscallError("close epoll_create1", unix.Close(epoller.epollfd))
		}
	}()
	epoller.eventfdBuf = make([]byte, 8)
	if err = os.NewSyscallError("epoll_ctl add",
		unix.EpollCtl(epoller.epollfd, unix.EPOLL_CTL_ADD, epoller.eventfd, &unix.EpollEvent{Fd: int32(epoller.eventfd), Events: readEvents})); err != nil {
		return nil, err
	}
	epoller.asyncQueue = NewEpollEventsQueue()
	return epoller, nil
}

func (e *epoller) trigger(fn TaskFunc, arg interface{}) (err error) {
	task := GetTask()
	task.RunFunc, task.RunArg = fn, arg
	e.asyncQueue.Push(task)
	if atomic.CompareAndSwapInt32(&e.wakeSignal, 0, 1) {
		for _, err = unix.Write(e.eventfd, b); err == unix.EINTR || err == unix.EAGAIN; _, err = unix.Write(e.eventfd, b) {
		}
	}
	return os.NewSyscallError(fmt.Sprintf("write  eventfd(%v)", e.eventfd), err)
}

func (e *epoller) close() error {
	if err := os.NewSyscallError("close epollfd", unix.Close(e.epollfd)); err != nil {
		return err
	}
	return os.NewSyscallError("close eventfd", unix.Close(e.eventfd))
}

func (e *epoller) addRead(fd int) error {
	return os.NewSyscallError("epoll_ctl add",
		unix.EpollCtl(e.epollfd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: readEvents}))
}

func (e *epoller) delete(fd int) error {
	return os.NewSyscallError("epoll_ctl del", unix.EpollCtl(e.epollfd, unix.EPOLL_CTL_DEL, fd, nil))
}

func (e *epoller) reactor(fn func(fd int, ev uint32) error) error {
	epollEvents := initEpollEvents(defauleEpollEventCap)
	var wake bool
	mesc := -1
	for {
		n, err := unix.EpollWait(e.epollfd, epollEvents.events, mesc)
		if n == 0 || (n < 0 && err == unix.EINTR) {
			runtime.Gosched()
			mesc = -1
			continue
		} else if err != nil {
			return err
		}
		mesc = 0
		for i := 0; i < n; i++ {
			event := epollEvents.events[i]
			if fd := int(event.Fd); fd != e.eventfd {
				switch err := fn(fd, event.Events); err {
				case nil:
				case ErrReactorShutdown:
					return err
				default:
					log.Warnf("occues errors in reactor:%v", err)
				}
			} else {
				wake = true
				_, _ = unix.Read(e.eventfd, e.eventfdBuf)
			}
		}
		if wake {
			wake = false
			for t, _ := e.asyncQueue.Pop(); t != nil; t, _ = e.asyncQueue.Pop() {
				task := t.(*Task)
				switch err = task.RunFunc(task.RunArg); err {
				case nil:
				case ErrReactorShutdown:
					return err
				default:
					log.Warnf("occues errors in custom function in reactor:%v", err)
				}
				PutTask(task)
			}
			atomic.StoreInt32(&e.wakeSignal, 0)
		}
	}
}
