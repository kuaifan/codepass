<template>
    <div class="workspaces">
        <Header/>
        <Banner/>
        <n-divider/>

        <!-- 搜索 -->
        <div class="search">
            <div class="wrapper" :class="{loading: loadIng && loadShow}">
                <div class="search-box">
                    <div class="input-box">
                        <n-input round v-model:value="searchKey" placeholder="">
                            <template #prefix>
                                <n-icon>
                                    <SearchOutline/>
                                </n-icon>
                            </template>
                        </n-input>
                        <div class="reload" @click="onLoad(true, true)">
                            <Loading v-if="loadIng && loadShow"/>
                            <n-icon v-else>
                                <ReloadIcon/>
                            </n-icon>
                        </div>
                    </div>
                    <div class="interval"></div>
                    <n-button type="success" @click="createModal = true">
                        <template #icon>
                            <n-icon>
                                <AddOutline/>
                            </n-icon>
                        </template>
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
                    <Create @onDone="createDone"/>
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
                                <div class="name">
                                    <ul>
                                        <li>{{ item.name }}</li>
                                        <li class="repos">
                                            <n-button
                                                text
                                                tag="a"
                                                :href="item['repos_url']"
                                                target="_blank">
                                                {{ item['repos_url'] }}
                                            </n-button>
                                        </li>
                                    </ul>
                                </div>
                                <div class="release">{{ item.release || '-' }}</div>
                                <div class="state" @click="onState(item)">
                                    <div v-if="stateJudge(item, 'Running')" class="run">
                                        <n-badge type="success" show-zero dot />
                                    </div>
                                    <div v-else-if="stateJudge(item, 'Stopped')" class="run">
                                        <n-badge color="grey" show-zero dot />
                                    </div>
                                    <div v-else-if="stateLoading(item)" class="load">
                                        <Loading/>
                                    </div>
                                    <div class="text" :style="stateStyle(item)">{{ stateText(item) }}</div>
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
                                                <EllipsisVertical/>
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
                    <Info :name="infoName" v-model:show="logModal"/>
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
                    <Log :name="logName" v-model:show="logModal"/>
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
import {AddOutline, EllipsisVertical, Reload as ReloadIcon, SearchOutline} from "@vicons/ionicons5";
import {useMessage, useDialog, NButton} from "naive-ui";
import call from "../call.js";
import utils from "../utils.js";

