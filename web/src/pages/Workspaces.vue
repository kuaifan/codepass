<template>
    <div class="workspaces">
        <Header/>
        <Banner/>
        <n-divider/>

        <!-- 搜索 -->
        <div class="search">
            <div class="wrapper" :class="{loading: loadIng}">
                <div class="search-box">
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
                <n-divider/>
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
                <template v-else>
                    <div class="item nav">
                        <div class="name">工作区名称</div>
                        <div class="release">系统版本</div>
                        <div class="state">状态</div>
                        <div class="menu">操作</div>
                    </div>
                    <n-list hoverable :show-divider="false">
                        <n-list-item v-for="item in searchList">
                            <div class="item">
                                <div class="name">{{ item.name }}</div>
                                <div class="release">{{ item.release || '-' }}</div>
                                <div class="state" @click="onState(item)">
                                    <div v-if="stateLoading(item)" class="load">
                                        <Loading/>
                                    </div>
                                    <div class="text">{{ stateText(item) }}</div>
                                </div>
                                <n-dropdown
                                        trigger="click"
                                        :show-arrow="true"
                                        :options="operationMenu"
                                        :render-label="operationLabel"
                                        @updateShow="operationShow($event, item)"
                                        @select="operationSelect($event, item)">
                                    <n-button quaternary class="menu">
                                        <template #icon>
                                            <n-icon>
                                                <ellipsis-vertical/>
                                            </n-icon>
                                        </template>
                                    </n-button>
                                </n-dropdown>
                            </div>
                        </n-list-item>
                    </n-list>
                </template>
            </div>
            <n-modal v-model:show="infoModal" :auto-focus="false">
                <n-card
                        style="width:600px;max-width:90%"
                        title="详情"
                        :bordered="false"
                        size="huge"
                        closable
                        @close="infoModal=false">
                    <Info :name="infoName"/>
                </n-card>
            </n-modal>
            <n-modal v-model:show="logModal" :auto-focus="false">
                <n-card
                        style="width:600px;max-width:90%"
                        title="日志"
                        :bordered="false"
                        size="huge"
                        closable
                        @close="logModal=false">
                    <Log :name="logName"/>
                </n-card>
            </n-modal>
        </div>
    </div>
</template>

<script lang="ts">
import {defineComponent, computed, h, ref, VNodeChild, onBeforeUnmount} from "vue";
import Header from "../components/Header.vue";
import Banner from "../components/Banner.vue";
import Create from "../components/Create.vue";
import Loading from "../components/Loading.vue";
import Log from "../components/Log.vue";
import Info from "../components/Info.vue";
import {AddOutline, EllipsisVertical, Reload, SearchOutline} from "@vicons/ionicons5";
import {useMessage, useDialog, NButton} from "naive-ui";
import call from "../call.js";

export default defineComponent({
    components: {
        Info,
        Log,
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
        const infoModal = ref(false);
        const infoName = ref("");
        const logModal = ref(false);
        const logName = ref("");
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
                label: '详情',
                key: 'info',
            }, {
                label: '日志',
                key: 'log',
            }, {
                type: 'divider',
                key: 'd1'
            }, {
                label: '删除',
                key: "delete",
            }
        ])
        const operationItem = ref({url: '', pass: '', repos_name: ''})

        const operationLabel = (option) => {
            if (option.disabled === true) {
                return option.label as VNodeChild
            }
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
                const action = operationItem.value.url + "/login?folder=/workspace/" + operationItem.value.repos_name
                return h('form',
                    {
                        action,
                        method: 'post',
                        target: '_blank'
                    },
                    [
                        h('input', {type: 'hidden', name: 'base', value: '.'}),
                        h('input', {type: 'hidden', name: 'href', value: action}),
                        h('input', {type: 'hidden', name: 'password', value: operationItem.value.pass}),
                        h(NButton, {
                            text: true,
                            attrType: 'submit'
                        }, {
                            default: () => option.label as VNodeChild
                        }),
                    ])
            }
            return option.label as VNodeChild
        }
        const operationShow = (show: boolean, item) => {
            if (show) {
                operationItem.value = item
                operationMenu.value[0]['disabled'] = !(item.create === "Success" && item.state === "Running" && /^https*:\/\//.test(item.url))
            }
        }
        const operationSelect = (key: string | number, item) => {
            if (key === 'info') {
                infoName.value = item.name
                infoModal.value = true
            } else if (key === 'log') {
                logName.value = item.name
                logModal.value = true
            } else if (key === 'delete') {
                const dd = dialog.warning({
                    title: '删除工作区',
                    content: '确定要删除该工作区吗？',
                    positiveText: '确定',
                    negativeText: '取消',
                    onPositiveClick: () => {
                        dd.loading = true
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
        }
        const addIcon = () => {
            return h(AddOutline);
        }
        const createDone = () => {
            createModal.value = false
            onLoad(true)
        }
        const stateLoading = (item) => {
            return ['Success', ''].indexOf(item.create) !== -1
                && ['Failed', ''].indexOf(item.state) !== -1
        }
        const stateText = (item) => {
            if (item.create === 'Success') {
                return item.state || 'Unknown'
            } else {
                return item.create || 'Error'
            }
        }
        const onState = (item) => {
            if (stateLoading(item)) {
                operationSelect('log', item)
            }
        }
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
        const loadInter = setInterval(_ => onLoad(false), 1000 * 30)
        onBeforeUnmount(() => {
            clearInterval(loadInter)
        })

        return {
            createModal,
            infoModal,
            infoName,
            logModal,
            logName,
            loadIng,
            searchKey,
            searchList,
            operationMenu,
            operationLabel,
            operationShow,
            operationSelect,
            addIcon,
            createDone,
            stateLoading,
            stateText,
            onState,
            onLoad
        };
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

        .search-box {
            display: flex;
            align-items: center;
            flex-direction: row;
            justify-content: space-between;
        }

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
                margin: 0 32px 0 16px;
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

        .empty {
            margin: 120px 0;
        }

        > ul {
            background-color: transparent;

            > li {
                border-radius: 18px;
            }
        }

        .item {
            display: flex;
            align-items: center;
            list-style: none;
            white-space: nowrap;
            justify-content: space-between;
            padding: 12px 0;

            &.nav {
                font-size: 16px;
                font-weight: 600;
                margin-bottom: 8px;
                padding-left: 20px;
                padding-right: 20px;
            }

            .name {
                width: 35%;
            }

            .release {
                width: 35%;
            }

            .state {
                width: 20%;
                display: flex;
                align-items: center;

                .load {
                    flex-shrink: 0;
                    width: 16px;
                    height: 16px;
                    margin-right: 6px;
                }

                .text {
                    flex: 1;
                }
            }

            .menu {
                min-width: 32px;
                padding: 0;
            }
        }

    }
}
</style>
