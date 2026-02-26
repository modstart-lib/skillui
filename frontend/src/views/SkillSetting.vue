<template>
  <div class="flex flex-col space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
        <Wrench :size="24" class="text-emerald-500"/>
        {{ $t('toolSettings.title') }}
      </h1>
    </div>

    <SkillSettingLocation v-model="defaultSkillDir" @changed="onSkillDirChanged"/>

    <SkillSettingIDE :defaultSkillDir="defaultSkillDir"/>
  </div>
</template>

<script setup lang="ts">
import { Wrench } from 'lucide-vue-next';
import { onMounted, ref } from 'vue';
import { GetSkillDir } from '../../wailsjs/go/main/App';
import SkillSettingIDE from './SkillSetting/IDE.vue';
import SkillSettingLocation from './SkillSetting/Location.vue';

const defaultSkillDir = ref('');

onMounted(async () => {
  try {
    defaultSkillDir.value = await GetSkillDir();
  } catch (e) {
    console.error('Failed to get skill dir', e);
  }
});

const onSkillDirChanged = (newDir: string) => {
  defaultSkillDir.value = newDir;
};
</script>
