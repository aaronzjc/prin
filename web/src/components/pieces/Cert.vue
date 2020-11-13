<template>
<div class="columns" id="cert">
    <div class="column">
        <textarea class="textarea" v-model="state.domains" placeholder="输入域名，一行一个"></textarea>
        <div class="buttons is-centered mt-4">
            <button class="button is-primary">获取证书</button>
        </div>
    </div>
</div>
</template>

<script>
import { reactive } from 'vue'
export default {
    name: "Cert",
    setup() {
        const state = reactive({
            domains: ""
        })

        return {
            state
        }
    }
}
</script>