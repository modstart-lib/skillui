<template>
  <!-- Default Config Section -->
  <div
      class="bg-white dark:bg-slate-800/80 rounded-xl border border-slate-200 dark:border-slate-700/50 p-6 shadow-sm">
    <h3 class="text-base font-bold text-slate-900 dark:text-white mb-4 flex items-center gap-2">
      <FolderOpen :size="18" class="text-emerald-500"/>
      {{ $t('toolSettings.globalTitle') }}
    </h3>
    <div class="max-w-xl">
      <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2 flex items-center gap-1.5">
        <Folder :size="14" class="text-slate-400"/>
        {{ $t('toolSettings.defaultSkillDir') }}
      </label>
      <div class="flex gap-2">
        <Input :value="modelValue" readonly
               class="bg-slate-50 dark:bg-slate-900 text-slate-500 cursor-not-allowed"/>
        <Button @click="openModifyDirModal">{{ $t('toolSettings.modify') }}</Button>
      </div>
      <p class="text-xs text-slate-500 mt-1">{{ $t('toolSettings.skillDirHint') }}</p>
    </div>
  </div>

  <!-- Modify Directory Modal -->
  <Modal
      v-model:open="isModifyDirModalOpen"
      :title="$t('toolSettings.modifyDirTitle')"
      :confirm-loading="saving"
      @ok="handleModifyDir"
      :okText="$t('toolSettings.confirmModify')"
      :cancelText="$t('actions.cancel')"
  >
    <div class="py-4 space-y-4">
      <div>
        <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">{{ $t('toolSettings.newDirLabel') }}</label>
        <div class="flex gap-2">
          <Input v-model:value="newSkillDir" :placeholder="$t('toolSettings.newDirPlaceholder')"/>
          <Button @click="handleSelectDir" :loading="selecting">{{ $t('toolSettings.selectDir') }}</Button>
        </div>
      </div>

      <div class="bg-amber-50 dark:bg-amber-900/20 p-4 rounded-lg border border-amber-200 dark:border-amber-800/50">
        <label class="flex items-start gap-2 cursor-pointer">
          <input type="checkbox" v-model="migrateSkills"
                 class="mt-1 w-4 h-4 text-emerald-600 rounded border-slate-300 focus:ring-emerald-500">
          <div class="text-sm flex-1">
            <div class="flex items-center gap-1.5 font-bold text-amber-800 dark:text-amber-400 mb-1">
              <AlertTriangle :size="14"/>
              {{ $t('toolSettings.migrateTitle') }}
            </div>
            <p class="text-amber-700/80 dark:text-amber-500/80 text-xs">
              {{ $t('toolSettings.migrateDesc') }}
            </p>
          </div>
        </label>
      </div>
    </div>
  </Modal>
</template>

<script setup lang="ts">
import { Button, Input, Modal, message } from 'ant-design-vue';
import { AlertTriangle, Folder, FolderOpen } from 'lucide-vue-next';
import { ref } from 'vue';
import { SelectDirectory, SetSkillDir } from '../../../wailsjs/go/main/App';
import { useAppStore } from '../../stores/app';

const appStore = useAppStore();

const props = defineProps<{
  modelValue: string;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void;
  (e: 'changed', value: string): void;
}>();

const isModifyDirModalOpen = ref(false);
const newSkillDir = ref('');
const migrateSkills = ref(true);
const saving = ref(false);
const selecting = ref(false);

const openModifyDirModal = () => {
  newSkillDir.value = props.modelValue;
  isModifyDirModalOpen.value = true;
};

const handleSelectDir = async () => {
  selecting.value = true;
  try {
    const dir = await SelectDirectory();
    if (dir) {
      newSkillDir.value = dir;
    }
  } catch (e) {
    message.error(appStore.t('toolSettings.selectDirFailed'));
  } finally {
    selecting.value = false;
  }
};

const handleModifyDir = async () => {
  if (!newSkillDir.value.trim()) {
    message.warning(appStore.t('toolSettings.enterValidDir'));
    return;
  }

  Modal.confirm({
    title: appStore.t('toolSettings.confirmModifyTitle'),
    content: migrateSkills.value
        ? appStore.t('toolSettings.confirmModifyWithMigrate')
        : appStore.t('toolSettings.confirmModifyWithoutMigrate'),
    okText: appStore.t('toolSettings.confirmSubmit'),
    cancelText: appStore.t('actions.cancel'),
    async onOk() {
      saving.value = true;
      try {
        await SetSkillDir(newSkillDir.value, migrateSkills.value);
        emit('update:modelValue', newSkillDir.value);
        emit('changed', newSkillDir.value);
        isModifyDirModalOpen.value = false;
        message.success(appStore.t('toolSettings.dirUpdated'));
      } catch (e: any) {
        message.error(appStore.t('toolSettings.dirUpdateFailed', { error: e?.message || String(e) }));
      } finally {
        saving.value = false;
      }
    }
  });
};
</script>
