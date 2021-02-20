<template>
  <div class="columns" id="jsontogo">
    <div class="column p-0">
      <div class="columns">
        <div class="column">
          <article class="message is-dark">
            <div class="message-body content is-small">
              <p>根据Json字符串，获取对应的Go结构体定义。基于<a href="https://github.com/mholt/json-to-go">mholt/json-to-go</a>包实现。</p>
            </div>
          </article>
        </div>
      </div>
      <div class="columns" id="code-block">
        <div class="column is-half">
          <textarea class="textarea" v-model="str" rows="20" placeholder="请输入有效的json字符串"></textarea>
        </div>
        <div class="column is0-half">
          <textarea class="textarea" v-model="code" rows="20" readonly></textarea>
        </div>
      </div>
      <div class="columns">
        <div class="column">
          <div id="go-panel">
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {JsonToGo} from "@/tools/jsontogo";
import {watch, ref, onMounted} from "vue";

export default {
  name: "JsonToGo.vue",
  setup() {
    const str = ref('{"name":"kate"}')
    const code = ref("")

    watch(str, () => {
      ToJson()
    })

    onMounted(() => {
      ToJson()
    })

    function ToJson() {
      console.log("change")
      code.value = ""
      if (str.value !== "") {
        let res = JsonToGo(str.value)
        if (res.go === "" && res.error !== "") {
          code.value = res.error
          return false
        }
        code.value = res.go
      }
    }

    return {
      str,
      code,
      ToJson
    }
  }
}
</script>

<style scoped>
#code-block textarea{
  font-size: 0.8rem;
  font-family: Consolas, Monaco, monospace;
}
</style>