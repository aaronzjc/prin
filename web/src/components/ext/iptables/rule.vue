<template>
<li>
    <div class="rule">
        <p class="rule-comment text-gray">{{ props.rule.comment }}</p>
        <p class="rule-text text-black">{{ props.rule.text }}</p>                            
        <p class="rule-match text-orange">
            <span v-for="(match, ii) in props.rule.matches" :key="ii">{{ match }}</span>
        </p>
        <p v-if="props.rule.chain" @click="toggle" class="rule-target text-blue rule-ref">{{ props.rule.target }}</p>
        <p v-else class="rule-target text-purple">{{ props.rule.target }}</p>
    </div>
    <chains :chains="props.rule.chain" v-if="state.open"></chains>
</li>
</template>

<script>
import { reactive, watch } from 'vue';

export default {
    name: "Rule",
    props: ['rule'],
    setup(props) {
        const state = reactive({
            open: false,
        })

        watch(props.rule, () => {
            state.open = false
        })

        function toggle() {
            state.open = !state.open
        }

        return {
            state,
            props,
            toggle
        }
    }
}
</script>