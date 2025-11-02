<template>
  <div class="dossier-detail">
    <div class="header-section">
      <div class="header-content">
        <button @click="goBack" class="back-btn">← Back to Dossiers</button>
        <h2>{{ isCreating ? "Create New Dossier" : "Edit Dossier" }}</h2>
        <p class="subtitle">
          {{
            isCreating
              ? "Configure your new dossier"
              : "Modify dossier settings"
          }}
        </p>
      </div>
      <div class="header-buttons">
        <button
          @click="saveDossier"
          :disabled="(!hasChanges && !isCreating) || loading || !hasFormData"
          class="primary save-btn"
        >
          {{ isCreating ? "Create Dossier" : "Save Changes" }}
        </button>
        <button v-if="!isCreating" @click="openArchive" class="secondary">
          Archive
        </button>
        <button
          @click="testEmail"
          :disabled="loading || isCreating"
          class="secondary"
        >
          Send Test Dossier
        </button>
        <button @click="testGenerate" :disabled="loading" class="secondary">
          Test Generate
        </button>
        <button v-if="!isCreating" @click="confirmDelete" class="danger">
          Delete
        </button>
      </div>
    </div>

    <!-- Error/Success Messages -->
    <div class="message-container">
      <div v-if="error" class="error">{{ error }}</div>
      <div v-if="success" class="success">{{ success }}</div>
    </div>

    <!-- Loading State -->
    <div v-if="loading && !isCreating" class="loading-state">
      <div class="loading-spinner"></div>
      <p>Loading dossier data...</p>
    </div>

    <div v-else class="detail-content">
      <div class="main-form">
        <!-- Basic Settings -->
        <div class="form-section">
          <h3>Basic Settings</h3>

          <div class="form-group">
            <label>Dossier Title*</label>
            <input
              v-model="formData.title"
              type="text"
              placeholder="e.g., Daily Tech News"
              @input="markChanged"
            />
          </div>

          <div class="form-group">
            <label>Email Addresses*</label>
            <div class="email-list">
              <div
                v-for="(email, index) in formData.emails"
                :key="index"
                class="email-row"
              >
                <input
                  v-model="formData.emails[index]"
                  type="email"
                  placeholder="recipient@email.com"
                  @input="markChanged"
                />
                <button
                  type="button"
                  @click="removeEmail(index)"
                  class="danger remove-btn"
                  :disabled="formData.emails.length === 1"
                >
                  Remove
                </button>
              </div>
              <button type="button" @click="addEmail" class="add-email-btn">
                + Add Email
              </button>
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>Article Count</label>
              <input
                v-model.number="formData.articleCount"
                type="number"
                min="1"
                max="50"
                @input="markChanged"
              />
            </div>

            <div class="form-group">
              <label>Active</label>
              <label class="toggle-switch">
                <input
                  v-model="formData.active"
                  type="checkbox"
                  @change="markChanged"
                />
                <span class="toggle-slider"></span>
              </label>
            </div>
          </div>
        </div>

        <!-- Schedule Settings -->
        <div class="form-section">
          <h3>Schedule</h3>

          <div class="form-row">
            <div class="form-group">
              <label>Frequency*</label>
              <select v-model="formData.frequency" @change="markChanged">
                <option value="daily">Daily</option>
                <option value="weekly">Weekly</option>
                <option value="monthly">Monthly</option>
              </select>
            </div>

            <div class="form-group">
              <label>Delivery Time*</label>
              <input
                v-model="formData.deliveryTime"
                type="time"
                @input="markChanged"
              />
            </div>

            <div class="form-group">
              <label>Timezone</label>
              <select v-model="formData.timezone" @change="markChanged">
                <option value="America/New_York">Eastern Time</option>
                <option value="America/Chicago">Central Time</option>
                <option value="America/Denver">Mountain Time</option>
                <option value="America/Los_Angeles">Pacific Time</option>
                <option value="UTC">UTC</option>
                <option value="Europe/London">London</option>
                <option value="Europe/Paris">Paris</option>
                <option value="Asia/Tokyo">Tokyo</option>
              </select>
            </div>
          </div>
        </div>

        <!-- Content Settings -->
        <div class="form-section">
          <h3>Content Settings</h3>

          <div class="form-row">
            <div class="form-group">
              <label>AI Tone*</label>
              <div class="tone-selection">
                <select v-model="formData.tone" @change="markChanged">
                  <option
                    v-for="tone in availableTones"
                    :key="tone.id"
                    :value="tone.name"
                  >
                    {{ tone.name }} {{ tone.isSystemDefault ? "" : "(Custom)" }}
                  </option>
                </select>
                <button
                  type="button"
                  @click="openCustomToneModal"
                  class="secondary add-tone-btn"
                >
                  + Custom Tone
                </button>
              </div>
            </div>

            <div class="form-group">
              <label>Language</label>
              <select v-model="formData.language" @change="markChanged">
                <option value="english">English</option>
                <option value="spanish">Spanish</option>
                <option value="french">French</option>
                <option value="german">German</option>
                <option value="italian">Italian</option>
                <option value="portuguese">Portuguese</option>
                <option value="japanese">Japanese</option>
                <option value="chinese">Chinese</option>
              </select>
            </div>
          </div>

          <div class="form-group">
            <label>Special Instructions</label>
            <textarea
              v-model="formData.specialInstructions"
              placeholder="Additional instructions for the AI..."
              @input="markChanged"
            ></textarea>
          </div>
        </div>

        <!-- RSS Feeds -->
        <div class="form-section">
          <h3>RSS Feeds</h3>

          <div class="feed-list">
            <div
              v-for="(feed, index) in formData.feedUrls"
              :key="index"
              class="feed-row"
            >
              <input
                v-model="formData.feedUrls[index]"
                type="url"
                placeholder="https://example.com/feed.xml"
                @input="markChanged"
              />
              <button
                type="button"
                @click="removeFeed(index)"
                class="danger remove-btn"
                :disabled="formData.feedUrls.length === 1"
              >
                Remove
              </button>
            </div>
            <button type="button" @click="addFeed" class="add-feed-btn">
              + Add RSS Feed
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Archive Modal -->
    <div
      v-if="showArchive"
      class="modal-overlay"
      @click.self="showArchive = false"
    >
      <div class="modal-content archive-modal">
        <div class="modal-header">
          <h3>Dossier Archive</h3>
          <button @click="showArchive = false" class="close-btn">
            &times;
          </button>
        </div>
        <div class="archive-content">
          <div class="archive-header">
            <p>Last {{ archiveLimit }} generated summaries</p>
            <div class="archive-settings">
              <label>History Limit:</label>
              <select v-model="archiveLimit" @change="updateArchiveLimit">
                <option value="50">50</option>
                <option value="100">100</option>
                <option value="200">200</option>
                <option value="500">500</option>
              </select>
              <small>Records beyond this limit will be removed (FIFO)</small>
            </div>
          </div>
          <div class="archive-list">
            <div v-if="archiveDeliveries.length === 0" class="empty-archive">
              <p>No delivery history found for this dossier.</p>
            </div>
            <div
              v-for="delivery in archiveDeliveries"
              :key="delivery.id"
              class="archive-item"
              @click="openDeliveryDetail(delivery)"
            >
              <div class="delivery-date">
                {{ formatDate(delivery.sentAt) }}
              </div>
              <div class="delivery-info">
                <div class="delivery-subject">{{ delivery.subject }}</div>
                <div class="delivery-summary">
                  {{ delivery.content.substring(0, 100) }}...
                </div>
              </div>
              <button class="view-details-btn">View Full Content</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Delivery Detail Modal -->
    <div
      v-if="showDeliveryDetail"
      class="modal-overlay"
      @click.self="showDeliveryDetail = false"
    >
      <div class="modal-content delivery-modal">
        <div class="modal-header">
          <h3>{{ selectedDelivery?.subject || "Dossier Content" }}</h3>
          <button @click="showDeliveryDetail = false" class="close-btn">
            &times;
          </button>
        </div>
        <div class="delivery-content">
          <div class="delivery-meta">
            <p>
              <strong>Sent:</strong> {{ formatDate(selectedDelivery?.sentAt) }}
            </p>
            <p><strong>Subject:</strong> {{ selectedDelivery?.subject }}</p>
          </div>
          <div class="delivery-body" v-html="selectedDelivery?.content"></div>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div
      v-if="showDeleteConfirmation"
      class="modal-overlay"
      @click.self="showDeleteConfirmation = false"
    >
      <div class="modal-content delete-modal">
        <div class="modal-header">
          <h3>Delete Dossier</h3>
          <button @click="showDeleteConfirmation = false" class="close-btn">
            &times;
          </button>
        </div>
        <div class="delete-content">
          <p><strong>Warning:</strong> This action cannot be undone.</p>
          <p>
            All configuration data and delivery history for this dossier will be
            permanently deleted.
          </p>
          <div class="delete-actions">
            <button @click="showDeleteConfirmation = false" class="secondary">
              Cancel
            </button>
            <button @click="deleteDossier" class="danger">
              Delete Permanently
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Test Generate Modal -->
    <div
      v-if="showTestGenerate"
      class="modal-overlay"
      @click.self="showTestGenerate = false"
    >
      <div class="modal-content test-generate-modal">
        <div class="modal-header">
          <h3>Test Generate Results</h3>
          <button @click="showTestGenerate = false" class="close-btn">
            &times;
          </button>
        </div>
        <div class="test-generate-content">
          <div class="test-section">
            <h4>Generated Prompt</h4>
            <pre class="test-prompt">{{ testResults.prompt }}</pre>
          </div>
          <div class="test-section">
            <h4>AI Response</h4>
            <div class="test-response" v-html="testResults.response"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Custom Tone Modal -->
    <div
      v-if="showCustomToneModal"
      class="modal-overlay"
      @click.self="showCustomToneModal = false"
    >
      <div class="modal-content custom-tone-modal">
        <div class="modal-header">
          <h3>Create Custom Tone</h3>
          <button @click="showCustomToneModal = false" class="close-btn">
            &times;
          </button>
        </div>
        <div class="custom-tone-content">
          <div class="form-group">
            <label>Tone Name*</label>
            <input
              v-model="customToneData.name"
              type="text"
              placeholder="e.g., Pirate, Shakespearean, Tech Bro"
              required
            />
          </div>
          <div class="form-group">
            <label>Tone Prompt*</label>
            <textarea
              v-model="customToneData.prompt"
              placeholder="Describe how the AI should write in this tone. Be specific about style, vocabulary, and personality..."
              rows="6"
              required
            ></textarea>
            <small>
              Example: "Write like a pirate captain. Use nautical terms, 'arrr'
              frequently, and refer to the reader as 'matey' or 'landlubber'."
            </small>
          </div>
          <div class="modal-actions">
            <button @click="showCustomToneModal = false" class="secondary">
              Cancel
            </button>
            <button
              @click="createCustomTone"
              :disabled="
                !customToneData.name.trim() ||
                !customToneData.prompt.trim() ||
                loading
              "
              class="primary"
            >
              Create Tone
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted, watch, nextTick } from "vue";
import { useStore } from "../store/index.js";
import { useRouter } from "vue-router";

