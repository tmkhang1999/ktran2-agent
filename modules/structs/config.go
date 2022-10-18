package structs

import "time"

type Config struct {
	Url    string        `yaml:"URL"`
	Method string        `yaml:"METHOD"`
	Query  string        `yaml:"QUERY"`
	Time   time.Duration `yaml:"TIME"`

	TableName string `yaml:"TABLENAME"`
	Region    string `yaml:"REGION"`

	AccessKey string
}
