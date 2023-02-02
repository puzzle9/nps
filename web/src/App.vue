<template>
  <n-config-provider :class="app_class" :locale="zhCN" :date-locale="dateZhCN" :theme="theme">
    <n-global-style/>
    <router-view></router-view>
  </n-config-provider>
</template>

<script setup lang="ts">
import {zhCN, dateZhCN, useOsTheme, lightTheme, darkTheme} from 'naive-ui'
import {useRouter, useRoute} from 'vue-router'
import {ref, computed} from 'vue'

const router = useRouter(),
    route = useRoute()

const app_class = computed(() => {
  return route.name == 'Index' ? 'app_display_flex' : null
})

const theme_os = useOsTheme(),
    theme_default = ref(null),
    theme = computed(() => {
      let theme_name = theme_default.value || theme_os.value
      switch (theme_name) {
        case 'dark':
          return darkTheme
        default:
          return lightTheme
      }
    })

// theme_default.value = 'dark'
</script>

<style lang="stylus">
body, html, #app
  height 100%
</style>

<style lang="stylus" scoped>
.app_display_flex
  height 100%
</style>
