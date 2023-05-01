import {createRouter, createWebHistory} from 'vue-router'

export default function createDemoRouter(app, routes) {
    const router = createRouter({
        history: createWebHistory(),
        routes
    })
    router.beforeEach(function (to, from, next) {
        next()
    })

    router.afterEach(function (to, from) {
        //
    })

    return router
}
