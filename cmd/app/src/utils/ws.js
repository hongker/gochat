import {userStore} from "../stores/counter";

/**
 * @description 创建实例并
 * @param {*} topic topic
 * @returns websocket实例
 */
let client = null
export const connectSocket = (wsUrl) => {
    if (client) {
        console.log(client);
        return client
    } else {
        client = new WebSocket(wsUrl)
        client.binaryType = 'arraybuffer';
    }
    return client

}

