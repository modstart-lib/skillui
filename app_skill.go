package main

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"time"
)

// SkillMeta represents metadata of an installed skill
type SkillMeta struct {
	Name         string   `json:"name"`
	Title        string   `json:"title"`
	TitleEn      string   `json:"titleEn"`
	TitleZh      string   `json:"titleZh"`
	Description  string   `json:"description"`
	DescEn       string   `json:"descEn"`
	DescZh       string   `json:"descZh"`
	Owner        string   `json:"owner"`
	Version      string   `json:"version"`
	Tags         []string `json:"tags"`
	IsMarket     bool     `json:"isMarket"`
	MarketID     int      `json:"marketId"`
	SyncedTools  []string `json:"syncedTools"`
	Location     string   `json:"location"`
	UpdatedAt    string   `json:"updatedAt"`
	SkillContent string   `json:"skillContent"`
}

// IDEToolInfo represents an IDE or AI coding tool detected on the system
type IDEToolInfo struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Installed     bool   `json:"installed"`
	Path          string `json:"path"`
	SkillRulesDir string `json:"skillRulesDir"`
}

// skillUIJson is the structure saved as skillui.json inside market-installed skills
type skillUIJson struct {
	MarketID    int    `json:"marketId"`
	Name        string `json:"name"`
	TitleEn     string `json:"titleEn"`
	TitleZh     string `json:"titleZh"`
	DescEn      string `json:"descEn"`
	DescZh      string `json:"descZh"`
	Owner       string `json:"owner"`
	Version     string `json:"version"`
	Md5         string `json:"md5"`
	InstalledAt string `json:"installedAt"`
}

// ideToolDef defines detection rules for an IDE tool per platform
type ideToolDef struct {
	ID            string
	Name          string
	CheckPathsMac []string
	CheckPathsWin []string
	CheckPathsLin []string
	RulesDirMac   string
	RulesDirWin   string
	RulesDirLin   string
}

// expandHome replaces leading ~ with the user's home directory
func expandHome(path string) string {
	if path == "" {
		return path
	}
	if strings.HasPrefix(path, "~/") || path == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(homeDir, path[1:])
	}
	return path
}

// getSkillDir returns the expanded skill directory path
func (a *App) getSkillDir() string {
	if a.config.SkillDir != "" {
		return expandHome(a.config.SkillDir)
	}
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".skillui", "skills")
}

// GetAutoSyncToolIDs returns the list of tool IDs with auto-sync enabled
func (a *App) GetAutoSyncToolIDs() []string {
	if a.config.AutoSyncToolIDs == nil {
		return []string{}
	}
	return a.config.AutoSyncToolIDs
}

// SetAutoSyncToolIDs saves the list of tool IDs with auto-sync enabled
func (a *App) SetAutoSyncToolIDs(ids []string) error {
	if ids == nil {
		ids = []string{}
	}
	a.config.AutoSyncToolIDs = ids
	return a.store.Save(a.config)
}

// syncToInstalledTools syncs a skill to all tools that are installed AND have auto-sync enabled
func (a *App) syncToInstalledTools(skillName string) {
	autoIDs := map[string]bool{}
	for _, id := range a.config.AutoSyncToolIDs {
		autoIDs[id] = true
	}
	if len(autoIDs) == 0 {
		return
	}
	scanned, err := a.ScanIDETools()
	if err != nil {
		return
	}
	targets := make([]string, 0)
	for _, t := range scanned {
		if t.Installed && autoIDs[t.ID] {
			targets = append(targets, t.ID)
		}
	}
	if len(targets) > 0 {
		_ = a.SyncSkillToTools(skillName, targets)
	}
}

// GetSkillDir returns the current skill directory to the frontend
func (a *App) GetSkillDir() string {
	return a.getSkillDir()
}

// SetSkillDir changes the skill directory, optionally migrating existing skills
func (a *App) SetSkillDir(newDir string, migrate bool) error {
	newDir = expandHome(newDir)
	if err := os.MkdirAll(newDir, 0755); err != nil {
		return fmt.Errorf("无法创建目录: %w", err)
	}
	if migrate {
		oldDir := a.getSkillDir()
		if oldDir != newDir {
			if err := moveDir(oldDir, newDir); err != nil {
				return fmt.Errorf("迁移技能失败: %w", err)
			}
		}
	}
	a.config.SkillDir = newDir
	return a.store.Save(a.config)
}

