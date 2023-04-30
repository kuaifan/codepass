<template>
    <div class="info">
        <n-log ref="nRef" :log="content" :rows="10" trim/>
        <div class="footer">
            <n-button :loading="loading" @click="getData">刷新</n-button>
        </div>
    </div>
</template>

<style lang="less" scoped>
.info {
    .footer {
        display: flex;
        align-items: center;
        justify-content: center;
        margin-top: 26px;
    }
}
</style>
<script lang="ts">
import {defineComponent, ref} from 'vue'
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
                url: 'workspaces/info',
                data: {
                    name: props.name,
                    format: "text"
                }
            }).then(({data}) => {
                content.value = data.info
            }).catch(err => {
                message.error(err.msg)
            }).finally(() => {
                loading.value = false
            })
        }
        getData()

        return {
            content,
            loading,
            nRef,

            getData
        }
    }
})
</script>
