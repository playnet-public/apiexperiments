package transport

import (
	"context"
	"net/http"

	"github.com/TheMysteriousVincent/apiexperiments/encode"
	"github.com/TheMysteriousVincent/libs/problem"
	"go.uber.org/zap"
)

//HandlerFunc of the main service functions
type HandlerFunc func(ctx context.Context, r *http.Request) (interface{}, error)

//MiddlewareFunc of middlewares surrounding the main service
type MiddlewareFunc func(ctx context.Context, r *http.Request) error

//Endpoint of an service
type Endpoint interface {
	//PutMiddlewareBefore the handler.
	//PutMiddlewareBefore can be called multiple as every middleware is appended one after another
	PutMiddlewareBefore(mw MiddlewareFunc) Endpoint
	//PutMiddlewareAfter the handler.
	//PutMiddlewareAfter can be called multiple as every middleware is appended one after another
	PutMiddlewareAfter(mw MiddlewareFunc) Endpoint
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

func (e *endpoint) PutMiddlewareBefore(mw MiddlewareFunc) Endpoint {
	if e.before == nil {
		e.before = mw
	} else {
		e.before = func(ctx context.Context, r *http.Request) error {
			if err := e.before(ctx, r); err != nil {
				return err
			}

			return mw(ctx, r)
		}
	}

	return e
}

func (e *endpoint) PutMiddlewareAfter(mw MiddlewareFunc) Endpoint {
	if e.after == nil {
		e.after = mw
	} else {
		e.after = func(ctx context.Context, r *http.Request) error {
			if err := e.after(ctx, r); err != nil {
				return err
			}

			return mw(ctx, r)
		}
	}

	return e
}

func (e *endpoint) HandlerFunc(l *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := e.before(ctx, r); err != nil {
			e.encodeError(w, err, l)
			return
		}

		v, err := e.hndl(ctx, r)
		if err != nil {
			e.encodeError(w, err, l)
			return
		}

		if err := e.after(ctx, r); err != nil {
			e.encodeError(w, err, l)
			return
		}

		e.enc.Encode(w, v)
	}
}

func (e *endpoint) encodeError(w http.ResponseWriter, err error, l *zap.Logger) {
	if err := e.enc.Encode(w, e.finalizeError(err)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		l.Error(
			"encoding failed",
			zap.Error(err),
		)
	}
}

func (e *endpoint) finalizeError(err error) *problem.Problem {
	p, ok := err.(*problem.Problem)
	if !ok {
		p = problem.New(
			http.StatusText(http.StatusInternalServerError),
			err.Error(),
			http.StatusInternalServerError,
		)
	}

	if p.Status == 0 {
		p.Status = http.StatusInternalServerError
	}

	return p
}
