<template>
    <div class="columns" id="jsontogo">
        <div class="column p-0">
            <div class="columns">
                <div class="column">
                    <article class="message is-dark">
                        <div class="message-body content is-small">
                            <p>
                                一个简单的Go代码执行工具，基于官方Playground移植。
                            </p>
                        </div>
                    </article>
                </div>
            </div>
            <div class="columns">
                <div class="column pt-0 pb-0">
                    <div class="buttons are-small">
                        <button
                            :class="[
                                'button',
                                'is-info',
                                { 'is-loading': state.running },
                            ]"
                            @click="run"
                        >
                            执行
                        </button>
                    </div>
                </div>
            </div>
            <div class="columns" id="code-block">
                <div class="column is-half">
                    <div class="ace-editor" ref="aceRef"></div>
                </div>
                <div class="column is-half">
                    <div id="output"><pre></pre></div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import qs from "qs";
import { onMounted, reactive, ref } from "vue";
import { PostOuterApi } from "../tools/http";
import ace from 'ace-builds'
import 'ace-builds/webpack-resolver' 
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/mode-golang'

export default {
    name: "GoPlay",
    setup() {
        const state = reactive({
            aceEditor: null,
            running: false,
        });
        const aceRef = ref(null)

        function playback(output, data) {
            // Backwards compatibility: default values do not affect the output.
            var events = data.Events || [];
            var errors = data.Errors || "";
            var status = data.Status || 0;
            var isTest = data.IsTest || false;
            var testsFailed = data.TestsFailed || 0;

			output.innerHTML = "<pre></pre>"

            var timeout;
            output({ Kind: "start" });
            function next() {
                if (!events || events.length === 0) {
                    if (isTest) {
                        if (testsFailed > 0) {
                            output({
                                Kind: "system",
                                Body:
                                    "\n" +
                                    testsFailed +
                                    " test" +
                                    (testsFailed > 1 ? "s" : "") +
                                    " failed.",
                            });
                        } else {
                            output({ Kind: "system", Body: "\nAll tests passed." });
                        }
                    } else {
                        if (status > 0) {
                            output({ Kind: "end", Body: "status " + status + "." });
                        } else {
                            if (errors !== "") {
                                // errors are displayed only in the case of timeout.
                                output({ Kind: "end", Body: errors + "." });
                            } else {
                                output({ Kind: "end" });
                            }
                        }
                    }
                    return;
                }
                var e = events.shift();
                if (e.Delay === 0) {
                    output({ Kind: e.Kind, Body: e.Message });
                    next();
                    return;
                }
                timeout = setTimeout(function () {
                    output({ Kind: e.Kind, Body: e.Message });
                    next();
                }, e.Delay / 1000000);
            }
            next();
            return {
                Stop: function () {
                    clearTimeout(timeout);
                },
            };
        }

        function PlaygroundOutput(el) {
            "use strict";

            return function (write) {
                if (write.Kind == "start") {
                    el.innerHTML = "";
                    return;
                }

                var cl = "system";
                if (write.Kind == "stdout" || write.Kind == "stderr") cl = write.Kind;

                var m = write.Body;
                if (write.Kind == "end") {
                    m = "\nProgram exited" + (m ? ": " + m : ".");
                }

                if (m.indexOf("IMAGE:") === 0) {
                    // TODO(adg): buffer all writes before creating image
                    var url = "data:image/png;base64," + m.substr(6);
                    var img = document.createElement("img");
                    img.src = url;
                    el.appendChild(img);
                    return;
                }

                // ^L clears the screen.
                var s = m.split("\x0c");
                if (s.length > 1) {
                    el.innerHTML = "";
                    m = s.pop();
                }

                m = m.replace(/&/g, "&amp;");
                m = m.replace(/</g, "&lt;");
                m = m.replace(/>/g, "&gt;");

                var needScroll = el.scrollTop + el.offsetHeight == el.scrollHeight;

                var span = document.createElement("span");
                span.className = cl;
                span.innerHTML = m;
                el.appendChild(span);

                if (needScroll) el.scrollTop = el.scrollHeight - el.offsetHeight;
            };
        }

        async function run() {
            if (state.running) {
                return false;
            }
			var output = document.getElementById("output")
			output.innerHTML = ""
			output = output.appendChild(document.createElement("pre"))
            let code = state.editor.session.getValue()
			if (code == "") {
				return false
			}
            state.running = true;
            let resp = await PostOuterApi(
                "https://goplay.memosa.cn",
                "/compile",
                qs.stringify({
                    body: code,
                    version: "2",
                    withVet: false,
                })
            );
            state.running = false;
            playback(PlaygroundOutput(output), resp.data);
        }

        const tpl = `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, playground")
    fmt.Println("By Ace Editor")
}`
        
        onMounted(() => {
            if (state.editor !== undefined) {
                return
            }
            state.editor = ace.edit(aceRef.value,{
                minLines: 40, // 最小行数，还未到最大行数时，编辑器会自动伸缩大小
                fontSize: 14, // 编辑器内字体大小
                themePath: "ace/theme/chrome",
                mode: "ace/mode/golang", // 默认设置的语言模式
                tabSize: 4, // 制表符设置为 4 个空格大小
                fontFamily: "Consolas, Monaco",
                highlightActiveLine: false,
            })
            state.editor.session.setValue(tpl)
        })

        return {
            state,
            aceRef,
            run,
        };
    },
};
</script>

<style scoped>
.ace-editor {
    border: 1px solid #ccc;
    height: 450px;
}
.ace-editor /deep/ .ace_print-margin{ 
    display:none;
}
.ace-editor /deep/ .ace_gutter {
    background: #fff !important;
}
.ace-editor /deep/ .ace_gutter-active-line {
    background: #fff !important;
}
#code-block textarea {
    font-size: 0.8rem;
    font-family: Consolas, Monaco, monospace;
}

#output {
    /* The default monospace font on OS X is ugly, so specify Menlo
	 * instead. On other systems the default monospace font will be used. */
    font-family: Menlo, monospace;
    font-size: 11pt;
    padding: 4px;
}
</style>