export default {
  name: "DossierDetailView",
  props: {
    dossierId: {
      type: [String, Number],
      required: true,
    },
  },
  setup(props) {
    const store = useStore();
    const router = useRouter();
    const loading = ref(false);
    const error = ref("");
    const success = ref("");
    const hasChanges = ref(false);
    const showArchive = ref(false);
    const showDeleteConfirmation = ref(false);
    const showDeliveryDetail = ref(false);
    const showTestGenerate = ref(false);
    const selectedDelivery = ref(null);
    const testResults = ref({ prompt: "", response: "" });
    const archiveDeliveries = ref([]);
    const archiveLimit = ref(100);
    const availableTones = ref([]);
    const showCustomToneModal = ref(false);
    const customToneData = reactive({
      name: "",
      prompt: "",
    });

    const isCreating = computed(() => props.dossierId === "new");

    // Watch for changes in form data to enable/disable save button
    const hasFormData = computed(() => {
      return (
        formData.title.trim().length > 0 &&
        formData.emails.some((email) => email.trim().length > 0) &&
        formData.feedUrls.some((url) => url.trim().length > 0)
      );
    });

    const formData = reactive({
      title: "",
      emails: [""],
      feedUrls: [""],
      articleCount: 10,
      frequency: "daily",
      deliveryTime: "08:00",
      timezone: "America/New_York",
      tone: "professional",
      language: "english",
      specialInstructions: "",
      active: true,
    });

    const originalFormData = ref({});

    const checkForChanges = () => {
      if (
        !originalFormData.value ||
        Object.keys(originalFormData.value).length === 0
      ) {
        hasChanges.value = isCreating.value;
        return;
      }

      const current = JSON.stringify({
        title: formData.title,
        emails: formData.emails.filter((email) => email.trim()),
        feedUrls: formData.feedUrls.filter((url) => url.trim()),
        articleCount: formData.articleCount,
        frequency: formData.frequency,
        deliveryTime: formData.deliveryTime,
        timezone: formData.timezone,
        tone: formData.tone,
        language: formData.language,
        specialInstructions: formData.specialInstructions,
        active: formData.active,
      });

      const original = JSON.stringify(originalFormData.value);
      hasChanges.value = current !== original;
    };

    const markChanged = () => {
      checkForChanges();
    };

    const addEmail = () => {
      formData.emails.push("");
      markChanged();
    };

    const removeEmail = (index) => {
      if (formData.emails.length > 1) {
        formData.emails.splice(index, 1);
        markChanged();
      }
    };

    const addFeed = () => {
      formData.feedUrls.push("");
      markChanged();
    };

    const removeFeed = (index) => {
      if (formData.feedUrls.length > 1) {
        formData.feedUrls.splice(index, 1);
        markChanged();
      }
    };

    const loadDossierData = async () => {
      if (!isCreating.value) {
        try {
          loading.value = true;
          const configs = await store.getDossierConfigs();
          const config = configs.find(
            (c) => c.id.toString() === props.dossierId.toString()
          );
          if (config) {
            // Reset form data to ensure clean state
            Object.assign(formData, {
              title: "",
              emails: [""],
              feedUrls: [""],
              articleCount: 10,
              frequency: "daily",
              deliveryTime: "08:00",
              timezone: "America/New_York",
              tone: "professional",
              language: "english",
              specialInstructions: "",
              active: true,
            });

            await nextTick();

            // Populate with actual data
            Object.assign(formData, {
              title: config.title || "",
              emails: config.email
                ? config.email.split(",").map((e) => e.trim())
                : [""],
              feedUrls:
                config.feedUrls && config.feedUrls.length > 0
                  ? [...config.feedUrls]
                  : [""],
              articleCount: config.articleCount || 10,
              frequency: config.frequency || "daily",
              deliveryTime: config.deliveryTime || "08:00",
              timezone: config.timezone || "America/New_York",
              tone: config.tone || "professional",
              language: config.language || "english",
              specialInstructions: config.specialInstructions || "",
              active: config.active !== undefined ? config.active : true,
            });

            const configData = {
              title: formData.title,
              emails: formData.emails,
              feedUrls: formData.feedUrls,
              articleCount: formData.articleCount,
              frequency: formData.frequency,
              deliveryTime: formData.deliveryTime,
              timezone: formData.timezone,
              tone: formData.tone,
              language: formData.language,
              specialInstructions: formData.specialInstructions,
              active: formData.active,
            };
            originalFormData.value = { ...configData };
            hasChanges.value = false;
          } else {
            error.value = `Dossier with ID ${props.dossierId} not found`;
          }
        } catch (err) {
          error.value = "Failed to load dossier data: " + (err.message || err);
        } finally {
          loading.value = false;
        }
      }
    };

    const goBack = () => {
      router.push("/");
    };

    const saveDossier = async () => {
      try {
        loading.value = true;
        error.value = "";

        const configData = {
          title: formData.title,
          email: formData.emails.filter((email) => email.trim()).join(", "),
          feedUrls: formData.feedUrls.filter((url) => url.trim()),
          articleCount: formData.articleCount,
          frequency: formData.frequency,
          deliveryTime: formData.deliveryTime,
          timezone: formData.timezone,
          tone: formData.tone,
          language: formData.language,
          specialInstructions: formData.specialInstructions,
        };

        if (isCreating.value) {
          await store.createDossierConfig(configData);
          success.value = "Dossier created successfully!";
        } else {
          await store.updateDossierConfig(
            props.dossierId.toString(),
            configData
          );
          success.value = "Dossier updated successfully!";
        }

        hasChanges.value = false;
        originalFormData.value = { ...formData };

        setTimeout(() => {
          success.value = "";
          router.push("/");
        }, 2000);
      } catch (err) {
        error.value = err.message || "Failed to save dossier";
      } finally {
        loading.value = false;
      }
    };

    const testEmail = async () => {
      if (isCreating.value) {
        error.value =
          "Please save the dossier first before sending a test email.";
        return;
      }

      try {
        loading.value = true;
        error.value = "";
        success.value = "Generating and sending test dossier email...";

        await store.generateAndSendDossier(props.dossierId.toString());

        success.value =
          "Test dossier email sent successfully! Check your inbox.";
        setTimeout(() => {
          success.value = "";
        }, 5000);
      } catch (err) {
        error.value = err.message || "Failed to send test email";
      } finally {
        loading.value = false;
      }
    };

    const testGenerate = async () => {
      try {
        loading.value = true;
        error.value = "";

        if (isCreating.value) {
          testResults.value = {
            prompt: "Cannot test generate for unsaved dossier",
            response:
              "Please save the dossier first before testing generation.",
          };
        } else {
          // Show sample prompt and response without actually generating
          const feedUrls = formData.feedUrls
            .filter((url) => url.trim())
            .join(", ");
          const emails = formData.emails
            .filter((email) => email.trim())
            .join(", ");

          testResults.value = {
            prompt: `AI Generation Test for "${formData.title}"
            
Configuration:
• Delivery Time: ${formData.deliveryTime} (${formData.timezone})
• Frequency: ${formData.frequency}
• Article Count: ${formData.articleCount}
• Tone: ${formData.tone}
• Language: ${formData.language}
• Email Recipients: ${emails}
• RSS Feeds: ${feedUrls}
• Special Instructions: ${formData.specialInstructions || "None"}

This test shows the configuration that would be used for AI generation. The actual generation would fetch recent articles from the RSS feeds and create a summary using the specified tone and language.`,
            response: `Sample AI Response Preview:

Subject: ${formData.title} - Daily Digest for ${new Date().toLocaleDateString()}

Dear Subscriber,

Here's your ${formData.tone} summary of today's top ${
              formData.articleCount
            } articles:

## Key Headlines:
• Breaking: Major technology breakthrough announced
• Market Update: Significant changes in the industry
• Analysis: Expert opinions on recent developments

[This is a preview. Actual generation would fetch real articles from your RSS feeds and create a personalized summary using ${
              formData.tone
            } tone in ${formData.language}.]

Best regards,
Your Dossier System`,
          };
        }

        showTestGenerate.value = true;
      } catch (err) {
        error.value = err.message || "Test generation failed";
      } finally {
        loading.value = false;
      }
    };

    const confirmDelete = () => {
      showDeleteConfirmation.value = true;
    };

    const deleteDossier = async () => {
      try {
        loading.value = true;
        error.value = "";
        await store.deleteDossierConfig(props.dossierId.toString());
        showDeleteConfirmation.value = false;
        success.value = "Dossier deleted successfully!";
        setTimeout(() => {
          router.push("/");
        }, 1500);
      } catch (err) {
        error.value = err.message || "Failed to delete dossier";
      } finally {
        loading.value = false;
      }
    };

    const openDeliveryDetail = (delivery) => {
      selectedDelivery.value = delivery;
      showDeliveryDetail.value = true;
    };

    const openArchive = async () => {
      showArchive.value = true;
      await loadArchiveDeliveries();
    };

    const updateArchiveLimit = async () => {
      // Reload deliveries with new limit
      await loadArchiveDeliveries();
    };

    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleString();
    };

    const loadArchiveDeliveries = async () => {
      if (!isCreating.value) {
        try {
          const deliveries = await store.getDossierDeliveries(
            props.dossierId.toString(),
            archiveLimit.value
          );
          archiveDeliveries.value = deliveries;
        } catch (err) {
          console.error("Failed to load archive deliveries:", err);
        }
      }
    };

    const loadTones = async () => {
      try {
        const response = await fetch("/graphql", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            query: `
              query GetTones {
                tones {
                  id
                  name
                  prompt
                  isSystemDefault
                  createdAt
                  updatedAt
                }
              }
            `,
          }),
        });

        const result = await response.json();
        if (result.errors) {
          throw new Error(result.errors[0].message);
        }

        availableTones.value = result.data.tones || [];
      } catch (err) {
        console.error("Failed to load tones:", err);
      }
    };

    const openCustomToneModal = () => {
      customToneData.name = "";
      customToneData.prompt = "";
      showCustomToneModal.value = true;
    };

    const createCustomTone = async () => {
      try {
        loading.value = true;
        error.value = "";

        const response = await fetch("/graphql", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            query: `
              mutation CreateTone($input: ToneInput!) {
                createTone(input: $input) {
                  id
                  name
                  prompt
                  isSystemDefault
                  createdAt
                  updatedAt
                }
              }
            `,
            variables: {
              input: {
                name: customToneData.name,
                prompt: customToneData.prompt,
              },
            },
          }),
        });

        const result = await response.json();
        if (result.errors) {
          throw new Error(result.errors[0].message);
        }

        const newTone = result.data.createTone;
        availableTones.value.push(newTone);
        formData.tone = newTone.name;
        showCustomToneModal.value = false;
        success.value = "Custom tone created successfully!";
        markChanged();

        setTimeout(() => {
          success.value = "";
        }, 3000);
      } catch (err) {
        error.value = err.message || "Failed to create custom tone";
      } finally {
        loading.value = false;
      }
    };

    // Watch for active toggle changes to sync with backend scheduler
    watch(
      () => formData.active,
      async (newActive, oldActive) => {
        if (
          oldActive !== undefined &&
          !isCreating.value &&
          originalFormData.value.active !== undefined
        ) {
          try {
            loading.value = true;
            await store.toggleDossierConfig(
              props.dossierId.toString(),
              newActive
            );
            success.value = newActive
              ? "Dossier activated successfully!"
              : "Dossier deactivated successfully!";
            setTimeout(() => {
              success.value = "";
            }, 2000);
          } catch (err) {
            error.value = err.message || "Failed to toggle dossier status";
            // Revert the toggle on error
            formData.active = oldActive;
          } finally {
            loading.value = false;
          }
        }
      }
    );

    // Watch for dossierId changes to reload data
    watch(
      () => props.dossierId,
      () => {
        loadDossierData();
        if (isCreating.value) {
          checkForChanges();
        }
      },
      { immediate: true }
    );

    onMounted(() => {
      loadTones();
    });

    return {
      loading,
      error,
      success,
      hasChanges,
      showArchive,
      showDeleteConfirmation,
      showDeliveryDetail,
      showTestGenerate,
      selectedDelivery,
      testResults,
      archiveDeliveries,
      archiveLimit,
      isCreating,
      hasFormData,
      formData,
      markChanged,
      checkForChanges,
      addEmail,
      removeEmail,
      addFeed,
      removeFeed,
      loadDossierData,
      goBack,
      saveDossier,
      testEmail,
      testGenerate,
      confirmDelete,
      deleteDossier,
      openArchive,
      openDeliveryDetail,
      updateArchiveLimit,
      formatDate,
      loadArchiveDeliveries,
      availableTones,
      showCustomToneModal,
      customToneData,
      loadTones,
      openCustomToneModal,
      createCustomTone,
    };
  },
};
</script>

