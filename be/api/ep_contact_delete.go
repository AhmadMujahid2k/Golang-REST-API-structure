package api

import (
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func ContactDelete(req *Req, resp *Resp) {

	_, _, isLogin, respCode := RequireLogin(req)
	if !isLogin {
		c.logger.Error("not login")
		resp.Send(respCode)
		return
	}

	contact_id := req.URL.Query().Get("contact_id")
	if contact_id == "" {
		resp.Send(http.StatusBadRequest)
		return
	}

	contactUUID, err := uuid.Parse(contact_id)
	if err != nil {
		c.logger.Error("invalid UUID format", zap.Error(err))
		resp.Send(http.StatusBadRequest)
		return
	}

	// Delete Contact by ID 
	err = c.Db.DeleteContact(nil, contactUUID)
	if err != nil {
		c.logger.Error("failed to delete accounts", zap.Error(err))
		resp.Send(http.StatusInternalServerError)
		return
	}

	resp.SendData(RC_CONTACT_DELETED, "Successfully Deleted")
}