// moveDir moves all immediate subdirectories from src to dst
func moveDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if err := os.Rename(srcPath, dstPath); err != nil {
			// cross-device: fallback to copy+delete
			if err2 := copyDir(srcPath, dstPath); err2 != nil {
				return err2
			}
			os.RemoveAll(srcPath)
		}
	}
	return nil
}

// copyDir recursively copies a directory
func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		dstPath := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}
		return copyFile(path, dstPath)
	})
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

// parseSkillMeta reads SKILL.md frontmatter and optionally skillui.json from a skill directory
func parseSkillMeta(dir string) SkillMeta {
	meta := SkillMeta{
		Name:     filepath.Base(dir),
		Location: dir,
	}

	// Parse SKILL.md frontmatter
	skillFile := filepath.Join(dir, "SKILL.md")
	if f, err := os.Open(skillFile); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		inFrontmatter := false
		firstLine := true
		for scanner.Scan() {
			line := scanner.Text()
			if firstLine {
				firstLine = false
				if line == "---" {
					inFrontmatter = true
					continue
				}
				break
			}
			if inFrontmatter {
				if line == "---" {
					break
				}
				if idx := strings.Index(line, ":"); idx > 0 {
					key := strings.TrimSpace(line[:idx])
					val := strings.TrimSpace(line[idx+1:])
					val = strings.Trim(val, "'\"")
					switch key {
					case "name":
						meta.Name = val
					case "title":
						meta.Title = val
					case "description":
						meta.Description = val
					case "owner":
						meta.Owner = val
					case "version":
						meta.Version = val
					}
				}
			}
		}
	}
	// Read full SKILL.md content
	if data, err := os.ReadFile(skillFile); err == nil {
		meta.SkillContent = string(data)
	}

	// Fallback: use dir name as title if empty
	if meta.Title == "" {
		meta.Title = meta.Name
	}

	// Read skillui.json if present (market-installed skill)
	skillUIFile := filepath.Join(dir, "skillui.json")
	if data, err := os.ReadFile(skillUIFile); err == nil {
		var sj skillUIJson
		if json.Unmarshal(data, &sj) == nil {
			meta.IsMarket = true
			meta.MarketID = sj.MarketID
			meta.TitleEn = sj.TitleEn
			meta.TitleZh = sj.TitleZh
			meta.DescEn = sj.DescEn
			meta.DescZh = sj.DescZh
			if sj.Owner != "" {
				meta.Owner = sj.Owner
			}
			if sj.Version != "" {
				meta.Version = sj.Version
			}
		}
	}

	// UpdatedAt from directory mod time
	if info, err := os.Stat(dir); err == nil {
		meta.UpdatedAt = info.ModTime().Format("2006/01/02 15:04:05")
	}

	return meta
}

// ListLocalSkills scans the skill directory and returns all installed skills
func (a *App) ListLocalSkills() ([]SkillMeta, error) {
	skillDir := a.getSkillDir()
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(skillDir)
	if err != nil {
		return nil, err
	}
	skills := make([]SkillMeta, 0)
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		dir := filepath.Join(skillDir, entry.Name())
		skill := parseSkillMeta(dir)
		// Detect synced tools
		skill.SyncedTools = detectSyncedTools(entry.Name(), a.getSkillDir())
		skills = append(skills, skill)
	}
	return skills, nil
}

// InstallSkillFromUrl downloads a zip from the given URL and installs it
func (a *App) InstallSkillFromUrl(url, name string) error {
	skillDir := a.getSkillDir()
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return err
	}
	tmpFile, err := os.CreateTemp("", "skill-*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()
	if _, err = io.Copy(tmpFile, resp.Body); err != nil {
		return fmt.Errorf("写入临时文件失败: %w", err)
	}
	tmpFile.Close()

	destDir := filepath.Join(skillDir, name)
	return extractZip(tmpFile.Name(), destDir)
}

// findSkillDirs recursively walks rootDir and returns all directories containing
// a SKILL.md file. It does NOT recurse into already-identified skill directories.
func findSkillDirs(rootDir string) ([]string, error) {
	var skillDirs []string
	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable paths
		}
		if !d.IsDir() {
			return nil
		}
		if path != rootDir && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}
		if _, e := os.Stat(filepath.Join(path, "SKILL.md")); e == nil {
			skillDirs = append(skillDirs, path)
			return filepath.SkipDir // don't recurse into a skill directory
		}
		return nil
	})
	return skillDirs, err
}

