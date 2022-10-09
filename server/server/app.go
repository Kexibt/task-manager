package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	authhttp "task_manager/auth/http"
	"task_manager/config"
	taskhttp "task_manager/task/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

type App struct {
	auth *authhttp.Auth
	task *taskhttp.Handler
	conn *pgxpool.Pool

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

	task := taskhttp.NewHandler(db)

	// fmt.Println(
	// 	task != nil,
	// 	task.Database != nil,
	// 	task.Localstore != nil,
	// )

	auth := authhttp.NewAuth(db)
	auth.AddNewMiddlewaredHandler("/create_task", task.CreateTask)
	auth.AddNewMiddlewaredHandler("/get_task", task.GetTask)
	auth.AddNewMiddlewaredHandler("/get_tasks", task.GetTasks)
	auth.AddNewMiddlewaredHandler("/get_tasks_after", task.GetTasksAfter)

	mux := http.NewServeMux()
	auth.AddHandlers(mux)
	// tasks todo

	return &App{
		task: task,
		auth: auth,
		conn: db,

		cfg: cfg,
		mux: mux,
	}
}

func (a *App) Run() {
	go func() {
		log.Printf("Server started at %s:%s", a.cfg.Host, a.cfg.Port)
		err := http.ListenAndServe(fmt.Sprintf("%s:%s", a.cfg.Host, a.cfg.Port), a.mux)
		log.Fatal(err)
	}()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign)

	<-sign
	a.conn.Close()
	log.Println("Server ended its serving")
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
