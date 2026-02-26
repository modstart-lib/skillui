[English](#skillui) | [中文](README.zh-CN.md)

# SkillUI

SkillUI is a cross-platform desktop GUI for managing AI coding assistant skills. Built with Wails, Go, and Vue 3 + TypeScript, it lets you install, organize, and sync skills to your favorite AI tools — all from one place.

## Screenshots

![](https://ms-assets.modstart.com/data/image/2026/02/26/14585_4g08_7720.png)

![](https://ms-assets.modstart.com/data/image/2026/02/26/14808_djax_1635.png)
![](https://ms-assets.modstart.com/data/image/2026/02/26/14763_mllp_1696.png)

![](https://ms-assets.modstart.com/data/image/2026/02/26/14625_csju_6992.png)

## Features

### Skill Management
- **Browse & Search**: View all locally installed skills with metadata parsed from `SKILL.md` frontmatter
- **Install from Marketplace**: One-click install from the [skillui.com](https://skillui.com) marketplace
- **Install from URL / Local Path**: Import skills directly from a download URL or a local directory/zip file
- **Remove Skills**: Uninstall skills and clean up synced tool rule files
- **Skill Details**: Read the full `SKILL.md` content in a built-in viewer

### IDE / AI Tool Sync
Sync skills to the rules directory of supported AI coding tools automatically. Supported tools:

| Tool | Platform |
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

- **Auto-sync**: Automatically sync to selected tools whenever a skill is installed from the marketplace
- **Manual sync**: Sync any skill to individual tools on demand
- **Sync status**: See which tools a skill has already been synced to

### Process Management
- **Add / Remove**: Register background processes with custom commands, arguments, and environment variables
- **Start / Stop / Restart**: Full lifecycle control with graceful shutdown
- **Auto-start**: Start processes automatically when the application launches
- **Restart Policies**: `always`, `on_failure`, or `never`
- **Monitoring**: Real-time status with PID, restart count, and error tracking

### Logging
- Rolling log files with configurable line and file retention limits
- Real-time log streaming per process
- Separate stdout / stderr capture

### Cross-Platform Support
- **Windows** (amd64)
- **macOS** (Intel and Apple Silicon)
- **Linux** (amd64, arm64)

### Auto-start on Boot
- **macOS**: LaunchAgent
- **Linux**: XDG Autostart
- **Windows**: Registry

### App Settings
- Light / Dark theme
- Simplified Chinese / English interface language
- Configurable skill storage directory with optional migration
- Version check and update prompt

## Build

### Prerequisites

- Go 1.23.0 or higher
- Node.js 18+ and npm
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

### Development

```bash
# Install frontend dependencies
cd frontend && npm install && cd ..

# Install Go dependencies
go mod download

# Run in development mode
wails dev
```

### Production Build

```bash
# Build for current platform
wails build

# The built application will be in build/bin directory
```

## 💬 Join the Community

> Add friend with note "SkillUI"

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

## License

This project is licensed under the [Apache 2.0 License](LICENSE).

