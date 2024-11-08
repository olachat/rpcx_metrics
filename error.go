package prom

func NewErrorInvalidParam(field string, reason string) (o *ErrorInvalidParam) {
	o = &ErrorInvalidParam{
		field:  field,
		reason: reason,
	}
	return
}

type ErrorInvalidParam struct {
	field  string
	reason string
}

func (o *ErrorInvalidParam) Error() string {
	return o.field + ": " + o.reason
}
