const {
  app,
  BrowserWindow,
  ipcMain,
  Menu,
  Tray,
  dialog,
  shell,
} = require("electron");
const { autoUpdater } = require("electron-updater");
const path = require("path");
const fs = require("fs");

// ============================================================================
// APPLICATION STATE
// ============================================================================

let mainWindow = null;
let tray = null;
let serviceProcess = null;
const isDev = process.argv.includes("--dev");
const servicePort = 8080;
const serviceBaseUrl = `http://localhost:${servicePort}`;

// ============================================================================
// WINDOW MANAGEMENT
// ============================================================================

/**
 * Creates the main application window.
 *
 * The window loads the Vue.js client either from:
 * - Development: Vite dev server (http://localhost:5173)
 * - Production: Bundled static files from resources/client
 */
function createWindow() {
  const windowIconPath = path.join(__dirname, "../build/icon.png");
  const windowConfig = {
    width: 1400,
    height: 900,
    minWidth: 800,
    minHeight: 600,
    title: "Dossier",
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      preload: path.join(__dirname, "preload.js"),
      webSecurity: true,
      allowRunningInsecureContent: false,
    },
    backgroundColor: "#1e1e1e",
    show: false, // Show after ready-to-show for smoother experience
  };

  // Add icon if it exists
  if (fs.existsSync(windowIconPath)) {
    windowConfig.icon = windowIconPath;
  }

  mainWindow = new BrowserWindow(windowConfig);

  // Load the appropriate URL based on environment
  const startUrl = isDev
    ? "http://localhost:5173" // Vite dev server
    : `file://${path.join(process.resourcesPath, "client/index.html")}`;

  console.log("Loading URL:", startUrl);

  mainWindow.loadURL(startUrl).catch((err) => {
    console.error("Failed to load URL:", err);
  });

  // Log when loading starts
  mainWindow.webContents.on("did-start-loading", () => {
    console.log("Started loading content...");
  });

  // Log when loading finishes
  mainWindow.webContents.on("did-finish-load", () => {
    console.log("Finished loading content");
  });

  // Log any loading failures
  mainWindow.webContents.on("did-fail-load", (event, errorCode, errorDescription) => {
    console.error("Failed to load:", errorCode, errorDescription);
  });

  // Show window when ready to prevent flickering
  mainWindow.once("ready-to-show", () => {
    console.log("Window ready to show");
    mainWindow.show();

    // Open DevTools in development mode
    if (isDev) {
      mainWindow.webContents.openDevTools();
    }
  });

  // Handle window close - minimize to tray instead of quitting
  mainWindow.on("close", (event) => {
    if (!app.isQuitting && tray) {
      event.preventDefault();
      mainWindow.hide();
    }
    return false;
  });

  // Clean up on window closed
  mainWindow.on("closed", () => {
    mainWindow = null;
  });

  // Handle external links - open in default browser
  mainWindow.webContents.setWindowOpenHandler(({ url }) => {
    shell.openExternal(url);
    return { action: "deny" };
  });

  // Create application menu
  createMenu();
}

/**
 * Creates the system tray icon and menu.
 *
 * Provides quick access to:
 * - Show/hide main window
 * - Check for updates
 * - Quit application
 */
function createTray() {
  const trayIconPath = path.join(__dirname, "../build/icon.png");

  // Skip tray creation if icon doesn't exist
  if (!fs.existsSync(trayIconPath)) {
    console.log("Tray icon not found, skipping tray creation");
    return;
  }

  tray = new Tray(trayIconPath);

  const contextMenu = Menu.buildFromTemplate([
    {
      label: "Show Dossier",
      click: () => {
        if (mainWindow) {
          mainWindow.show();
          mainWindow.focus();
        } else {
          createWindow();
        }
      },
    },
    { type: "separator" },
    {
      label: "Check for Updates",
      click: () => {
        autoUpdater.checkForUpdatesAndNotify();
      },
    },
    { type: "separator" },
    {
      label: "Quit",
      click: () => {
        app.isQuitting = true;
        app.quit();
      },
    },
  ]);

  tray.setToolTip("Dossier - News Digest");
  tray.setContextMenu(contextMenu);

  // Double-click to show window
  tray.on("double-click", () => {
    if (mainWindow) {
      mainWindow.show();
      mainWindow.focus();
    } else {
      createWindow();
    }
  });
}

