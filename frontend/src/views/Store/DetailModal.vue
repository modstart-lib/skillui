<template>
  <Modal
      :open="open"
      :title="displayTitle"
      @update:open="$emit('update:open', $event)"
      :footer="null"
      width="680px"
      destroyOnClose
      centered
  >
    <div v-if="skill" class="pt-4">
      <!-- Basic Info - 垂直排列 -->
      <div class="space-y-4 mb-6">
        <div class="flex flex-col gap-3 text-sm">
          <!-- 名称 -->
          <div>
            <div class="text-[11px] uppercase tracking-wider text-slate-400 mb-1">{{ $t('skill.detail.labelTitle') }}</div>
            <div class="font-semibold text-slate-900 dark:text-white">{{ displayTitle }}</div>
          </div>
          <!-- 标识 -->
          <div>
            <div class="text-[11px] uppercase tracking-wider text-slate-400 mb-1">{{ $t('skill.detail.labelName') }}</div>
            <div class="font-mono text-emerald-600 dark:text-emerald-400 bg-slate-100 dark:bg-slate-800 px-2 py-0.5 rounded inline-block text-xs">
              {{ skill.name }}
            </div>
          </div>
          <!-- 来源 -->
          <div>
            <div class="text-[11px] uppercase tracking-wider text-slate-400 mb-1">{{ $t('skill.detail.labelOwner') }}</div>
            <div class="font-mono text-xs text-slate-600 dark:text-slate-300">{{ skill.owner }}</div>
          </div>
          <!-- 版本 -->
          <div v-if="skill.version">
            <div class="text-[11px] uppercase tracking-wider text-slate-400 mb-1">{{ $t('skill.detail.labelVersion') }}</div>
            <div class="font-mono text-xs text-slate-600 dark:text-slate-300">{{ skill.version }}</div>
          </div>
          <!-- 下载次数 -->
          <div v-if="skill.downloadCount != null">
            <div class="text-[11px] uppercase tracking-wider text-slate-400 mb-1">{{ $t('skill.detail.labelDownloadCount') }}</div>
            <div class="text-xs text-slate-600 dark:text-slate-300">{{ skill.downloadCount }}</div>
          </div>
        </div>

        <!-- Tags -->
        <div v-if="displayTags && displayTags.length > 0">
          <div class="text-[11px] uppercase tracking-wider text-slate-400 mb-2">{{ $t('skill.detail.labelTags') }}</div>
          <div class="flex flex-wrap gap-1.5">
            <span
                v-for="tag in displayTags"
                :key="tag"
                class="px-2 py-0.5 text-[11px] rounded-full bg-slate-100 dark:bg-slate-700 text-slate-600 dark:text-slate-300"
            >{{ tag }}</span>
          </div>
        </div>

        <!-- Description -->
        <div>
          <div class="text-[11px] uppercase tracking-wider text-slate-400 mb-2">{{ $t('skill.detail.labelDescription') }}</div>
          <div class="bg-slate-50 dark:bg-slate-800/50 p-4 rounded-lg border border-slate-100 dark:border-slate-800">
            <p class="text-sm text-slate-700 dark:text-slate-300 leading-relaxed whitespace-pre-wrap">{{ displayDesc }}</p>
          </div>
        </div>
      </div>

      <!-- Preview Content (Markdown from detail API) -->
      <div class="mt-4">
        <div class="text-[11px] uppercase tracking-wider text-slate-400 mb-2">{{ $t('skill.detail.labelContent') }}</div>
        <div v-if="detailLoading" class="bg-slate-50 dark:bg-slate-800/50 rounded-lg border border-slate-100 dark:border-slate-800 p-4 space-y-2 animate-pulse">
          <div class="h-3 bg-slate-200 dark:bg-slate-700 rounded w-3/4"></div>
          <div class="h-3 bg-slate-200 dark:bg-slate-700 rounded w-full"></div>
          <div class="h-3 bg-slate-200 dark:bg-slate-700 rounded w-5/6"></div>
        </div>
        <div
            v-else-if="previewHtml"
            class="bg-slate-50 dark:bg-slate-800/50 rounded-lg border border-slate-100 dark:border-slate-800 p-4 max-h-72 overflow-y-auto markdown-body text-sm"
            v-html="previewHtml"
        />
        <div v-else class="bg-slate-50 dark:bg-slate-800/50 rounded-lg border border-slate-100 dark:border-slate-800 p-4">
          <p class="text-sm text-slate-400 text-center">{{ $t('skill.detail.noContent') }}</p>
        </div>
      </div>

      <div class="mt-6 pt-5 border-t border-slate-100 dark:border-slate-800 flex justify-end gap-3">
        <Button @click="$emit('update:open', false)">{{ $t('skill.detail.close') }}</Button>
        <Button v-if="!isInstalled" type="primary" class="px-6" @click="handleInstall" :loading="installing">
          {{ $t('skill.detail.install') }}
        </Button>
        <span v-else class="flex items-center gap-1.5 text-sm text-slate-400 bg-slate-100 dark:bg-slate-700 px-4 py-1.5 rounded">
          <CheckCircle2 :size="14" class="text-emerald-500"/>
          {{ $t('skill.detail.installed') }}
        </span>
      </div>
    </div>
  </Modal>
