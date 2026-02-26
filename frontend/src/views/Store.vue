<script setup lang="ts">
import { setApiBaseUrl, skillUiDownload, skillUiPaginate } from '@/api';
import SkillDetailsModal from '@/views/Store/DetailModal.vue';
import { Button, Input, Pagination, message } from 'ant-design-vue';
import { Clock, Compass, Download, Flame, Search, SearchX, Tag } from 'lucide-vue-next';
import { onMounted, ref, watch } from 'vue';
import { GetAppConfig, InstallSkillFromMarket, ListLocalSkills } from '../../wailsjs/go/main/App';
import { useAppStore } from '../stores/app';

const appStore = useAppStore();
const isDetailsModalOpen = ref(false);
const selectedSkill = ref<any>(null);
const activeFilter = ref('hot');
const keywords = ref('');
const skills = ref<any[]>([]);
const loading = ref(false);
const installing = ref<Record<number, boolean>>({});
const total = ref(0);
const currentPage = ref(1);
const pageSize = 24;
const installedNames = ref<Set<string>>(new Set());

let searchTimer: ReturnType<typeof setTimeout> | null = null;

onMounted(async () => {
  try {
    const cfg = await GetAppConfig();
    setApiBaseUrl((cfg.apiBaseUrl as string) || 'https://skillui.com/api');
  } catch {
    // 使用默认 baseUrl
  }
  await Promise.all([fetchSkills(1), refreshInstalledList()]);
});

const refreshInstalledList = async () => {
  try {
    const result = await ListLocalSkills();
    installedNames.value = new Set((result || []).map((s: any) => s.name));
  } catch {
    installedNames.value = new Set();
  }
};

const isInstalled = (skill: any) => installedNames.value.has(skill.name);

const fetchSkills = async (page: number) => {
  loading.value = true;
  try {
    const json = await skillUiPaginate({
      sort: activeFilter.value === 'latest' ? 'time' : 'hot',
      keywords: keywords.value.trim(),
      page,
      pageSize,
    });
    if (json.code === 0) {
      skills.value = json.data.records || [];
      total.value = json.data.total || 0;
      currentPage.value = page;
    } else {
      message.error(appStore.t('store.loadFailed', { error: json.message || 'Unknown error' }));
    }
  } catch (e: any) {
    message.error(appStore.t('store.networkFailed'));
  } finally {
    loading.value = false;
  }
};

watch(activeFilter, () => fetchSkills(1));

watch(keywords, () => {
  if (searchTimer) clearTimeout(searchTimer);
  searchTimer = setTimeout(() => fetchSkills(1), 300);
});

const viewDetails = (skill: any) => {
  selectedSkill.value = skill;
  isDetailsModalOpen.value = true;
};

const handleInstall = async (skill: any) => {
  if (isInstalled(skill) || installing.value[skill.id]) return;
  installing.value = {...installing.value, [skill.id]: true};
  try {
    // Step 1: get download URL
    const json = await skillUiDownload(skill.id);
    if (json.code !== 0) throw new Error(json.message || appStore.t('store.getDownloadFailed'));
    const downloadUrl: string = json.data.url;

    // Step 2: download & install via Go, writing skillui.json metadata
    await InstallSkillFromMarket(downloadUrl, {
      name: skill.name,
      marketId: skill.id,
      titleEn: skill.titleEn || '',
      titleZh: skill.titleZh || '',
      descEn: skill.descEn || '',
      descZh: skill.descZh || '',
      owner: skill.owner || '',
      version: skill.version || '',
      title: skill.titleEn || skill.name,
      description: skill.descEn || '',
      tags: skill.tagsEn || [],
      isMarket: true,
      syncedTools: [],
      location: '',
      updatedAt: '',
      skillContent: '',
    });

    installedNames.value = new Set([...installedNames.value, skill.name]);
    message.success(appStore.t('store.installSuccess', { title: appStore.skillTitle(skill) }));
    appStore.notifySkillInstalled();
    isDetailsModalOpen.value = false;
  } catch (e: any) {
    message.error(appStore.t('store.installFailed', { error: e?.message || String(e) }));
  } finally {
    const next = {...installing.value};
    delete next[skill.id];
    installing.value = next;
  }
};
</script>

