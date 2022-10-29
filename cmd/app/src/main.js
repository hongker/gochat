import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

// import './assets/main.css'
import socket from './utils/ws'

// import ElementUI from 'element-plus'
// import 'element-plus/dist/index.css'

const app = createApp(App)

// app.use(ElementUI)
app.use(createPinia())
app.use(router)
app.provide('socket', socket)
app.provide('packet', {
    rawHeaderLen: 10,
    packetOffset: 0,
    opOffset: 4,
    contentTypeOffset: 6,
    seqOffset: 8,
})
app.provide('operation', {
    login: 2,
})


app.mount('#app')
