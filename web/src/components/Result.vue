<template>
    <div class="result">
        <n-result
                :status="status"
                :title="title"
                :description="desc">
            <template #footer>
                <n-button @click="goHome">返回首页</n-button>
            </template>
        </n-result>
    </div>
</template>

<style scoped lang="less">
.result {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    padding: 0 24px;
    display: flex;
    align-items: center;
    justify-content: center;
}
</style>
<script lang="ts">
import {defineComponent, ref} from 'vue'
import utils from '../utils.js'

export default defineComponent({
    setup() {
        let resultCode = utils.resultCode()
        const resultMsg = utils.resultMsg()

        if (resultCode === 400) {
            resultCode = "info"
        }

        const status = ref<"500" | "error" | "info" | "success" | "warning" | "404" | "403" | "418">(utils.inArray(resultCode, ["500", "error", "info", "success", "warning", "404", "403", "418"]) ? resultCode : "info")
        const title = ref(resultMsg.length <= 10 ? resultMsg : '')
        const desc = ref(resultMsg.length > 10 ? resultMsg : '')
        const goHome = () => {
            window.location.href = "/"
        }

        return {
            resultCode,
            resultMsg,

            status,
            title,
            desc,

            goHome
        }
    }
})
</script>
