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
                <n-thing>
                  <template #header>
                    <n-button class="n-thing-header__title" text @click="connectCopy">连接命令</n-button>
                  </template>
                  <template #header-extra>
                    <template v-if="base.nps_is_connect">在线</template>
                    <n-button v-else text @click="updateVKey">离线</n-button>
                  </template>
                  <template #description v-show="base.nps_bridge_domain">
                    <code class="connect">{{ connect }}</code>
                  </template>
                </n-thing>
              </n-list-item>
            </n-list>
          </n-tab-pane>
          <n-tab-pane name="tunnel" tab="隧道列表">
            <n-data-table :columns="tunnel_columns" :data="tunnel_data" :single-line="false" :scroll-x="950"/>
            <br/>
            <n-button type="info" @click="tunnel_model_show = true" block>新增隧道</n-button>
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
        <p>示例:</p>
        <ul>
          隧道类型: TCP/UDP 服务端端口: 9080 目标地址: :8080 -> 访问服务器的 9080 端口相当于以 TCP/UDP 模式访问
          本地的 8080 端口
        </ul>
        <ul>
          隧道类型: HTTP/SOCKET5 服务端端口: 9090 目标地址: :9099 -> 将设备代理 IP 设为服务端 IP 端口设为 9090
          相当于走了 HTTP/SOCKET5 模式的客户端代理
        </ul>
      </div>
    </n-card>

    <n-modal v-model:show="tunnel_model_show" preset="card" title="新增隧道" size="huge" :bordered="false" :mask-closable="false" class="tunnel_model" style="width: 90%; max-width: 800px">
      <n-form ref="tunnel_form" :model="tunnel_model" :rules="tunnel_rules">
        <n-form-item label="隧道类型" path="type">
          <n-select v-model:value="tunnel_model.type" filterable placeholder="隧道类型" :options="tunnel_types" :render-option="tunnelTypeRenderOption"/>
        </n-form-item>
        <n-form-item label="服务端端口" path="port">
          <n-input-number
              v-model:value="tunnel_model.port"
              :placeholder="`可用区间为: ${base.nps_allow_ports}, 留空默认生成`"
              clearable
              :min="port_min"
              :max="port_max"
              :show-button="false"
              style="width: 100%"
          />
        </n-form-item>
        <n-form-item label="目标地址" path="target">
          <n-input
              v-model:value="tunnel_model.target"
              placeholder="代理到本地可以只填写端口号,在 TCP 类型下可以填写多行支持负载均衡, 如
:1024
127.0.0.1:1314
局域网IP:5210
"
              :rows="4"
              type="textarea"
              clearable
          />
        </n-form-item>
        <n-form-item label="备注" path="remark">
          <n-input v-model:value="tunnel_model.remark" placeholder="给自己看的备注" type="textarea" maxlength="30" show-count clearable/>
        </n-form-item>
      </n-form>

      <template #action>
        <div class="action">
          <n-button type="info" quaternary @click="tunnel_model_show = false">取消</n-button>
          <n-button type="info" class="save" :loading="tunnel_loading" @click="tunnelSubmit">保存</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import {
  createDiscreteApi,
  NTooltip,
  NButton,
  NPopconfirm,
  FormInst,
  FormRules,
  DataTableColumns,
  SelectOption
} from 'naive-ui'
import {useRouter, useRoute} from 'vue-router'
import {ref, computed, h, VNode} from 'vue'
import useClipboard from 'vue-clipboard3'
import fly from '../utils/fly'

const {message} = createDiscreteApi(['message'])
const {toClipboard} = useClipboard()

const router = useRouter(),
    route = useRoute()

const tab_name = ref('info')

