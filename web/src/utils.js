import localforage from "localforage";

localforage.config({name: 'web', storeName: 'common'});

const utils = {
    /**
     * 是否数组
     * @param obj
     * @returns {boolean}
     */
    isArray(obj) {
        return typeof (obj) == "object" && Object.prototype.toString.call(obj).toLowerCase() == '[object array]' && typeof obj.length == "number";
    },

    /**
     * 是否数组对象
     * @param obj
     * @returns {boolean}
     */
    isJson(obj) {
        return typeof (obj) == "object" && Object.prototype.toString.call(obj).toLowerCase() == "[object object]" && typeof obj.length == "undefined";
    },

    /**
     * 克隆对象
     * @param myObj
     * @returns {*}
     */
    cloneJSON(myObj) {
        if (typeof (myObj) !== 'object') return myObj;
        if (myObj === null) return myObj;
        //
        return utils.jsonParse(utils.jsonStringify(myObj))
    },

    /**
     * 将一个 JSON 字符串转换为对象（已try）
     * @param str
     * @param defaultVal
     * @returns {*}
     */
    jsonParse(str, defaultVal = undefined) {
        if (str === null) {
            return defaultVal ? defaultVal : {};
        }
        if (typeof str === "object") {
            return str;
        }
        try {
            return JSON.parse(str.replace(/\n/g, "\\n").replace(/\r/g, "\\r"));
        } catch (e) {
            return defaultVal ? defaultVal : {};
        }
    },

    /**
     * 将 JavaScript 值转换为 JSON 字符串（已try）
     * @param json
     * @param defaultVal
     * @returns {string}
     */
    jsonStringify(json, defaultVal = undefined) {
        if (typeof json !== 'object') {
            return json;
        }
        try {
            return JSON.stringify(json);
        } catch (e) {
            return defaultVal ? defaultVal : "";
        }
    },

    /**
     * 字符串是否包含
     * @param string
     * @param find
     * @param lower
     * @returns {boolean}
     */
    strExists(string, find, lower = false) {
        string += "";
        find += "";
        if (lower !== true) {
            string = string.toLowerCase();
            find = find.toLowerCase();
        }
        return (string.indexOf(find) !== -1);
    },

    /**
     * 字符串是否左边包含
     * @param string
     * @param find
     * @param lower
     * @returns {boolean}
     */
    leftExists(string, find, lower = false) {
        string += "";
        find += "";
        if (lower !== true) {
            string = string.toLowerCase();
            find = find.toLowerCase();
        }
        return (string.substring(0, find.length) === find);
    },

    /**
     * 删除左边字符串
     * @param string
     * @param find
     * @param lower
     * @returns {string}
     */
    leftDelete(string, find, lower = false) {
        string += "";
        find += "";
        if (utils.leftExists(string, find, lower)) {
            string = string.substring(find.length)
        }
        return string ? string : '';
    },

    /**
     * 字符串是否右边包含
     * @param string
     * @param find
     * @param lower
     * @returns {boolean}
     */
    rightExists(string, find, lower = false) {
        string += "";
        find += "";
        if (lower !== true) {
            string = string.toLowerCase();
            find = find.toLowerCase();
        }
        return (string.substring(string.length - find.length) === find);
    },

    /**
     * 删除右边字符串
     * @param string
     * @param find
     * @param lower
     * @returns {string}
     */
    rightDelete(string, find, lower = false) {
        string += "";
        find += "";
        if (utils.rightExists(string, find, lower)) {
            string = string.substring(0, string.length - find.length)
        }
        return string ? string : '';
    },

    /**
     * 删除地址中的参数
     * @param url
     * @param parameter
     * @returns {string|*}
     */
    removeURLParameter(url, parameter) {
        if (parameter instanceof Array) {
            parameter.forEach((key) => {
                url = utils.removeURLParameter(url, key)
            });
            return url;
        }
        let urlparts = url.split('?');
        if (urlparts.length >= 2) {
            //参数名前缀
            let prefix = encodeURIComponent(parameter) + '=';
            let pars = urlparts[1].split(/[&;]/g);

            //循环查找匹配参数
            for (let i = pars.length; i-- > 0;) {
                if (pars[i].lastIndexOf(prefix, 0) !== -1) {
                    //存在则删除
                    pars.splice(i, 1);
                }
            }

            return urlparts[0] + (pars.length > 0 ? '?' + pars.join('&') : '');
        }
        return url;
    },

    /**
     * 连接加上参数
     * @param url
     * @param params
     * @returns {*}
     */
    urlAddParams(url, params) {
        if (utils.isJson(params)) {
            if (url) {
                url = utils.removeURLParameter(url, Object.keys(params))
            }
            url += "";
            url += url.indexOf("?") === -1 ? '?' : '';
            for (let key in params) {
                if (!params.hasOwnProperty(key)) {
                    continue;
                }
                url += '&' + key + '=' + params[key];
            }
        } else if (params) {
            url += (url.indexOf("?") === -1 ? '?' : '&') + params;
        }
        if (!url) {
            return ""
        }
        return utils.rightDelete(url.replace("?&", "?"), '?');
    },

    /**
     * =============================================================================
     * *****************************   localForage   ******************************
     * =============================================================================
     */
    __IDBTimer: {},

    IDBSave(key, value, delay = 100) {
        if (typeof utils.__IDBTimer[key] !== "undefined") {
            clearTimeout(utils.__IDBTimer[key])
            delete utils.__IDBTimer[key]
        }
        utils.__IDBTimer[key] = setTimeout(async _ => {
            await localforage.setItem(key, value)
        }, delay)
    },

    IDBDel(key) {
        localforage.removeItem(key).then(_ => {
        })
    },

    IDBSet(key, value) {
        return localforage.setItem(key, value)
    },

    IDBRemove(key) {
        return localforage.removeItem(key)
    },

    IDBClear() {
        return localforage.clear()
    },

    IDBValue(key) {
        return localforage.getItem(key)
    },

    async IDBString(key, def = "") {
        const value = await utils.IDBValue(key)
        return typeof value === "string" || typeof value === "number" ? value : def;
    },

    async IDBInt(key, def = 0) {
        const value = await utils.IDBValue(key)
        return typeof value === "number" ? value : def;
    },

    async IDBBoolean(key, def = false) {
        const value = await utils.IDBValue(key)
        return typeof value === "boolean" ? value : def;
    },

    async IDBArray(key, def = []) {
        const value = await utils.IDBValue(key)
        return utils.isArray(value) ? value : def;
    },

    async IDBJson(key, def = {}) {
        const value = await utils.IDBValue(key)
        return utils.isJson(value) ? value : def;
    }
}

export default utils
