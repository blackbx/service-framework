package test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BlackBX/service-framework/response"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// TodoModel is a struct that represents the API response
type TodoModel struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// NewHTTPHandler produces a new instance of the HTTPHandler type
func NewHTTPHandler(provider response.ResponderProvider, logger *zap.Logger, client *http.Client) HTTPHandler {
	return HTTPHandler{
		ResponseProvider: provider,
		Logger:           logger,
		Client:           client,
	}
}

// HTTPHandler is a type that will reach out to a third party service
type HTTPHandler struct {
	ResponseProvider response.ResponderProvider
	Logger           *zap.Logger
	Client           *http.Client
}

// Get is the function that is called when the route is hit
func (h HTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
	responder := h.ResponseProvider.Responder(w, r)
	id, ok := mux.Vars(r)["id"]
	if !ok {
		responder.RespondWithProblem(http.StatusBadRequest, "Please specify an ID")
		return
	}
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/todos/%s", id)
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		responder.RespondWithProblem(http.StatusInternalServerError, "Could not build request")
		return
	}
	req = req.WithContext(r.Context())
	resp, err := h.Client.Do(req)
	if err != nil {
		responder.RespondWithProblem(http.StatusInternalServerError, "Could not reach the external service")
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			h.Logger.Error("could not close response body", zap.Error(err))
		}
	}()
	todo := TodoModel{}
	if err := json.NewDecoder(resp.Body).Decode(&todo); err != nil {
		responder.RespondWithProblem(http.StatusInternalServerError, "Could not parse response to TodoModel")
		return
	}
	responder.Respond(http.StatusOK, todo)
}
