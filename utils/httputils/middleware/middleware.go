package middleware

import "applicationDesignTest/utils/httputils"

type Handler func(next httputils.Handler) httputils.Handler

func WrapHandler(h httputils.Handler, mws ...Handler) httputils.Handler {
	for _, mw := range mws {
		h = mw(h)
	}
	return h
}
