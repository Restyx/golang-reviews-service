package messagehandler

type Config struct {
	PgUser      string `toml:"postgres_user"`
	PgPassword  string `toml:"postgres_pass"`
	PgHost      string `toml:"postgres_host"`
	PgPort      string `toml:"postgres_port"`
	PgDB        string `toml:"postgres_db"`
	RmqUser     string `toml:"rabbitmq_user"`
	RmqPassword string `toml:"rabbitmq_pass"`
	RmqHost     string `toml:"rabbitmq_host"`
	RmqPort     string `toml:"rabbitmq_port"`
}

func NewConfig() *Config {
	return &Config{
		PgUser:      "postgres",
		PgPassword:  "postgres",
		PgHost:      "localhost",
		PgPort:      "5432",
		PgDB:        "reviews",
		RmqUser:     "guest",
		RmqPassword: "guest",
		RmqHost:     "localhost",
		RmqPort:     "5672",
	}
}
