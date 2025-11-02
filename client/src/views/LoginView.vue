<template>
  <div class="login-container">
    <div class="card login-card">
      <h2>{{ isLogin ? "Login" : "Register" }}</h2>

      <div v-if="error" class="error">{{ error }}</div>

      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label>Email</label>
          <input
            v-model="email"
            type="email"
            required
            placeholder="your@email.com"
          />
        </div>

        <div v-if="!isLogin" class="form-group">
          <label>Name</label>
          <input v-model="name" type="text" required placeholder="Your Name" />
        </div>

        <div class="form-group">
          <label>Password</label>
          <input
            v-model="password"
            type="password"
            required
            placeholder="Password"
          />
        </div>

        <button type="submit" class="primary" :disabled="loading">
          {{ loading ? "Loading..." : isLogin ? "Login" : "Register" }}
        </button>
      </form>

      <p class="toggle">
        {{ isLogin ? "Don't have an account?" : "Already have an account?" }}
        <a href="#" @click.prevent="isLogin = !isLogin">
          {{ isLogin ? "Register" : "Login" }}
        </a>
      </p>
    </div>
  </div>
</template>

<script>
import { ref } from "vue";
import { useStore } from "../store";

export default {
  name: "LoginView",
  emits: ["login"],
  setup(props, { emit }) {
    const store = useStore();
    const isLogin = ref(true);
    const email = ref("");
    const password = ref("");
    const name = ref("");
    const error = ref("");
    const loading = ref(false);

    const handleSubmit = async () => {
      error.value = "";
      loading.value = true;

      try {
        if (isLogin.value) {
          console.log("Attempting login...");
          await store.login(email.value, password.value);
          console.log("Login successful, user:", store.state.user);
        } else {
          console.log("Attempting registration...");
          await store.register(email.value, password.value, name.value);
          console.log("Registration successful, user:", store.state.user);
        }

        // Only emit login event if we actually have a user
        if (store.state.user) {
          console.log("Emitting login event");
          emit("login");
        } else {
          error.value = "Authentication failed - no user returned";
        }
      } catch (err) {
        console.error("Authentication error:", err);
        error.value = err.message || "An error occurred";
      } finally {
        loading.value = false;
      }
    };

    return {
      isLogin,
      email,
      password,
      name,
      error,
      loading,
      handleSubmit,
    };
  },
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
}

.login-card {
  max-width: 400px;
  width: 100%;
}

h2 {
  margin-bottom: 1.5rem;
  text-align: center;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
}

button[type="submit"] {
  width: 100%;
  padding: 0.75rem;
  margin-top: 1rem;
}

.toggle {
  text-align: center;
  margin-top: 1rem;
}

.toggle a {
  color: #4caf50;
  text-decoration: none;
}

.toggle a:hover {
  text-decoration: underline;
}
</style>
