<template>
<ul class="rules">
    <rule :rule="rule" v-for="(rule, idx) in props.rules" :key="idx"></rule>
</ul>
</template>

<script>
export default {
    name: "Rules",
    props: ['rules'],
    setup(props) {
        return {
            props
        }
    }
}
</script>