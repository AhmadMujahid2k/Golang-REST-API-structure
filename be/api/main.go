package api

import (
	"Golang-REST-API-structure/be"
	util "Golang-REST-API-structure/be/lib"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

type Req struct {
	be.Req
	Session *sessions.Session
}

type Resp struct {
	be.Resp
}

func Main() {
	c = Init()
	listenUrl := util.MustOsGetEnv("LISTEN_URL")
	server := http.NewServeMux()
	server.HandleFunc("/", MainHandler)
	c.logger.Info(
		"listening at ",
		zap.String("listenUrl", listenUrl),
	)

	err := http.ListenAndServe(listenUrl, server)
	if err != nil {
		c.logger.Info(
			"failed to listen and serve",
			zap.String("listenUrl", listenUrl),
		)
	}
}

func MainHandler(w http.ResponseWriter, httpreq *http.Request) {
	req := &Req{Req: be.Req{Request: httpreq}}
	resp := &Resp{Resp: be.Resp{ResponseWriter: w}}

	req.Session, _ = c.SessionStore.Get(req.Request, c.SessionName)

	url := req.URL.Path
	if url != "/" {
		url = strings.TrimRight(req.URL.Path, "/")
	}

	c.logger.Info(
		"request received",
		zap.String("method", req.Method),
		zap.String("url", url),
	)

	if req.Method == "GET" {
		if url == "/api/ahmadmujahid/v1/contact/list" {
			ContactList(req, resp)
		} else if url == "/api/ahmadmujahid/v1/auditlogs/list" {
			AuditLogList(req, resp)
		} else {
			resp.Send(http.StatusNotFound)
		}
	} else if req.Method == "POST" {
		if url == "/api/ahmadmujahid/v1/user/signup" {
			UserSignup(req, resp)
		} else if url == "/api/ahmadmujahid/v1/user/login" {
			UserLogin(req, resp)
		} else if url == "/api/ahmadmujahid/v1/user/pw_change" {
			UserPwChange(req, resp)
		} else if url == "/api/ahmadmujahid/v1/user/profile_upd" {
			UserProfileUpd(req, resp)
		} else if url == "/api/ahmadmujahid/v1/contact/upload" {
			ContactUpload(req, resp)
		} else if url == "/api/ahmadmujahid/v1/contact/bulk/upload" {
			ContactBulkUpload(req, resp)
		} else {
			resp.Send(http.StatusNotFound)
		}
	} else if req.Method == "DELETE" {
		if url == "/api/ahmadmujahid/v1/contact/delete" {
			ContactDelete(req, resp)
		} else {
			resp.Send(http.StatusNotFound)
		}
	} else {
		resp.Send(http.StatusMethodNotAllowed)
	}
}
