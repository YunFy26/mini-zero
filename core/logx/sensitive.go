package logx

type Sensitive interface {
	MaskSensitive() any
}

func maskSensitive(v any) any {
	if s, ok := v.(Sensitive); ok {
		return s.MaskSensitive()
	}
	return v
}
