//go:build !autotest

package main

// parseAutoTestPort 非测试构建（未指定 autotest tag）时返回 0，不启动测试服务。
func parseAutoTestPort() int { return 0 }

// startAutoTestServer 非测试构建时为空操作。
func (a *App) startAutoTestServer(port int) {}
