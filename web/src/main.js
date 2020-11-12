import "bulma"

import { createApp } from 'vue'
const app = createApp(App)

import router from "./router/router";

app.use(router)

import App from './App.vue'

app.mount('#app')
