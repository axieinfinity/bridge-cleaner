package configs

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

var AppConfig Config

type Config struct {
	DB       Postgres
	Cleaner  Cleaner
	LogLevel int `default:"3" envconfig:"LOG_LEVEL"`
}

type Postgres struct {
	Host     string `default:"localhost" envconfig:"DB_HOST"`
	User     string `default:"root" envconfig:"DB_USER"`
	Password string `default:"" envconfig:"DB_PASS"`
	DBName   string `default:"" envconfig:"DB_NAME"`
	Port     int    `default:"5432" envconfig:"DB_PORT"`
}

type Cleaner struct {
	ClearSuccessTaskScheduler string `default:"0 0 1 * *" envconfig:"CLEAR_SUCCESS_TASK_SCHEDULER"`
	ClearFailedTaskScheduler  string `default:"0 1 1 * *" envconfig:"CLEAR_FAILED_TASK_SCHEDULER"`
	ClearEventScheduler       string `default:"0 0 * * 0" envconfig:"CLEAR_EVENT_SCHEDULER"`
	CLearSuccessJobScheduler  string `default:"0 0 * * 0" envconfig:"CLEAR_SUCCESS_JOB_SCHEDULER"`
	CLearFailedJobScheduler   string `default:"0 0 * * 0" envconfig:"CLEAR_FAILED_JOB_SCHEDULER"`

	ClearSuccessTaskThreshold int64 `default:"1000" envconfig:"CLEAR_SUCCESS_TASK_THRESHOLD"`
	ClearFailedTaskThreshold  int64 `default:"1000" envconfig:"CLEAR_FAILED_TASK_THRESHOLD"`
	ClearEventThreshold       int64 `default:"1000" envconfig:"CLEAR_EVENT_THRESHOLD"`
	CLearSuccessJobThreshold  int64 `default:"1000" envconfig:"CLEAR_SUCCESS_JOB_THRESHOLD"`
	CLearFailedJobThreshold   int64 `default:"1000" envconfig:"CLEAR_FAILED_JOB_THRESHOLD"`

	ClearSuccessTaskExpiration int64 `default:"7776000" envconfig:"CLEAR_SUCCESS_TASK_EXPIRATION"`
	ClearFailedTaskExpiration  int64 `default:"31104000" envconfig:"CLEAR_FAILED_TASK_EXPIRATION"`
	ClearEventExpiration       int64 `default:"604800" envconfig:"CLEAR_EVENT_EXPIRATION"`
	ClearSuccessJobExpiration  int64 `default:"604800" envconfig:"CLEAR_SUCCESS_JOB_EXPIRATION"`
	ClearFailedJobExpiration   int64 `default:"604800" envconfig:"CLEAR_FAILED_JOB_EXPIRATION"`
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
	return &AppConfig, nil
}
