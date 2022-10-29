import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { storeToRefs } from 'pinia'


export const useCounterStore = defineStore('counter', () => {
  const count = ref(0)
  const doubleCount = computed(() => count.value * 2)
  function increment() {
    count.value++
  }

  return { count, doubleCount, increment }
})

export const userStore = defineStore('user', {
  state: () => ({
    uid:'',
    token : '',
  }),
  actions: {
    changeUid(id) {
      this.uid = id
    },
    changeToken(t) {
      this.token = t
    }
  },
  persist: {
    enabled: true,
  }
})

export const profileStore = defineStore('profile', () => {
  let name = ''
  function setName(str) {
    name = str
  }

  return {name, setName}
})