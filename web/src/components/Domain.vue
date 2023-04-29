<template>
    <n-form
            ref="formRef"
            :model="formData"
            :rules="formRules"
            size="large">
        <n-form-item path="domain" label="域名">
            <n-input v-model:value="formData.domain" placeholder="请输入域名"/>
        </n-form-item>
        <n-form-item path="key" label="私钥（key）">
            <n-input
                    v-model:value="formData.key"
                    type="textarea"
                    :autosize="{ minRows: 2, maxRows: 5 }"
                    placeholder="请输入域名私钥"/>
        </n-form-item>
        <n-form-item path="crt" label="证书（crt）">
            <n-input
                    v-model:value="formData.crt"
                    type="textarea"
                    :autosize="{ minRows: 2, maxRows: 5 }"
                    placeholder="请输入域名证书"/>
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
    domain: string | null
    key: string | null
    crt: string | null
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
            domain: null,
            key: null,
            crt: null,
        })
        const formRules: FormRules = {
            domain: [
                {
                    required: true,
                    message: '请输入名称',
                    trigger: ['input', 'blur']
                }
            ],
            key: [
                {
                    required: true,
                    message: '请输入名称',
                    trigger: ['input', 'blur']
                }
            ],
            crt: [
                {
                    required: true,
                    message: '请输入名称',
                    trigger: ['input', 'blur']
                }
            ],
        }

        call({
            method: "get",
            url: 'cert/info',
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
                        url: 'cert/save',
                        data: formData.value
                    }).then(({msg}) => {
                        message.success(msg);
                        emit('close')
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
