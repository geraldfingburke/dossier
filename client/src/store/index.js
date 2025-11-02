import { reactive } from "vue";

// ============================================================================
// Constants
// ============================================================================

const GRAPHQL_ENDPOINT = "/graphql";
const HTTP_METHOD_POST = "POST";
const CONTENT_TYPE_JSON = "application/json";

const HTTP_HEADERS = {
  "Content-Type": CONTENT_TYPE_JSON,
};

// ============================================================================
// GraphQL Query Fragments
// ============================================================================

/**
 * Complete dossier configuration fields for GraphQL queries
 */
const DOSSIER_CONFIG_FIELDS = `
  id
  title
  email
  feedUrls
  articleCount
  frequency
  deliveryTime
  timezone
  tone
  language
  specialInstructions
  active
  createdAt
`;

/**
 * Complete tone fields for GraphQL queries
 */
const TONE_FIELDS = `
  id
  name
  prompt
  isSystemDefault
  createdAt
  updatedAt
`;

/**
 * Complete delivery/dossier fields for GraphQL queries
 */
const DELIVERY_FIELDS = `
  id
  configId
  subject
  content
  sentAt
`;

// ============================================================================
// Helper Functions
// ============================================================================

/**
 * Executes a GraphQL request to the backend API
 * @param {string} query - The GraphQL query or mutation string
 * @param {Object} variables - Variables for the GraphQL operation
 * @returns {Promise<Object>} The data from the GraphQL response
 * @throws {Error} If the request fails or returns errors
 */
async function executeGraphQLRequest(query, variables = {}) {
  const response = await fetch(GRAPHQL_ENDPOINT, {
    method: HTTP_METHOD_POST,
    headers: HTTP_HEADERS,
    body: JSON.stringify({ query, variables }),
  });

  const result = await response.json();

  if (result.errors) {
    throw new Error(result.errors[0].message);
  }

  return result.data;
}

/**
 * Finds and replaces an item in an array by ID
 * @param {Array} array - The array to search
 * @param {string|number} id - The ID to match
 * @param {Object} newItem - The new item to replace with
 * @returns {boolean} True if item was found and replaced
 */
function replaceItemById(array, id, newItem) {
  const itemIndex = array.findIndex((item) => item.id === id);
  if (itemIndex !== -1) {
    array[itemIndex] = newItem;
    return true;
  }
  return false;
}

/**
 * Removes an item from an array by ID
 * @param {Array} array - The array to filter
 * @param {string|number} id - The ID to remove
 * @returns {Array} New array without the item
 */
function removeItemById(array, id) {
  return array.filter((item) => item.id !== id);
}

// ============================================================================
// Store Definition
// ============================================================================

/**
 * Global application store for managing dossier configurations,
 * deliveries, tones, and API communication
 */
