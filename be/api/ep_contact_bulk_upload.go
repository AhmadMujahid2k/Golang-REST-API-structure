package api

import (
	"bufio"
	"Golang-REST-API-structure/be/lib/db"
	"mime"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ContactBulkResp struct {
	ContactNumber string     `json:"contact_number"`
	ContactID     uuid.UUID  `json:"contactId"`
}

func ContactBulkUpload(req *Req, resp *Resp) {

	userId, _, isLogin, respCode := RequireLogin(req)
	if !isLogin {
		c.logger.Error("not login")
		resp.Send(respCode)
		return
	}

	contentType := req.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		c.logger.Error("failed to parse media type", zap.Error(err))
		resp.Send(RC_CONTACT_MEDIA_PARSE_FAIL)
		return
	}

	if mediaType != "multipart/form-data" {
		resp.Send(http.StatusUnsupportedMediaType)
		return
	}

	const maxSizeInBytes = 10 * 1024 * 1024
	if err = req.ParseMultipartForm(
		maxSizeInBytes,
	); err != nil || len(req.MultipartForm.File) == 0 {
		c.logger.Error("failed to parse multi part form", zap.Error(err))
		resp.Send(http.StatusBadRequest)
		return
	}

	file, _, err := req.FormFile("file")
	if err != nil {
		c.logger.Error("No file found in the request", zap.Error(err))
	    return
	}
	defer file.Close()

	uploadedContacts := make([]*ContactBulkResp, 0)

	// Scan each contact line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		contactID := uuid.New()
        contact := db.Contact{
			ID:             contactID,
            UserID:         userId,
            ContactNumber:  line,
        }

		// Upload Contacts
		err = c.Db.UploadContact(nil, &contact)
		if err != nil {
			c.logger.Error(
				"failed to upload accounts",
				zap.Error(err),
			)
			resp.Send(http.StatusInternalServerError)
			return
		}

		uploadedContacts = append(uploadedContacts, &ContactBulkResp{
			ContactNumber: contact.ContactNumber,
			ContactID:     contact.ID,
		})
	}

	if err := scanner.Err(); err != nil {
		c.logger.Error(
			"Error reading file:",
			zap.Error(err),
		)
		return
	}

	resp.SendData(RC_CONTACT_UPLOAD, uploadedContacts)
}
