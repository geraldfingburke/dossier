<template>
  <div class="digests-view">
    <div class="header">
      <h2>AI Digests</h2>
      <button @click="handleGenerateDigest" class="primary" :disabled="store.state.loading">
        {{ store.state.loading ? 'Generating...' : 'Generate New Digest' }}
      </button>
    </div>

    <div v-if="store.state.error" class="error">{{ store.state.error }}</div>
    
    <div v-if="generateSuccess" class="success">Digest generated successfully!</div>

    <div v-if="store.state.loading && store.state.digests.length === 0" class="card">
      <p>Loading digests...</p>
    </div>

    <div v-if="store.state.digests.length === 0 && !store.state.loading" class="card">
      <p>No digests yet. Generate your first AI-powered digest by clicking the button above!</p>
    </div>

    <div class="digests-list">
      <div v-for="digest in store.state.digests" :key="digest.id" class="card digest-item">
        <div class="digest-header">
          <h3>📅 {{ formatDate(digest.date) }}</h3>
          <span class="article-count">{{ digest.articles.length }} articles</span>
        </div>
        
        <div class="digest-summary">
          <pre>{{ digest.summary }}</pre>
        </div>

        <div v-if="digest.articles.length > 0" class="digest-articles">
          <h4>Articles in this digest:</h4>
          <ul>
            <li v-for="article in digest.articles" :key="article.id">
              <a :href="article.link" target="_blank" rel="noopener noreferrer">
                {{ article.title }}
              </a>
              <span class="article-date">{{ formatTime(article.publishedAt) }}</span>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useStore } from '../store'

export default {
  name: 'DigestsView',
  setup() {
    const store = useStore()
    const generateSuccess = ref(false)

    onMounted(async () => {
      await store.fetchDigests()
    })

    const handleGenerateDigest = async () => {
      try {
        await store.generateDigest()
        generateSuccess.value = true
        setTimeout(() => {
          generateSuccess.value = false
        }, 3000)
      } catch (error) {
        console.error('Error generating digest:', error)
      }
    }

    const formatDate = (dateString) => {
      const date = new Date(dateString)
      return date.toLocaleDateString('en-US', { 
        year: 'numeric', 
        month: 'long', 
        day: 'numeric'
      })
    }

    const formatTime = (dateString) => {
      const date = new Date(dateString)
      return date.toLocaleString('en-US', { 
        month: 'short', 
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    return {
      store,
      generateSuccess,
      handleGenerateDigest,
      formatDate,
      formatTime
    }
  }
}
</script>

<style scoped>
.digests-view {
  padding: 1rem 0;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.digests-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.digest-item {
  padding: 2rem;
}

.digest-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  padding-bottom: 1rem;
  border-bottom: 2px solid #e0e0e0;
}

.digest-header h3 {
  color: #2c3e50;
}

.article-count {
  background: #4CAF50;
  color: white;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.85rem;
}

.digest-summary {
  margin-bottom: 1.5rem;
  background: #f9f9f9;
  padding: 1.5rem;
  border-radius: 4px;
  border-left: 4px solid #4CAF50;
}

.digest-summary pre {
  white-space: pre-wrap;
  font-family: inherit;
  line-height: 1.6;
  color: #333;
}

.digest-articles {
  margin-top: 1.5rem;
}

.digest-articles h4 {
  margin-bottom: 0.75rem;
  color: #555;
}

.digest-articles ul {
  list-style: none;
  padding: 0;
}

.digest-articles li {
  padding: 0.5rem 0;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.digest-articles li:last-child {
  border-bottom: none;
}

.digest-articles a {
  color: #2c3e50;
  text-decoration: none;
  flex: 1;
}

.digest-articles a:hover {
  color: #4CAF50;
}

.article-date {
  font-size: 0.85rem;
  color: #666;
  margin-left: 1rem;
}
</style>