const store = reactive({
  // ============================================================================
  // State
  // ============================================================================

  /** Indicates if an async operation is in progress */
  loading: false,

  /** Current error message, null if no error */
  error: null,

  /** Array of all dossier configurations */
  dossierConfigs: [],

  // ============================================================================
  // Dossier Configuration Methods
  // ============================================================================

  /**
   * Fetches all dossier configurations from the backend
   * @returns {Promise<Array>} Array of dossier configuration objects
   * @throws {Error} If the fetch operation fails
   */
  /**
   * Fetches all dossier configurations from the backend
   * @returns {Promise<Array>} Array of dossier configuration objects
   * @throws {Error} If the fetch operation fails
   */
  async getDossierConfigs() {
    try {
      const data = await executeGraphQLRequest(`
        query GetDossierConfigs {
          dossierConfigs {
            ${DOSSIER_CONFIG_FIELDS}
          }
        }
      `);

      this.dossierConfigs = data.dossierConfigs || [];
      return this.dossierConfigs;
    } catch (error) {
      console.error("Failed to fetch dossier configs:", error);
      throw error;
    }
  },

  /**
   * Creates a new dossier configuration
   * @param {Object} configData - The configuration data for the new dossier
   * @returns {Promise<Object>} The newly created dossier configuration
   * @throws {Error} If the creation fails
   */
  async createDossierConfig(configData) {
    try {
      const data = await executeGraphQLRequest(
        `
          mutation CreateDossierConfig($input: DossierConfigInput!) {
            createDossierConfig(input: $input) {
              ${DOSSIER_CONFIG_FIELDS}
            }
          }
        `,
        { input: configData }
      );

      const newConfig = data.createDossierConfig;
      this.dossierConfigs.push(newConfig);
      return newConfig;
    } catch (error) {
      console.error("Failed to create dossier config:", error);
      throw error;
    }
  },

  /**
   * Updates an existing dossier configuration
   * @param {string} configId - The ID of the configuration to update
   * @param {Object} configData - The updated configuration data
   * @returns {Promise<Object>} The updated dossier configuration
   * @throws {Error} If the update fails
   */
  /**
   * Updates an existing dossier configuration
   * @param {string} configId - The ID of the configuration to update
   * @param {Object} configData - The updated configuration data
   * @returns {Promise<Object>} The updated dossier configuration
   * @throws {Error} If the update fails
   */
  async updateDossierConfig(configId, configData) {
    try {
      const data = await executeGraphQLRequest(
        `
          mutation UpdateDossierConfig($id: ID!, $input: DossierConfigInput!) {
            updateDossierConfig(id: $id, input: $input) {
              ${DOSSIER_CONFIG_FIELDS}
            }
          }
        `,
        { id: configId, input: configData }
      );

      const updatedConfig = data.updateDossierConfig;
      replaceItemById(this.dossierConfigs, configId, updatedConfig);
      return updatedConfig;
    } catch (error) {
      console.error("Failed to update dossier config:", error);
      throw error;
    }
  },

  /**
   * Deletes a dossier configuration
   * @param {string} configId - The ID of the configuration to delete
   * @returns {Promise<boolean>} True if deletion was successful
   * @throws {Error} If the deletion fails
   */
  async deleteDossierConfig(configId) {
    try {
      await executeGraphQLRequest(
        `
          mutation DeleteDossierConfig($id: ID!) {
            deleteDossierConfig(id: $id)
          }
        `,
        { id: configId }
      );

      this.dossierConfigs = removeItemById(this.dossierConfigs, configId);
      return true;
    } catch (error) {
      console.error("Failed to delete dossier config:", error);
      throw error;
    }
  },

  /**
   * Toggles the active status of a dossier configuration
   * @param {string} configId - The ID of the configuration to toggle
   * @param {boolean} isActive - The new active status
   * @returns {Promise<Object>} The updated dossier configuration
   * @throws {Error} If the toggle operation fails
   */
  async toggleDossierConfig(configId, isActive) {
    try {
      const data = await executeGraphQLRequest(
        `
          mutation ToggleDossierConfig($id: ID!, $active: Boolean!) {
            toggleDossierConfig(id: $id, active: $active) {
              ${DOSSIER_CONFIG_FIELDS}
            }
          }
        `,
        { id: configId, active: isActive }
      );

      const updatedConfig = data.toggleDossierConfig;
      replaceItemById(this.dossierConfigs, configId, updatedConfig);
      return updatedConfig;
    } catch (error) {
      console.error("Failed to toggle dossier config:", error);
      throw error;
    }
  },

  // ============================================================================
  // Dossier Generation & Delivery Methods
  // ============================================================================

  /**
   * Generates and sends a dossier immediately for the specified configuration
   * @param {string} configId - The ID of the configuration to generate from
   * @returns {Promise<boolean>} True if generation and sending was successful
   * @throws {Error} If the operation fails
   */
  /**
   * Generates and sends a dossier immediately for the specified configuration
   * @param {string} configId - The ID of the configuration to generate from
   * @returns {Promise<boolean>} True if generation and sending was successful
   * @throws {Error} If the operation fails
   */
  async generateAndSendDossier(configId) {
    try {
      const data = await executeGraphQLRequest(
        `
          mutation GenerateAndSendDossier($configId: ID!) {
            generateAndSendDossier(configId: $configId)
          }
        `,
        { configId }
      );

      return data.generateAndSendDossier;
    } catch (error) {
      console.error("Failed to generate and send dossier:", error);
      throw error;
    }
  },

  /**
   * Retrieves the delivery history for a specific dossier configuration
   * @param {string} configId - The ID of the configuration
   * @param {number} limit - Maximum number of deliveries to retrieve (default: 100)
   * @returns {Promise<Array>} Array of delivery objects
   * @throws {Error} If the fetch operation fails
   */
  async getDossierDeliveries(configId, limit = 100) {
    try {
      const data = await executeGraphQLRequest(
        `
          query GetDossierDeliveries($configId: ID, $limit: Int) {
            dossiers(configId: $configId, limit: $limit) {
              ${DELIVERY_FIELDS}
            }
          }
        `,
        { configId, limit }
      );

      return data.dossiers || [];
    } catch (error) {
      console.error("Failed to fetch dossier deliveries:", error);
      throw error;
    }
  },

  // ============================================================================
  // Email Testing Methods
  // ============================================================================

  /**
   * Tests the email connection settings without sending an actual dossier
   * @returns {Promise<boolean>} True if email connection test succeeds
   * @throws {Error} If the connection test fails
   */
  async testEmailConnection() {
    try {
      const data = await executeGraphQLRequest(`
        mutation TestEmailConnection {
          testEmailConnection
        }
      `);

      return data.testEmailConnection;
    } catch (error) {
      console.error("Failed to test email connection:", error);
      throw error;
    }
  },

  // ============================================================================
  // Tone Management Methods
  // ============================================================================

  /**
   * Fetches all available tones (both system defaults and custom)
   * @returns {Promise<Array>} Array of tone objects
   * @throws {Error} If the fetch operation fails
   */
  async getTones() {
    try {
      const data = await executeGraphQLRequest(`
        query GetTones {
          tones {
            ${TONE_FIELDS}
          }
        }
      `);

      return data.tones || [];
    } catch (error) {
      console.error("Failed to fetch tones:", error);
      throw error;
    }
  },

  /**
   * Creates a new custom tone
   * @param {Object} toneData - The tone data (name and prompt)
   * @returns {Promise<Object>} The newly created tone object
   * @throws {Error} If the creation fails
   */
  async createTone(toneData) {
    try {
      const data = await executeGraphQLRequest(
        `
          mutation CreateTone($input: ToneInput!) {
            createTone(input: $input) {
              ${TONE_FIELDS}
            }
          }
        `,
        { input: toneData }
      );

      return data.createTone;
    } catch (error) {
      console.error("Failed to create tone:", error);
      throw error;
    }
  },

  /**
   * Updates an existing custom tone
   * @param {string} toneId - The ID of the tone to update
   * @param {Object} toneData - The updated tone data
   * @returns {Promise<Object>} The updated tone object
   * @throws {Error} If the update fails
   */
  async updateTone(toneId, toneData) {
    try {
      const data = await executeGraphQLRequest(
        `
          mutation UpdateTone($id: ID!, $input: ToneInput!) {
            updateTone(id: $id, input: $input) {
              ${TONE_FIELDS}
            }
          }
        `,
        { id: toneId, input: toneData }
      );

      return data.updateTone;
    } catch (error) {
      console.error("Failed to update tone:", error);
      throw error;
    }
  },

  /**
   * Deletes a custom tone (system default tones cannot be deleted)
   * @param {string} toneId - The ID of the tone to delete
   * @returns {Promise<boolean>} True if deletion was successful
   * @throws {Error} If the deletion fails
   */
  async deleteTone(toneId) {
    try {
      const data = await executeGraphQLRequest(
        `
          mutation DeleteTone($id: ID!) {
            deleteTone(id: $id)
          }
        `,
        { id: toneId }
      );

      return data.deleteTone;
    } catch (error) {
      console.error("Failed to delete tone:", error);
      throw error;
    }
  },

  // ============================================================================
  // Utility Methods
  // ============================================================================

  /**
   * Clears any error state in the store
   */
  clearError() {
    this.error = null;
  },
});

// ============================================================================
// Store Export
// ============================================================================

/**
 * Hook to access the global store instance
 * @returns {Object} The reactive store object
 */
export function useStore() {
  return store;
}
