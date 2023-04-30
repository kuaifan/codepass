<template>
    <div class="workspaces">
        <Header @domainSave="domainSave"/>
        <Banner/>

        <!-- 搜索 -->
        <div class="search">
            <div class="wrapper" :class="{loading: loadIng}">
                <div class="input-box">
                    <n-input round v-model:value="searchKey" placeholder="">
                        <template #prefix>
                            <n-icon :component="SearchOutline"/>
                        </template>
                    </n-input>
                    <div class="reload" @click="onLoad(true)">
                        <Loading v-if="loadIng"/>
                        <n-icon v-else>
                            <reload/>
                        </n-icon>
                    </div>
                </div>
                <div class="interval"></div>
                <n-button type="success" :render-icon="addIcon" @click="createModal = true">
                    创建工作区
                </n-button>
            </div>
            <n-modal v-model:show="createModal" :auto-focus="false">
                <n-card
                        style="width:600px;max-width:90%"
                        title="创建工作区"
                        :bordered="false"
                        size="huge"
                        closable
                        @close="createModal=false">
                    <Create @createDone="createDone"/>
                </n-card>
            </n-modal>
        </div>

        <!-- 列表 -->
        <div class="list">
            <div class="wrapper">
                <n-empty v-if="searchList.length === 0" class="empty" size="huge" description="没有工作区"/>
                <ul v-else>
                    <li class="nav">
                        <div class="name">工作区名称</div>
                        <div class="release">系统版本</div>
                        <div class="state">状态</div>
                        <div class="menu">操作</div>
                    </li>
                    <li v-for="item in searchList">
                        <div class="name">{{ item.name }}</div>
                        <div class="release">{{ item.release }}</div>
                        <div class="state">{{ state(item) }}</div>
                        <n-dropdown
                                trigger="click"
                                :show-arrow="true"
                                :options="operationMenu"
                                :render-label="operationLabelRender"
                                @updateShow="operationShow($event, item)"
                                @select="operationSelect($event, item)">
                            <n-icon class="menu" size="20">
                                <ellipsis-vertical/>
                            </n-icon>
                        </n-dropdown>
                    </li>
                </ul>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import {defineComponent, computed, h, ref, VNodeChild} from "vue";
import Header from "../components/Header.vue";
import Banner from "../components/Banner.vue";
import Create from "../components/Create.vue";
import Loading from "../components/Loading.vue";
import {AddOutline, EllipsisVertical, Reload, SearchOutline} from "@vicons/ionicons5";
import {useMessage, useDialog} from "naive-ui";
import call from "../call.js";

export default defineComponent({
    components: {
        EllipsisVertical,
        Reload, Loading, Create,
        Banner,
        Header,
    },
    computed: {
        SearchOutline() {
            return SearchOutline
        }
    },
    setup() {
        const message = useMessage()
        const dialog = useDialog()
        const createModal = ref(false);
        const loadIng = ref(false);
        const items = ref([])
        const searchKey = ref("");
        const searchList = computed(() => {
            if (searchKey.value === "") {
                return items.value
            }
            return items.value.filter(item => item.name.indexOf(searchKey.value) !== -1)
        })
        const operationMenu = ref([
            {
                label: '打开',
                key: 'open',
                disabled: false,
            }, {
                type: 'divider',
                key: 'd1'
            }, {
                label: '删除',
                key: "delete",
            }
        ])
        const operationItem = ref({})

        const onLoad = (tip) => {
            if (loadIng.value) {
                return
            }
            loadIng.value = true
            //
            call({
                method: "get",
                url: 'workspaces/list',
            }).then(({data}) => {
                items.value = data.list
            }).catch(err => {
                if (tip === true) {
                    message.error(err.msg)
                } else {
                    console.log(err)
                }
            }).finally(() => {
                loadIng.value = false
            })
        }
        onLoad(false)
        setInterval(_ => onLoad(false), 1000 * 30)

        return {
            createModal,
            loadIng,
            searchKey,
            searchList,
            operationMenu,
            operationLabelRender(option) {
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
                if (option.key === 'open' && operationItem.value.url) {
                    return h(
                        'a',
                        {
                            href: operationItem.value.url,
                            target: '_blank'
                        },
                        {
                            default: () => option.label as VNodeChild
                        }
                    )
                }
                return option.label as VNodeChild
            },
            operationShow(show: boolean, item) {
                if (show) {
                    operationItem.value = item
                    operationMenu.value[0]['disabled'] = !(item.create === "Success" && item.state === "Running" && /^https*:\/\//.test(item.url))
                }
            },
            operationSelect(key: string | number, item) {
                if (key === 'delete') {
                    const d = dialog.warning({
                        title: '删除工作区',
                        content: '确定要删除该工作区吗？',
                        positiveText: '确定',
                        negativeText: '取消',
                        onPositiveClick: () => {
                            d.loading = true
                            return new Promise((resolve) => {
                                call({
                                    method: "get",
                                    url: 'workspaces/delete',
                                    data: {
                                        name: item.name
                                    }
                                }).then(({msg}) => {
                                    message.success(msg)
                                    items.value = items.value.filter(i => i.name !== item.name)
                                }).catch(({msg}) => {
                                    message.error(msg)
                                    onLoad(false)
                                }).finally(resolve)
                            })
                        }
                    })
                }
            },
            addIcon() {
                return h(AddOutline);
            },
            createDone() {
                createModal.value = false
                onLoad(true)
            },
            domainSave() {
                onLoad(true)
            },
            onLoad
        };
    },
    methods: {
        state(item) {
            if (item.create === 'Success') {
                return item.state || 'Unknown'
            } else {
                return item.create || 'Error'
            }
        }
    }
})
</script>

<style lang="less" scoped>
.search {
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;

    .wrapper {
        flex: 1;
        padding: 20px 0;
        border-bottom: 1px solid rgba(41, 37, 36, 0.8);
        display: flex;
        align-items: center;
        flex-direction: row;
        justify-content: space-between;

        &.loading,
        &:hover {
            .input-box {
                .reload {
                    > i,
                    .loading {
                        opacity: 1;
                    }
                }
            }
        }

        .input-box {
            display: flex;
            align-items: center;

            .reload {
                margin-left: 16px;
                width: 30px;
                height: 30px;
                display: flex;
                align-items: center;
                justify-items: center;

                > i,
                .loading {
                    transition: all 0.3s;
                    opacity: 0.5;
                    font-size: 20px;
                    width: 20px;
                    height: 20px;
                }
            }
        }

        .interval {
            flex: 1;
        }
    }

}

.list {
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;

    .wrapper {
        flex: 1;
        padding: 18px 0;

        .empty {
            margin: 120px 0;
        }

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
