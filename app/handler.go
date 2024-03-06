package app

import (
	"html/template"
	"net/http"
	"stdrifa/app/company"
	"stdrifa/model"

	"github.com/go-chi/httplog/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jritsema/gotoolbox/web"
)

type App struct {
	DB     *pgxpool.Pool
	Q      *model.Queries
	Logger *httplog.Logger
	Templ  *template.Template
}

func New(db *pgxpool.Pool, logger *httplog.Logger, templ *template.Template) *App {
	return &App{
		DB:     db,
		Q:      model.New(db),
		Logger: logger,
		Templ:  templ,
	}
}

func (a *App) Index(w http.ResponseWriter, r *http.Request) {
	res := web.HTML(http.StatusOK, a.Templ, "index.html", company.GetCompanies(), nil)
	res.Write(w)
}