// InstallSkillFromGit clones a git repository and installs all detected skills (directories containing SKILL.md)
func (a *App) InstallSkillFromGit(repoUrl string) ([]string, error) {
	// Verify git is available
	if _, err := exec.LookPath("git"); err != nil {
		return nil, fmt.Errorf("未找到 git 命令，请先安装 Git")
	}

	// Create temp dir for clone
	tmpDir, err := os.MkdirTemp("", "skill-git-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	// git clone --depth=1 <repoUrl> <tmpDir>
	cmd := exec.Command("git", "clone", "--depth=1", repoUrl, tmpDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git clone 失败: %w\n%s", err, strings.TrimSpace(string(output)))
	}

	skillDir := a.getSkillDir()
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return nil, err
	}

	var installed []string

	// Recursively find all directories containing SKILL.md
	skillPaths, err := findSkillDirs(tmpDir)
	if err != nil {
		return nil, fmt.Errorf("扫描仓库失败: %w", err)
	}

	for _, subDir := range skillPaths {
		name := filepath.Base(subDir)
		// If the repo root itself is detected, use the repo name
		if subDir == tmpDir {
			name = filepath.Base(strings.TrimSuffix(repoUrl, ".git"))
		}
		destDir := filepath.Join(skillDir, name)
		if err := copyDir(subDir, destDir); err != nil {
			continue
		}
		installed = append(installed, name)
	}

	if len(installed) == 0 {
		return nil, fmt.Errorf("未在仓库中找到任何有效技能（含 SKILL.md 的目录）")
	}

	// Auto-sync each installed skill
	for _, name := range installed {
		if len(a.config.AutoSyncToolIDs) > 0 {
			a.syncToInstalledTools(name)
		}
	}

	return installed, nil
}

// InstallSkillFromMarket downloads a zip and writes skillui.json metadata
func (a *App) InstallSkillFromMarket(url string, meta SkillMeta) error {
	if err := a.InstallSkillFromUrl(url, meta.Name); err != nil {
		return err
	}
	// Write skillui.json
	sj := skillUIJson{
		MarketID:    meta.MarketID,
		Name:        meta.Name,
		TitleEn:     meta.TitleEn,
		TitleZh:     meta.TitleZh,
		DescEn:      meta.DescEn,
		DescZh:      meta.DescZh,
		Owner:       meta.Owner,
		Version:     meta.Version,
		InstalledAt: time.Now().Format("2006-01-02T15:04:05Z"),
	}
	data, err := json.MarshalIndent(sj, "", "  ")
	if err != nil {
		return err
	}
	destDir := filepath.Join(a.getSkillDir(), meta.Name)
	if err := os.WriteFile(filepath.Join(destDir, "skillui.json"), data, 0644); err != nil {
		return err
	}
	// Auto-sync to configured IDE tools if any
	if len(a.config.AutoSyncToolIDs) > 0 {
		a.syncToInstalledTools(meta.Name)
	}
	return nil
}

// InstallSkillFromLocalPath installs a skill from a local directory or file path
func (a *App) InstallSkillFromLocalPath(srcPath string) error {
	skillDir := a.getSkillDir()
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return err
	}
	name := filepath.Base(srcPath)
	// strip extension for zip/tar files
	name = strings.TrimSuffix(name, ".zip")
	name = strings.TrimSuffix(name, ".tar.gz")

	info, err := os.Stat(srcPath)
	if err != nil {
		return err
	}
	destDir := filepath.Join(skillDir, name)

	if info.IsDir() {
		return copyDir(srcPath, destDir)
	}
	// treat as zip
	return extractZip(srcPath, destDir)
}

// InstallSkillFromText creates a skill from pasted markdown text
func (a *App) InstallSkillFromText(name, content string) error {
	skillDir := a.getSkillDir()
	destDir := filepath.Join(skillDir, name)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(destDir, "SKILL.md"), []byte(content), 0644)
}

// DeleteSkill removes a skill directory and cleans up all synced tool files
func (a *App) DeleteSkill(name string) error {
	// Remove synced files from all IDE tool rules directories
	defs := ideToolDefs()
	for _, def := range defs {
		rulesDir := def.getRulesDir()
		if rulesDir == "" {
			continue
		}
		// Remove skillName.md symlink/file
		skillFile := filepath.Join(rulesDir, name+".md")
		if _, err := os.Lstat(skillFile); err == nil {
			os.Remove(skillFile)
		}
		// Remove skillName dir/link (fallback path)
		skillLink := filepath.Join(rulesDir, name)
		if _, err := os.Lstat(skillLink); err == nil {
			os.RemoveAll(skillLink)
		}
	}

	skillDir := filepath.Join(a.getSkillDir(), name)
	return os.RemoveAll(skillDir)
}

