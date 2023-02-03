import {createDiscreteApi} from 'naive-ui'
import router from '../router'
import fly from 'flyio'

fly.config.headers.Accept = 'application/json'
fly.config.timeout = 10000
fly.config.withCredentials = true
fly.config.baseURL = '/api'

const naive = createDiscreteApi(
    ['message', 'dialog', 'loadingBar'],
)

fly.interceptors.request.use(request => {
    naive.loadingBar.start()
    return request
})

fly.interceptors.response.use(
    (res) => {
        naive.loadingBar.finish()
        return res.data
    },
    (err: any) => {
        naive.loadingBar.error()
        let status = err.status,
            message = err.response?.data
        console.log(status, message)
        switch (err.status) {
            case 401:
                naive.message.warning('登录失效 正在返回首页')
                router.replace({
                    name: 'Index',
                })
                break

            case 422:
                naive.message.error(message)
                break

            case 500:
                naive.dialog.error({
                    title: '貌似服务器错误了',
                    positiveText: '我会报告BUG的',
                })
                break

            default:
                naive.message.error(`好像哪里出现了错误 - ${status}`)
                break
        }
        return err
    },
)

export default fly
