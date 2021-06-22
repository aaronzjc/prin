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
                <div
                    v-show="!state.hideLeft"
                    :class="['column', { 'is-half': !state.hideLeft }]"
                >
                    <textarea
                        class="textarea"
                        v-model="str"
                        placeholder="输入Json字符串"
                        style="height: 50rem; max-height: 50rem"
                    ></textarea>
                </div>
                <div :class="['column', { 'is-half': !state.hideLeft }]">
                    <div class="output panelContent" style="height: 50rem">
                        <template v-if="state.err != ''"
                            ><p>{{ state.err }}</p></template
                        >
                        <template v-else>
                            <table class="treeTable" cellspacing="0" cellpadding="0">
                                <thead>
                                    <tr>
                                        <td></td>
                                        <td style="width: 100%"></td>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr class="treeRow">
                                        <td>{{ state.arr ? "[" : "{" }}</td>
                                        <td></td>
                                    </tr>
                                    <template
                                        v-for="(item, idx) in state.rows"
                                        :key="idx"
                                    >
                                        <tr
                                            :class="[
                                                'treeRow',
                                                { hasChildren: item.ar > 0 },
                                            ]"
                                            v-if="item.sh"
                                        >
                                            <td
                                                class="treeLabelCell"
                                                :style="{
                                                    '--tree-label-cell-indent':
                                                        item.flr * 16 + 'px',
                                                }"
                                                @click="toggle(idx, !item.op)"
                                            >
                                                <span
                                                    :class="[
                                                        'treeIcon',
                                                        { open: item.ar > 0 && item.op },
                                                        {
                                                            close:
                                                                item.ar > 0 && !item.op,
                                                        },
                                                    ]"
                                                ></span>
                                                <span class="treeLabel">{{
                                                    item.k
                                                }}</span>
                                            </td>
                                            <td class="treeValueCell">
                                                <span :class="['jsonT' + item.t]">{{
                                                    item.v
                                                }}</span>
                                            </td>
                                        </tr>
                                    </template>
                                    <tr class="treeRow">
                                        <td>{{ state.arr ? "]" : "}" }}</td>
                                        <td></td>
                                    </tr>
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

const AR_NONE = -1 // 无展开
const AR_ARRAY = 1 // 数组展开
const AR_OBJ = 2 // 对象展开

export default {
    name: "MyJson",
    setup() {
        const state = reactive({
            err: "",
            arr: false,

            rows: [],
            idx: {},

            hideLeft: false,
        });
        const str = ref('{"name":"kate"}');

        function subPath(path, subKey) {
            return path + "/" + String(subKey).replace(/[\\/]/g, "\\$&");
        }

        function getArrow(val) {
            if (typeof val != "object" || val == null || val == undefined) {
                return -1;
            }
            let children = Object.keys(val).length > 0;
            if (!children) {
                // 不能展开
                return AR_NONE;
            }
            if (val instanceof Array) {
                // 数组展开
                return AR_ARRAY;
            }
            // 对象展开
            return AR_OBJ;
        }

        function getMemberVal(member) {
            if (member.type != "object") {
                if (member.type == "string") {
                    return '"' + member.value + '"'
                }
                if (member.type == "number") {
                    return member.value
                }
            }
            if (member.value == null) {
                return "null"
            }
            if (member.arrow < 0) {
                if (member.type instanceof Array) {
                    return "[]";
                }
                return "{}";
            }
            return "";
        }

        function getMembers(parent, level, path) {
            if (typeof parent == "string") {
                return [];
            }
            let members = [];
            for (const prop in parent) {
                members.push({
                    name: prop,
                    value: parent[prop],
                    type: typeof parent[prop],
                    level: level,
                    path: subPath(path, prop),
                    arrow: getArrow(parent[prop]),
                });
            }
            return members;
        }

        function getRows(parent, level = 0, path = "") {
            let rows = [];
            const members = getMembers(parent, level, path);
            members.map((member) => {
                rows.push({
                    id: member.path,
                    k: member.name,
                    v: getMemberVal(member),
                    t: member.type,
                    flr: member.level + 1,
                    ar: member.arrow,
                    op: member.arrow > 0, // 默认打开
                    sh: true, // 默认展示
                });
                if (member.arrow > 0) {
                    const childRows = getRows(member.value, level + 1, member.path);
                    rows = rows.concat(childRows);
                }
            });
            return rows;
        }

        function toggle(startIdx, open) {
            if (state.rows[startIdx]["ar"] == AR_NONE) {
                return;
            }
            function handle(nodeIdx, open) {
                var fatherNode = state.rows[nodeIdx];
                var i = nodeIdx;
                while (i < state.rows.length) {
                    // 如果操作的是当前节点，那么修改打开状态
                    if (fatherNode["id"] == state.rows[i]["id"]) {
                        state.rows[i]["op"] = open;
                        i++;
                        continue;
                    }
                    // 如果和id不匹配了，说明遍历完了
                    if (state.rows[i]["id"].indexOf(fatherNode["id"]) == -1) {
                        return i;
                    }
                    if (fatherNode["op"] == false) {
                        state.rows[i]["sh"] = false;
                    } else {
                        state.rows[i]["sh"] = true;
                        if (state.rows[i]["ar"] != AR_NONE) {
                            // 如果这个也是一个父节点，那么，递归处理它的子节点
                            i = handle(i, state.rows[i]["op"]);
                            continue;
                        }
                    }
                    i++;
                }
                return i;
            }
            var lastIdx = handle(startIdx, open);
            if (open == false) {
                if (state.rows[startIdx]["ar"] == AR_ARRAY) {
                    state.rows[startIdx]["v"] = "[…]";
                } else if (state.rows[startIdx]["ar"] == AR_OBJ) {
                    state.rows[startIdx]["v"] = "{…}";
                }
            } else {
                state.rows[startIdx]["v"] = "";
            }
            return lastIdx
        }

        /*
        function expandAll(maxDepth = 0) {
            var i = 0
            while (i < state.rows.length) {
                if (state.rows[i].ar == AR_NONE) {
                    continue
                }
                if (maxDepth > 0 && state.rows[i].flr > maxDepth) {
                    break
                }
                i = toggle(i, true)
            }
        }
        */

        function parse() {
            if (str.value == "") {
                return;
            }
            state.err = "";
            state.arr = false;
            try {
                var jsonObj = JSON.parse(str.value);
                state.rows = getRows(jsonObj);
                for (var i in state.rows) {
                    state.idx[state.rows[i].id] = i;
                }
                if (typeof jsonObj == "object") {
                    if (jsonObj instanceof Array) {
                        state.arr = true;
                    }
                }
            } catch (err) {
                state.err = err;
            }
        }

        onMounted(() => {
            parse();
        });
        watch(str, () => {
            parse();
        });

        return {
            str,
            state,
            toggle,
        };
    },
};
</script>

