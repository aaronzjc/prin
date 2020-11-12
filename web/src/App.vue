<template>
<section class="section" id="app">
    <Main></Main>
</section>
</template>

<script>
import Main from "./components/Main"

export default {
  name: 'App',
  components: {
    Main
  }
}
</script>

<style lang="scss">
html {
  background-color: white;
  font-size: 16px;
  -moz-osx-font-smoothing: grayscale;
  -webkit-font-smoothing: antialiased;
  min-width: 300px;
  overflow-x: hidden;
  overflow-y: scroll;
  text-rendering: optimizeLegibility;
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}
section.section {
  margin: 0 0;
  padding: 0 0;
  .columns {
    margin: 0;
  }
}
</style>