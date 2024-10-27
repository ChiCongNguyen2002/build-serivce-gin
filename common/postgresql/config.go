package postgresql

type PostgresqlConfig struct {
	Host     string `env:"HOST,required,notEmpty"`
	Port     int    `env:"PORT,required,notEmpty"`
	Username string `env:"USERNAME,required,notEmpty"`
	Password string `env:"PASSWORD,required,notEmpty"`
	DBName   string `env:"DBNAME,required,notEmpty"`
	LogLevel int    `env:"LOG_LEVEL,required,notEmpty"`
}
