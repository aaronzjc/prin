<template>
<main class="container">
    <nav class="navbar" role="navigation" aria-label="main navigation">
        <div class="navbar-brand">
            <a class="navbar-item" href="/">
                <img alt="一些好用的小公举" src="../assets/logo.png" />
                <p>一些好用的小公举</p>
            </a>
        </div>
    </nav>
    <div class="columns">
        <div class="column">
            <div class="tabs">
                <ul>
                    <router-link v-for="(tab, idx) in state.tabs" v-bind="$attrs" :key="idx" :to="tab.name" custom v-slot="{ isActive, navigate }">
                    <li :class="isActive ? 'is-active' : '' " @click="navigate" role="link">
                        <a href="#">{{ tab.title }}</a>
                    </li>
                    <slot />
                    </router-link>
                </ul>
            </div>
        </div>
    </div>
    <router-view v-slot="{ Component }">
        <keep-alive>
            <component :is="Component" />
        </keep-alive>
    </router-view>
    <div class="columns" id="footer">
        <div class="column copyright has-text-centered">
            <p>
                <a href="https://github.com/aaronzjc">@aaronzjc</a>
                开发，源码<a href="https://github.com/aaronzjc/prin">在此</a>，欢迎Star 
                v{{ version }}. 
                <span v-if="state.count != ''">当前在线 <strong class="online">{{ state.count }}</strong> 人</span>
            </p>
        </div>
    </div>
</main>
</template>

<script>
import { reactive, readonly } from 'vue'
import { mainRoutes } from '../router/router'
import {Get} from "../tools/http"

export default {
    name: "Main",
    setup() {
        const state = reactive({
            tabs: readonly(mainRoutes),
            count: "",
        })
        const version = readonly(process.env.VUE_APP_VERSION)
        let ticker = false
        async function GetOnline() {
            let resp = await Get("/api/stat/online")
            if (resp.data.code === 10000) {
                state.count = resp.data.data.count
                ticker = true
            }
        }
        
        GetOnline() // 初始化
        if (ticker) {
            setInterval(GetOnline, 30 * 1000)
        }

        return {
            state,
            version
        }
    }
}
</script>

<style lang="scss">
    nav.navbar {
    .navbar-item {
        p {
            font-size: .8rem;
            color: #666;
            padding: 0 0 0 1rem;
        }
        img {
            max-height: 2rem;
        }
    }
    .navbar-brand {
        margin-left: 0 !important;
    }
}
.container:not(.is-max-desktop) {
    max-width: 960px;
}
#footer {
    margin-top: 100px;
    background: none;
}
</style>