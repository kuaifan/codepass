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
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
}
</style>
<script lang="ts">
import {defineComponent, ref} from 'vue'
import Cookies from 'js-cookie'
import utils from '../utils.js'

export default defineComponent({
    setup() {
        let resultCode = Cookies.get('result_code')
        let resultMsg = `${Cookies.get('result_msg') || ''}`
        if (utils.parseInt(resultCode) === 400) {
            resultCode = "warning"
        }

        const status = ref<"500" | "error" | "info" | "success" | "warning" | "404" | "403" | "418">(utils.inArray(resultCode, ["500", "error", "info", "success", "warning", "404", "403", "418"]) ? resultCode : "info")
        const title = ref(resultMsg.length <= 10 ? resultMsg : '')
        const desc = ref(resultMsg.length > 10 ? resultMsg : '')
        const goHome = () => {
            Cookies.remove('result_code')
            Cookies.remove('result_msg')
            window.location.href = ""
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
