package reader

import (
	"github.com/NYTimes/gizmo/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides all of the dependencies required by the reader module
var Module = fx.Provide(
	NewReader,
)

// NewReader is the constructor of the reader type
func NewReader(
	queue pubsub.Subscriber,
	logger *zap.Logger,
	shutdowner fx.Shutdowner,
) Reader {
	return Reader{
		Queue:      queue,
		Logger:     logger,
		Shutdowner: shutdowner,
	}
}

// Reader reads from an SQS Queue and logs out the messages
type Reader struct {
	Queue      pubsub.Subscriber
	Logger     *zap.Logger
	Shutdowner fx.Shutdowner
}

// Run starts the reader reading messages and logs them out, it will
// shutdown the server if the queue doesn't return any more messages
func (r Reader) Run() {
	messages := r.Queue.Start()
	for m := range messages {
		body := string(m.Message())
		r.Logger.Info("Got message:", zap.String("body", body))
		if err := m.Done(); err != nil {
			r.Logger.Error("Error when setting message as done", zap.Error(err))
		}
	}
	if err := r.Queue.Err(); err != nil {
		r.Logger.Error("Error received after queue stopped", zap.Error(err))
	}
	if err := r.Shutdowner.Shutdown(); err != nil {
		r.Logger.Error("Got error shutting down server", zap.Error(err))
	}
}

// Run is a function to be used in invoke to start the server running.
func Run(reader Reader) {
	go reader.Run()
}
