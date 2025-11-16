package logx

type (
	Writer interface {
		Alert(v any)
		Close() error
		Debug(v any, fields ...LogField)
		Error(v any, fields ...LogField)
		Info(v any, fields ...LogField)
		Severe(v any)
		Slow(v any, fields ...LogField)
		Stack(v any)
		Stat(v any, fields ...LogField)
	}
)
