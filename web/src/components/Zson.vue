<template>
    <div class="columns" id="jsontogo">
        <div class="column p-0">
            <div class="columns">
                <div class="column pt-0 pb-0">
                    <div class="buttons are-small">
                        <button class="button is-light"  @click="format">格式化</button>
                        <button class="button is-light"  @click="tar">压缩</button>
                        <button class="button is-light"  @click="unesc">去除转义</button>
                    </div>
                </div>
            </div>
            <div class="columns" id="code-block">
                <div class="column json-editor">
                    <v-ace-editor
                        v-model:value="state.content"
                        lang="json"
                        @init="initEditor"
                        theme="tomorrow"
                        style="height: 600px" />
                </div>
                <div class="column is-narrow pl-0">
                    <div id="tag-tab" class="field is-grouped is-flex-direction-column is-align-items-center">
                        <div class="control" v-for="(page, idx) in state.pages" :key="idx">
                            <div class="tags has-addons">
                                <span @click="swit(idx)" :class="[ 'tag', { 'is-grey': idx != state.current }, { 'is-info' : idx == state.current } ]">草稿</span>
                                <a class="tag is-delete" @click="del(idx)"></a>
                            </div>
                        </div>
                        <div class="control add" @click="add" v-if="state.pages.length < 9">
                            <button class="button is-small is-warning">新建草稿</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { reactive } from '@vue/runtime-core'
import { VAceEditor } from 'vue3-ace-editor';
import 'ace-builds/src-noconflict/mode-json';
import 'ace-builds/src-noconflict/theme-tomorrow';

export default {
    name: "Zson",
    components: {
        VAceEditor
    },
    setup() {
        const state = reactive({
            pages: [],
            current: -1,
            content: "",
            editor: null
        })

        function add() {
            if (state.pages.length > 0) {
                state.pages[state.current] = state.content
            }
            state.pages.push(JSON.stringify({"hello": "hello from prin ."}, null, 4))
            swit(state.pages.length-1)
        } // 新增
        function del(idx) {
            if (state.pages.length <= 1) {
                return 
            }
            state.pages.splice(idx, 1)
            state.current--
            if (state.current < 0) {
                state.current = 0
            }
            state.content = state.pages[state.current]
        } // 删除
        function swit(idx) {
            state.pages[state.current] = state.content
            state.content = state.pages[idx]
            state.current = idx
        } // 切换
        // 格式化
        function format() {
            try {
                state.content = JSON.stringify(JSON.parse(state.content), null, 4)
            } catch (e) {
                alert("格式化失败")
            }
        }

        function tar() {
            try {
                state.content = JSON.stringify(JSON.parse(state.content))
            } catch (e) {
                alert("压缩失败")
            }
        }

        function unesc() {
            state.content = state.content.replace(/\\\\/g, '\\').replace(/\\\"/g, '"')
        }

        function initEditor(editor) {
            state.editor = editor
            editor.setOption("showPrintMargin", false);
            editor.setOption("fontSize", 12.5);
            editor.setOption("wrap", true);

            add()
        }
        return {
            state,
            add,
            del,
            swit,
            format,
            tar,
            unesc,
            initEditor
        }
    }
}
</script>

<style>
div.json-editor > .ace_editor{
    border: 1px solid #dadada;
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