/**
 * Creates the application menu bar.
 */
function createMenu() {
  const template = [
    {
      label: "File",
      submenu: [
        {
          label: "Preferences",
          accelerator: "CmdOrCtrl+,",
          click: () => {
            // Navigate to settings view
            mainWindow?.webContents.send("navigate", "/settings");
          },
        },
        { type: "separator" },
        {
          label: "Exit",
          accelerator: "CmdOrCtrl+Q",
          click: () => {
            app.isQuitting = true;
            app.quit();
          },
        },
      ],
    },
    {
      label: "Edit",
      submenu: [
        { role: "undo" },
        { role: "redo" },
        { type: "separator" },
        { role: "cut" },
        { role: "copy" },
        { role: "paste" },
        { role: "selectAll" },
      ],
    },
    {
      label: "View",
      submenu: [
        { role: "reload" },
        { role: "forceReload" },
        { type: "separator" },
        { role: "resetZoom" },
        { role: "zoomIn" },
        { role: "zoomOut" },
        { type: "separator" },
        { role: "togglefullscreen" },
      ],
    },
    {
      label: "Help",
      submenu: [
        {
          label: "Documentation",
          click: async () => {
            await shell.openExternal(
              "https://github.com/geraldfingburke/dossier"
            );
          },
        },
        {
          label: "Report Issue",
          click: async () => {
            await shell.openExternal(
              "https://github.com/geraldfingburke/dossier/issues"
            );
          },
        },
        { type: "separator" },
        {
          label: "About",
          click: () => {
            dialog.showMessageBox(mainWindow, {
              type: "info",
              title: "About Dossier",
              message: "Dossier Desktop",
              detail: `Version: ${app.getVersion()}\n\nAI-powered news digest application\n\n© 2025 Gerald Fingburke`,
              buttons: ["OK"],
            });
          },
        },
      ],
    },
  ];

  // Add DevTools menu in development
  if (isDev) {
    template.push({
      label: "Developer",
      submenu: [
        { role: "toggleDevTools" },
        { type: "separator" },
        {
          label: "Open App Data",
          click: () => {
            shell.openPath(app.getPath("userData"));
          },
        },
      ],
    });
  }

  const menu = Menu.buildFromTemplate(template);
  Menu.setApplicationMenu(menu);
}

// ============================================================================
// SERVICE MANAGEMENT
// ============================================================================

/**
 * Checks if the Go backend service is running.
 *
 * @returns {Promise<boolean>} True if service is accessible
 */
async function checkServiceHealth() {
  try {
    const http = require("http");

    return new Promise((resolve) => {
      const req = http.get(`${serviceBaseUrl}/health`, (res) => {
        resolve(res.statusCode === 200);
      });

      req.on("error", () => {
        resolve(false);
      });

      req.setTimeout(2000, () => {
        req.destroy();
        resolve(false);
      });
    });
  } catch (error) {
    return false;
  }
}

/**
 * Starts the Go backend service.
 *
 * In development: Assumes service is already running via docker-compose
 * In production: Will launch bundled Go binary (future implementation)
 */
async function startService() {
  if (isDev) {
    console.log("Development mode: Expecting service at", serviceBaseUrl);

    // Check if service is running
    const isHealthy = await checkServiceHealth();
    if (!isHealthy) {
      dialog.showErrorBox(
        "Service Not Running",
        "The Dossier backend service is not running.\n\n" +
          "Please start it with: docker-compose up\n\n" +
          "The application will continue, but features will be unavailable."
      );
    }
    return;
  }

  // TODO: Production - launch bundled Go service binary
  // This will be implemented when creating the desktop-service package
  console.log("Production mode: Service management not yet implemented");
}

/**
 * Stops the Go backend service.
 */
function stopService() {
  if (serviceProcess) {
    serviceProcess.kill();
    serviceProcess = null;
  }
}

// ============================================================================
// IPC HANDLERS
// ============================================================================

