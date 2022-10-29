<script setup>

</script>

<template>
  <div class="limiter">
    <div class="container-login100">
      <div class="wrap-login100 p-l-55 p-r-55 p-t-65 p-b-54">
        <div class="login100-form validate-form">
					<span class="login100-form-title p-b-49">
						Login
					</span>

          <div class="wrap-input100 validate-input m-b-23" data-validate = "Username is reauired">
            <span class="label-input100">Username</span>
            <input class="input100" type="text" name="username" v-model="user.name" placeholder="Type your username">
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
      user : {
        name: "admin",
      }

    }
  },
  inject: ["socket", "packet", "operation"],
  mounted() {
    console.log(this.packet)
  },
  methods: {

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

      var textEncoder = new TextEncoder();
      var headerBuf = new ArrayBuffer(this.packet.rawHeaderLen);
      var headerView = new DataView(headerBuf, 0);
      var bodyBuf = textEncoder.encode(JSON.stringify(this.user));
      headerView.setInt32(this.packet.packetOffset, this.packet.rawHeaderLen + bodyBuf.byteLength);
      headerView.setInt16(this.packet.opOffset, this.operation.login);
      headerView.setInt16(this.packet.contentTypeOffset, 1);
      headerView.setInt16(this.packet.seqOffset, 1);
      var buf = this.mergeArrayBuffer(headerBuf, bodyBuf)

      var ws = this.socket()
      ws.send(buf);

      this.$router.push({path:'/'})
    }
  }
}
</script>

<style scoped>
.container-login100 {
  background-image: url("@/assets/static/images/bg-01.jpg");
}
</style>