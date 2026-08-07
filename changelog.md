## [unreleased]

- 新增：AI Studio 识别支持手动设置工具规则目录，自动扫描识别不到已安装的 AI IDE 时，可手动指定路径（如 `~/.cursor/rules`），指定后该工具视为已安装并作为技能同步目标。
- 新增：补充识别主流 AI 编程工具（Continue、Aider、Tabby、Coco、MarsCode），并完善各平台默认检测路径。
- 新增：前端工具卡片提供"手动设置 / 修改路径 / 清除"操作入口与弹窗，展示手动标记。

## v0.2.0 App Store 适配完成，跨平台打包全面升级

- 新增：`scripts/package-linux.sh` 用于 Linux 打包（deb + AppImage，不再产出 tar.gz），产物文件名带版本号（如 `SkillUI-0.2.0-linux-amd64.deb` / `.AppImage`）。
- 新增：CI 打包产物文件名统一带版本号，Linux `.deb` / `.AppImage`、Windows 独立 `.exe`、macOS `.dmg` 均包含 `<version>` 字段；Windows 仅产出独立 `.exe`（不再产出 NSIS 安装器与 `.zip`）。
- 新增：在 Makefile 中加入 App Store 签名配置（BUNDLE_ID / 描述文件 / 签名身份）与 `build-devtools` 目标。
- 新增：为 App Store 构建添加标记 (VITE_APPSTORE_BUILD)，禁用版本检查并隐藏相关界面。
- 优化：统一所有操作系统的数据目录为 `~/.skillui`：macOS 不再使用 `~/Library/Application Support/SkillUI`，首次启动时自动将旧目录一次性迁移回 `~/.skillui`。
- 优化：统一默认技能目录为 `~/.skillui/skills`：macOS 不再使用 `~/Documents/SkillUI`，用户未自定义技能目录时自动将旧目录迁移到新位置。
- 优化：CI 流水线重构为 Main / Test / Release 三个流程。
- 优化：将 macOS 应用数据目录从 `~/.skillui` 迁移至 `~/Library/Application Support/SkillUI`，符合 App Sandbox 规范；首次启动时自动执行一次性数据迁移。
- 优化：将 macOS 默认技能目录改为 `~/Documents/SkillUI`，确保用户可访问且符合 App Sandbox 指引 2.4.5(i)。
- 优化：删除未使用的前端 README 和资产文件。
- 优化：更新前端 index.html 和 style.css。
- 优化：将 Dock 图标处理重构到新的 `internal/platform` 包，并更新了引用/调用。
- 优化：更新 macOS Info.plist 中的 bundle identifier 为 `com.skillui`。
- 优化：更新应用和托盘图标；添加新的 logo 资源并移除过时的 logo 文件。
- 优化：从项目根目录删除遗留的 dock_* 文件 (清理)。
- 优化：重命名页面专属 Vue 组件为带有父页面前缀，并更新导入路径。
- 优化：在 macOS Info.plist 中添加 LSApplicationCategoryType 和 ITSAppUsesNonExemptEncryption 键，用于正确的应用分类和加密声明。
- 修复：修复 macOS 在 App Sandbox 下因使用系统临时目录导致技能安装失败的问题，改为将临时 ZIP 写入用户技能目录。
