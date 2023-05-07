<template>
    <div class="info">
        <n-log ref="nRef" :log="content" :rows="10" trim/>
        <div class="footer">
            <n-button :loading="loading" @click="getData">刷新</n-button>
            <n-button type="primary" @click="onModify">修改</n-button>
        </div>
        <n-modal v-model:show="modifyModal" :auto-focus="false">
            <n-card
                style="width:600px;max-width:90%"
                title="修改工作区"
                :bordered="false"
                size="huge"
                closable
                @close="modifyModal=false">
                <Modify :name="name" v-model:show="modifyModal"/>
            </n-card>
        </n-modal>
    </div>
</template>

<style lang="less" scoped>
.info {
    .footer {
        display: flex;
        align-items: center;
        justify-content: center;
        margin-top: 26px;
        > * {
            margin: 0 8px;
        }
    }
}
</style>
<script lang="ts">
import {defineComponent, ref} from 'vue'
import call from "../call.js";
import {useDialog} from "naive-ui";
import Modify from "./Modify.vue";

export default defineComponent({
    components: {Modify},
    props: {
        name: {
            type: String,
            required: true
        },
        show: {
            type: Boolean,
        },
    },
    setup(props, {emit}) {
        const dialog = useDialog()

        const content = ref("");
        const loading = ref(false);
        const nRef = ref(null);
        const modifyModal = ref(false);

        const getData = () => {
            if (loading.value) {
                return
            }
            loading.value = true
            call({
                method: "get",
                url: 'workspaces/info',
                data: {
                    name: props.name,
                    format: "text"
                }
            }).then(({data}) => {
                content.value = data.info
            }).catch(({msg}) => {
                dialog.error({
                    title: '请求错误',
                    content: msg,
                    positiveText: '确定',
                    onPositiveClick: () => {
                        emit("update:show", false)
                    }
                })
            }).finally(() => {
                loading.value = false
            })
        }
        getData()
        const onModify = () => {
            modifyModal.value = true
        }

        return {
            content,
            loading,
            nRef,
            modifyModal,

            getData,
            onModify,
        }
    }
})
</script>
