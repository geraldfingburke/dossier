import { reactive, inject } from 'vue'
import { GraphQLClient } from 'graphql-request'

const API_URL = 'http://localhost:8080/graphql'

export function createStore() {
  const state = reactive({
    user: null,
    token: localStorage.getItem('token') || null,
    feeds: [],
    articles: [],
    digests: [],
    loading: false,
    error: null
  })

  const client = new GraphQLClient(API_URL, {
    headers: () => ({
      Authorization: state.token ? `Bearer ${state.token}` : ''
    })
  })

  async function login(email, password) {
    try {
      state.loading = true
      state.error = null
      
      const mutation = `
        mutation Login($email: String!, $password: String!) {
          login(email: $email, password: $password) {
            token
            user {
              id
              email
              name
            }
          }
        }
      `
      
      const data = await client.request(mutation, { email, password })
      state.token = data.login.token
      state.user = data.login.user
      localStorage.setItem('token', data.login.token)
    } catch (error) {
      state.error = error.message
      throw error
    } finally {
      state.loading = false
    }
  }

  async function register(email, password, name) {
    try {
      state.loading = true
      state.error = null
      
      const mutation = `
        mutation Register($email: String!, $password: String!, $name: String!) {
          register(email: $email, password: $password, name: $name) {
            token
            user {
              id
              email
              name
            }
          }
        }
      `
      
      const data = await client.request(mutation, { email, password, name })
      state.token = data.register.token
      state.user = data.register.user
      localStorage.setItem('token', data.register.token)
    } catch (error) {
      state.error = error.message
      throw error
    } finally {
      state.loading = false
    }
  }

  function logout() {
    state.user = null
    state.token = null
    state.feeds = []
    state.articles = []
    state.digests = []
    localStorage.removeItem('token')
  }

  async function fetchFeeds() {
    try {
      state.loading = true
      const query = `
        query {
          feeds {
            id
            url
            title
            description
            active
            createdAt
            updatedAt
          }
        }
      `
      
      const data = await client.request(query)
      state.feeds = data.feeds
    } catch (error) {
      state.error = error.message
    } finally {
      state.loading = false
    }
  }

  async function addFeed(url) {
    try {
      state.loading = true
      const mutation = `
        mutation AddFeed($url: String!) {
          addFeed(url: $url) {
            id
            url
            title
            description
            active
            createdAt
            updatedAt
          }
        }
      `
      
      const data = await client.request(mutation, { url })
      state.feeds.unshift(data.addFeed)
      return data.addFeed
    } catch (error) {
      state.error = error.message
      throw error
    } finally {
      state.loading = false
    }
  }

  async function deleteFeed(id) {
    try {
      state.loading = true
      const mutation = `
        mutation DeleteFeed($id: ID!) {
          deleteFeed(id: $id)
        }
      `
      
      await client.request(mutation, { id })
      state.feeds = state.feeds.filter(f => f.id !== id)
    } catch (error) {
      state.error = error.message
      throw error
    } finally {
      state.loading = false
    }
  }

  async function refreshAllFeeds() {
    try {
      state.loading = true
      const mutation = `
        mutation {
          refreshAllFeeds
        }
      `
      
      await client.request(mutation)
      // Refresh the feeds list after a short delay
      setTimeout(() => fetchFeeds(), 2000)
    } catch (error) {
      state.error = error.message
      throw error
    } finally {
      state.loading = false
    }
  }

  async function fetchArticles(limit = 50, offset = 0) {
    try {
      state.loading = true
      const query = `
        query Articles($limit: Int, $offset: Int) {
          articles(limit: $limit, offset: $offset) {
            id
            feedId
            title
            link
            description
            content
            author
            publishedAt
            createdAt
          }
        }
      `
      
      const data = await client.request(query, { limit, offset })
      state.articles = data.articles
    } catch (error) {
      state.error = error.message
    } finally {
      state.loading = false
    }
  }

  async function fetchDigests(limit = 10) {
    try {
      state.loading = true
      const query = `
        query Digests($limit: Int) {
          digests(limit: $limit) {
            id
            userId
            date
            summary
            createdAt
            articles {
              id
              title
              link
              publishedAt
            }
          }
        }
      `
      
      const data = await client.request(query, { limit })
      state.digests = data.digests
    } catch (error) {
      state.error = error.message
    } finally {
      state.loading = false
    }
  }

  async function generateDigest() {
    try {
      state.loading = true
      const mutation = `
        mutation {
          generateDigest {
            id
            userId
            date
            summary
            createdAt
            articles {
              id
              title
              link
              publishedAt
            }
          }
        }
      `
      
      const data = await client.request(mutation)
      state.digests.unshift(data.generateDigest)
      return data.generateDigest
    } catch (error) {
      state.error = error.message
      throw error
    } finally {
      state.loading = false
    }
  }

  // Initialize user if token exists
  if (state.token) {
    const query = `
      query {
        me {
          id
          email
          name
        }
      }
    `
    client.request(query)
      .then(data => {
        state.user = data.me
      })
      .catch(() => {
        logout()
      })
  }

  return {
    state,
    login,
    register,
    logout,
    fetchFeeds,
    addFeed,
    deleteFeed,
    refreshAllFeeds,
    fetchArticles,
    fetchDigests,
    generateDigest
  }
}

export function useStore() {
  const store = inject('store')
  return store
}
