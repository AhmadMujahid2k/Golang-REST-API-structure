package api

import (
	"Golang-REST-API-structure/be"
	"net/http"

	"Golang-REST-API-structure/be/lib/db"

	"github.com/google/uuid"
)

func RequireLogin(req *Req) (uuid.UUID, string, bool, be.RespCode) {

	var userId uuid.UUID
	var role string
	var authSessVal db.SessionVals
	sessVal := req.Session.Values[0]
	if sessVal != nil {
		authSessVal = sessVal.(db.SessionVals)
		userId = authSessVal.Id
		role = authSessVal.Role
	}

	if req.Session.IsNew {
		c.logger.Error("unauthorize access")
		return userId, role, false, http.StatusUnauthorized
	}

	return userId, role, true, 0
}
