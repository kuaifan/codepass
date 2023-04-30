<template>
    <div class="create-log">
        <n-log ref="nRef" :log="content" trim/>
        <div class="footer">
            <n-button :loading="loading" @click="getData">刷新</n-button>
        </div>
    </div>
</template>

<style lang="less" scoped>
.create-log {
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
import {useMessage} from "naive-ui";

export default defineComponent({
    props: {
        name: {
            type: String,
            required: true
        }
    },
    setup(props) {
        const message = useMessage()

        const content = ref("");
        const loading = ref(false);
        const nRef = ref(null);

        const getData = () => {
            if (loading.value) {
                return
            }
            loading.value = true
            call({
                method: "get",
                url: 'workspaces/create/log',
                data: {
                    name: props.name
                }
            }).then(({data}) => {
                content.value = data.log
                nextTick(() => {
                    nRef.value?.scrollTo({ position: 'bottom', slient: true })
                })
            }).catch(err => {
                message.error(err.msg)
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
