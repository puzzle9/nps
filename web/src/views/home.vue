<template>
  <div class="Home">
    <n-card title="NPS" hoverable>
      <template #header-extra>
        <a :href="base.user_account_url" target="_blank" v-if="base.user_avatar">
          <n-avatar round size="small" :src="base.user_avatar"/>
        </a>
        <n-button style="margin-left: 10px" text tag="a" :href="`${fly.config.baseURL}/sign/out`">退出登录</n-button>
      </template>
    </n-card>

    <br/>
    <n-card class="info" hoverable>
      <n-spin :show="loading">
        <n-tabs type="line" v-model:value="tab_name" animated>
          <n-tab-pane name="info" tab="参数信息">
            <n-list>
              <n-list-item>
                <n-thing title="限速" :title-extra="base.nps_rate_limit ? `${base.nps_rate_limit} KB/S` : '不限速'"/>
              </n-list-item>
              <n-list-item>
                <n-thing title="已用流量" :title-extra="`↑ ${(base.nps_flow_inlet / 1024 / 1024).toFixed(2)} MB / ↓ ${(base.nps_flow_export / 1024 / 1024).toFixed(2)} MB`"/>
              </n-list-item>
              <n-list-item>
                <n-thing title="连接命令" :title-extra="base.nps_is_connect ? '在线' : '离线'">
                  <template #description v-show="base.nps_bridge_domain">
                    <code class="connect">./npc -server={{ base.nps_bridge_domain }}:{{ base.nps_bridge_port }}
                      -vkey={{ base.nps_vkey }} -type={{ base.nps_bridge_type }}</code>
                  </template>
                </n-thing>
              </n-list-item>
            </n-list>
          </n-tab-pane>
          <!--          <n-tab-pane name="domain" tab="域名列表"></n-tab-pane>-->
          <n-tab-pane name="tunnel" tab="隧道列表">
            <n-form ref="tunnel" :model="tunnel_model">
              <n-dynamic-input v-model:value="tunnel_model" :on-create="tunnel_create" :max="8" #="{ index, value }">
                <n-form-item :show-label="false" :path="`[${index}].type`" :rule="tunnel_rules.type">
                  <n-select v-model:value="tunnel_model[index].type" filterable placeholder="隧道类型" :options="tunnel_types" style="width: 135px"/>
                </n-form-item>

                <n-form-item :show-label="false">
                  <n-input-number v-model:value="tunnel_model[index].port" placeholder="服务端端口" clearable :min="port_min" :max="port_max" :show-button="false" style="width: 120px"/>
                </n-form-item>

                <n-form-item :show-label="false" :path="`[${index}].target`" :rule="tunnel_rules.target">
                  <n-input v-model:value="tunnel_model[index].target" placeholder="目标地址 如 :8080" type="textarea" rows="1" clearable/>
                </n-form-item>

                <n-form-item :show-label="false">
                  <n-input v-model:value="tunnel_model[index].remark" placeholder="备注" type="textarea" maxlength="30" rows="1" show-count clearable/>
                </n-form-item>
              </n-dynamic-input>
            </n-form>
            <br/>
            <n-button type="info" attr-type="button" :loading="tunnel_loading" @click="tunnelSubmit" block>保存
            </n-button>
          </n-tab-pane>
        </n-tabs>
      </n-spin>
    </n-card>

    <br/>
    <!-- prettier-ignore-->
    <n-card title="使用方法" class="use" hoverable>
      <div v-show="tab_name == 'info'">
        <p>下载客户端软件:</p>
        <ul>已有 NPS 可以略过此部分</ul>
        <ul v-for="(urls, system_type) in downloads">
          <b>{{ system_type }}:</b>
          <n-button v-for="(download_url, download_name) in urls" tertiary type="info" size="small" tag="a" :href="download_url" target="_blank">
            {{ download_name }} 下载
          </n-button>
        </ul>

        <p>解压后在当前目录打开终端 并 输入连接命令 进行连接 如</p>
        <ul><img src="/npc_linux.png" alt="linux connect"/></ul>
      </div>

      <div v-show="tab_name == 'tunnel'">
        <p>隧道类型:</p>
        <ul v-for="info in tunnel_types" :key="info.value"><b>{{ info.label }}</b>: {{ info.placeholder }}</ul>
        <p>服务端端口:</p>
        <ul>可用区间为: {{ base.nps_allow_ports }}, 留空默认生成</ul>
        <p>目标地址:</p>
        <ul>代理到本地可以只填写端口号,在 TCP 隧道下可以填写多行支持负载均衡, 如 127.0.0.1:8080</ul>
        <p>示例:</p>
        <ul>
          <small>隧道类型: TCP/UDP 服务端端口: 9080 目标地址: :8080 -> 访问服务器的 9080 端口相当于以 TCP/UDP 模式访问
            本地的 8080 端口</small>
        </ul>
        <ul>
          <small>隧道类型: HTTP/SOCKET5 服务端端口: 9090 目标地址: :9099 -> 将设备代理 IP 设为服务端 IP 端口设为 9090
            相当于走了 HTTP/SOCKET5 模式的客户端代理</small>
        </ul>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import {createDiscreteApi, FormInst, FormRules} from 'naive-ui'
