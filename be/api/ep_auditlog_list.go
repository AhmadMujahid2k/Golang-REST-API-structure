package api

import (
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

const AUDIT_LOG_LIMIT = 10

func AuditLogList(req *Req, resp *Resp) {

	_, _, isLogin, respCode := RequireLogin(req)
	if !isLogin {
		c.logger.Error("not login")
		resp.Send(respCode)
		return
	}

	pageNumberStr := req.URL.Query().Get("page_no")
	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil {
	    pageNumber = 1
	}

	limitStr := req.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if limit > AUDIT_LOG_LIMIT || err != nil {
		limit = AUDIT_LOG_LIMIT
	}

	offset := (pageNumber - 1) * limit

	// Get Audit Logs with pagination
	logsTups, err := c.Db.GetAuditLogList(nil, offset, limit)
	if err != nil {
		c.logger.Error("failed to get logs", zap.Error(err))
		resp.Send(http.StatusInternalServerError)
		return
	}

	resp.SendData(RC_AUDIT_LOG_DETAILS, logsTups)
}
