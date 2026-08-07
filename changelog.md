## [Unreleased]

### Features
- 新增 `scripts/package-linux.sh` 用于 Linux 打包（deb + AppImage，不再产出 tar.gz），产物文件名带版本号（如 `SkillUI-0.2.0-linux-amd64.deb` / `.AppImage`）。
- CI 打包产物文件名统一带版本号：Linux `.deb` / `.AppImage`、Windows 独立 `.exe`、macOS `.dmg` 均包含 `<version>` 字段；Windows 仅产出独立 `.exe`（不再产出 NSIS 安装器与 `.zip`）。
- 在 Makefile 中加入 App Store 签名配置（BUNDLE_ID / 描述文件 / 签名身份）与 `build-devtools` 目标。
- 为 App Store 构建添加标记 (VITE_APPSTORE_BUILD)，禁用版本检查并隐藏相关界面。

### Improvements
- 统一所有操作系统的数据目录为 `~/.skillui`：macOS 不再使用 `~/Library/Application Support/SkillUI`，首次启动时自动将旧目录一次性迁移回 `~/.skillui`。
- 统一默认技能目录为 `~/.skillui/skills`：macOS 不再使用 `~/Documents/SkillUI`，用户未自定义技能目录时自动将旧目录迁移到新位置。
- CI 流水线重构为 Main / Test / Release 三个流程

### Fixes
- 修复 macOS 在 App Sandbox 下因使用系统临时目录导致技能安装失败的问题，改为将临时 ZIP 写入用户技能目录。

### Improvements
- 将 macOS 应用数据目录从 `~/.skillui` 迁移至 `~/Library/Application Support/SkillUI`，符合 App Sandbox 规范；首次启动时自动执行一次性数据迁移。
- 将 macOS 默认技能目录改为 `~/Documents/SkillUI`，确保用户可访问且符合 App Sandbox 指引 2.4.5(i)。
- 删除未使用的前端 README 和资产文件。
- 更新前端 index.html 和 style.css。
- 将 Dock 图标处理重构到新的 `internal/platform` 包，并更新了引用/调用。
- 更新 macOS Info.plist 中的 bundle identifier 为 `com.skillui`。
- 更新应用和托盘图标；添加新的 logo 资源并移除过时的 logo 文件。
- 从项目根目录删除遗留的 dock_* 文件 (清理)。
- 重命名页面专属 Vue 组件为带有父页面前缀，并更新导入路径。
- 在 macOS Info.plist 中添加 LSApplicationCategoryType 和 ITSAppUsesNonExemptEncryption 键，用于正确的应用分类和加密声明。
