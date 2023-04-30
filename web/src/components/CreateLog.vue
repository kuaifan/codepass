<template>
    <div class="create-log">
        <n-log ref="logRef" :log="logInfo" trim/>
        <div class="footer">
            <n-button :loading="logLoad" @click="getData">刷新</n-button>
        </div>
    </div>
</template>

<style lang="less" scoped>
.create-log {
    .footer {
        display: flex;
        align-items: center;
        justify-content: center;
    }
}
</style>
<script lang="ts">
import {defineComponent, nextTick, ref} from 'vue'
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

        const logInfo = ref("");
        const logLoad = ref(false);
        const logRef = ref(null);

        const getData = () => {
            if (logLoad.value) {
                return
            }
            logLoad.value = true
            call({
                method: "get",
                url: 'workspaces/create/log',
                data: {
                    name: props.name
                }
            }).then(({data}) => {
                logInfo.value = data.log
                nextTick(() => {
                    logRef.value?.scrollTo({ position: 'bottom', slient: true })
                })
            }).catch(err => {
                message.error(err.msg)
            }).finally(() => {
                logLoad.value = false
            })
        }
        getData()
        setInterval(getData, 30 * 1000)

        return {
            logInfo,
            logLoad,
            logRef,

            getData
        }
    }
})
</script>
