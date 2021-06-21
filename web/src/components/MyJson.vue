<template>
    <div class="columns" id="jsontogo">
        <div class="column p-0">
            <div class="columns">
                <div class="column pt-0 pb-0">
                    <div class="buttons are-small">
                        <button
                            class="button is-light"
                            @click="state.hideLeft = !state.hideLeft"
                        >
                            {{ state.hideLeft ? "显示左栏" : "隐藏左栏" }}
                        </button>
                    </div>
                </div>
            </div>
            <div class="columns" id="code-block">
                <div v-show="!state.hideLeft" :class="[ 'column', { 'is-half' : !state.hideLeft } ]">
                    <textarea
                        class="textarea scroller"
                        v-model="str"
                        placeholder="输入Json字符串"
                        style="height: 50rem;max-height: 50rem;"
                    ></textarea>
                </div>
                <div :class="[ 'column', { 'is-half' : !state.hideLeft } ]">
                    <div class="output panelContent scroller" style="height: 50rem;">
                        <template v-if="state.err != ''"><p>{{ state.err }}</p></template>
                        <template v-else>
                            <table class="jsonTable" cellspacing="0" cellpadding="0">
                                <tbody>
                                    <tr class="jsonRow" v-if="state.arr"><td>[</td><td></td></tr>
                                    <tr class="jsonRow" v-else><td>{</td><td></td></tr>
                                    <template v-for="(item, id) in state.format" :key="id">
                                        <tr class="jsonRow" v-show="item.sh">
                                            <td
                                                class="jsonCell jsonCellK"
                                                :style="{
                                                    '--tree-label-cell-indent':
                                                        item.flr * 16 + 'px',
                                                }"
                                                @click="toggle(item.id, !item.op)"
                                            >
                                                <span
                                                    :class="[ 'jsonIcon', { 'open': item.ar > 0 && !item.op }, { 'close': item.ar > 0 && item.op } ]"></span><span class="jsonK">{{ item.k }}</span>
                                            </td>
                                            <td class="jsonCell">
                                                <span :class="[ 'jsonV', 'jsonT' + item.t]">{{ item.t == 'string' ? '"' + item.v + '"' : ( item.v === null ? "null" : item.v) }}</span>
                                            </td>
                                        </tr>
                                    </template>
                                    <tr class="jsonRow" v-if="state.arr"><td>]</td><td></td></tr>
                                    <tr class="jsonRow" v-else><td>}</td><td></td></tr>
                                </tbody>
                            </table>
                        </template>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, reactive, watch, onMounted } from "vue";

