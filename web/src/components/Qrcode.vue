<template>
<div class="columns" id="qrcode">
    <div class="column">
        <textarea class="textarea" v-model="state.content" placeholder="请输入内容"></textarea>
        <div class="buttons is-centered mt-4">
            <button :class="[ 'button', 'is-info', { 'is-loading' : state.loading } ]" @click="genCode">生成二维码</button>
        </div>
        <div class="history">
            <div class="title is-size-6">历史记录</div>
            <div class="item" v-for="(log, idx) in state.history" :key="idx" @click="setData(idx)"><span class="data">{{ log }}</span><span class="tag is-warning refresh" @click.stop="removeHistory(idx)">删除</span></div>
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
import {Post} from "../tools/http"
import {add, get, remove} from "../tools/lru"

export default {
    name: "Qrcode",
    setup() {
        const state = reactive({
            history: get(),
            loading: false,
            content: "",
            img: ""
        })

        async function genCode() {
            if (state.content == "") {
                state.img = "";
                return false;
            }
            state.loading = true
            let resp = await Post("/api/qrcode", {content: state.content})
            state.loading = false
            if (resp.data.code === 10000) {
                state.img = "data:image/png;base64," + resp.data.data.qrcode;
            }
            state.history = add(state.content)
        }

        async function removeHistory(idx) {
            state.history.splice(idx, 1)
            remove(idx)
        }

        function setData(idx) {
            state.content = state.history[idx]
        }

        return {
            state,
            genCode,
            removeHistory,
            setData
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
    .history {
        .item {
            &:hover {
                cursor: pointer;
            }
            margin: 0 0 8px 0px;
            padding: 4px;;
            display: flex;
            flex-direction: row;
            justify-content: space-between;
            align-items: center;

            border: 1px dashed #096;

            .data {
                margin-right: 4px;
                word-break: break-all;
            }
        }
    }
}
</style>
