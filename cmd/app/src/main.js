import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

import './assets/main.css'
import socket from './utils/ws'

const app = createApp(App)


app.use(createPinia())
app.use(router)
app.provide('socket', socket)


app.mount('#app')
