<template>
  <div class="Index" :style="{
    background: style_background
  }">
    <div class="info">
      <div class="icon">
        <img src="/favicon.svg">
      </div>
      <div class="version">
        <span>NPS<small>0.26.10</small></span>
      </div>
      <p class="title">一款轻量级、高性能、功能强大的内网穿透代理的服务端</p>
      <div class="introduction">
        <p>基于 <code>NPS</code> 适配 <code>NPC <small>0.26.10</small></code> 之后版本</p>
        <p>微信公众号、小程序开发等</p>
        <p>本地项目公网展示</p>
        <p>各种支付回调调试</p>
        <p><code>TCP</code>、<code>UDP</code> 转发等</p>
      </div>
      <div class="buttons">
        <n-button class="button" type="primary" size="large" ghost round tag="a" href="https://github.com/puzzle9/nps" target="_blank">
          GITHUB
        </n-button>
        <n-button class="button" type="primary" size="large" round tag="a" :href="`${fly.config.baseURL}/sign/in`">
          开始使用
        </n-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import fly from '../utils/fly'
import {useRouter, useRoute} from 'vue-router'
import {ref, computed} from 'vue'

const router = useRouter(),
    route = useRoute()

// todo: 渐变色
const style_backgrounds = [
      'linear-gradient(to left bottom, hsl(246, 100%, 85%) 0%,hsl(45, 100%, 85%) 100%)',
      'linear-gradient(to left bottom, hsl(80, 100%, 85%) 0%, hsl(199, 100%, 85%) 100%)',
      'linear-gradient(to left bottom, hsl(89, 100%, 85%) 0%,hsl(7, 100%, 85%) 100%)',
      'linear-gradient(to left bottom, hsl(201, 100%, 85%) 0%,hsl(133, 100%, 85%) 100%)',
      'linear-gradient(to left bottom, hsl(74, 100%, 85%) 0%,hsl(73, 100%, 85%) 100%)',
    ],
    // https://stackoverflow.com/a/5915122/6686061
    style_background = style_backgrounds[Math.floor(Math.random() * style_backgrounds.length)]
</script>

<style lang="stylus" scoped>
.Index
  height 100%
  width 100%
  display flex

  .info
    margin auto

    .icon
      text-align center

    .version
      font-size: 2.5rem
      font-weight: 300
      position: relative
      text-align: center

      small
        bottom: -0.4375rem
        font-size: 1rem
        position: absolute

    .title
      font-size 1.5rem

    .introduction
      text-align center
      line-height: 1

    .buttons
      display: flex
      justify-content: center

      .button
        margin 10px
</style>