<template>
  <div class="flex flex-col space-y-6 h-full relative">
    <!-- Header -->
    <div class="flex flex-col gap-4 shrink-0">
      <div class="flex justify-between items-center">
        <h1 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
          <Compass :size="24" class="text-emerald-500"/>
          {{ $t('store.title') }}
        </h1>
      </div>

      <!-- Filter and Search Bar -->
      <div class="flex flex-wrap items-center justify-between gap-4">
        <!-- Filters -->
        <div class="flex bg-slate-100 dark:bg-slate-800 p-1 rounded-lg">
          <button
              @click="activeFilter = 'hot'"
              class="px-3 py-1.5 text-sm font-medium rounded-md transition-all flex items-center gap-2"
              :class="activeFilter === 'hot' ? 'bg-white dark:bg-slate-700 text-slate-900 dark:text-white shadow-sm' : 'text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200'"
          >
            <Flame :size="16" class="text-orange-500"/>
            {{ $t('store.filterHot') }}
          </button>
          <button
              @click="activeFilter = 'latest'"
              class="px-3 py-1.5 text-sm font-medium rounded-md transition-all flex items-center gap-2"
              :class="activeFilter === 'latest' ? 'bg-white dark:bg-slate-700 text-slate-900 dark:text-white shadow-sm' : 'text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200'"
          >
            <Clock :size="16" class="text-blue-500"/>
            {{ $t('store.filterLatest') }}
          </button>
        </div>

        <!-- Search -->
        <div class="w-full sm:w-auto flex-1 max-w-md">
          <Input v-model:value="keywords" :placeholder="$t('store.searchPlaceholder')" allowClear>
            <template #prefix>
              <Search :size="16" class="text-slate-400 mr-1"/>
            </template>
          </Input>
        </div>
      </div>
    </div>

    <!-- Loading Skeleton -->
    <div v-if="loading" class="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-4 pb-6 overflow-y-auto min-h-0">
      <div v-for="i in 8" :key="i"
           class="bg-white dark:bg-slate-800/80 rounded-xl border border-slate-200 dark:border-slate-700/50 p-4 shadow-sm h-44 animate-pulse">
        <div class="h-4 bg-slate-200 dark:bg-slate-700 rounded mb-2 w-3/4"></div>
        <div class="h-3 bg-slate-100 dark:bg-slate-700/50 rounded mb-4 w-1/2"></div>
        <div class="h-3 bg-slate-100 dark:bg-slate-700/50 rounded mb-1 w-full"></div>
        <div class="h-3 bg-slate-100 dark:bg-slate-700/50 rounded mb-1 w-5/6"></div>
        <div class="h-3 bg-slate-100 dark:bg-slate-700/50 rounded w-4/6"></div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="skills.length === 0"
         class="bg-white dark:bg-slate-800/80 rounded-xl border border-slate-200 dark:border-slate-700/50 p-12 flex flex-col items-center justify-center text-center">
      <div class="w-16 h-16 bg-slate-100 dark:bg-slate-700/50 rounded-full flex items-center justify-center mb-4">
        <SearchX :size="28" class="text-slate-400"/>
      </div>
      <h3 class="text-lg font-bold text-slate-900 dark:text-white mb-1">{{ $t('store.emptyTitle') }}</h3>
      <p class="text-slate-500 text-sm max-w-xs mx-auto">{{ $t('store.emptyDesc') }}</p>
      <Button class="mt-4" @click="keywords = ''; activeFilter = 'hot'">{{ $t('store.clearFilter') }}</Button>
    </div>

    <!-- Grid -->
    <template v-else>
      <div class="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-4 overflow-y-auto min-h-0 pb-4">
        <div
            v-for="skill in skills"
            :key="skill.id || skill.name"
            class="bg-white dark:bg-slate-800/80 rounded-xl border border-slate-200 dark:border-slate-700/50 p-4 shadow-sm transition-all hover:shadow-md hover:border-emerald-500/30 flex flex-col cursor-pointer relative"
            @click="viewDetails(skill)"
        >
          <!-- Installed dot -->
          <div v-if="isInstalled(skill)" class="absolute top-3 right-3">
            <span class="w-2 h-2 rounded-full bg-emerald-500 block"></span>
          </div>

          <!-- Title -->
          <div class="mb-2 pr-4">
            <h3 class="font-bold text-slate-900 dark:text-white text-sm leading-tight">
              {{ appStore.skillTitle(skill) }}
            </h3>
          </div>

          <!-- Name & Owner -->
          <div class="flex flex-col gap-0.5 mb-3 truncate">
            <span
                class="font-mono text-[11px] text-emerald-600 dark:text-emerald-400 bg-slate-100 dark:bg-slate-700/50 px-1.5 py-0.5 rounded w-fit">{{
                skill.name
              }}</span>
          </div>

          <!-- Description -->
          <p class="text-xs text-slate-500 dark:text-slate-400 line-clamp-3 flex-1 mb-3">
            {{ appStore.skillDesc(skill) }}
          </p>

          <!-- Footer: version/download + button -->
          <div class="pt-3 border-t border-slate-100 dark:border-slate-700/50 flex items-center justify-between mt-auto">
            <div class="flex items-center gap-3">
              <div v-if="skill.version" class="flex items-center gap-0.5 text-[10px] text-slate-400">
                <Tag :size="10" class="shrink-0"/>
                <span class="font-mono">{{ skill.version }}</span>
              </div>
              <div v-if="skill.downloadCount != null" class="flex items-center gap-0.5 text-[10px] text-slate-400">
                <Download :size="10" class="shrink-0"/>
                <span>{{ skill.downloadCount }}</span>
              </div>
            </div>
            <Button
                v-if="!isInstalled(skill)"
                type="primary"
                size="small"
                class="text-[11px] px-3 h-7"
                :loading="installing[skill.id]"
                @click.stop="handleInstall(skill)"
            >{{ $t('store.install') }}
            </Button>
            <span v-else
                  class="text-[11px] text-slate-400 bg-slate-100 dark:bg-slate-700/50 px-2 py-1 rounded select-none">{{ $t('store.installed') }}</span>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div class="flex justify-center shrink-0 pb-2">
        <Pagination
            v-model:current="currentPage"
            :total="total"
            :page-size="pageSize"
            :show-size-changer="false"
            show-less-items
            @change="(page: number) => fetchSkills(page)"
        />
      </div>
    </template>

    <!-- Details Modal -->
    <SkillDetailsModal
        v-model:open="isDetailsModalOpen"
        :skill="selectedSkill"
        :is-installed="selectedSkill ? isInstalled(selectedSkill) : false"
        :installing="selectedSkill ? !!installing[selectedSkill?.id] : false"
        @install="handleInstall"
    />
  </div>
</template>