// extractZip extracts a zip file to destDir, auto-handling single-root nesting
func extractZip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("无法打开ZIP: %w", err)
	}
	defer r.Close()

	// Detect common root prefix
	rootPrefix := ""
	for _, f := range r.File {
		name := filepath.ToSlash(f.Name)
		parts := strings.SplitN(name, "/", 2)
		if len(parts) < 2 {
			// file at root level, no common prefix
			rootPrefix = ""
			break
		}
		if rootPrefix == "" {
			rootPrefix = parts[0]
		} else if rootPrefix != parts[0] {
			rootPrefix = ""
			break
		}
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	for _, f := range r.File {
		name := filepath.ToSlash(f.Name)
		// Strip common root prefix
		if rootPrefix != "" {
			prefix := rootPrefix + "/"
			if strings.HasPrefix(name, prefix) {
				name = name[len(prefix):]
			} else {
				continue
			}
		}
		if name == "" {
			continue
		}

		targetPath := filepath.Join(destDir, filepath.FromSlash(name))

		// Prevent zip-slip
		if !strings.HasPrefix(targetPath, filepath.Clean(destDir)+string(os.PathSeparator)) &&
			targetPath != filepath.Clean(destDir) {
			return fmt.Errorf("非法路径: %s", targetPath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(targetPath, f.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.Create(targetPath)
		if err != nil {
			rc.Close()
			return err
		}
		_, err = io.Copy(out, rc)
		rc.Close()
		out.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// ideToolDefs returns the list of known IDE tools with per-platform detection rules
func ideToolDefs() []ideToolDef {
	home, _ := os.UserHomeDir()
	return []ideToolDef{
		{
			ID:   "cursor",
			Name: "Cursor",
			CheckPathsMac: []string{
				filepath.Join(home, ".cursor"),
				"/Applications/Cursor.app",
			},
			CheckPathsWin: []string{
				filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "cursor", "Cursor.exe"),
				filepath.Join(os.Getenv("APPDATA"), "Cursor"),
			},
			CheckPathsLin: []string{
				filepath.Join(home, ".config", "Cursor"),
				"/usr/share/applications/cursor.desktop",
			},
			RulesDirMac: filepath.Join(home, ".cursor", "rules"),
			RulesDirWin: filepath.Join(os.Getenv("APPDATA"), "Cursor", "User", "rules"),
			RulesDirLin: filepath.Join(home, ".config", "Cursor", "User", "rules"),
		},
		{
			ID:   "claude_code",
			Name: "Claude Code",
			CheckPathsMac: []string{},
			CheckPathsWin: []string{},
			CheckPathsLin: []string{},
			RulesDirMac:   filepath.Join(home, ".claude", "commands"),
			RulesDirWin:   filepath.Join(os.Getenv("APPDATA"), ".claude", "commands"),
			RulesDirLin:   filepath.Join(home, ".claude", "commands"),
		},
		{
			ID:   "windsurf",
			Name: "Windsurf",
			CheckPathsMac: []string{
				"/Applications/Windsurf.app",
				filepath.Join(home, "Library", "Application Support", "Windsurf"),
			},
			CheckPathsWin: []string{
				filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "windsurf", "Windsurf.exe"),
				filepath.Join(os.Getenv("APPDATA"), "Windsurf"),
			},
			CheckPathsLin: []string{
				filepath.Join(home, ".config", "Windsurf"),
				"/usr/share/applications/windsurf.desktop",
			},
			RulesDirMac: filepath.Join(home, ".codeium", "windsurf", "memories"),
			RulesDirWin: filepath.Join(os.Getenv("APPDATA"), "Codeium", "windsurf", "memories"),
			RulesDirLin: filepath.Join(home, ".codeium", "windsurf", "memories"),
		},
		{
			ID:   "trae",
			Name: "TRAE IDE",
			CheckPathsMac: []string{
				"/Applications/Trae.app",
				filepath.Join(home, "Library", "Application Support", "Trae"),
			},
			CheckPathsWin: []string{
				filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "trae", "Trae.exe"),
				filepath.Join(os.Getenv("APPDATA"), "Trae"),
			},
			CheckPathsLin: []string{
				filepath.Join(home, ".config", "Trae"),
			},
			RulesDirMac: filepath.Join(home, "Library", "Application Support", "Trae", "User", "rules"),
			RulesDirWin: filepath.Join(os.Getenv("APPDATA"), "Trae", "User", "rules"),
			RulesDirLin: filepath.Join(home, ".config", "Trae", "User", "rules"),
		},
		{
			ID:   "zed",
			Name: "Zed",
			CheckPathsMac: []string{
				"/Applications/Zed.app",
				filepath.Join(home, "Library", "Application Support", "Zed"),
			},
			CheckPathsWin: []string{
				filepath.Join(os.Getenv("APPDATA"), "Zed"),
			},
			CheckPathsLin: []string{
				filepath.Join(home, ".config", "zed"),
			},
			RulesDirMac: filepath.Join(home, "Library", "Application Support", "Zed", "rules"),
			RulesDirWin: filepath.Join(os.Getenv("APPDATA"), "Zed", "rules"),
			RulesDirLin: filepath.Join(home, ".config", "zed", "rules"),
		},
		{
			ID:   "kilo_code",
			Name: "Kilo Code",
			CheckPathsMac: []string{
				filepath.Join(home, ".kilo"),
			},
			CheckPathsWin: []string{
				filepath.Join(os.Getenv("APPDATA"), "Kilo"),
			},
			CheckPathsLin: []string{
				filepath.Join(home, ".kilo"),
			},
			RulesDirMac: filepath.Join(home, ".kilo", "rules"),
			RulesDirWin: filepath.Join(os.Getenv("APPDATA"), "Kilo", "rules"),
			RulesDirLin: filepath.Join(home, ".kilo", "rules"),
		},
		{
			ID:   "roo_code",
			Name: "Roo Code",
			CheckPathsMac: []string{
				filepath.Join(home, ".roo"),
			},
			CheckPathsWin: []string{
				filepath.Join(os.Getenv("APPDATA"), "Roo"),
			},
			CheckPathsLin: []string{
				filepath.Join(home, ".roo"),
			},
			RulesDirMac: filepath.Join(home, ".roo", "rules"),
			RulesDirWin: filepath.Join(os.Getenv("APPDATA"), "Roo", "rules"),
			RulesDirLin: filepath.Join(home, ".roo", "rules"),
		},
		{
			ID:   "goose",
			Name: "Goose",
			CheckPathsMac: []string{
				filepath.Join(home, ".config", "goose"),
			},
			CheckPathsWin: []string{
				filepath.Join(os.Getenv("APPDATA"), "goose"),
			},
			CheckPathsLin: []string{
				filepath.Join(home, ".config", "goose"),
			},
			RulesDirMac: filepath.Join(home, ".config", "goose", "rules"),
			RulesDirWin: filepath.Join(os.Getenv("APPDATA"), "goose", "rules"),
			RulesDirLin: filepath.Join(home, ".config", "goose", "rules"),
		},
		{
			ID:   "gemini_cli",
			Name: "Gemini CLI",
			CheckPathsMac: []string{},
			CheckPathsWin: []string{},
			CheckPathsLin: []string{},
			RulesDirMac:   filepath.Join(home, ".gemini", "rules"),
			RulesDirWin:   filepath.Join(os.Getenv("APPDATA"), "gemini", "rules"),
			RulesDirLin:   filepath.Join(home, ".gemini", "rules"),
		},
		{
			ID:   "github_copilot",
			Name: "GitHub Copilot",
			CheckPathsMac: []string{
				filepath.Join(home, "Library", "Application Support", "GitHub Copilot"),
				filepath.Join(home, ".config", "github-copilot"),
			},
			CheckPathsWin: []string{
				filepath.Join(os.Getenv("APPDATA"), "GitHub Copilot"),
			},
			CheckPathsLin: []string{
				filepath.Join(home, ".config", "github-copilot"),
			},
			RulesDirMac: filepath.Join(home, ".github", "copilot", "rules"),
			RulesDirWin: filepath.Join(os.Getenv("USERPROFILE"), ".github", "copilot", "rules"),
			RulesDirLin: filepath.Join(home, ".github", "copilot", "rules"),
		},
		{
			ID:   "opencode",
			Name: "OpenCode",
			CheckPathsMac: []string{
				filepath.Join(home, ".config", "opencode"),
			},
			CheckPathsWin: []string{
				filepath.Join(os.Getenv("APPDATA"), "opencode"),
			},
			CheckPathsLin: []string{
				filepath.Join(home, ".config", "opencode"),
			},
			RulesDirMac: filepath.Join(home, ".config", "opencode", "rules"),
			RulesDirWin: filepath.Join(os.Getenv("APPDATA"), "opencode", "rules"),
			RulesDirLin: filepath.Join(home, ".config", "opencode", "rules"),
		},
		{
			ID:   "amp",
			Name: "Amp",
			CheckPathsMac: []string{
				filepath.Join(home, ".amp"),
			},
			CheckPathsWin: []string{
				filepath.Join(os.Getenv("APPDATA"), "Amp"),
			},
			CheckPathsLin: []string{
				filepath.Join(home, ".amp"),
			},
			RulesDirMac: filepath.Join(home, ".amp", "rules"),
			RulesDirWin: filepath.Join(os.Getenv("APPDATA"), "Amp", "rules"),
			RulesDirLin: filepath.Join(home, ".amp", "rules"),
		},
		{
			ID:            "codex",
			Name:          "Codex",
			CheckPathsMac: []string{},
			CheckPathsWin: []string{},
			CheckPathsLin: []string{},
			RulesDirMac:   filepath.Join(home, ".codex", "rules"),
			RulesDirWin:   filepath.Join(os.Getenv("APPDATA"), "codex", "rules"),
			RulesDirLin:   filepath.Join(home, ".codex", "rules"),
		},
		{
			ID:   "amazon_q",
			Name: "Amazon Q",
			CheckPathsMac: []string{
				filepath.Join(home, ".aws", "amazonq"),
			},
			CheckPathsWin: []string{
				filepath.Join(os.Getenv("USERPROFILE"), ".aws", "amazonq"),
			},
			CheckPathsLin: []string{
				filepath.Join(home, ".aws", "amazonq"),
			},
			RulesDirMac: filepath.Join(home, ".aws", "amazonq", "rules"),
			RulesDirWin: filepath.Join(os.Getenv("USERPROFILE"), ".aws", "amazonq", "rules"),
			RulesDirLin: filepath.Join(home, ".aws", "amazonq", "rules"),
		},
		{
			ID:   "cline",
			Name: "Cline",
			CheckPathsMac: []string{
				filepath.Join(home, ".cline"),
			},
			CheckPathsWin: []string{
				filepath.Join(os.Getenv("APPDATA"), "cline"),
			},
			CheckPathsLin: []string{
				filepath.Join(home, ".cline"),
			},
			RulesDirMac: filepath.Join(home, ".cline", "rules"),
			RulesDirWin: filepath.Join(os.Getenv("APPDATA"), "cline", "rules"),
			RulesDirLin: filepath.Join(home, ".cline", "rules"),
		},
		{
			ID:   "antigravity",
			Name: "Antigravity",
			CheckPathsMac: []string{
				filepath.Join(home, ".antigravity"),
			},
			CheckPathsWin: []string{
				filepath.Join(os.Getenv("APPDATA"), "antigravity"),
			},
			CheckPathsLin: []string{
				filepath.Join(home, ".antigravity"),
			},
			RulesDirMac: filepath.Join(home, ".antigravity", "rules"),
			RulesDirWin: filepath.Join(os.Getenv("APPDATA"), "antigravity", "rules"),
			RulesDirLin: filepath.Join(home, ".antigravity", "rules"),
		},
		{
			ID:   "qoder",
			Name: "Qoder",
			CheckPathsMac: []string{
				filepath.Join(home, ".qoder"),
			},
			CheckPathsWin: []string{
				filepath.Join(os.Getenv("APPDATA"), "qoder"),
			},
			CheckPathsLin: []string{
				filepath.Join(home, ".qoder"),
			},
			RulesDirMac: filepath.Join(home, ".qoder", "rules"),
			RulesDirWin: filepath.Join(os.Getenv("APPDATA"), "qoder", "rules"),
			RulesDirLin: filepath.Join(home, ".qoder", "rules"),
		},
		{
			ID:            "auggie_cli",
			Name:          "Auggie CLI",
			CheckPathsMac: []string{},
			CheckPathsWin: []string{},
			CheckPathsLin: []string{},
			RulesDirMac:   filepath.Join(home, ".augment", "rules"),
			RulesDirWin:   filepath.Join(os.Getenv("APPDATA"), "augment", "rules"),
			RulesDirLin:   filepath.Join(home, ".augment", "rules"),
		},
		{
			ID:            "qwen_code",
			Name:          "Qwen Code",
			CheckPathsMac: []string{},
			CheckPathsWin: []string{},
			CheckPathsLin: []string{},
			RulesDirMac:   filepath.Join(home, ".qwen-code", "rules"),
			RulesDirWin:   filepath.Join(os.Getenv("APPDATA"), "qwen-code", "rules"),
			RulesDirLin:   filepath.Join(home, ".qwen-code", "rules"),
		},
		{
			ID:   "codebuddy",
			Name: "CodeBuddy",
			CheckPathsMac: []string{
				filepath.Join(home, "Library", "Application Support", "CodeBuddy"),
			},
			CheckPathsWin: []string{
				filepath.Join(os.Getenv("APPDATA"), "CodeBuddy"),
			},
			CheckPathsLin: []string{
				filepath.Join(home, ".config", "CodeBuddy"),
			},
			RulesDirMac: filepath.Join(home, ".codebuddy", "rules"),
			RulesDirWin: filepath.Join(os.Getenv("APPDATA"), "CodeBuddy", "rules"),
			RulesDirLin: filepath.Join(home, ".codebuddy", "rules"),
		},
		{
			ID:   "costrict",
			Name: "CoStrict",
			CheckPathsMac: []string{
				filepath.Join(home, ".costrict"),
			},
			CheckPathsWin: []string{
				filepath.Join(os.Getenv("APPDATA"), "costrict"),
			},
			CheckPathsLin: []string{
				filepath.Join(home, ".costrict"),
			},
			RulesDirMac: filepath.Join(home, ".costrict", "rules"),
			RulesDirWin: filepath.Join(os.Getenv("APPDATA"), "costrict", "rules"),
			RulesDirLin: filepath.Join(home, ".costrict", "rules"),
		},
		{
			ID:            "crush",
			Name:          "Crush",
			CheckPathsMac: []string{},
			CheckPathsWin: []string{},
			CheckPathsLin: []string{},
			RulesDirMac:   filepath.Join(home, ".crush", "rules"),
			RulesDirWin:   filepath.Join(os.Getenv("APPDATA"), "crush", "rules"),
			RulesDirLin:   filepath.Join(home, ".crush", "rules"),
		},
		{
			ID:   "factory_droid",
			Name: "Factory Droid",
			CheckPathsMac: []string{
				filepath.Join(home, ".factory"),
			},
			CheckPathsWin: []string{
				filepath.Join(os.Getenv("APPDATA"), "factory"),
			},
			CheckPathsLin: []string{
				filepath.Join(home, ".factory"),
			},
			RulesDirMac: filepath.Join(home, ".factory", "rules"),
			RulesDirWin: filepath.Join(os.Getenv("APPDATA"), "factory", "rules"),
			RulesDirLin: filepath.Join(home, ".factory", "rules"),
		},
		{
			ID:   "iflow",
			Name: "iFlow",
			CheckPathsMac: []string{
				filepath.Join(home, ".iflow"),
			},
			CheckPathsWin: []string{
				filepath.Join(os.Getenv("APPDATA"), "iflow"),
			},
			CheckPathsLin: []string{
				filepath.Join(home, ".iflow"),
			},
			RulesDirMac: filepath.Join(home, ".iflow", "rules"),
			RulesDirWin: filepath.Join(os.Getenv("APPDATA"), "iflow", "rules"),
			RulesDirLin: filepath.Join(home, ".iflow", "rules"),
		},
	}
}

// getCheckPaths returns platform-specific check paths for a tool def
func (def *ideToolDef) getCheckPaths() []string {
	switch goruntime.GOOS {
	case "windows":
		return def.CheckPathsWin
	case "linux":
		return def.CheckPathsLin
	default: // darwin
		return def.CheckPathsMac
	}
}

// getRulesDir returns platform-specific rules directory for a tool def
func (def *ideToolDef) getRulesDir() string {
	switch goruntime.GOOS {
	case "windows":
		return def.RulesDirWin
	case "linux":
		return def.RulesDirLin
	default: // darwin
		return def.RulesDirMac
	}
}

// isCommandInPath checks if a command is available in PATH
func isCommandInPath(cmd string) (bool, string) {
	path, err := exec.LookPath(cmd)
	if err == nil {
		return true, path
	}
	return false, ""
}

// ScanIDETools detects installed AI coding tools on the current system
func (a *App) ScanIDETools() ([]IDEToolInfo, error) {
	defs := ideToolDefs()
	result := make([]IDEToolInfo, 0, len(defs))

	for _, def := range defs {
		tool := IDEToolInfo{
			ID:            def.ID,
			Name:          def.Name,
			Installed:     false,
			Path:          "",
			SkillRulesDir: def.getRulesDir(),
		}

		// Check command in PATH first (works cross-platform for CLI tools)
		switch def.ID {
		case "claude_code":
			if found, p := isCommandInPath("claude"); found {
				tool.Installed = true
				tool.Path = p
			}
		case "gemini_cli":
			if found, p := isCommandInPath("gemini"); found {
				tool.Installed = true
				tool.Path = p
			}
		case "goose":
			if found, p := isCommandInPath("goose"); found {
				tool.Installed = true
				tool.Path = p
			}
		case "codex":
			if found, p := isCommandInPath("codex"); found {
				tool.Installed = true
				tool.Path = p
			}
		case "amazon_q":
			if found, p := isCommandInPath("q"); found {
				tool.Installed = true
				tool.Path = p
			}
		case "auggie_cli":
			if found, p := isCommandInPath("auggie"); found {
				tool.Installed = true
				tool.Path = p
			}
		case "qwen_code":
			if found, p := isCommandInPath("qwen-code"); found {
				tool.Installed = true
				tool.Path = p
			}
		case "crush":
			if found, p := isCommandInPath("crush"); found {
				tool.Installed = true
				tool.Path = p
			}
		case "factory_droid":
			if found, p := isCommandInPath("droid"); found {
				tool.Installed = true
				tool.Path = p
			}
		}

		// Check filesystem paths
		if !tool.Installed {
			for _, checkPath := range def.getCheckPaths() {
				if checkPath == "" {
					continue
				}
				if _, err := os.Stat(checkPath); err == nil {
					tool.Installed = true
					tool.Path = checkPath
					break
				}
			}
		}

		result = append(result, tool)
	}
	return result, nil
}

// detectSyncedTools checks which tools have this skill synced to their rules dir
func detectSyncedTools(skillName, _ string) []string {
	defs := ideToolDefs()
	synced := make([]string, 0)
	for _, def := range defs {
		rulesDir := def.getRulesDir()
		if rulesDir == "" {
			continue
		}
		// Check if skill file or dir exists inside tools rules dir
		skillFile := filepath.Join(rulesDir, skillName+".md")
		skillLink := filepath.Join(rulesDir, skillName)
		if _, err := os.Lstat(skillFile); err == nil {
			synced = append(synced, def.Name)
		} else if _, err := os.Lstat(skillLink); err == nil {
			synced = append(synced, def.Name)
		}
	}
	return synced
}

// SyncSkillToTools copies (or symlinks on unix) the skill's SKILL.md to each tool's rules dir
func (a *App) SyncSkillToTools(skillName string, toolIds []string) error {
	skillDir := a.getSkillDir()
	skillMdPath := filepath.Join(skillDir, skillName, "SKILL.md")
	if _, err := os.Stat(skillMdPath); err != nil {
		return fmt.Errorf("技能文件不存在: %s", skillMdPath)
	}

	defs := ideToolDefs()
	defMap := make(map[string]ideToolDef, len(defs))
	for _, d := range defs {
		defMap[d.ID] = d
	}

	var errs []string
	for _, toolID := range toolIds {
		def, ok := defMap[toolID]
		if !ok {
			continue
		}
		rulesDir := def.getRulesDir()
		if rulesDir == "" {
			continue
		}
		if err := os.MkdirAll(rulesDir, 0755); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", def.Name, err))
			continue
		}
		destFile := filepath.Join(rulesDir, skillName+".md")

		// Remove existing
		os.Remove(destFile)
		if _, err := os.Lstat(destFile); err == nil {
			os.RemoveAll(destFile)
		}

		if goruntime.GOOS == "windows" {
			// Windows: copy file
			if err := copyFile(skillMdPath, destFile); err != nil {
				errs = append(errs, fmt.Sprintf("%s: %v", def.Name, err))
			}
		} else {
			// macOS / Linux: symlink
			if err := os.Symlink(skillMdPath, destFile); err != nil {
				// fallback to copy
				if err2 := copyFile(skillMdPath, destFile); err2 != nil {
					errs = append(errs, fmt.Sprintf("%s: %v", def.Name, err2))
				}
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("部分工具同步失败: %s", strings.Join(errs, "; "))
	}
	return nil
}

// UnsyncSkillFromTools removes the skill's synced file from each tool's rules dir
func (a *App) UnsyncSkillFromTools(skillName string, toolIds []string) error {
	defs := ideToolDefs()
	defMap := make(map[string]ideToolDef, len(defs))
	for _, d := range defs {
		defMap[d.ID] = d
	}

	var errs []string
	for _, toolID := range toolIds {
		def, ok := defMap[toolID]
		if !ok {
			continue
		}
		rulesDir := def.getRulesDir()
		if rulesDir == "" {
			continue
		}
		destFile := filepath.Join(rulesDir, skillName+".md")
		if err := os.Remove(destFile); err != nil && !os.IsNotExist(err) {
			errs = append(errs, fmt.Sprintf("%s: %v", def.Name, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("部分工具取消同步失败: %s", strings.Join(errs, "; "))
	}
	return nil
}
