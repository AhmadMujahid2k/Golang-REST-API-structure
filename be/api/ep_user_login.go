package api

import (
	util "Golang-REST-API-structure/be/lib"
	"Golang-REST-API-structure/be/lib/db"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type ULoginReq struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
}

type ULoginResp struct {
	ID        uuid.UUID  `json:"id"`
	Fullname  string     `json:"fullname"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
}

func UserLogin(req *Req, resp *Resp) {
	defer req.Body.Close()
	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		c.logger.Error(
			"failed to parse req body",
			zap.Error(err),
		)
		resp.Send(RC_E_NO_BODY)
		return
	}

	body := &ULoginReq{}
	if err := json.Unmarshal(rawBody, body); err != nil {
		c.logger.Error(
			"failed to parse JSON object",
			zap.Error(err),
		)
		resp.Send(RC_E_MALFORMED)
		return
	}

	if body.Email == "" || !util.IsEmail(body.Email) {
		resp.Send(RC_USER_NO_EMAIL)
		return
	}

	// Check User Existence
	tup, err := c.Db.ULogin(nil, body.Email)
	if err != nil {
		c.logger.Error("failed to fetch from db", zap.Error(err))
		resp.Send(RC_USER_NOT_EXIST)
		return
	}

	// Compared Passwords
	err = bcrypt.CompareHashAndPassword([]byte(tup.Password), []byte(body.Password))
	if err != nil {
		
		action := ""
		if tup.Role == "USER" {
			action = "user.failed_login"
		}else {
			action = "agent.failed_login"
		}

		// Audit Logs
		if tup.ID == uuid.Nil {
			c.Db.AuditLogs(nil, tup.ID, action)
		}
	
		if err == bcrypt.ErrMismatchedHashAndPassword {
			c.logger.Error("invalid credentials", zap.Error(err))
			resp.Send(RC_USER_INVALID_CREDENTIALS)
		} else {
			c.logger.Error("internal server error", zap.Error(err))
			resp.Send(http.StatusInternalServerError)
		}
		return
	}

	// Create Session with ID and Role
	req.Session.Values[0] = db.SessionVals{
		Id:   tup.ID,
		Role: tup.Role,
	}
	err = req.Session.Save(req.Request, resp.Resp)
	if err != nil {
		c.logger.Error(
			"failed to save session",
			zap.Error(err),
		)
		resp.Send(RC_USER_SAVE_SESSION_FAIL)
		return
	}

	action := ""
	if tup.Role == "USER" {
		action = "user.successful_login"
	} else {
		action = "agent.successful_login"
	}

	// Audit Logs
	err = c.Db.AuditLogs(nil, tup.ID, action)
	if err != nil {
		c.logger.Error(
			"failed to update logs",
			zap.Error(err),
		)
		resp.Send(http.StatusInternalServerError)
		return
	}

	resp.SendData(RC_USER_LOGIN, ULoginResp{
		ID:       tup.ID,
		Fullname: tup.Fullname,
		Role:     tup.Role,
		Email:    tup.Email,
	})
}
