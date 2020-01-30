package test

import (
	"net/http"

	"github.com/BlackBX/service-framework/redis"
	"github.com/BlackBX/service-framework/response"
)

// NewRedisHandler is the constructor for RedisHandler
func NewRedisHandler(redisClient redis.Cmdable, responder response.ResponderProvider) RedisHandler {
	return RedisHandler{
		Redis:            redisClient,
		ResponseProvider: responder,
	}
}

// RedisHandler is the handler
type RedisHandler struct {
	Redis            redis.Cmdable
	ResponseProvider response.ResponderProvider
}

// Get is the function that is called when the route is routed to
func (h RedisHandler) Get(rw http.ResponseWriter, r *http.Request) {
	responder := h.ResponseProvider.Responder(rw, r)
	key := r.URL.Query().Get("key")
	if key == "" {
		key = "key"
	}
	count, err := h.Redis.WithContext(r.Context()).Incr(key).Result()
	if err != nil {
		responder.RespondWithProblem(http.StatusInternalServerError, err.Error())
		return
	}
	responder.Respond(http.StatusOK, RedisResponse{Count: count})
}

// RedisResponse is the response returned by the RedisHandler
type RedisResponse struct {
	Count int64 `json:"count"`
}
