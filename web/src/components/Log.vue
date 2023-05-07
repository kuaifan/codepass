<template>
    <div class="log">
        <n-log ref="nRef" :log="content" trim/>
        <div class="footer">
            <n-button :loading="loading" @click="getData">刷新</n-button>
        </div>
    </div>
</template>

<style lang="less" scoped>
.log {
    .footer {
        display: flex;
        align-items: center;
        justify-content: center;
        margin-top: 26px;
    }
}
</style>
<script lang="ts">
import {defineComponent, onBeforeUnmount, nextTick, ref} from 'vue'
import call from "../call.js";
import {useDialog} from "naive-ui";

export default defineComponent({
    props: {
        name: {
            type: String,
            required: true
        },
        show: {
            type: Boolean,
        },
    },
    setup(props, {emit}) {
        const dialog = useDialog()

        const nRef = ref(null);
        const dLog = ref(null);
        const loading = ref(false);
        const content = ref("");

        const scrollToBottom = () => {
            const {scrollbarRef} = nRef.value
            const { containerRef, contentRef } = scrollbarRef
            if (containerRef && contentRef) {
                const containerHeight = containerRef.offsetHeight
                const containerScrollTop = containerRef.scrollTop
                const contentHeight = contentRef.offsetHeight
                const scrollBottom = contentHeight - containerScrollTop - containerHeight
                return scrollBottom < 10
            }
            return true
        }

        const getData = () => {
            if (loading.value) {
                return
            }
            loading.value = true
            call({
                method: "get",
                url: 'workspaces/log',
                data: {
                    name: props.name
                }
            }).then(({data}) => {
                const isBottom = scrollToBottom()
                content.value = data.log
                isBottom && nextTick(() => {
                    nRef.value?.scrollTo({ position: 'bottom' })
                })
            }).catch(({msg}) => {
                if (dLog.value) {
                    dLog.value.destroy()
                    dLog.value = null
                }
                dLog.value = dialog.error({
                    title: '请求错误',
                    content: msg,
                    positiveText: '确定',
                    onPositiveClick: () => {
                        emit("update:show", false)
                    }
                })
            }).finally(() => {
                loading.value = false
            })
        }
        getData()
        const getInter = setInterval(getData, 15 * 1000)

        onBeforeUnmount(() => {
            clearInterval(getInter)
        })

        return {
            content,
            loading,
            nRef,

            getData
        }
    }
})
</script>
