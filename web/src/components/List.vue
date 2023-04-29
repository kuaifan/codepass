<template>
    <div class="list">
        <div class="wrapper">
            <ul>
                <li class="nav">
                    <div class="name">实例名称</div>
                    <div class="release">实例版本</div>
                    <div class="state">状态</div>
                    <div class="menu">操作</div>
                </li>
                <li v-for="item in data">
                    <div class="name">{{ item.name }}</div>
                    <div class="release">{{ item.release }}</div>
                    <div class="state">{{ item.state }}</div>
                    <n-dropdown
                            trigger="click"
                            :show-arrow="true"
                            :options="options"
                            :render-label="renderDropdownLabel"
                            @select="handleSelect($event, item)">
                        <n-icon class="menu" size="20">
                            <ellipsis-vertical/>
                        </n-icon>
                    </n-dropdown>
                </li>
            </ul>
        </div>
    </div>
</template>

<style lang="less" scoped>
.list {
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;

    .wrapper {
        flex: 1;
        padding: 18px 0;

        ul {
            margin: 0;
            padding: 0;

            > li {
                display: flex;
                align-items: center;
                list-style: none;
                white-space: nowrap;
                justify-content: space-between;
                padding: 22px;
                border-radius: 12px;

                &.nav {
                    font-size: 16px;
                    font-weight: 600;

                    &:hover {
                        background-color: transparent;
                    }

                    .menu {
                        cursor: default;
                        &:hover {
                            background-color: transparent;
                        }
                    }
                }

                &:hover {
                    background-color: rgb(41, 37, 36);
                }

                .name {
                    width: 35%;
                }

                .release {
                    width: 35%;
                }

                .state {
                    width: 20%;
                }

                .menu {
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    min-width: 32px;
                    height: 32px;
                    border-radius: 6px;
                    cursor: pointer;
                    &:hover {
                        background-color: rgb(68, 64, 60);
                    }
                }
            }
        }
    }

}
</style>
<script lang="ts">
import {defineComponent, h} from 'vue'
import type {VNodeChild} from 'vue'
import {EllipsisVertical} from '@vicons/ionicons5'
import {useMessage} from 'naive-ui'

type ItemData = {
    name: string
    release: string
    state: string
}

const data: ItemData[] = [
    {name: "aabb1", release: 'Ubuntu 20.04 LTS', state: 'Running'},
    {name: "aabb2", release: 'Ubuntu 20.04 LTS', state: 'Running'}
]

export default defineComponent({
    components: {
        EllipsisVertical
    },
    setup() {
        const message = useMessage()
        return {
            data,
            options: [
                {
                    label: '打开',
                    key: 'open',
                },
                {
                    type: 'divider',
                    key: 'd1'
                },
                {
                    label: '删除',
                    key: "delete",
                }
            ],
            renderDropdownLabel(option) {
                if (option.key === 'delete') {
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
                return h(
                    'a',
                    {
                        href: '',
                        target: '_blank'
                    },
                    {
                        default: () => option.label as VNodeChild
                    }
                )
            },
            handleSelect(key: string | number, item: ItemData) {
                if (key === 'delete') {
                    // console.log(item);
                    message.info(String(key))
                }
            },
        }
    }
})
</script>
