package server

import "net/http"

type Ctx struct {
	w http.ResponseWriter
	r *http.Request

	Data map[string]interface{}
}

type Middleware = func(ctx *Ctx) (err error, code int)

type HttpHandler struct {
	middleware []Middleware
}

func (h *HttpHandler) Use(m Middleware) {
	h.middleware = append(h.middleware, m)
}

func (h *HttpHandler) GetHandleFunc() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &Ctx{
			w:    w,
			r:    r,
			Data: make(map[string]interface{}, 8),
		}
		for _, m := range h.middleware {
			err, code := m(ctx)
			if err != nil {
				http.Error(w, err.Error(), code)
				return
			}
		}
	}
}
