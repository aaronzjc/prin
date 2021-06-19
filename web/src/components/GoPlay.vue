<template>
    <div class="columns" id="jsontogo">
        <div class="column p-0">
            <div class="columns">
                <div class="column">
                    <article class="message is-dark">
                        <div class="message-body content is-small">
                            <p>
                                一个简单的Go
                                Playground工具。和官方Playground类似，用于执行Go代码片段。
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
                    <textarea
                        class="textarea"
                        v-model="state.code"
                        rows="20"
                        placeholder="随便写一点Go代码"
                    ></textarea>
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
import { reactive } from "vue";
import { PostOuterApi } from "../tools/http";

export default {
    name: "GoPlay",
    setup() {
        const state = reactive({
            code: `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, playground")
}
`,
            running: false,
        });
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
			if (state.code == "") {
				return false
			}
            state.running = true;
            let resp = await PostOuterApi(
                "http://49.51.163.143:7980",
                "/compile",
                qs.stringify({
                    body: state.code,
                    version: "2",
                    withVet: false,
                })
            );
            state.running = false;
            console.log(resp);
            playback(PlaygroundOutput(output), resp.data);
        }
        return {
            state,
            run,
        };
    },
};
</script>

<style scoped>
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
