<template>
    <div class="login">
        <div class="buttons">
            <n-button
                    v-for="item in items"
                    tag="a"
                    :href="item.href"
                    size="large"
                    type="success">
                <template #icon>
                    <n-icon v-if="item.type === 'github'">
                        <logo-github />
                    </n-icon>
                </template>
                {{item.label}}
            </n-button>
        </div>
        <div class="policy">登录即表示您同意我们的服务条款和隐私政策。</div>
    </div>
</template>

<style lang="less" scoped>
.login {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    .buttons {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-direction: column;
        > a + a {
            margin-top: 12px;
        }
    }
    .policy {
        padding: 12px 32px 32px;
    }
}
</style>
<script>
import {defineComponent, ref} from "vue";
import {LogoGithub, AddCircleOutline} from "@vicons/ionicons5";
import utils from "../utils.js";

export default defineComponent({
    components: {
        LogoGithub,
        AddCircleOutline
    },
    setup() {
        const items = ref(utils.jsonParse(utils.resultMsg()))
        if (!utils.isArray(items.value)) {
            items.value = []
        }
        return {
            items,
        }
    }
})
</script>
