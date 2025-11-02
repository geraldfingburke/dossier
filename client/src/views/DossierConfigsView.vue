<template>
  <div class="dossier-configs">
    <div class="header-section">
      <div class="header-content">
        <h2>Dossiers</h2>
        <p class="subtitle">Manage your dossier subscriptions</p>
      </div>
      <div class="header-buttons">
        <button @click="createNewDossier" class="primary">
          <span>+ New Dossier</span>
        </button>
      </div>
    </div>

    <!-- Error/Success Messages -->
    <div v-if="error" class="error">{{ error }}</div>
    <div v-if="success" class="success">{{ success }}</div>

    <!-- Dossier Tiles -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Loading dossiers...</p>
    </div>

    <div v-else-if="configs.length === 0" class="empty-state">
      <div class="empty-content">
        <h3>No dossiers yet</h3>
        <p>Create your first dossier to start receiving automated summaries</p>
        <button @click="createNewDossier" class="primary">
          Create Your First Dossier
        </button>
      </div>
    </div>

    <div v-else class="configs-list">
      <div
        v-for="config in configs"
        :key="config.id"
        class="dossier-tile"
        @click="viewDossier(config.id)"
      >
        <div class="tile-header">
          <h3>{{ config.title }}</h3>
          <span :class="['status-indicator', { active: config.active }]">
            {{ config.active ? "●" : "○" }}
          </span>
        </div>

        <div class="tile-summary">
          <div class="summary-item">
            <span class="icon">📡</span>
            <span>{{ config.feedUrls.length }} feeds</span>
          </div>
          <div class="summary-item">
            <span class="icon">📅</span>
            <span>{{ config.frequency }} at {{ config.deliveryTime }}</span>
          </div>
          <div class="summary-item">
            <span class="icon">🎯</span>
            <span>{{ config.tone }} tone</span>
          </div>
        </div>

        <div class="tile-footer">
          <span class="email-indicator">{{ config.email }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from "vue";
import { useStore } from "../store";
import { useRouter } from "vue-router";

export default {
  name: "DossierConfigsView",
  setup() {
    const store = useStore();
    const router = useRouter();
    const configs = ref([]);
    const loading = ref(false);
    const error = ref("");
    const success = ref("");

    /**
     * Load all dossier configurations from the store
     */
    const loadConfigs = async () => {
      try {
        loading.value = true;
        error.value = "";
        configs.value = await store.getDossierConfigs();
      } catch (err) {
        error.value = "Failed to load configurations";
        console.error("Load configs error:", err);
      } finally {
        loading.value = false;
      }
    };

    /**
     * Navigate to dossier detail view
     * @param {number} dossierId - The ID of the dossier to view
     */
    const viewDossier = (dossierId) => {
      router.push(`/dossier/${dossierId}`);
    };

    /**
     * Navigate to create new dossier view
     */
    const createNewDossier = () => {
      router.push("/dossier/new");
    };

    onMounted(() => {
      loadConfigs();
    });

    return {
      configs,
      loading,
      error,
      success,
      viewDossier,
      createNewDossier,
    };
  },
};
</script>

<style scoped>
/**
 * DossierConfigsView Component Styles
 * View-specific styles only - common patterns in global styles
 */

.dossier-configs {
  max-width: 100%;
}

/* Header Gradient - Specific to this view */
.header-content h2 {
  background: linear-gradient(135deg, #e5e5e7, #a3a3a3);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

/* Dossier Tiles - Unique to list view */
.configs-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: var(--spacing-lg);
}

.dossier-tile {
  background: linear-gradient(
    135deg,
    rgba(26, 26, 30, 0.8) 0%,
    rgba(18, 18, 22, 0.9) 100%
  );
  backdrop-filter: var(--blur-md);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  padding: var(--spacing-lg);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.dossier-tile::before {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, #3b82f6, #8b5cf6, #06b6d4);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.dossier-tile:hover {
  border-color: rgba(255, 255, 255, 0.2);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
  transform: translateY(-4px);
}

.dossier-tile:hover::before {
  opacity: 1;
}

.tile-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--spacing-md);
}

.tile-header h3 {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0;
  color: var(--color-text-primary);
  line-height: 1.3;
}

.status-indicator {
  font-size: 1.2rem;
  opacity: 0.7;
  transition: all 0.2s ease;
}

.status-indicator.active {
  color: var(--color-success);
  opacity: 1;
}

.tile-summary {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  margin-bottom: 1.25rem;
}

.summary-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  color: rgba(229, 229, 231, 0.8);
  font-size: 0.9rem;
}

.summary-item .icon {
  font-size: 1rem;
  opacity: 0.8;
}

.tile-footer {
  padding-top: var(--spacing-md);
  border-top: 1px solid var(--color-border-subtle);
}

.email-indicator {
  font-size: 0.85rem;
  color: var(--color-text-secondary);
  font-family: "SF Mono", "Monaco", "Cascadia Code", monospace;
}
</style>
