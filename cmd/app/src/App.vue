<script setup>
import { RouterLink, RouterView } from 'vue-router'
import HelloWorld from './components/HelloWorld.vue'

import { onMounted, inject } from 'vue'
onMounted(() => {
  const socket = inject('socket')
  const ws = socket('')
  ws.onopen = () => {
    console.log('init websocket successfully')
  }
  ws.onmessage = ({ data }) => {
    console.log(data)
  }
})
</script>

<template>
  <RouterView />
</template>

<style scoped>

</style>
