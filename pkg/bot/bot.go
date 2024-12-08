package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-joe/joe"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/bonaysoft/gofbot/pkg/robot"
)

type Server struct {
	bot *joe.Bot

	botAnswers *robot.Robot
}

func NewServer(adapter Adapter) (*Server, error) {
	bot := joe.New("gofbot", joe.WithLogLevel(zap.DebugLevel), adapter.Adapter())
	bot.Brain.RegisterHandler(adapter.GetHandler(bot))
	bot.Respond("ping", func(msg joe.Message) error {
		msg.Respond("pong")
		return nil
	})
	answers, err := robot.Load("./robots")
	if err != nil {
		return nil, err
	}
	answer, ok := lo.Find(answers, func(item *robot.Robot) bool { return item.Kind == adapter.Name() })
	if !ok {
		return nil, fmt.Errorf("not found robot for %s", adapter)
	}

	return &Server{bot: bot, botAnswers: answer}, nil
}

func (b *Server) Run(addr string) error {
	go func() {
		http.HandleFunc("POST /api/webhooks/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			var chatID string
			exists, _ := b.bot.Store.Get(id, &chatID)
			if !exists {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			output, err := b.makeRobotResponse(r)
			if err != nil {
				b.bot.Logger.Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// b.bot.Adapter
			b.bot.Say(chatID, string(output))
			w.WriteHeader(http.StatusOK)
		})
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err)
			return
		}
	}()

	return b.bot.Run()
}

func (b *Server) makeRobotResponse(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	params := make(robot.Map)
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, err
	}

	msg, ok := b.botAnswers.MatchMsg(body)
	if !ok {
		return nil, fmt.Errorf("not found")
	}

	return []byte(msg.Build(params)), nil
}
