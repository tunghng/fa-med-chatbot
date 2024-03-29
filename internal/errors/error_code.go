package errors

/*
Example error code is 5000102:
- 500 is HTTP status code (400, 401, 403, 500, ...)
- 01 is module represents for each handlers
  - 00 for common error for all handlers
  - 01 for health check handlers

- 02 is actual error code, just auto increment and start at 1
*/
const ERR_3RD string = "3rd-code"

var (
	// Errors of module common
	// Format: ErrCommon<ERROR_NAME> = xxx00yy
	ErrCommonInternalServer    = ErrorCode("50000001")
	ErrCommonInvalidRequest    = ErrorCode("40000001")
	ErrCommonBindRequestError  = ErrorCode("40000002")
	ErrCommonExpiredToken      = ErrorCode("40100006")
	ErrAuthorizedNotPermission = ErrorCode("40100108")

	ErrUserAlreadyExist      = ErrorCode("40000103")
	ErrInvalidOldPassword    = ErrorCode("40000104")
	ErrInvalidOldNewPassword = ErrorCode("40000105")
	ErrCommonUnauthorized    = ErrorCode("40100103")

	ErrEmailNotExist    = ErrorCode("40100109")
	ErrPhoneNotExist    = ErrorCode("40100104")
	ErrGmailNotExist    = ErrorCode("40100105")
	ErrFacebookNotExist = ErrorCode("40100106")
	ErrAppleNotExist    = ErrorCode("40100107")

	ErrorFileUploadNotNull  = ErrorCode("40000201")
	ErrorFileUploadMaximum  = ErrorCode("40000202")
	ErrorMimeTypeNotSupport = ErrorCode("40000203")
	ErrorUploadImageDrive   = ErrorCode("40000204")

	ErrorFlashcardDetailNull = ErrorCode("40300301")
	ErrorFlashcardPrivate    = ErrorCode("40000302")
	ErrorFlashcardNeedTopic  = ErrorCode("40000303")
)
