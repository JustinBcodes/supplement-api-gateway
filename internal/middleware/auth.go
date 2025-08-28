package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"supplx-gateway-marketplace/internal/common"
	"supplx-gateway-marketplace/internal/middleware/userctx"
)

type AuthConfig struct {
	Issuer string
	Audience string
	JWKSURL string
}

type AuthMiddleware struct {
	JWKS *common.JWKSCache
	Cfg  AuthConfig
}

func (a *AuthMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authz := r.Header.Get("Authorization")
		if !strings.HasPrefix(authz, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tok := strings.TrimPrefix(authz, "Bearer ")
		jwks, err := a.JWKS.Get(a.Cfg.JWKSURL)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		parser := jwt.NewParser()
		claims := jwt.MapClaims{}
		_, err = parser.ParseWithClaims(tok, claims, func(token *jwt.Token) (any, error) {
			kid, _ := token.Header["kid"].(string)
			for _, k := range jwks.Keys {
				if k["kid"] == kid && k["kty"] == "RSA" {
					nStr, _ := k["n"].(string)
					eStr, _ := k["e"].(string)
					nBytes, _ := base64.RawURLEncoding.DecodeString(nStr)
					eBytes, _ := base64.RawURLEncoding.DecodeString(eStr)
					var nn big.Int
					nn.SetBytes(nBytes)
					var ee int
					if len(eBytes) > 0 {
						ee = 0
						for _, b := range eBytes { ee = ee<<8 + int(b) }
					} else { ee = 65537 }
					return &rsa.PublicKey{N: &nn, E: ee}, nil
				}
			}
			return nil, jwt.ErrTokenUnverifiable
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if iss, _ := claims["iss"].(string); a.Cfg.Issuer != "" && iss != a.Cfg.Issuer { w.WriteHeader(http.StatusUnauthorized); return }
		if aud, _ := claims["aud"].(string); a.Cfg.Audience != "" && aud != a.Cfg.Audience { w.WriteHeader(http.StatusUnauthorized); return }
		sub, _ := claims["sub"].(string)
		next.ServeHTTP(w, r.WithContext(userctx.WithSub(r.Context(), sub)))
	})
}


