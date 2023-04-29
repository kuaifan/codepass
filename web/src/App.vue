<template>
    <n-loading-bar-provider>
        <n-message-provider>
            <n-notification-provider>
                <n-dialog-provider>
                    <n-config-provider :theme="darkTheme">
                        <Header/>
                        <Banner/>
                        <Search/>
                        <List/>
                        <n-global-style/>
                    </n-config-provider>
                </n-dialog-provider>
            </n-notification-provider>
        </n-message-provider>
    </n-loading-bar-provider>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {darkTheme} from 'naive-ui'

import Header from "./components/Header.vue";
import Banner from "./components/Banner.vue";
import Search from "./components/Search.vue";
import List from "./components/List.vue";

export default defineComponent({
    components: {
        Header,
        Banner,
        Search,
        List,
    },
    setup() {
        return {
            darkTheme,
        }
    }
});
</script>
