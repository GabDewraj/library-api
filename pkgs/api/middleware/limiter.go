package middleware

import (
	"encoding/binary"
	"errors"
	"net/http"
	"strings"
	"time"

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
			payload, err := s.Cache.Retrieve(r.Context(), []string{clientKey})
			if err != nil {
				logrus.Error(err)
				http.Error(w, "Could not get user request number from cache", http.StatusInternalServerError)
				return
			}

			if len(payload) > 0 {
				count, err := s.binaryVarint(payload[0].Value)
				if err != nil {
					logrus.Error(err)
					http.Error(w, "Could not convert payload value to int64", http.StatusInternalServerError)
					return
				}

				if count > 5 {
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
			}
		default:
			if err := s.Cache.Store(r.Context(), []*cache.CachePayload{
				{
					Key:        clientKey,
					Value:      []byte{1},
					Expiration: 3 * time.Minute,
				},
			}); err != nil {
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

func (s *service) binaryVarint(data []byte) (int64, error) {
	value, n := binary.Varint(data)
	if n <= 0 {
		return 0, errors.New("binaryVarint: decoding failure")
	}
	return value, nil
}
