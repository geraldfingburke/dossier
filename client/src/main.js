import { createApp } from 'vue'
import App from './App.vue'
import { createStore } from './store'

const app = createApp(App)
const store = createStore()

app.config.globalProperties.$store = store
app.provide('store', store)

app.mount('#app')
