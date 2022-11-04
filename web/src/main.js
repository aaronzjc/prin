import "bulma"

import { createApp } from 'vue'
import App from './App.vue'

const app = createApp(App)

import router from "./router/router";
app.use(router)

// 注册全局组件，因为循环引用的关系，需要注册为全局的
import Rules from './components/ext/iptables/rules.vue'
import Chains from './components/ext/iptables/chains.vue'
app.component('rules', Rules)
app.component('chains', Chains)

app.mount('#app')
