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
      <!-- Basic Info -->
      <div class="space-y-4 mb-6">
        <div class="grid grid-cols-3 gap-x-4 gap-y-3 text-sm">
          <div>
            <div class="text-[11px] uppercase tracking-wider text-slate-400 mb-1">{{ $t('skill.detail.labelTitle') }}</div>
            <div class="font-semibold text-slate-900 dark:text-white">{{ displayTitle }}</div>
          </div>
          <div>
            <div class="text-[11px] uppercase tracking-wider text-slate-400 mb-1">{{ $t('skill.detail.labelName') }}</div>
            <div class="font-mono text-emerald-600 dark:text-emerald-400 bg-slate-100 dark:bg-slate-800 px-2 py-0.5 rounded inline-block text-xs">
              {{ skill.name }}
            </div>
          </div>
          <div>
            <div class="text-[11px] uppercase tracking-wider text-slate-400 mb-1">{{ $t('skill.detail.labelOwner') }}</div>
            <div class="font-mono text-xs text-slate-600 dark:text-slate-300 flex items-center gap-1 flex-wrap">
              <span v-if="skill.isMarket"
                    class="px-1.5 py-0.5 rounded bg-indigo-50 dark:bg-indigo-900/30 text-indigo-500 text-[10px] font-bold border border-indigo-200 dark:border-indigo-800">
                MARKET
              </span>
              {{ skill.owner }}
            </div>
          </div>
        </div>

        <!-- Version & Updated At -->
        <div class="flex gap-6 text-xs text-slate-500 flex-wrap">
          <div v-if="skill.version">
            <span class="text-slate-400">{{ $t('skill.detail.labelVersion') }}：</span>
            <span class="font-mono">{{ skill.version }}</span>
          </div>
          <div v-if="skill.updatedAt">
            <span class="text-slate-400">{{ $t('skill.detail.labelUpdatedAt') }}：</span>
            <span>{{ skill.updatedAt }}</span>
          </div>
        </div>

          <!-- Location -->
          <div v-if="skill.location">
            <div class="text-[11px] uppercase tracking-wider text-slate-400 mb-1 flex items-center gap-1">
              <FolderOpen :size="10"/>
              {{ $t('skill.detail.labelLocation') }}
            </div>
          <div class="font-mono text-xs text-slate-500 dark:text-slate-400 bg-slate-100 dark:bg-slate-800 px-3 py-1.5 rounded break-all">
            {{ skill.location }}
          </div>
        </div>

        <!-- Synced Tools -->
        <div>
            <div class="text-[11px] uppercase tracking-wider text-slate-400 mb-2 flex items-center gap-1">
              <Wrench :size="10"/>
              {{ $t('skill.detail.labelSyncedTools') }}
            </div>
            <div class="flex flex-wrap gap-1.5">
              <template v-if="skill.syncedTools && skill.syncedTools.length > 0">
                <span
                    v-for="tool in skill.syncedTools"
                    :key="tool"
                    class="px-2 py-0.5 text-[11px] rounded-full bg-emerald-50 dark:bg-emerald-900/20 text-emerald-600 dark:text-emerald-400 border border-emerald-100 dark:border-emerald-800 flex items-center gap-1"
                >
                  <CheckCircle2 :size="10"/>
                  {{ tool }}
                </span>
              </template>
              <span v-else class="text-xs text-slate-400">{{ $t('skill.detail.noSyncedTools') }}</span>
            </div>
          </div>

          <!-- Tags -->
          <div v-if="displayTags && displayTags.length > 0">
            <div class="text-[11px] uppercase tracking-wider text-slate-400 mb-2 flex items-center gap-1">
              <Tag :size="10"/>
              {{ $t('skill.detail.labelTags') }}
            </div>
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
            <div class="text-[11px] uppercase tracking-wider text-slate-400 mb-2 flex items-center gap-1">
              <AlignLeft :size="10"/>
              {{ $t('skill.detail.labelDescription') }}
            </div>
          <div class="bg-slate-50 dark:bg-slate-800/50 p-4 rounded-lg border border-slate-100 dark:border-slate-800">
            <p class="text-sm text-slate-700 dark:text-slate-300 leading-relaxed whitespace-pre-wrap">{{ displayDesc }}</p>
          </div>
        </div>
      </div>

      <!-- Skill Content (Markdown) -->
      <div v-if="skillMarkdownHtml" class="mt-6">
        <div class="text-[11px] uppercase tracking-wider text-slate-400 mb-2 flex items-center gap-1">
          <FileCode2 :size="10"/>
          {{ $t('skill.detail.labelContent') }}
        </div>
        <div
            class="bg-slate-50 dark:bg-slate-800/50 rounded-lg border border-slate-100 dark:border-slate-800 p-4 max-h-72 overflow-y-auto markdown-body text-sm"
            v-html="skillMarkdownHtml"
        />
      </div>

      <div class="mt-6 pt-5 border-t border-slate-100 dark:border-slate-800 flex justify-end">
        <Button @click="$emit('update:open', false)">{{ $t('skill.detail.close') }}</Button>
      </div>
    </div>
  </Modal>
</template>

<script setup lang="ts">
import { Button, Modal } from 'ant-design-vue';
import { AlignLeft, CheckCircle2, FileCode2, FolderOpen, Tag, Wrench } from 'lucide-vue-next';
import { marked } from 'marked';
import { computed } from 'vue';
import { useAppStore } from '../../stores/app';

const appStore = useAppStore();

const props = defineProps<{
  open: boolean;
  skill: Record<string, any> | null;
}>();

defineEmits(['update:open']);

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

const skillMarkdownHtml = computed(() => {
  if (!props.skill?.skillContent) return '';
  let content = props.skill.skillContent as string;
  const frontmatterRegex = /^---[\s\S]*?---\s*/;
  content = content.replace(frontmatterRegex, '').trim();
  if (!content) return '';
  return marked.parse(content) as string;
});
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
