<script setup>
import { inject } from 'vue'
const socket = inject('socket')
const ws = socket()
</script>

<template>

  <link rel="icon" type="image/png" href="src/static/images/icons/favicon.ico"/>
  <!--===============================================================================================-->
  <link rel="stylesheet" type="text/css" href="src/static/vendor/bootstrap/css/bootstrap.min.css">
  <!--===============================================================================================-->
  <link rel="stylesheet" type="text/css" href="src/static/fonts/font-awesome-4.7.0/css/font-awesome.min.css">
  <!--===============================================================================================-->
  <link rel="stylesheet" type="text/css" href="src/static/fonts/iconic/css/material-design-iconic-font.min.css">
  <!--===============================================================================================-->
  <link rel="stylesheet" type="text/css" href="src/static/vendor/animate/animate.css">
  <!--===============================================================================================-->
  <link rel="stylesheet" type="text/css" href="src/static/vendor/css-hamburgers/hamburgers.min.css">
  <!--===============================================================================================-->
  <link rel="stylesheet" type="text/css" href="src/static/vendor/animsition/css/animsition.min.css">
  <!--===============================================================================================-->
  <link rel="stylesheet" type="text/css" href="src/static/vendor/select2/select2.min.css">
  <!--===============================================================================================-->
  <!--	<link rel="stylesheet" type="text/css" href="src/static/vendor/daterangepicker/daterangepicker.css">-->
  <!--===============================================================================================-->
  <link rel="stylesheet" type="text/css" href="src/static/css/util.css">
  <link rel="stylesheet" type="text/css" href="src/static/css/main.css">
  
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
    }
  }
}
</script>

<style scoped>

</style>