/**
 * Set up IPC handlers for communication between renderer and main process.
 */
function setupIpcHandlers() {
  // Get service URL for frontend
  ipcMain.handle("get-service-url", () => {
    return serviceBaseUrl;
  });

  // Check service health
  ipcMain.handle("check-service-health", async () => {
    return await checkServiceHealth();
  });

  // Open external URL
  ipcMain.handle("open-external", async (event, url) => {
    await shell.openExternal(url);
  });

  // Get app version
  ipcMain.handle("get-app-version", () => {
    return app.getVersion();
  });

  // Minimize to tray
  ipcMain.handle("minimize-to-tray", () => {
    if (mainWindow) {
      mainWindow.hide();
    }
  });

  // Show notification
  ipcMain.handle("show-notification", (event, { title, body }) => {
    const { Notification } = require("electron");
    new Notification({ title, body }).show();
  });
}

// ============================================================================
// AUTO-UPDATER
// ============================================================================

/**
 * Configure auto-updater for checking and installing updates.
 */
function setupAutoUpdater() {
  // Disable auto-download in development
  autoUpdater.autoDownload = !isDev;
  autoUpdater.autoInstallOnAppQuit = true;

  autoUpdater.on("checking-for-update", () => {
    console.log("Checking for updates...");
  });

  autoUpdater.on("update-available", (info) => {
    console.log("Update available:", info.version);

    dialog.showMessageBox(mainWindow, {
      type: "info",
      title: "Update Available",
      message: `A new version (${info.version}) is available!`,
      detail: "The update will be downloaded in the background.",
      buttons: ["OK"],
    });
  });

  autoUpdater.on("update-not-available", () => {
    console.log("No updates available");
  });

  autoUpdater.on("error", (err) => {
    console.error("Update error:", err);
  });

  autoUpdater.on("download-progress", (progressObj) => {
    console.log(
      `Download speed: ${progressObj.bytesPerSecond} - Downloaded ${progressObj.percent}%`
    );
  });

  autoUpdater.on("update-downloaded", (info) => {
    console.log("Update downloaded:", info.version);

    dialog
      .showMessageBox(mainWindow, {
        type: "info",
        title: "Update Ready",
        message: "A new version has been downloaded.",
        detail: "The application will restart to install the update.",
        buttons: ["Restart Now", "Later"],
      })
      .then((result) => {
        if (result.response === 0) {
          autoUpdater.quitAndInstall();
        }
      });
  });

  // Check for updates on startup (after 3 seconds delay)
  if (!isDev) {
    setTimeout(() => {
      autoUpdater.checkForUpdatesAndNotify();
    }, 3000);
  }
}

// ============================================================================
// APPLICATION LIFECYCLE
// ============================================================================

// Disable GPU acceleration to prevent crashes on some Windows systems
app.disableHardwareAcceleration();

/**
 * Initialize the application when Electron is ready.
 */
app.whenReady().then(async () => {
  console.log("Dossier Desktop starting...");
  console.log("App path:", app.getAppPath());
  console.log("User data:", app.getPath("userData"));
  console.log("Development mode:", isDev);

  // Start backend service
  await startService();

  // Set up IPC handlers
  setupIpcHandlers();

  // Create window and tray
  createWindow();
  createTray();

  // Set up auto-updater
  setupAutoUpdater();

  // macOS: Re-create window when dock icon is clicked
  app.on("activate", () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    } else if (mainWindow) {
      mainWindow.show();
    }
  });
});

/**
 * Quit when all windows are closed (except on macOS).
 */
app.on("window-all-closed", () => {
  if (process.platform !== "darwin") {
    app.quit();
  }
});

/**
 * Clean up before quitting.
 */
app.on("before-quit", () => {
  app.isQuitting = true;
});

/**
 * Stop service on quit.
 */
app.on("quit", () => {
  stopService();
});

/**
 * Handle uncaught exceptions.
 */
process.on("uncaughtException", (error) => {
  console.error("Uncaught exception:", error);

  dialog.showErrorBox(
    "Application Error",
    `An unexpected error occurred:\n\n${error.message}\n\nThe application will continue running.`
  );
});