<style scoped>
/**
 * DossierDetailView Component Styles
 * View-specific styles only - common patterns in global styles
 */

.dossier-detail {
  background: var(--color-bg-elevated);
  backdrop-filter: var(--blur-md);
  border: 1px solid var(--color-border-subtle);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  max-width: 1200px;
  margin: 0 auto;
  color: var(--color-text-primary);
  overflow: hidden;
}

/* Header Section - Specific styling for detail view */
.header-section {
  border-bottom: 1px solid var(--color-border-subtle);
  padding: var(--spacing-xl) var(--spacing-xl) var(--spacing-lg)
    var(--spacing-xl);
  margin-bottom: 0;
}

.header-content h2 {
  margin: 0;
  background: linear-gradient(
    135deg,
    var(--color-accent-blue),
    var(--color-accent-green)
  );
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.back-btn {
  background: none;
  border: none;
  color: var(--color-accent-blue);
  cursor: pointer;
  font-size: 1rem;
  margin-bottom: var(--spacing-md);
  padding: var(--spacing-sm) 0;
  transition: var(--transition-fast);
}

.back-btn:hover {
  color: var(--color-accent-green);
  transform: translateX(-2px);
  box-shadow: none;
}

/* Detail Content Layout */
.detail-content {
  display: grid;
  grid-template-columns: 1fr;
  gap: var(--spacing-xl);
  padding: 0 var(--spacing-xl) var(--spacing-xl) var(--spacing-xl);
}

/* Toggle Switch - Unique component */
.toggle-switch {
  position: relative;
  display: inline-block;
  width: 60px;
  height: 34px;
  margin-bottom: 0;
}

.toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.toggle-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(55, 65, 81, 0.8);
  transition: all 0.3s ease;
  border-radius: 34px;
  border: 1px solid var(--color-border-subtle);
}

