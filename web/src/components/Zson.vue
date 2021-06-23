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
            "json": {
                "string": "Hello World"
            },
            editor: undefined,
            mode: "code"
        })
        
        onMounted(() => {
            if (state.editor != undefined) {
                return true
            }
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

            state.editor = new JSONEditor(container, options, state.json)
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
        return {
            state,
            switchMode,
            expand,
            collapse,
            format
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
</style>