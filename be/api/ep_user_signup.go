package api

import (
	util "Golang-REST-API-structure/be/lib"
	"Golang-REST-API-structure/be/lib/db"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type USignupReq struct {
	Fullname  string  `json:"fullname"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Role      string  `json:"role"`
}

func UserSignup(req *Req, resp *Resp) {
	fmt.Println("body")
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

	body := &USignupReq{}
	if err := json.Unmarshal(rawBody, body); err != nil {
		c.logger.Error(
			"failed to parse JSON object",
			zap.Error(err),
		)
		resp.Send(RC_E_MALFORMED)
		return
	}

	if body.Email == "" || !util.IsEmail(body.Email) {
		c.logger.Error(
			"invalid email",
			zap.Error(err),
		)
		resp.Send(RC_USER_EMAIL_NOT_FOUND)
		return
	}

	uID := uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.logger.Error("failed to hash the password")
		resp.Send(http.StatusInternalServerError)
		return
	}

	// Begin transaction.
	tx, err := c.Pg.Begin()
	if err != nil {
		c.logger.Error(
			"failed to begin transaction",
			zap.Error(err),
		)
		resp.Send(http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() 

	err = c.Db.Signup(nil, &db.USignup{
		ID:       uID,
		Fullname: body.Fullname,
		Email:    body.Email,
		Password: string(hashedPassword),
		Role:     body.Role,
	})
	if err != nil {
		c.logger.Error(
			"failed to sign up",
			zap.Error(err),
		)
		resp.Send(http.StatusInternalServerError)
		return
	}

	// Audit Logs
	action := "user.signup"
	err = c.Db.AuditLogs(nil, uID, action)
	if err != nil {
		c.logger.Error(
			"failed to update logs",
			zap.Error(err),
		)
		resp.Send(http.StatusInternalServerError)
		return
	}

	// Commit
	if err := tx.Commit(); err != nil {
		c.logger.Error(
			"failed to commit the transaction",
			zap.Error(err),
		)
		resp.Send(http.StatusInternalServerError)
        return
	}
	
	resp.SendData(RC_USER_SIGNUP, "Successful signup")
}