.toggle-slider:before {
  position: absolute;
  content: "";
  height: 26px;
  width: 26px;
  left: 4px;
  bottom: 4px;
  background: linear-gradient(135deg, #e5e5e7, #d1d5db);
  transition: all 0.3s ease;
  border-radius: 50%;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

input:checked + .toggle-slider {
  background: linear-gradient(
    135deg,
    var(--color-primary),
    var(--color-primary-dark)
  );
  border-color: rgba(59, 130, 246, 0.3);
}

input:checked + .toggle-slider:before {
  transform: translateX(26px);
  background: linear-gradient(135deg, #ffffff, #f8fafc);
}

/* Message container positioning */
.message-container {
  padding: 0 var(--spacing-xl);
}

/* Archive Modal - Specific widths */
.archive-modal {
  width: 800px;
}

.archive-content {
  padding: var(--spacing-lg);
}

.archive-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
}

.archive-settings {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.archive-settings label {
  color: var(--color-text-primary);
  font-weight: 500;
}

.archive-settings small {
  color: var(--color-text-secondary);
  margin-left: var(--spacing-sm);
}

.archive-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  max-height: 60vh;
  overflow-y: auto;
}

.empty-archive {
  text-align: center;
  padding: 3rem 1rem;
  color: var(--color-text-secondary);
}

.archive-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md);
  background: rgba(17, 17, 19, 0.6);
  border: 1px solid var(--color-border-subtle);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: var(--transition-fast);
}

