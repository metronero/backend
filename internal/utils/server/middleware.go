package server

import (
	"context"
	"net/http"
	"strings"

	"gitea.com/go-chi/session"
	"github.com/mr-tron/base58"
	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/internal/utils/helpers"
	"gitlab.com/metronero/backend/pkg/apierror"
)

const version = "0.1.0"

type ApiKeyData struct {
	Accepted  bool
	AccountId string
	Username  string
	Role      string
}

var ApiKeyCache map[string]ApiKeyData

func middlewareServerHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "Metronero/"+version)
		next.ServeHTTP(w, r)
	})
}

func middlewareAdminArea(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctxRole := ctx.Value("role")
		if role, ok := ctxRole.(string); !ok || role != "admin" {
			http.Error(w, http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func middlewareMerchantArea(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctxRole := ctx.Value("role")
		if role, ok := ctxRole.(string); !ok || role != "merchant" {
			http.Error(w, http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Check if the session is authenticated regardless of role
func middlewareAuthArea(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := session.GetSession(r)
		userId := sess.Get("account_id")
		userIdStr, ok := userId.(string)
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized)
			return
		}
		c1 := context.WithValue(r.Context(), "account_id", userIdStr)
		c2 := context.WithValue(c1, "role", sess.Get("role").(string))
		c3 := context.WithValue(c2, "username", sess.Get("username").(string))
		next.ServeHTTP(w, r.WithContext(c3))
	})
}

func ApiKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized)
			return
		}
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			helpers.WriteError(w, apierror.ErrBadRequest, nil)

			return
		}
		apiKey := strings.TrimPrefix(authHeader, bearerPrefix)
		const tokenPrefix = "mnero_"
		if !strings.HasPrefix(apiKey, tokenPrefix) {
			helpers.WriteError(w, apierror.ErrBadRequest, nil)
			return
		}
		parts := strings.SplitN(strings.TrimPrefix(apiKey, tokenPrefix), "_", 2)
		if len(parts) != 2 {
			helpers.WriteError(w, apierror.ErrBadRequest, nil)
			return
		}
		keyId := parts[0]
		keySecret := parts[1]
		ctx := r.Context()
		if data, ok := ApiKeyCache[keyId]; ok {
			if !data.Accepted {
				helpers.WriteError(w, apierror.ErrUnauthorized, nil)
				return
			}
			c1 := context.WithValue(r.Context(), "account_id", data.AccountId)
			c2 := context.WithValue(c1, "role", data.Role)
			c3 := context.WithValue(c2, "username", data.Username)
			next.ServeHTTP(w, r.WithContext(c3))
		}

		keyUuid, err := base58.Decode(keyId)
		if err != nil {
			helpers.WriteError(w, apierror.ErrBadRequest, err)
			return
		}

		accepted, accountId, role, name, appErr := queries.CheckKey(ctx, string(keyUuid), keySecret)
		if appErr != nil {
			helpers.WriteError(w, apierror.ErrDatabase, appErr)
			return
		}
		ApiKeyCache[keyId] = ApiKeyData{
			Accepted:  accepted,
			AccountId: accountId,
			Role:      role,
			Username:  name,
		}
		if !accepted {
			helpers.WriteError(w, apierror.ErrUnauthorized, nil)
			return
		}
		c1 := context.WithValue(r.Context(), "account_id", accountId)
		c2 := context.WithValue(c1, "role", role)
		c3 := context.WithValue(c2, "username", name)
		next.ServeHTTP(w, r.WithContext(c3))
	})
}

func disableDirListing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" || strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
