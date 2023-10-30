package merrors

const (
	ErrMsgSuccess = "success"

	ErrMsgRequestBindFailed = "request_bind_failed"
	ErrMsgDataNotSupport    = "data_not_support"
)

var ErnoMsgMap = map[int]string{
	ErnoSuccess:           ErrMsgSuccess,
	ErnoRequestBindFailed: ErrMsgRequestBindFailed,
	ErnoDataNotSupport:    ErrMsgDataNotSupport,
}
