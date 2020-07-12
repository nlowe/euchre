package hosting

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/gammazero/nexus/v3/router"
	"github.com/go-chi/chi"
	"github.com/markbates/pkger"
	"github.com/nlowe/euchre/hosting/middleware"
	"github.com/nlowe/euchre/hosting/views"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type GameServer struct {
	http        *http.Server
	tableRouter router.Router

	tableServers map[string]*TableServer
}

func NewServer() (*GameServer, error) {
	tableServers := map[string]*TableServer{}

	nxr, err := router.NewRouter(&router.Config{
		RealmTemplate: &router.RealmConfig{
			AnonymousAuth: true,
			AllowDisclose: true,
		},
	}, logrus.WithField("prefix", "wamp::server"))
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(middleware.LogrusLogger())

	viewTemplates, err := views.Parse()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load views")
	}

	r.Get("/_play", router.NewWebsocketServer(nxr).ServeHTTP)

	r.Get("/table/{table}/", func(w http.ResponseWriter, r *http.Request) {
		table := chi.URLParam(r, "table")

		if _, ok := tableServers[table]; !ok {
			if tableServers[table], err = NewTableServer(nxr, table); err != nil {
				logrus.WithError(err).WithField("table", table).Error("Failed to create table server for table")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		log := logrus.WithField("prefix", fmt.Sprintf("%s::Table", table))
		log.WithFields(logrus.Fields{
			"referer": r.Referer(),
		}).Info("Someone joined the room!")
		w.Header().Set("Content-Type", "text/html")

		w.WriteHeader(http.StatusOK)
		if err := viewTemplates.ExecuteTemplate(w, "table", &views.TableViewModel{Table: table}); err != nil {
			log.WithError(err).Error("Failed to render page")
		}
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		name := fmt.Sprintf(
			"%s%s%s",
			strings.Title(petname.Adverb()),
			strings.Title(petname.Adjective()),
			strings.Title(petname.Name()),
		)

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		if err := viewTemplates.ExecuteTemplate(w, "index", &views.IndexViewModel{Table: name}); err != nil {
			logrus.WithError(err).Error("Failed to render page")
		}
	})

	static := pkger.Dir("/static")
	_ = pkger.Walk("/static", func(p string, _ os.FileInfo, _ error) error {
		logrus.WithFields(logrus.Fields{"prefix": "static", "file": p}).Debug("Found file")
		return nil
	})

	r.Get("/*", http.FileServer(static).ServeHTTP)

	s := &http.Server{
		Addr:         "0.0.0.0:5000",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return &GameServer{
		http:        s,
		tableRouter: nxr,

		tableServers: tableServers,
	}, nil
}

func (g *GameServer) ListenAndServe() error {
	return g.http.ListenAndServe()
}

func (g *GameServer) Shutdown(ctx context.Context) error {
	e := errgroup.Group{}
	for _, t := range g.tableServers {
		e.Go(t.Close)
	}

	e2 := errgroup.Group{}
	e2.Go(e.Wait)
	e2.Go(func() error {
		return g.http.Shutdown(ctx)
	})

	g.tableRouter.Close()
	return e2.Wait()
}
