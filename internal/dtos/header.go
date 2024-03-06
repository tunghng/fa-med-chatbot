package dtos

// Headers definition.
const (
	HeaderAuthorization     = "Authorization"
	HeaderXRequestID        = "X-Request-ID"
	HeaderXFcmToken         = "X_FCM_TOKEN"
	GinContextBasicUsername = "basic_username"
)

// Context key to transfer to next layer
const (
	ContextResponse = "x-response"
	ID              = "id"
	UserId          = "user_id"
	Avatar          = "avatar"
	FullName        = "full_name"
)

// Authorization method header
const (
	BearerAuth = "bearer"
	BasicAuth  = "basic"
)
