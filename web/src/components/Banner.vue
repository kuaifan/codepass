<template>
    <div class="banner">
        <div class="wrapper">
            <div class="title">工作区</div>
            <div class="sub-title">管理最近和停止的工作区</div>
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
            font-size: 42px;
            line-height: 64px;
        }

        .sub-title {
            font-size: 16px;
            line-height: 28px;
        }
    }

}
</style>
<script>
import {defineComponent} from "vue";

export default defineComponent({})
</script>
