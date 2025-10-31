<template>
  <div class="feeds-view">
    <div class="header">
      <h2>RSS Feeds</h2>
      <button @click="handleRefreshAll" class="primary" :disabled="store.state.loading">
        {{ store.state.loading ? 'Refreshing...' : 'Refresh All Feeds' }}
      </button>
    </div>

    <div v-if="store.state.error" class="error">{{ store.state.error }}</div>

    <div class="add-feed card">
      <h3>Add New Feed</h3>
      <form @submit.prevent="handleAddFeed">
        <input 
          v-model="newFeedUrl" 
          type="url" 
          placeholder="https://example.com/feed.xml"
          required
        />
        <button type="submit" class="primary" :disabled="store.state.loading">
          Add Feed
        </button>
      </form>
      <div v-if="addFeedSuccess" class="success">Feed added successfully!</div>
    </div>

    <div class="feeds-list">
      <div v-if="store.state.feeds.length === 0" class="card">
        <p>No feeds yet. Add your first RSS feed above!</p>
      </div>
      
      <div v-for="feed in store.state.feeds" :key="feed.id" class="card feed-item">
        <div class="feed-header">
          <div>
            <h3>{{ feed.title || feed.url }}</h3>
            <p class="feed-url">{{ feed.url }}</p>
            <p v-if="feed.description" class="feed-description">{{ feed.description }}</p>
          </div>
          <div class="feed-actions">
            <button @click="handleDeleteFeed(feed.id)" class="delete">
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useStore } from '../store'

export default {
  name: 'FeedsView',
  setup() {
    const store = useStore()
    const newFeedUrl = ref('')
    const addFeedSuccess = ref(false)

    onMounted(async () => {
      await store.fetchFeeds()
    })

    const handleAddFeed = async () => {
      try {
        await store.addFeed(newFeedUrl.value)
        newFeedUrl.value = ''
        addFeedSuccess.value = true
        setTimeout(() => {
          addFeedSuccess.value = false
        }, 3000)
      } catch (error) {
        console.error('Error adding feed:', error)
      }
    }

    const handleDeleteFeed = async (id) => {
      if (confirm('Are you sure you want to delete this feed?')) {
        try {
          await store.deleteFeed(id)
        } catch (error) {
          console.error('Error deleting feed:', error)
        }
      }
    }

    const handleRefreshAll = async () => {
      try {
        await store.refreshAllFeeds()
      } catch (error) {
        console.error('Error refreshing feeds:', error)
      }
    }

    return {
      store,
      newFeedUrl,
      addFeedSuccess,
      handleAddFeed,
      handleDeleteFeed,
      handleRefreshAll
    }
  }
}
</script>

<style scoped>
.feeds-view {
  padding: 1rem 0;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.add-feed h3 {
  margin-bottom: 1rem;
}

.add-feed form {
  display: flex;
  gap: 1rem;
}

.add-feed input {
  flex: 1;
}

.feeds-list {
  margin-top: 2rem;
}

.feed-item {
  margin-bottom: 1rem;
}

.feed-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.feed-item h3 {
  margin-bottom: 0.5rem;
  color: #2c3e50;
}

.feed-url {
  color: #666;
  font-size: 0.9rem;
  margin-bottom: 0.5rem;
}

.feed-description {
  color: #555;
  font-size: 0.95rem;
}

.delete {
  background: #f44336;
  color: white;
}

.delete:hover {
  background: #d32f2f;
}
</style>