.archive-item:hover {
  background: rgba(55, 65, 81, 0.6);
  border-color: rgba(59, 130, 246, 0.3);
  transform: translateY(-2px);
}

.delivery-date {
  font-weight: 500;
  color: var(--color-accent-blue);
  min-width: 180px;
}

.delivery-info {
  flex: 1;
  margin-left: var(--spacing-md);
}

.delivery-subject {
  font-weight: 500;
  color: var(--color-text-primary);
  margin-bottom: 0.25rem;
}

.delivery-summary {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
}

.view-details-btn {
  background: linear-gradient(
    135deg,
    var(--color-primary),
    var(--color-primary-dark)
  );
  color: white;
  border: none;
  padding: var(--spacing-sm) var(--spacing-md);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: var(--transition-fast);
  font-weight: 500;
}

.view-details-btn:hover {
  background: linear-gradient(
    135deg,
    var(--color-primary-darker),
    var(--color-primary-dark)
  );
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
  transform: translateY(-1px);
}

/* Delivery Modal */
.delivery-modal {
  width: 900px;
}

.delivery-content {
  padding: var(--spacing-lg);
  max-height: 70vh;
  overflow: auto;
}

.delivery-meta {
  background: rgba(17, 17, 19, 0.6);
  border: 1px solid var(--color-border-subtle);
  border-radius: var(--radius-md);
  padding: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
}

