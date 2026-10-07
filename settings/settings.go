package settings

import (
	_ "embed"
	"encoding/json"
	"hash/fnv"
	"log"
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
}

func (v *VarConfig[T]) GetValue(ip string) T {
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
	ImmediateTlsWindow       int
	AuthRetriesBeforeSuccess int
}

type Settings struct {
	ImmediateTlsWindow       VarConfig[int] `json:"immediateTlsWindow"`
	AuthRetriesBeforeSuccess VarConfig[int] `json:"authRetriesBeforeSuccess"`
}

func (s *Settings) ToEvaluatedSettings(ip string) *EvaluatedSettings {
	return &EvaluatedSettings{
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
