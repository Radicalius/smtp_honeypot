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
		Value      T   `json:"value"`
		Percentage int `json:"percentage"`
	} `json:"rollout"`
}

func (v *VarConfig[T]) GetValue(ip string) T {
	if v.Rollout != nil {
		h := fnv.New32a()
		h.Write([]byte(ip))
		hashInt := h.Sum32()
		if hashInt%100 < uint32(v.Rollout.Percentage) {
			return v.Rollout.Value
		}
	}

	return v.Default
}

type Settings struct {
	ImmediateTlsWindow VarConfig[int] `json:"immediateTlsWindow"`
}

func GetSettings() *Settings {
	var settings Settings
	err := json.Unmarshal(settingsFile, &settings)
	if err != nil {
		log.Fatalf("error parsing config: %s\n", err.Error())
	}

	return &settings
}
