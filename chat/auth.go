package main

import (
	"net/http"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		// 인증 불가
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		// 다른 에러
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 성공 - 다음 핸들러 호출
	h.next.ServeHTTP(w, r)
}

// MustAuth ...
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
