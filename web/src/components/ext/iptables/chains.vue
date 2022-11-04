<template>
<ul class="chains">
    <li class="chain-item" v-for="(chain, idx) in props.chains" :key="idx">
        <div :class="[ 'chain-info', { 'empty' : chain.rules.length == 0 } ]">
            <span class="chain-name text-blue">{{ chain.name }}</span>
            <span class="chain-policy">policy: <span class="text-red">{{ chain.policy ? chain.policy : '-' }}</span></span>
        </div>
        <rules :rules="chain.rules"></rules>
    </li>
</ul>
</template>

<script>
export default {
    name: "Chains",
    props: ['chains'],
    setup(props) {
        return {
            props
        }
    }
}
</script>