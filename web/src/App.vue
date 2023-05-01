<template>
    <n-config-provider :theme-name="themeName" :theme="theme">
        <n-loading-bar-provider>
            <n-message-provider>
                <n-notification-provider>
                    <n-dialog-provider>
                        <Workspaces/>
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

export default defineComponent({
    components: {
        Workspaces,
    },
    setup() {
        return siteSetup()
    }
});
</script>