.delivery-meta p {
  margin: var(--spacing-sm) 0;
  color: var(--color-text-primary);
}

.delivery-meta strong {
  color: var(--color-accent-blue);
}

.delivery-body {
  background: rgba(17, 17, 19, 0.4);
  border: 1px solid var(--color-border-subtle);
  border-radius: var(--radius-md);
  padding: var(--spacing-lg);
  line-height: 1.7;
  color: var(--color-text-primary);
}

/* Delete Modal */
.delete-modal {
  width: 400px;
}

.delete-content {
  padding: var(--spacing-lg);
}

.delete-content p {
  margin: 0 0 var(--spacing-md) 0;
  color: var(--color-text-primary);
  line-height: 1.6;
}

.delete-content p strong {
  color: var(--color-danger-light);
}

.delete-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-md);
  margin-top: var(--spacing-lg);
}

/* Test Generate Modal */
.test-generate-modal {
  width: 900px;
}

.test-generate-content {
  padding: var(--spacing-lg);
}

.test-section {
  margin-bottom: var(--spacing-xl);
}

.test-section h4 {
  margin: 0 0 var(--spacing-md) 0;
  color: var(--color-accent-blue);
  font-weight: 600;
}

.test-prompt {
  background: rgba(17, 17, 19, 0.6);
  border: 1px solid var(--color-border-subtle);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  font-family: "Monaco", "Menlo", "Courier New", monospace;
  white-space: pre-wrap;
  font-size: 0.875rem;
  color: var(--color-success-light);
  line-height: 1.5;
}

