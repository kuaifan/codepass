import {computed, watch, ref} from 'vue'
import {darkTheme, useOsTheme} from 'naive-ui'
import utils from "../utils.js";

const themeNameRef = ref('light')
const themeRef = computed(() => {
    const {value} = themeNameRef
    return value === 'dark' ? darkTheme : null
})
watch(themeNameRef, name => {
    utils.IDBSave("themeName", name)
})

export function useThemeName() {
    return themeNameRef
}

export function siteSetup() {
    return {
        themeName: themeNameRef,
        theme: themeRef,
    }
}

export function init() {
    return new Promise(async (resolve) => {
        themeNameRef.value = await utils.IDBString("themeName")
        if (['light', 'dark'].indexOf(themeNameRef.value) === -1) {
            themeNameRef.value = useOsTheme().value
        }
        resolve()
    })
}
