package test

import (
	"net/http"

	"github.com/BlackBX/service-framework/response"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// PGHandler is the handler that communicates with a postgres database
type PGHandler struct {
	DB               *sqlx.DB
	logger           *zap.Logger
	ResponseProvider response.ResponderProvider
}

// NewPGHandler is a function that creates a new instance of the PGHandler type
func NewPGHandler(db *sqlx.DB, logger *zap.Logger, responseProvider response.ResponderProvider) PGHandler {
	return PGHandler{DB: db, logger: logger, ResponseProvider: responseProvider}
}

// Get is a function that is called to pull data from the database
func (h PGHandler) Get(w http.ResponseWriter, r *http.Request) {
	responder := h.ResponseProvider.Responder(w, r)
	res := &DBResponse{}
	err := h.DB.GetContext(r.Context(), res, "SELECT 1 + 1 as result")

	if err != nil {
		h.logger.Error("The database is borked", zap.Error(err))
		responder.RespondWithProblem(http.StatusInternalServerError, "The database is broken :(")
		return
	}
	responder.Respond(http.StatusOK, res)
}

// DBResponse is the model that represents the response from the database
type DBResponse struct {
	Result int64 `json:"result" db:"result"`
}
