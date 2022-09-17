package config

import (
	"fmt"
	"sync"

	"github.com/devesh2997/consequent/errorx"
	"github.com/spf13/viper"
)

// AppConfig represents the application config that are defined in env files.
type AppConfig struct {
	Log  LogConfig `mapstructure:"log"`
	SQL  SQLConfig `mapstructure:"sql"`
	Port string    `mapstructure:"port"`
}

func (appConfig AppConfig) Validate() error {
	if err := appConfig.Log.Validate("log"); err != nil {
		return err
	}
	if appConfig.Port == "" {
		return errorx.NewValidationError(2, "(appconfig)port not found")
	}
	if err := appConfig.SQL.Validate(); err != nil {
		return err
	}

	return nil
}

// LogConfig represents logger handler
// Logger has many parameters can be set or changed. Currently, only three are listed here. Can add more into it to
// according to your needs.
type LogConfig struct {
	// log library name
	Code string `mapstructure:"code"`
	// log level
	Level string `mapstructure:"level"`
	// show caller in log message
	EnableCaller bool `mapstructure:"enableCaller"`
}

func (logConfig LogConfig) Validate(Type string) error {
	if logConfig.Code == "" {
		message := fmt.Sprintf("(%s)code not found", Type)
		return errorx.NewValidationError(11, message)
	}
	if logConfig.Level == "" {
		message := fmt.Sprintf("(%s)level not found", Type)
		return errorx.NewValidationError(12, message)
	}
	if !logConfig.EnableCaller {
		message := fmt.Sprintf("(%s)enableCaller not found", Type)
		return errorx.NewValidationError(13, message)
	}

	return nil
}

// SQLConnConfig represents config values for SQL data store.
type SQLConnConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	DB       string `mapstructure:"db"`
}

func (sqlConnConfig SQLConnConfig) Validate() error {
	if sqlConnConfig.User == "" {
		return errorx.NewValidationError(331, "(sqlconnconfig)user not found")
	}
	if sqlConnConfig.Password == "" {
		return errorx.NewValidationError(332, "(sqlconnconfig)password not found")
	}
	if sqlConnConfig.Host == "" {
		return errorx.NewValidationError(333, "(sqlconnconfig)host not found")
	}
	if sqlConnConfig.DB == "" {
		return errorx.NewValidationError(333, "(sqlconnconfig)db not found")
	}
	return nil
}

type SQLConfig struct {
	DriverName           string          `mapstructure:"driverName"`
	DataSourceNameFormat string          `mapstructure:"dataSourceNameFormat"`
	Master               SQLConnConfig   `mapstructure:"master"`
	Slaves               []SQLConnConfig `mapstructure:"slaves"`
}

func (sqlConfig SQLConfig) Validate() error {
	if sqlConfig.DriverName == "" {
		return errorx.NewValidationError(31, "(sqlconfig)drivername not found")
	}
	if sqlConfig.DataSourceNameFormat == "" {
		return errorx.NewValidationError(32, "(sqlconfig)datasourcenameformat not found")
	}
	if err := sqlConfig.Master.Validate(); err != nil {
		return err
	}

	for _, slaveConfig := range sqlConfig.Slaves {
		if err := slaveConfig.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Config is ...
var Config AppConfig

// LoadConfig loads config file and marshals in Config var
func LoadConfig(env string, path string) {
	var configOnce sync.Once

	configOnce.Do(func() {
		viper.SetConfigName(env + ".config")
		viper.AddConfigPath(path)
		viper.SetConfigType("yaml")
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		err := viper.Unmarshal(&Config)

		if err != nil {
			panic(err)
		}
	})
}
