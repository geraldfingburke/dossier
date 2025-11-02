<template>
  <div class="dossier-configs">
    <div class="header-section">
      <div class="header-content">
        <h2>Dossiers</h2>
        <p class="subtitle">Manage your dossier subscriptions</p>
      </div>
      <div class="header-buttons">
        <button
          @click="testEmailConnection"
          class="secondary"
          :disabled="loading"
        >
          Test Email
        </button>
        <button @click="showCreateForm = true" class="primary">
          <span>+ New Dossier</span>
        </button>
      </div>
    </div>

    <!-- Error/Success Messages -->
    <div v-if="error" class="error">{{ error }}</div>
    <div v-if="success" class="success">{{ success }}</div>

    <!-- Create/Edit Form Modal -->
    <div
      v-if="showCreateForm || editingConfig"
      class="modal-overlay"
      @click.self="closeForm"
    >
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ editingConfig ? "Edit" : "Create" }} Dossier Configuration</h3>
          <button @click="closeForm" class="close-btn">&times;</button>
        </div>

        <form @submit.prevent="saveConfig" class="config-form">
          <div class="form-group">
            <label>Title*</label>
            <input
              v-model="title"
              type="text"
              placeholder="e.g., Daily Tech News"
              :class="{ error: titleError }"
            />
            <div v-if="titleError" class="field-error">{{ titleError }}</div>
          </div>

          <div class="form-group">
            <label>Email Address*</label>
            <input
              v-model="emailField"
              type="email"
              placeholder="your@email.com"
              :class="{ error: emailError }"
            />
            <div v-if="emailError" class="field-error">{{ emailError }}</div>
          </div>

          <div class="form-group">
            <label>RSS Feed URLs*</label>
            <div class="feed-urls" :class="{ error: feedUrlsError }">
              <div
                v-for="(url, index) in feedUrls"
                :key="index"
                class="feed-url-row"
              >
                <input
                  v-model="feedUrls[index]"
                  type="url"
                  placeholder="https://example.com/feed.xml"
                  :class="{ error: feedUrlsError }"
                />
                <button
                  type="button"
                  @click="removeFeedUrl(index)"
                  class="danger remove-btn"
                  :disabled="feedUrls.length === 1"
                >
                  Remove
                </button>
              </div>
              <button type="button" @click="addFeedUrl" class="add-feed-btn">
                + Add Feed
              </button>
            </div>
            <div v-if="feedUrlsError" class="field-error">
              {{ feedUrlsError }}
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>Article Count (1-50)*</label>
              <input
                v-model.number="articleCount"
                type="number"
                min="1"
                max="50"
                :class="{ error: articleCountError }"
              />
              <div v-if="articleCountError" class="field-error">
                {{ articleCountError }}
              </div>
            </div>

            <div class="form-group">
              <label>Frequency*</label>
              <select v-model="frequency" :class="{ error: frequencyError }">
                <option value="daily">Daily</option>
                <option value="weekly">Weekly</option>
                <option value="monthly">Monthly</option>
              </select>
              <div v-if="frequencyError" class="field-error">
                {{ frequencyError }}
              </div>
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>Delivery Time*</label>
              <input
                v-model="deliveryTime"
                type="time"
                :class="{ error: deliveryTimeError }"
              />
              <div v-if="deliveryTimeError" class="field-error">
                {{ deliveryTimeError }}
              </div>
            </div>

            <div class="form-group">
              <label>Timezone</label>
              <select v-model="timezone">
                <option value="UTC">UTC</option>
                <option value="America/New_York">Eastern Time</option>
                <option value="America/Chicago">Central Time</option>
                <option value="America/Denver">Mountain Time</option>
                <option value="America/Los_Angeles">Pacific Time</option>
                <option value="Europe/London">London</option>
                <option value="Europe/Paris">Paris</option>
                <option value="Asia/Tokyo">Tokyo</option>
              </select>
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>Tone</label>
              <select v-model="tone">
                <option value="professional">Professional</option>
                <option value="casual">Casual</option>
                <option value="humorous">Humorous</option>
                <option value="analytical">Analytical</option>
                <option value="orc">Orc (Warcraft style)</option>
                <option value="sweary">Sweary (Explicit)</option>
                <option value="robot">Robot (Beep boop)</option>
                <option value="southern-belle">Southern Belle</option>
                <option value="doomsayer">Doomsayer (Apocalyptic)</option>
                <option value="apologetic">Apologetic</option>
              </select>
            </div>

            <div class="form-group">
              <label>Language</label>
              <select v-model="language">
                <option value="English">English</option>
                <option value="Spanish">Spanish</option>
                <option value="French">French</option>
                <option value="German">German</option>
                <option value="Italian">Italian</option>
                <option value="Portuguese">Portuguese</option>
              </select>
            </div>
          </div>

          <div class="form-group">
            <label>Special Instructions</label>
            <textarea
              v-model="specialInstructions"
              placeholder="Any specific instructions for the AI summary (optional)"
              rows="3"
            ></textarea>
          </div>

          <div class="form-actions">
            <button type="button" @click="closeForm">Cancel</button>
            <button type="submit" class="primary" :disabled="loading">
              {{ loading ? "Saving..." : editingConfig ? "Update" : "Create" }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Configs List -->
    <div class="configs-grid">
      <div v-if="loading && configs.length === 0" class="loading-state">
        <div class="loading">Loading configurations...</div>
      </div>

      <div v-else-if="configs.length === 0" class="empty-state">
        <div class="empty-icon">📧</div>
        <h3>No Dossiers</h3>
        <p>Create your first automated dossier to get started</p>
        <button @click="showCreateForm = true" class="primary">
          Create First Dossier
        </button>
      </div>

      <div v-else class="configs-list">
        <div v-for="config in configs" :key="config.id" class="config-card">
          <div class="config-header">
            <div>
              <h3>{{ config.title }}</h3>
              <p class="config-email">{{ config.email }}</p>
            </div>
            <div class="config-actions">
              <button
                @click="testSendDossier(config.id)"
                class="test-btn"
                :disabled="loading"
              >
                Test Send
              </button>
              <button @click="editConfig(config)" class="edit-btn">Edit</button>
              <button
                @click="deleteConfig(config.id)"
                class="danger delete-btn"
              >
                Delete
              </button>
            </div>
          </div>

          <div class="config-details">
            <div class="detail-row">
              <span class="label">Feeds:</span>
              <span>{{ config.feedUrls.length }} feed(s)</span>
            </div>
            <div class="detail-row">
              <span class="label">Frequency:</span>
              <span>{{ config.frequency }} at {{ config.deliveryTime }}</span>
            </div>
            <div class="detail-row">
              <span class="label">Articles:</span>
              <span>{{ config.articleCount }} per dossier</span>
            </div>
            <div class="detail-row">
              <span class="label">Style:</span>
              <span>{{ config.tone }} tone, {{ config.language }}</span>
            </div>
          </div>

          <div v-if="config.specialInstructions" class="special-instructions">
            <strong>Instructions:</strong> {{ config.specialInstructions }}
          </div>

          <div class="config-status">
            <span :class="['status-badge', { active: config.active }]">
              {{ config.active ? "Active" : "Inactive" }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from "vue";
import { useStore } from "../store";
import { useForm, useField } from "vee-validate";

export default {
  name: "DossierConfigsView",
  setup() {
    const store = useStore();
    const configs = ref([]);
    const loading = ref(false);
    const error = ref("");
    const success = ref("");
    const showCreateForm = ref(false);
    const editingConfig = ref(null);

    // Validation schema
    const validationSchema = {
      title: (value) => {
        if (!value) return "Title is required";
        if (value.length < 3) return "Title must be at least 3 characters";
        if (value.length > 100) return "Title must be less than 100 characters";
        return true;
      },
      email: (value) => {
        if (!value) return "Email is required";
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(value)) return "Must be a valid email address";
        return true;
      },
      feedUrls: (value) => {
        if (!value || value.length === 0)
          return "At least one feed URL is required";
        const nonEmptyUrls = value.filter((url) => url.trim() !== "");
        if (nonEmptyUrls.length === 0)
          return "At least one feed URL is required";

        for (const feedUrl of nonEmptyUrls) {
          try {
            new URL(feedUrl);
          } catch (e) {
            return `"${feedUrl}" is not a valid URL`;
          }
        }
        return true;
      },
      articleCount: (value) => {
        if (!value) return "Article count is required";
        if (value < 1) return "Article count must be at least 1";
        if (value > 50) return "Article count must be at most 50";
        return true;
      },
      frequency: (value) => {
        if (!value) return "Frequency is required";
        if (!["daily", "weekly", "monthly"].includes(value))
          return "Invalid frequency";
        return true;
      },
      deliveryTime: (value) => {
        if (!value) return "Delivery time is required";
        const timePattern = /^([01]?[0-9]|2[0-3]):[0-5][0-9]$/;
        if (!timePattern.test(value)) return "Time must be in HH:MM format";
        return true;
      },
    };

    // Initialize VeeValidate form
    const {
      handleSubmit,
      resetForm: resetVeeValidateForm,
      setFieldValue,
    } = useForm({
      validationSchema,
      initialValues: {
        title: "",
        email: "",
        feedUrls: [""],
        articleCount: 20,
        frequency: "daily",
        deliveryTime: "09:00",
        timezone: "UTC",
        tone: "professional",
        language: "English",
        specialInstructions: "",
      },
    });

    // Form fields with validation
    const { value: title, errorMessage: titleError } = useField("title");
    const { value: emailField, errorMessage: emailError } = useField("email");
    const { value: feedUrls, errorMessage: feedUrlsError } =
      useField("feedUrls");
    const { value: articleCount, errorMessage: articleCountError } =
      useField("articleCount");
    const { value: frequency, errorMessage: frequencyError } =
      useField("frequency");
    const { value: deliveryTime, errorMessage: deliveryTimeError } =
      useField("deliveryTime");
    const { value: timezone } = useField("timezone");
    const { value: tone } = useField("tone");
    const { value: language } = useField("language");
    const { value: specialInstructions } = useField("specialInstructions");

    const resetForm = () => {
      resetVeeValidateForm();
      setFieldValue("feedUrls", [""]);
    };

    const closeForm = () => {
      showCreateForm.value = false;
      editingConfig.value = null;
      resetForm();
    };

    const addFeedUrl = () => {
      const currentUrls = feedUrls.value || [];
      setFieldValue("feedUrls", [...currentUrls, ""]);
    };

    const removeFeedUrl = (index) => {
      const currentUrls = feedUrls.value || [];
      if (currentUrls.length > 1) {
        const newUrls = currentUrls.filter((_, i) => i !== index);
        setFieldValue("feedUrls", newUrls);
      }
    };

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

    const saveConfig = handleSubmit(async (values) => {
      try {
        loading.value = true;
        error.value = "";
        success.value = "";

        // Filter out empty feed URLs
        const cleanedFormData = {
          ...values,
          feedUrls: values.feedUrls.filter((url) => url.trim() !== ""),
        };

        if (editingConfig.value) {
          await store.updateDossierConfig(
            editingConfig.value.id,
            cleanedFormData
          );
          success.value = "Configuration updated successfully!";
        } else {
          await store.createDossierConfig(cleanedFormData);
          success.value = "Configuration created successfully!";
        }

        closeForm();
        loadConfigs();

        // Clear success message after 3 seconds
        setTimeout(() => {
          success.value = "";
        }, 3000);
      } catch (err) {
        error.value = editingConfig.value
          ? "Failed to update configuration"
          : "Failed to create configuration";
        console.error("Save config error:", err);
      } finally {
        loading.value = false;
      }
    });

    const editConfig = (config) => {
      editingConfig.value = config;
      // Populate VeeValidate form fields
      setFieldValue("title", config.title);
      setFieldValue("email", config.email);
      setFieldValue("feedUrls", config.feedUrls);
      setFieldValue("articleCount", config.articleCount);
      setFieldValue("frequency", config.frequency);
      setFieldValue("deliveryTime", config.deliveryTime);
      setFieldValue("timezone", config.timezone);
      setFieldValue("tone", config.tone);
      setFieldValue("language", config.language);
      setFieldValue("specialInstructions", config.specialInstructions);
      showCreateForm.value = true;
    };

    const deleteConfig = async (id) => {
      if (!confirm("Are you sure you want to delete this configuration?")) {
        return;
      }

      try {
        loading.value = true;
        await store.deleteDossierConfig(id);
        success.value = "Configuration deleted successfully!";
        loadConfigs();

        setTimeout(() => {
          success.value = "";
        }, 3000);
      } catch (err) {
        error.value = "Failed to delete configuration";
        console.error("Delete config error:", err);
      } finally {
        loading.value = false;
      }
    };

    const testSendDossier = async (configId) => {
      if (
        !confirm(
          "This will generate and send a test dossier to the configured email address. Continue?"
        )
      ) {
        return;
      }

      try {
        loading.value = true;
        error.value = "";
        success.value = "";

        await store.generateAndSendDossier(configId);
        success.value = "Test dossier sent successfully! Check your email.";

        setTimeout(() => {
          success.value = "";
        }, 5000);
      } catch (err) {
        error.value = "Failed to send test dossier: " + err.message;
        console.error("Test send error:", err);
      } finally {
        loading.value = false;
      }
    };

    const testEmailConnection = async () => {
      try {
        loading.value = true;
        error.value = "";
        success.value = "";

        await store.testEmailConnection();
        success.value = "Email connection test successful!";

        setTimeout(() => {
          success.value = "";
        }, 3000);
      } catch (err) {
        error.value = "Email connection test failed: " + err.message;
        console.error("Email test error:", err);
      } finally {
        loading.value = false;
      }
    };

    onMounted(() => {
      loadConfigs();
    });

    return {
      configs,
      loading,
      error,
      success,
      showCreateForm,
      editingConfig,
      closeForm,
      addFeedUrl,
      removeFeedUrl,
      saveConfig,
      editConfig,
      deleteConfig,
      // VeeValidate fields
      title,
      titleError,
      emailField,
      emailError,
      feedUrls,
      feedUrlsError,
      articleCount,
      articleCountError,
      frequency,
      frequencyError,
      deliveryTime,
      deliveryTimeError,
      timezone,
      tone,
      language,
      specialInstructions,
      testSendDossier,
      testEmailConnection,
    };
  },
};
</script>

<style scoped>
.dossier-configs {
  max-width: 100%;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 2rem;
  gap: 2rem;
}

.header-content h2 {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  background: linear-gradient(135deg, #e5e5e7, #a3a3a3);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.header-content .subtitle {
  color: rgba(229, 229, 231, 0.6);
  font-size: 1rem;
}

.header-buttons {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.header-buttons button {
  white-space: nowrap;
}

.header-buttons .secondary {
  background: rgba(255, 255, 255, 0.05);
  color: rgba(229, 229, 231, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.header-buttons .secondary:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.1);
  color: #e5e5e7;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(8px);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.modal-content {
  background: rgba(26, 26, 30, 0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 1rem;
  width: 100%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 2rem 2rem 1rem 2rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-header h3 {
  font-size: 1.5rem;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  font-size: 2rem;
  cursor: pointer;
  color: rgba(229, 229, 231, 0.6);
  padding: 0;
  line-height: 1;
}

.close-btn:hover {
  color: #e5e5e7;
  background: none;
  transform: none;
}

/* Form Styles */
.config-form {
  padding: 2rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #e5e5e7;
}

/* Validation Error Styles */
.field-error {
  color: #ef4444;
  font-size: 0.875rem;
  margin-top: 0.25rem;
  font-weight: 500;
}

input.error,
select.error,
textarea.error {
  border-color: #ef4444 !important;
  background-color: rgba(239, 68, 68, 0.1) !important;
}

.feed-urls.error {
  border-color: #ef4444 !important;
  background-color: rgba(239, 68, 68, 0.1) !important;
}

input.error:focus,
select.error:focus,
textarea.error:focus {
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.2) !important;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.feed-urls {
  border: 1px solid rgba(55, 65, 81, 0.6);
  border-radius: 0.75rem;
  padding: 1rem;
  background: rgba(17, 17, 19, 0.5);
}

.feed-url-row {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
  align-items: stretch;
}

.feed-url-row:last-child {
  margin-bottom: 0;
}

.remove-btn {
  flex-shrink: 0;
  padding: 0.875rem 1rem;
  font-size: 0.75rem;
}

.add-feed-btn {
  background: rgba(55, 65, 81, 0.6);
  border: 1px dashed rgba(255, 255, 255, 0.3);
  color: rgba(229, 229, 231, 0.8);
  width: 100%;
  margin-top: 0.75rem;
}

.add-feed-btn:hover {
  background: rgba(75, 85, 99, 0.7);
  border-color: rgba(255, 255, 255, 0.5);
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

/* Config Cards */
.configs-grid {
  min-height: 400px;
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  text-align: center;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
  opacity: 0.6;
}

.empty-state h3 {
  margin-bottom: 0.5rem;
  color: #e5e5e7;
}

.empty-state p {
  color: rgba(229, 229, 231, 0.6);
  margin-bottom: 2rem;
}

.configs-list {
  display: grid;
  gap: 1.5rem;
}

.config-card {
  background: rgba(26, 26, 30, 0.6);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 1rem;
  padding: 2rem;
  transition: all 0.2s ease;
}

.config-card:hover {
  border-color: rgba(255, 255, 255, 0.15);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.4);
  transform: translateY(-2px);
}

.config-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.5rem;
  gap: 1rem;
}

.config-header h3 {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 0.25rem;
}

.config-email {
  color: rgba(229, 229, 231, 0.6);
  font-size: 0.875rem;
}

.config-actions {
  display: flex;
  gap: 0.5rem;
  flex-shrink: 0;
}

.edit-btn,
.delete-btn,
.test-btn {
  padding: 0.5rem 1rem;
  font-size: 0.75rem;
}

.test-btn {
  background: rgba(59, 130, 246, 0.1);
  color: #60a5fa;
  border: 1px solid rgba(59, 130, 246, 0.3);
}

.test-btn:hover:not(:disabled) {
  background: rgba(59, 130, 246, 0.2);
  border-color: rgba(59, 130, 246, 0.5);
  color: #93c5fd;
}

.config-details {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 1rem;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  padding: 0.75rem;
  background: rgba(17, 17, 19, 0.4);
  border-radius: 0.5rem;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.detail-row .label {
  color: rgba(229, 229, 231, 0.6);
  font-weight: 500;
}

.special-instructions {
  padding: 1rem;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 0.5rem;
  margin-bottom: 1rem;
  font-size: 0.875rem;
}

.config-status {
  display: flex;
  justify-content: flex-end;
}

.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 500;
  background: rgba(55, 65, 81, 0.6);
  color: rgba(229, 229, 231, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.status-badge.active {
  background: rgba(16, 185, 129, 0.2);
  color: #86efac;
  border-color: rgba(16, 185, 129, 0.3);
}

/* Responsive Design */
@media (max-width: 768px) {
  .header-section {
    flex-direction: column;
    align-items: stretch;
    gap: 1rem;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .config-header {
    flex-direction: column;
    gap: 1rem;
  }

  .config-actions {
    align-self: flex-start;
  }

  .config-details {
    grid-template-columns: 1fr;
  }

  .modal-overlay {
    padding: 1rem;
  }

  .modal-content {
    max-height: 95vh;
  }
}
</style>
