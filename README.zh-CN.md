[English](README.md) | [中文](#skillui)

# SkillUI

SkillUI 是一款跨平台的 AI 编程助手技能管理桌面应用，基于 Wails、Go 和 Vue 3 + TypeScript 构建。你可以在一个界面中安装、管理和同步技能到各种 AI 编程工具。

## 截图

![](https://ms-assets.modstart.com/data/image/2026/02/26/14585_4g08_7720.png)

![](https://ms-assets.modstart.com/data/image/2026/02/26/14808_djax_1635.png)
![](https://ms-assets.modstart.com/data/image/2026/02/26/14763_mllp_1696.png)

![](https://ms-assets.modstart.com/data/image/2026/02/26/14625_csju_6992.png)

## 功能

### 技能管理
- **浏览与搜索**：查看所有本地已安装的技能，自动解析 `SKILL.md` 中的元数据
- **从应用商店安装**：一键从 [skillui.com](https://skillui.com) 应用商店安装技能
- **从 URL / 本地路径安装**：支持从下载链接或本地目录/压缩包导入技能
- **删除技能**：卸载技能并同步清理已同步工具中的规则文件
- **技能详情**：内置阅读器直接查看 `SKILL.md` 完整内容

### IDE / AI 工具同步
将技能自动同步至支持的 AI 编程工具的 rules 目录。支持的工具：

| 工具 | 平台 |
|---|---|
| Cursor | Windows / macOS / Linux |
| Claude Code | Windows / macOS / Linux |
| Windsurf | Windows / macOS / Linux |
| TRAE IDE | Windows / macOS / Linux |
| Zed | Windows / macOS / Linux |
| Kilo Code | Windows / macOS / Linux |
| Roo Code | Windows / macOS / Linux |
| Goose | Windows / macOS / Linux |
| Gemini CLI | Windows / macOS / Linux |
| GitHub Copilot | Windows / macOS / Linux |
| OpenCode | Windows / macOS / Linux |
| Amp | Windows / macOS / Linux |

- **自动同步**：从应用商店安装技能时，自动同步到已选中的工具
- **手动同步**：随时将任意技能同步到指定工具
- **同步状态**：查看技能已被同步到哪些工具

### 进程管理
- **添加/删除**：注册后台进程，支持自定义命令、参数和环境变量
- **启动/停止/重启**：完整的进程生命周期控制，支持优雅退出
- **自动启动**：应用启动时自动拉起已配置的进程
- **重启策略**：支持 `always`、`on_failure`、`never` 三种策略
- **进程监控**：实时查看运行状态、PID、重启次数和错误信息

### 日志
- 支持按行数和文件数滚动的日志文件
- 每个进程独立的实时日志流
- 分离 stdout / stderr 输出

### 跨平台支持
- **Windows** (amd64)
- **macOS**（Intel 和 Apple Silicon）
- **Linux** (amd64, arm64)

### 开机自启
- **macOS**：使用 LaunchAgent
- **Linux**：使用 XDG Autostart
- **Windows**：使用注册表

### 应用设置
- 浅色 / 深色主题
- 简体中文 / English 界面语言
- 可配置技能存储目录，支持迁移已有技能
- 版本检测与更新提示

## 构建

### 前置要求

- Go 1.23.0 或更高版本
- Node.js 18+ 及 npm
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

### 开发模式

```bash
# 安装前端依赖
cd frontend && npm install && cd ..

# 安装 Go 依赖
go mod download

# 以开发模式运行
wails dev
```

### 生产构建

```bash
# 构建当前平台版本
wails build

# 构建产物位于 build/bin 目录
```

## 💬 交流沟通

> 添加好友时备注 "SkillUI"

<table width="100%">
    <thead>
        <tr>
            <th width="50%">WeChat Group</th>
            <th>QQ Group</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>
                <img style="width:100%;"
                     src="https://modstart.com/code_dynamic/modstart_wx" />
            </td>
            <td>
                <img style="width:100%;"
                     src="https://modstart.com/code_dynamic/modstart_qq" />
            </td>
        </tr>
    </tbody>
</table>

## 许可证

本项目基于 [Apache 2.0 License](LICENSE) 开源。
