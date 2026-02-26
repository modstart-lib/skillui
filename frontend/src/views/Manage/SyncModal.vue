<template>
  <Modal
      :open="open"
      :title="$t('sync.modalTitle')"
      :ok-button-props="{ loading: syncing }"
      :okText="$t('sync.startSync')"
      :cancelText="$t('actions.cancel')"
      @ok="handleSync"
      @update:open="$emit('update:open', $event)"
  >
    <div class="py-2">
      <p class="text-sm text-slate-500 mb-4">{{ $t('sync.selectTarget') }}</p>
      <div v-if="loadingTools" class="flex justify-center items-center gap-2 py-6">
        <Loader2 :size="18" class="text-slate-400 animate-spin"/>
        <span class="text-slate-400 text-sm">{{ $t('sync.scanning') }}</span>
      </div>
      <div v-else-if="localTools.length === 0" class="flex flex-col items-center py-8 gap-3">
        <div class="w-14 h-14 bg-slate-100 dark:bg-slate-700/50 rounded-full flex items-center justify-center">
          <Bot :size="24" class="text-slate-400"/>
        </div>
        <span class="text-center text-slate-400 text-sm">{{ $t('sync.noTools') }}</span>
      </div>
      <div v-else class="grid grid-cols-2 gap-3">
        <div
            v-for="tool in localTools"
            :key="tool.id"
            class="border rounded-lg p-3 flex items-center gap-3 cursor-pointer transition-all"
            :class="tool.checked ? 'border-emerald-500 bg-emerald-50 dark:bg-emerald-900/10' : 'border-slate-200 dark:border-slate-700 hover:border-slate-300'"
            @click="tool.checked = !tool.checked"
        >
          <div
              class="w-4 h-4 rounded border flex items-center justify-center transition-colors shrink-0"
              :class="tool.checked ? 'bg-emerald-500 border-emerald-500' : 'bg-white dark:bg-slate-800 border-slate-300 dark:border-slate-600'"
          >
            <Check v-if="tool.checked" :size="12" class="text-white"/>
          </div>
          <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ tool.name }}</span>
        </div>
      </div>
    </div>
  </Modal>
</template>

<script setup lang="ts">
import { Modal, message } from 'ant-design-vue';
import { Bot, Check, Loader2 } from 'lucide-vue-next';
import { ref, watch } from 'vue';
import { ScanIDETools, SyncSkillToTools, UnsyncSkillFromTools } from '../../../wailsjs/go/main/App';
import { useAppStore } from '../../stores/app';

const appStore = useAppStore();

const props = defineProps<{
  open: boolean;
  skill: any;
}>();

const emit = defineEmits(['update:open', 'synced']);

interface LocalTool {
  id: string;
  name: string;
  checked: boolean;
}

const localTools = ref<LocalTool[]>([]);
const loadingTools = ref(false);
const syncing = ref(false);
// 记录打开时已勾选的 tool id，用于判断哪些需要取消同步
const initialCheckedIds = ref<Set<string>>(new Set());

watch(() => props.open, async (val) => {
  if (!val) return;
  loadingTools.value = true;
  try {
    const result = await ScanIDETools();
    // Only show installed tools
    const installed = result.filter((t: any) => t.installed);
    const syncedNames = props.skill?.syncedTools || [];
    localTools.value = installed.map((t: any) => ({
      id: t.id,
      name: t.name,
      checked: syncedNames.includes(t.name),
    }));
    initialCheckedIds.value = new Set(
      localTools.value.filter(t => t.checked).map(t => t.id)
    );
  } catch (e) {
    message.error(appStore.t('sync.scanFailed'));
  } finally {
    loadingTools.value = false;
  }
});

const handleSync = async () => {
  if (!props.skill) return;
  const checkedIds = localTools.value.filter(t => t.checked).map(t => t.id);
  const uncheckedIds = localTools.value
    .filter(t => !t.checked && initialCheckedIds.value.has(t.id))
    .map(t => t.id);
  syncing.value = true;
  try {
    const tasks: Promise<void>[] = [];
    if (checkedIds.length > 0) {
      tasks.push(SyncSkillToTools(props.skill.name, checkedIds));
    }
    if (uncheckedIds.length > 0) {
      tasks.push(UnsyncSkillFromTools(props.skill.name, uncheckedIds));
    }
    await Promise.all(tasks);
    message.success(appStore.t('sync.syncSuccess'));
    const checkedNames = localTools.value.filter(t => t.checked).map(t => t.name);
    emit('synced', checkedNames);
    emit('update:open', false);
  } catch (e: any) {
    message.error(appStore.t('sync.syncFailed', { error: e?.message || String(e) }));
  } finally {
    syncing.value = false;
  }
};
</script>
