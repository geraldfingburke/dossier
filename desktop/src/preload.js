const { contextBridge, ipcRenderer } = require("electron");

// ============================================================================
// CONTEXT BRIDGE
// ============================================================================

/**
 * Preload script for Electron context isolation.
 *
 * This script creates a secure bridge between the renderer process (Vue.js app)
 * and the main Electron process. It exposes only specific, safe APIs to the
 * frontend while keeping Node.js and Electron APIs isolated.
 *
 * Security:
 * - contextIsolation: true (in main.js)
 * - nodeIntegration: false (in main.js)
 * - Only whitelisted APIs exposed via contextBridge
 */

contextBridge.exposeInMainWorld("electronAPI", {
  // ============================================================================
  // SERVICE COMMUNICATION
  // ============================================================================

  /**
   * Get the base URL for the backend service.
   *
   * @returns {Promise<string>} Service URL (e.g., 'http://localhost:8080')
   */
  getServiceUrl: () => ipcRenderer.invoke("get-service-url"),

  /**
   * Check if the backend service is healthy/running.
   *
   * @returns {Promise<boolean>} True if service is accessible
   */
  checkServiceHealth: () => ipcRenderer.invoke("check-service-health"),

  // ============================================================================
  // EXTERNAL LINKS
  // ============================================================================

  /**
   * Open a URL in the default system browser.
   *
   * @param {string} url - URL to open
   * @returns {Promise<void>}
   */
  openExternal: (url) => ipcRenderer.invoke("open-external", url),

  // ============================================================================
  // APPLICATION INFO
  // ============================================================================

  /**
   * Get the application version.
   *
   * @returns {Promise<string>} Version string (e.g., '1.0.0')
   */
  getAppVersion: () => ipcRenderer.invoke("get-app-version"),

  // ============================================================================
  // WINDOW MANAGEMENT
  // ============================================================================

  /**
   * Minimize the application window to system tray.
   *
   * @returns {Promise<void>}
   */
  minimizeToTray: () => ipcRenderer.invoke("minimize-to-tray"),

  // ============================================================================
  // NOTIFICATIONS
  // ============================================================================

  /**
   * Show a native system notification.
   *
   * @param {Object} options - Notification options
   * @param {string} options.title - Notification title
   * @param {string} options.body - Notification body text
   * @returns {Promise<void>}
   */
  showNotification: (options) =>
    ipcRenderer.invoke("show-notification", options),

  // ============================================================================
  // EVENT LISTENERS
  // ============================================================================

  /**
   * Listen for navigation events from main process.
   *
   * @param {Function} callback - Callback function (route: string) => void
   */
  onNavigate: (callback) => {
    ipcRenderer.on("navigate", (event, route) => callback(route));
  },

  /**
   * Remove navigation event listener.
   *
   * @param {Function} callback - Previously registered callback
   */
  removeNavigateListener: (callback) => {
    ipcRenderer.removeListener("navigate", callback);
  },
});

/**
 * Expose platform information to the renderer.
 *
 * This helps the Vue.js app adjust UI/behavior based on the operating system.
 */
contextBridge.exposeInMainWorld("platform", {
  /**
   * Get the current operating system.
   *
   * @returns {string} 'win32', 'darwin', or 'linux'
   */
  os: process.platform,

  /**
   * Check if running on Windows.
   *
   * @returns {boolean}
   */
  isWindows: process.platform === "win32",

  /**
   * Check if running on macOS.
   *
   * @returns {boolean}
   */
  isMac: process.platform === "darwin",

  /**
   * Check if running on Linux.
   *
   * @returns {boolean}
   */
  isLinux: process.platform === "linux",
});

console.log("Preload script loaded - Context bridge established");
