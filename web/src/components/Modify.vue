<template>
    <n-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        :disabled="readIng"
        size="large"
        label-placement="left"
        label-width="auto"
        require-mark-placement="right-hanging">
        <n-form-item path="name" label="名称">
            <n-input v-model:value="formData.name" disabled></n-input>
        </n-form-item>
        <n-form-item path="cpus" label="CPU">
            <n-input v-model:value="formData.cpus" placeholder="请输入CPU核数">
                <template #suffix>
                    核
                </template>
            </n-input>
        </n-form-item>
        <n-form-item path="memory" label="内存">
            <n-input v-model:value="formData.memory" placeholder="请输入内存大小">
                <template #suffix>
                    GB
                </template>
            </n-input>
        </n-form-item>
        <n-form-item path="disk" label="磁盘">
            <n-input v-model:value="formData.disk" placeholder="请输入磁盘大小">
                <template #suffix>
                    GB
                </template>
            </n-input>
        </n-form-item>
        <n-row :gutter="[0, 24]">
            <n-col :span="24">
                <div class="button-group">
                    <n-button :loading="loadIng" round type="primary" @click="handleSubmit">
                        修改
                    </n-button>
                </div>
            </n-col>
        </n-row>
    </n-form>
</template>

<script lang="ts">
import {defineComponent, ref} from 'vue'
import {
    FormInst,
    FormItemRule,
    FormRules, useMessage, useDialog
} from 'naive-ui'
import call from "../call.js";
import utils from "../utils.js";

interface ModelType {
    name: string | null
    cpus: string | null
    memory: string | null
    disk: string | null
}

export default defineComponent({
    props: {
        name: {
            type: String,
            required: true
        },
        show: {
            type: Boolean,
        }
    },
    setup(props, {emit}) {
        const message = useMessage()
        const dialog = useDialog()
        const readIng = ref<boolean>(true)
        const loadIng = ref<boolean>(false)
        const formRef = ref<FormInst | null>(null)
        const formData = ref<ModelType>({
            name: props.name,
            cpus: "...",
            memory: "...",
            disk: "...",
        })

        call({
            method: "get",
            url: 'workspaces/info',
            data: {
                name: props.name,
                format: "hard"
            }
        }).then(({data}) => {
            formData.value.cpus = data.cpus
            formData.value.memory = String(utils.parseInt(data.memory))
            formData.value.disk = String(utils.parseInt(data.disk))
        }).catch(({msg}) => {
            dialog.error({
                title: '请求错误',
                content: msg,
                positiveText: '确定',
            })
        }).finally(() => {
            readIng.value = false
        })

        const formRules: FormRules = {
            cpus: [
                {
                    validator(rule: FormItemRule, value: string) {
                        if (value) {
                            if (!/^\d+$/.test(value)) {
                                return new Error('CPU必须是整数')
                            } else if (Number(value) < 1 || Number(value) > 8) {
                                return new Error('CPU应该是1-8之间的整数')
                            }
                        }
                        return true
                    },
                    trigger: ['input', 'blur']
                }
            ],
            disk: [
                {
                    validator(rule: FormItemRule, value: string) {
                        if (value) {
                            if (!/^\d+$/.test(value)) {
                                return new Error('硬盘必须是整数')
                            } else if (Number(value) < 10 || Number(value) > 1024) {
                                return new Error('硬盘应该是10-1024之间的整数')
                            }
                        }
                        return true
                    },
                    trigger: ['input', 'blur']
                }
            ],
            memory: [
                {
                    validator(rule: FormItemRule, value: string) {
                        if (value) {
                            if (!/^[1-9][0-9]*(\.[0-9]{1,2})?$/.test(value)) {
                                return new Error('内存格式输入错误')
                            } else if (Number(value) < 0.5 || Number(value) > 64) {
                                return new Error('内存应该是0.5-64之间的整数')
                            }
                        }
                        return true
                    },
                    trigger: ['input', 'blur']
                }
            ],
        }

        const handleSubmit = (e: MouseEvent) => {
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
                const data = utils.cloneJSON(formData.value)
                data.disk = data.disk + "GB"
                data.memory = data.memory + "GB"
                call({
                    method: "post",
                    url: 'workspaces/modify',
                    data
                }).then(({msg}) => {
                    message.success(msg);
                    emit('update:show', false)
                }).catch(({msg}) => {
                    dialog.error({
                        title: '请求错误',
                        content: msg,
                        positiveText: '确定',
                    })
                }).finally(() => {
                    loadIng.value = false
                })
            })
        }
        return {
            readIng,
            loadIng,
            formRef,
            formData,
            formRules,
            handleSubmit
        }
    }
})
</script>
