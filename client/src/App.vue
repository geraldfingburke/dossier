<template>
  <div id="app">
    <header v-if="store.user">
      <div class="container">
        <h1>📰 Dossier</h1>
        <nav>
          <button @click="currentView = 'feeds'" :class="{ active: currentView === 'feeds' }">
            Feeds
          </button>
          <button @click="currentView = 'articles'" :class="{ active: currentView === 'articles' }">
            Articles
          </button>
          <button @click="currentView = 'digests'" :class="{ active: currentView === 'digests' }">
            Digests
          </button>
          <button @click="logout" class="logout">Logout</button>
        </nav>
      </div>
    </header>

    <main class="container">
      <LoginView v-if="!store.user" @login="handleLogin" />
      <FeedsView v-else-if="currentView === 'feeds'" />
      <ArticlesView v-else-if="currentView === 'articles'" />
      <DigestsView v-else-if="currentView === 'digests'" />
    </main>
  </div>
</template>

<script>
import { reactive, ref } from 'vue'
import LoginView from './views/LoginView.vue'
import FeedsView from './views/FeedsView.vue'
import ArticlesView from './views/ArticlesView.vue'
import DigestsView from './views/DigestsView.vue'
import { useStore } from './store'

export default {
  name: 'App',
  components: {
    LoginView,
    FeedsView,
    ArticlesView,
    DigestsView
  },
  setup() {
    const store = useStore()
    const currentView = ref('feeds')

    const handleLogin = () => {
      currentView.value = 'feeds'
    }

    const logout = () => {
      store.logout()
    }

    return {
      store,
      currentView,
      handleLogin,
      logout
    }
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
  background: #f5f5f5;
  color: #333;
  line-height: 1.6;
}

#app {
  min-height: 100vh;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

header {
  background: white;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  padding: 1rem 0;
  margin-bottom: 2rem;
}

header .container {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

h1 {
  font-size: 1.5rem;
  color: #2c3e50;
}

nav {
  display: flex;
  gap: 1rem;
}

button {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  background: #e0e0e0;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background 0.2s;
}

button:hover {
  background: #d0d0d0;
}

button.active {
  background: #4CAF50;
  color: white;
}

button.logout {
  background: #f44336;
  color: white;
}

button.logout:hover {
  background: #d32f2f;
}

button.primary {
  background: #4CAF50;
  color: white;
}

button.primary:hover {
  background: #45a049;
}

button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

input, textarea {
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  font-family: inherit;
  width: 100%;
}

input:focus, textarea:focus {
  outline: none;
  border-color: #4CAF50;
}

.card {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 1rem;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.error {
  color: #f44336;
  background: #ffebee;
  padding: 0.5rem;
  border-radius: 4px;
  margin-bottom: 1rem;
}

.success {
  color: #4CAF50;
  background: #e8f5e9;
  padding: 0.5rem;
  border-radius: 4px;
  margin-bottom: 1rem;
}
</style>
