<template>
<ul class="chains">
    <chain :chain="chain" v-for="(chain, idx) in props.chains" :key="idx"></chain>
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