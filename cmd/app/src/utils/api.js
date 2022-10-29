import axios from 'axios'

const instance = axios.create({
    baseURL:"http://localhost:8080",//该部分写的是我们使用接口的公共代理
    timeout:2500//请求时间，在超时前，所有请求都会等待 2.5 秒，由于网络原因，我们可以把时间修改的长一点
});
instance.interceptors.request.use(function (config) {
    // 在发送请求之前做些什么
    //如果我们的使用的接口需要配置headers请求头或者body请求，可以再改部分添加
    //headers请求头:config.headers["字段名"]="字段值" + token值
    config.headers["Content-Type"] = "application/json"
    return config;
}, function (error) {
    // 对请求错误做些什么
    return Promise.reject(error);
});

//**该部分的instance是与上面的实例化的名字是一致的，实例化名字是什么，这里就是什么(必需修改)
instance.interceptors.response.use(function (response) {
    // 对响应数据做点什么
    //对数据进行处理,如：脱壳
    console.log(response.data)
    return response.data;
}, function (error) {
    // 对响应错误做点什么
    return Promise.reject(error);
});

export const get = (uri) => {
    return new Promise((resolve, reject) => {
        instance.get(uri).then(res=>{
            resolve(res)
        }).catch(err=>{
            reject(err)
        })
    })
}

export const sendRequest = (uri, param) => {
    return new Promise((resolve, reject) => {
        instance.post(uri, param).then(res=>{
            resolve(res)
        }).catch(err=>{
            reject(err)
        })
    })
}