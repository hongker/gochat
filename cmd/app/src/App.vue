<script setup>
import { RouterLink, RouterView } from 'vue-router'

import { onMounted, inject } from 'vue'
onMounted(() => {
  const socket = inject('socket')
  const packet = inject('packet')
  const operation = inject('operation')
  const ws = socket('')
  ws.onopen = () => {
    console.log('init websocket successfully')
  }

  var textDecoder = new TextDecoder();

  ws.onmessage = ({ data }) => {
    console.log(data)
    var dataView = new DataView(data, 0);
    var packetLen = dataView.getInt32(packet.packetOffset);
    var op = dataView.getInt16(packet.opOffset);
    var contentType = dataView.getInt16(packet.contentTypeOffset);
    var seq = dataView.getInt16(packet.seqOffset);
    var msgBody = textDecoder.decode(data.slice(packet.rawHeaderLen, packetLen));

    console.log("receiveHeader: packetLen=" + packetLen,  "op=" + op,  "contentType=" + contentType, "seq=" + seq, "msgBody=" + msgBody);

    switch (op) {
      case operation.login:
        console.log("login success");
        break;
      default:
        console.log("unknown operation")
    }
  }
})
</script>

<template>
  <RouterView />
</template>

<style scoped>

</style>
