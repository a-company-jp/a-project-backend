package config

type (
	Application struct {
		Server Server `yaml:"server"`
		GCS    GCS    `yaml:"gcs"`
	}
	Server struct {
		OnProduction bool     `yaml:"on_production"`
		Frontend     Frontend `yaml:"frontend"`
		Backend      Backend  `yaml:"backend"`
	}
	Frontend struct {
		Protocol string `yaml:"protocol"`
		Domain   string `yaml:"domain"`
		Port     string `yaml:"port"`
	}
	Backend struct {
		Protocol string `yaml:"protocol"`
		Domain   string `yaml:"domain"`
		Port     string `yaml:"port"`
	}
	//	Concatenating protocol + domain + port should form a valid URL

	GCS struct {
		BucketName string `yaml:"bucket_name"`
	}
)