import {useRouter, useRoute} from 'vue-router'
import {ref, computed} from 'vue'
import fly from '../utils/fly'

const {message} = createDiscreteApi(['message'])

const router = useRouter(),
    route = useRoute()

const tab_name = ref('info')

const loading = ref(true),
    base = ref({
      nps_flow_inlet: 0,
      nps_flow_export: 0,
    }),
    port_max = ref(0),
    port_min = ref(0),
    getBase = () => {
      fly.post('/base').then((res: any) => {
        message.info(`欢迎回来 ${res.user_nickname}`)
        base.value = res
        loading.value = false

        // "8000-9999,30000-60000"
        let allow_ports: any = (res.nps_allow_ports as string).replaceAll('-', ',').split(',')
        // https://www.cnblogs.com/zhouyangla/p/8482010.html
        port_min.value = Math.min.apply(null, allow_ports)
        port_max.value = Math.max.apply(null, allow_ports)

        tunnelGet()
      })
    }

getBase()

const tunnel = ref<FormInst | null>(null),
    tunnel_model = ref([{}]),
    tunnel_rules: FormRules = {
      type: {
        required: true,
        message: ' ',
        trigger: ['input', 'change'],
      },
      target: {
        required: true,
        message: ' ',
        trigger: ['input', 'blur'],
      },
    },
    tunnel_types = [
      {
        label: 'TCP 隧道',
        value: 'tcp',
        placeholder: 'TCP 协议,如 HTTP(S)、SSH,接口调试通常选这个.',
      },
      {
        label: 'UDP 隧道',
        value: 'udp',
        placeholder: 'UTP 协议,如 DNS,通常用不到.',
      },
      {
        label: 'HTTP 代理',
        value: 'http',
        placeholder: '让服务器作为 HTTP 代理,访问内网资源等.',
      },
      {
        label: 'SOCKS5 代理',
        value: 'socks5',
        placeholder: '让服务器作为 SOCKS5 代理,访问内网资源等.',
      },
    ],
    tunnel_create = () => {
      return {
        type: null,
        value: null,
      }
    },
    tunnel_loading = ref(false),
    tunnelGet = () => {
      tunnel_loading.value = true
      fly
          .get('/tunnel')
          .then((res: any) => {
            tunnel_model.value = res.rows
          })
          .finally(() => {
            tunnel_loading.value = false
          })
    },
    tunnelSubmit = async () => {
      await tunnel.value?.validate((errors) => (errors ? message.error('貌似哪里不对') : null))

      tunnel_loading.value = true
      fly
          .post('/tunnel', tunnel_model.value)
          .then((res: any) => {
            message.success(res.message)
            tunnelGet()
          })
          .finally(() => {
            tunnel_loading.value = false
          })
    }

const downloads = {
  win64: {
    github: 'https://github.com/ehang-io/nps/releases/download/v0.26.10/windows_amd64_client.tar.gz',
  },
  linux64: {
    github: 'https://github.com/ehang-io/nps/releases/download/v0.26.10/linux_amd64_client.tar.gz',
  },
}
</script>

<style lang="stylus" scoped>
.info, .use
  width 98%
  max-width 800px
  margin auto

  .connect
    background-color: #F9F2F4;
    border-radius: 4px;
    color: #ca4440;
    font-size: 90%;
    padding: 2px 4px;

  img
    width 100%
</style>
