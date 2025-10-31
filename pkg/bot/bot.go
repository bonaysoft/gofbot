package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-joe/file-memory"
	"github.com/go-joe/joe"
	"go.uber.org/zap"

	"github.com/bonaysoft/gofbot/pkg/messenger"
)

type Server struct {
	bot *joe.Bot

	messenger messenger.Manager
}

func NewServer(adapter Adapter, messenger messenger.Manager) (*Server, error) {
	bot := joe.New("gofbot",
		joe.WithLogLevel(zap.DebugLevel),
		file.Memory("./data/memory.json"),
		adapter.Adapter())
	bot.Respond("ping", func(msg joe.Message) error {
		msg.Respond("pong")
		return nil
	})

	for _, handler := range adapter.GetHandlers(bot) {
		bot.Brain.RegisterHandler(handler)
	}
	return &Server{bot: bot, messenger: messenger}, nil
}

func (b *Server) Run(addr string) error {
	go func() {
		http.HandleFunc("POST /api/webhooks/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := fmt.Sprintf("WH_%s", r.PathValue("id"))
			var chat Chat
			exists, _ := b.bot.Store.Get(id, &chat)
			if !exists {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			output, err := b.makeRobotResponse(r, chat)
			if err != nil {
				b.bot.Logger.Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// b.bot.Adapter
			b.bot.Say(chat.Channel, string(output))
			w.WriteHeader(http.StatusOK)
		})
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(fmt.Errorf("server start failed: %s", err))
			return
		}
	}()

	return b.bot.Run()
}

func (b *Server) makeRobotResponse(r *http.Request, chat Chat) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	params := make(map[string]any)
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, fmt.Errorf("decode body: %w", err)
	}
	b.bot.Logger.Sugar().Debugf("received body", zap.Any("params", params))

	params["chatProvider"] = chat.Provider
	msg, err := b.messenger.Match(params)
	if err != nil {
		return nil, fmt.Errorf("match message: %w", err)
	}

	return b.messenger.BuildReply(msg, params)
}
