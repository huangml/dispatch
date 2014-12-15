package dispatch

import (
	"fmt"

	"github.com/huangml/mux"
)

type LockedDest struct {
	h HandlerFunc
	m Mutex
}

func NewLockedDest(h HandlerFunc) *LockedDest {
	return &LockedDest{h: h, m: NewMutex()}
}

func (d *LockedDest) Call(ctx *Context, r Request) Response {
	return d.h(ctx, d.m, r)
}

func (d *LockedDest) Send(r Request) error {
	go d.h(NewContext(), d.m, r)
	return nil
}

type MuxDest struct {
	mu  *mux.Mux
	mtx Mutex
}

func NewMuxDest() *MuxDest {
	return &MuxDest{mu: mux.NewPathMux(), mtx: NewMutex()}
}

func (d *MuxDest) HandleFunc(pattern string, h HandlerFunc) {
	d.mu.Map(pattern, h)
}

func (d *MuxDest) Call(ctx *Context, r Request) Response {
	if h := d.mu.Match(r.Protocol()); h != nil {
		return h.(HandlerFunc)(ctx, d.mtx, r)
	}
	return ErrResponse(protocolNotFoundErr(r.Protocol()))
}

func (d *MuxDest) Send(r Request) error {
	if h := d.mu.Match(r.Protocol()); h != nil {
		go h.(HandlerFunc)(NewContext(), d.mtx, r)
		return nil
	}
	return protocolNotFoundErr(r.Protocol())
}

func protocolNotFoundErr(protocol string) error {
	return fmt.Errorf("protocol %v not implemented", protocol)
}
