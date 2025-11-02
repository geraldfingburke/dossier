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

  clearError() {
    this.error = null;
  },
});

export function useStore() {
  return store;
}
