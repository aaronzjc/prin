<template>
  <li class="font-mono text-sm">
    <div :class="{ bold: state.isFolder }" @click="toggle">
      <p class="chain-name">{{ state.item.name }} <span class="chain-raw">{{ state.item.raw }}</span></p>
      <p class="chain-format">{{ state.item.raw_format }}</p>
    </div>
    <ul v-show="state.isOpen" v-if="state.isFolder">
      <tree-item class="item" v-for="(child, index) in state.item.children" :key="index" :item="child"></tree-item>
    </ul>
  </li>
</template>

<script>
import { computed, reactive } from "vue"

export default {
  name: "TreeItem",
  props: ["item"],
  setup(props) {
    const state = reactive({
      isOpen: false,
      item: props.item,
      isFolder: computed(() => {
        return state.item && state.item.children && state.item.children.length > 0;
      }),
    });

    function toggle() {
      if (state.isFolder) {
        state.isOpen = !state.isOpen;
      }
    }

    return {
      state,
      toggle,
    };
  },
};
</script>
