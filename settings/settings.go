package settings

import (
	_ "embed"
	"encoding/json"
	"hash/fnv"
	"log"
	"os"
	"strconv"
)

//go:embed settings.json
var settingsFile []byte

type VarConfig[T any] struct {
	Default T `json:"default"`
	Rollout *struct {
		Value      T      `json:"value"`
		Percentage int    `json:"percentage"`
		Key        string `json:"key"`
	} `json:"rollout"`
	Envvar string `json:"envvar"`
}

func (v *VarConfig[T]) GetValue(ip string) T {
	env := os.Getenv(v.Envvar)

	var t T
	if env != "" {
		switch any(t).(type) {
		case string:
			return any(env).(T)
		case int:
			i, err := strconv.ParseInt(env, 10, 64)
			if err == nil {
				log.Fatalf("error parsing int variable: %s; %s\n", env, err.Error())
			}
			return any(i).(T)
		case float64:
			i, err := strconv.ParseFloat(env, 64)
			if err == nil {
				log.Fatalf("error parsing float variable: %s; %s\n", env, err.Error())
			}
			return any(i).(T)
		case bool:
			i, err := strconv.ParseBool(env)
			if err == nil {
				log.Fatalf("error parsing bool variable: %s; %s\n", env, err.Error())
			}
			return any(i).(T)
		default:
			return t
		}
	}

	if v.Rollout != nil {
		h := fnv.New32a()
		h.Write([]byte(v.Rollout.Key + "|" + ip))
		hashInt := h.Sum32()
		if hashInt%100 < uint32(v.Rollout.Percentage) {
			return v.Rollout.Value
		}
	}

	return v.Default
}

type EvaluatedSettings struct {
	HoneypotIdentity         string
	ImmediateTlsWindow       int
	AuthRetriesBeforeSuccess int
}

type Settings struct {
	Ports               VarConfig[string] `json:"ports"`
	LogRoot             VarConfig[string] `json:"logRoot"`
	CertPath            VarConfig[string] `json:"certPath"`
	PublishUrl          VarConfig[string] `json:"publishUrl"`
	PublishToken        VarConfig[string] `json:"publishToken"`
	OperatorIpWhitelist VarConfig[string] `json:"operatorIpWhitelist"`

	HoneypotIdentity         VarConfig[string] `json:"honeypotIdentity"`
	ImmediateTlsWindow       VarConfig[int]    `json:"immediateTlsWindow"`
	AuthRetriesBeforeSuccess VarConfig[int]    `json:"authRetriesBeforeSuccess"`
}

func (s *Settings) ToEvaluatedSettings(ip string) *EvaluatedSettings {
	return &EvaluatedSettings{
		HoneypotIdentity:         s.HoneypotIdentity.GetValue(ip),
		ImmediateTlsWindow:       s.ImmediateTlsWindow.GetValue(ip),
		AuthRetriesBeforeSuccess: s.AuthRetriesBeforeSuccess.GetValue(ip),
	}
}

func GetSettings() *Settings {
	var settings Settings
	err := json.Unmarshal(settingsFile, &settings)
	if err != nil {
		log.Fatalf("error parsing config: %s\n", err.Error())
	}

	return &settings
}
