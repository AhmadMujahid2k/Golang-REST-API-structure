package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type UProfileUpdReq struct {
	Name    string 	`json:"fullname"`
	Gender  string 	`json:"gender"`
	Dob     string  `json:"dob"`
	Phone   string  `json:"phone"`
}

func UserProfileUpd(req *Req, resp *Resp) {

	userId, _, isLogin, respCode := RequireLogin(req)
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

	body := &UProfileUpdReq{}
	if err := json.Unmarshal(rawBody, body); err != nil {
		c.logger.Error(
			"failed to parse JSON object",
			zap.Error(err),
		)
		resp.Send(RC_E_MALFORMED)
		return
	}

	// Parse Date of birth
	dob, err := time.Parse(time.RFC3339, body.Dob)
	if err != nil {
		c.logger.Error(
			"failed to parse dob",
			zap.Error(err),
		)
		resp.Send(RC_E_MALFORMED)
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

	// Update Profile
	err = c.Db.UpdateProfileById(nil, userId, body.Phone, dob, body.Gender, body.Name)
	if err != nil {
		c.logger.Error(
			"failed to update user profile",
			zap.Error(err),
		)
		resp.Send(http.StatusInternalServerError)
		return
	}

	// Audit Logs
	action := "user.profile_update"
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
	
	resp.Send(RC_USER_PROFILE_UPDATE)
}
