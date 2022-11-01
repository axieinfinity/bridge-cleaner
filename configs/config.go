package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/kelseyhightower/envconfig"
)

var AppConfig Config

type Config struct {
	DB                Postgres
	LogLevel          int    `default:"3" envconfig:"LOG_LEVEL"`
	CronJobConfigPath string `envconfig:"CRON_JOB_CONFIG_PATH"`
	CronJob           map[string]*Cronjob
}

type Postgres struct {
	Host     string `default:"localhost" envconfig:"DB_HOST"`
	User     string `default:"root" envconfig:"DB_USERNAME"`
	Password string `default:"" envconfig:"DB_PASSWORD"`
	DBName   string `default:"" envconfig:"DB_NAME"`
	Port     int    `default:"5432" envconfig:"DB_PORT"`
}

type CleanerV2 struct {
	DB Postgres `json:"database"`
}

type Cronjob struct {
	Cron string `default:"0 0 1 * *" json:"cron"`
	// Conditions
	Threshold  int64 `default:"1000" json:"threshold"`
	Expiration int64 `default:"1000" json:"expiration"`
}

func (p *Postgres) ConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", p.Host, p.User, p.Password, p.DBName, p.Port)
}

func (p *Postgres) ConnectionStringURL() string {
	return fmt.Sprintf("postgres://%v:%v/%v?sslmode=disable&user=%v&password=%v", p.Host, p.Port, p.DBName, p.User, p.Password)
}

func New() (*Config, error) {
	if err := envconfig.Process("", &AppConfig); err != nil {
		return nil, err
	}

	AppConfig.CronJob = map[string]*Cronjob{}
	if AppConfig.CronJobConfigPath != "" {
		data, _ := ioutil.ReadFile(AppConfig.CronJobConfigPath)
		if err := json.Unmarshal(data, &AppConfig.CronJob); err != nil {
			panic(err)
		}
	}

	return &AppConfig, nil
}