</template>

<script setup lang="ts">
import { Button, Modal } from 'ant-design-vue';
import { CheckCircle2 } from 'lucide-vue-next';
import { computed, ref, watch } from 'vue';
import { skillUiDetail } from '../../api';
import { useAppStore } from '../../stores/app';

const appStore = useAppStore();

const props = defineProps<{
  open: boolean;
  skill: Record<string, any> | null;
  isInstalled?: boolean;
  installing?: boolean;
}>();

const emit = defineEmits(['update:open', 'install']);

const detailLoading = ref(false);
const previewHtml = ref('');

// 打开时调用详情接口
watch(() => props.open, async (val) => {
  if (!val || !props.skill) {
    previewHtml.value = '';
    return;
  }
  detailLoading.value = true;
  previewHtml.value = '';
  try {
    const json = await skillUiDetail(props.skill.id);
    if (json.code === 0 && json.data?.record?._preview) {
      previewHtml.value = json.data.record._preview;
    }
  } catch {
    // ignore
  } finally {
    detailLoading.value = false;
  }
});

const displayTitle = computed(() => {
  if (!props.skill) return '';
  return appStore.skillTitle(props.skill);
});

const displayDesc = computed(() => {
  if (!props.skill) return '';
  return appStore.skillDesc(props.skill);
});

const displayTags = computed(() => {
  if (!props.skill) return [];
  const locale = appStore.locale;
  if (locale === 'zh') {
    return props.skill.tagsZh?.length ? props.skill.tagsZh : (props.skill.tagsEn || props.skill.tags || []);
  }
  return props.skill.tagsEn?.length ? props.skill.tagsEn : (props.skill.tagsZh || props.skill.tags || []);
});

const handleInstall = () => {
  if (!props.isInstalled) {
    emit('install', props.skill);
  }
};
</script>

<style scoped>
.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3),
.markdown-body :deep(h4) {
  font-weight: 600;
  margin-top: 1em;
  margin-bottom: 0.4em;
  color: inherit;
}
.markdown-body :deep(h1) { font-size: 1.2em; }
.markdown-body :deep(h2) { font-size: 1.1em; }
.markdown-body :deep(h3) { font-size: 1em; }
.markdown-body :deep(p) {
  margin-top: 0.4em;
  margin-bottom: 0.4em;
  line-height: 1.6;
  color: inherit;
}
.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  padding-left: 1.4em;
  margin: 0.4em 0;
}
.markdown-body :deep(li) {
  margin: 0.2em 0;
  line-height: 1.6;
}
.markdown-body :deep(code) {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 0.85em;
  background: rgba(0,0,0,0.06);
  padding: 0.1em 0.4em;
  border-radius: 3px;
}
.markdown-body :deep(pre) {
  background: rgba(0,0,0,0.06);
  border-radius: 6px;
  padding: 0.75em 1em;
  overflow-x: auto;
  margin: 0.6em 0;
}
.markdown-body :deep(pre code) {
  background: none;
  padding: 0;
}
.markdown-body :deep(blockquote) {
  border-left: 3px solid #10b981;
  padding-left: 0.8em;
  margin: 0.6em 0;
  color: #64748b;
}
.markdown-body :deep(hr) {
  border: none;
  border-top: 1px solid #e2e8f0;
  margin: 0.8em 0;
}
.markdown-body :deep(strong) { font-weight: 600; }
.markdown-body :deep(a) { color: #10b981; text-decoration: underline; }
.markdown-body :deep(table) { width: 100%; border-collapse: collapse; font-size: 0.9em; }
.markdown-body :deep(th),
.markdown-body :deep(td) {
  border: 1px solid #e2e8f0;
  padding: 0.3em 0.6em;
  text-align: left;
}
.markdown-body :deep(th) { background: rgba(0,0,0,0.04); font-weight: 600; }
</style>
