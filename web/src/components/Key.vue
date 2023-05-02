<template>
    <n-form
            ref="formRef"
            :model="formData"
            :rules="formRules"
            size="large">
        <n-form-item path="key" label="Key">
            <n-input
                    v-model:value="formData.key"
                    type="textarea"
                    :autosize="{ minRows: 5, maxRows: 8 }"
                    placeholder="开始以 'ssh-rsa', 'ecdsa-sha2-nistp256', 'ecdsa-sha2-nistp384', 'ecdsa-sha2-nistp521', 'ssh-ed25519', 'sk-ecdsa-sha2-nistp256@openssh.com', or 'sk-ssh-ed25519@openssh.com'"/>
        </n-form-item>
        <n-form-item path="title" label="标题">
            <n-input v-model:value="formData.title" placeholder="例如：笔记本电脑"/>
        </n-form-item>
        <n-row :gutter="[0, 24]">
            <n-col :span="24">
                <div style="display:flex; justify-content:flex-end">
                    <n-button round type="primary" :loading="loadIng" @click="handleSubmit">
                        保存
                    </n-button>
                </div>
            </n-col>
        </n-row>
    </n-form>
</template>

<script lang="ts">
import {defineComponent, ref} from 'vue'
import {FormInst, FormRules, useMessage} from 'naive-ui'
import call from "../call.js";

interface ModelType {
    title: string | null
    key: string | null
}

export default defineComponent({
    emits: {
        onCreate: () => true,
    },
    setup(props, {emit}) {
        const message = useMessage()
        const loadIng = ref<boolean>(true)
        const formRef = ref<FormInst | null>(null)
        const formData = ref<ModelType>({
            title: null,
            key: null,
        })
        const formRules: FormRules = {
            key: [
                {
                    required: true,
                    message: '请输入KEY',
                    trigger: ['input', 'blur']
                }
            ],
        }

        call({
            method: "get",
            url: 'keys/info',
        }).then(({data}) => {
            formData.value = data
        }).catch(err => {
            console.log(err);
        }).finally(() => {
            loadIng.value = false
        })

        return {
            loadIng,
            formRef,
            formData,
            formRules,
            handleSubmit(e: MouseEvent) {
                e.preventDefault()
                formRef.value?.validate((errors) => {
                    if (errors) {
                        return;
                    }
                    //
                    if (loadIng.value) {
                        return
                    }
                    loadIng.value = true
                    call({
                        method: "post",
                        url: 'keys/save',
                        data: formData.value
                    }).then(({msg}) => {
                        message.success(msg);
                        emit('keySave')
                    }).catch(({msg}) => {
                        message.error(msg);
                    }).finally(() => {
                        loadIng.value = false
                    })
                })
            }
        }
    }
})
</script>