.test-response {
  background: rgba(17, 17, 19, 0.4);
  border: 1px solid var(--color-border-subtle);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  line-height: 1.7;
  color: var(--color-text-primary);
}

/* Custom Tone Modal */
.custom-tone-modal {
  max-width: 600px;
}

.custom-tone-content {
  padding: var(--spacing-lg);
}

.custom-tone-content .form-group {
  margin-bottom: var(--spacing-lg);
}

.custom-tone-content label {
  color: var(--color-text-primary);
  font-weight: 500;
  margin-bottom: var(--spacing-sm);
  display: block;
}

.custom-tone-content textarea {
  min-height: 120px;
  resize: vertical;
}

.custom-tone-content small {
  display: block;
  margin-top: var(--spacing-sm);
  color: var(--color-text-secondary);
  font-style: italic;
}

/* Tone Selection */
.tone-selection {
  display: flex;
  gap: var(--spacing-sm);
  align-items: center;
}

.tone-selection select {
  flex: 1;
}

/* Custom Scrollbars for modals */
.archive-list::-webkit-scrollbar,
.delivery-content::-webkit-scrollbar {
  width: 8px;
}

.archive-list::-webkit-scrollbar-track,
.delivery-content::-webkit-scrollbar-track {
  background: rgba(17, 17, 19, 0.4);
  border-radius: 4px;
}

.archive-list::-webkit-scrollbar-thumb,
.delivery-content::-webkit-scrollbar-thumb {
  background: rgba(96, 165, 250, 0.4);
  border-radius: 4px;
}

.archive-list::-webkit-scrollbar-thumb:hover,
.delivery-content::-webkit-scrollbar-thumb:hover {
  background: rgba(96, 165, 250, 0.6);
}
</style>
