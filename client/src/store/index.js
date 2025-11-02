import { reactive } from "vue";

const store = reactive({
  // Application state
  loading: false,
  error: null,
  dossierConfigs: [],

  // Dossier Configuration Methods
  async getDossierConfigs() {
    try {
      const response = await fetch("/graphql", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          query: `
            query GetDossierConfigs {
              dossierConfigs {
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
              }
            }
          `,
        }),
      });

      const result = await response.json();
      if (result.errors) {
        throw new Error(result.errors[0].message);
      }

      this.dossierConfigs = result.data.dossierConfigs || [];
      return this.dossierConfigs;
    } catch (error) {
      console.error("Failed to fetch dossier configs:", error);
      throw error;
    }
  },

  async createDossierConfig(configData) {
    try {
      const response = await fetch("/graphql", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          query: `
            mutation CreateDossierConfig($input: DossierConfigInput!) {
              createDossierConfig(input: $input) {
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
              }
            }
          `,
          variables: { input: configData },
        }),
      });

      const result = await response.json();
      if (result.errors) {
        throw new Error(result.errors[0].message);
      }

      const newConfig = result.data.createDossierConfig;
      this.dossierConfigs.push(newConfig);
      return newConfig;
    } catch (error) {
      console.error("Failed to create dossier config:", error);
      throw error;
    }
  },

  async updateDossierConfig(id, configData) {
    try {
      const response = await fetch("/graphql", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          query: `
            mutation UpdateDossierConfig($id: ID!, $input: DossierConfigInput!) {
              updateDossierConfig(id: $id, input: $input) {
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
              }
            }
          `,
          variables: { id, input: configData },
        }),
      });

      const result = await response.json();
      if (result.errors) {
        throw new Error(result.errors[0].message);
      }

      const updatedConfig = result.data.updateDossierConfig;
      const index = this.dossierConfigs.findIndex((config) => config.id === id);
      if (index !== -1) {
        this.dossierConfigs[index] = updatedConfig;
      }
      return updatedConfig;
    } catch (error) {
      console.error("Failed to update dossier config:", error);
      throw error;
    }
  },

  async deleteDossierConfig(id) {
    try {
      const response = await fetch("/graphql", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          query: `
            mutation DeleteDossierConfig($id: ID!) {
              deleteDossierConfig(id: $id)
            }
          `,
          variables: { id },
        }),
      });

      const result = await response.json();
      if (result.errors) {
        throw new Error(result.errors[0].message);
      }

      this.dossierConfigs = this.dossierConfigs.filter(
        (config) => config.id !== id
      );
      return true;
    } catch (error) {
      console.error("Failed to delete dossier config:", error);
      throw error;
    }
  },

  async generateAndSendDossier(configId) {
    try {
      const response = await fetch("/graphql", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          query: `
            mutation GenerateAndSendDossier($configId: ID!) {
              generateAndSendDossier(configId: $configId)
            }
          `,
          variables: { configId },
        }),
      });

      const result = await response.json();
      if (result.errors) {
        throw new Error(result.errors[0].message);
      }

      return result.data.generateAndSendDossier;
    } catch (error) {
      console.error("Failed to generate and send dossier:", error);
      throw error;
    }
  },

  async testEmailConnection() {
    try {
      const response = await fetch("/graphql", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          query: `
            mutation TestEmailConnection {
              testEmailConnection
            }
          `,
        }),
      });

      const result = await response.json();
      if (result.errors) {
        throw new Error(result.errors[0].message);
      }

      return result.data.testEmailConnection;
    } catch (error) {
      console.error("Failed to test email connection:", error);
      throw error;
    }
  },

  async getDossierDeliveries(configId, limit = 100) {
    try {
      const response = await fetch("/graphql", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          query: `
            query GetDossierDeliveries($configId: ID, $limit: Int) {
              dossiers(configId: $configId, limit: $limit) {
                id
                configId
                subject
                content
                sentAt
              }
            }
          `,
          variables: { configId, limit },
        }),
      });

      const result = await response.json();
      if (result.errors) {
        throw new Error(result.errors[0].message);
      }

      return result.data.dossiers || [];
    } catch (error) {
      console.error("Failed to fetch dossier deliveries:", error);
      throw error;
    }
  },

  async toggleDossierConfig(id, active) {
    try {
      const response = await fetch("/graphql", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          query: `
            mutation ToggleDossierConfig($id: ID!, $active: Boolean!) {
              toggleDossierConfig(id: $id, active: $active) {
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
              }
            }
          `,
          variables: { id, active },
        }),
      });

      const result = await response.json();
      if (result.errors) {
        throw new Error(result.errors[0].message);
      }

      const updatedConfig = result.data.toggleDossierConfig;
      const index = this.dossierConfigs.findIndex((config) => config.id === id);
      if (index !== -1) {
        this.dossierConfigs[index] = updatedConfig;
      }
      return updatedConfig;
    } catch (error) {
      console.error("Failed to toggle dossier config:", error);
      throw error;
    }
  },

  async getTones() {
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

      return result.data.tones || [];
    } catch (error) {
      console.error("Failed to fetch tones:", error);
      throw error;
    }
  },

  async createTone(toneData) {
    try {
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
          variables: { input: toneData },
        }),
      });

      const result = await response.json();
      if (result.errors) {
        throw new Error(result.errors[0].message);
      }

      return result.data.createTone;
    } catch (error) {
      console.error("Failed to create tone:", error);
      throw error;
    }
  },

  async updateTone(id, toneData) {
    try {
      const response = await fetch("/graphql", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          query: `
            mutation UpdateTone($id: ID!, $input: ToneInput!) {
              updateTone(id: $id, input: $input) {
                id
                name
                prompt
                isSystemDefault
                createdAt
                updatedAt
              }
            }
          `,
          variables: { id, input: toneData },
        }),
      });

      const result = await response.json();
      if (result.errors) {
        throw new Error(result.errors[0].message);
      }

      return result.data.updateTone;
    } catch (error) {
      console.error("Failed to update tone:", error);
      throw error;
    }
  },

  async deleteTone(id) {
    try {
      const response = await fetch("/graphql", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          query: `
            mutation DeleteTone($id: ID!) {
              deleteTone(id: $id)
            }
          `,
          variables: { id },
        }),
      });

      const result = await response.json();
      if (result.errors) {
        throw new Error(result.errors[0].message);
      }

      return result.data.deleteTone;
    } catch (error) {
      console.error("Failed to delete tone:", error);
      throw error;
    }
  },

  clearError() {
    this.error = null;
  },
});

export function useStore() {
  return store;
}
