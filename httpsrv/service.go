package httpsrv

import (
	"context"
	"net"
	"net/http"
	"sync"
)

type Server struct {
	HttpServer              http.Server
	TLSCertFile, TLSKeyFile string

	// Private fields
	mut               sync.RWMutex
	onReloadListeners []func(context.Context, interface{}) error
}

func (srv *Server) Serve(ctx context.Context) (err error) {
	srv.HttpServer.BaseContext = func(_ net.Listener) context.Context { return ctx }

	if len(srv.TLSCertFile) > 0 && len(srv.TLSKeyFile) > 0 {
		err = srv.HttpServer.ListenAndServeTLS(srv.TLSCertFile, srv.TLSKeyFile)
	} else {
		err = srv.HttpServer.ListenAndServe()
	}
	if err == http.ErrServerClosed {
		err = nil
	}
	return
}

func (srv *Server) Reload(ctx context.Context, data interface{}) (err error) {
	srv.mut.RLock()
	defer srv.mut.RUnlock()

	for _, f := range srv.onReloadListeners {
		if err = f(ctx, data); err != nil {
			break
		}
	}
	return
}

func (srv *Server) Shutdown(ctx context.Context) (err error) {
	return srv.HttpServer.Shutdown(ctx)
}

func (srv *Server) RegisterOnReload(f func(context.Context, interface{}) error) {
	if f != nil {
		srv.mut.Lock()
		srv.onReloadListeners = append(srv.onReloadListeners, f)
		srv.mut.Unlock()
	}
}