<style lang="scss">
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */
 /* copy & modify from mozilla firefox devtools/client/shared/components/tree/TreeView.css */
.output,
.textarea {
    font-family: "dejavu sans mono","droid sans mono",Menlo,monaco,"lucida console","courier new",courier,monospace,sans-serif;
    font-size: 12px;
    scrollbar-width: 2px;
}
.output {
    direction: ltr;
    overflow: auto;
    border: 1px solid #dadada;
    padding: 8px 4px;
}

.treeTable .treeRow:hover {
    background-color: rgb(240, 249, 254);
}

.treeTable .treeLabelCell,
.treeTable .treeValueCell {
    padding: 2px 0;
    padding-inline-start: 4px;
    line-height: 16px; /* make rows 20px tall */
    vertical-align: top;
    overflow: hidden;
}

.treeTable .treeLabelCell {
    white-space: nowrap;
    cursor: default;
    padding-inline-start: var(--tree-label-cell-indent);
}

.treeTable .treeLabelCell::after {
    content: ":";
    color: var(--object-color);
}

.treeTable .treeValueCell > [aria-labelledby],
.treeTable .treeLabelCell > .treeLabel {
    unicode-bidi: plaintext;
    text-align: match-parent;
}

/* No padding if there is actually no label */
.treeTable .treeLabel:empty {
    padding-inline-start: 0;
}

.treeTable .treeRow.hasChildren > .treeLabelCell > .treeLabel:hover {
    cursor: pointer;
    text-decoration: underline;
}

.treeTable .treeRow .treeIcon {
    box-sizing: content-box;
    height: 14px;
    width: 14px;
    padding: 1px;
    /* Set the size of loading spinner (see .devtools-throbber) */
    font-size: 10px;
    line-height: 16px;
    display: inline-block;
    vertical-align: bottom;
    /* Use a total width of 20px (margins + padding + width) */
    margin-inline: 3px 1px;
    background-size: 10px;
    background-position: center;
    color: #fff;
    &.open {
        background-image: url("../assets/arrow.svg");
        fill:rgba(135, 135, 137, 0.9);
    }
    &.close {
        transform: rotate(-90deg);
        background-image: url("../assets/arrow.svg");
        fill:rgba(135, 135, 137, 0.9);
    }
}

/* All expanded/collapsed styles need to apply on immediate children
   since there might be nested trees within a tree. */
.treeTable .treeRow.hasChildren > .treeLabelCell > .treeIcon {
    cursor: pointer;
    background-repeat: no-repeat;
}
.treeLabelCell {
    .treeLabel {
        color: rgb(0, 116, 232);
    }
}
.treeValueCell {
    .jsonTstring {
        color: rgb(221, 0, 169);
    }
    .jsonTnumber {
        color: rgb(5, 139, 0);
    }
    .jsonTboolean {
        color: blue;
    }
    .jsonTobject {
        color: #333;
    }
}
</style>
