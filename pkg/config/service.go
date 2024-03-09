package config

type (
	Service struct {
		Authentication Authentication `yaml:"authentication"`
	}
	Authentication struct {
		JwtSecret string `yaml:"jwt_secret"`
	}
)
