<template>
    <n-form
            ref="formRef"
            :model="formData"
            :rules="formRules"
            size="large"
            label-placement="left"
            label-width="auto"
            require-mark-placement="right-hanging">
        <n-form-item path="name" label="名称">
            <n-input v-model:value="formData.name" placeholder="请输入工作区名称"/>
        </n-form-item>
        <n-form-item v-show="advanced" path="cpus" label="CPU">
            <n-input v-model:value="formData.cpus" placeholder="请输入CPU核数">
                <template #suffix>
                    核
                </template>
            </n-input>
        </n-form-item>
        <n-form-item v-show="advanced" path="memory" label="内存">
            <n-input v-model:value="formData.memory" placeholder="请输入内存大小">
                <template #suffix>
                    GB
                </template>
            </n-input>
        </n-form-item>
        <n-form-item v-show="advanced" path="disk" label="磁盘" placeholder="请输入磁盘大小">
            <n-input v-model:value="formData.disk">
                <template #suffix>
                    GB
                </template>
            </n-input>
        </n-form-item>
        <n-row :gutter="[0, 24]">
            <n-col :span="24">
                <div class="button-group">
                    <n-button round quaternary type="default" @click="advanced=!advanced">
                        高级
                    </n-button>
                    <n-button round type="primary" @click="handleSubmit">
                        创建
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
    FormRules, useMessage
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
    emits: {
        onCreate: () => true,
    },
    setup(props, {emit}) {
        const message = useMessage()
        const advanced = ref<boolean>(false)
        const loadIng = ref<boolean>(false)
        const formRef = ref<FormInst | null>(null)
        const formData = ref<ModelType>({
            name: null,
            cpus: "2",
            memory: "2",
            disk: "20",
        })

        const formRules: FormRules = {
            name: [
                {
                    required: true,
                    message: '请输入名称',
                    trigger: ['input', 'blur']
                }
            ],
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
        return {
            advanced,
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
                    const data = utils.cloneJSON(formData.value)
                    data.disk = data.disk + "GB"
                    data.memory = data.memory + "GB"
                    call({
                        method: "get",
                        url: 'workspaces/create',
                        data
                    }).then(({msg}) => {
                        message.success(msg);
                        emit('createDone')
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
