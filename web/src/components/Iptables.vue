<template>
<div class="columns" id="iptable">
    <div class="column p-0">
        <div class="columns">
            <div class="column">
                <article class="message is-dark">
                <div class="message-body content is-small">
                    <p>格式化IPTables规则，使其更加直观，使用方法: </p>
                    <ol type="1">
                        <li>登录服务器执行`iptables-save`，复制输出的内容。</li>
                        <li>粘贴到文本框，选择一个数据流方向，然后点击美化。</li>
                        <li>数据流一般有三种：发往本机(in)，转发(forward)，本机发出(out)。</li>
                    </ol>
                    <p>附<a href="https://thermalcircle.de/lib/exe/fetch.php?media=linux:nf-hooks-iptables1.png">Netfilter中数据包流图</a>。</p>
                </div>
                </article>
            </div>
        </div>
        <div class="columns">
            <div class="column pt-0 pb-0">
                <div class="tags">
                    <span :class="[ 'tag', { 'is-dark' : idx == state.current } ]" v-for="(tag, idx) in tags" :key="idx" @click="switchTag(idx)">{{ tag }}</span>
                </div>
            </div>
        </div>
        <div class="columns">
            <div class="column">
                <textarea class="textarea" v-model="state.input" placeholder="请输入原始内容"></textarea>
                <div class="buttons is-centered mt-4">
                    <button :class="[ 'button', 'is-info', { 'is-loading' : state.loading } ]" @click="formatTable">美化一下试试</button>
                </div>
            </div>
        </div>
        <div class="columns" v-if="state.chainData.length > 0">
            <div class="column">
            <div id="iptable-panel">
                <chains :chains="state.chainData" :key="state.counter"></chains>
            </div>
            </div>
        </div>
    </div>
</div>
</template>
    
<script>
import { reactive, readonly } from 'vue'
import {Post} from "../tools/http"

export default {
    name: "Iptables",
    setup() {
        const state = reactive({
            loading: false,
            current: 0,
            counter: 0,
            input: "",
            chainData: []
        })
        const tags = readonly(["in", "forward", "out"])

        function switchTag(idx) {
            state.current = idx;
        }

        async function formatTable() {
            if (state.input == "") {
                return
            }
            state.loading = true
            let resp = await Post("/api/iptables", {
                data: state.input,
                type: tags[state.current]
            })
            state.loading = false
            if (resp.data.code === 10000) {
                state.counter++ // just force component re-render
                state.chainData = resp.data.data
            }
        }
        return {
            state,
            tags,
            switchTag,
            formatTable
        }
    }
}
</script>

<style css>
textarea {
    font-family: ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,Liberation Mono,Courier New,monospace;
    font-size: 0.75rem;
}
#iptable-panel {
    border: 1px solid #ccc;
    padding: 0.5rem 1rem 0.5rem 0rem;
    font-family: ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,Liberation Mono,Courier New,monospace;
}
.chain-info .chain-name, .chain-info .chain-policy{
    display: inline-block;
    margin-right: 2rem;
}
.chain-info {
    background-color: rgb(226 232 240);
    padding-left: 0.25rem;
    margin: 0.25rem 0;
}
.chain-info.empty span {
    color: #aaa;
}
.rule {
    font-size: 0.8rem;
    margin: 4px 0px;
    background-color: rgb(248 250 252);
    padding-left: 0.25rem;
    border-left: 4px solid #79bda8;
}
.rule-ref {
    text-decoration-line: underline;
    text-decoration-style: wavy;
    cursor: pointer;
}
.rule-match span.text-orange {
    border: 1px solid #F59E0B;
    border-radius: 0.25rem;
    margin: 0px 0.25rem;
    padding: 0 0.25rem;
    font-size: 0.75rem;
    line-height: 0.75rem;
}
#iptable-panel .rules,#iptable-panel .chain-info {
    margin-left: 1rem;
}
.text-gray {
    color: #6B7280;
}
.text-blue {
    color: #3B82F6;
}
.text-black {
    color: #1F2937;
}
.text-orange {
    color: #F59E0B;
}
.text-red {
    color: #EF4444;
}
.text-purple {
    color: #8B5CF6;
}
</style>