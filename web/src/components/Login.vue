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
                {{label(item.label)}}
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
    width: 100%;
    height: 100%;
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
import Cookies from "js-cookie";

export default defineComponent({
    components: {
        LogoGithub,
        AddCircleOutline
    },
    methods: {
        label(value) {
            if (!value) return ''
            return value.replace(/\+/g, ' ')
        }
    },
    setup() {
        const items = ref(utils.jsonParse(Cookies.get('result_msg')))
        if (!utils.isArray(items.value)) {
            items.value = []
        }
        return {
            items,
        }
    }
})
</script>
