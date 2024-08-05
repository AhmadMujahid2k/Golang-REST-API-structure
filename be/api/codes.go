package api

import "Golang-REST-API-structure/be"

const (
	// Reserved
	RC_E_NO_BODY                be.RespCode = 2999
	RC_E_MALFORMED              be.RespCode = 2998
	RC_SDK_ERROR                be.RespCode = 2997

	// POST /api/ahmadmujahid/v1/user/signup
	RC_USER_SIGNUP              be.RespCode = 1000
	RC_USER_EMAIL_NOT_FOUND     be.RespCode = 2001

	// POST /api/ahmadmujahid/v1/user/login
	RC_USER_LOGIN               be.RespCode = 1000
	RC_USER_NO_EMAIL            be.RespCode = 2010
	RC_USER_NOT_EXIST           be.RespCode = 2011
	RC_USER_INVALID_CREDENTIALS be.RespCode = 2012
	RC_USER_SAVE_SESSION_FAIL   be.RespCode = 2013

	// POST /api/ahmadmujahid/v1/user/pw_change
	RC_USER_PW_CHANGED          be.RespCode = 1000
	RC_USER_PW_NOT_EXIST        be.RespCode = 2014
	RC_USER_PW_HASHING_FAIL     be.RespCode = 2015
	RC_USER_PW_NOT_MATCH        be.RespCode = 2016
	RC_USER_PW_UPDATE_FAIL      be.RespCode = 2017

	//POST /api/ahmadmujahid/v1/user/profile_upd
	RC_USER_PROFILE_UPDATE      be.RespCode = 1000
	
	// POST /api/ahmadmujahid/v1/contact/bulk
	// POST /api/ahmadmujahid/v1/contact/bulk/upload
	RC_CONTACT_UPLOAD           be.RespCode = 1000
	RC_CONTACT_MEDIA_PARSE_FAIL be.RespCode = 2024
	RC_CONTACT_INVALID          be.RespCode = 2025

	// DELETE /api/ahmadmujahid/v1/contact/delete
	RC_CONTACT_DELETED          be.RespCode = 1000

	// GET /api/ahmadmujahid/v1/contact/list
	RC_CONTACT_DETAILS         be.RespCode = 1000

	// GET /api/ahmadmujahid/v1/auditlogs/list
	RC_AUDIT_LOG_DETAILS       be.RespCode = 1000  
)
