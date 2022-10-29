<script setup>
import { inject } from 'vue'
const socket = inject('socket')
const ws = socket()
</script>

<template>


  
  <div class="limiter" id="vue">
    <div class="container-login100" style="background-image: url('src/static/images/bg-01.jpg');">
      <div class="wrap-login100 p-l-55 p-r-55 p-t-65 p-b-54">
        <div class="login100-form validate-form">
					<span class="login100-form-title p-b-49">
						Login
					</span>

          <div class="wrap-input100 validate-input m-b-23" data-validate = "Username is reauired">
            <span class="label-input100">Username</span>
            <input class="input100" type="text" name="username" v-model="username" placeholder="Type your username">
            <span class="focus-input100" data-symbol="&#xf206;"></span>
          </div>

          <div class="wrap-input100 validate-input" data-validate="Password is required">
            <span class="label-input100">Password</span>
            <input class="input100" type="password" name="pass" placeholder="Type your password">
            <span class="focus-input100" data-symbol="&#xf190;"></span>
          </div>

          <div class="text-right p-t-8 p-b-31">
            <a href="#">
              Forgot password?
            </a>
          </div>

          <div class="container-login100-form-btn">
            <div class="wrap-login100-form-btn">
              <div class="login100-form-bgbtn"></div>
              <button class="login100-form-btn" @click="submit">
                登录
              </button>
            </div>
          </div>

          <div class="txt1 text-center p-t-54 p-b-20">
						<span>
							Or Sign Up Using
						</span>
          </div>



          <div class="flex-col-c p-t-155">
						<span class="txt1 p-b-17">
							Or Sign Up Using
						</span>

            <a href="#" class="txt2">
              Sign Up
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>

</template>
<!--===============================================================================================-->


<script>
export default {
  name: "Login",
  data() {
    return {
      username: "admin",
    }
  },
  methods: {
    test() {
      console.log("Hello")
    },
    mergeArrayBuffer(ab1, ab2) {
      var u81 = new Uint8Array(ab1),
          u82 = new Uint8Array(ab2),
          res = new Uint8Array(ab1.byteLength + ab2.byteLength);
      res.set(u81, 0);
      res.set(u82, ab1.byteLength);
      return res.buffer;
    },
    submit() {
      console.log("vue submit")

      var that = this;

      const rawHeaderLen = 10;
      const packetOffset = 0;
      const opOffset = 4;
      const contentTypeOffset = 6;
      const seqOffset = 8;
      var textEncoder = new TextEncoder();


      var token = '{"name":"foo"}'
      var headerBuf = new ArrayBuffer(rawHeaderLen);
      var headerView = new DataView(headerBuf, 0);
      var bodyBuf = textEncoder.encode(token);
      headerView.setInt32(packetOffset, rawHeaderLen + bodyBuf.byteLength);
      headerView.setInt16(opOffset, 2);
      headerView.setInt16(contentTypeOffset, 1);
      headerView.setInt16(seqOffset, 1);
      var buf = that.mergeArrayBuffer(headerBuf, bodyBuf)
      that.ws.send(buf);
    }
  }
}
</script>

<style scoped>

</style>