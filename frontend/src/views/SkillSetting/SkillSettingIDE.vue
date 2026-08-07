<template>
  <!-- Scan Section -->
  <div
      class="bg-white dark:bg-slate-800/80 rounded-xl border border-slate-200 dark:border-slate-700/50 p-6 shadow-sm">
    <div class="flex justify-between items-center mb-6">
      <div>
        <h3 class="text-lg font-bold text-slate-900 dark:text-white flex items-center gap-2">
          <Monitor :size="20" class="text-emerald-500"/>
          {{ $t('toolSettings.ideTitle') }}
        </h3>
        <p class="text-slate-500 text-sm mt-1">{{ $t('toolSettings.ideDesc') }}</p>
      </div>
      <Button type="primary" :loading="scanning" @click="scanTools">
        <template #icon>
          <Scan :size="16"/>
        </template>
        {{ scanning ? $t('toolSettings.scanning') : $t('toolSettings.rescan') }}
      </Button>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      <div
          v-for="tool in tools"
          :key="tool.id"
          class="ide-card flex flex-col p-4 rounded-lg border transition-all h-full"
          :class="tool.installed ? 'border-emerald-500/30 bg-emerald-500/5 dark:bg-emerald-900/10' : 'border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/50 opacity-70'"
      >
        <div class="flex items-start justify-between mb-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 flex items-center justify-center rounded-lg shrink-0"
                 :class="tool.installed ? 'bg-white dark:bg-slate-800 text-emerald-500 shadow-sm' : 'bg-slate-200 dark:bg-slate-700 text-slate-400'">
              <component :is="iconForTool(tool.id)" :size="20"/>
            </div>
            <div class="min-w-0">
              <div class="flex flex-col">
                <h4 class="font-bold text-slate-900 dark:text-white truncate" :title="tool.name">{{ tool.name }}</h4>
                <span class="text-xs text-slate-500 font-mono truncate">{{ tool.id }}</span>
              </div>
            </div>
          </div>

          <div class="shrink-0 flex items-center gap-1.5">
            <span v-if="tool.manual" class="text-[10px] bg-sky-500 text-white px-1.5 py-0.5 rounded font-bold shadow-sm flex items-center gap-1">
              <Wrench :size="10"/> {{ $t('toolSettings.manual') }}
            </span>
            <span v-if="tool.installed" class="text-[10px] bg-emerald-500 text-white px-1.5 py-0.5 rounded font-bold shadow-sm flex items-center gap-1">
              <Check :size="10"/> {{ $t('toolSettings.installed') }}
            </span>
            <span v-else class="text-[10px] text-slate-400 dark:text-slate-500">{{ $t('toolSettings.notInstalled') }}</span>
          </div>
        </div>

        <!-- Path Config -->
        <div class="mt-auto pt-3 border-t border-slate-100 dark:border-slate-700/50 space-y-2">
          <div v-if="tool.installed">
            <div v-if="tool.path" class="space-y-1">
              <div class="text-[10px] text-slate-400 font-mono bg-slate-100 dark:bg-slate-900/50 px-2 py-1 rounded truncate" :title="tool.path">
                {{ tool.path }}
              </div>
            </div>
            <div v-if="tool.skillRulesDir" class="space-y-1">
              <div class="text-[10px] text-slate-400 mb-0.5">{{ $t('toolSettings.skillPath') }}</div>
              <div class="text-[10px] text-emerald-600 dark:text-emerald-400 font-mono bg-emerald-50 dark:bg-emerald-900/20 px-2 py-1 rounded truncate" :title="tool.skillRulesDir">
                {{ tool.skillRulesDir }}
              </div>
            </div>
            <!-- Auto-sync toggle -->
            <div class="flex items-center justify-between pt-1">
              <span class="text-xs text-slate-500">{{ $t('toolSettings.autoSync') }}</span>
              <Switch
                  :checked="isAutoSync(tool.id)"
                  :loading="togglingTools[tool.id]"
                  size="small"
                  @change="(v: any) => setAutoSync(tool.id, !!v)"
              />
            </div>
          </div>
          <div v-else class="h-6 flex items-center text-xs text-slate-400">
            {{ $t('toolSettings.noPath') }}
          </div>

          <!-- Manual path actions -->
          <div class="flex items-center gap-2 pt-1">
            <Button size="small" class="!text-xs" @click="openToolPathModal(tool)">
              {{ tool.manual ? $t('toolSettings.modifyPath') : $t('toolSettings.setPath') }}
            </Button>
            <Button v-if="tool.manual" size="small" danger class="!text-xs" @click="handleClearToolPath(tool.id)">
              {{ $t('toolSettings.clearPath') }}
            </Button>
          </div>
        </div>
      </div>
    </div>

    <!-- Manual path modal -->
    <SkillSettingToolPathModal
        v-model:visible="toolPathModalVisible"
        :tool-id="editingToolId"
        :tool-name="editingToolName"
        :default-value="editingToolDefaultPath"
        @saved="onToolPathSaved"
    />
  </div>
