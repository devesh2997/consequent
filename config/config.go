package config

import (
	"errors"
	"fmt"
	"sync"

	"github.com/devesh2997/consequent/errorx"
	"github.com/spf13/viper"
)

// AppConfig represents the application config that are defined in env files.
type AppConfig struct {
	Log           LogConfig     `mapstructure:"log"`
	SQL           SQLConfig     `mapstructure:"sql"`
	Factor2Config Factor2Config `mapstructure:"2factor"`
	Port          string        `mapstructure:"port"`
}

func (appConfig AppConfig) Validate() error {
	if err := appConfig.Log.Validate("log"); err != nil {
		return err
	}
	if appConfig.Port == "" {
		return errorx.NewSystemError(-1, errors.New("(appconfig)port not found"))
	}
	if err := appConfig.SQL.Validate(); err != nil {
		return err
	}
	if err := appConfig.Factor2Config.Validate(); err != nil {
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
		return errorx.NewSystemError(-1, errors.New(message))
	}
	if logConfig.Level == "" {
		message := fmt.Sprintf("(%s)level not found", Type)
		return errorx.NewSystemError(-1, errors.New(message))
	}
	if !logConfig.EnableCaller {
		message := fmt.Sprintf("(%s)enableCaller not found", Type)
		return errorx.NewSystemError(-1, errors.New(message))
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
		return errorx.NewSystemError(-1, errors.New("(sqlconnconfig)user not found"))
	}
	if sqlConnConfig.Password == "" {
		return errorx.NewSystemError(-1, errors.New("(sqlconnconfig)password not found"))
	}
	if sqlConnConfig.Host == "" {
		return errorx.NewSystemError(-1, errors.New("(sqlconnconfig)host not found"))
	}
	if sqlConnConfig.DB == "" {
		return errorx.NewSystemError(-1, errors.New("(sqlconnconfig)db not found"))
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
		return errorx.NewSystemError(-1, errors.New("(sqlconfig)drivername not found"))
	}
	if sqlConfig.DataSourceNameFormat == "" {
		return errorx.NewSystemError(-1, errors.New("(sqlconfig)datasourcenameformat not found"))
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

type Factor2Config struct {
	APIKey          string `mapstructure:"api_key"`
	OTPTemplateName string `mapstructure:"otp_template_name"`
}

func (factorConfig Factor2Config) Validate() error {
	if factorConfig.APIKey == "" {
		return errorx.NewSystemError(-1, errors.New("2factor.in api key not found"))
	}
	if factorConfig.OTPTemplateName == "" {
		return errorx.NewSystemError(-1, errors.New("2factor.in otp template name not found"))
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
