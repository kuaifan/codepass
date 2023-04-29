<template>
    <div class="banner">
        <div class="wrapper">
            <div class="title">Workspaces</div>
            <div class="sub-title">Manage recent and stopped workspaces.</div>
        </div>
    </div>
</template>

<style lang="less" scoped>
.banner {
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    border-bottom: 1px solid rgba(41, 37, 36, 0.8);

    .wrapper {
        flex: 1;
        padding: 24px 0 32px;

        .title {
            font-weight: 700;
            font-size: 48px;
            line-height: 64px;
        }

        .sub-title {
            font-size: 16px;
            line-height: 24px;
        }
    }

}
</style>
<script>
import {defineComponent} from "vue";

export default defineComponent({})
</script>