</template>

<script setup lang="ts">
import { Button, Switch, message } from 'ant-design-vue';
import { Box, Check, Code, Command, Monitor, Scan, Terminal, Wrench, Zap } from 'lucide-vue-next';
import { onMounted, onUnmounted, ref, watch } from 'vue';
import { ClearToolPath, GetAutoSyncToolIDs, ScanIDETools, SetAutoSyncToolIDs, SetToolPath } from '../../../wailsjs/go/main/App';
import { useAppStore } from '../../stores/app';
import SkillSettingToolPathModal from './SkillSettingToolPathModal.vue';
import { testActionSet, testActionUnset } from '../../utils/test';

const appStore = useAppStore();

const props = defineProps<{
  defaultSkillDir: string;
}>();

const scanning = ref(false);
const autoSyncIDs = ref<string[]>([]);
const togglingTools = ref<Record<string, boolean>>({});
const toolPathModalVisible = ref(false);
const editingToolId = ref('');
const editingToolName = ref('');
const editingToolDefaultPath = ref('');

onMounted(async () => {
  try {
    autoSyncIDs.value = await GetAutoSyncToolIDs();
  } catch {
    // ignore
  }

  // ── 自动化测试 action（open 版 testActionSet 为空函数，直接保留）──────────
  testActionSet('ToolSettings.getToolCount', () => tools.value.length);
  testActionSet('ToolSettings.getInstalledCount', () => tools.value.filter((t) => t.installed).length);
  testActionSet('ToolSettings.getTools', () => tools.value);
  testActionSet('ToolSettings.scanTools', async () => {
    await scanTools();
    return tools.value.length;
  });
  testActionSet('ToolSettings.getAutoSyncIDs', () => autoSyncIDs.value);
  testActionSet('ToolSettings.setAutoSync', async (params: unknown) => {
    const { toolId, enabled } = params as { toolId: string; enabled: boolean };
    await setAutoSync(toolId, enabled);
    return autoSyncIDs.value.includes(toolId);
  });
  testActionSet('ToolSettings.openToolPathModal', (params: unknown) => {
    const { toolId } = params as { toolId: string };
    const tool = tools.value.find((t) => t.id === toolId);
    if (!tool) throw new Error(`工具不存在: ${toolId}`);
    openToolPathModal(tool);
    return true;
  });
  testActionSet('ToolSettings.isToolPathModalOpen', () => toolPathModalVisible.value);
  testActionSet('ToolSettings.closeToolPathModal', () => {
    toolPathModalVisible.value = false;
    return true;
  });
  testActionSet('ToolSettings.setToolPath', async (params: unknown) => {
    const { toolId, path } = params as { toolId: string; path: string };
    const tool = tools.value.find((t) => t.id === toolId);
    if (!tool) throw new Error(`工具不存在: ${toolId}`);
    await saveToolPath(toolId, path);
    await scanTools();
    return true;
  });
  testActionSet('ToolSettings.clearToolPath', async (params: unknown) => {
    const { toolId } = params as { toolId: string };
    await handleClearToolPath(toolId);
    return true;
  });
});

