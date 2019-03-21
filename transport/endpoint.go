package transport

import (
	"context"
	"net/http"
	"github.com/playnet-public/apiexperiments/problems"
	"github.com/playnet-public/apiexperiments/encode"
	"go.uber.org/zap"
)

//HandlerFunc of the main service functions
type HandlerFunc func(ctx context.Context, r *http.Request) (interface{}, error)

//MiddlewareFunc of middlewares surrounding the main service
type MiddlewareFunc func(ctx context.Context) error

//Endpoint of an service
type Endpoint interface {
	//WithBefore places a middleware before the handler.
	//WithBefore can be called multiple as every middleware is appended one after another
	WithBefore(mw MiddlewareFunc) Endpoint
	//WithAfter places a middleware after the handler.
 	//WithAfter can be called multiple as every middleware is appended one after another
	WithAfter(mw MiddlewareFunc) Endpoint
	//HandlerFunc for the http endpoint
	HandlerFunc(*zap.Logger) http.HandlerFunc
}

type endpoint struct {
	before MiddlewareFunc
	after  MiddlewareFunc
	hndl   HandlerFunc
	enc    encode.Encoder
}

//NewEndpoint of an http service
func NewEndpoint(e encode.Encoder, hndl HandlerFunc) Endpoint {
	return &endpoint{
		hndl: hndl,
		enc:  e,
	}
}

func (e *endpoint) WithBefore(mw MiddlewareFunc) Endpoint {
	if e.before == nil {
		e.before = mw
	} else {
		e.before = func(ctx context.Context) error {
			if err := e.before(ctx); err != nil {
				return err
			}

			return mw(ctx)
		}
	}

	return e
}

func (e *endpoint) WithAfter(mw MiddlewareFunc) Endpoint {
	if e.after == nil {
		e.after = mw
	} else {
		e.after = func(ctx context.Context) error {
			if err := e.after(ctx); err != nil {
				return err
			}

			return mw(ctx)
		}
	}

	return e
}

func (e *endpoint) HandlerFunc(l *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := e.before(ctx); err != nil {
			e.encodeError(w, err, l)
			return
		}

		v, err := e.hndl(ctx, r)
		if err != nil {
			e.encodeError(w, err, l)
			return
		}

		if err := e.after(ctx); err != nil {
			e.encodeError(w, err, l)
			return
		}

		e.enc.Encode(w, v)
	}
}

func (e *endpoint) encodeError(w http.ResponseWriter, err error, l *zap.Logger) {
	if err := e.enc.Encode(w, problemError(err)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		l.Error(
			"encoding failed",
			zap.Error(err),
		)
	}
}
