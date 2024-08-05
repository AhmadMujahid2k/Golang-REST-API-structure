package api

import (
	"Golang-REST-API-structure/be/lib/db"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ContactReq struct {
	ContactNumber string `json:"contact_number"`
}

type ContactResp struct {
	ContactNumber string     `json:"contact_number"`
	ContactID     uuid.UUID  `json:"contactId"`
}

func ContactUpload(req *Req, resp *Resp) {

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

	body := &ContactReq{}
	if err := json.Unmarshal(rawBody, body); err != nil {
		c.logger.Error(
			"failed to parse JSON object",
			zap.Error(err),
		)
		resp.Send(RC_E_MALFORMED)
		return
	}


	if len(body.ContactNumber) != 11 && len(body.ContactNumber) != 12 {
		resp.Send(RC_CONTACT_INVALID)
		return
	}

	contactID := uuid.New()
	contact := db.Contact{
		ID:             contactID,
        UserID:         userId,
        ContactNumber:  body.ContactNumber,
    }

	// Upload contact to DB
	err = c.Db.UploadContact(nil, &contact)
	if err != nil {
		c.logger.Error(
			"failed to upload accounts",
			zap.Error(err),
		)
		resp.Send(http.StatusInternalServerError)
		return
	}

	data := ContactResp{
		ContactNumber:  contact.ContactNumber,
		ContactID:      contact.ID,
	}

	resp.SendData(RC_CONTACT_UPLOAD, data)
}
