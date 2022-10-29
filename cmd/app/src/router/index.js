import { createRouter, createWebHistory } from 'vue-router'
import {userStore} from "../stores/counter";


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/Home.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/Login.vue')
    }
  ]
})

router.beforeEach((to, from, next) => {
  let user = userStore()
  if (to.path !== "/login" && user.uid === "") {
    next("/login")
  }else {
    next()
  }


})
export default router
