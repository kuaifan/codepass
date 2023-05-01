<template>
    <n-layout-header class="header" bordered>
        <div class="wrapper">
            <n-button tertiary round>
                工作区管理
            </n-button>
            <div class="user">
                <n-button
                        size="small"
                        tag="a"
                        quaternary
                        class="name"
                        target="_blank">
                    Codepass
                </n-button>
                <n-dropdown
                        trigger="click"
                        :show-arrow="true"
                        :options="options"
                        :render-label="renderDropdownLabel"
                        @select="handleSelect">
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
            margin-right: 12px;
        }

        .avatar {
            cursor: pointer;
        }
    }
}
</style>

<script lang="ts">
import {defineComponent, h, ref, VNodeChild} from "vue";
import { useMessage } from 'naive-ui'
import {EllipsisVertical} from "@vicons/ionicons5";

export default defineComponent({
    components: {EllipsisVertical},
    setup(props, {emit}) {
        const message = useMessage()
        return {
            options: [
                {
                    label: '退出登录',
                    key: "logout",
                }
            ],
            renderDropdownLabel(option) {
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
            },
            handleSelect(key: string) {
                if (key === 'logout') {
                    message.warning('没有实现')
                }
            }
        }
    }
})
</script>
