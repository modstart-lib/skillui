<script lang="ts" setup>
import { ConfigProvider, theme } from 'ant-design-vue'
import enUS from 'ant-design-vue/es/locale/en_US'
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import { Compass, Layers, Settings, Wrench } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import BrandLogo from './components/BrandLogo.vue'
import { useAppStore } from './stores/app'
import { trackVisit } from './utils/analytics'
import { autoCheckVersion } from './utils/version'
import SkillsManagement from './views/Manage.vue'
import SettingsPage from './views/Setting.vue'
import ToolSettings from './views/SkillSetting.vue'
import SkillsDiscovery from './views/Store.vue'

const appStore = useAppStore()

// Track main page visit on mount
onMounted(async () => {
  // Initialize settings from backend
  await appStore.initSettings()

  trackVisit('Main')

  // Auto check version after 5 seconds
  autoCheckVersion(5000)
})

const antLocale = computed(() => {
  return appStore.locale === 'zh' ? zhCN : enUS
})

const themeConfig = computed(() => ({
  algorithm: appStore.isDark ? theme.darkAlgorithm : theme.defaultAlgorithm,
  token: {
    colorPrimary: '#10b981',
    borderRadius: 8,
    fontFamily: 'Inter, system-ui, -apple-system, BlinkMacSystemFont, sans-serif',
  },
}))

// Active tab state
const activeTab = ref('management')
</script>

<template>
  <ConfigProvider :locale="antLocale" :theme="themeConfig">
    <div class="app-shell">
      <div class="app-layout">
        <!-- Left Sidebar with Vertical Tabs - RESTORED ORIGINAL STYLE -->
        <div class="sidebar">
          <!-- Brand -->
          <div class="flex flex-col items-center justify-center py-6 mb-2">
            <BrandLogo class="w-10 h-10"/>
            <div class="mt-2 font-bold text-lg tracking-wide">
              <span class="text-transparent bg-clip-text bg-gradient-to-r from-gray-800 to-gray-900">Skill</span>
              <span class="text-transparent bg-clip-text bg-gradient-to-r from-gray-800 to-gray-900">UI</span>
            </div>
          </div>

          <div class="sidebar-tabs">
            <button
                class="tab-button"
                :class="{ active: activeTab === 'management' }"
                @click="activeTab = 'management'"
            >
              <Layers :size="20"/>
              <span class="tab-label">{{ appStore.t('sidebar.nav.management') }}</span>
            </button>
            <button
                class="tab-button"
                :class="{ active: activeTab === 'discovery' }"
                @click="activeTab = 'discovery'"
            >
              <Compass :size="20"/>
              <span class="tab-label">{{ appStore.t('sidebar.nav.discovery') }}</span>
            </button>
            <button
                class="tab-button"
                :class="{ active: activeTab === 'tools' }"
                @click="activeTab = 'tools'"
            >
              <Wrench :size="20"/>
              <span class="tab-label">{{ appStore.t('sidebar.nav.tools') }}</span>
            </button>

            <!-- Spacer to push settings to bottom -->
            <div class="flex-1"></div>

            <button
                class="tab-button mt-auto"
                :class="{ active: activeTab === 'settings' }"
                @click="activeTab = 'settings'"
            >
              <Settings :size="20"/>
              <span class="tab-label">{{ appStore.t('sidebar.nav.settings') }}</span>
            </button>
          </div>
        </div>

        <!-- Main Content Area -->
        <div class="main-content">
          <div v-show="activeTab === 'management'" class="content-view">
            <SkillsManagement/>
          </div>
          <div v-show="activeTab === 'discovery'" class="content-view">
            <SkillsDiscovery/>
          </div>
          <div v-show="activeTab === 'tools'" class="content-view">
            <ToolSettings/>
          </div>
          <div v-show="activeTab === 'settings'" class="content-view">
            <SettingsPage/>
          </div>
        </div>
      </div>
    </div>
  </ConfigProvider>
</template>

<style scoped>
/* Any additional overrides used to match skillui style if strictly needed,
   but Tailwind utility classes should cover it.
   The user asked to keep style consistent with SkillUI but provided images.
   I am using the dark theme from images "SkillSLM" (black sidebar, dark content).
   The "style.css" has dark mode media query, so this fits. */
</style>
