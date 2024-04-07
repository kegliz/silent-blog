package config

type (
	configVar struct {
		Type    configVarType
		Default interface{}
		EnvVar  string
	}
	configVarType string
)

var (
	stringType configVarType = "string"
	intType    configVarType = "int"
	boolType   configVarType = "bool"
)

var configVars = map[string]configVar{
	"port": {
		Type:    intType,
		Default: 3000,
		EnvVar:  "PORT",
	},
	"gracefulshutdowntimeout": {
		Type:    intType,
		Default: 10,
		EnvVar:  "GRACEFULSHUTDOWNTIMEOUT",
	},
	"debug": {
		Type:    boolType,
		Default: false,
		EnvVar:  "DEBUG",
	},
	"localonly": {
		Type:    boolType,
		Default: true,
		EnvVar:  "LOCALONLY",
	},
	"tls": {
		Type:    boolType,
		Default: false,
		EnvVar:  "TLS",
	},
	"domain": {
		Type:    stringType,
		Default: "localhost",
		EnvVar:  "DOMAIN",
	},
	"posts.file": {
		Type:    stringType,
		Default: "posts.json",
		EnvVar:  "POSTS_FILE",
	},
	"posts.mddir": {
		Type:    stringType,
		Default: "md",
		EnvVar:  "POSTS_MDDIR",
	},
}
