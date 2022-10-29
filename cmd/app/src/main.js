import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import piniaPluginPersist from 'pinia-plugin-persist'

// import './assets/main.css'

// import ElementUI from 'element-plus'
// import 'element-plus/dist/index.css'

const app = createApp(App)

// app.use(ElementUI)
const store = createPinia()
store.use(piniaPluginPersist)
app.use(store)
app.use(router)
app.provide('operation', {
    connect: 2,
    updateProfile: 3,
    profile: 4,
})


app.mount('#app')
