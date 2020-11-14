<template>
<div class="columns" id="cert">
    <div class="column p-0">
        <div class="columns">
            <div class="column">
                <article class="message is-dark">
                <div class="message-body content is-small">
                    <p>使用方法: </p>
                    <ol type="1">
                        <li>下载根证书，添加本地信任，不同系统方法不同</li>
                        <li>在下面生成你的域名证书和私有Key</li>
                        <li>使用生成的证书和私钥部署HTTPS即可</li>
                        <li>证书仅供开发测试环境使用，请谨慎使用</li>
                    </ol>
                    <p>本站生成的所有证书使用的是同一的根证书，因此，只需要信任一次根证书即可。如后续变更，会在此通知。</p>
                </div>
                </article>
            </div>
        </div>
        <div class="columns">
            <div class="column">
            <div class="buttons">
                <a href="/assets/ca.pem" target="_blank" class="button is-light is-small">查看根证书</a>
                <a href="/assets/ca.key" target="_blank" class="button is-light is-small">查看根证书私钥</a>
            </div>
            </div>
        </div>
        <div class="columns">
            <div class="column">
            <div class="field is-grouped">
                <p class="control is-expanded">
                    <input class="input" type="text" v-model="state.domains" placeholder="输入域名，多个域名用英文逗号区分">
                </p>
                <p class="control">
                    <a :class="[ 'button', 'is-info', { 'is-loading' : state.loading } ]" @click="generate">生成证书和私钥</a>
                </p>
            </div>
            </div>
        </div>
        <div class="columns">
            <div class="column">
                <span>domain.crt</span>
                <textarea class="textarea" rows="10" v-model="state.cert" readonly placeholder="域名证书"></textarea>
            </div>
            <div class="column">
                <span>domain.key</span>
                <textarea class="textarea" rows="10" v-model="state.key" readonly placeholder="域名私钥"></textarea>
            </div>
        </div>
    </div>
</div>
</template>

<script>
import { reactive } from 'vue'
import {Post} from "../../tools/http"

export default {
    name: "Cert",
    setup() {
        const state = reactive({
            domains: "",
            cert: "",
            key: "",
            loading: false,
        })

        async function generate() {
            if (state.domains == "") {
                state.cert = ""
                state.key = ""
                return
            }
            state.loading = true
            let resp = await Post("/api/cert", {
                domains: state.domains
            })
            state.loading = false
            if (resp.data.code === 10000) {
                state.cert = resp.data.data.cert
                state.key = resp.data.data.key
            }
        }

        return {
            state,
            generate,
        }
    }
}
</script>