<script setup>
import { RouterLink, RouterView } from 'vue-router'

import { onMounted, inject } from 'vue'
onMounted(() => {
  const socket = inject('socket')
  const ws = socket('')
  ws.onopen = () => {
    console.log('init websocket successfully')
  }

  var textDecoder = new TextDecoder();
  const rawHeaderLen = 10;
  const packetOffset = 0;
  const opOffset = 4;
  const contentTypeOffset = 6;
  const seqOffset = 8;
  ws.onmessage = ({ data }) => {
    console.log(data)
    var dataView = new DataView(data, 0);
    var packetLen = dataView.getInt32(packetOffset);
    var op = dataView.getInt16(opOffset);
    var contentType = dataView.getInt16(contentTypeOffset);
    var seq = dataView.getInt16(seqOffset);

    console.log("receiveHeader: packetLen=" + packetLen,  "op=" + op,  "contentType=" + contentType, "seq=" + seq);

    var msgBody = textDecoder.decode(data.slice(rawHeaderLen, packetLen));
    console.log(msgBody)
  }
})
</script>

<template>
  <RouterView />
</template>

<style scoped>

</style>
