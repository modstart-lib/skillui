<template>
  <div class="flex flex-col space-y-6 relative h-full">
    <!-- Header Area -->
    <div class="flex flex-col gap-4">
      <h1 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
        <Layers :size="24" class="text-emerald-500"/>
        {{ $t('manage.title') }}
      </h1>

      <!-- Toolbar -->
      <div class="flex flex-wrap justify-between items-center gap-4 py-2">
        <div class="flex items-center gap-2">
          <Button @click="refreshList" :loading="loading" class="flex items-center gap-2">
            <template #icon>
              <RefreshCw :size="16"/>
            </template>
            {{ $t('manage.refreshList') }}
          </Button>
          <Button type="primary" @click="isInstallModalOpen = true" class="flex items-center gap-2">
            <template #icon>
              <Plus :size="16"/>
            </template>
            {{ $t('manage.manualAdd') }}
          </Button>
        </div>

        <div class="w-full sm:w-72">
          <Input v-model:value="searchQuery" :placeholder="$t('manage.searchPlaceholder')" allowClear>
            <template #prefix>
              <Search :size="16" class="text-slate-400 mr-1"/>
            </template>
          </Input>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex-1 flex items-center justify-center">
      <span class="text-slate-400 text-sm">{{ $t('manage.loading') }}</span>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredSkills.length === 0"
         class="flex-1 bg-white dark:bg-slate-800/80 rounded-xl border border-slate-200 dark:border-slate-700/50 p-12 flex flex-col items-center justify-center text-center">
      <div class="w-16 h-16 bg-slate-100 dark:bg-slate-700/50 rounded-full flex items-center justify-center mb-4">
        <PackageSearch :size="28" class="text-slate-400"/>
      </div>
      <h3 class="text-lg font-bold text-slate-900 dark:text-white mb-1">
        {{ searchQuery ? $t('manage.emptySearchTitle') : $t('manage.emptyTitle') }}
      </h3>
      <p class="text-slate-500 text-sm max-w-xs mx-auto">
        {{ searchQuery ? $t('manage.emptySearchDesc') : $t('manage.emptyDesc') }}
      </p>
    </div>

    <!-- List Content -->
    <div v-else class="flex-1 overflow-y-auto space-y-3 min-h-0">
      <div
          v-for="skill in filteredSkills"
          :key="skill.name"
          class="bg-white dark:bg-slate-800/80 rounded-xl border border-slate-200 dark:border-slate-700/50 p-5 shadow-sm hover:shadow-md transition-all group">
        <div class="flex flex-col gap-4">
          <!-- Top Row: Info + Actions -->
          <div class="flex justify-between items-start">
            <div class="flex flex-col gap-1">
              <div class="flex items-center gap-2">
                <Zap :size="16" class="text-emerald-500 shrink-0"/>
                <h3 class="text-base font-bold text-slate-900 dark:text-white leading-tight cursor-pointer hover:text-emerald-500 transition-colors"
                    @click="viewDetails(skill)">
                  {{ appStore.skillTitle(skill) }}
                </h3>
                <span v-if="skill.isMarket"
                      class="px-1.5 py-0.5 rounded bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400 text-[10px] font-bold border border-indigo-200 dark:border-indigo-800">
                  MARKET
                </span>
              </div>
              <div class="flex items-center gap-3 text-xs text-slate-500 mt-2">
                <span class="flex items-center gap-1 font-mono bg-slate-100 dark:bg-slate-700/50 px-1.5 py-0.5 rounded text-emerald-600 dark:text-emerald-400">
                  <Hash :size="10" class="shrink-0"/>
                  {{ skill.name }}
                </span>
                <span v-if="skill.owner" class="flex items-center gap-1 text-slate-400">
                  <User :size="11" class="shrink-0"/>
                  {{ skill.owner }}
                </span>
              </div>
              <div class="flex items-center gap-1 truncate text-xs font-mono text-slate-500 mt-2" :title="skill.location">
                <FolderOpen :size="12" class="shrink-0 text-slate-400"/>
                <span class="truncate">{{ skill.location }}</span>
              </div>
            </div>

            <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
              <Button type="text" size="small" @click="viewDetails(skill)" :title="$t('actions.view')"
                      class="text-slate-400 hover:text-emerald-500">
                <Eye :size="18"/>
              </Button>
              <Button type="text" danger size="small" :title="$t('actions.delete')" class="opacity-80 hover:opacity-100"
                      @click="confirmDelete(skill)">
                <Trash2 :size="18"/>
              </Button>
            </div>
          </div>

          <!-- Bottom Row: Metadata -->
          <div
              class="flex flex-col gap-y-2 text-xs text-slate-500 border-t border-slate-100 dark:border-slate-700/50 pt-3 mt-1">
            <div class="flex items-center gap-2">
              <span class="flex items-center gap-1 text-slate-400 shrink-0"><Wrench :size="11"/>{{ $t('manage.syncedTools') }}</span>
              <div class="flex flex-wrap items-center gap-2">
                <template v-if="skill.syncedTools && skill.syncedTools.length > 0">
                  <span v-for="tool in skill.syncedTools" :key="tool"
                       class="px-2 py-0.5 text-[11px] rounded-full bg-emerald-50 dark:bg-emerald-900/20 text-emerald-600 dark:text-emerald-400 border border-emerald-100 dark:border-emerald-800 flex items-center gap-1">
                    <CheckCircle2 :size="10"/>
                    {{ tool }}
                  </span>
                </template>
                <span v-else class="text-slate-400 text-[10px]">{{ $t('manage.noSyncedTools') }}</span>
                <Button type="text" size="small"
                        class="text-slate-400 hover:text-emerald-500 px-1 h-5 flex items-center"
                        @click="openSyncModal(skill)">
                  <RefreshCw :size="12"/>
                  <span class="ml-1 text-[10px]">{{ $t('manage.syncBtn') }}</span>
                </Button>
              </div>
            </div>
            <div v-if="skill.updatedAt" class="flex items-center gap-2">
              <span class="flex items-center gap-1 text-slate-400"><Clock :size="11"/>{{ $t('manage.updatedAt') }}</span>
              <span>{{ skill.updatedAt }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Install Modal -->
    <Modal
        v-model:open="isInstallModalOpen"
        :title="$t('manage.installModalTitle')"
        :footer="null"
        :width="800"
        destroyOnClose
    >
      <div class="pt-2">
        <!-- Tabs -->
        <div class="flex p-1 bg-slate-100 dark:bg-slate-800 rounded-xl mb-6">
          <button
              v-for="tab in installTabs"
              :key="tab.id"
              @click="activeInstallTab = tab.id"
              class="flex-1 flex items-center justify-center gap-2 py-2 px-3 rounded-lg text-sm font-medium transition-all"
              :class="activeInstallTab === tab.id ? 'bg-white dark:bg-emerald-600 text-emerald-600 dark:text-white shadow-sm ring-1 ring-black/5' : 'text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200'"
          >
            <component :is="tab.icon" :size="16"/>
            {{ tab.name }}
          </button>
        </div>

        <!-- Tab Content -->
        <div class="min-h-[220px]">
          <!-- Local Folder -->
          <div v-if="activeInstallTab === 'local'" class="space-y-4">
            <div class="grid grid-cols-2 gap-3">
              <div
                  class="border-2 border-dashed border-slate-300 dark:border-slate-700 rounded-xl p-6 flex flex-col items-center justify-center text-center hover:border-emerald-500 dark:hover:border-emerald-500 transition-colors cursor-pointer group bg-slate-50/50 dark:bg-slate-800/30"
                  @click="handleSelectDirectory">
                <div
                    class="w-10 h-10 bg-white dark:bg-slate-800 rounded-full flex items-center justify-center shadow-sm mb-2 group-hover:scale-110 transition-transform">
                  <FolderOpen :size="18" class="text-emerald-500"/>
                </div>
                <p class="text-slate-900 dark:text-white font-medium text-sm mb-0.5">{{ $t('manage.selectFolderTitle') }}</p>
                <p class="text-xs text-slate-500">{{ $t('manage.selectFolderDesc') }}</p>
              </div>
              <div
                  class="border-2 border-dashed border-slate-300 dark:border-slate-700 rounded-xl p-6 flex flex-col items-center justify-center text-center hover:border-emerald-500 dark:hover:border-emerald-500 transition-colors cursor-pointer group bg-slate-50/50 dark:bg-slate-800/30"
                  @click="handleSelectZipFile">
                <div
                    class="w-10 h-10 bg-white dark:bg-slate-800 rounded-full flex items-center justify-center shadow-sm mb-2 group-hover:scale-110 transition-transform">
                  <Upload :size="18" class="text-emerald-500"/>
                </div>
                <p class="text-slate-900 dark:text-white font-medium text-sm mb-0.5">{{ $t('manage.selectZipTitle') }}</p>
                <p class="text-xs text-slate-500">{{ $t('manage.selectZipDesc') }}</p>
              </div>
            </div>
            <p v-if="localPathSelected"
               class="text-xs text-emerald-600 font-mono bg-emerald-50 dark:bg-emerald-900/20 px-3 py-2 rounded flex items-center gap-1">
              <CheckCircle2 :size="12"/>
              {{ localPathSelected }}
            </p>
          </div>

          <!-- Remote URL -->
          <div v-if="activeInstallTab === 'remote'" class="space-y-5">
            <div>
              <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">{{ $t('manage.remoteUrlLabel') }}</label>
              <Input v-model:value="remoteUrl" :placeholder="$t('manage.remoteUrlPlaceholder')">
                <template #prefix>
                  <GitBranch v-if="isGitUrl" :size="18" class="text-emerald-500"/>
                  <Link2 v-else :size="18" class="text-slate-400"/>
                </template>
              </Input>
              <p v-if="isGitUrl" class="text-xs text-emerald-600 dark:text-emerald-400 mt-2 ml-1 flex items-center gap-1">
                <GitBranch :size="12"/>
                {{ $t('manage.remoteGitHint') }}
              </p>
              <p v-else class="text-xs text-slate-500 mt-2 ml-1">{{ $t('manage.remoteUrlHint') }}</p>
            </div>
            <div v-if="!isGitUrl">
              <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">{{ $t('manage.nameLabel') }}</label>
              <Input v-model:value="remoteName" placeholder="my-skill"/>
            </div>
          </div>

          <!-- Paste Text -->
          <div v-if="activeInstallTab === 'text'" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">{{ $t('manage.nameLabel') }}</label>
              <Input v-model:value="textName" placeholder="my-skill"/>
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">{{ $t('manage.textContentLabel') }}</label>
              <Input.TextArea v-model:value="textContent" :rows="7" :placeholder="$t('manage.textContentPlaceholder')"/>
            </div>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="mt-6 flex justify-end gap-3 pt-4 border-t border-slate-100 dark:border-slate-800">
          <Button @click="isInstallModalOpen = false">{{ $t('actions.cancel') }}</Button>
          <Button type="primary" :loading="installing" @click="handleInstall">{{ $t('manage.confirmInstall') }}</Button>
        </div>
      </div>
    </Modal>

    <!-- Details Modal -->
    <SkillDetailsModal
        v-model:open="isDetailsModalOpen"
        :skill="selectedSkill"
        :is-installed="true"
    />

    <!-- Sync Modal -->
    <SkillSync
        v-model:open="isSyncModalOpen"
        :skill="syncTargetSkill"
        @synced="handleSynced"
    />
  </div>
</template>

<script setup lang="ts">
import SkillDetailsModal from '@/views/Manage/DetailModal.vue';
import SkillSync from '@/views/Manage/SyncModal.vue';
import { Button, Input, Modal, message } from 'ant-design-vue';
import {
  CheckCircle2,
  Clock,
  Eye,
  FileText,
  FolderOpen,
  GitBranch,
  Hash,
  Layers,
  Link2,
  PackageSearch,
  Plus,
  RefreshCw,
  Search,
  Trash2,
  Upload,
  User,
  Wrench,
  Zap
} from 'lucide-vue-next';
import { computed, onMounted, ref, watch } from 'vue';
import {
  DeleteSkill,
  InstallSkillFromGit,
  InstallSkillFromLocalPath,
  InstallSkillFromText,
  InstallSkillFromUrl,
  ListLocalSkills,
  SelectDirectory,
  SelectZipFile
} from '../../wailsjs/go/main/App';
import { useAppStore } from '../stores/app';

const appStore = useAppStore();

const loading = ref(false);
const installing = ref(false);
const isInstallModalOpen = ref(false);
const isDetailsModalOpen = ref(false);
const activeInstallTab = ref('local');
const searchQuery = ref('');
const skills = ref<any[]>([]);

// Local install
const localPathSelected = ref('');
// Remote install
const remoteUrl = ref('');
const remoteName = ref('');
// Text install
const textName = ref('');
const textContent = ref('');

const isSyncModalOpen = ref(false);
const syncTargetSkill = ref<any>(null);
const selectedSkill = ref<any>(null);

// 检测是否是 git 仓库 URL
const isGitUrl = computed(() => {
  const url = remoteUrl.value.trim();
  if (!url) return false;
  if (url.endsWith('.git')) return true;
  const gitHosts = ['github.com', 'gitlab.com', 'gitee.com', 'bitbucket.org'];
  if (gitHosts.some(h => url.includes(h)) && !url.match(/\.(zip|tar\.gz|tgz)$/i)) return true;
  return false;
});

const filteredSkills = computed(() => {
  if (!searchQuery.value.trim()) return skills.value;
  const q = searchQuery.value.toLowerCase();
  return skills.value.filter(s =>
      (s.name || '').toLowerCase().includes(q) ||
      (s.title || '').toLowerCase().includes(q) ||
      (s.titleZh || '').toLowerCase().includes(q) ||
      (s.titleEn || '').toLowerCase().includes(q) ||
      (s.owner || '').toLowerCase().includes(q)
  );
});

const refreshList = async () => {
  loading.value = true;
  try {
    const result = await ListLocalSkills();
    skills.value = result || [];
  } catch (e: any) {
    message.error(appStore.t('manage.loadFailed', { error: e?.message || String(e) }));
  } finally {
    loading.value = false;
  }
};

onMounted(() => refreshList());

// 监听来自技能市场的安装事件，自动刷新列表
watch(() => appStore.skillInstallVersion, (val) => {
  if (val > 0) refreshList();
});

const viewDetails = (skill: any) => {
  selectedSkill.value = skill;
  isDetailsModalOpen.value = true;
};

const openSyncModal = (skill: any) => {
  syncTargetSkill.value = skill;
  isSyncModalOpen.value = true;
};

const handleSynced = (checkedNames: string[]) => {
  if (syncTargetSkill.value) {
    const idx = skills.value.findIndex(s => s.name === syncTargetSkill.value.name);
    if (idx >= 0) {
      skills.value[idx] = {...skills.value[idx], syncedTools: checkedNames};
    }
  }
};

const confirmDelete = (skill: any) => {
    Modal.confirm({
    title: appStore.t('manage.confirmDeleteTitle'),
    content: appStore.t('manage.confirmDeleteContent', { title: appStore.skillTitle(skill), name: skill.name }),
    okText: appStore.t('manage.confirmDeleteOk'),
    okType: 'danger',
    cancelText: appStore.t('actions.cancel'),
    async onOk() {
      try {
        await DeleteSkill(skill.name);
        message.success(appStore.t('manage.skillDeleted'));
        await refreshList();
      } catch (e: any) {
        message.error(appStore.t('manage.deleteFailed', { error: e?.message || String(e) }));
      }
    },
  });
};

const handleSelectDirectory = async () => {
  try {
    const dir = await SelectDirectory();
    if (dir) localPathSelected.value = dir;
  } catch (e: any) {
    message.error(appStore.t('manage.selectFailed', { error: e?.message || String(e) }));
  }
};

const handleSelectZipFile = async () => {
  try {
    const file = await SelectZipFile();
    if (file) localPathSelected.value = file;
  } catch (e: any) {
    message.error(appStore.t('manage.selectFailed', { error: e?.message || String(e) }));
  }
};

const installTabs = computed(() => [
  {id: 'local', name: appStore.t('manage.tabLocal'), icon: Upload},
  {id: 'remote', name: appStore.t('manage.tabRemote'), icon: Link2},
  {id: 'text', name: appStore.t('manage.tabText'), icon: FileText},
]);

const handleInstall = async () => {
  installing.value = true;
  try {
    if (activeInstallTab.value === 'local') {
      if (!localPathSelected.value) {
        message.warning(appStore.t('manage.selectFileFirst'));
        return;
      }
      await InstallSkillFromLocalPath(localPathSelected.value);
    } else if (activeInstallTab.value === 'remote') {
      if (!remoteUrl.value.trim()) {
        message.warning(appStore.t('manage.enterRemoteUrl'));
        return;
      }
      if (isGitUrl.value) {
        const installedNames = await InstallSkillFromGit(remoteUrl.value.trim());
        message.success(appStore.t('manage.gitInstallSuccess', { count: installedNames.length }));
        isInstallModalOpen.value = false;
        remoteUrl.value = '';
        remoteName.value = '';
        await refreshList();
        return;
      }
      const name = remoteName.value.trim() || remoteUrl.value.split('/').pop()?.replace(/\.zip$/, '') || 'skill';
      await InstallSkillFromUrl(remoteUrl.value.trim(), name);
    } else {
      if (!textName.value.trim()) {
        message.warning(appStore.t('manage.enterSkillName'));
        return;
      }
      if (!textContent.value.trim()) {
        message.warning(appStore.t('manage.enterSkillContent'));
        return;
      }
      await InstallSkillFromText(textName.value.trim(), textContent.value);
    }
    message.success(appStore.t('manage.skillInstalled'));
    isInstallModalOpen.value = false;
    localPathSelected.value = '';
    remoteUrl.value = '';
    remoteName.value = '';
    textName.value = '';
    textContent.value = '';
    await refreshList();
  } catch (e: any) {
    message.error(appStore.t('manage.installFailed', { error: e?.message || String(e) }));
  } finally {
    installing.value = false;
  }
};
</script>
