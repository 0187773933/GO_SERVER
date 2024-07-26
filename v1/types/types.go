package types

import (
	logger_types "github.com/0187773933/Logger/v1/types"
)

type URLS struct {
	Local  string `yaml:"local"`
	Private string `yaml:"private"`
	Public  string `yaml:"public"`
	AdminLogin string `yaml:"admin_login"`
	AdminPrefix string `yaml:"admin_prefix"`
	Login   string `yaml:"login"`
	Prefix  string `yaml:"prefix"`
}

type CookieInfo struct {
	Name    string `yaml:"name"`
	Message string `yaml:"message"`
}

type Cookie struct {
	User  CookieInfo `yaml:"user"`
	Admin CookieInfo `yaml:"admin"`
	Secret  string `yaml:"secret"`
}

type Docker struct {
	Name string `yaml:"name"`
}

type Git struct {
	URL string `yaml:"url"`
	SSHURL string `yaml:"ssh_url"`
}

type Go struct {
	Version string `yaml:"version"`
	OS string `yaml:"os"`
	Arch string `yaml:"arch"`
}

type Bolt struct {
	Path string `yaml:"path"`
	Prefix string `yaml:"prefix"`
}

type Redis struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Password string `yaml:"password"`
	Number int `yaml:"number"`
	Prefix string `yaml:"prefix"`
	Enabled bool `yaml:"enabled"`
}

type Creds struct {
	APIKey         string `yaml:"api_key"`
	AdminUsername  string `yaml:"admin_username"`
	AdminPassword  string `yaml:"admin_password"`
	EncryptionKey  string `yaml:"encryption_key"`
	OpenAIKey      string `yaml:"openai_key"`
}

type Config struct {
	Name      string `yaml:"name"`
	Port      string `yaml:"port"`
	URLS      URLS   `yaml:"urls"`
	Cookie    Cookie `yaml:"cookie"`
	Creds     Creds  `yaml:"creds"`
	TimeZone  string `yaml:"time_zone"`
	AllowOrigins []string `yaml:"allow_origins"`
	SaveFilesPath string `yaml:"save_files_path"`
	Bolt Bolt `yaml:"bolt"`
	Redis Redis `yaml:"redis"`
	Docker Docker `yaml:"docker"`
	Git Git `yaml:"git"`
	Go Go `yaml:"go"`
	Log logger_types.ConfigFile `yaml:"log"`
}

type ConfigGeneric map[string]interface{}