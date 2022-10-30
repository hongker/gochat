<script setup>
</script>

<template>
  <div class="limiter">
    <div class="container-login100">
      <div class="wrap-login100 p-l-55 p-r-55 p-t-65 p-b-54">
        <div class="login100-form validate-form">
					<span class="login100-form-title p-b-49">
						ImTalk
					</span>

          <div class="wrap-input100 validate-input m-b-23" data-validate = "Email is reauired">
            <span class="label-input100">Email</span>
            <input class="input100" type="email" name="email" v-model="user.email" placeholder="Type your Email">
            <span class="focus-input100" data-symbol="&#xf206;"></span>
          </div>
          <div class="wrap-input100 validate-input m-b-23" data-validate = "Username is reauired">
            <span class="label-input100">Username</span>
            <input class="input100" type="text" name="username" v-model="user.name" placeholder="Type your username">
            <span class="focus-input100" data-symbol="&#xf206;"></span>
          </div>

<!--          <div class="wrap-input100 validate-input" data-validate="Password is required">-->
<!--            <span class="label-input100">Password</span>-->
<!--            <input class="input100" type="password" name="pass" placeholder="Type your password">-->
<!--            <span class="focus-input100" data-symbol="&#xf190;"></span>-->
<!--          </div>-->

          <div class="text-right p-t-8 p-b-31">
            <a href="#">
              Forgot password?
            </a>
          </div>

          <div class="container-login100-form-btn">
            <div class="wrap-login100-form-btn">
              <div class="login100-form-bgbtn"></div>
              <button class="login100-form-btn" @click="submit" v-loading.fullscreen.lock="fullscreenLoading">
                Login
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
import {sendRequest} from "../utils/api";
import {userStore} from "../stores/counter";
import {ElMessage} from "element-plus";

export default {
  name: "Login",
  data() {
    return {
      fullscreenLoading: false,
      user : {
        name: "",
        email: "",
      }

    }
  },
  mounted() {
    console.log(import.meta.env)
  },
  methods: {
    submit() {
      this.fullscreenLoading = true
      var that = this;
      console.log("vue submit");


      sendRequest("/user/auth", this.user).then((res) => {

        if (res.code === 0) {
          let store = userStore();
          store.changeUid(res.data.uid);
          store.changeToken(res.data.token)

          setTimeout(function () {
            that.fullscreenLoading = false
            that.$router.push({path:'/'})
          }, 1000)
        }else {
          ElMessage.error(res.msg)
          that.fullscreenLoading = false
        }

      })





    }
  }
}
</script>

<style scoped>
.container-login100 {
  background-image: url("/static/images/bg-01.jpg");
}
</style>