const loading = ref(true),
    base: any = ref({
      nps_vkey: null,
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

const connect = computed(
    // @prettier-ignore
    () =>
        `${/windows|win32/i.test(navigator.userAgent) ? 'npc.exe' : './npc'} -server=${base.value.nps_bridge_domain}:${base.value.nps_bridge_port} -vkey=${base.value.nps_vkey} -type=${
            base.value.nps_bridge_type
        }`,
)

const connectCopy = () => {
  toClipboard(connect.value)
  message.success('复制成功')
}

const updateVKey = () => {
  fly.post('/updateVKey').then((res: any) => {
    message.success('更新成功')
    base.value.nps_vkey = res.nps_vkey
  })
}

const tunnel_form = ref<FormInst | null>(null),
    tunnel_model_show = ref(false),
    tunnel_columns: DataTableColumns = [
      {
        title: '#',
        key: 'key',
        width: 30,
        render: (_, index) => `${index + 1}`,
        fixed: 'left',
      },
      {
        title: '隧道类型',
        key: 'type',
        width: 120,
        render: (row: any) => tunnel_types.find((tunnel) => tunnel.value == row.type)?.label,
      },
      {
        title: '访问地址',
        key: 'url',
        align: 'center',
        children: [
          {
            title: '域名',
            key: 'url_domain',
            width: 200,
            render: (row: any) =>
                h(
                    NButton,
                    {
                      text: true,
                      onClick: () => {
                        toClipboard(row.url_domain)
                        message.success('复制成功 暂只支持 http 或 https 访问')
                      },
                    },
                    {
                      default: () => row.url_domain,
                    },
                ),
          },
          {
            title: 'IP',
            key: 'url_ip',
            width: 180,
            render: (row: any) =>
                h(
                    NButton,
                    {
                      text: true,
                      onClick: () => {
                        toClipboard(row.url_ip)
                        message.success('复制成功')
                      },
                    },
                    {
                      default: () => row.url_ip,
                    },
                ),
          },
        ],
      },
      {
        title: '目标地址',
        key: 'target',
        resizable: true,
        minWidth: 150,
      },
      {
        title: '备注',
        key: 'remark',
        resizable: true,
        minWidth: 100,
        ellipsis: {
          tooltip: true,
        },
      },
      {
        title: '操作',
        key: 'actions',
        width: 70,
        fixed: 'right',
        render: (row: any, _index) =>
            h(
                NPopconfirm,
                {
                  showIcon: false,
                  trigger: 'click',
                  negativeText: '点错啦',
                  positiveText: '没点错',
                  onPositiveClick: async () => {
                    await fly.delete('/tunnel', {
                      id: row.id,
                    })
                    message.success('删除成功')
                    tunnelGet()
                  },
                },
                {
                  trigger: () =>
                      h(
                          NButton,
                          {
                            type: 'warning',
                            size: 'small',
                          },
                          {
                            default: () => '删除',
                          },
                      ),
                  default: () => `将删除 #${_index + 1} ${row.type} 隧道`,
                },
            ),
      },
    ],
    tunnel_model = ref({}),
    tunnel_data = ref([]),
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
        value: 'httpProxy',
        placeholder: '让服务器作为 HTTP 代理,访问内网资源等.',
      },
      {
        label: 'SOCKS5 代理',
        value: 'socks5',
        placeholder: '让服务器作为 SOCKS5 代理,访问内网资源等.',
      },
    ],
    tunnelTypeRenderOption = ({node, option}: { node: VNode; option: SelectOption }) =>
        h(NTooltip, null, {
          trigger: () => node,
          default: () => option.placeholder,
        }),
    tunnel_loading = ref(false),
    tunnelGet = () => {
      tunnel_loading.value = true
      fly
          .get('/tunnel')
          .then((res: any) => {
            tunnel_data.value = res.rows || []
          })
          .finally(() => {
            tunnel_loading.value = false
          })
    },
    tunnelSubmit = async () => {
      await tunnel_form.value?.validate((errors) => (errors ? message.error('貌似哪里不对') : null))

      tunnel_loading.value = true
      fly
          .post('/tunnel', tunnel_model.value)
          .then((res: any) => {
            message.success(res.message)
            tunnel_model_show.value = false
            tunnelGet()
            tunnel_model.value = {}
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
  max-width 1024px
  margin auto

  .connect
    background-color: #F9F2F4;
    border-radius: 4px;
    color: #ca4440;
    font-size: 90%;
    padding: 2px 4px;

  img
    width 100%

.tunnel_model
  .action
    float right

    .save
      margin-left 10px
</style>
