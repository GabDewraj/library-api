package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/GabDewraj/library-api/pkgs/infrastructure/cache"
	"github.com/sirupsen/logrus"
)

func (s *service) RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIp := s.getClientIp(r)
		clientKey := "rate:" + clientIp
		// First do an existence check
		exists, err := s.Cache.ExistenceCheck(r.Context(), clientKey)
		if err != nil {
			logrus.Error(err)
			http.Error(w, "Could not check for key existence in cache", http.StatusInternalServerError)
			return
		}

		switch exists {
		case true:
			count, err := s.Cache.RetrieveInteger(r.Context(), clientKey)
			if err != nil {
				logrus.Error(err)
				http.Error(w, "Could not convert payload value to int64", http.StatusInternalServerError)
				return
			}

			if count >= 5 {
				err := errors.New("Client has hit rate limit")
				http.Error(w, err.Error(), http.StatusTooManyRequests)
				return
			}
			// Increment request count
			_, err = s.Cache.KeyIncrement(r.Context(), clientKey)
			if err != nil {
				logrus.Error(err)
				http.Error(w, "Could not increment client key", http.StatusInternalServerError)
				return
			}

		default:
			if err := s.Cache.StoreInteger(r.Context(), cache.CacheIntegerPayload{Key: clientKey, Value: 1}); err != nil {
				logrus.Error(err)
				http.Error(w, "Could not store user request ip in cache", http.StatusInternalServerError)
				return
			}
		}

		// Serve the request
		next.ServeHTTP(w, r)
	})
}

func (s *service) getClientIp(r *http.Request) string {
	ipAddress := r.RemoteAddr
	// Remove the port number if present
	if colonIndex := strings.LastIndex(ipAddress, ":"); colonIndex != -1 {
		ipAddress = ipAddress[:colonIndex]
	}
	return ipAddress
}
