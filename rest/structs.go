package rest

type RestResponse struct {
	Error interface{}
	Value interface{}
}

type RestErrorResponse struct {
	Msg    string
	Status int
}

func (e RestErrorResponse) Error() string {
	return e.Msg
}

func GetRestErr(err error) *RestErrorResponse {
	var r *RestErrorResponse = nil
	if err != nil {
		return &RestErrorResponse{
			Msg:    err.Error(),
			Status: 0,
		}
	}
	return r
}
