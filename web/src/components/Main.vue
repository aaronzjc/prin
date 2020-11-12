<template>
<main class="container">
    <nav class="navbar" role="navigation" aria-label="main navigation">
        <div class="navbar-brand">
            <a class="navbar-item" href="/">
                <img alt="简单的小公举" src="../assets/logo.png" />
                <p>一些简单的小公举</p>
            </a>
        </div>
    </nav>
    <div class="columns">
        <div class="column">
            <div class="tabs">
                <ul>
                    <router-link v-for="(tab, idx) in state.tabs" v-bind="$attrs" :key="idx" :to="tab.link" custom v-slot="{ isActive, navigate }">
                    <li :class="isActive ? 'is-active' : '' " @click="navigate" role="link">
                        <a href="#">{{ tab.name }}</a>
                    </li>
                    <slot />
                    </router-link>
                </ul>
            </div>
        </div>
    </div>
    <router-view></router-view>
    <div class="columns" id="footer">
        <div class="column copyright has-text-centered">
            <p><a href="https://github.com/aaronzjc">@aaronzjc</a>开发，源码<a href="https://github.com/aaronzjc/prin">在此</a>，欢迎Star v0.1</p>
        </div>
    </div>
</main>
</template>

<script>
import { reactive } from 'vue'

export default {
    name: "Main",
    setup() {
        const state = reactive({
            tabs: [
                {
                    "name": "二维码",
                    "link": "qrcode"
                },
                {
                    "name": "时间",
                    "link": "time"
                }
            ],
        })

        return {
            state
        }
    }
}
</script>

<style lang="scss">
    .container {
        max-width: 960px !important;
    }
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
#footer {
    margin-top: 100px;
    background: none;
}
</style>