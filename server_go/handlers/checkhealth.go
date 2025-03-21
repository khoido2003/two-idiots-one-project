package handlers

import "net/http"

//  Create new context for handler
func (handlerContext *HandlerContext) Ping(w http.ResponseWriter, r *http.Request) {
	respondWithString(w, 200, "PING")
}

func (HandlerContext *HandlerContext) CheckHealth(w http.ResponseWriter, r *http.Request) {
    respondWithString(w, 200, "PONG")
}


