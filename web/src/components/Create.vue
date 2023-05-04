<template>
    <n-form
            ref="formRef"
            :model="formData"
            :rules="formRules"
            size="large"
            label-placement="left"
            label-width="auto"
            require-mark-placement="right-hanging">
        <n-form-item path="repos" label="储存库">
            <n-select
                    v-model:value="formData.repos"
                    :options="reposComputed"
                    :show-arrow="false"
                    :on-blur="reposBlur"
                    filterable
                    tag
                    placeholder="请输入或选择储存库地址"/>
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
        <n-form-item v-show="advanced" path="disk" label="磁盘">
            <n-input v-model:value="formData.disk" placeholder="请输入磁盘大小">
                <template #suffix>
                    GB
                </template>
            </n-input>
        </n-form-item>
        <n-form-item v-show="advanced" path="image" label="系统">
            <n-select v-model:value="formData.image" :options="imageList" placeholder="请选择系统版本"/>
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
import {computed, defineComponent, ref} from 'vue'
import {
    FormInst,
    FormItemRule,
    FormRules, useMessage
} from 'naive-ui'
import call from "../call.js";
import utils from "../utils.js";

interface ModelType {
    repos: string | null
    cpus: string | null
    memory: string | null
    disk: string | null
    image: string | null
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
            repos: null,
            cpus: "2",
            memory: "2",
            disk: "20",
            image: "20.04",
        })

        const reposRef = ref(null)
        utils.IDBArray("userRepos").then((data) => {
            reposRef.value = data
            call({
                method: "get",
                url: 'user/repos',
            }).then(({data}) => {
                utils.IDBSave("userRepos", reposRef.value = data.list)
            })
        })
        const reposComputed = computed(() => {
            if (reposRef.value == null) {
                return []
            }
            return reposRef.value.map(item => {
                return {
                    label: item['html_url'],
                    value: item['html_url']
                }
            })
        })
        const imageList = ["18.04", "20.04", "22.04", "22.10"].map(item => {
            return {
                label: `Ubuntu ${item}`,
                value: item
            }
        })

        const formRules: FormRules = {
            repos: [
                {
                    required: true,
                    message: '请选择储存库地址',
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

        const reposBlur = (e) => {
            const value = `${e.target.value}`.trim()
            if (reposRef.value === null) {
                reposRef.value = []
            }
            if (/^https*:\/\//.test(value) && !reposRef.value.find(item => item['html_url'] === value)) {
                reposRef.value.unshift({
                    html_url: value
                })
                formData.value.repos = value
            }
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
        return {
            advanced,
            loadIng,
            formRef,
            formData,
            formRules,
            reposComputed,
            imageList,
            reposBlur,
            handleSubmit
        }
    }
})
</script>
