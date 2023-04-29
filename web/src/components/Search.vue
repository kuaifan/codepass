<template>
    <div class="search">
        <div class="wrapper" :class="{loading: loadIng}">
            <div class="input-box">
                <n-input round placeholder="">
                    <template #prefix>
                        <n-icon :component="SearchOutline"/>
                    </template>
                </n-input>
                <div class="reload" @click="loadIng=!loadIng">
                    <Loading v-if="loadIng"/>
                    <n-icon v-else>
                        <reload/>
                    </n-icon>
                </div>
            </div>
            <div class="interval"></div>
            <n-button type="success" :render-icon="addIcon" @click="showModal = true">
                创建实例
            </n-button>
        </div>
        <n-modal v-model:show="showModal" :auto-focus="false">
            <n-card
                    style="width:600px;max-width:90%"
                    title="创建实例"
                    :bordered="false"
                    size="huge"
                    closable
                    @close="showModal=false">
                <Create/>
            </n-card>
        </n-modal>
    </div>
</template>

<style lang="less" scoped>
.search {
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;

    .wrapper {
        flex: 1;
        padding: 20px 0;
        border-bottom: 1px solid rgba(41, 37, 36, 0.8);
        display: flex;
        align-items: center;
        flex-direction: row;
        justify-content: space-between;

        &.loading,
        &:hover {
            .input-box {
                .reload {
                    > i,
                    .loading {
                        display: flex;
                    }
                }
            }
        }

        .input-box {
            display: flex;
            align-items: center;

            .reload {
                margin-left: 16px;
                width: 30px;
                height: 30px;
                display: flex;
                align-items: center;
                justify-items: center;

                > i,
                .loading {
                    display: none;
                    font-size: 20px;
                    width: 20px;
                    height: 20px;
                }
            }
        }

        .interval {
            flex: 1;
        }
    }

}
</style>
<script>
import {defineComponent, ref, h} from "vue";
import {SearchOutline, AddOutline, Reload} from "@vicons/ionicons5";
import Loading from "./Loading.vue";
import Create from "./Create.vue";

export default defineComponent({
    components: {
        Create,
        Loading,
        Reload
    },
    computed: {
        SearchOutline() {
            return SearchOutline
        }
    },
    setup() {
        return {
            loadIng: ref(false),
            showModal: ref(false),
            addIcon() {
                return h(AddOutline);
            }
        };
    }
})
</script>