onUnmounted(() => {
  testActionUnset([
    'ToolSettings.getToolCount',
    'ToolSettings.getInstalledCount',
    'ToolSettings.getTools',
    'ToolSettings.scanTools',
    'ToolSettings.getAutoSyncIDs',
    'ToolSettings.setAutoSync',
    'ToolSettings.openToolPathModal',
    'ToolSettings.isToolPathModalOpen',
    'ToolSettings.closeToolPathModal',
    'ToolSettings.setToolPath',
    'ToolSettings.clearToolPath',
  ]);
});

const isAutoSync = (toolId: string) => autoSyncIDs.value.includes(toolId);

const setAutoSync = async (toolId: string, enabled: boolean) => {
  togglingTools.value = { ...togglingTools.value, [toolId]: true };
  try {
    const next = enabled
        ? [...new Set([...autoSyncIDs.value, toolId])]
        : autoSyncIDs.value.filter(id => id !== toolId);
    await SetAutoSyncToolIDs(next);
    autoSyncIDs.value = next;
  } catch (e: any) {
    message.error(appStore.t('toolSettings.setFailed', { error: e?.message || String(e) }));
  } finally {
    const t = { ...togglingTools.value };
    delete t[toolId];
    togglingTools.value = t;
  }
};

interface Tool {
  id: string;
  name: string;
  installed: boolean;
  path: string;
  skillRulesDir: string;
  manual: boolean;
}

const iconMap: Record<string, any> = {
  cursor: Zap,
  claude_code: Terminal,
  windsurf: Terminal,
  trae: Box,
  zed: Code,
  kilo_code: Command,
  roo_code: Box,
  goose: Command,
  gemini_cli: Terminal,
  github_copilot: Code,
  opencode: Box,
  amp: Zap,
  codex: Terminal,
  amazon_q: Box,
  cline: Code,
  antigravity: Zap,
  qoder: Box,
  auggie_cli: Command,
  qwen_code: Terminal,
  codebuddy: Code,
  costrict: Box,
  crush: Zap,
  factory_droid: Box,
  iflow: Code,
  continue: Box,
  aider: Terminal,
  tabby: Box,
  coco: Command,
  mars_code: Box,
};

const iconForTool = (id: string) => iconMap[id] || Box;

const tools = ref<Tool[]>([]);

const scanTools = async () => {
  scanning.value = true;
  try {
    const result = await ScanIDETools();
    tools.value = result.map((t: any) => ({
      id: t.id,
      name: t.name,
      installed: t.installed,
      path: t.path,
      skillRulesDir: t.skillRulesDir,
      manual: t.manual || false,
    }));
    const installedCount = tools.value.filter(t => t.installed).length;
    message.success(appStore.t('toolSettings.scanSuccess', { count: installedCount }));
  } catch (e: any) {
    message.error(appStore.t('toolSettings.scanFailed', { error: e?.message || String(e) }));
  } finally {
    scanning.value = false;
  }
};

const openToolPathModal = (tool: Tool) => {
  editingToolId.value = tool.id;
  editingToolName.value = tool.name;
  editingToolDefaultPath.value = tool.manual ? tool.skillRulesDir : '';
  toolPathModalVisible.value = true;
};

const saveToolPath = async (toolId: string, path: string) => {
  await SetToolPath(toolId, path);
};

const onToolPathSaved = async (toolId: string) => {
  await scanTools();
};

const handleClearToolPath = async (toolId: string) => {
  try {
    await ClearToolPath(toolId);
    await scanTools();
  } catch (e: any) {
    message.error(appStore.t('toolSettings.toolPathFailed', { error: e?.message || String(e) }));
  }
};

watch(() => props.defaultSkillDir, (val) => {
  if (val) scanTools();
}, { immediate: true });
</script>
