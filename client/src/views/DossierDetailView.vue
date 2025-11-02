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
        <button @click="testEmail" :disabled="loading" class="secondary">
          Test Email
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
import {
  ref,
  reactive,
  computed,
  onMounted,
  watch,
  nextTick,
  toRefs,
} from "vue";
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
      email: "",
      emails: [""], // Support for multiple emails in the UI
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
    const debugInfo = ref({
      loaded: false,
      configId: null,
      formTitle: computed(() => formData.title),
    });

    const checkForChanges = () => {
      if (
        !originalFormData.value ||
        Object.keys(originalFormData.value).length === 0
      ) {
        hasChanges.value = isCreating.value; // For new dossiers, always show save button
        return;
      }

      // Deep compare objects
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
            // First reset form data to ensure clean state
            Object.assign(formData, {
              title: "",
              email: "",
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

            // Force a Vue update cycle
            await nextTick();

            // Now populate with actual data
            Object.assign(formData, {
              title: config.title || "",
              email: config.email || "",
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

            debugInfo.value.loaded = true;
            debugInfo.value.configId = config.id;
            const configData = {
              title: formData.title,
              email: formData.email,
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
      try {
        loading.value = true;
        error.value = "";
        await store.testEmailConnection();
        success.value = "Email connection test successful!";
        setTimeout(() => {
          success.value = "";
        }, 3000);
      } catch (err) {
        error.value = err.message || "Email test failed";
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
        // Direct GraphQL call to bypass store issue
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

        const tones = result.data.tones || [];
        availableTones.value = tones;
      } catch (err) {
        // Silently fail - tones will remain empty array
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

        // Direct GraphQL call to create tone
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
      ...toRefs(props),
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
      debugInfo,
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
.dossier-detail {
  background: rgba(26, 26, 30, 0.8);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 1rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  max-width: 1200px;
  margin: 0 auto;
  color: #e5e5e7;
  overflow: hidden;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  padding: 2rem 2rem 1.5rem 2rem;
}

.header-content h2 {
  margin: 0;
  font-size: 2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #60a5fa, #34d399);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header-content .subtitle {
  margin: 0.5rem 0 0 0;
  color: rgba(229, 229, 231, 0.6);
  font-size: 0.95rem;
}

.back-btn {
  background: none;
  border: none;
  color: #60a5fa;
  cursor: pointer;
  font-size: 1rem;
  margin-bottom: 0.75rem;
  padding: 0.5rem 0;
  transition: all 0.2s ease;
}

.back-btn:hover {
  color: #34d399;
  transform: translateX(-2px);
}

.header-buttons {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.header-buttons button {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 0.75rem;
  cursor: pointer;
  font-weight: 500;
  font-size: 0.875rem;
  transition: all 0.2s ease;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.primary {
  background: linear-gradient(135deg, #3b82f6, #1d4ed8);
  color: white;
  border: 1px solid rgba(59, 130, 246, 0.3);
}

.primary:hover:not(:disabled) {
  background: linear-gradient(135deg, #2563eb, #1e40af);
  box-shadow: 0 8px 24px rgba(59, 130, 246, 0.3);
  transform: translateY(-1px);
}

.secondary {
  background: rgba(55, 65, 81, 0.8);
  color: #e5e5e7;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.secondary:hover:not(:disabled) {
  background: rgba(75, 85, 99, 0.9);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  transform: translateY(-1px);
}

.danger {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: white;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.danger:hover:not(:disabled) {
  background: linear-gradient(135deg, #dc2626, #b91c1c);
  box-shadow: 0 8px 24px rgba(239, 68, 68, 0.3);
  transform: translateY(-1px);
}

button:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.save-btn:disabled {
  background: rgba(55, 65, 81, 0.4);
}

.detail-content {
  display: grid;
  grid-template-columns: 1fr;
  gap: 2rem;
  padding: 0 2rem 2rem 2rem;
}

.form-section {
  background: rgba(26, 26, 30, 0.8);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 1rem;
  padding: 2rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  margin-top: 1.5rem;
}

.form-section:first-child {
  margin-top: 0;
}

.form-section h3 {
  margin: 0 0 1.5rem 0;
  color: #e5e5e7;
  font-size: 1.25rem;
  font-weight: 600;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  font-weight: 500;
  margin-bottom: 0.75rem;
  color: #e5e5e7;
  font-size: 0.9rem;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 0.875rem;
  border: 1px solid rgba(55, 65, 81, 0.6);
  border-radius: 0.75rem;
  background: rgba(17, 17, 19, 0.8);
  color: #e5e5e7;
  font-size: 0.875rem;
  font-family: inherit;
  backdrop-filter: blur(10px);
  transition: all 0.2s ease;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  background: rgba(17, 17, 19, 0.9);
}

.form-group input::placeholder,
.form-group textarea::placeholder {
  color: rgba(229, 229, 231, 0.5);
}

.form-group textarea {
  min-height: 120px;
  resize: vertical;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

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
  border: 1px solid rgba(255, 255, 255, 0.1);
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
  background: linear-gradient(135deg, #3b82f6, #1d4ed8);
  border-color: rgba(59, 130, 246, 0.3);
}

input:checked + .toggle-slider:before {
  transform: translateX(26px);
  background: linear-gradient(135deg, #ffffff, #f8fafc);
}

/* Loading State */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  text-align: center;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(59, 130, 246, 0.2);
  border-top: 3px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.loading-state p {
  color: rgba(229, 229, 231, 0.8);
  font-size: 0.95rem;
}

.email-list,
.feed-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.email-row,
.feed-row {
  display: flex;
  gap: 0.75rem;
  align-items: center;
}

.email-row input,
.feed-row input {
  flex: 1;
}

.remove-btn {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: white;
  border: none;
  padding: 0.75rem 1rem;
  border-radius: 0.75rem;
  cursor: pointer;
  white-space: nowrap;
  font-weight: 500;
  font-size: 0.875rem;
  transition: all 0.2s ease;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.remove-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #dc2626, #b91c1c);
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
  transform: translateY(-1px);
}

.remove-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.add-email-btn,
.add-feed-btn {
  background: rgba(55, 65, 81, 0.8);
  color: #e5e5e7;
  border: 1px dashed rgba(255, 255, 255, 0.2);
  padding: 0.75rem;
  border-radius: 0.75rem;
  cursor: pointer;
  transition: all 0.2s ease;
  backdrop-filter: blur(10px);
  font-weight: 500;
  font-size: 0.875rem;
}

.add-email-btn:hover,
.add-feed-btn:hover {
  background: rgba(75, 85, 99, 0.9);
  border-color: rgba(255, 255, 255, 0.3);
  transform: translateY(-1px);
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: rgba(26, 26, 30, 0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 1rem;
  max-width: 90vw;
  max-height: 90vh;
  overflow: auto;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.5);
  color: #e5e5e7;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-header h3 {
  margin: 0;
  color: #2c3e50;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #7f8c8d;
}

.close-btn:hover {
  color: #2c3e50;
}

.archive-modal {
  width: 800px;
}

.archive-content {
  padding: 1.5rem;
}

.archive-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.archive-settings {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.archive-settings small {
  color: #7f8c8d;
  margin-left: 0.5rem;
}

.archive-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.archive-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border: 1px solid #e1e5e9;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.archive-item:hover {
  background-color: #f8f9fa;
}

.delivery-date {
  font-weight: 500;
  color: #2c3e50;
}

.delivery-info {
  flex: 1;
  margin-left: 1rem;
}

.delivery-info div {
  font-size: 0.9rem;
  color: #7f8c8d;
}

.view-details-btn {
  background: #3498db;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
}

.view-details-btn:hover {
  background: #2980b9;
}

.delivery-modal {
  width: 900px;
}

.delivery-content {
  padding: 1.5rem;
  max-height: 70vh;
  overflow: auto;
}

.delete-modal {
  width: 400px;
}

.delete-content {
  padding: 1.5rem;
}

.delete-content p {
  margin: 0 0 1rem 0;
  color: #2c3e50;
}

.delete-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 1.5rem;
}

.test-generate-modal {
  width: 900px;
}

.test-generate-content {
  padding: 1.5rem;
}

.test-section {
  margin-bottom: 2rem;
}

.test-section h4 {
  margin: 0 0 1rem 0;
  color: #2c3e50;
}

.test-prompt {
  background: #f8f9fa;
  padding: 1rem;
  border-radius: 4px;
  font-family: monospace;
  white-space: pre-wrap;
  font-size: 0.9rem;
}

.test-response {
  background: white;
  border: 1px solid #e1e5e9;
  padding: 1rem;
  border-radius: 4px;
}

.message-container {
  padding: 0 2rem;
}

.error {
  color: #fca5a5;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.2);
  padding: 1rem;
  border-radius: 0.75rem;
  margin-bottom: 1.5rem;
  backdrop-filter: blur(10px);
}

.success {
  color: #86efac;
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.2);
  padding: 1rem;
  border-radius: 0.75rem;
  margin-bottom: 1.5rem;
  backdrop-filter: blur(10px);
}

/* Custom tone styles */
.tone-selection {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.tone-selection select {
  flex: 1;
}

.add-tone-btn {
  background: #28a745;
  color: white;
  border: none;
  padding: 0.75rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  white-space: nowrap;
  font-size: 0.9rem;
}

.add-tone-btn:hover:not(:disabled) {
  background: #218838;
}

.custom-tone-modal {
  max-width: 600px;
}

.custom-tone-content {
  padding: 1.5rem;
}

.custom-tone-content .form-group {
  margin-bottom: 1.5rem;
}

.custom-tone-content textarea {
  min-height: 120px;
  resize: vertical;
}

.custom-tone-content small {
  display: block;
  margin-top: 0.5rem;
  color: #6c757d;
  font-style: italic;
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid #e1e5e9;
}

/* Dark Mode Styles */
@media (prefers-color-scheme: dark) {
  .dossier-detail {
    background: #1a1a1a;
    color: #ffffff;
  }

  .header-section {
    background: #2d2d2d;
    border-bottom-color: #404040;
  }

  .secondary {
    background: #404040;
    color: #ffffff;
    border-color: #555555;
  }

  .secondary:hover:not(:disabled) {
    background: #505050;
  }

  .form-section {
    background: #2d2d2d;
    border-color: #404040;
  }

  .form-group label {
    color: #ffffff;
  }

  .form-group input,
  .form-group select,
  .form-group textarea {
    background: #404040;
    border-color: #555555;
    color: #ffffff;
  }

  .form-group input:focus,
  .form-group select:focus,
  .form-group textarea:focus {
    border-color: #3498db;
  }

  .custom-tone-content small {
    color: #aaaaaa;
  }

  .modal-overlay {
    background: rgba(0, 0, 0, 0.8);
  }

  .modal-content {
    background: #2d2d2d;
    color: #ffffff;
  }

  .modal-actions {
    border-top-color: #404040;
  }

  .archive-section,
  .test-section {
    background: #2d2d2d;
    border-color: #404040;
  }

  .delivery-item {
    background: #404040;
    border-color: #555555;
  }

  .delivery-item:hover {
    background: #505050;
  }
}
</style>