export default defineComponent({
    components: {
        Info,
        Log,
        AddOutline,
        EllipsisVertical,
        ReloadIcon,
        SearchOutline,
        Loading,
        Create,
        Banner,
        Header,
    },
    setup() {
        const message = useMessage()
        const dialog = useDialog()
        const dLog = ref(null);
        const createModal = ref(false);
        const infoModal = ref(false);
        const infoName = ref("");
        const logModal = ref(false);
        const logName = ref("");
        const loadIng = ref(false);
        const loadShow = ref(false);
        const items = ref([])
        const searchKey = ref("");
        const searchList = computed(() => {
            if (searchKey.value === "") {
                return items.value
            }
            return items.value.filter(item => item.name.indexOf(searchKey.value) !== -1)
        })
        const setItemStatus = (name, status) => {
            items.value.forEach(item => {
                if (item.name === name) {
                    item.status = status
                }
            })
        }
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
                label: '启动',
                key: "start",
                disabled: false,
            }, {
                label: '停止',
                key: "stop",
                disabled: false,
            }, {
                label: '重启',
                key: "restart",
                disabled: false,
            }, {
                label: '删除',
                key: "delete",
            }
        ])
        const operationItem = ref({url: '', password: '', repos_name: ''})

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
                        h('input', {type: 'hidden', name: 'password', value: operationItem.value.password}),
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
                const state = stateText(item)
                operationSetDisabled('open', !/^https*:\/\//.test(item.url))
                operationSetDisabled('start', state !== 'Stopped')
                operationSetDisabled('stop', state !== "Running")
                operationSetDisabled('restart', state !== "Running")
            }
        }
        const operationSetDisabled = (key: string, disabled: boolean) => {
            operationMenu.value.forEach(item => {
                if (item.key === key) {
                    item['disabled'] = disabled
                }
            })
        }
        const operationSelect = (key: string | number, item) => {
            if (key === 'info') {
                infoName.value = item.name
                infoModal.value = true
            } else if (key === 'log') {
                logName.value = item.name
                logModal.value = true
            } else if (key === 'start') {
                const dd = dialog.warning({
                    title: '启动工作区',
                    content: '确定要启动该工作区吗？',
                    positiveText: '确定',
                    negativeText: '取消',
                    onPositiveClick: () => {
                        dd.loading = true
                        return operationInstance('start', item.name)
                    }
                })
            } else if (key === 'stop') {
                const dd = dialog.warning({
                    title: '停止工作区',
                    content: '确定要停止该工作区吗？',
                    positiveText: '确定',
                    negativeText: '取消',
                    onPositiveClick: () => {
                        dd.loading = true
                        return operationInstance('stop', item.name)
                    }
                })
            } else if (key === 'restart') {
                const dd = dialog.warning({
                    title: '重启工作区',
                    content: '确定要重启该工作区吗？',
                    positiveText: '确定',
                    negativeText: '取消',
                    onPositiveClick: () => {
                        dd.loading = true
                        return operationInstance('restart', item.name)
                    }
                })
            } else if (key === 'delete') {
                const dd = dialog.warning({
                    title: '删除工作区',
                    content: '确定要删除该工作区吗？',
                    positiveText: '确定',
                    negativeText: '取消',
                    onPositiveClick: () => {
                        dd.loading = true
                        return operationInstance('delete', item.name)
                    }
                })
            }
        }
        const operationInstance = (operation, name) => {
            return new Promise((resolve) => {
                call({
                    method: "get",
                    url: 'workspaces/operation',
                    data: {operation, name}
                }).then(({data, msg}) => {
                    message.success(msg)
                    setItemStatus(name, data.status)
                    onLoad(false, true)
                }).catch(({msg}) => {
                    dialog.error({
                        title: '请求错误',
                        content: msg,
                        positiveText: '确定',
                    })
                }).finally(resolve)
            })
        }
        const createDone = () => {
            createModal.value = false
            onLoad(true, true)
        }
        const stateJudge = (item, state) => {
            return stateText(item) === state
        }
        const stateLoading = (item) => {
            const state = stateText(item)
            return state.slice(-3) === 'ing' && state != "Running"
        }
        const stateStyle = (item) => {
            const state = stateText(item)
            switch (state) {
                case 'Success':
                    return {
                        color: 'rgb(82,196,26)'
                    }
                case 'Failed':
                    return {
                        color: 'rgb(248,113,113)'
                    }
                case 'Unknown':
                case 'Error':
                    return {
                        color: 'rgb(252,211,77)'
                    }
                default:
                    return {}
            }
        }
        const stateText = (item) => {
            if (item.status === 'Success') {
                return item.state || 'Unknown'
            } else {
                return item.status || 'Error'
            }
        }
        const onState = (item) => {
            if (stateLoading(item)) {
                operationSelect('log', item)
            }
        }
        const onLoad = (tip, showLoad) => {
            if (loadIng.value) {
                if (showLoad === true) {
                    loadShow.value = tip
                }
                return
            }
            loadIng.value = true
            loadShow.value = showLoad
            //
            call({
                method: "get",
                url: 'workspaces/list',
            }).then(({data}) => {
                if (!utils.isArray(data.list)) {
                    if (tip === true) {
                        message.warning("暂无数据")
                    }
                    items.value = []
                    return
                }
                items.value = data.list
            }).catch(({msg}) => {
                if (tip) {
                    if (dLog.value) {
                        dLog.value.destroy()
                        dLog.value = null
                    }
                    dLog.value = dialog.error({
                        title: '请求错误',
                        content: msg,
                        positiveText: '确定',
                    })
                }
            }).finally(() => {
                loadIng.value = false
            })
        }

        onLoad(false, true)
        const loadInter = setInterval(() => onLoad(false, false), 1000 * 30)

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
            loadShow,
            searchKey,
            searchList,
            operationMenu,
            operationLabel,
            operationShow,
            operationSelect,
            createDone,
            stateJudge,
            stateLoading,
            stateStyle,
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
                width: 40%;
                ul {
                    display: flex;
                    flex-direction: column;
                    justify-content: center;
                    padding: 0;
                    margin: 0;
                    > li {
                        list-style: none;
                        padding: 0 6px 0 0;
                        margin: 0;
                        font-size: 16px;
                        font-weight: 600;
                        overflow: hidden;
                        text-overflow: ellipsis;
                        white-space: nowrap;
                        &.repos {
                            font-size: 14px;
                            font-weight: normal;
                            > a {
                                opacity: 0.5;
                                user-select: auto;
                                &:hover {
                                    opacity: 1;
                                }
                            }
                        }
                    }
                }
            }

            .release {
                width: 30%;
            }

            .state {
                width: 20%;
                display: flex;
                align-items: center;

                .load,
                .run {
                    flex-shrink: 0;
                    width: 16px;
                    height: 16px;
                    margin-right: 6px;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                }
                .run {
                    .n-badge {
                        transform: scale(1.2);
                    }
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
