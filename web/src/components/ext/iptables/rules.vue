<template>
<ul class="rules">
    <li v-for="(rule, idx) in props.rules" :key="idx">
        <div class="rule">
            <p class="rule-comment text-gray">{{ rule.comment }}</p>
            <p class="rule-text text-black">{{ rule.text }}</p>
            <p class="rule-match" v-if="rule.matches.length > 0">
                <span class="text-gray">Matches: </span>
                <span class="text-orange" v-for="(match, ii) in rule.matches" :key="ii">{{ match }}</span>
            </p>
            <p v-if="rule.chain" @click="toggle(idx)" class="rule-target">
                <span class="text-gray">Target: </span>
                <span class=" text-blue rule-ref">{{ rule.target }}</span>
            </p>
            <p v-else class="rule-target">
                <span class="text-gray">Target: </span>
                <span class="text-purple">{{ rule.target }}</span>
            </p>
        </div>
        <chains :chains="rule.chain" v-if="state.opens[idx]"></chains>
    </li>
</ul>
</template>

<script>
import { reactive } from 'vue';

export default {
    name: "Rules",
    props: ['rules'],
    setup(props) {
        const state = reactive({
            opens: new Array(props.rules.length)
        })
        function toggle(idx) {
            state.opens[idx] = !state.opens[idx]
        }
        return {
            state,
            props,
            toggle
        }
    }
}
</script>