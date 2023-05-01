export const routes = [
    {
        name: 'workspaces',
        path: '/workspaces',
        component: () => import('../pages/workspaces.vue')
    },
    {
        name: 'not-found',
        path: '/:pathMatch(.*)*',
        redirect: {
            name: 'workspaces',
            params: {
                theme: 'os-theme'
            }
        }
    }
]
