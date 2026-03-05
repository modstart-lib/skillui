## [Unreleased]

### Features
- Added App Store build flag (VITE_APPSTORE_BUILD) to disable version checks and hide related UI. 为 App Store 构建添加标记 (VITE_APPSTORE_BUILD)，禁用版本检查并隐藏相关界面。

### Improvements
- Migrated macOS app data directory from `~/.skillui` to `~/Library/Application Support/SkillUI` to comply with App Sandbox guidelines; added one-time automatic migration on first launch. 将 macOS 应用数据目录从 `~/.skillui` 迁移至 `~/Library/Application Support/SkillUI`，符合 App Sandbox 规范；首次启动时自动执行一次性数据迁移。
- Changed default skill directory on macOS to `~/Documents/SkillUI` for user-accessible storage compliant with App Sandbox Guideline 2.4.5(i). 将 macOS 默认技能目录改为 `~/Documents/SkillUI`，确保用户可访问且符合 App Sandbox 指引 2.4.5(i)。
- Bumped version to v0.2.1. 版本号升级至 v0.2.1。
- Removed unused frontend README and asset files. 删除未使用的前端 README 和资产文件。
- Updated frontend index.html and style.css. 更新前端 index.html 和 style.css。
- Refactored Dock icon handling into a new `internal/platform` package and updated imports/calls. 将 Dock 图标处理重构到新的 `internal/platform` 包，并更新了引用/调用。
- Updated bundle identifier in macOS Info.plist to `com.skillui`. 更新 macOS Info.plist 中的 bundle identifier 为 `com.skillui`。
- Updated application and tray icons; added new logo assets and removed outdated logo files. 更新应用和托盘图标；添加新的 logo 资源并移除过时的 logo 文件。
- Removed legacy dock_* files from project root (cleanup). 从项目根目录删除遗留的 dock_* 文件 (清理)。
- Renamed page-specific Vue components with parent page prefixes and updated import paths. 重命名页面专属 Vue 组件为带有父页面前缀，并更新导入路径。
- Added LSApplicationCategoryType and ITSAppUsesNonExemptEncryption keys to macOS Info.plist for proper app categorization and encryption declaration. 在 macOS Info.plist 中添加 LSApplicationCategoryType 和 ITSAppUsesNonExemptEncryption 键，用于正确的应用分类和加密声明。
