package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"task_manager/auth"
	"task_manager/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

type App struct {
	auth *auth.Auth
	// tasks task.Repository // todo

	cfg config.Config
	mux *http.ServeMux
}

func NewApp() *App {
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GetConnectionTimeout())
	defer cancel()

	db, err := connect(ctx, cfg.GetConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	// mux.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {}) todo

	return &App{
		// tasks: ,
		auth: auth.NewAuth(db),

		cfg: cfg,
		mux: mux,
	}
}

func (a *App) Run(port string) error {
	return http.ListenAndServe(fmt.Sprintf("%s:%s", a.cfg.Host, a.cfg.Port), a.mux)
}

func connect(ctx context.Context, connectionStr string) (conn *pgxpool.Pool, err error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err = pgxpool.Connect(context.Background(), connectionStr)
			if err == nil {
				return
			}
		}
	}
}
