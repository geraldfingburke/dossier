import { createRouter, createWebHistory } from "vue-router";
import DossierConfigsView from "../views/DossierConfigsView.vue";
import DossierDetailView from "../views/DossierDetailView.vue";

const routes = [
  {
    path: "/",
    name: "DossierList",
    component: DossierConfigsView,
    meta: {
      title: "Dossiers",
      showAsList: true,
    },
  },
  {
    path: "/dossier/new",
    name: "DossierCreate",
    component: DossierDetailView,
    props: { dossierId: "new" },
    meta: {
      title: "Create New Dossier",
    },
  },
  {
    path: "/dossier/:dossierId",
    name: "DossierDetail",
    component: DossierDetailView,
    props: true,
    meta: {
      title: "Edit Dossier",
    },
  },
  // Redirect old paths for backward compatibility
  {
    path: "/dossiers",
    redirect: "/",
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Update document title based on route meta
router.beforeEach((to, from, next) => {
  if (to.meta.title) {
    document.title = `${to.meta.title} - Dossier`;
  }
  next();
});

export default router;
