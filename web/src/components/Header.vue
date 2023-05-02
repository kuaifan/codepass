<template>
    <n-layout-header class="header" bordered>
        <div class="wrapper">
            <n-button tertiary round>
                工作区管理
            </n-button>
            <div class="user">
                <n-button
                        size="small"
                        quaternary
                        class="name"
                        @click="handleThemeUpdate">
                    {{ themeLabelMap[themeName] }}
                </n-button>
                <n-dropdown
                        trigger="click"
                        :show-arrow="true"
                        :options="userMenuOptions"
                        :render-label="renderDropdownLabel"
                        @select="handleMenuSelect">
                    <n-avatar class="avatar" round>C</n-avatar>
                </n-dropdown>
            </div>
        </div>
    </n-layout-header>
</template>

<style lang="less" scoped>
.header {
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    height: 64px;

    .wrapper {
        flex: 1;
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        align-items: center;
    }

    .user {
        display: flex;
        align-items: center;

        .name {
            margin-right: 18px;
        }

        .avatar {
            cursor: pointer;
        }
    }
}
</style>

<script lang="ts">
import {defineComponent, computed, h, ref, VNodeChild} from "vue";
import {useMessage} from 'naive-ui'
import {EllipsisVertical} from "@vicons/ionicons5";
import {useThemeName} from '../store'

export default defineComponent({
    components: {EllipsisVertical},
    setup() {
        const message = useMessage()
        const themeLabelMap = computed(() => ({
            dark: "浅色",
            light: "深色"
        }))
        const themeName = useThemeName()
        const handleThemeUpdate = () => {
            if (themeName.value === 'dark') {
                themeName.value = 'light'
            } else {
                themeName.value = 'dark'
            }
        }
        const userMenuOptions = ref([{
            label: '退出登录',
            key: "logout",
        }])
        const renderDropdownLabel = (option) => {
            if (option.key === 'logout') {
                return h(
                    'span',
                    {
                        style: 'color:rgb(248,113,113);height:34px;display:flex;align-items:center',
                    },
                    {
                        default: () => option.label as VNodeChild
                    }
                )
            }
            return option.label as VNodeChild
        }
        const handleMenuSelect = (key: string) => {
            if (key === 'logout') {
                message.warning('没有实现')
            }
        }
        return {
            // theme
            themeName,
            themeLabelMap,
            handleThemeUpdate,
            //
            userMenuOptions,
            handleMenuSelect,
            renderDropdownLabel,
        }
    }
})
</script>
