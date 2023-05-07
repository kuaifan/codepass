<template>
    <n-config-provider :theme-name="themeName" :theme="theme">
        <n-loading-bar-provider>
            <n-message-provider>
                <n-notification-provider>
                    <n-dialog-provider>
                        <Login v-if="resultCode === 401"/>
                        <Result v-else-if="resultCode > 0"/>
                        <router-view v-else/>
                        <n-global-style/>
                    </n-dialog-provider>
                </n-notification-provider>
            </n-message-provider>
        </n-loading-bar-provider>
    </n-config-provider>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {siteSetup} from './store'

import Workspaces from "./pages/Workspaces.vue";
import Login from "./components/Login.vue";
import Result from "./components/Result.vue";

export default defineComponent({
    components: {
        Result,
        Login,
        Workspaces,
    },
    setup() {
        return siteSetup()
    }
});
</script>
