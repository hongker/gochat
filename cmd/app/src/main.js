import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import piniaPluginPersist from 'pinia-plugin-persist'
import {connectSocket} from "./utils/ws";

// import './assets/main.css'

// import ElementUI from 'element-plus'
// import 'element-plus/dist/index.css'


const app = createApp(App)

// app.use(ElementUI)
const store = createPinia()
store.use(piniaPluginPersist)
app.use(store)
app.use(router)
app.provide('socket', connectSocket)
app.provide('operation', {
    connect: 2,
    updateProfile: 3,
    profile: 4,
    sendMessage: 6,
    queryContacts: 12,
    newMessage: 101,
})


app.mount('#app')
