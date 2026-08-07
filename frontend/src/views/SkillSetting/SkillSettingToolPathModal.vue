<template>
  <Modal
      :open="visible"
      :title="$t('toolSettings.toolPathTitle')"
      :okText="$t('actions.save')"
      :cancelText="$t('actions.cancel')"
      :confirm-loading="saving"
      width="min(600px, 90vw)"
      @ok="handleSave"
      @update:open="onOpenChange"
  >
    <div class="py-4 space-y-4">
      <div v-if="toolName" class="flex items-center gap-2 pb-2 border-b border-slate-100 dark:border-slate-700/50">
        <component :is="iconForTool(toolId)" :size="18" class="text-emerald-500"/>
        <span class="font-bold text-slate-900 dark:text-white">{{ toolName }}</span>
        <span class="text-xs text-slate-400 font-mono">{{ toolId }}</span>
      </div>

      <div>
        <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">
          {{ $t('toolSettings.toolPathLabel') }}
        </label>
        <div class="flex gap-2">
          <Input v-model:value="path" :placeholder="$t('toolSettings.toolPathPlaceholder')"/>
          <Button @click="handleSelectDir" :loading="selecting">{{ $t('toolSettings.selectDir') }}</Button>
        </div>
        <p class="text-xs text-slate-400 mt-2">{{ $t('toolSettings.toolPathDesc') }}</p>
      </div>
    </div>
  </Modal>
</template>

<script setup lang="ts">
import { Button, Input, Modal, message } from 'ant-design-vue';
import { Box, Code, Command, Terminal, Zap } from 'lucide-vue-next';
import { ref, watch } from 'vue';
import { SelectDirectory, SetToolPath } from '../../../wailsjs/go/main/App';
import { useAppStore } from '../../stores/app';

const appStore = useAppStore();

const props = defineProps<{
  visible: boolean;
  toolId: string;
  toolName: string;
  defaultValue: string;
}>();

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
  (e: 'saved', toolId: string, path: string): void;
}>();

const path = ref('');
const saving = ref(false);
const selecting = ref(false);

watch(() => props.visible, (val) => {
  if (val) {
    path.value = props.defaultValue || '';
  }
});

const onOpenChange = (val: boolean) => {
  emit('update:visible', val);
};

const iconForTool = (id: string) => {
  switch (id) {
    case 'cursor': return Zap;
    case 'claude_code': case 'gemini_cli': case 'codex': case 'qwen_code': case 'crush': return Terminal;
    case 'kilo_code': case 'roo_code': case 'goose': case 'auggie_cli': return Command;
    case 'zed': case 'github_copilot': case 'cline': case 'codebuddy': case 'iflow': return Code;
    default: return Box;
  }
};

const handleSelectDir = async () => {
  selecting.value = true;
  try {
    const dir = await SelectDirectory();
    if (dir) {
      path.value = dir;
    }
  } catch (e) {
    message.error(appStore.t('toolSettings.selectDirFailed'));
  } finally {
    selecting.value = false;
  }
};

const handleSave = async () => {
  if (!path.value.trim()) {
    message.warning(appStore.t('toolSettings.enterValidDir'));
    return;
  }
  saving.value = true;
  try {
    await SetToolPath(props.toolId, path.value.trim());
    emit('saved', props.toolId, path.value.trim());
    emit('update:visible', false);
    message.success(appStore.t('toolSettings.toolPathUpdated'));
  } catch (e: any) {
    message.error(appStore.t('toolSettings.toolPathFailed', { error: e?.message || String(e) }));
  } finally {
    saving.value = false;
  }
};
</script>
