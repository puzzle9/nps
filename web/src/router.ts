import {createRouter, createWebHashHistory} from 'vue-router'

const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: '/',
            name: 'Index',
            component: () => import('./views/index.vue'),
        },
        {
            path: '/home',
            name: 'Home',
            component: () => import('./views/home.vue'),
        },
    ],
})

export default router
