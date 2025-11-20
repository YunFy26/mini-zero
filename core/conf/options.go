package conf

type (
	// Option defines the method to customize the config options.
	Option func(opt *options)

	options struct {
		env bool
	}
)

// UseEnv sets whether to use environment variables for configuration
func UseEnv(env bool) Option {
	return func(opt *options) {
		opt.env = env
	}
}
