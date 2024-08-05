package api

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UPwChangeReq struct {
	OldPw  string  `json:"old_pw"`
	NewPw  string  `json:"new_pw"`
}

func UserPwChange(req *Req, resp *Resp) {

	userId , _, isLogin, respCode := RequireLogin(req)
	if !isLogin {
		c.logger.Error("not login")
		resp.Send(respCode)
		return
	}

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

	body := &UPwChangeReq{}
	if err := json.Unmarshal(rawBody, body); err != nil {
		c.logger.Error(
			"failed to parse JSON object",
			zap.Error(err),
		)
		resp.Send(RC_E_MALFORMED)
		return
	}

	// Get Old Password
	tup, err := c.Db.GetPwByID(nil, userId) 
	if err != nil {
		c.logger.Error("",zap.Error(err))
		resp.Send(RC_USER_PW_NOT_EXIST)
		return
	}

	// Hashed New Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPw), bcrypt.DefaultCost)
	if err != nil {
		c.logger.Error("failed to hash the password")
		resp.Send(RC_USER_PW_HASHING_FAIL)
		return
	}

	// Compare Old and New Passwords
	err = bcrypt.CompareHashAndPassword([]byte(tup.Password), []byte(body.OldPw))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			c.logger.Error("password not matched", zap.Error(err))
			resp.Send(RC_USER_PW_NOT_MATCH)
		} else {
			c.logger.Error("internal server error", zap.Error(err))
			resp.Send(http.StatusInternalServerError)
		}
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

	// Update Password
	if err = c.Db.UpdatePwById(
		nil,
		userId,
		string(hashedPassword),
	); err != nil {
		c.logger.Error(
			"failed to update password",
			zap.Error(err),
		)
		resp.Send(http.StatusInternalServerError)
		return
	}

	// Audit Logs
	action := "user.pw_reset"
	err = c.Db.AuditLogs(nil, userId, action)
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

	resp.Send(RC_USER_PW_CHANGED)
}

