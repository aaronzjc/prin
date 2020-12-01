<template>
<div class="columns" id="iptable">
    <div class="column p-0">
        <div class="columns">
            <div class="column">
                <article class="message is-dark">
                <div class="message-body content is-small">
                    <p>格式化IPTables规则，使其更加直观。使用方法: </p>
                    <ol type="1">
                        <li>登录服务器执行：iptables -t {{ tags[state.current] }} -L</li>
                        <li>粘贴输出的内容，然后点击美化</li>
                    </ol>
                    <p>因为IPTable都是从默认的链开始匹配。所以，美化后的输出，根链为每个表的默认链，不会展示孤儿链。附<a target="_blank" href="https://thermalcircle.de/lib/exe/fetch.php?media=linux:nf-hooks-iptables1.png">Netfilter中数据包流图。</a></p>
                </div>
                </article>
            </div>
        </div>
        <div class="columns">
            <div class="column pt-0">
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
        <div class="columns" v-if="state.tree.children.length > 0">
            <div class="column">
            <div id="iptable-panel">
                <ul>
                    <tree-item
                        class="item"
                        :item="state.tree"
                    ></tree-item>
                </ul>
            </div>
            </div>
        </div>
    </div>
</div>
</template>

<script>
import { reactive, readonly } from 'vue'
import {Post} from "../tools/http"
import TreeItem from "./ext/TreeItem"

export default {
    name: "Iptable",
    setup() {
        const state = reactive({
            loading: false,
            current: 0,
            input: "",
            tree: {
                name: "iptables",
                children: []
            },
        })
        const tags = readonly(["NAT", "FILTER", "MANGLE"])

        function switchTag(idx) {
            state.current = idx;
        }

        async function formatTable() {
            if (state.input == "") {
                return
            }
            if (state.current != 0) {
                alert("暂时只支持NAT，其他链在适配中，敬请期待")
                return
            }
            state.loading = true
            state.tree.children = []
            let resp = await Post("/api/iptable", {
                data: state.input,
                table: tags[state.current]
            })
            state.loading = false
            if (resp.data.code === 10000) {
                state.tree.children = resp.data.data
            }
        }
        return {
            state,
            tags,
            switchTag,
            formatTable
        }
    },
    components: {
        TreeItem
    }
}
</script>

<style lang="scss">
#iptable-panel {
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
    border: 1px solid #ccc;
    li {
        margin: 4px 0 4px 1rem;
    }
    .bold .chain-name {
        cursor: pointer;
        display: inline-block;
    }
    .chain-name, .chain-raw, .chain-format {
        font-size: 12px;
    }
    .chain-name {
        color: hsl(217, 71%, 53%);
    }
    .chain-raw {
        color: hsl(0, 0%, 71%);
    }
    .chain-format {
        color: hsl(0, 0%, 29%);
    }
}
</style>
