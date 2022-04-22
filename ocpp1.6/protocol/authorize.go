package protocol

type AuthorizeRequest struct {
	IdTag IdToken `json:"idTag" validate:"required,max=20"`
}

func (r *AuthorizeRequest) Reset() {
	r.IdTag = ""
}

func (AuthorizeRequest) Action() string {
	return AuthorizeName
}

type AuthorizeResponse struct {
	IdTagInfo IdTagInfo `json:"idTagInfo" validate:"required"`
}

func (AuthorizeResponse) Action() string {
	return AuthorizeName
}

func (r *AuthorizeResponse) Reset() {
	r.IdTagInfo = IdTagInfo{}
}
