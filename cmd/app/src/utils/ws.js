/**
 * @description 创建实例并
 * @param {*} topic topic
 * @returns websocket实例
 */
let client = null
const connectSocket = (topic) => {
    const baseUrl = import.meta.env.VITE_APP_WS_URL
    const wsUrl = `ws://127.0.0.1:8082`
    if (client) {
        console.log(client);
        return client
    } else {
        client = new WebSocket(wsUrl)
        client.binaryType = 'arraybuffer';
    }
    return client

}

export default connectSocket