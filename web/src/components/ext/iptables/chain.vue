<template>
<li class="chain-item">
    <div :class="[ 'chain-info', { 'empty' : isEmpty } ]">
        <span class="chain-name text-blue">{{ props.chain.name }}</span>
        <span class="chain-policy">policy: <span class="text-red">{{ props.chain.policy ? props.chain.policy : '-' }}</span></span>
    </div>
    <rules :rules="props.chain.rules"></rules>
</li>
</template>

<script>
import { computed } from 'vue';

export default {
   name: "Chain",
   props: ['chain'],
   setup(props) {
        const isEmpty = computed(() => {
            return props.chain.rules.length == 0
        })
       return {
           isEmpty,
           props
       }
   }
}
</script>