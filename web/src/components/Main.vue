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
            <p><a href="https://github.com/aaronzjc">@aaronzjc</a>开发，源码<a href="https://github.com/aaronzjc/prin">在此</a>，欢迎Star v1.6</p>
        </div>
    </div>
</main>
</template>

<script>
import { reactive, readonly } from 'vue'
import { mainRoutes } from '../router/router'

export default {
    name: "Main",
    setup() {
        const state = reactive({
            tabs: readonly(mainRoutes)
        })

        return {
            state
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