export default {
    name: "MyJson",
    setup() {
        const state = reactive({
            err: "",
            arr: false,
            format: {},
            idx: {},
            hideLeft: false,
        });
        const str = ref('{"name":"kate"}')

        function subPath(path, subKey) {
            return path + "/" + String(subKey).replace(/[\\/]/g, "\\$&")
        }

        function getArrow(val) {
            if (typeof val != "object" || val == null || val == undefined) {
                return -1
            }
            var children = Object.keys(val).length > 0
            if (!children) {
                return -1
            }
            if (val instanceof Array) {
                return 1
            }
            return 2
        }

        function getMemberVal(member) {
            if (member.value == null || member.type != 'object') {
                return member.value
            }
            return ""
        }

        function getMembers(parent, level, path) {
            if (typeof parent == "string") {
                return []
            }
            let members = []
            for (const prop in parent) {
                members.push({
                    name: prop,
                    value: parent[prop],
                    type: typeof parent[prop],
                    level: level,
                    path: subPath(path, prop),
                    arrow: getArrow(parent[prop]),
                })
            }
            return members
        }

        function getRows(parent, level = 0, path = "") {
            let rows = []
            const members = getMembers(parent, level, path)
            members.map(member => {
                rows.push({
                    id: member.path,
                    k: member.name,
                    v: getMemberVal(member),
                    t: member.type,
                    flr: member.level + 1,
                    ar: member.arrow,
                    op: member.arrow > 0, // 默认打开
                    sh: true, // 默认展示
                })
                if (member.arrow > 0) {
                    const childRows = getRows(member.value, level+1, member.path)
                    rows = rows.concat(childRows)
                }
            })
            return rows
        }

        function toggle(id, open) {
            var startIdx = state.idx[id]
            if (state.format[startIdx]["ar"] < 0) {
                return
            }
            function handle(id, open) {
                console.log(id, open);
                var fatherIdx = state.idx[id];
                var fatherNode = state.format[fatherIdx];
                var i = fatherIdx;
                while (i < state.format.length) {
                    // 如果操作的是当前节点，那么修改打开状态
                    if (id == state.format[i]["id"]) {
                        state.format[i]["op"] = open;
                        i++;
                        continue;
                    }
                    // 如果和id不匹配了，说明遍历完了
                    if (state.format[i]["id"].indexOf(id) == -1) {
                        return i;
                    }
                    if (fatherNode["op"] == false) {
                        state.format[i]["sh"] = false;
                    } else {
                        if (fatherNode["sh"]) {
                            state.format[i]["sh"] = true;
                        }
                        if (state.format[i]["t"] == 'object' && state.format[i]["ar"] > 0) {
                            // 如果这个也是一个父节点，那么，递归处理它的子节点
                            i = handle(state.format[i]["id"], state.format[i]["op"]);
                            continue;
                        }
                    }
                    i++;
                }
                return i;
            }
            handle(id, open);
            if (open == false) {
                if (state.format[startIdx]["ar"] == 1) {
                    state.format[startIdx]["v"] = "[...]";
                } else if (state.format[startIdx]["ar"] == 2) {
                    state.format[startIdx]["v"] = "{...}";
                }
            } else {
                state.format[startIdx]["v"] = "";
            }
        }

        function run() {
            if (str.value == "") {
                return
            }
            state.err = ""
            try {
                var jsonObj = JSON.parse(str.value)
                state.format = getRows(jsonObj)
                for (var i in state.format) {
                    state.idx[state.format[i].id] = i
                }
                if (typeof jsonObj == 'object') {
                    if (jsonObj instanceof Array) {
                        state.arr = true
                    }
                }
            } catch (err) {
                state.err = err
            }
        }

        onMounted(() => {
            run()
        })
        watch(str, () => {
            run()
        })

        return {
            str,
            state,
            toggle,
        };
    },
};
</script>

<style lang="scss">
.output, .textarea {
    font-family: Menlo, monospace;
    font-size: 12px;
    scrollbar-width: 2px;
}
.output {
    direction: ltr;
    overflow: auto;
    border: 1px solid #dadada;
    padding: 8px 4px;
}
.jsonTable {
    width: 100%;
}
.jsonRow {
    width: 100%;
    height: 20px;
    &:hover {
        background: rgba(221, 225, 228, 0.66);
    }
    .jsonCell {
        padding-right: 12px;
        height: 16px;
        line-height: 16px;;
        padding: 2px 0px;

        &.jsonCellK {
            width: 1%;
            white-space: nowrap;
            padding-right: 8px;
            color: rgb(0, 116, 232);
            padding-inline-start: var(--tree-label-cell-indent);
        }
    }
    .jsonIcon {
        box-sizing: content-box;
        height: 14px;
        width: 14px;
        padding: 1px;
        font-size: 10px;
        line-height: 14px;
        display: inline-block;
        vertical-align: bottom;
        margin-inline: 3px 1px;
        cursor: pointer;

        background-repeat: no-repeat;
        background-position: center;
        background-size: 10px;
        
        &.open {
            transform: rotate(-90deg);
            background-image: url("../assets/arrow.svg");
        }
        &.close {
            background-image: url("../assets/arrow.svg");
        }
    }
    .jsonK {
        display: inline-block;
        white-space: nowrap;
        line-height: 16px;
        text-align: match-parent;
        unicode-bidi:plaintext;
    }
    .jsonK::after {
        content: ":";
    }
    .jsonV {
        color: #096;
        &.jsonTstring {
            color: rgb(221, 0, 169);
        }
        &.jsonTnumber {
            color: rgb(5, 139, 0);
        }
        &.jsonTboolean {
            color:blue;
        }
        &.jsonTobject {
            color: grey;
        }
    }
}
</style>
