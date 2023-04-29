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
        <n-modal v-model:show="showModal" :auto-focus="false">
            <n-card
                    style="width:600px;max-width:90%"
                    title="域名证书"
                    :bordered="false"
                    size="huge"
                    closable
                    @close="showModal=false">
                <Domain @close="showModal=false"/>
            </n-card>
        </n-modal>
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
import Domain from "./Domain.vue";

const showModal = ref(false)

export default defineComponent({
    components: {Domain, EllipsisVertical},
    setup() {
        const message = useMessage()
        return {
            showModal,
            options: [
                {
                    label: '域名证书',
                    key: 'domain',
                },
                {
                    type: 'divider',
                    key: 'd1'
                },
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
                if (key === 'domain') {
                    showModal.value = true
                } else if (key === 'logout') {
                    message.warning('没有实现')
                }
            },
        }
    }
})
</script>
