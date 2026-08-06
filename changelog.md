## [Unreleased]

### Features
- 新增 `make screenshot` 截图流水线（Vite dev server + Playwright + mock 绑定 + ss-image-mockup），覆盖 4 个页面及关键弹窗，参考 prochub-pro 的实现。
- 截图输出采用 linkandroid 的 I18nDocument 格式：`storage/I18nDocument/asset/image/`（常规版 + `.trans.png` 透明版），image.json 含 `url` / `name` / `urlTrans` / `nameTrans`，CDN 路径为 `https://cdn.skillui.com/theme/I18nDocument/asset/image`。
- 补全 I18nDocument 中文文档：快速开始、技能管理、工具同步、设置四个分类，含 `_README.md` 索引与操作指南，文档引用 CDN 截图。
- 新增 `make publish-website`，将 `storage/I18nDocument/` 同步到 SkillUI 独立演示站。
- 新增 `scripts/package-linux.sh` 用于 Linux 打包（tar.gz + deb + AppImage）。
- 在 Makefile 中加入 App Store 签名配置（BUNDLE_ID / 描述文件 / 签名身份）与 `build-devtools` 目标。
- 补充 `make build-seed-test` 说明：全量测试流水线用于验证构建与全部测试。
- 为 App Store 构建添加标记 (VITE_APPSTORE_BUILD)，禁用版本检查并隐藏相关界面。

### Improvements
- CI 流水线重构为 Main / Test / Release 三个流程，参考 linkandroid 的开源/pro 模式；Main 和 Release 流程克隆 `skillui-pro.git` 并构建 Pro 版本发布产物。
- 创建 `skillui-pro` 仓库作为 Pro 版真源，并通过 `publish/` 机制（ss-publish 白名单 + 替换规则）同步到开源版。
- 新增 `make publish` 目标（仅 Pro 版），通过 ss-publish 将 Pro 仓库同步到开源版。

### Fixes
- 修复 macOS 在 App Sandbox 下因使用系统临时目录导致技能安装失败的问题，改为将临时 ZIP 写入用户技能目录。

### Improvements
- 将 macOS 应用数据目录从 `~/.skillui` 迁移至 `~/Library/Application Support/SkillUI`，符合 App Sandbox 规范；首次启动时自动执行一次性数据迁移。
- 将 macOS 默认技能目录改为 `~/Documents/SkillUI`，确保用户可访问且符合 App Sandbox 指引 2.4.5(i)。
- 版本号升级至 v0.2.1。
- 删除未使用的前端 README 和资产文件。
- 更新前端 index.html 和 style.css。
- 将 Dock 图标处理重构到新的 `internal/platform` 包，并更新了引用/调用。
- 更新 macOS Info.plist 中的 bundle identifier 为 `com.skillui`。
- 更新应用和托盘图标；添加新的 logo 资源并移除过时的 logo 文件。
- 从项目根目录删除遗留的 dock_* 文件 (清理)。
- 重命名页面专属 Vue 组件为带有父页面前缀，并更新导入路径。
- 在 macOS Info.plist 中添加 LSApplicationCategoryType 和 ITSAppUsesNonExemptEncryption 键，用于正确的应用分类和加密声明。
