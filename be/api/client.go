package api

import (
	util "Golang-REST-API-structure/be/lib"
	"Golang-REST-API-structure/be/lib/db"
	"Golang-REST-API-structure/be/lib/psql"
	"context"
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

var c *Client

type Client struct {
    Pg             *psql.Postgres
    Db             *db.DbC
    logger         *zap.Logger
    SessionStore   *sessions.CookieStore
    SessionName    string
    SessionDomain  string
    SessionKey     string
}

func Init() *Client {
	var err error
	c = &Client{}
	c.Pg, err = psql.Init(context.Background(), util.MustOsGetEnv("DB_URL"))
	if err != nil {
		panic("unable to connect to database")
	}

	c.Db = db.Init(c.Pg)

	c.logger = zap.Must(zap.NewProduction())
	defer c.logger.Sync()

	c.SessionKey = util.MustOsGetEnv("USER_SESSION_KEY")
	c.SessionDomain = util.MustOsGetEnv("USER_SESSION_DOMAIN")
	c.SessionName = util.MustOsGetEnv("USER_SESSION_NAME")
	c.SessionStore = sessions.NewCookieStore([]byte(c.SessionKey))

	c.SessionStore.Options = &sessions.Options{
		Path:     "/",
		Domain:   c.SessionDomain,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	}
	gob.Register(db.SessionVals{})

	return c
}
