package api

import (
	"net/http"

	"go.uber.org/zap"
)

func ContactList(req *Req, resp *Resp) {

	userId, _, isLogin, respCode := RequireLogin(req)
	if !isLogin {
		c.logger.Error("not login")
		resp.Send(respCode)
		return
	}

	// Get Contact by Account ID
	contactTups, err := c.Db.GetContactList(nil, userId)
	if err != nil {
		c.logger.Error("failed to get contact", zap.Error(err))
		resp.Send(http.StatusInternalServerError)
		return
	}

	resp.SendData(RC_CONTACT_DETAILS, contactTups)
}
