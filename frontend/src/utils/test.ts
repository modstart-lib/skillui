/**
 * src/utils/test.ts — UI 自动化测试辅助工具
 *
 * 核心思路：action 注册时机即页面就绪时机。
 * 测试侧通过 callAction 等待 action 被注册，无需额外的 ready 状态。
 *
 * 用法（页面组件 setup 中）：
 *   import { testActionSet, testActionUnset } from '../utils/test'
 *   import { onMounted, onUnmounted } from 'vue'
 *
 *   onMounted(() => {
 *     testActionSet('Manage.getSkillCount', () => skills.value.length)
 *   })
 *   onUnmounted(() => {
 *     testActionUnset('Manage.getSkillCount')
 *   })
 *
 * 注意：ss-publish 的条件编译裁剪不支持嵌套，所有 TYPE_PRO 块必须互相独立。
 * 发布到 Open 版时 TYPE_PRO 块被移除，函数变为空操作，零开销。
 */

export type TestAction = (arg?: unknown) => Promise<unknown> | unknown

export interface TestRegistry {
  setAction(name: string, fn: TestAction): void
  unsetAction(name: string | string[]): void
  callAction(name: string, arg?: unknown): Promise<unknown>
  listActions(): string[]
}



/**
 * 主动上报一条错误日志（仅测试模式激活后有效）。
 */
export function reportTestError(msg: string): void {
  
}

/**
 * 主动上报一条警告日志（仅测试模式激活后有效）。
 */
export function reportTestWarn(msg: string): void {
  
}

/**
 * 设置一个 test action。在 onMounted 中调用。
 * action 设置时机即页面就绪时机——测试侧 callAction 会等待 action 出现。
 * 发布到 Open 版时函数体为空，零开销。
 */
export function testActionSet(name: string, fn: TestAction): void {
  
}

/**
 * 移除一个或多个 test action。在 onUnmounted 中调用。
 */
export function testActionUnset(nameOrNames: string | string[]): void {
  
}

/**
 * 挂载 window.__test，在 App.vue 的 onMounted 中调用。
 * 注册 App.getConsoleLogs / App.clearConsoleLogs / App.startTestMode action。
 * 拦截逻辑不在此处激活，需测试程序调用 App.startTestMode 后才生效。
 */
export function initTestRegistry(opts: {
  onStartTestMode: () => void
}): void {
  
}
