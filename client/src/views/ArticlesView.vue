<template>
  <div class="articles-view">
    <h2>Recent Articles</h2>

    <div v-if="store.state.loading" class="card">
      <p>Loading articles...</p>
    </div>

    <div v-if="store.state.articles.length === 0 && !store.state.loading" class="card">
      <p>No articles yet. Add some feeds and refresh them to see articles here!</p>
    </div>

    <div class="articles-list">
      <div v-for="article in store.state.articles" :key="article.id" class="card article-item">
        <h3>
          <a :href="article.link" target="_blank" rel="noopener noreferrer">
            {{ article.title }}
          </a>
        </h3>
        <div class="article-meta">
          <span v-if="article.author" class="author">{{ article.author }}</span>
          <span class="date">{{ formatDate(article.publishedAt) }}</span>
        </div>
        <p v-if="article.description" class="description">
          {{ article.description }}
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import { onMounted } from 'vue'
import { useStore } from '../store'

export default {
  name: 'ArticlesView',
  setup() {
    const store = useStore()

    onMounted(async () => {
      await store.fetchArticles()
    })

    const formatDate = (dateString) => {
      const date = new Date(dateString)
      return date.toLocaleDateString('en-US', { 
        year: 'numeric', 
        month: 'short', 
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    return {
      store,
      formatDate
    }
  }
}
</script>

<style scoped>
.articles-view {
  padding: 1rem 0;
}

h2 {
  margin-bottom: 2rem;
}

.articles-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.article-item h3 {
  margin-bottom: 0.5rem;
}

.article-item h3 a {
  color: #2c3e50;
  text-decoration: none;
}

.article-item h3 a:hover {
  color: #4CAF50;
  text-decoration: underline;
}

.article-meta {
  display: flex;
  gap: 1rem;
  margin-bottom: 0.5rem;
  font-size: 0.85rem;
  color: #666;
}

.author {
  font-weight: 500;
}

.description {
  color: #555;
  line-height: 1.5;
}
</style>
