package config

import "skillui/internal/process"

type AppConfig struct {
	Locale          string   `json:"locale"`
	AutoStart       bool     `json:"autoStart"`
	LogDir          string   `json:"logDir"`
	MaxLogLines     int      `json:"maxLogLines"`
	MaxLogFiles     int      `json:"maxLogFiles"`
	MaxRestart      int      `json:"maxRestart"`
	RestartPolicy   string   `json:"restartPolicy"`
	DeviceUUID      string   `json:"deviceUUID"`
	SkillDir        string   `json:"skillDir"`
	AutoSyncToolIDs []string `json:"autoSyncToolIDs"`
	// ToolPaths 记录用户手动指定的工具规则目录（toolID -> 目录路径）。
	// 自动扫描识别不到时，可手动指定以覆盖默认检测结果。
	ToolPaths map[string]string    `json:"toolPaths"`
	Processes []process.Definition `json:"processes"`
}

func DefaultConfig() AppConfig {
	return AppConfig{
		Locale:          "zh",
		AutoStart:       false,
		LogDir:          "logs",
		MaxLogLines:     1000,
		MaxLogFiles:     5,
		MaxRestart:      5,
		RestartPolicy:   "on_failure",
		SkillDir:        "",
		AutoSyncToolIDs: []string{},
		ToolPaths:       map[string]string{},
	}
}
