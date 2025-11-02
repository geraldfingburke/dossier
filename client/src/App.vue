<template>
  <div id="app">
    <header>
      <div class="header-container">
        <div class="logo">
          <h1>📰 Dossier</h1>
          <span class="subtitle">An AI RSS Digest</span>
        </div>
      </div>
    </header>

    <main class="main-container">
      <router-view />
    </main>
  </div>
</template>

<script>
export default {
  name: "App",
};
</script>

<style scoped>
/**
 * App Component Styles
 * Scoped styles specific to the App layout structure
 */

#app {
  min-height: 100vh;
  background: linear-gradient(135deg, #0a0a0b 0%, #1a1a1e 100%);
}

/* Header Styles */
header {
  background: rgba(26, 26, 30, 0.8);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  padding: 1rem 0;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 2rem;
}

.logo {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.logo h1 {
  font-size: 2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #60a5fa, #34d399);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.subtitle {
  font-size: 0.875rem;
  color: rgba(229, 229, 231, 0.6);
  font-weight: 400;
}

/* Main Content Container */
.main-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem;
}
</style>
