<template>
    <n-config-provider :theme="darkTheme">
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
import {darkTheme} from 'naive-ui'

import Workspaces from "./pages/Workspaces.vue";

export default defineComponent({
    components: {
        Workspaces,
    },
    setup() {
        return {
            darkTheme,
        }
    }
});
</script>
