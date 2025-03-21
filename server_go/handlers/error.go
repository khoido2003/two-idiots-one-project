package handlers

import "net/http"

func (handlerContext *HandlerContext) HandleError(w http.ResponseWriter, r *http.Request) {
        respondWithError(w, 400, "Something went wrong")
}
