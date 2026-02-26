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
          class="flex flex-col p-4 rounded-lg border transition-all h-full"
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

          <div v-if="tool.installed" class="shrink-0">
            <span class="text-[10px] bg-emerald-500 text-white px-1.5 py-0.5 rounded font-bold shadow-sm flex items-center gap-1">
              <Check :size="10"/> {{ $t('toolSettings.installed') }}
            </span>
          </div>
          <div v-else class="shrink-0">
            <span class="text-[10px] text-slate-400 dark:text-slate-500">{{ $t('toolSettings.notInstalled') }}</span>
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
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Button, Switch, message } from 'ant-design-vue';
import { Box, Check, Code, Command, Monitor, Scan, Terminal, Zap } from 'lucide-vue-next';
import { onMounted, ref, watch } from 'vue';
import { GetAutoSyncToolIDs, ScanIDETools, SetAutoSyncToolIDs } from '../../../wailsjs/go/main/App';
import { useAppStore } from '../../stores/app';

const appStore = useAppStore();

const props = defineProps<{
  defaultSkillDir: string;
}>();

const scanning = ref(false);
const autoSyncIDs = ref<string[]>([]);
const togglingTools = ref<Record<string, boolean>>({});

onMounted(async () => {
  try {
    autoSyncIDs.value = await GetAutoSyncToolIDs();
  } catch {
    // ignore
  }
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
    }));
    const installedCount = tools.value.filter(t => t.installed).length;
    message.success(appStore.t('toolSettings.scanSuccess', { count: installedCount }));
  } catch (e: any) {
    message.error(appStore.t('toolSettings.scanFailed', { error: e?.message || String(e) }));
  } finally {
    scanning.value = false;
  }
};

watch(() => props.defaultSkillDir, (val) => {
  if (val) scanTools();
}, { immediate: true });
</script>
