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
      <DossierConfigsView />
    </main>
  </div>
</template>

<script>
import DossierConfigsView from "./views/DossierConfigsView.vue";

export default {
  name: "App",
  components: {
    DossierConfigsView,
  },
};
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen,
    Ubuntu, Cantarell, sans-serif;
  background: #0a0a0b;
  color: #e5e5e7;
  line-height: 1.6;
  overflow-x: hidden;
}

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

/* Main Content */
.main-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem;
}

/* Global Components */
button {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 0.75rem;
  background: rgba(55, 65, 81, 0.8);
  color: #e5e5e7;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
  transition: all 0.2s ease;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

button:hover {
  background: rgba(75, 85, 99, 0.9);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

button.primary {
  background: linear-gradient(135deg, #3b82f6, #1d4ed8);
  border: 1px solid rgba(59, 130, 246, 0.3);
}

button.primary:hover {
  background: linear-gradient(135deg, #2563eb, #1e40af);
  box-shadow: 0 8px 24px rgba(59, 130, 246, 0.3);
}

button.success {
  background: linear-gradient(135deg, #10b981, #059669);
  border: 1px solid rgba(16, 185, 129, 0.3);
}

button.success:hover {
  background: linear-gradient(135deg, #059669, #047857);
  box-shadow: 0 8px 24px rgba(16, 185, 129, 0.3);
}

button.danger {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  border: 1px solid rgba(239, 68, 68, 0.3);
}

button.danger:hover {
  background: linear-gradient(135deg, #dc2626, #b91c1c);
  box-shadow: 0 8px 24px rgba(239, 68, 68, 0.3);
}

button:disabled {
  background: rgba(55, 65, 81, 0.4);
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

input,
textarea,
select {
  padding: 0.875rem;
  border: 1px solid rgba(55, 65, 81, 0.6);
  border-radius: 0.75rem;
  background: rgba(17, 17, 19, 0.8);
  color: #e5e5e7;
  font-size: 0.875rem;
  font-family: inherit;
  width: 100%;
  backdrop-filter: blur(10px);
  transition: all 0.2s ease;
}

input:focus,
textarea:focus,
select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  background: rgba(17, 17, 19, 0.9);
}

input::placeholder,
textarea::placeholder {
  color: rgba(229, 229, 231, 0.5);
}

.card {
  background: rgba(26, 26, 30, 0.8);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 1rem;
  padding: 2rem;
  margin-bottom: 1.5rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

.error {
  color: #fca5a5;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.2);
  padding: 1rem;
  border-radius: 0.75rem;
  margin-bottom: 1rem;
  backdrop-filter: blur(10px);
}

.success {
  color: #86efac;
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.2);
  padding: 1rem;
  border-radius: 0.75rem;
  margin-bottom: 1rem;
  backdrop-filter: blur(10px);
}

/* Loading Animation */
.loading {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.loading::after {
  content: "";
  width: 1rem;
  height: 1rem;
  border: 2px solid rgba(59, 130, 246, 0.3);
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Utility Classes */
.text-sm {
  font-size: 0.875rem;
}

.text-xs {
  font-size: 0.75rem;
}

.font-medium {
  font-weight: 500;
}

.font-semibold {
  font-weight: 600;
}

.opacity-60 {
  opacity: 0.6;
}

.flex {
  display: flex;
}

.flex-col {
  flex-direction: column;
}

.items-center {
  align-items: center;
}

.justify-between {
  justify-content: space-between;
}

.gap-2 {
  gap: 0.5rem;
}

.gap-4 {
  gap: 1rem;
}

.mb-2 {
  margin-bottom: 0.5rem;
}

.mb-4 {
  margin-bottom: 1rem;
}

.mt-4 {
  margin-top: 1rem;
}
</style>
