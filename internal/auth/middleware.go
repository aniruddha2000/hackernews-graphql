package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aniruddha2000/hackernews/internal/users"
	"github.com/aniruddha2000/hackernews/pkg/jwt"
)

type ContextKey struct {
	name string
}

var usrCtxKey = &ContextKey{"user"}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			tokenStr := header
			username, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid Token", http.StatusForbidden)
				return
			}

			user := users.Users{Username: username}
			id, err := users.GetUserIdByUsername(user.Username)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			user.ID = strconv.Itoa(int(id))
			ctx := context.WithValue(r.Context(), usrCtxKey, &user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) *users.Users {
	raw, _ := ctx.Value(usrCtxKey).(*users.Users)
	return raw
}
