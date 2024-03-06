package server

import (
	"context"
	"embed"
	"html/template"
	"log/slog"
	"stdrifa/config"
	"stdrifa/model"
	"time"

	"github.com/go-chi/httplog/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jritsema/gotoolbox/web"
)

type Server struct {
	Name string
	Host string
	Port string

	Logger *httplog.Logger
	DB     *pgxpool.Pool
	Q      *model.Queries

	PublicFS   embed.FS
	TemplateFS embed.FS
	Templ      *template.Template
}

func New() *Server {
	s := &Server{
		Name: config.Name,
		Host: config.Host,
		Port: config.Port,

		Logger: InitLogger(),
	}

	return s
}

func InitLogger() *httplog.Logger {
	logLevel := slog.LevelDebug
	if config.IsProduction {
		logLevel = slog.LevelInfo
	}

	return httplog.NewLogger(config.Name, httplog.Options{
		JSON:             config.IsProduction,
		LogLevel:         logLevel,
		Concise:          !config.IsProduction,
		RequestHeaders:   config.IsProduction,
		MessageFieldName: "message",
		Tags: map[string]string{
			"env":     config.Env,
			"version": config.Version,
		},
		QuietDownRoutes: []string{
			"/",
			"/ping",
		},
		QuietDownPeriod: 10 * time.Second,
	})
}

func (s *Server) ParseTemplates(fileSystem embed.FS) {
	html, err := web.TemplateParseFSRecursive(s.TemplateFS, ".html", true, nil)
	if err != nil {
		panic(err)
	}

	s.Templ = html
}

func (s *Server) ConnectDB(ctx context.Context) {
	db, err := pgxpool.New(ctx, config.DatabaseURL)
	if err != nil {
		s.Logger.Error("Failed to connect to database:", err)
		panic(err)
	}

	s.DB = db
	s.Q = model.New(db)
}
