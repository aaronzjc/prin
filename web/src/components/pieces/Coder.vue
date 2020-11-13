<template>
<div class="columns" id="coder">
    <div class="column p-0">
    <div class="columns">
        <div class="column pt-0">
        <div class="tags">
            <span :class="[ 'tag', { 'is-dark' : idx == state.current } ]" v-for="(tag, idx) in tags" :key="idx" @click="switchTag(idx)">{{ tag.name }}</span>
        </div>
        </div>
    </div>
    <div class="columns">
        <div class="column">
            <textarea class="textarea" v-model="state.input" placeholder="请输入原始内容"></textarea>
            <div class="buttons is-centered mt-4">
                <button class="button is-primary" @click="Decode(opt.type)" v-for="(opt, idx) in tags[state.current].opts" :key="idx">{{ opt.name }}</button>
            </div>
            <textarea class="textarea" v-model="state.output" readonly placeholder="转换后的内容"></textarea>
        </div>
    </div>
    </div>
</div>
</template>

<script>
import { reactive, readonly } from 'vue'
import {Post} from "../../tools/http"

export default {
    name: "Coder",
    setup() {
        const state = reactive({
            current: 0,
            input: "",
            output: ""
        })
        const tags = readonly([
            {
                "name": "URLCoder",
                "opts": [
                    {
                        "name": "URLEncode",
                        "type": "urlencode"
                    },
                    {
                        "name": "URLDecode",
                        "type": "urldecode"
                    },
                ]
            },
            {
                "name": "Unicode",
                "opts": [
                    {
                        "name": "Unicode->中文",
                        "type": "unicodedecode"
                    },
                    {
                        "name": "中文->Unicode",
                        "type": "unicodeencode"
                    }
                ]
            },
            {
                "name": "Base64",
                "opts": [
                    {
                        "name": "Base64Encode",
                        "type": "base64encode"
                    },
                    {
                        "name": "Base64Decode",
                        "type": "base64decode"
                    }
                ]
            }
        ])

        async function Decode(type) {
            if (state.input == "") {
                state.output = "";
                return
            }

            let resp = await Post("/api/coder", {input: state.input, type: type})
            if (resp.data.code === 10000) {
                state.output = resp.data.data.output;
            }
        }

        function switchTag(idx) {
            state.current = idx;
        }

        return {
            state,
            tags,
            switchTag,
            Decode
        }
    }
}
</script>

<style lang="scss">
#coder {
    .tag {
        cursor: pointer;
    }
}
</style>
