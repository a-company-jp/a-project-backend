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
	Protocol   string `yaml:"protocol"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	UnixSocket string `yaml:"unix_socket"`
	User       string `yaml:"username"`
	Password   string `yaml:"password"`
	DBName     string `yaml:"db_name"`
}
