<template>
    <div class="columns" id="jsontogo">
        <div class="column p-0">
            <div class="columns">
                <div class="column pt-0 pb-0">
                    <div class="buttons are-small">
                        <button class="button is-info" @click="switchMode">{{ state.mode == "tree" ? "源码" : "树形" }}</button>
                        <button class="button is-light" v-if="state.mode == 'code'" @click="format">格式化</button>
                        <template v-if="state.mode == 'tree'">
                            <button class="button is-light" @click="collapse">折叠全部</button>
                            <button class="button is-light" @click="expand">展开全部</button>
                        </template>
                    </div>
                </div>
            </div>
            <div class="columns" id="code-block">
                <div class="column">
                    <div id="jsoneditor" style="height: 600px"></div>
                </div>
                <div class="column is-narrow pl-0">
                    <div id="tag-tab" class="field is-grouped is-flex-direction-column is-align-items-center">
                        <div class="control" v-for="(page, idx) in state.pages" :key="idx">
                            <div class="tags has-addons">
                                <span @click="switchPage(idx)" :class="[ 'tag', { 'is-grey': idx != state.current }, { 'is-info' : idx == state.current } ]">{{ page.title }}</span>
                                <a class="tag is-delete" @click="delPage(idx)"></a>
                            </div>
                        </div>
                        <div class="control add" @click="addPage" v-if="state.pages.length < 9">
                            <button class="button is-small is-warning">ADD  PAGE</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { onMounted, reactive } from '@vue/runtime-core'
import JSONEditor from "jsoneditor"
import "jsoneditor/dist/jsoneditor.css"

export default {
    name: "Zson",
    setup() {
        const state = reactive({
            pages: [],
            current: -1,
            idx: 0,
            editor: undefined,
            mode: "code"
        })
        
        onMounted(() => {
            document.getElementById('jsoneditor').innerHTML = ""
            const container = document.getElementById('jsoneditor')
            const options = {
                mode: state.mode,
                indentation: 4,
                mainMenuBar: false,
                navigationBar: false,
                statusBar: false,
                modes: ['code', 'tree'],
                onError: function (err) {
                    alert(err.toString())
                },
                onModeChange: function (newMode, oldMode) {
                    console.log('Mode switched from', oldMode, 'to', newMode)
                },
                timestampTag: function() {
                    return false
                }
            }
            state.editor = new JSONEditor(container, options)
            addPage()
        })

        function format() {
            try {
                JSON.parse(state.editor.getText())
            } catch (err) {
                return false
            }
            state.editor.set(state.editor.get())
        }
        function switchMode() {
            state.mode = state.mode == "code" ? "tree" : "code"
            state.editor.setMode(state.mode)
            
        }
        function expand() {
            state.editor.expandAll()
        }
        function collapse() {
            state.editor.collapseAll()
        }

        function addPage() {
            let maxIdx = 0
            if (state.pages.length > 0) {
                maxIdx = state.pages[state.pages.length-1]["idx"]
            }
            state.idx = ++maxIdx
            let title = "第 " + (state.idx < 10 ? "0" + state.idx : state.idx) + " 页" 
            let tpl = {
                title: title,
                idx: state.idx,
                data: {"hi": "this is page " + state.idx }
            }
            state.pages.push(tpl)
            state.current++
            state.editor.set(state.pages[state.current]["data"])
        }
        function delPage(idx) {
            if (state.pages.length <= 1) {
                return 
            }
            state.pages.splice(idx, 1)
            state.current--
            if (state.current < 0) {
                state.current = 0
            }
            state.editor.set(state.pages[state.current]["data"])
        }
        function switchPage(idx) {
            // 暂存修改过的数据
            let validJson = {}
            try {
                validJson = state.editor.get()
            } catch(err) {
                if (!confirm("当前JSON无效，切换页签将丢失数据，是否继续")) {
                    return
                }
            }
            state.mode = "code"
            state.editor.setMode(state.mode)
            state.pages[state.current]["data"] = validJson

            // 切换到新的页签
            state.editor.set(state.pages[idx]["data"])
            state.current = idx
        }
        return {
            state,
            switchMode,
            expand,
            collapse,
            format,
            addPage,
            delPage,
            switchPage,
        }
    }
}
</script>

<style>
div.jsoneditor,
div.jsoneditor-menu {
    border-color: #dadada;
}
.ace-jsoneditor .ace_gutter {
    background: #fff;
    color: #333
}
.ace-jsoneditor .ace_gutter-active-line {
    background: #fff;;
}
#tag-tab .control {
    margin-right: 0.5rem;
    margin-bottom: 0.5rem;
}
#tag-tab .tag {
    border-radius: 0;
}
#tag-tab .tag:hover {
    cursor: pointer;
}
</style>