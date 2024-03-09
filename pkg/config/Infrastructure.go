package config

type Infrastructure struct {
	MySQLDB     MySQLDB     `yaml:"mysql"`
	GoogleCloud GoogleCloud `yaml:"google_cloud"`
}

type GoogleCloud struct {
	ProjectID           string `yaml:"project_id"`
	UseCredentialsFile  bool   `yaml:"use_credentials_file"`
	CredentialsFilePath string `yaml:"credentials_file_path"`
}

type MySQLDB struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}
