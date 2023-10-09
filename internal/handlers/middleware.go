package handlers

import (
	"context"
	"log"
	"net/http"
)

func (h *Handler) CheckUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		token := request.Header.Get("Authentication")
		id, err := h.Service.TokenCheck(&token)
		log.Println(id)
		if err != nil {
			h.logger.Error(err)
			return
		}
		ctx := context.WithValue(request.Context(), "id", id)
		request = request.WithContext(ctx)
		next.ServeHTTP(response, request)
	})
}
