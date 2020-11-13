<template>
<div class="columns" id="qrcode">
    <div class="column">
        <textarea class="textarea" v-model="state.content" placeholder="请输入内容"></textarea>
        <div class="buttons is-centered mt-4">
            <button class="button is-primary" @click="genCode">生成二维码</button>
        </div>
    </div>
    <div class="column">
        <div class="img-box">
            <img :src="state.img" v-if="state.img != ''">
            <span v-else>创建一个试试</span>
        </div>
    </div>
</div>
</template>

<script>
import { reactive } from 'vue'
import {Post} from "../../tools/http"

export default {
    name: "Qrcode",
    setup() {
        const state = reactive({
            content: "",
            img: ""
        })

        async function genCode() {
            if (state.content == "") {
                state.img = "";
                return false;
            }
            let resp = await Post("/api/qrcode", {content: state.content})
            if (resp.data.code === 10000) {
                state.img = "data:image/png;base64," + resp.data.data.qrcode;
            }
        }

        return {
            state,
            genCode
        }
    }
}
</script>

<style lang="scss">
#qrcode {
    .img-box {
        border:2px dashed #096;
        width: 256px;
        height: 256px;
        text-align: center;
        margin: 0 auto;
        span {
            line-height: 256px;
        }
    }
}
</style>