<template>
  <div class="auth">
    <form @submit.prevent="submit">
      <input
        type="text"
        placeholder="username"
        v-model="username"
      />
      <input
        type="password"
        placeholder="password"
        v-model="password"
      />
      <button
        :disabled="!valid">{{mode === 'signin' ? 'Login' : 'Register'}}</button>
      <a
        href="#"
        @click.prevent="toggleMode"
      >
        {{mode == 'signin' ? 'Register' : 'Login'}}
      </a>
    </form>
  </div>
</template>

<script>
import { computed, ref } from 'vue'
export default {
  name: 'Auth',
  emits: ['token'],
  setup(props, { emit }) {
    const username = ref('')
    const password = ref('')
    const mode = ref('signin')
    return {
      username,
      password,
      mode,
      toggleMode: () => {
        if (mode.value === 'signin') mode.value = 'signup'
        else mode.value = 'signin'
      },
      valid: computed(() => username.value.trim() !== '' && password.value !== ''),
      submit: async () => {
        try {
          const resp = await fetch(`http://localhost:8080/${mode.value}`,
            {
              method: 'POST',
              body: JSON.stringify({
                username: username.value.trim(),
                password: password.value
              })
            })
          if (resp.status !== 200) {
            throw 'expected 200'
          }
          const parsed = await resp.json()
          emit('token', parsed.token)
        } catch {
          alert(mode.value + ' error')
        }
      }
    }
  }
}
</script>

<style scoped>
.auth {
  box-sizing: border-box;
  width: 100%;
  font-family: sans-serif;
}
a {
    font-size: 12px;
}
input {
  box-sizing: border-box;
  margin: 10px 0;
  width: 100%;
  padding: 2px 5px;
  display: block;
}
button {
  box-sizing: border-box;
  width: 100%;
  display: block;
  margin: auto;
}
</style>>