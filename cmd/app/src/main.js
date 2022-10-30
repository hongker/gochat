import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import piniaPluginPersist from 'pinia-plugin-persist'
import {connectSocket} from "./utils/ws";

// import './assets/main.css'

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'


const app = createApp(App)

app.use(ElementPlus)
const store = createPinia()
store.use(piniaPluginPersist)
app.use(store)
app.use(router)
app.provide('socket', connectSocket)
app.provide('operation', {
    heartbeat:1,
    connect: 2,
    updateProfile: 3,
    profile: 4,
    listSession: 5,
    sendMessage: 6,
    createGroup: 7,
    joinGroup: 8,
    leaveGroup:9,
    broadcastGroup: 10,

    queryHistory: 11,
    queryContacts: 12,
    queryChannel: 13,
    newMessage: 101,
})


app.mount('#app')
