package config

type Config struct {
	Server ServerConfig `mapstructure:"server" json:"server" yaml:"server"`
	Mysql  MysqlConfig  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis  RedisConfig  `mapstructure:"redis" json:"redis" yaml:"redis"`
	Chat   ChatConfig   `mapstructure:"chat" json:"chat" yaml:"chat"`
	Upload UploadConfig `mapstructure:"upload" json:"upload" yaml:"upload"`
	Log    LogConfig    `mapstructure:"log" json:"log" yaml:"log"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
	Mode string `mapstructure:"mode" json:"mode" yaml:"mode"`
}

type MysqlConfig struct {
	Path         string `mapstructure:"path" json:"path" yaml:"path"`
	Port         int    `mapstructure:"port" json:"port" yaml:"port"`
	Config       string `mapstructure:"config" json:"config" yaml:"config"`
	DbName       string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`
	Username     string `mapstructure:"username" json:"username" yaml:"username"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"`
}

type RedisConfig struct {
	Db       int    `mapstructure:"db" json:"db" yaml:"db"`
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Prefix   string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
}

type ChatConfig struct {
	FaqMatchThreshold  float64  `mapstructure:"faq_match_threshold" json:"faq_match_threshold" yaml:"faq_match_threshold"`
	MaxUploadSize      int64    `mapstructure:"max_upload_size" json:"max_upload_size" yaml:"max_upload_size"`
	AllowedUploadTypes []string `mapstructure:"allowed_upload_types" json:"allowed_upload_types" yaml:"allowed_upload_types"`
	SessionTimeout     int      `mapstructure:"session_timeout" json:"session_timeout" yaml:"session_timeout"`
	AgentAssignment    string   `mapstructure:"agent_assignment" json:"agent_assignment" yaml:"agent_assignment"`
}

type UploadConfig struct {
	Path    string `mapstructure:"path" json:"path" yaml:"path"`
	MaxSize int64  `mapstructure:"max_size" json:"max_size" yaml:"max_size"`
}

type LogConfig struct {
	Level      string `mapstructure:"level" json:"level" yaml:"level"`               // debug/info/warn/error
	FilePath   string `mapstructure:"file_path" json:"file_path" yaml:"file_path"`   // logs/app.log
	MaxSize    int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"`      // MB, rotate when exceeded
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"` // keep N backups
	MaxAge     int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`         // max days to keep
	Compress   bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
	Console    bool   `mapstructure:"console" json:"console" yaml:"console"`         // also output to console
}
