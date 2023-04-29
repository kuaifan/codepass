import axios from 'axios'
import utils from "./utils.js";


// 创建一个 axios 实例
const call = axios.create({
    baseURL: '/api', // 所有的请求地址前缀部分
    timeout: 60000, // 请求超时时间毫秒
    withCredentials: true, // 异步请求携带cookie
    headers: {
        // 设置后端需要的传参类型
        'Content-Type': 'application/x-www-form-urlencoded',
    },
})

// 添加请求拦截器
call.interceptors.request.use(
    function (config) {
        // 在发送请求之前做些什么
        if (config.method === 'get' && config.data) {
            config.url = utils.urlAddParams(config.url, config.data)
        }
        return config
    },
    function (error) {
        // 对请求错误做些什么
        console.log(error)
        return Promise.reject(error)
    }
)

// 添加响应拦截器
call.interceptors.response.use(
    function (response) {
        // console.log(response)
        // 2xx 范围内的状态码都会触发该函数。
        // 对响应数据做点什么
        // dataAxios 是 axios 返回数据中的 data
        const dataAxios = response.data
        // 这个状态码是和后端约定的
        const code = dataAxios.reset
        //
        if (!utils.isJson(dataAxios)) {
            return Promise.reject({ret: 0, msg: "返回数据格式错误"})
        }
        if (dataAxios.ret !== 1) {
            return Promise.reject(dataAxios)
        }
        return dataAxios
    },
    function (error) {
        // 超出 2xx 范围的状态码都会触发该函数。
        // 对响应错误做点什么
        console.log(error)
        return Promise.reject({ret: 0, msg: "请求失败", "data": error})
    }
)

export default call
