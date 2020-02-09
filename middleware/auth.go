package middleware

import (
	"net/http"

	"github.com/AI-Research-HIT/2019-nCoV-Service/auth"
)

func JwtAuthMw(next http.Handler) http.Handler {
	return auth.Jwt.Handler(next)
}
