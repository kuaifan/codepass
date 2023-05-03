export const routes = [
    {
        name: 'workspaces',
        path: '/workspaces',
        component: () => import('../pages/Workspaces.vue')
    },
    {
        name: 'not-found',
        path: '/:pathMatch(.*)*',
        redirect: {
            name: 'workspaces',
        }
    }
]
